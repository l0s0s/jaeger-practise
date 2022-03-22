package config

import (
	"log"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	Jaeger struct {
		URL string `env:"URL"`
	} `env:"JAEGER"`
	Service struct {
		Name string `env:"NAME"`
		ID   string `env:"ID"`
		Port string `env:"PORT" default:":8080"`
	} `env:"SERVICE"`
	NextURL string `env:"NEXT_URL"`
}

func Parse() Config {
	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{})
	if err := loader.Load(); err != nil {
		log.Fatal(err)
	}

	return cfg
}
