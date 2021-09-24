package command

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/lunjon/http/client"
	"github.com/lunjon/http/util"
)

var (
	aliasPattern = regexp.MustCompile(`^[a-zA-Z_]\w{0,19}$`)
)

type AliasHandler struct {
	aliasFilepath string
	// Output of infos
	infos io.Writer
	// Output of errors
	errors io.Writer
}

func (handler *AliasHandler) listAlias() error {
	alias, err := readAliasFile(handler.aliasFilepath)
	if err != nil {
		return err
	}

	// Sort by name
	names := []string{}
	for name := range alias {
		names = append(names, name)
	}
	sort.Strings(names)

	taber := util.NewTaber("")
	for _, name := range names {
		taber.WriteLine(name+":", alias[name])
	}
	fmt.Fprintln(handler.infos, taber.String())
	return nil
}

func (handler *AliasHandler) setAlias(alias, url string) error {
	if !aliasPattern.MatchString(alias) {
		return fmt.Errorf("invalid alias name: %s", alias)
	}

	u, err := client.ParseURL(url, nil)
	if err != nil {
		return fmt.Errorf("invalid alias URL: %s", url)
	}

	aliases, err := readAliasFile(handler.aliasFilepath)
	if err != nil {
		return err
	}

	aliases[alias] = u.String()
	return writeAliasFile(handler.aliasFilepath, aliases)
}

func readAliasFile(filepath string) (map[string]string, error) {
	alias := make(map[string]string)
	file, err := os.Open(filepath)
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

func writeAliasFile(filepath string, aliases map[string]string) error {
	dir := path.Dir(filepath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()
	for alias, url := range aliases {
		_, err = file.WriteString(fmt.Sprintf("%s %s\n", alias, url))
		if err != nil {
			return err
		}
	}

	return nil
}
