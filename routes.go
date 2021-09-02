package todosvc

import "net/http"

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.Router.HandleFunc("/listing", s.handleListing()).Methods("GET") // json response
	s.Router.HandleFunc("/new", s.handleNew()).Methods("GET", "POST")
	s.Router.HandleFunc("/done/{id:[0-9]+}", s.adminOnly(s.handleDone())).Methods("GET")
	// s.Router.HandleFunc("/delete/{id:[0-9]+}", s.handleDelete()).Methods("DELETE")
}

func (s *Server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAdmin {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}
