package stores

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jonasroussel/hyve/tools"
)

type FileStore struct {
	Directory string
}

func NewFileStore() FileStore {
	dir := os.Getenv("STORE_DIR")
	if dir == "" {
		dir = tools.Env.DataDir + "/certificates"
	}

	return FileStore{
		Directory: dir,
	}
}

func (store FileStore) Load() error {
	sd := store.Directory
	if sd[len(sd)-1] == '/' {
		store.Directory = sd[:len(sd)-1]
	}

	err := os.MkdirAll(store.Directory, 0700)
	if err != nil {
		return err
	}

	return nil
}

func (store FileStore) AddCertificate(domain string, cert Certificate) error {
	domainDir := fmt.Sprintf("%s/%s", store.Directory, domain)

	err := os.MkdirAll(domainDir, 0700)
	if err != nil {
		return err
	}

	// Certificate
	err = os.WriteFile(fmt.Sprintf("%s/certificate.crt", domainDir), cert.CertificateData, 0600)
	if err != nil {
		return err
	}

	// Private Key
	err = os.WriteFile(fmt.Sprintf("%s/private.key", domainDir), cert.PrivateKeyData, 0600)
	if err != nil {
		return err
	}

	// Info
	info, err := json.Marshal(map[string]interface{}{
		"issuer":     cert.Issuer,
		"expires_at": cert.ExpiresAt,
		"created_at": cert.CreatedAt,
	})
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/info.json", domainDir), info, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (store FileStore) GetCertificate(domain string) (*Certificate, error) {
	domainDir := fmt.Sprintf("%s/%s", store.Directory, domain)

	var cert Certificate

	info, err := os.ReadFile(fmt.Sprintf("%s/info.json", domainDir))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(info, &cert)
	if err != nil {
		return nil, err
	}

	cert.CertificateData, err = os.ReadFile(fmt.Sprintf("%s/certificate.crt", domainDir))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	cert.PrivateKeyData, err = os.ReadFile(fmt.Sprintf("%s/private.key", domainDir))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	cert.Domain = domain

	return &cert, nil
}

func (store FileStore) GetAllCertificates(exp int64) []Certificate {
	certs := []Certificate{}

	entries, err := os.ReadDir(store.Directory)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		domainDir := fmt.Sprintf("%s/%s", store.Directory, entry.Name())

		info, err := os.ReadFile(fmt.Sprintf("%s/info.json", domainDir))
		if err != nil {
			continue
		}

		var cert Certificate
		err = json.Unmarshal(info, &cert)
		if err != nil {
			continue
		}

		if cert.ExpiresAt > exp {
			continue
		}

		cert.CertificateData, err = os.ReadFile(fmt.Sprintf("%s/certificate.crt", domainDir))
		if err != nil {
			continue
		}

		cert.PrivateKeyData, err = os.ReadFile(fmt.Sprintf("%s/private.key", domainDir))
		if err != nil {
			continue
		}

		cert.Domain = entry.Name()

		certs = append(certs, cert)
	}

	return certs
}

func (store FileStore) UpdateCertificate(domain string, cert Certificate) error {
	domainDir := fmt.Sprintf("%s/%s", store.Directory, domain)

	// Certificate
	if cert.CertificateData != nil {
		err := os.WriteFile(fmt.Sprintf("%s/certificate.crt", domainDir), cert.CertificateData, 0600)
		if err != nil {
			return err
		}
	}

	// Private Key
	if cert.PrivateKeyData != nil {
		err := os.WriteFile(fmt.Sprintf("%s/private.key", domainDir), cert.PrivateKeyData, 0600)
		if err != nil {
			return err
		}
	}

	// Info
	var oldCert Certificate

	infoData, err := os.ReadFile(fmt.Sprintf("%s/info.json", domainDir))
	if err != nil {
		return err
	}
	err = json.Unmarshal(infoData, &oldCert)
	if err != nil {
		return err
	}

	if cert.Issuer == "" {
		cert.Issuer = oldCert.Issuer
	}
	if cert.ExpiresAt == 0 {
		cert.ExpiresAt = oldCert.ExpiresAt
	}
	if cert.CreatedAt == 0 {
		cert.CreatedAt = oldCert.CreatedAt
	}

	info, err := json.Marshal(map[string]interface{}{
		"issuer":     cert.Issuer,
		"expires_at": cert.ExpiresAt,
		"created_at": cert.CreatedAt,
	})
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/info.json", domainDir), info, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (store FileStore) RemoveCertificate(domain string) error {
	domainDir := fmt.Sprintf("%s/%s", store.Directory, domain)

	err := os.RemoveAll(domainDir)
	if err != nil {
		return err
	}

	return nil
}
