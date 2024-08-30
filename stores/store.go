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
	Domain          string `json:"domain" bson:"domain"`
	CertificateData []byte `json:"-" bson:"certificate"`
	PrivateKeyData  []byte `json:"-" bson:"private_key"`
	Issuer          string `json:"issuer" bson:"issuer"`
	ExpiresAt       int64  `json:"expires_at" bson:"expires_at"`
	CreatedAt       int64  `json:"created_at" bson:"created_at"`
}

func Load() {
	switch tools.Env.StoreType {
	case "sql":
		Active = NewSQLStore()
	case "mongo":
		Active = NewMongoStore()
	case "file":
		Active = NewFileStore()
	default:
		log.Fatal("STORE_TYPE not supported")
	}

	err := Active.Load()
	if err != nil {
		log.Fatal(err)
	}
}
