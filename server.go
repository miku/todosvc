package todosvc

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Todo item.
type Todo struct {
	Id      int
	Title   string
	Done    bool
	Created time.Time
}

type Server struct {
	DB     *sqlx.DB
	Router *mux.Router
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		err := s.DB.Select(&todos, "SELECT * FROM todo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v := range todos {
			log.Printf("%v", v)
		}
		tmpl := template.Must(template.ParseFiles("./views/index.html"))
		tmpl.Execute(w, todos)
	}
}

func (s *Server) handleNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("./views/new.html"))
			tmpl.Execute(w, nil)
			return
		}
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			title := r.FormValue("title")
			stmt := `insert into todo (title, done, created) values (?, ?, ?)`
			if _, err := s.DB.Exec(stmt, title, false, time.Now()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", 302)
		}
	}
}

func (s *Server) handleDone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var todo Todo
		err := s.DB.Get(&todo, "SELECT * FROM todo where id = ?", vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(todo)
		stmt := `update todo set done = 1 where id = ?`
		if _, err := s.DB.Exec(stmt, vars["id"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", 302)
	}
}

func (s *Server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}
