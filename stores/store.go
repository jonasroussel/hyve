package stores

import (
	"errors"
)

var ErrCertNotFound = errors.New("certificate not found")

type Store interface {
	Load() error
	AddCertificate(domain string, cert Certificate) error
	GetCertificate(domain string) (*Certificate, error)
	GetAllCertificates(exp int64) []Certificate
	UpdateCertificate(domain string, cert Certificate) error
	RemoveCertificate(domain string) error
	Close() error
}

type Certificate struct {
	Domain          string `json:"domain" bson:"domain"`
	CertificateData []byte `json:"-" bson:"certificate"`
	PrivateKeyData  []byte `json:"-" bson:"private_key"`
	Issuer          string `json:"issuer" bson:"issuer"`
	ExpiresAt       int64  `json:"expires_at" bson:"expires_at"`
	CreatedAt       int64  `json:"created_at" bson:"created_at"`
}
