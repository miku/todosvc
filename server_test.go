package todosvc

import (
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func TestHandleIndex(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	srv := &Server{
		DB:     sqlx.NewDb(db, "sqlmock"),
		Router: mux.NewRouter(),
	}
	srv.Routes()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mock.ExpectQuery("SELECT (.+) FROM todo")
	srv.ServeHTTP(w, req)
	io.Copy(os.Stdout, w.Body)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
