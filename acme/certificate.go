package acme

import (
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"

	"github.com/jonasroussel/proxbee/stores"
)

func RegisterDomain(domain string) error {
	config := lego.NewConfig(ActiveUser)

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	err = client.Challenge.SetHTTP01Provider(HTTP01Provider)
	if err != nil {
		return err
	}

	req := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	rawCert, err := client.Certificate.Obtain(req)
	if err != nil {
		return err
	}

	cert := stores.Certificate{
		CertificateData: rawCert.Certificate,
		PrivateKeyData:  rawCert.PrivateKey,
		Issuer:          rawCert.CertStableURL,
		CreatedAt:       time.Now().Unix(),
		ExpiresAt:       time.Now().Add((90 - 1) * (24 * time.Hour)).Unix(), // -1 is just for safety
	}

	err = stores.Active.AddCertificate(rawCert.Domain, cert)
	if err != nil {
		return err
	}

	return nil
}

func RenewDomain(domain string) error {
	// TODO

	return nil
}
