package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	SecretKey   string        `env:"JWT_SECRET_KEY"`
	JwtTTL      time.Duration `env:"JWT_TTL"`
	DatabaseCfg database
	ServerCfg   server
}

type database struct {
	DatabaseName     string `env:"PG_DATABASE_NAME"`
	DatabaseUser     string `env:"PG_USER"`
	DatabasePassword string `env:"PG_PASSWORD"`
	DatabasePort     string `env:"PG_PORT"`
	DatabaseURL      string `env:"DATABASE_URL"`
}

type server struct {
	HttpPort         string        `env:"HTTP_PORT"`
	HttpIdleTimeOut  time.Duration `env:"HTTP_IDLE_TIMEOUT"`
	HttpWriteTimeOut time.Duration `env:"HTTP_WRITE_TIMEOUT"`
	HttpReadTimeOut  time.Duration `env:"HTTP_READ_TIMEOUT"`
}

func LoadConfig() *Config {
	path := os.Getenv("ENV_PATH")
	if path == "" {
		path = ".env"
	}
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		log.Fatalf("Error read config")
	}
	return &cfg
}
