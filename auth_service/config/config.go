package config

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	MaxFindingEnvDepth = 100
	SessionLength      = 32

	POSTGRES_USER = "POSTGRES_USER"
	DB_PASSWORD   = "DB_PASSWORD"
	POSTGRES_DB   = "POSTGRES_DB"
	POSTGRES_HOST = "POSTGRES_HOST"
	POSTGRES_PORT = "POSTGRES_PORT"
	REDIS_HOST    = "REDIS_HOST"
	REDIS_PORT    = "REDIS_PORT"

	KinolkAvatarsFolder     = "KINOLK_AVATARS_FOLDER"
	KinolkAvatarsStaticPath = "KINOLK_AVATARS_STATIC_PATH"
)

type Config struct {
	Listener Listener `yaml:"cookie" mapstructure:"listener"`
}

type Cookie struct {
	SessionName   string        `yaml:"session_name" mapstructure:"session_name"`
	SessionLength int           `yaml:"session_length" mapstructure:"session_length"`
	HTTPOnly      bool          `yaml:"http_only" mapstructure:"http_only"`
	Secure        bool          `yaml:"secure" mapstructure:"secure"`
	SameSite      http.SameSite `yaml:"same_site" mapstructure:"same_site"`
	Path          string        `yaml:"path" mapstructure:"path"`
	ExpirationAge int           `yaml:"expiration_age" mapstructure:"expiration_age"`
}

func New() (*Config, error) {
	log.Info().Msg("Initializing config")

	var config Config

	config.Listener = Listener{
		Port: ":8081",
	}

	log.Info().Msg("Config initialized")
	return &config, nil
}
