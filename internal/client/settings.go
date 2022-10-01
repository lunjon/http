package client

import (
	"net/http"
	"time"
)

type checkRedirectFunc func(*http.Request, []*http.Request) error

type Settings struct {
	Timeout         time.Duration
	TLS             TLSOptions
	FollowRedirects bool
}

func NewSettings() Settings {
	return Settings{
		Timeout:         time.Second * 30,
		FollowRedirects: true,
		TLS:             NewTLSOptions(),
	}
}

func (s Settings) WithTLSOptions(opts TLSOptions) Settings {
	s.TLS = opts
	return s
}

func (s Settings) WithTimeout(t time.Duration) Settings {
	s.Timeout = t
	return s
}

func (s Settings) WithNoFollowRedirects(b bool) Settings {
	s.FollowRedirects = !b
	return s
}

func (s Settings) BuildHTTPClient() (*http.Client, error) {
	tlsConfig, err := s.TLS.getTLSConfig()
	if err != nil {
		return nil, err
	}

	var redirect checkRedirectFunc
	if !s.FollowRedirects {
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
