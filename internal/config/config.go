// Package config builds app config
package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config parses env vars.
type Config struct {
	Host     string `required:"true"`
	GRPCPort string `required:"true"`
}

// NewConfig constructs app config and returns config.
func NewConfig() *Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	return &cfg
}
