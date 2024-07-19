package acme

import (
	"log"
	"time"

	"github.com/jonasroussel/hyve/stores"
	"github.com/robfig/cron/v3"
)

func ActivateAutoRenew() {
	jobs := cron.New()

	jobs.AddFunc("@daily", renewAllNearlyExpired)

	jobs.Start()
}

func renewAllNearlyExpired() {
	now1week := time.Now().Add(7 * (24 * time.Hour)).Unix()

	for _, cert := range stores.Active.GetAllCertificates(now1week) {
		// A check just in case the `GetAllCertificates` function fails to compare the expiry date.
		if cert.ExpiresAt > now1week {
			continue
		}

		err := RenewDomain(cert.Domain)
		if err != nil {
			log.Printf("Failed to renew certificate (%s): %s", cert.Domain, err)
			continue
		}

		log.Printf("Certificate renewed: %s", cert.Domain)
	}
}
