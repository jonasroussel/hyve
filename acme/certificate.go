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

var legoClient *lego.Client

func InitLego() {
	config := lego.NewConfig(ActiveUser)

	config.Certificate.KeyType = certcrypto.EC256

	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Challenge.SetHTTP01Provider(HTTP01Provider)
	if err != nil {
		log.Fatal(err)
	}

	legoClient = client
}

func RegisterDomain(domain string) error {
	req := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	rawCert, err := legoClient.Certificate.Obtain(req)
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
	cert, err := stores.Active.GetCertificate(domain)
	if err != nil {
		return err
	}

	res := certificate.Resource{
		Domain:        domain,
		CertURL:       cert.Issuer,
		CertStableURL: cert.Issuer,
		Certificate:   cert.CertificateData,
	}

	rawCert, err := legoClient.Certificate.RenewWithOptions(res, &certificate.RenewOptions{
		Bundle: true,
	})
	if err != nil {
		return err
	}

	cert = &stores.Certificate{
		CertificateData: rawCert.Certificate,
		PrivateKeyData:  rawCert.PrivateKey,
		Issuer:          rawCert.CertStableURL,
		CreatedAt:       time.Now().Unix(),
		ExpiresAt:       time.Now().Add((90 - 1) * (24 * time.Hour)).Unix(), // -1 is just for safety
	}

	err = stores.Active.UpdateCertificate(domain, *cert)
	if err != nil {
		return err
	}

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
