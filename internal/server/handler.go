package server

import (
	"net/http"
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

func (h *requestHandler) handle(w http.ResponseWriter, r *http.Request) {
	h.ch <- r
	w.WriteHeader(h.statusCode)
}
