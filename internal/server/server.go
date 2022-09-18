package server

import (
	"fmt"
	"net/http"

	"github.com/lunjon/http/internal/style"
)

type Options struct {
	Port        uint
	ShowSummary bool
}

type Server struct {
	server  *http.Server
	handler *requestHandler
	cb      callback
	ch      chan *http.Request
	options Options
}

func New(opts Options) *Server {
	ch := make(chan *http.Request)

	handler := newHandler(ch)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.handle)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: mux,
	}

	return &Server{
		server:  s,
		handler: handler,
		cb:      defaultCallback{},
		ch:      ch,
		options: opts,
	}
}

func (s *Server) Serve() error {
	err := s.cb.start()
	if err != nil {
		return err
	}

	go s.onRequest()

	fmt.Printf("Starting server on :%s.\n", style.Bold.Render(fmt.Sprint(s.options.Port)))
	fmt.Printf("Press %s to exit.\n", style.Bold.Render("CTRL-C"))
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	close(s.ch)
	return s.server.Close()
}

func (s *Server) onRequest() {
	for r := range s.ch {
		s.cb.handle(r)
	}
}
