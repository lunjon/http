package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lunjon/http/internal/style"
)

type requestHandler struct {
	statusCode int
}

func newHandler() *requestHandler {
	return &requestHandler{
		statusCode: 200,
	}
}

func (h *requestHandler) handle(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(h.statusCode)
}
