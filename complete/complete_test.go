package complete

import "testing"

func TestComplete(t *testing.T) {
	items := []string{"http://local", "https://example.com", "balloon", "https://golang.org"}
	prefix, matches := Complete("http", items)

	if prefix != "http" {
		t.Errorf("was %s expected %s", prefix, "http")
	}

	if len(matches) != 3 {
		t.Errorf("was %d expected 3", len(matches))
	}
}
