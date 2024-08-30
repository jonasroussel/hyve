package tools

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"
)

var adminIPS []net.IPAddr
var resolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{Timeout: time.Millisecond * time.Duration(10000)}
		return d.DialContext(ctx, network, "1.1.1.1:53")
	},
}

func IsDNSValid(domain string) bool {
	err := lazyLoadAdminIPS()
	if err != nil {
		return false
	}

	var ips []net.IPAddr

	if strings.HasPrefix(domain, "*.") {
		ips, err = resolver.LookupIPAddr(context.Background(), strings.Replace(domain, "*", "_hyve", 1))
		if err != nil {
			return false
		}
	} else {
		ips, err = resolver.LookupIPAddr(context.Background(), domain)
		if err != nil {
			return false
		}
	}

	for _, ip := range ips {
		for _, adminIP := range adminIPS {
			if ip.IP.Equal(adminIP.IP) {
				return true
			}
		}
	}

	return false
}

func lazyLoadAdminIPS() error {
	if adminIPS != nil {
		return nil
	}

	if Env.AdminDomain == "" {
		return errors.New("ADMIN_DOMAIN environment variable is not set")
	}

	ips, err := resolver.LookupIPAddr(context.Background(), Env.AdminDomain)
	if err != nil {
		return err
	}

	adminIPS = ips

	return nil
}
