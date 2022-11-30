package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"

	"github.com/lunjon/http/cli/options"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/history"
	"github.com/lunjon/http/internal/types"
)

var (
	newline = []byte("\r\n")
)

const (
	userAgentHeader     = "User-Agent"
	contentTypeHeader   = "Content-Type"
	contentLengthHeader = "Content-Length"
)

// RequestHandler handles all commands.
type RequestHandler struct {
	cfg            config.Config
	client         *client.Client
	headers        http.Header
	formatter      Formatter
	signer         client.RequestSigner
	historyHandler history.Handler
	output         io.Writer
	logger         *log.Logger
	failFunc       FailFunc
	outputFile     types.Option[string]
}

func newRequestHandler(
	client *client.Client,
	formatter Formatter,
	signer client.RequestSigner,
	historyHandler history.Handler,
	logger *log.Logger,
	cfg config.Config,
	headers http.Header,
	output io.Writer,
	outputFile string,
	failFunc FailFunc,
) *RequestHandler {
	outfile := types.Option[string]{}
	if outputFile != "" {
		outfile = outfile.Set(outputFile)
	}

	return &RequestHandler{
		cfg:            cfg,
		client:         client,
		headers:        headers,
		formatter:      formatter,
		signer:         signer,
		historyHandler: historyHandler,
		output:         output,
		logger:         logger,
		failFunc:       failFunc,
		outputFile:     outfile,
	}
}

func (handler *RequestHandler) handleRequest(method, url string, dataOptions options.DataOptions) error {
	headers, err := handler.getHeaders()
	if err != nil {
		return err
	}

	var body []byte
	data, mime, err := dataOptions.GetData()
	if err != nil {
		return err
	}

	if data.IsSome() {
		body = data.MustGet()

		setContentType := headers.Get(contentTypeHeader) == "" && mime != client.MIMETypeUnknown
		if setContentType {
			handler.logger.Printf("Detected MIME type: %s", mime)
			headers.Set(contentTypeHeader, mime.String())
		}

		setContentLength := headers.Get(contentLengthHeader) == "" && len(body) > 0
		if setContentLength {
			handler.logger.Printf("Adding %s header", contentLengthHeader)
			headers.Set(contentLengthHeader, fmt.Sprint(len(body)))
		}
	}

	u, err := client.ParseURL(url, handler.cfg.Aliases)
	if err != nil {
		return err
	}

	req, err := handler.buildRequest(method, u, body, headers)
	if err != nil {
		return err
	}

	// Add to history
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := handler.historyHandler.Append(req, body)
		if err != nil {
			handler.logger.Printf("Error building history entry: %s", err)
			return
		}

		err = handler.historyHandler.Write()
		if err != nil {
			handler.logger.Printf("Error writing history file: %s", err)
		}
	}()

	res, err := handler.client.Send(req)
	if err != nil {
		return err
	}

	err = handler.outputResults(res)
	if err != nil {
		return err
	}

	// Wait for history to be written
	wg.Wait()
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
	b, err := handler.formatter.FormatResponse(r)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		var err error
		if filepath, ok := handler.outputFile.Get(); ok {
			err = os.WriteFile(filepath, b, 0644)
		} else {
			_, err = handler.output.Write(b)
			_, _ = handler.output.Write(newline)
		}

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
