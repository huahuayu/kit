package server

import "net/http"

type HttpServer interface {
	Serve(host, port string, routes map[string]http.HandlerFunc) error
}

type Server struct {
	http.Server
}

func New() HttpServer {
	return &Server{
		Server: http.Server{},
	}
}

func (s *Server) Serve(host, port string, routes map[string]http.HandlerFunc) error {
	s.Addr = host + ":" + port
	for route, handler := range routes {
		http.HandleFunc(route, handler)
	}
	return s.ListenAndServe()
}
