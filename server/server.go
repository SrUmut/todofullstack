package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/srumut/todofullstack/storage"
)

type Server struct {
	store storage.Storage
	fs    http.Handler
}

func NewServer(store storage.Storage) *Server {
	return &Server{
		store,
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	}
}

func (s *Server) initialize() *http.Server {
	mux := mux.NewRouter()

	mux.HandleFunc("/", Make(s.HandleIndexPage)).Methods("GET")

	mux.HandleFunc("/register", Make(s.HandleRegisterPage)).Methods("GET")
	mux.HandleFunc("/register", Make(s.HandleRegisterRequest)).Methods("POST")
	mux.HandleFunc("/login", Make(s.HandleLoginPage)).Methods("GET")
	mux.HandleFunc("/login", Make(s.HandleLoginRequest)).Methods("POST")

	mux.HandleFunc("/todo/add", Make(s.HandleTodoAdd)).Methods("POST")
	mux.HandleFunc("/todo/del/{id:[0-9]+}", Make(s.HandleTodoDel)).Methods("DELETE")
	mux.PathPrefix("/static/").HandlerFunc(s.HandleStaticFiles)

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
	return os.Getenv("LISTEN_ADDR")
}
