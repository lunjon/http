package options

import (
	"crypto/tls"
	"fmt"
	"path"
	"strings"

	"github.com/lunjon/http/internal/util"
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

type FileOption struct {
	file string
}

func (h *FileOption) Type() string   { return "File" }
func (h *FileOption) String() string { return "" }

func (h *FileOption) Value() (string, bool) {
	return h.file, h.file != ""
}

func (h *FileOption) Set(s string) error {
	exists, isdir, err := util.FileExists(s)
	if err != nil {
		return err
	}

	if !exists || isdir {
		return fmt.Errorf("file not found: %s", s)
	}
	h.file = s
	return nil
}

type CertKind string

const (
	CertKindX509   = "x509"
	CertKindPKCS12 = "PKCS#12"
)

type CertKindOption struct {
	kind string
}

func (h *CertKindOption) Type() string   { return "CertKind" }
func (h *CertKindOption) String() string { return CertKindX509 }
func (h *CertKindOption) Update(certFile string) {
	if h.kind != "" {
		return
	}
	if path.Ext(certFile) == ".pfx" {
		h.kind = CertKindPKCS12
	}
}

func (h *CertKindOption) Value() string {
	if h.kind == "" {
		return CertKindX509
	}
	return h.kind
}

func (h *CertKindOption) Set(s string) error {
	switch strings.ToLower(s) {
	case "x509":
		h.kind = CertKindX509
	case "pkcs12", "pkcs#12":
		h.kind = CertKindPKCS12
	default:
		return fmt.Errorf("invalid cert kind: %s", s)
	}
	return nil
}
