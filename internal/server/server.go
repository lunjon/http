package server

import (
	"fmt"
	"net/http"

	"github.com/lunjon/http/internal/style"
)

type Options struct {
	Port       uint
	ShowStatus bool
}

type Server struct {
	server  *http.Server
	handler *requestHandler
	cb      callback
	ch      chan *http.Request
	done    chan bool
	options Options
}

func New(opts Options) *Server {
	ch := make(chan *http.Request)

	handler := newHandler(ch)
	mux := http.NewServeMux()
	mux.HandleFunc("/~/success/", handler.handleSuccess)
	mux.HandleFunc("/~/client-errors/", handler.handleClientErrors)
	mux.HandleFunc("/~/server-errors/", handler.handleServerErrors)
	mux.HandleFunc("/", handler.handleDefault)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: mux,
	}

	var cb callback
	if opts.ShowStatus {
		cb = newStatusCallback()
	} else {
		cb = defaultCallback{}
	}

	return &Server{
		server:  s,
		handler: handler,
		cb:      cb,
		ch:      ch,
		options: opts,
	}
}

func (s *Server) Serve() error {
	go s.cb.loop(s.ch, s.done)

	fmt.Printf("Starting server on :%s.\n", style.Bold.Render(fmt.Sprint(s.options.Port)))
	fmt.Printf("Press %s to exit.\n", style.Bold.Render("CTRL-C"))
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	close(s.ch)
	go func() {
		s.done <- true
	}()
	return s.server.Close()
}

func ListRoutes() {
	fmt.Println("Success routes:")
	for route, methods := range successResponses {
		for method, resp := range methods {
			fmt.Printf("  %s %s (status %d)\n", method, route, resp.status)
		}
	}

	// FIXME: sort by status
	fmt.Println("\nClient Errors (accepts any method)")
	for route, status := range clientErrors {
		fmt.Printf("  %s %d\n", route, status)
	}

	// FIXME: sort by status
	fmt.Println("\nServer Errors (accepts any method)")
	for route, status := range serverErrors {
		fmt.Printf("  %s %d\n", route, status)
	}
}
