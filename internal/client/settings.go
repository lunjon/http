package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/lunjon/http/internal/types"
)

type checkRedirectFunc func(*http.Request, []*http.Request) error

type certPair struct {
	certFile string
	keyFile  string
}

type Settings struct {
	timeout         time.Duration
	cert            types.Option[certPair]
	followRedirects bool
}

func NewSettings() Settings {
	return Settings{
		timeout: time.Second * 30,
		cert:    types.Option[certPair]{},
	}
}

func (s Settings) WithTimeout(t time.Duration) Settings {
	s.timeout = t
	return s
}

func (s Settings) WithCert(certFile, keyFile string) Settings {
	s.cert = s.cert.Set(certPair{certFile, keyFile})
	return s
}

func (s Settings) WithNoFollowRedirects(b bool) Settings {
	s.followRedirects = !b
	return s
}

func (s Settings) buildHTTPClient() (*http.Client, error) {
	var tlsConfig tls.Config
	if s.cert.IsSome() {
		pair := s.cert.Value()
		cert, err := tls.LoadX509KeyPair(pair.certFile, pair.keyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig = tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	var redirect checkRedirectFunc
	if s.followRedirects {
		redirect = func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &http.Client{
		Timeout:       s.timeout,
		CheckRedirect: redirect,
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tlsConfig,
		},
	}, nil
}
