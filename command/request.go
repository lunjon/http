package command

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
	"strings"

	"github.com/lunjon/http/client"
)

var (
	newline           = []byte("\n")
	errCertFlags      = errors.New("--cert-pub-file requires --cert-key-file and vice versa")
	errBriefAndSilent = errors.New("cannot specify both --brief and --silent")
	emptyRequestBody  = requestBody{mime: client.MIMETypeUnknown}
)

// RequestHandler handles all commands.
type RequestHandler struct {
	client         *client.Client
	aliasManager   AliasManager
	formatter      Formatter
	signer         client.RequestSigner
	infos          io.Writer
	errors         io.Writer
	logger         *log.Logger
	fail           bool
	failFunc       FailFunc
	repeat         int
	defaultHeaders string
	headerOpt      *HeaderOption
}

func newHandler(
	client *client.Client,
	aliasManager AliasManager,
	formatter Formatter,
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
		return newUserError(err)
	}

	headers, err := handler.getHeaders()
	if err != nil {
		return newUserError(err)
	}

	setContentType := headers.Get("content-type") == "" && body.mime != client.MIMETypeUnknown
	if setContentType {
		handler.logger.Printf("Detected MIME type: %s", body.mime)
		headers.Add("Content-Type", body.mime.String())
	}

	for i := 0; i < handler.repeat; i++ {
		req, err := handler.buildRequest(method, u, body.bytes, headers)
		if err != nil {
			return newUserError(err)
		}

		res, err := handler.client.Send(req)
		if err != nil {
			return err
		}

		err = handler.outputResults(res)
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

func (handler *RequestHandler) outputResults(r *http.Response) error {
	b, err := handler.formatter.Format(r)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		_, err = handler.infos.Write(b)
		if err != nil {
			return err
		}

		_, err = handler.infos.Write(newline)
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

// Get the request headers from the handler header field as well as
// the environment variable for default headers.
func (handler *RequestHandler) getHeaders() (http.Header, error) {
	headers := handler.headerOpt.values
	if handler.defaultHeaders == "" {
		return headers, nil
	}

	// val is a string containing headers separated by a vertical pipe: |
	for _, h := range strings.Split(handler.defaultHeaders, "|") {
		key, value, err := parseHeader(h)
		if err != nil {
			return headers, fmt.Errorf("invalid header format in %s: %w", defaultHeadersEnv, err)
		}
		headers.Add(key, value)
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
