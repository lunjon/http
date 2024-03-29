package server

import (
	"fmt"
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
	}
}

type statusCallback struct{}

func newStatusCallback() *statusCallback {
	return &statusCallback{}
}

func (s statusCallback) loop(reqs chan *http.Request, stop chan bool) {
	ticker := time.NewTicker(time.Second)
	t := ticker.C

	count := 0 // Count of requests per second

	for {
		select {
		case <-reqs:
			count++
		case <-t:
			fmt.Printf("%d requests/s\n", count)
			count = 0
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
