package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Env Environment

func LoadEnvs() {
	var defaultEnv Environment = Environment{
		Name: "Delphos Server API",
		Port: ":8080",
	}
	err := godotenv.Load()
	if err != nil {
		Env = defaultEnv
	}

	name, NameOk := os.LookupEnv("NAME")
	port, PortOk := os.LookupEnv("PORT")
	if NameOk && PortOk {
		Env.Name = name
		Env.Port = port
	} else {
		Env = defaultEnv
	}

}

func init() {
	LoadEnvs()
}
