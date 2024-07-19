package stores

import (
	"errors"
	"log"

	"github.com/jonasroussel/hyve/tools"
)

var Active Store

var ErrNotFound = errors.New("certificate not found")

type Store interface {
	Load() error
	AddCertificate(domain string, cert Certificate) error
	GetCertificate(domain string) (*Certificate, error)
	GetAllCertificates(exp int64) []Certificate
	UpdateCertificate(domain string, cert Certificate) error
	RemoveCertificate(domain string) error
}

type Certificate struct {
	Domain          string `json:"domain"`
	CertificateData []byte
	PrivateKeyData  []byte
	Issuer          string `json:"issuer"`
	ExpiresAt       int64  `json:"expires_at"`
	CreatedAt       int64  `json:"created_at"`
}

func Load() {
	switch tools.Env.StoreType {
	case "sql":
		Active = NewSQLStore()
	case "file":
		Active = NewFileStore()
	default:
		log.Fatal("STORE_TYPE not supported")
	}
}
