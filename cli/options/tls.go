package options

import (
	"crypto/tls"
	"fmt"
)

var tlsVersions map[string]uint16 = map[string]uint16{
	"1.0": tls.VersionTLS10,
	"1.1": tls.VersionTLS11,
	"1.2": tls.VersionTLS12,
	"1.3": tls.VersionTLS13,
}

type TLSVersionOption struct {
	value uint16
}

func NewTLSVersionOption(initial uint16) *TLSVersionOption {
	return &TLSVersionOption{
		value: initial,
	}
}

func (h *TLSVersionOption) Value() uint16 {
	return h.value
}

func (h *TLSVersionOption) Set(s string) error {
	version, found := tlsVersions[s]
	if !found {
		return fmt.Errorf("invalid TLS version: %s", s)
	}

	h.value = version
	return nil
}

func (h *TLSVersionOption) Type() string {
	return "TLSVersion"
}

func (h *TLSVersionOption) String() string {
	for s, v := range tlsVersions {
		if h.value == v {
			return s
		}
	}
	return fmt.Sprint(h.value)
}
