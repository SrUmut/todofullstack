package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/srumut/todofullstack/storage"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) HandleIndexPage(w http.ResponseWriter, r *http.Request) error {
	claims, err := s.handleJWT(w, r)
	if err != nil {
		fmt.Printf("HandleIndexPage | error handling jwt | err: %s\n", err)
		return nil
	}
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	t := template.Must(template.ParseFiles("./views/index.html"))

	todos, err := s.store.GetTodos(claims["sub"].(string))
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(w, "index", todos); err != nil {
		log.Printf("handleIndexPage | error executing template | err: %s", err)
		return err
	}

	return nil
}

// Adds a new todo to the database, and responses with the new todo
func (s *Server) HandleTodoAdd(w http.ResponseWriter, r *http.Request) error {
	claims, err := s.handleJWT(w, r)
	if err != nil {
		fmt.Printf("HandleTodoAdd | error handling jwt | err: %s\n", err)
		return nil
	}

	username := claims["sub"].(string)

	// add new todo to database
	title := r.FormValue("title")
	description := r.FormValue("description")

	// check values
	if err := s.checkInputs(title, description, claims); err != nil {
		log.Printf("HandleTodoAdd | error checking inputs | err: %s", err)
		return err
	}

	todoid, err := s.store.AddTodo(username, title, description)
	if err != nil {
		return err
	}

	// display new todo in index page
	t := template.Must(template.ParseFiles("./views/index.html"))
	todo := storage.NewTodo(int(todoid), title, description)
	if err := t.ExecuteTemplate(w, "todo", todo); err != nil {
		log.Printf("HandleTodoAdd | error executing template | err: %s", err)
		return err
	}

	return nil
}

// Deletes a todo with id in the url from the database
func (s *Server) HandleTodoDel(w http.ResponseWriter, r *http.Request) error {
	claims, err := s.handleJWT(w, r)
	if err != nil {
		fmt.Printf("HandleTodoDel | error handling jwt | err: %s\n", err)
		return nil
	}
	idStr := mux.Vars(r)["id"]
	if err := s.store.RemoveTodo(claims["sub"].(string), idStr); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

type serverFunc func(http.ResponseWriter, *http.Request) error

// Make wraps a serverFunc into a http.HandlerFunc
func Make(f serverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("error in %s: %s\n", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Check if any of the inputs are empty or if a todo with the same title already exists
func (s *Server) checkInputs(title, description string, claims jwt.MapClaims) error {
	if title == "" || description == "" {
		return fmt.Errorf("title and description are required")
	}

	// check if todo with same title already exists
	todos, err := s.store.GetTodos(claims["sub"].(string))
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

func (s *Server) HandleStaticFiles(w http.ResponseWriter, r *http.Request) {
	if filepath.Ext(r.URL.Path) == ".css" {
		w.Header().Set("content-type", "text/css")
	} else if filepath.Ext(r.URL.Path) == ".js" {
		w.Header().Set("content-type", "text/javascript")
	}
	s.fs.ServeHTTP(w, r)
}

func (s *Server) HandleRegisterPage(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	t := template.Must(template.ParseFiles("./views/register.html"))

	if err := t.ExecuteTemplate(w, "register", nil); err != nil {
		log.Printf("HandleRegisterPage | error executing template | err: %s", err)
		return err
	}

	return nil
}

func (s *Server) HandleRegisterRequest(w http.ResponseWriter, r *http.Request) error {
	username, password, password_repeated := r.FormValue("username"), r.FormValue("password"), r.FormValue("password-repeated")

	if err := checkRegisterReq(username, password, password_repeated); err != nil {
		w.Header().Set("Content-Type", "text/html charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return nil
	}

	enc_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.store.AddUser(username, string(enc_password)); err != nil {
		w.Header().Set("Content-Type", "text/html charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return nil
	}

	w.Header().Set("HX-Redirect", "/login")

	return nil
}

func (s *Server) HandleLoginPage(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	t := template.Must(template.ParseFiles("./views/login.html"))

	if err := t.ExecuteTemplate(w, "login", nil); err != nil {
		log.Printf("HandleLoginPage | error executing template | err: %s", err)
		return err
	}

	return nil
}

func (s *Server) HandleLoginRequest(w http.ResponseWriter, r *http.Request) error {
	username, password := r.FormValue("username"), r.FormValue("password")

	enc_pass, err := s.store.GetPassword(username)
	if err != nil {
		s.wrongCredentials(w, err)
		return nil
	}

	// if no user with given username
	if enc_pass == "" {
		s.wrongCredentials(w, err)
		return nil
	}

	// if pasword is not correct
	if err := bcrypt.CompareHashAndPassword([]byte(enc_pass), []byte(password)); err != nil {
		s.wrongCredentials(w, err)
		return nil
	}

	expeDur := time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(expeDur).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return err
	}

	w.Header().Set("jwt-token", tokenStr)
	w.Header().Set("jwt-exp", strconv.FormatInt(expeDur.Milliseconds(), 10))

	return nil
}

func (s *Server) handleJWT(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, error) {
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}
	tokenString := strings.TrimPrefix(cookie.String(), "jwt_token=")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || time.Now().Compare(time.Unix(int64(claims["exp"].(float64)), 0)) >= 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, err
	}

	return claims, nil
}

func (s *Server) wrongCredentials(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("wrong username or password"))
	log.Printf("HandleLoginRequest | error getting password | err: %s", err)
}
