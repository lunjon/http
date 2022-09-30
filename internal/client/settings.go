package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/lunjon/http/internal/types"
)

type checkRedirectFunc func(*http.Request, []*http.Request) error

type CertKind int

const (
	CertKindPEM CertKind = iota
	CertKindPFX CertKind = iota
)

type CertOptions struct {
	CertFile string
	Kind     CertKind
	KeyFile  string
}

func CertOptionsFrom(certfile, keyfile string) (types.Option[CertOptions], error) {
	opt := types.Option[CertOptions]{}
	if certfile == "" && keyfile == "" {
		return opt, nil
	}
	opts := CertOptions{
		CertFile: certfile,
		KeyFile:  keyfile,
	}
	return opt.Set(opts), opts.Validate()
}

func (opts CertOptions) getTLSConfig() (tls.Config, error) {
	var err error
	if err = opts.Validate(); err != nil {
		return tls.Config{}, err
	}

	var cert tls.Certificate
	switch opts.Kind {
	case CertKindPEM:
		cert, err = tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return tls.Config{}, err
		}
	case CertKindPFX:
		return tls.Config{}, fmt.Errorf("not implemented support for PFX yet")
	}

	return tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func (opts CertOptions) Validate() error {
	invalid := (opts.Kind == CertKindPEM && opts.KeyFile == "") || (opts.Kind == CertKindPFX && opts.KeyFile != "") || (opts.CertFile == "" && opts.KeyFile == "") || (opts.CertFile == "" && opts.KeyFile != "")
	if invalid {
		return fmt.Errorf("invalid combination of certificate options")
	}

	return nil
}

type Settings struct {
	Timeout         time.Duration             `json:"timeout"`
	Cert            types.Option[CertOptions] `json:"cert"`
	FollowRedirects bool                      `json:"followRedirects"`
}

func NewSettings() Settings {
	return Settings{
		Timeout: time.Second * 30,
		Cert:    types.Option[CertOptions]{},
	}
}

func (s Settings) WithTimeout(t time.Duration) Settings {
	s.Timeout = t
	return s
}

func (s Settings) WithCertPEM(certFile, keyFile string) Settings {
	s.Cert = s.Cert.Set(CertOptions{certFile, CertKindPEM, keyFile})
	return s
}

func (s Settings) WithCert(cert types.Option[CertOptions]) Settings {
	s.Cert = cert
	return s
}

func (s Settings) WithNoFollowRedirects(b bool) Settings {
	s.FollowRedirects = !b
	return s
}

func (s Settings) BuildHTTPClient() (*http.Client, error) {
	var tlsConfig tls.Config
	if s.Cert.IsSome() {
		info := s.Cert.Value()
		cfg, err := info.getTLSConfig()
		if err != nil {
			return nil, err
		}

		tlsConfig = cfg
	}

	var redirect checkRedirectFunc
	if s.FollowRedirects {
		redirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &http.Client{
		Timeout:       s.Timeout,
		CheckRedirect: redirect,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tlsConfig,
		},
	}, nil
}
