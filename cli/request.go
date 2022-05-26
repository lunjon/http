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

	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/format"
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
	cfg        config.Config
	client     *client.Client
	headers    http.Header
	aliases    map[string]string
	formatter  format.ResponseFormatter
	signer     client.RequestSigner
	output     io.Writer
	logger     *log.Logger
	failFunc   FailFunc
	headerOpt  *HeaderOption
	outputFile string
}

func newRequestHandler(
	client *client.Client,
	aliases map[string]string,
	formatter format.ResponseFormatter,
	signer client.RequestSigner,
	logger *log.Logger,
	cfg config.Config,
	headers http.Header,
	output io.Writer,
	outputFile string,
	failFunc FailFunc,
) *RequestHandler {
	return &RequestHandler{
		cfg:        cfg,
		client:     client,
		output:     output,
		aliases:    aliases,
		headers:    headers,
		logger:     logger,
		signer:     signer,
		formatter:  formatter,
		failFunc:   failFunc,
		outputFile: outputFile,
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

	u, err := client.ParseURL(url, handler.aliases)
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

	req, err := handler.buildRequest(method, u, body.bytes, headers)
	if err != nil {
		return err
	}

	res, err := handler.client.Send(req)
	if err != nil {
		return err
	}

	err = handler.outputResults(res)
	if err != nil {
		return err
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

func (handler *RequestHandler) outputResults(r *http.Response) error {
	b, err := handler.formatter.Format(r)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		_, err = handler.output.Write(b)
		if err != nil {
			return err
		}

		_, err = handler.output.Write(newline)
		if err != nil {
			return err
		}
	}

	doFail := handler.cfg.Fail && r.StatusCode >= 400
	if doFail {
		handler.logger.Printf("Request failed with status %s", r.Status)
		handler.failFunc(1)
	}

	return nil
}

// Get request headers passed as parameters and defaultHeaders.
// Also sets the User-Agent header if not set by the client.
func (handler *RequestHandler) getHeaders() (http.Header, error) {
	if handler.headers.Get(userAgentHeader) == "" {
		s := fmt.Sprintf("go-http-cli (%s; %s)", runtime.GOOS, runtime.GOARCH)
		handler.headers.Set(userAgentHeader, s)
	}
	return handler.headers, nil
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
