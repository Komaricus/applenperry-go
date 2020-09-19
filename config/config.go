package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
	PRODUCTION  bool
}

func GetConfig() Configuration {
	configuration := Configuration{}
	if err := gonfig.GetConf("config/config.json", &configuration); err != nil {
		panic("Can't read config")
	}
	return configuration
}
