package cli

import (
	"testing"
)

func TestHeaderOption(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		// Valid
		{"Name: value", false},
		{"Name = value", false},
		{"Name : 12-122?!=!92", false},
		{"Name : Bearer 1234-abcd", false},

		// Invalid
		{"", true},
		{"\n", true},
		{"A B", true},
		{"Name: ", true},
		{": value", true},
	}

	header := newHeaderOption()
	for _, tt := range tests {
		t.Run("Parse header: "+tt.value, func(t *testing.T) {
			err := header.Set(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
