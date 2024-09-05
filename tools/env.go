package tools

import (
	"log"
	"os"
	"regexp"
)

var Env struct {
	Target         string         `env:"TARGET"`
	DYNAMIC_TARGET string         `env:"DYNAMIC_TARGET"`
	DataDir        string         `env:"DATA_DIR"`
	UserDir        string         `env:"USER_DIR"`
	StoreType      string         `env:"STORE"`
	DNSProvider    string         `env:"DNS_PROVIDER"`
	AdminDomain    string         `env:"ADMIN_DOMAIN"`
	AdminKey       string         `env:"ADMIN_KEY"`
	Blacklist      *regexp.Regexp `env:"BLACKLIST"`
}

func LoadEnv() {
	Env.Target = os.Getenv("TARGET")
	Env.DYNAMIC_TARGET = os.Getenv("DYNAMIC_TARGET")
	if Env.Target == "" && Env.DYNAMIC_TARGET == "" {
		log.Fatal("TARGET or DYNAMIC_TARGET environment variable must be set to start Hyve")
	}

	Env.DataDir = os.Getenv("DATA_DIR")
	if Env.DataDir == "" {
		Env.DataDir = "/var/lib/hyve"
	}
	err := os.MkdirAll(Env.DataDir, 0700)
	if err != nil {
		log.Fatal(err)
	}

	Env.UserDir = os.Getenv("USER_DIR")
	if Env.UserDir == "" {
		Env.UserDir = Env.DataDir + "/user"
	}
	err = os.MkdirAll(Env.UserDir, 0700)
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("STORE") == "sql" {
		Env.StoreType = "sql"
	} else if os.Getenv("STORE") == "mongo" {
		Env.StoreType = "mongo"
	} else {
		Env.StoreType = "file"
	}

	Env.DNSProvider = os.Getenv("DNS_PROVIDER")

	Env.AdminDomain = os.Getenv("ADMIN_DOMAIN")
	Env.AdminKey = os.Getenv("ADMIN_KEY")

	if Env.AdminDomain == "" || Env.AdminKey == "" {
		log.Println("\033[33mWARNING: ADMIN(s) environment variables are not set, admin API will not be available\033[0m")
	}

	if os.Getenv("BLACKLIST") != "" {
		Env.Blacklist = regexp.MustCompile(os.Getenv("BLACKLIST"))
	}
}
