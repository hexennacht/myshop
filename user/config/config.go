package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Environment   string `envconfig:"ENVIRONMENT"`
	Port          string `envconfig:"PORT"`
	SecretJWT     string `envconfig:"SECRET_JWT"`
	TokenLifeTime int64  `envconfig:"TOKEN_LIFE_TIME"` // in hour

	DbHost     string `envconfig:"DB_HOST"`
	DbPort     string `envconfig:"DB_PORT"`
	DbUser     string `envconfig:"DB_USER"`
	DbPassword string `envconfig:"DB_PASSWORD"`
	DbName     string `envconfig:"DB_NAME"`

	RedisHost string `envconfig:"REDIS_HOST"`
	RedisPort string `envconfig:"REDIS_PORT"`
}

func Read() *Config {
	cfg := Config{}

	envconfig.MustProcess("", &cfg)

	return &cfg
}
