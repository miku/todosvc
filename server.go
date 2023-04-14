package todosvc

import (
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//go:embed views/*
var views embed.FS

var ErrInvalid = errors.New("invalid")

// Todo item.
type Todo struct {
	Id      int
	Title   string
	Done    bool
	Created time.Time
}

func (t *Todo) Finish() error {
	if time.Now().After(time.Now().Add(-240 * time.Hour)) {
		t.Done = true
		return nil
	}
	return ErrInvalid
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
		tmpl := mustTemplate("views/index.html")
		tmpl.Execute(w, todos)
	}
}

func (s *Server) handleListing() http.HandlerFunc {
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
		enc := json.NewEncoder(w)
		if err := enc.Encode(todos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := mustTemplate("views/new.html")
			tmpl.Execute(w, map[string]interface{}{csrf.TemplateTag: csrf.TemplateField(r)})
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
		err = todo.Finish()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// TODO: addition business logic, e.g. todo only if older then 10 days
		stmt := `update todo set done = 1 where id = ?`
		if _, err := s.DB.Exec(stmt, vars["id"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", 302)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func mustTemplate(path string) *template.Template {
	b, err := views.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return template.Must(template.New("t").Parse(string(b)))
}
