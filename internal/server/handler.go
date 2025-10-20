package server

import (
	"net/http"
	"strconv"
	"time"
)

type requestHandler struct {
	ch         chan *http.Request
	statusCode int
}

func newHandler(ch chan *http.Request) *requestHandler {
	return &requestHandler{
		ch:         ch,
		statusCode: 200,
	}
}

func (h *requestHandler) handleDefault(w http.ResponseWriter, r *http.Request) {
	h.ch <- r
	w.WriteHeader(h.statusCode)
}

func (h *requestHandler) handleWithCode(w http.ResponseWriter, r *http.Request) {
	h.ch <- r

	code := r.PathValue("code")
	status := http.StatusOK
	if code != "" {
		if n, err := strconv.Atoi(code); err == nil {
			status = n
		}
	}

	w.WriteHeader(status)
}

func (h *requestHandler) handleTimeout(w http.ResponseWriter, r *http.Request) {
	h.ch <- r
	time.Sleep(time.Minute * 5)
}
