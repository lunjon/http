package cli

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/lunjon/http/internal/client"
	"github.com/spf13/cobra"
)

const (
	headerFlagName                = "header"
	awsSigV4FlagName              = "aws-sigv4"
	awsRegionFlagName             = "aws-region"
	dataStringFlagName            = "data"
	dataStdinFlagName             = "data-stdin"
	dataFileFlagName              = "data-file"
	displayFlagName               = "display"
	failFlagName                  = "fail"
	detailsFlagName               = "details"
	timeoutFlagName               = "timeout"
	verboseFlagName               = "verbose"
	certfileFlagName              = "cert"
	certkeyFlagName               = "key"
	certKindFlagName              = "cert-kind"
	outfileFlagName               = "outfile"
	noFollowRedirectsFlagName     = "no-follow-redirects"
	aliasHeadingFlagName          = "no-heading"
	historyNumberFlagName         = "num"
	tlsTraceFlagName              = "tls-trace"
	tlsMinVersionFlagName         = "tls-min-version"
	tlsMaxVersionFlagName         = "tls-max-version"
	tlsInsecureSkipVerifyFlagName = "tls-skip-verify-insecure"
)

var (
	headerReg = regexp.MustCompile(`^([a-zA-Z0-9\-_]+)\s*[:=]\s*(\S[\s\S]*)+$`)
)

type HeaderOption struct {
	values http.Header
}

func newHeaderOption() *HeaderOption {
	return &HeaderOption{
		values: make(http.Header),
	}
}

func (h *HeaderOption) Header() http.Header {
	return h.values
}

// Append adds the provided value as a header if it is valid
func (h *HeaderOption) Set(s string) error {
	key, value, err := parseHeader(s)
	if err != nil {
		return err
	}
	h.values.Add(key, value)
	return nil
}

func (h *HeaderOption) Type() string {
	return "Header"
}

func (h *HeaderOption) String() string {
	return ""
}

// Parse string s into a header name and value.
func parseHeader(h string) (string, string, error) {
	h = strings.TrimSpace(h)
	if len(h) == 0 {
		return "", "", fmt.Errorf("empty header")
	}

	match := headerReg.FindAllStringSubmatch(h, -1)
	if match == nil {
		return "", "", fmt.Errorf("invalid header format: %s", h)
	}

	key := strings.TrimSpace(match[0][1])
	value := strings.TrimSpace(match[0][2])
	return key, value, nil
}

type portOption struct {
	port uint
}

func newPortOption() *portOption {
	return &portOption{
		port: 8080,
	}
}

func (o *portOption) value() uint {
	return o.port
}

func (o *portOption) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid number or outside")
	}

	if v > math.MaxUint16 {
		return fmt.Errorf("outside valid range for port")
	}
	if v <= 1024 {
		return fmt.Errorf("reserved port number")
	}

	o.port = uint(v)
	return nil
}

func (h *portOption) Type() string {
	return "Port"
}

func (h *portOption) String() string {
	return ""
}

// Container for the data flags/options.
// Every field is mutually exclusive.
type dataOptions struct {
	dataString string
	dataFile   string
	dataStdin  bool
}

func dataOptionsFromFlags(cmd *cobra.Command) (dataOptions, error) {
	flags := cmd.Flags()
	dataString, _ := flags.GetString(dataStringFlagName)
	dataFile, _ := flags.GetString(dataFileFlagName)
	dataStdin, _ := flags.GetBool(dataStdinFlagName)

	opts := dataOptions{
		dataString: dataString,
		dataFile:   dataFile,
		dataStdin:  dataStdin,
	}
	return opts, opts.validate()
}

func (opts dataOptions) validate() error {
	invalid := (opts.dataString != "" && opts.dataFile != "") || (opts.dataString != "" && opts.dataStdin) || (opts.dataFile != "" && opts.dataStdin)
	if invalid {
		return fmt.Errorf("invalid combination of --data* options: must only specify one")
	}
	return nil
}

func (opts dataOptions) getRequestBody() (requestBody, error) {
	if err := opts.validate(); err != nil {
		return emptyRequestBody, err
	}

	mime := client.MIMETypeUnknown
	if opts.dataString != "" {
		return requestBody{[]byte(opts.dataString), mime}, nil
	} else if opts.dataFile != "" {
		body, err := os.ReadFile(opts.dataFile)
		if err != nil {
			return emptyRequestBody, err
		}

		// Try detecting filetype in order to set MIME type
		switch path.Ext(opts.dataFile) {
		case ".html":
			mime = client.MIMETypeHTML
		case ".csv":
			mime = client.MIMETypeCSV
		case ".json":
			mime = client.MIMETypeJSON
		case ".xml":
			mime = client.MIMETypeXML
		}
		return requestBody{body, mime}, nil
	} else if opts.dataStdin {
		b, err := io.ReadAll(os.Stdin)
		return requestBody{b, mime}, err
	}

	return emptyRequestBody, nil
}
