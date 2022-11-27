package options

import (
	"net/http"

	"github.com/lunjon/http/internal/util"
)

type HeaderOption struct {
	splitter *util.Splitter
	values   http.Header
}

func NewHeaderOption() *HeaderOption {
	return &HeaderOption{
		splitter: util.NewSplitter(":"),
		values:   make(http.Header),
	}
}

func (h *HeaderOption) Header() http.Header {
	return h.values
}

// Append adds the provided value as a header if it is valid
func (h *HeaderOption) Set(s string) error {
	key, value, err := h.splitter.Parse(s)
	if err != nil {
		return err
	}
	h.values.Add(key, value)
	return nil
}

func (h *HeaderOption) Type() string {
	return "Header"
}

func (h *HeaderOption) String() string {
	return ""
}
