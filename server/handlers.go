package server

import (
	"fmt"
	"net/http"
	"text/template"

	"log"

	"github.com/gorilla/mux"
	"github.com/srumut/todofullstack/storage"
)

func (s *Server) HandleIndexPage(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	t := template.Must(template.ParseFiles("./views/index.html"))

	todos, err := s.store.GetTodos()
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(w, "index", todos); err != nil {
		log.Printf("handleIndexPage | error executing template | err: %s", err)
		return err
	}
	return nil
}

func (s *Server) HandleTodoAdd(w http.ResponseWriter, r *http.Request) error {
	nextID := s.store.GetNextID()
	// add new todo to database
	title := r.FormValue("title")
	description := r.FormValue("description")

	// check values
	if err := s.checkInputs(title, description); err != nil {
		log.Printf("HandleTodoAdd | error checking inputs | err: %s", err)
		return err
	}

	if err := s.store.AddTodo(nextID, title, description); err != nil {
		return err
	}

	// display new todo in index page
	t := template.Must(template.ParseFiles("./views/index.html"))
	todo := storage.NewTodo(nextID, title, description)
	if err := t.ExecuteTemplate(w, "todo", todo); err != nil {
		log.Printf("HandleTodoAdd | error executing template | err: %s", err)
		return err
	}

	return nil
}

func (s *Server) HandleTodoDel(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	if err := s.store.RemoveByID(idStr); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

type serverFunc func(http.ResponseWriter, *http.Request) error

func Make(f serverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) checkInputs(title, description string) error {
	if title == "" || description == "" {
		return fmt.Errorf("title and description are required")
	}

	// check if todo with same title already exists
	todos, err := s.store.GetTodos()
	if err != nil {
		return err
	}
	for _, todo := range todos {
		if todo.Title == title {
			return fmt.Errorf("todo with title %s already exists", title)
		}
	}

	return nil
}
