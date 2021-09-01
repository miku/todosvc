package todosvc

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.Router.HandleFunc("/new", s.handleNew()).Methods("GET", "POST")
	s.Router.HandleFunc("/done/{id:[0-9]+}", s.handleDone()).Methods("GET")
	s.Router.HandleFunc("/delete/{id:[0-9]+}", s.handleDelete()).Methods("DELETE")
}
