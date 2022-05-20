package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lunjon/http/format"
)

type Server struct {
	styler *format.Styler
	server *http.Server
}

func New(port uint, styler *format.Styler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", createHandler(styler))

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

func createHandler(styler *format.Styler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Incoming request:")
		fmt.Printf("  Method:  %s\n", styler.GreenB(r.Method))
		fmt.Printf("  Path:    %s\n", styler.GreenB(r.URL.Path))
		if len(r.Header) > 0 {
			fmt.Println("  Headers:")
			for name, values := range r.Header {
				v := strings.Join(values, styler.Cyan("; "))
				v = fmt.Sprintf("%s %s %s", styler.CyanB("["), v, styler.CyanB("]"))
				fmt.Printf("    %s: %s\n", styler.WhiteB(name), v)
			}
		}
	}
}
