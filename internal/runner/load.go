package runner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lunjon/httpreq/internal/rest"
	"gopkg.in/yaml.v2"
)

func Load(filepath string) (runner *Runner, err error) {
	var spec *RunnerSpec
	switch {
	case strings.HasSuffix(filepath, ".json"):
		spec, err = loadJSONSpec(filepath)
	case strings.HasSuffix(filepath, ".json"):
		spec, err = loadJSONSpec(filepath)
	case strings.HasSuffix(filepath, ".yaml"), strings.HasSuffix(filepath, ".yml"):
		spec, err = loadYAMLSpec(filepath)
	default:
		err = fmt.Errorf("invalid file ending: requires JSON or YAML file")
		return
	}

	if err != nil {
		return
	}

	err = validateSpec(spec)
	if err != nil {
		return
	}

	runner = &Runner{Spec: spec}
	return
}

func validateSpec(spec *RunnerSpec) error {
	if spec.Requests == nil {
		return fmt.Errorf("missing required field 'requests'")
	}

	if len(spec.Requests) == 0 {
		return fmt.Errorf("requests cannot be empty")
	}

	// Keep track of IDs to guarantee that they are unique
	ids := map[string]bool{}

	for _, req := range spec.Requests {
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
	}

	return nil
}

func loadJSONSpec(filepath string) (*RunnerSpec, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var spec RunnerSpec
	err = json.Unmarshal(bytes, &spec)
	return &spec, err
}

func loadYAMLSpec(filepath string) (*RunnerSpec, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var spec RunnerSpec
	err = yaml.Unmarshal(bytes, &spec)
	return &spec, err
}
