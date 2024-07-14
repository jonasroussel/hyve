package stores

import "errors"

var ErrNotFound = errors.New("certificate not found")

type Store interface {
	Load() error
	AddCertificate(sni string, cert Certificate) error
	GetCertificate(sni string) (*Certificate, error)
	UpdateCertificate(sni string, cert Certificate) error
	RemoveCertificate(sni string) error
}

type Certificate struct {
	CertificateData []byte
	PrivateKeyData  []byte
	Issuer          string `json:"issuer"`
	ExpiresAt       int64  `json:"expires_at"`
	CreatedAt       int64  `json:"created_at"`
}
