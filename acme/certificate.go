package acme

import (
	"log"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"

	"github.com/jonasroussel/hyve/stores"
	"github.com/jonasroussel/hyve/tools"
)

func RegisterDomain(domain string) error {
	config := lego.NewConfig(ActiveUser)

	config.Certificate.KeyType = certcrypto.EC256

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

func RegisterAdminDomain() {
	if tools.Env.AdminDomain == "" {
		return
	}

	cert, err := stores.Active.GetCertificate(tools.Env.AdminDomain)
	if cert != nil {
		return
	} else if err != nil && err != stores.ErrNotFound {
		log.Fatal(err)
	}

	err = RegisterDomain(tools.Env.AdminDomain)
	if err != nil {
		log.Fatal(err)
	}
}
