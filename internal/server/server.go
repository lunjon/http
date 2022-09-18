package server

import (
	"fmt"
	"net/http"

	"github.com/lunjon/http/internal/style"
)

type Server struct {
	server    *http.Server
	handler   *requestHandler
	formatter formatter
	callback  chan *http.Request
	port      uint
}

func New(port uint) *Server {
	handler := newHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.handle)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return &Server{
		server:    s,
		handler:   handler,
		formatter: defaultFormatter{},
		callback:  make(chan *http.Request),
		port:      port,
	}
}

func (s *Server) Close() error {
	close(s.callback)
	return s.server.Close()
}

func (s *Server) Serve() error {
	go onRequest(s.formatter, s.callback)

	fmt.Printf("Starting server on :%s.\n", style.Bold.Render(fmt.Sprint(s.port)))
	fmt.Printf("Press %s to exit.\n", style.Bold.Render("CTRL-C"))
	return s.server.ListenAndServe()
}

func onRequest(f formatter, ch chan *http.Request) {
	for r := range ch {
		f.format(r)
	}
}
