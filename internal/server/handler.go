package server

import (
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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
	var status int
	if strings.ToLower(code) == "random" {
		// Get a random status code
		status = statuses[rnd.Intn(len(statuses))]
	} else if code != "" {
		if n, err := strconv.Atoi(code); err == nil {
			status = n
		} else {
			status = http.StatusTeapot
		}
	} else {
		status = http.StatusNotFound
	}

	w.WriteHeader(status)
}

func (h *requestHandler) handleTimeout(w http.ResponseWriter, r *http.Request) {
	h.ch <- r
	time.Sleep(time.Minute * 5)
}

var (
	rnd      = rand.New(rand.NewSource(time.Now().Unix()))
	statuses = []int{
		// 2XX
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNoContent,
		// 4XX
		http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusConflict,
		http.StatusTeapot,
		// 5XX
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
)
