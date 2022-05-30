package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lunjon/http/internal/style"
)

type Server struct {
	server *http.Server
	port   uint
}

func New(port uint) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return &Server{
		server: s,
		port:   port,
	}
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Serve() error {
	fmt.Printf("Starting server on :%s.\n", style.Bold.Render(fmt.Sprint(s.port)))
	fmt.Printf("Press %s to exit.\n", style.Bold.Render("CTRL-C"))
	return s.server.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request:")
	fmt.Printf("  Method:  %s\n", style.GreenB.Render(r.Method))
	fmt.Printf("  Path:    %s\n", style.GreenB.Render(r.URL.Path))

	if len(r.Header) > 0 {
		fmt.Println("  Headers:")
		for name, values := range r.Header {
			v := strings.Join(values, style.Grey.Render("; "))
			v = fmt.Sprintf("%s %s %s", style.GreenB.Render("["), v, style.GreenB.Render("]"))
			fmt.Printf("    %s: %s\n", style.Bold.Render(name), v)
		}
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("  Body:    %s: %s\n", style.RedB.Render("error"), err)
	} else if len(b) > 0 {
		s := fmt.Sprintf("%d bytes", len(b))
		fmt.Printf("  Body:    %s\n", style.GreenB.Render(s))
	}
}
