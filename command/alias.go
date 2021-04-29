package command

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	filepath, err := os.UserHomeDir()
	checkErr(err)
	aliasFilepath = path.Join(filepath, ".http", "alias")
}

var aliasFilepath string

func (handler *Handler) handleAlias(_ *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		handler.listAlias()
	case 2:
		handler.setAlias(args[0], args[1])
	default:
		fmt.Fprintln(handler.errors, "unknown number of arguments")
	}
}

func (handler *Handler) listAlias() {
	alias, err := handler.readAliasFile()
	checkErr(err)
	for a, url := range alias {
		fmt.Fprintf(handler.infos, "%s  ->  %s\n", a, url)
	}
}

func (handler *Handler) setAlias(alias, url string) {
	aliases, err := handler.readAliasFile()
	aliases[alias] = url

	file, err := os.OpenFile(aliasFilepath, os.O_WRONLY|os.O_CREATE, 0600)
	checkErr(err)
	defer file.Close()
	for alias, url := range aliases {
		_, err = file.WriteString(fmt.Sprintf("%s %s\n", alias, url))
		checkErr(err)
	}
}

func (handler *Handler) readAliasFile() (map[string]string, error) {
	alias := make(map[string]string)
	file, err := os.Open(aliasFilepath)
	if os.IsNotExist(err) {
		return alias, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		s := strings.Split(line, " ")
		if len(s) != 2 {
			continue
		}
		alias[s[0]] = s[1]
	}
	return alias, nil
}
