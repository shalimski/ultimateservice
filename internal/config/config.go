package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host            string        `env:"HOST" env-default:"0.0.0.0:9020"`
	DebugURI        string        `env:"DEBUGURI" env-default:"0.0.0.0:9021"`
	ReadTimeout     time.Duration `env:"READTIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `env:"WRITETIMEOUT" env-default:"20s"`
	IdleTimeout     time.Duration `env:"IDLETIMEOUT" env-default:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWNTIMEOUT" env-default:"20s"`
}

func New() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
