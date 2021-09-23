package command

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lunjon/http/client"
	"github.com/spf13/cobra"
)

var (
	newline           = []byte("\n")
	errBriefAndSilent = errors.New("cannot specify both --brief and --silent")
	errCertFlags      = errors.New("--cert-pub-file requires --cert-key-file and vice versa")
	emptyRequestBody  = requestBody{}
)

// Handler handles all commands.
type Handler struct {
	// The directory of the CLI in home (or other for testing)
	gohttpDir     string
	aliasFilePath string
	// General logger
	logger *log.Logger
	// Specific logger for TLS tracing
	traceLogger *log.Logger
	// Output of infos
	infos io.Writer
	// Output of errors
	errors io.Writer
	// HTTP client. Configured in Handler.Init()
	client *client.Client
	// Field that all headers are set
	header *HeaderOption
	// Function to invoke on errors which cannot be recovered from
	exitFunc func()
	// Set on init
	cmd       *cobra.Command
	fail      bool
	formatter Formatter
}

func NewHandler(
	client *client.Client,
	logger *log.Logger,
	traceLogger *log.Logger,
	infos io.Writer,
	errors io.Writer,
	dir string,
	exitFunc func(),
) *Handler {
	return &Handler{
		gohttpDir:     dir,
		aliasFilePath: path.Join(dir, "alias"),
		logger:        logger,
		traceLogger:   traceLogger,
		infos:         infos,
		errors:        errors,
		client:        client,
		header:        NewHeaderOption(),
		formatter:     DefaultFormatter{},
		exitFunc:      exitFunc,
	}
}

func (handler *Handler) init(cmd *cobra.Command) {
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
	handler.handleRequest(http.MethodGet, emptyRequestBody, cmd, args)
}

func (handler *Handler) Head(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodHead, emptyRequestBody, cmd, args)
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
	handler.handleRequest(http.MethodDelete, emptyRequestBody, cmd, args)
}

func (handler *Handler) handleRequest(method string, body requestBody, cmd *cobra.Command, args []string) {
	alias, err := handler.readAliasFile()
	handler.checkExecutionError(err)

	url, err := client.ParseURL(args[0], alias)
	handler.checkUserError(err, cmd)

	headers, err := handler.getHeaders()
	handler.checkUserError(err, cmd)
	setContentType := headers.Get("content-type") == "" && body.mime != client.MIMETypeUnknown
	if setContentType {
		handler.logger.Printf("Detected MIME type: %s", body.mime)
		headers.Add("Content-Type", body.mime.String())
	}

	repeat, _ := cmd.Flags().GetInt(repeatFlagName)
	for i := 0; i < repeat; i++ {
		req, err := handler.buildRequest(cmd, method, url, body.bytes, headers)
		handler.checkUserError(err, cmd)

		res := handler.client.SendRequest(req)
		handler.checkExecutionError(res.Error())
		handler.outputResults(cmd, res)
	}
}

func (handler *Handler) buildRequest(
	cmd *cobra.Command,
	method string,
	url *client.URL,
	body []byte,
	header http.Header,
) (*http.Request, error) {
	req, err := handler.client.BuildRequest(method, url, body, header)
	if err != nil {
		return nil, err
	}

	signRequest, _ := cmd.Flags().GetBool(awsSigV4FlagName)
	if signRequest {
		region, _ := cmd.Flags().GetString(awsRegionFlagName)

		err = handler.client.SignRequest(req, body, region)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (handler *Handler) outputResults(cmd *cobra.Command, r *client.Result) {
	b, err := handler.formatter.Format(r)
	handler.checkExecutionError(err)

	if len(b) > 0 {
		_, err = handler.infos.Write(b)
		handler.checkExecutionError(err)
		_, _ = handler.infos.Write(newline)
	}

	if handler.fail && !r.Successful() {
		handler.logger.Printf("Request failed with status %s", r.Status())
		handler.exitFunc()
	}
}

// Get the request headers from the handler header field as well as
// the environment variable for default headers.
func (handler *Handler) getHeaders() (http.Header, error) {
	headers := handler.header.values
	val, set := os.LookupEnv(defaultHeadersEnv)
	if !set {
		return headers, nil
	}

	// val is a string containing headers separated by a vertical pipe: |
	for _, h := range strings.Split(val, "|") {
		key, value, err := parseHeader(h)
		if err != nil {
			return headers, fmt.Errorf("invalid header format in %s: %w", defaultHeadersEnv, err)
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
	handler.exitFunc()
}

func (handler *Handler) checkUserError(err error, cmd *cobra.Command) {
	if err == nil {
		return
	}
	fmt.Fprintf(handler.errors, "error: %v\n\n", err)
	cmd.Usage()
	handler.exitFunc()
}

type requestBody struct {
	bytes []byte
	mime  client.MIMEType
}

func (handler *Handler) getRequestBody(cmd *cobra.Command) requestBody {
	bodyFlag, _ := cmd.Flags().GetString(bodyFlagName)
	bodyFlag = strings.TrimSpace(bodyFlag)

	mime := client.MIMETypeUnknown

	if bodyFlag == "" {
		// Not provided via flags, check stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			handler.logger.Print("Reading body from stdin")
			b, err := io.ReadAll(os.Stdin)
			handler.checkExecutionError(err)

			return requestBody{b, mime}
		}

		handler.logger.Print("No body provided")
		return requestBody{nil, mime}
	}

	// We first try to read as a file
	body, err := ioutil.ReadFile(bodyFlag)
	if err != nil && !os.IsNotExist(err) {
		handler.checkUserError(err, cmd)
	}

	if body != nil {
		// Try detecting filetype in order to set MIME type
		switch path.Ext(bodyFlag) {
		case ".html":
			mime = client.MIMETypeHTML
		case ".csv":
			mime = client.MIMETypeCSV
		case ".json":
			mime = client.MIMETypeJSON
		case ".xml":
			mime = client.MIMETypeXML
		}
	}

	if body == nil {
		// Assume that the content was given as string
		body = []byte(bodyFlag)
	}

	return requestBody{body, mime}
}
