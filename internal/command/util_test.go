package command

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGetHeaders(t *testing.T) {
	tests := []struct {
		args    []string
		want    http.Header
		wantErr bool
	}{
		{[]string{}, nil, false},
		{nil, nil, false},
		{[]string{"a=b"}, http.Header{"A": []string{"b"}}, false},
		{[]string{"a:b"}, http.Header{"A": []string{"b"}}, false},
		{[]string{"a:b", "c=d"}, http.Header{"A": []string{"b"}, "C": []string{"d"}}, false},
		{[]string{"lol"}, nil, true},
	}
	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			got, err := getHeaders(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
