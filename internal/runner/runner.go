package runner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lunjon/httpreq/internal/rest"
)

type Runner struct {
	Spec   *RunnerSpec
	hasRun bool
}

func NewRunner(spec *RunnerSpec) *Runner {
	return &Runner{spec, false}
}

func (runner *Runner) Run(client *rest.Client, targets ...string) ([]*rest.Result, error) {
	if runner.hasRun {
		panic("a runner may only runner once")
	}

	var requests []*RequestTarget
	if len(targets) > 0 {
		for _, t := range targets {
			req, err := runner.findTarget(t)
			if err != nil {
				return nil, err
			}
			requests = append(requests, req)
		}

	} else {
		requests = runner.Spec.Requests
	}

	results := []*rest.Result{}
	for _, req := range requests {
		res, err := run(req, client)
		if err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	runner.hasRun = true

	return results, nil
}

func (runner *Runner) findTarget(target string) (*RequestTarget, error) {
	for _, req := range runner.Spec.Requests {
		if req.ID == target {
			return req, nil
		}
	}

	return nil, fmt.Errorf("unknown request ID: %s", target)
}

func run(req *RequestTarget, client *rest.Client) (res *rest.Result, err error) {
	header := http.Header{}
	for k, v := range req.Headers {
		header.Add(k, v)
	}

	var body []byte
	if req.Method == http.MethodPost {
		body, err = json.Marshal(req.Body)
		if err != nil {
			return
		}
		header.Add("Content-type", "application/json")
	}

	r, err := client.BuildRequest(req.Method, req.URL, body, header)
	if err != nil {
		return
	}

	if req.AWS != nil {
		err = client.SignRequest(r, body, req.AWS.Region, req.AWS.Profile)
		if err != nil {
			return
		}
	}

	res = client.SendRequest(r)
	return
}
