package tools

import (
	"log"
	"os"
)

var Env struct {
	Target      string `env:"TARGET"`
	DataDir     string `env:"DATA_DIR"`
	UserDir     string `env:"USER_DIR"`
	StoreType   string `env:"STORE"`
	AdminDomain string `env:"ADMIN_DOMAIN"`
	AdminKey    string `env:"ADMIN_KEY"`
}

func LoadEnv() {
	Env.Target = os.Getenv("TARGET")
	if Env.Target == "" {
		log.Fatal("TARGET environment variable is not set")
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
	} else {
		Env.StoreType = "file"
	}

	Env.AdminDomain = os.Getenv("ADMIN_DOMAIN")
	Env.AdminKey = os.Getenv("ADMIN_KEY")

	if Env.AdminDomain == "" || Env.AdminKey == "" {
		log.Println("\033[33mWARNING: ADMIN(s) environment variables are not set, admin API will not be available\033[0m")
	}
}
