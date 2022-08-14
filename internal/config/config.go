package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port string `env:"PORT" env-default:"9020"`
	Host string `env:"HOST" env-default:"0.0.0.0"`
}

func New() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
