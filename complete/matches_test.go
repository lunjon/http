package complete

import "testing"

func TestMatches(t *testing.T) {
	items := []string{"turtle", "Car", "balloon", "Pennywise"}
	matches := Matches("a", items)
	if len(matches) != 2 {
		t.Errorf("was %d expected 2", len(matches))
	}
}
