package server

import (
	"encoding/json"
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

func (h *requestHandler) handleDefault(w http.ResponseWriter, r *http.Request) {
	h.ch <- r
	w.WriteHeader(h.statusCode)
}

func (h *requestHandler) handleSuccess(w http.ResponseWriter, r *http.Request) {
	h.ch <- r

	response, ok := successResponses[r.URL.Path]
	if ok {
		res, ok := response[r.Method]
		if ok {
			w.WriteHeader(res.status)
			if res.body != nil {
				writeJSON(w, res.body)
			}
			return
		}
	}

	w.WriteHeader(404)
}

func (h *requestHandler) handleClientErrors(w http.ResponseWriter, r *http.Request) {
	h.ch <- r

	var status int
	status, found := clientErrors[r.URL.Path]
	if !found {
		status = 404
	}

	w.WriteHeader(status)
}

func (h *requestHandler) handleServerErrors(w http.ResponseWriter, r *http.Request) {
	h.ch <- r

	var status int
	status, found := serverErrors[r.URL.Path]
	if !found {
		status = 404
	}

	w.WriteHeader(status)
}

func writeJSON(w http.ResponseWriter, body any) {
	b, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

var clientErrors map[string]int = map[string]int{
	"/~/client-errors/bad-request":                   400,
	"/~/client-errors/unauthorized":                  401,
	"/~/client-errors/payment-required":              402,
	"/~/client-errors/forbidden":                     403,
	"/~/client-errors/not-found":                     404,
	"/~/client-errors/method-not-allowed":            405,
	"/~/client-errors/not-acceptable":                406,
	"/~/client-errors/proxy-authentication-required": 407,
	"/~/client-errors/request-timeout":               408,
	"/~/client-errors/conflict":                      409,
	"/~/client-errors/gone":                          410,
	"/~/client-errors/too-many-requests":             429,
}

var serverErrors map[string]int = map[string]int{
	"/~/server-errors/internal-server-error":      500,
	"/~/server-errors/not-implemented":            501,
	"/~/server-errors/bad-gateway":                502,
	"/~/server-errors/service-unavailable":        503,
	"/~/server-errors/gateway-timeout":            504,
	"/~/server-errors/http-version-not-supported": 505,
}

var successResponses = map[string]methodResponse{
	"/~/success/cats": map[string]response{
		"GET": {
			status: 200,
			body:   cat{Name: "Kitty", Color: "Brown"},
		},
	},
}
