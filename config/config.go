package config

import (
  "time"

  "github.com/spf13/viper"
)

type Config struct {
  Server Server `yaml:"server" mapstructure:"server"`
}

type Server struct {
  Address         string        `yaml:"address" mapstructure:"address"`
  Port            int           `yaml:"port" mapstructure:"port"`
  ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
  WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
  ShutdownTimeout time.Duration `yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
  IdleTimeout     time.Duration `yaml:"idle_timeout" mapstructure:"idle_timeout"`
}

func New() (*Config, error) {
  if err := setupViper(); err != nil {
    return nil, err
  }
  var config Config
  if err := viper.Unmarshal(&config); err != nil {
    return nil, err
  }
  return &config, nil
}

func setupViper() error {
  viper.SetConfigName("config")
  viper.SetConfigType("yml")
  viper.AddConfigPath("./internal/config")

  viper.SetDefault("server.address", "localhost")
  viper.SetDefault("server.port", 8080)
  viper.SetDefault("server.read_timeout", time.Second*5)
  viper.SetDefault("server.write_timeout", time.Second*5)
  viper.SetDefault("server.shutdown_timeout", time.Second*30)
  viper.SetDefault("server.idle_timeout", time.Second*60)

  if err := viper.ReadInConfig(); err != nil {
    return err
  }
  return nil
}
