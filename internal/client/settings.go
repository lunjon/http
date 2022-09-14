package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/lunjon/http/internal/types"
)

type checkRedirectFunc func(*http.Request, []*http.Request) error

type CertPair struct {
	CertFile string
	KeyFile  string
}

type Settings struct {
	Timeout         time.Duration          `json:"timeout"`
	Cert            types.Option[CertPair] `json:"cert"`
	FollowRedirects bool                   `json:"followRedirects"`
}

func NewSettings() Settings {
	return Settings{
		Timeout: time.Second * 30,
		Cert:    types.Option[CertPair]{},
	}
}

func (s Settings) WithTimeout(t time.Duration) Settings {
	s.Timeout = t
	return s
}

func (s Settings) WithCert(certFile, keyFile string) Settings {
	s.Cert = s.Cert.Set(CertPair{certFile, keyFile})
	return s
}

func (s Settings) WithNoFollowRedirects(b bool) Settings {
	s.FollowRedirects = !b
	return s
}

func (s Settings) BuildHTTPClient() (*http.Client, error) {
	var tlsConfig tls.Config
	if s.Cert.IsSome() {
		pair := s.Cert.Value()
		cert, err := tls.LoadX509KeyPair(pair.CertFile, pair.KeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig = tls.Config{
			Certificates: []tls.Certificate{cert},
		}
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
