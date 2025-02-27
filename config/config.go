package config

import (
  "net/http"
  "time"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config/defaults"
  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"
  "github.com/spf13/viper"
)

type Config struct {
  Server Server `yaml:"server" mapstructure:"server"`
  Cookie Cookie `yaml:"cookie" mapstructure:"cookie"`
}

type Server struct {
  Address         string        `yaml:"address" mapstructure:"address"`
  Port            int           `yaml:"port" mapstructure:"port"`
  ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
  WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
  ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
  IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

type Cookie struct {
  SessionName   string        `yaml:"session_name" mapstructure:"session_name"`
  SessionLength int           `yaml:"session_length" mapstructure:"session_length"`
  HTTPOnly      bool          `yaml:"http_only" mapstructure:"http_only"`
  Secure        bool          `yaml:"secure" mapstructure:"secure"`
  SameSite      http.SameSite `yaml:"same_site" mapstructure:"same_site"`
  Path          string        `yaml:"path" mapstructure:"path"`
  ResetMaxAge   int           `yaml:"reset_max_age" mapstructure:"reset_max_age"`
}

func New() (*Config, error) {
  log.Info().Msg("Initializing config")

  if err := setupViper(); err != nil {
    log.Error().Err(err).Msg("Error initializing config")
    return nil, errors.Wrap(err, "failed to initialize config")
  }

  var config Config
  if err := viper.Unmarshal(&config); err != nil {
    log.Error().Err(err).Msg("Error unmarshalling config")
    return nil, errors.Wrap(err, "failed to unmarshal config")
  }

  log.Info().Msg("Config initialized")
  return &config, nil
}

func setupServer() {
  viper.SetDefault("server.address", defaults.Address)
  viper.SetDefault("server.port", defaults.Port)
  viper.SetDefault("server.read_timeout", defaults.ReadTimeout)
  viper.SetDefault("server.write_timeout", defaults.WriteTimeout)
  viper.SetDefault("server.shutdown_timeout", defaults.ShutdownTimeout)
  viper.SetDefault("server.idle_timeout", defaults.IdleTimeout)
}

func setupCookie() {
  viper.SetDefault("cookie.session_name", defaults.SessionName)
  viper.SetDefault("cookie.session_length", defaults.SessionLength)
  viper.SetDefault("cookie.http_only", defaults.HTTPOnly)
  viper.SetDefault("cookie.secure", defaults.Secure)
  viper.SetDefault("cookie.same_site", defaults.SameSite)
  viper.SetDefault("cookie.path", defaults.Path)
  viper.SetDefault("cookie.reset_max_age", defaults.ResetMaxAge)
}

func setupViper() error {
  log.Info().Msg("Initializing viper")

  viper.SetConfigName("config")
  viper.SetConfigType("yml")
  viper.AddConfigPath("./internal/config")

  setupServer()
  setupCookie()

  if err := viper.ReadInConfig(); err != nil {
    log.Error().Err(err).Msg("Error reading config")
    return errors.Wrap(err, "failed to read config")
  }

  log.Info().Msg("Viper initialized")
  return nil
}
