package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port     string `env:"PORT" env-default:"9020"`
	Host     string `env:"HOST" env-default:"0.0.0.0"`
	DebugURI string `env:"DEBUGURI" env-default:"0.0.0.0:9021"`
}

func New() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
