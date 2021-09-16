package command

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lunjon/http/rest"
	"github.com/spf13/cobra"
)

var (
	newline           = []byte("\n")
	errBriefAndSilent = errors.New("cannot specify both --brief and --silent")
	errCertFlags      = errors.New("--cert-pub-file requires --cert-key-file and vice versa")
)

// Handler ...
type Handler struct {
	logger      *log.Logger
	traceLogger *log.Logger
	infos       io.Writer
	errors      io.Writer
	client      *rest.Client
	header      *HeaderOption
	// Set on init
	cmd       *cobra.Command
	fail      bool
	formatter Formatter
}

func NewHandler(
	client *rest.Client,
	logger *log.Logger,
	traceLogger *log.Logger,
	h *HeaderOption) *Handler {
	return &Handler{
		logger:      logger,
		traceLogger: traceLogger,
		infos:       os.Stdout,
		errors:      os.Stderr,
		client:      client,
		header:      h,
		formatter:   DefaultFormatter{},
	}
}

func (handler *Handler) Init(cmd *cobra.Command) {
	if cmd == nil {
		return
	}

	timeout, _ := cmd.Flags().GetDuration(timeoutFlagName)
	handler.Timeout(timeout)

	certPub, _ := cmd.Flags().GetString(certpubFlagName)
	certKey, _ := cmd.Flags().GetString(certkeyFlagName)
	if certPub != "" && certKey == "" {
		handler.checkUserError(errCertFlags, cmd)
	} else if certPub == "" && certKey != "" {
		handler.checkUserError(errCertFlags, cmd)
	} else if certPub != "" && certKey != "" {
		err := handler.client.Cert(certPub, certKey)
		handler.checkUserError(err, cmd)
	}

	verbose, _ := cmd.Flags().GetBool(verboseFlagName)
	if verbose {
		handler.logger.SetOutput(os.Stderr)
	} else {
		handler.logger.SetOutput(ioutil.Discard)
	}

	trace, _ := cmd.Flags().GetBool(traceFlagName)
	if trace {
		handler.traceLogger.SetOutput(os.Stderr)
	} else {
		handler.traceLogger.SetOutput(ioutil.Discard)
	}

	handler.fail, _ = cmd.Flags().GetBool(failFlagName)

	brief, _ := cmd.Flags().GetBool(briefFlagName)
	silent, _ := cmd.Flags().GetBool(silentFlagName)

	if brief && silent {
		handler.checkUserError(errBriefAndSilent, cmd)
	}

	if silent {
		handler.formatter = NullFormatter{}
	}
	if brief {
		handler.formatter = BriefFormatter{}
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
	body := handler.getRequestBody(cmd)
	handler.handleRequest(http.MethodPost, body, cmd, args)
}

func (handler *Handler) Patch(cmd *cobra.Command, args []string) {
	body := handler.getRequestBody(cmd)
	handler.handleRequest(http.MethodPatch, body, cmd, args)
}

func (handler *Handler) Put(cmd *cobra.Command, args []string) {
	body := handler.getRequestBody(cmd)
	handler.handleRequest(http.MethodPut, body, cmd, args)
}

func (handler *Handler) Delete(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodDelete, nil, cmd, args)
}

func (handler *Handler) handleRequest(method string, body []byte, cmd *cobra.Command, args []string) {
	alias, err := handler.readAliasFile()
	handler.checkExecutionError(err)

	url, err := rest.ParseURL(args[0], alias)
	handler.checkUserError(err, cmd)

	headers, err := handler.getHeaders()
	handler.checkUserError(err, cmd)

	repeat, _ := cmd.Flags().GetInt(repeatFlagName)
	for i := 0; i < repeat; i++ {
		req, err := handler.buildRequest(cmd, method, url, body, headers)
		handler.checkUserError(err, cmd)

		res := handler.client.SendRequest(req)
		handler.checkExecutionError(res.Error())
		handler.outputResults(cmd, res)
	}
}

func (handler *Handler) buildRequest(cmd *cobra.Command, method string, url *rest.URL, body []byte, header http.Header) (*http.Request, error) {
	req, err := handler.client.BuildRequest(method, url, body, header)
	if err != nil {
		return nil, err
	}

	signRequest, _ := cmd.Flags().GetBool(awsSigV4FlagName)
	if signRequest {
		region, _ := cmd.Flags().GetString(awsRegionFlagName)
		profile, _ := cmd.Flags().GetString(awsProfileFlagName)

		err = handler.client.SignRequest(req, body, region, profile)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (handler *Handler) outputResults(cmd *cobra.Command, r *rest.Result) {
	b, err := handler.formatter.Format(r)
	handler.checkExecutionError(err)

	if len(b) > 0 {
		_, err = handler.infos.Write(b)
		handler.checkExecutionError(err)
		_, _ = handler.infos.Write(newline)
	}

	if handler.fail && !r.Successful() {
		handler.logger.Printf("Request failed with status %s", r.Status())
		os.Exit(1)
	}
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
	fmt.Fprintf(handler.errors, "error: %v\n", err)
	os.Exit(1)
}

func (handler *Handler) checkUserError(err error, cmd *cobra.Command) {
	if err == nil {
		return
	}
	fmt.Fprintf(handler.errors, "error: %v\n\n", err)
	cmd.Usage()
	os.Exit(1)
}

func (handler *Handler) getRequestBody(cmd *cobra.Command) []byte {
	bodyFlag, _ := cmd.Flags().GetString(bodyFlagName)
	bodyFlag = strings.TrimSpace(bodyFlag)

	if bodyFlag == "" {
		// Not provided via flags, check stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			handler.logger.Print("Reading body from stdin")
			b, err := io.ReadAll(os.Stdin)
			handler.checkExecutionError(err)
			return b
		}

		handler.logger.Print("No body provided")
		return nil
	}

	// We first try to read as a file
	body, err := ioutil.ReadFile(bodyFlag)
	if err != nil && !os.IsNotExist(err) {
		handler.checkUserError(err, cmd)
	}

	if body == nil {
		// Assume that the content was given as string
		body = []byte(bodyFlag)
	}

	return body
}
