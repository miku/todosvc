package main

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mattn/go-sqlite3"
	"github.com/miku/todosvc"
)

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
	svr := &todosvc.Server{
		DB:     db,
		Router: mux.NewRouter(),
	}
	svr.Routes()
	h := csrf.Protect([]byte("abcdabcdabcdabcdabcdabcdabcdabcd"))(svr)
	log.Printf("starting service at http://%v", config.HostPort)
	log.Fatal(http.ListenAndServe(config.HostPort, h))
}
