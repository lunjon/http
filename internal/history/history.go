package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/util"
)

var ErrNoHistory = errors.New("no history")

type Handler interface {
	GetAll() ([]Entry, error)
	GetByIndex(index uint16) (Entry, error)
	Append(*http.Request, []byte, client.Settings) (Entry, error)
	Latest() (Entry, error)
	Write() error
	Clear() error
}

// real implementation of Handler.
type fileHandler struct {
	filepath string
	changes  []Entry
	// History is loaded lazily
	history []Entry
}

func NewHandler(filepath string) *fileHandler {
	return &fileHandler{
		filepath: filepath,
	}
}

func (h *fileHandler) GetAll() ([]Entry, error) {
	return h.load()
}

func (h *fileHandler) GetByIndex(i uint16) (Entry, error) {
	entries, err := h.load()
	if err != nil {
		return Entry{}, err
	}
	if int(i) > len(entries)-1 {
		return Entry{}, fmt.Errorf("invalid history index")
	}

	if len(entries) == 0 {
		return Entry{}, ErrNoHistory
	}
	return entries[i], nil
}

func (h *fileHandler) Append(
	req *http.Request,
	body []byte,
	settings client.Settings,
) (Entry, error) {
	entry := Entry{
		Timestamp:      time.Now(),
		Method:         req.Method,
		URL:            req.URL.String(),
		Header:         req.Header,
		Body:           body,
		ClientSettings: settings,
	}

	h.changes = append(h.changes, entry)
	return entry, nil
}

func (h *fileHandler) Latest() (Entry, error) {
	if len(h.changes) > 0 {
		return h.changes[len(h.changes)-1], nil
	}

	hist, err := h.load()
	if err != nil {
		return Entry{}, nil
	}

	if len(hist) == 0 {
		return Entry{}, ErrNoHistory
	}
	return hist[len(hist)-1], nil
}

func (h *fileHandler) Clear() error {
	f, err := os.Create(h.filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	h.history = []Entry{}
	h.changes = []Entry{}
	return nil
}

func (h *fileHandler) Write() error {
	if len(h.changes) == 0 {
		return nil
	}

	f, err := os.OpenFile(h.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, entry := range h.changes {
		b, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}

		_, err = f.WriteString("\n")
		if err != nil {
			return err
		}
	}
	return nil
}

// Loads the history file if not already loaded.
// Sorts the requests from newest to oldest.
func (h *fileHandler) load() ([]Entry, error) {
	if h.history != nil {
		return h.history, nil
	}

	exists, isdir, err := util.FileExists(h.filepath)
	if err != nil {
		return nil, err
	}

	if exists && isdir {
		return nil, fmt.Errorf("expected file but was directory: %s", h.filepath)
	}

	if !exists {
		h.history = []Entry{}
		return h.history, nil
	}

	b, err := os.ReadFile(h.filepath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(string(b), "\r\n", "\n"), "\n")

	// Each line contains an entry serialized to a JSON object
	entries := []Entry{}
	for _, line := range util.Filter(lines, stringNotEmpty) {
		var e Entry
		err = json.Unmarshal([]byte(line), &e)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	h.history = entries
	return h.history, nil
}

func stringNotEmpty(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}
