package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/lunjon/http/internal/alias"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/format"
	"github.com/lunjon/http/internal/util"
)

var (
	newline          = []byte("\n")
	errCertFlags     = errors.New("--cert must be used with --key")
	emptyRequestBody = requestBody{mime: client.MIMETypeUnknown}
)

const (
	userAgentHeader   = "User-Agent"
	contentTypeHeader = "Content-Type"
)

// RequestHandler handles all commands.
type RequestHandler struct {
	client         *client.Client
	aliasManager   alias.Manager
	formatter      format.ResponseFormatter
	signer         client.RequestSigner
	infos          io.Writer
	errors         io.Writer
	logger         *log.Logger
	fail           bool
	failFunc       FailFunc
	repeat         int
	defaultHeaders string
	headerOpt      *HeaderOption
	version        string
	outputFile     string
}

func newHandler(
	client *client.Client,
	aliasManager alias.Manager,
	formatter format.ResponseFormatter,
	signer client.RequestSigner,
	logger *log.Logger,
	failFunc FailFunc,
	cfg *config,
) *RequestHandler {
	return &RequestHandler{
		client:         client,
		aliasManager:   aliasManager,
		defaultHeaders: cfg.defaultHeaders,
		headerOpt:      cfg.headerOpt,
		logger:         logger,
		infos:          cfg.infos,
		errors:         cfg.errs,
		signer:         signer,
		formatter:      formatter,
		fail:           cfg.fail,
		failFunc:       failFunc,
		repeat:         cfg.repeat,
		version:        cfg.version,
		outputFile:     cfg.output,
	}
}

func (handler *RequestHandler) handleRequest(method, url, bodyflag string) error {
	body := emptyRequestBody
	if strings.Contains("post put patch", strings.ToLower(method)) {
		b, err := handler.getRequestBody(bodyflag)
		if err != nil {
			return err
		}
		body = b
	}

	alias, err := handler.aliasManager.Load()
	if err != nil {
		return err
	}

	u, err := client.ParseURL(url, alias)
	if err != nil {
		return err
	}

	headers, err := handler.getHeaders()
	if err != nil {
		return err
	}

	setContentType := headers.Get(contentTypeHeader) == "" && body.mime != client.MIMETypeUnknown
	if setContentType {
		handler.logger.Printf("Detected MIME type: %s", body.mime)
		headers.Set(contentTypeHeader, body.mime.String())
	}

	var output io.Writer = handler.infos
	if handler.outputFile != "" {
		file, err := os.Create(handler.outputFile)
		if err != nil {
			return err
		}
		output = file
		defer file.Close()
	}

	for i := 0; i < handler.repeat; i++ {
		req, err := handler.buildRequest(method, u, body.bytes, headers)
		if err != nil {
			return err
		}

		res, err := handler.client.Send(req)
		if err != nil {
			return err
		}

		err = handler.outputResults(res, output)
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler *RequestHandler) buildRequest(
	method string,
	url *url.URL,
	body []byte,
	header http.Header,
) (*http.Request, error) {
	req, err := handler.client.BuildRequest(method, url, body, header)
	if err != nil {
		return nil, err
	}

	err = handler.signer.Sign(req, bytes.NewReader(body))
	return req, err
}

func (handler *RequestHandler) outputResults(r *http.Response, w io.Writer) error {
	b, err := handler.formatter.Format(r)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		_, err = w.Write(b)
		if err != nil {
			return err
		}

		_, err = w.Write(newline)
		if err != nil {
			return err
		}
	}

	doFail := handler.fail && r.StatusCode >= 400
	if doFail {
		handler.logger.Printf("Request failed with status %s", r.Status)
		handler.failFunc(1)
	}

	return nil
}

// Get request headers passed as parameters and defaultHeaders.
// Also sets the User-Agent header if not set by the client.
func (handler *RequestHandler) getHeaders() (http.Header, error) {
	headers := handler.headerOpt.values

	// handler.defaultHeaders must be a string containing headers separated by |
	values := strings.Split(handler.defaultHeaders, "|") // Split by |
	values = util.Map(values, strings.TrimSpace)         // Remove whitespace
	values = util.Filter(values, func(s string) bool {   // Filter empty
		return len(s) > 0
	})
	for _, h := range values {
		key, value, err := parseHeader(h)
		if err != nil {
			return headers, fmt.Errorf("invalid header format in %s: %w", defaultHeadersEnv, err)
		}
		headers.Add(key, value)
	}

	if headers.Get(userAgentHeader) == "" {
		s := fmt.Sprintf("go-http-cli/%s (%s; %s)", handler.version, runtime.GOOS, runtime.GOARCH)
		headers.Set(userAgentHeader, s)
	}

	return headers, nil
}

type requestBody struct {
	bytes []byte
	mime  client.MIMEType
}

func (handler *RequestHandler) getRequestBody(bodyFlag string) (requestBody, error) {
	bodyFlag = strings.TrimSpace(bodyFlag)
	mime := client.MIMETypeUnknown

	if bodyFlag == "" {
		// Not provided via flags, check stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			handler.logger.Print("Reading body from stdin")

			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				return emptyRequestBody, err
			}

			return requestBody{b, mime}, nil
		}

		handler.logger.Print("No body provided")
		return requestBody{nil, mime}, nil
	}

	// We first try to read as a file
	body, err := ioutil.ReadFile(bodyFlag)
	if err != nil && !os.IsNotExist(err) {
		return emptyRequestBody, err
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

	return requestBody{body, mime}, nil
}
