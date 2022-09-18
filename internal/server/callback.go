package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lunjon/http/internal/style"
)

type callback interface {
	loop(chan *http.Request, chan bool)
}

type defaultCallback struct{}

func (f defaultCallback) loop(ch chan *http.Request, stop chan bool) {
	for r := range ch {
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
}

type statusCallback struct{}

func newStatusCallback() *statusCallback {
	return &statusCallback{}
}

func (s statusCallback) loop(reqs chan *http.Request, stop chan bool) {
	ticker := time.NewTicker(time.Second)
	t := ticker.C
	curr := 0
	prev := 0

	for {
		select {
		case <-reqs:
			curr++
		case <-t:
			if curr == prev {
				fmt.Println("0 requests/s")
				continue
			}

			fmt.Printf("~ %d requests/s\n", curr)
			prev = curr
			curr = 0
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
