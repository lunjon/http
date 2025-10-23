package server

import (
	"fmt"
	"net/http"

	"github.com/lunjon/http/internal/style"
)

type Options struct {
	Port       uint
	ShowStatus bool
	StaticRoot string
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

	if opts.StaticRoot != "" {
		root := http.Dir(opts.StaticRoot)
		fs := http.FileServer(root)
		mux.Handle("/", fs)
	} else {
		mux.HandleFunc("/~/status/{code}", handler.handleWithCode)
		mux.HandleFunc("/~/timeout", handler.handleTimeout)
		mux.HandleFunc("/", handler.handleDefault)
	}

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
	fmt.Println("/~/status/{code}  Respond with the given code as status.")
	fmt.Println("                  Send 'random' as the path parameter to get a random status.")
	fmt.Println("/~/timeout        Endpoint that hangs the request (for 5 min).")
	fmt.Println("/*                Echo the request with 200 OK status code.")
}
