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

func (store FileStore) AddCertificate(sni string, cert Certificate) error {
	sniDir := fmt.Sprintf("%s/%s", store.Directory, sni)

	err := os.MkdirAll(sniDir, 0700)
	if err != nil {
		return err
	}

	// Certificate
	err = os.WriteFile(fmt.Sprintf("%s/certificate.crt", sniDir), cert.CertificateData, 0600)
	if err != nil {
		return err
	}

	// Private Key
	err = os.WriteFile(fmt.Sprintf("%s/private.key", sniDir), cert.PrivateKeyData, 0600)
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
	err = os.WriteFile(fmt.Sprintf("%s/info.json", sniDir), info, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (store FileStore) GetCertificate(sni string) (*Certificate, error) {
	sniDir := fmt.Sprintf("%s/%s", store.Directory, sni)

	var cert Certificate

	info, err := os.ReadFile(fmt.Sprintf("%s/info.json", sniDir))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(info, &cert)
	if err != nil {
		return nil, err
	}

	cert.CertificateData, err = os.ReadFile(fmt.Sprintf("%s/certificate.crt", sniDir))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	cert.PrivateKeyData, err = os.ReadFile(fmt.Sprintf("%s/private.key", sniDir))
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (store FileStore) UpdateCertificate(sni string, cert Certificate) error {
	sniDir := fmt.Sprintf("%s/%s", store.Directory, sni)

	// Certificate
	if cert.CertificateData != nil {
		err := os.WriteFile(fmt.Sprintf("%s/certificate.crt", sniDir), cert.CertificateData, 0600)
		if err != nil {
			return err
		}
	}

	// Private Key
	if cert.PrivateKeyData != nil {
		err := os.WriteFile(fmt.Sprintf("%s/private.key", sniDir), cert.PrivateKeyData, 0600)
		if err != nil {
			return err
		}
	}

	// Info
	var oldCert Certificate

	infoData, err := os.ReadFile(fmt.Sprintf("%s/info.json", sniDir))
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
	err = os.WriteFile(fmt.Sprintf("%s/info.json", sniDir), info, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (store FileStore) RemoveCertificate(sni string) error {
	sniDir := fmt.Sprintf("%s/%s", store.Directory, sni)

	err := os.RemoveAll(sniDir)
	if err != nil {
		return err
	}

	return nil
}
