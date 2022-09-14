package history

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/util"
	"github.com/stretchr/testify/require"
)

type fixture struct {
	handler  *fileHandler
	settings client.Settings
	filepath string
}

func (f *fixture) fileExists() bool {
	exists, _, _ := util.FileExists(f.filepath)
	return exists
}

func setupTest(t *testing.T) *fixture {
	filepath := "./history.txt"
	t.Cleanup(func() {
		_ = os.Remove(filepath)
	})
	return &fixture{
		handler:  NewHandler(filepath),
		filepath: filepath,
		settings: client.NewSettings(),
	}
}

func TestAppend(t *testing.T) {
	f := setupTest(t)
	req := newRequest(http.MethodGet, nil)

	_, err := f.handler.Append(req, nil, f.settings)

	require.NoError(t, err)
	require.False(t, f.fileExists())
	require.NotEmpty(t, f.handler.changes)
}

func TestLatestEmpty(t *testing.T) {
	f := setupTest(t)
	_, err := f.handler.Latest()
	require.Error(t, err)
	require.EqualError(t, err, ErrNoHistory.Error())
}

func TestLatest(t *testing.T) {
	// Arrange
	f := setupTest(t)
	req := newRequest(http.MethodGet, nil)
	_, err := f.handler.Append(req, nil, f.settings)

	// Act
	require.NoError(t, err)

	// Assert
	entry, err := f.handler.Latest()
	require.NoError(t, err)
	require.Equal(t, req.Method, entry.Method)
}

func TestWrite(t *testing.T) {
	// Arrange
	f := setupTest(t)

	// Act
	req := newRequest(http.MethodGet, nil)
	_, appendErr := f.handler.Append(req, nil, f.settings)
	writeErr := f.handler.Write()

	// Assert
	require.NoError(t, appendErr)
	require.NoError(t, writeErr)
	require.True(t, f.fileExists())
}

func TestClear(t *testing.T) {
	// Arrange
	f := setupTest(t)
	_, err := f.handler.Append(newRequest(http.MethodGet, nil), nil, f.settings)
	require.NoError(t, err)
	err = f.handler.Write()
	require.NoError(t, err)

	// Act
	err = f.handler.Clear()

	// Assert
	require.NoError(t, err)
	entries, err := f.handler.GetAll()
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestLoadNoHistory(t *testing.T) {
	// Arrange
	f := setupTest(t)

	// Act
	entries, err := f.handler.load()

	// Assert
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestLoad(t *testing.T) {
	// Arrange
	f := setupTest(t)

	// Act
	requests := []struct {
		method   string
		body     []byte
		settings client.Settings
	}{
		{http.MethodGet, nil, client.NewSettings()},
		{http.MethodPost, []byte("test"), client.NewSettings().WithCert("cert", "key")},
	}
	for _, s := range requests {
		r := newRequest(s.method, s.body)
		_, err := f.handler.Append(r, nil, s.settings)
		require.NoError(t, err)
	}
	_ = f.handler.Write()
	entries, err := f.handler.load()

	// Assert
	require.NoError(t, err)
	require.Len(t, entries, len(requests))
	expected := requests[len(requests)-1]
	actual := entries[len(entries)-1]
	require.Equal(t, expected.method, actual.Method)
}

func newRequest(method string, body []byte) *http.Request {
	var b io.Reader
	if body != nil {
		b = bytes.NewReader(body)
	}
	r, _ := http.NewRequest(method, "http://localhost:8080/path", b)
	return r
}
