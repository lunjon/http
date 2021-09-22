package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
	checkErr(err)
	aliases[alias] = url
	handler.writeAliasFile(aliases)
}

func (handler *Handler) readAliasFile() (map[string]string, error) {
	alias := make(map[string]string)
	file, err := os.Open(handler.aliasFilePath)
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

func (handler *Handler) writeAliasFile(aliases map[string]string) {
	if _, err := os.Stat(handler.gohttpDir); os.IsNotExist(err) {
		err := os.MkdirAll(handler.gohttpDir, 0700)
		checkErr(err)
	}

	file, err := os.OpenFile(handler.aliasFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	checkErr(err)
	defer file.Close()
	for alias, url := range aliases {
		_, err = file.WriteString(fmt.Sprintf("%s %s\n", alias, url))
		checkErr(err)
	}
}
