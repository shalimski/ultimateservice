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
	AuthKeysFolder  string        `env:"AUTH_KEYFOLDER" env-default:"infra/keys/"`
	AuthActiveKID   string        `env:"AUTH_ACTIVEKID" env-default:"c2e055bb-f637-4cc3-9b4a-916a8b31304a"`

	// Database PostgreSQL
	DB DB
}

type DB struct {
	User         string `env:"DB_USER" env-default:"postgres"`
	Password     string `env:"DB_PASSWORD" env-default:"postgres"`
	Host         string `env:"DB_HOST" env-default:"localhost"`
	Name         string `env:"DB_NAME" env-default:"postgres"`
	MaxIdleConns int    `env:"DB_MAXIDLE" env-default:"0"`
	MaxOpenConns int    `env:"DB_MAXOPEN" env-default:"0"`
	DisableTLS   bool   `env:"DB_DISABLETLS" env-default:"true"`
}

type auth struct{}

func New() Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
