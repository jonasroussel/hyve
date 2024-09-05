package caching

import (
	"container/list"
	"crypto/tls"
	"errors"

	"github.com/jonasroussel/hyve/stores"
	"github.com/jonasroussel/hyve/tools"
)

var cachingQueue = list.New()

type CahcedData struct {
	sni         string
	certificate *tls.Certificate
}

func CertificateRetriever(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if chi.ServerName == "" {
		return nil, errors.New("server name (sni) is empty")
	}

	if tools.Env.Blacklist != nil && tools.Env.Blacklist.MatchString(chi.ServerName) {
		return nil, errors.New("server name (" + chi.ServerName + ") is blacklisted")
	}

	cert := loadFromCache(chi.ServerName)
	if cert != nil {
		return cert, nil
	}

	cert, err := loadFromStore(chi.ServerName)
	if err != nil {
		return nil, err
	}

	addToCache(chi.ServerName, cert)

	return cert, nil
}

func loadFromStore(sni string) (*tls.Certificate, error) {
	cert, err := stores.Active.GetCertificate(sni)
	if err != nil {
		return nil, err
	}

	x509Cert, err := tls.X509KeyPair(cert.CertificateData, cert.PrivateKeyData)
	if err != nil {
		return nil, err
	}

	return &x509Cert, nil
}

func loadFromCache(sni string) *tls.Certificate {
	for e := cachingQueue.Front(); e != nil; e = e.Next() {
		cert := e.Value.(CahcedData)
		if cert.sni == sni {
			return cert.certificate
		}
	}

	return nil
}

func addToCache(sni string, cert *tls.Certificate) {
	cachingQueue.PushFront(CahcedData{
		sni:         sni,
		certificate: cert,
	})

	// Remove the oldest certificate if the cache is too big
	if cachingQueue.Len() > 1000 {
		cachingQueue.Remove(cachingQueue.Back())
	}
}
