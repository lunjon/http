package runner

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/lunjon/httpreq/pkg/parse"
)

/*
AWSSignConfig represents the configuration
for signing requests using AWS signature V4.
*/
type AWSSignConfig struct {
	Profile string
	Region  string
}

func NewAWSSign(profile, region string) *AWSSignConfig {
	return &AWSSignConfig{
		Profile: profile,
		Region:  region,
	}
}

/*
RequestTarget describe the model of the requests in the request file.
*/
type RequestTarget struct {
	ID     string
	Method string
	URL    string

	Headers map[string]string

	Body map[string]interface{}

	AWS interface{}
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

/*
Set the base URL of the request to the new URL,
but keep the path of the original request.
*/
func (req *RequestTarget) SetBaseURL(url string) error {
	u, err := parse.ParseURL(url)
	if err != nil {
		return err
	}
	r, err := parse.ParseURL(req.URL)
	if err != nil {
		return err
	}

	req.URL = u.BaseURL() + r.Path
	return nil
}

/*
Validate that the request is valid.
Should be called before anything else.
*/
func (req *RequestTarget) Validate(ids map[string]bool) error {
	// ID

	if req.ID == "" {
		return fmt.Errorf("invalid or missing ID in request")
	}

	if strings.ContainsAny(req.ID, " ") {
		return fmt.Errorf("IDs cannot contain any whitespace")
	}

	if _, found := ids[req.ID]; found {
		return fmt.Errorf("duplicate ID: %s", req.ID)
	}

	// Method

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

	// AWS Signing
	switch req.AWS.(type) {
	case nil:
		req.AWS = nil
	case bool, string:
		req.AWS = NewAWSSign("", constants.DefaultAWSRegion)
	case map[interface{}]interface{}:
		v := req.AWS.(map[interface{}]interface{})
		profile := "default"
		region := constants.DefaultAWSRegion
		if p, found := v["profile"]; found {
			profile = p.(string)
		}

		if r, found := v["region"]; found {
			region = r.(string)
		}

		req.AWS = NewAWSSign(profile, region)
	}

	ids[req.ID] = true
	return nil
}

func (req *RequestTarget) GetAWSSign() *AWSSignConfig {
	if req.AWS != nil {
		return req.AWS.(*AWSSignConfig)
	}

	return nil
}

/*
Spec is the specification of runner files.
It's only used to load files from the system.
*/
type Spec struct {
	Headers  map[string]string
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
