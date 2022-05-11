package server

import (
	"fmt"
	"net/http"
)

type Config struct {
	Port uint
}

func (c Config) addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

type Server struct {
	server *http.Server
	cfg    Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	s := &http.Server{
		Addr:    cfg.addr(),
		Handler: mux,
	}
	return &Server{s, cfg}
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {}
