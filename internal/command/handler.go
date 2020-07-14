package command

import (
	"fmt"
	"github.com/lunjon/httpreq/internal/constants"
	"github.com/lunjon/httpreq/internal/rest"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Handler ...
type Handler struct {
	logger *log.Logger
	// Output is where the handler will write the results.
	// It is initialized to os.Stdout as default
	Output io.Writer
}

func NewHandler() *Handler {
	return &Handler{
		logger: log.New(os.Stdout, "", 0),
		Output: os.Stdout,
	}
}

func (handler *Handler) Verbose(v bool) {
	if !v {
		handler.logger.SetOutput(ioutil.Discard)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	h.logger.Printf("Handling request: %s %s", r.Method, r.URL)

	body := fmt.Sprintf(`{"url": "%s", "method": "%s"}`, r.URL, r.Method)
	_, _ = w.Write([]byte(body))
}

func (handler *Handler) Get(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodGet, nil, cmd, args)
}

func (handler *Handler) Post(cmd *cobra.Command, args []string) {
	bodyFlag, _ := cmd.Flags().GetString(constants.BodyFlagName)
	if bodyFlag == "" {
		fmt.Println("No or invalid JSON body specified")
		cmd.Usage()
		os.Exit(2)
	}

	// We first try to read as a file
	body, err := ioutil.ReadFile(bodyFlag)
	if err != nil && !os.IsNotExist(err) {
		log.Print("Failed to open input file: ", err)
		checkError(err, 1, false, cmd)
	}

	if body == nil {
		// Assume that the content was given as string
		log.Print("Assuming body was given as content string")
		body = []byte(bodyFlag)
	}

	handler.handleRequest(http.MethodPost, body, cmd, args)
}

func (handler *Handler) Delete(cmd *cobra.Command, args []string) {
	handler.handleRequest(http.MethodDelete, nil, cmd, args)
}

func (handler *Handler) handleRequest(method string, body []byte, cmd *cobra.Command, args []string) {
	url := args[0]
	headerString, _ := cmd.Flags().GetStringSlice(constants.HeaderFlagName)
	header, err := getHeaders(headerString)
	checkError(err, 2, true, cmd)

	timeout, _ := cmd.Flags().GetDuration(constants.TimeoutFlagName)
	log.Printf("Using timeout: %v", timeout)

	httpClient := &http.Client{
		Timeout: timeout,
	}
	client := rest.NewClient(httpClient, handler.logger)

	req, err := client.BuildRequest(method, url, body, header)
	checkError(err, 2, false, cmd)

	signRequest, _ := cmd.Flags().GetBool(constants.AWSSigV4FlagName)
	if signRequest {
		region, _ := cmd.Flags().GetString(constants.AWSRegionFlagName)
		profile, _ := cmd.Flags().GetString(constants.AWSProfileFlagName)

		err = client.SignRequest(req, nil, region, profile)
		checkError(err, 2, true, cmd)
	}

	res := client.SendRequest(req)
	checkError(res.Error(), 1, false, cmd)
	handler.outputResults(cmd, res)
}

func (handler *Handler) outputResults(cmd *cobra.Command, r *rest.Result) {
	printResponseBodyOnly, _ := cmd.Flags().GetBool(constants.ResponseBodyOnlyFlagName)
	filename, _ := cmd.Flags().GetString(constants.OutputFileFlagName)

	var writeToFile bool
	output := handler.Output

	if filename != "" {
		log.Printf("Writing results to file: %s", filename)
		writeToFile = true
		f, err := os.Create(filename)
		checkError(err, 1, false, cmd)
		output = f
	}

	var body string
	if !printResponseBodyOnly {
		fmt.Fprintln(output, r.Info())
	}

	if writeToFile {
		b, err := r.Body()
		checkError(err, 1, false, cmd)
		body = string(b)
	} else {
		body, _ = r.BodyFormatString()
	}

	_, err := fmt.Fprintln(output, body)
	checkError(err, 1, false, cmd)
}
