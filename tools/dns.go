package tools

import (
	"errors"
	"net"
)

var adminIPS []net.IP

func IsDNSValid(domain string) bool {
	err := lazyLoadAdminIPS()
	if err != nil {
		return false
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return false
	}

	for _, ip := range ips {
		for _, adminIP := range adminIPS {
			if ip.Equal(adminIP) {
				return true
			}
		}
	}

	return true
}

func lazyLoadAdminIPS() error {
	if adminIPS != nil {
		return nil
	}

	if Env.AdminDomain == "" {
		return errors.New("ADMIN_DOMAIN environment variable is not set")
	}

	ips, err := net.LookupIP(Env.AdminDomain)
	if err != nil {
		return err
	}

	adminIPS = ips

	return nil
}
