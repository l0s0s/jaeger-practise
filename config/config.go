package config

import (
	"log"

	"github.com/cristalhq/aconfig"
)

// Config describes application config.
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

// Parse configs from environment. Shuts down the application if something goes wrong.
func Parse() Config {
	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{})
	if err := loader.Load(); err != nil {
		log.Fatal(err)
	}

	return cfg
}
