package client

import (
	"crypto/tls"
	"encoding/pem"

	"github.com/lunjon/http/internal/types"
	"golang.org/x/crypto/pkcs12"
)

type TLSOptions struct {
	Cert               types.Option[certificate]
	MinVersion         uint16
	MaxVersion         uint16
	SkipVerifyInsecure bool
}

func NewTLSOptions() TLSOptions {
	return TLSOptions{
		Cert:               types.Option[certificate]{},
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
		SkipVerifyInsecure: false,
	}
}

func (tlsOptions TLSOptions) WithVersions(min, max uint16) TLSOptions {
	tlsOptions.MinVersion = min
	tlsOptions.MaxVersion = max
	return tlsOptions
}

func (tlsOptions TLSOptions) WithX509Cert(certfile, keyfile string) TLSOptions {
	tlsOptions.Cert = tlsOptions.Cert.Set(x509Cert{
		certfile: certfile,
		keyfile:  keyfile,
	})
	return tlsOptions
}

func (tlsOptions TLSOptions) WithPKCS12Cert(pfx []byte, passw string) TLSOptions {
	tlsOptions.Cert = tlsOptions.Cert.Set(pkcs12Cert{
		pfx:   pfx,
		passw: passw,
	})
	return tlsOptions
}

func (tlsOptions TLSOptions) getTLSConfig() (tls.Config, error) {
	certs := []tls.Certificate{}
	if tlsOptions.Cert.IsSome() {
		certOpt := tlsOptions.Cert.Value()
		crt, err := certOpt.Load()
		if err != nil {
			return tls.Config{}, err
		}
		certs = append(certs, crt)
	}

	return tls.Config{
		Certificates:       certs,
		InsecureSkipVerify: tlsOptions.SkipVerifyInsecure,
		MinVersion:         tlsOptions.MinVersion,
		MaxVersion:         tlsOptions.MaxVersion,
	}, nil
}

type certificate interface {
	Load() (tls.Certificate, error)
}

type x509Cert struct {
	certfile string
	keyfile  string
}

func (cert x509Cert) Load() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(cert.certfile, cert.keyfile)
}

type pkcs12Cert struct {
	pfx   []byte
	passw string
}

func (cert pkcs12Cert) Load() (tls.Certificate, error) {
	blocks, err := pkcs12.ToPEM(cert.pfx, cert.passw)
	if err != nil {
		return tls.Certificate{}, err
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	return tls.X509KeyPair(pemData, pemData)
}
