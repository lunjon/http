package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lunjon/http/style"
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
		fmt.Println("Incoming request:")
		fmt.Printf("  Method:  %s\n", style.GreenB(r.Method))
		fmt.Printf("  Path:    %s\n", style.GreenB(r.URL.Path))
		if len(r.Header) > 0 {
			fmt.Println("  Headers:")
			for name, values := range r.Header {
				v := strings.Join(values, style.Cyan("; "))
				v = fmt.Sprintf("%s %s %s", style.CyanB("["), v, style.CyanB("]"))
				fmt.Printf("    %s: %s\n", style.WhiteB(name), v)
			}
		}
	}
}
