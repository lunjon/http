package client

import (
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

// TODO: use i handler or client
type Signer interface {
	Sign(r *http.Request, body io.ReadSeeker) error
}

type AWSigner struct {
	region string
	sgn    *v4.Signer
}

func NewAWSigner(sgn *v4.Signer, region string) *AWSigner {
	return &AWSigner{
		sgn:    sgn,
		region: region,
	}
}

func (s *AWSigner) Sign(r *http.Request, body io.ReadSeeker) error {
	creds := credentials.NewCredentials(&credentials.EnvProvider{})
	signer := v4.NewSigner(creds)
	_, err := signer.Sign(r, body, "execute-api", s.region, time.Now())
	return err
}
