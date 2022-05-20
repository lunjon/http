package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/lunjon/http/client"
	"github.com/lunjon/http/format"
	"github.com/lunjon/http/util"
)

var (
	aliasPattern = regexp.MustCompile(`^[a-zA-Z_]\w{0,19}$`)
)

type AliasManager interface {
	Load() (map[string]string, error)
	Save(map[string]string) error
}

// An implementation of AliasLoader that loads from the configured filepath.
type fileAliasManager struct {
	aliases  map[string]string
	filepath string
}

func newAliasLoader(filepath string) *fileAliasManager {
	return &fileAliasManager{
		filepath: filepath,
	}
}

type AliasHandler struct {
	manager AliasManager
	// Output of infos
	infos io.Writer
	// Output of errors
	errors io.Writer
	styler *format.Styler
}

func NewAliasHandler(m AliasManager, styler *format.Styler, infos, errors io.Writer) *AliasHandler {
	return &AliasHandler{
		manager: m,
		infos:   infos,
		errors:  errors,
	}
}

func (handler *AliasHandler) listAlias(noHeading bool) error {
	aliases, err := handler.manager.Load()
	if err != nil {
		return err
	}

	// Sort by name
	names := []string{}
	for name := range aliases {
		names = append(names, name)
	}
	sort.Strings(names)

	taber := util.NewTaber("")

	if !noHeading {
		taber.WriteLine(handler.styler.WhiteB("Name\t"), handler.styler.WhiteB("URL"))
	}

	for _, name := range names {
		taber.WriteLine(name, aliases[name])
	}
	fmt.Fprintln(handler.infos, taber.String())

	return nil
}

func (handler *AliasHandler) removeAlias(name string) error {
	name = strings.TrimSpace(name)
	if !aliasPattern.MatchString(name) {
		return fmt.Errorf("impossible alias name: %s", name)
	}

	aliases, err := handler.manager.Load()
	if err != nil {
		return err
	}

	_, found := aliases[name]
	if !found {
		return fmt.Errorf("unknown alias: %s\n", name)
	}

	delete(aliases, name)
	err = handler.manager.Save(aliases)
	if err != nil {
		return err
	}
	fmt.Fprintf(handler.infos, "Removed alias %s\n", name)
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

	aliases, err := handler.manager.Load()
	if err != nil {
		return err
	}

	aliases[alias] = u.String()
	return handler.manager.Save(aliases)
}

func (f *fileAliasManager) Load() (map[string]string, error) {
	if f.aliases != nil {
		return f.aliases, nil
	}

	var alias map[string]string
	file, err := os.Open(f.filepath)
	if os.IsNotExist(err) {
		return make(map[string]string), nil
	}

	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &alias)
	return alias, err
}

func (f *fileAliasManager) Save(aliases map[string]string) error {
	dir := path.Dir(f.filepath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return err
		}
	}

	// Write to temporary file
	tmp := fmt.Sprintf("%s.new", f.filepath)
	file, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.MarshalIndent(aliases, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	// Move to real name
	return os.Rename(tmp, f.filepath)
}
