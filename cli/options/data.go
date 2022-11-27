package options

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/types"
	"github.com/lunjon/http/internal/util"
	"github.com/spf13/cobra"
)

// Container for the data flags/options.
// Every field is mutually exclusive.
type DataOptions struct {
	dataString     string
	dataFile       string
	dataStdin      bool
	dataURLEncoded []string
}

func NewDataOptions(dataString, dataFile string, dataStdin bool, urlEncoded []string) DataOptions {
	return DataOptions{
		dataString:     dataString,
		dataFile:       dataFile,
		dataStdin:      dataStdin,
		dataURLEncoded: urlEncoded,
	}
}

func DataOptionsFromFlags(cmd *cobra.Command) (DataOptions, error) {
	flags := cmd.Flags()
	dataString, _ := flags.GetString(DataStringFlagName)
	dataFile, _ := flags.GetString(DataFileFlagName)
	dataStdin, _ := flags.GetBool(DataStdinFlagName)
	dataURLEncoded, _ := flags.GetStringArray(DataURLEncodeFlagName)

	opts := DataOptions{
		dataString:     dataString,
		dataFile:       dataFile,
		dataStdin:      dataStdin,
		dataURLEncoded: dataURLEncoded,
	}
	return opts, nil
}

func (opts DataOptions) GetData() (types.Option[[]byte], client.MIMEType, error) {
	body := types.Option[[]byte]{}
	mime := client.MIMETypeUnknown

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
	} else if len(opts.dataURLEncoded) > 0 {
		splitter := util.NewSplitter("=")
		m, err := splitter.ParseMany(opts.dataURLEncoded)
		if err != nil {
			return body, mime, err
		}

		values := []string{}
		for k, v := range m {
			values = append(values, fmt.Sprintf("%s=%s", k, v))
		}

		b := strings.Join(values, "&")
		return body.Set([]byte(b)), client.MIMETypeFormURLEncoded, nil
	}

	return body, mime, nil
}
