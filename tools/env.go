package tools

import "os"

var Env struct {
	Target      string `env:"TARGET"`
	DataDir     string `env:"DATA_DIR"`
	StoreType   string `env:"STORE"`
	AdminDomain string `env:"ADMIN_DOMAIN"`
	AdminKey    string `env:"ADMIN_KEY"`
}

func LoadEnv() {
	Env.Target = os.Getenv("TARGET")
	if Env.Target == "" {
		panic("TARGET environment variable is not set")
	}

	Env.DataDir = os.Getenv("DATA_DIR")
	if Env.DataDir == "" {
		Env.DataDir = "/var/lib/proxbee"
	}

	if os.Getenv("STORE") == "sql" {
		Env.StoreType = "sql"
	} else {
		Env.StoreType = "file"
	}

	Env.AdminDomain = os.Getenv("ADMIN_DOMAIN")
	Env.AdminKey = os.Getenv("ADMIN_KEY")

	if Env.AdminDomain == "" || Env.AdminKey == "" {
		panic("\033[33mWARNING: ADMIN(s) environment variables are not set, admin API will not be available\033[0m")
	}
}
