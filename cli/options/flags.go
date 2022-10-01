package options

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
	"github.com/lunjon/http/internal/types"
	"github.com/spf13/cobra"
)

const (
	HeaderFlagName                = "header"
	AWSSigV4FlagName              = "aws-sigv4"
	AWSRegionFlagName             = "aws-region"
	DataStringFlagName            = "data"
	DataStdinFlagName             = "data-stdin"
	DataFileFlagName              = "data-file"
	DisplayFlagName               = "display"
	FailFlagName                  = "fail"
	DetailsFlagName               = "details"
	TimeoutFlagName               = "timeout"
	VerboseFlagName               = "verbose"
	OutfileFlagName               = "outfile"
	NoFollowRedirectsFlagName     = "no-follow-redirects"
	AliasHeadingFlagName          = "no-heading"
	CertfileFlagName              = "cert"
	CertkeyFlagName               = "key"
	CertKindFlagName              = "cert-kind"
	TLSTraceFlagName              = "tls-trace"
	TLSMinVersionFlagName         = "tls-min-version"
	TLSMaxVersionFlagName         = "tls-max-version"
	TLSInsecureSkipVerifyFlagName = "tls-skip-verify-insecure"
)

var (
	headerReg = regexp.MustCompile(`^([a-zA-Z0-9\-_]+)\s*[:=]\s*(\S[\s\S]*)+$`)
)

type HeaderOption struct {
	values http.Header
}

func NewHeaderOption() *HeaderOption {
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

func NewPortOption() *portOption {
	return &portOption{
		port: 8080,
	}
}

func (o *portOption) Value() uint {
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
type DataOptions struct {
	dataString string
	dataFile   string
	dataStdin  bool
}

func NewDataOptions(dataString, dataFile string, dataStdin bool) DataOptions {
	return DataOptions{
		dataString: dataString,
		dataFile:   dataFile,
		dataStdin:  dataStdin,
	}
}

func DataOptionsFromFlags(cmd *cobra.Command) (DataOptions, error) {
	flags := cmd.Flags()
	dataString, _ := flags.GetString(DataStringFlagName)
	dataFile, _ := flags.GetString(DataFileFlagName)
	dataStdin, _ := flags.GetBool(DataStdinFlagName)

	opts := DataOptions{
		dataString: dataString,
		dataFile:   dataFile,
		dataStdin:  dataStdin,
	}
	return opts, opts.validate()
}

func (opts DataOptions) validate() error {
	invalid := (opts.dataString != "" && opts.dataFile != "") || (opts.dataString != "" && opts.dataStdin) || (opts.dataFile != "" && opts.dataStdin)
	if invalid {
		return fmt.Errorf("invalid combination of --data* options: must only specify one")
	}
	return nil
}

func (opts DataOptions) GetData() (types.Option[[]byte], client.MIMEType, error) {
	body := types.Option[[]byte]{}
	mime := client.MIMETypeUnknown
	if err := opts.validate(); err != nil {
		return body, mime, err
	}

	if opts.dataString != "" {
		return body.Set([]byte(opts.dataString)), mime, nil
	} else if opts.dataFile != "" {
		b, err := os.ReadFile(opts.dataFile)
		if err != nil {
			return body, mime, err
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
		return body.Set(b), mime, nil
	} else if opts.dataStdin {
		b, err := io.ReadAll(os.Stdin)
		return body.Set(b), mime, err
	}

	return body, mime, nil
}
