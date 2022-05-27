package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunjon/http/internal/format"
)

var (
	bold      = lipgloss.NewStyle().Bold(true)
	redB      = bold.Copy().Foreground(lipgloss.Color("1")).Render
	greenB    = bold.Copy().Foreground(lipgloss.Color("2")).Render
	greenish  = bold.Copy().Foreground(lipgloss.Color("48"))
	greenishB = greenish.Copy().Bold(true)
	grey      = bold.Copy().Foreground(lipgloss.Color("245")).Render
)

type Server struct {
	styler *format.Styler
	server *http.Server
}

func New(port uint, styler *format.Styler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	return &Server{
		server: s,
	}
}

func (s *Server) Close() error {
	return s.server.Close()
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming request:")
	fmt.Printf("  Method:  %s\n", greenB(r.Method))
	fmt.Printf("  Path:    %s\n", greenB(r.URL.Path))

	if len(r.Header) > 0 {
		fmt.Println("  Headers:")
		for name, values := range r.Header {
			v := strings.Join(values, grey("; "))
			v = fmt.Sprintf("%s %s %s", greenishB.Render("["), v, greenishB.Render("]"))
			fmt.Printf("    %s: %s\n", bold.Render(name), v)
		}
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("  Body:    %s: %s\n", redB("error"), err)
	} else if len(b) > 0 {
		body := string(b)
		if len(body) > 100 {
			body = body[:90] + "..."
		}
		fmt.Printf("  Body:    %s\n", body)
	}
}
