package config

import (
	"log"
	"os"

	"github.com/jonasroussel/proxbee/stores"
)

var (
	TARGET       string       // Mandatory
	DATA_DIR     string       // Optional
	STORE        stores.Store // Optional
	ADMIN_DOMAIN string       // Optional but mandatory to enable admin api
	ADMIN_KEY    string       // Optional but mandatory to enable admin api
)

func Load() {
	TARGET = os.Getenv("TARGET")
	if TARGET == "" {
		panic("TARGET environment variable is not set")
	}

	DATA_DIR = os.Getenv("DATA_DIR")
	if DATA_DIR == "" {
		DATA_DIR = "/var/lib/proxbee"
	}

	if os.Getenv("STORE") == "sql" {
		STORE = stores.NewSQLStore()
	} else {
		STORE = stores.NewFileStore(DATA_DIR)
	}

	ADMIN_DOMAIN = os.Getenv("ADMIN_DOMAIN")
	ADMIN_KEY = os.Getenv("ADMIN_KEY")

	if ADMIN_DOMAIN == "" || ADMIN_KEY == "" {
		log.Println("\033[33mWARNING: ADMIN(s) environment variables are not set, admin API will not be available\033[0m")
	}
}
