package config

import "os"

type Config struct {
	KrakenHost string
	HttpAddr   string
}

func Read() Config {
	config := Config{
		KrakenHost: "https://api.kraken.com",
		HttpAddr:   ":8090",
	}

	krakenHost, exists := os.LookupEnv("KRAKEN_HOST")

	if exists {
		config.KrakenHost = krakenHost
	}

	httpAddr, exists := os.LookupEnv("HTTP_ADDR")

	if exists {
		config.HttpAddr = httpAddr
	}

	return config
}
