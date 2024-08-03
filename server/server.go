package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/srumut/todofullstack/storage"
)

type Server struct {
	store storage.Storage
	fs    http.Handler
}

func NewServer(store storage.Storage) *Server {
	return &Server{
		store,
		http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))),
	}
}

func (s *Server) initialize() *http.Server {
	mux := mux.NewRouter()

	mux.HandleFunc("/", Make(s.HandleIndexPage)).Methods("GET")
	mux.HandleFunc("/todo/add", Make(s.HandleTodoAdd)).Methods("POST")
	mux.HandleFunc("/todo/del/{id:[0-9]+}", Make(s.HandleTodoDel)).Methods("DELETE")
	mux.PathPrefix("/static/").HandlerFunc(s.handlerForStatic)

	server := &http.Server{
		Addr:    getAddr(),
		Handler: mux,
	}

	return server
}

func (s *Server) Start() error {
	server := s.initialize()
	fmt.Println("Server is up and listening on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func getAddr() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading.env file")
	}
	return os.Getenv("LISTEN_ADDR")
}

func (s *Server) handlerForStatic(w http.ResponseWriter, r *http.Request) {
	if filepath.Ext(r.URL.Path) == ".css" {
		w.Header().Set("content-type", "text/css")
	}
	s.fs.ServeHTTP(w, r)
}
