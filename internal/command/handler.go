package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
)

// Handler ...
type Handler struct {
	logger *log.Logger
	// output is where the handler will write the results.
	// It is initialized to os.Stdout as default
	output io.Writer
	client *rest.Client
	// A pointer to the header flag instance, i.e. headers
	// provided as a flag will be inserted here (or into it's values)
	header *Header
}

func NewHandler(client *rest.Client, logger *log.Logger, h *Header) *Handler {
	return &Handler{
		logger: logger,
		output: os.Stdout,
		client: client,
		header: h,
	}
}

func (handler *Handler) Verbose(v bool) {
	if !v {
		handler.logger.SetOutput(ioutil.Discard)
	}
}

func (handler *Handler) Timeout(timeout time.Duration) {
	handler.client.Timeout(timeout)
}

func (handler *Handler) Get(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodGet, nil, cmd, args)
}

func (handler *Handler) Head(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodHead, nil, cmd, args)
}

func (handler *Handler) Post(cmd *cobra.Command, args []string) {
	body := handler.expectBody(cmd)
	handler.handleRequest(http.MethodPost, body, cmd, args)
}

func (handler *Handler) Patch(cmd *cobra.Command, args []string) {
	body := handler.expectBody(cmd)
	handler.handleRequest(http.MethodPatch, body, cmd, args)
}

func (handler *Handler) Put(cmd *cobra.Command, args []string) {
	body := handler.expectBody(cmd)
	handler.handleRequest(http.MethodPut, body, cmd, args)
}

func (handler *Handler) Delete(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodDelete, nil, cmd, args)
}

func (handler *Handler) handleRequest(method string, body []byte, cmd *cobra.Command, args []string) {
	url := args[0]

	headers, err := handler.getHeaders()
	handler.checkUserError(err, cmd)

	req, err := handler.client.BuildRequest(method, url, body, headers)
	handler.checkUserError(err, cmd)

	signRequest, _ := cmd.Flags().GetBool(constants.AWSSigV4FlagName)
	if signRequest {
		region, _ := cmd.Flags().GetString(constants.AWSRegionFlagName)
		profile, _ := cmd.Flags().GetString(constants.AWSProfileFlagName)

		err = handler.client.SignRequest(req, nil, region, profile)
		handler.checkExecutionError(err)
	}

	res := handler.client.SendRequest(req)
	handler.checkExecutionError(res.Error())
	handler.outputResults(cmd, res)
}

func (handler *Handler) outputResults(cmd *cobra.Command, r *rest.Result) {
	silent, _ := cmd.Flags().GetBool(constants.SilentFlagName)
	if silent {
		return
	}

	body, err := r.Body()
	handler.checkExecutionError(err)
	_, err = handler.output.Write(body)
	handler.checkExecutionError(err)
}

// Get the request headers from the handler header field as well as
// the environment variable for default headers.
func (handler *Handler) getHeaders() (http.Header, error) {
	headers := handler.header.values
	val, set := os.LookupEnv("DEFAULT_HEADERS")
	if !set {
		return headers, nil
	}

	// val is a string containing headers separated by a vertical pipe: |
	for _, h := range strings.Split(val, "|") {
		key, value, err := parseHeader(strings.TrimSpace(h))
		if err != nil {
			return headers, fmt.Errorf("invalid header format in DEFAULT_HEADERS: %w", err)
		}

		headers.Add(key, value)
	}

	return headers, nil
}

func (handler *Handler) checkExecutionError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func (handler *Handler) checkUserError(err error, cmd *cobra.Command) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	cmd.Usage()
	os.Exit(1)
}

func (handler *Handler) expectBody(cmd *cobra.Command) []byte {
	bodyFlag, _ := cmd.Flags().GetString(constants.BodyFlagName)
    bodyFlag = strings.TrimSpace(bodyFlag)

	if bodyFlag == "" {
		handler.logger.Printf("Empty body")
        return nil
	}

	// We first try to read as a file
	body, err := ioutil.ReadFile(bodyFlag)
	if err != nil && !os.IsNotExist(err) {
		handler.logger.Printf("Failed to open input file: %v", err)
		handler.checkUserError(err, cmd)
	}

	if body == nil {
		// Assume that the content was given as string
		handler.logger.Print("Assuming body was given as content string")
		body = []byte(bodyFlag)
	}

	return body
}
