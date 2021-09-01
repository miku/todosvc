package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
)

// Todo item.
type Todo struct {
	Id      int
	Title   string
	Done    bool
	Created time.Time
}

type Server struct {
	db     *sqlx.DB
	router *mux.Router
}

func (s *Server) routes() {
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/new", s.handleNew()).Methods("GET", "POST")
	s.router.HandleFunc("/done/{id:[0-9]+}", s.handleDone()).Methods("GET")
	s.router.HandleFunc("/delete/{id:[0-9]+}", s.handleDelete()).Methods("DELETE")
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		err := s.db.Select(&todos, "SELECT * FROM todo")
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
			if _, err := s.db.Exec(stmt, title, false, time.Now()); err != nil {
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
		err := s.db.Get(&todo, "SELECT * FROM todo where id = ?", vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(todo)
		stmt := `update todo set done = 1 where id = ?`
		if _, err := s.db.Exec(stmt, vars["id"]); err != nil {
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
	s.router.ServeHTTP(w, r)
}

// func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if !currentUser(r).IsAdmin {
// 			http.NotFound(w, r)
// 			return
// 		}
// 		h(w, r)
// 	}
// }

type Config struct {
	HostPort string `default:"localhost:3000"`
	Database string `default:"todo.db"` // path to sqlite3
}

func main() {
	var config Config
	err := envconfig.Process("todo", &config)
	if err != nil {
		log.Fatal(err)
	}
	db, err := sqlx.Open("sqlite3", config.Database)
	if err != nil {
		log.Fatal(err)
	}
	svr := &Server{
		db:     db,
		router: mux.NewRouter(),
	}
	svr.routes()
	log.Printf("starting service at http://%v", config.HostPort)
	log.Fatal(http.ListenAndServe(config.HostPort, svr))
}
