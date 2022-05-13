package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Config struct {
	Port   uint
	infos  io.Writer
	errors io.Writer
}

func (c Config) addr() string {
	return fmt.Sprintf(":%d", c.Port)
}

type Server struct {
	logger *log.Logger
	server *http.Server
	cfg    Config
	infos  io.Writer
	errors io.Writer
}

func New(cfg Config, l *log.Logger, infos, errors io.Writer) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", createHandler(l))

	s := &http.Server{
		Addr:    cfg.addr(),
		Handler: mux,
	}
	return &Server{
		server: s,
		cfg:    cfg,
		logger: l,
		infos:  infos,
		errors: errors,
	}
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func createHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
	}
}
