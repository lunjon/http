package runner

import (
	"fmt"
	"github.com/lunjon/httpreq/internal/rest"
	"net/http"
	"strings"
)

/*
RequestTarget describe the model in the request files.
*/
type RequestTarget struct {
	ID     string
	Method string
	URL    string

	Headers map[string]string

	Body map[string]interface{}

	AWS *struct {
		Profile string
		Region  string
	}
}

/*
TrySetHeader tries to set the header in the request.
If it is already set, it does not override the value.
*/
func (req *RequestTarget) TrySetHeader(key, value string) {
	if req.Headers == nil {
		req.Headers = map[string]string{}
	}

	if _, found := req.Headers[key]; !found {
		req.Headers[key] = value
	}
}

// Validate that the request is valid
func (req *RequestTarget) Validate(ids map[string]bool) error {
	if req.ID == "" {
		return fmt.Errorf("invalid or missing ID in request")
	}

	if strings.ContainsAny(req.ID, " ") {
		return fmt.Errorf("IDs cannot contain any whitespace")
	}

	method := strings.ToUpper(req.Method)
	if method == "" {
		method = http.MethodGet
	}

	if !(method == http.MethodGet || method == http.MethodPost || method == http.MethodDelete) {
		return fmt.Errorf("invalid HTTP method: %s", req.Method)
	}
	req.Method = method

	if _, err := rest.ParseURL(req.URL); err != nil {
		return err
	}

	if req.Method == http.MethodPost && req.Body == nil {
		return fmt.Errorf("missing body in POST request with ID '%s'", req.ID)
	}

	if _, found := ids[req.ID]; found {
		return fmt.Errorf("duplicate ID: %s", req.ID)
	}

	ids[req.ID] = true
	return nil
}


/*
Spec is the specification of runner files.
It's only used to load files from the system.
*/
type Spec struct {
	Headers map[string]string
	Requests []*RequestTarget
}

// Validate that the specification is valid
func (spec *Spec) Validate() error {
	if spec.Requests == nil {
		return fmt.Errorf("missing required field 'requests'")
	}

	if len(spec.Requests) == 0 {
		return fmt.Errorf("requests cannot be empty")
	}

	// Keep track of IDs to guarantee that they are unique
	ids := map[string]bool{}

	for _, req := range spec.Requests {
		err := req.Validate(ids)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetHeaders tries to set the default headers in each request
func (spec *Spec) SetHeaders() {
	for name, value := range spec.Headers {
		for _, req := range spec.Requests {
			req.TrySetHeader(name, value)
		}
	}
}