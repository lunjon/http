package options

import (
	"crypto/tls"
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

	header := NewHeaderOption()
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

func TestPortOption(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		{"1234", false},
		{"12345", false},
		{"30333", false},
		{"80", true},
		{"443", true},
		{"1024", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			opt := NewPortOption()
			err := opt.Set(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("portOption.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDataOptions(t *testing.T) {
	tests := []struct {
		value   DataOptions
		wantErr bool
	}{
		{DataOptions{}, false},
		{DataOptions{dataString: ""}, false},
		{DataOptions{dataFile: ""}, false},
		{DataOptions{dataStdin: true}, false},

		{DataOptions{dataString: "yes", dataFile: "yes"}, true},
		{DataOptions{dataString: "yes", dataStdin: true}, true},
		{DataOptions{dataFile: "yes", dataStdin: true}, true},
		{DataOptions{dataString: "yes", dataFile: "yes", dataStdin: true}, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := tt.value.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("dataOptions.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTLSVersionOption(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		{"1.0", false},
		{"1.1", false},
		{"1.2", false},
		{"1.3", false},
		{"", true},
		{"2.1", true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			opt := NewTLSVersionOption(tls.VersionTLS13)
			err := opt.Set(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLSVersionOptin.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
