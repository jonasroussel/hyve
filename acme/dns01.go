package acme

import (
	"errors"

	"github.com/go-acme/lego/providers/dns/easydns"
	"github.com/go-acme/lego/providers/dns/gandi"
	"github.com/go-acme/lego/providers/dns/godaddy"
	"github.com/go-acme/lego/providers/dns/linode"
	"github.com/go-acme/lego/providers/dns/namecheap"
	"github.com/go-acme/lego/providers/dns/namedotcom"
	"github.com/go-acme/lego/providers/dns/oraclecloud"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/arvancloud"
	"github.com/go-acme/lego/v4/providers/dns/autodns"
	"github.com/go-acme/lego/v4/providers/dns/bunny"
	"github.com/go-acme/lego/v4/providers/dns/clouddns"
	"github.com/go-acme/lego/v4/providers/dns/digitalocean"
	"github.com/go-acme/lego/v4/providers/dns/ionos"
	"github.com/go-acme/lego/v4/providers/dns/ovh"
	"github.com/go-acme/lego/v4/providers/dns/scaleway"
	"github.com/go-acme/lego/v4/providers/dns/vercel"
	"github.com/jonasroussel/hyve/tools"
)

var DNS01Provider challenge.Provider

func LoadDNS01Provider() error {
	switch tools.Env.DNSProvider {
	case "arvancloud":
		provider, err := arvancloud.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "autodns":
		provider, err := autodns.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "bunny":
		provider, err := bunny.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "clouddns":
		provider, err := clouddns.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "digitalocean":
		provider, err := digitalocean.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "easydns":
		provider, err := easydns.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "gandi":
		provider, err := gandi.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "godaddy":
		provider, err := godaddy.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "ionos":
		provider, err := ionos.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "linode":
		provider, err := linode.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "namedotcom":
		provider, err := namedotcom.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "namecheap":
		provider, err := namecheap.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "oraclecloud":
		provider, err := oraclecloud.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "ovh":
		provider, err := ovh.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "scaleway":
		provider, err := scaleway.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	case "vercel":
		provider, err := vercel.NewDNSProvider()
		if err != nil {
			return err
		}

		DNS01Provider = provider
	default:
		return errors.New("DNS_PROVIDER not supported")
	}

	return nil
}
