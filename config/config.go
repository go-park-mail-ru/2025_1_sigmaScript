package config

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config/defaults"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	MaxFindingEnvDepth = 100

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
	Server         Server     `yaml:"server" mapstructure:"server"`
	Cookie         Cookie     `yaml:"cookie" mapstructure:"cookie"`
	PostgresConfig ConfigPgDB `yaml:"cookie" mapstructure:"configpgdb"`
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
	ExpirationAge int           `yaml:"expiration_age" mapstructure:"expiration_age"`
}

func New() (*Config, error) {
	log.Info().Msg("Initializing config")

	if err := setupViper(); err != nil {
		log.Error().Err(errors.Wrap(err, errs.ErrInitializeConfig)).Msg(errors.Wrap(err, errs.ErrInitializeConfig).Error())
		return nil, errors.Wrap(err, errs.ErrInitializeConfig)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Error().Err(errors.Wrap(err, errs.ErrUnmarshalConfig)).Msg(errors.Wrap(err, errs.ErrUnmarshalConfig).Error())
		return nil, errors.Wrap(err, errs.ErrUnmarshalConfig)
	}

	config.PostgresConfig = *setupPgConfig()

	log.Info().Msg("Config initialized")
	return &config, nil
}

func setupPgConfig() *ConfigPgDB {
	// TODO fix config: it`s test database test password
	postgres := Postgres{
		Host:            viper.GetString(POSTGRES_HOST),
		Port:            viper.GetInt(POSTGRES_PORT),
		User:            viper.GetString(POSTGRES_USER),
		Password:        viper.GetString(DB_PASSWORD),
		Name:            viper.GetString(POSTGRES_DB),
		MaxOpenConns:    100,
		MaxIdleConns:    30,
		ConnMaxLifetime: 30,
		ConnMaxIdleTime: 5,
	}

	avatarLocalStorage := LocalAvatarsStorage{
		UserAvatarsFullPath:   viper.GetString(KinolkAvatarsFolder),
		UserAvatarsStaticPath: viper.GetString(KinolkAvatarsStaticPath),
	}

	pgDatabase := Databases{
		Postgres:     postgres,
		LocalStorage: avatarLocalStorage,
	}

	pgListener := Listener{
		Port: viper.GetString(POSTGRES_PORT),
	}

	cfgDB := ConfigPgDB{
		Listener:  pgListener,
		Databases: pgDatabase,
	}

	return &cfgDB
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
	viper.SetDefault("cookie.expiration_age", defaults.ExpirationAge)
}

func findEnvDir() (string, error) {
	log.Info().Msg("Finding environment dir")
	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, errs.ErrGetDirectory)
	}

	for i := 0; i < MaxFindingEnvDepth; i++ {
		path := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(path); err == nil {
			log.Info().Msg("Found .env file")
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", errors.Wrap(err, errs.ErrDirectoryNotFound)
		}
		currentDir = parentDir
	}

	return "", errors.Wrap(err, errs.ErrDirectoryNotFound)
}

func setupViper() error {
	log.Info().Msg("Initializing viper")

	envDir, err := findEnvDir()
	if err != nil {
		wrapped := errors.Wrap(err, errs.ErrDirectoryNotFound)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(envDir)

	if err := viper.ReadInConfig(); err != nil {
		wrapped := errors.Wrap(err, errs.ErrReadEnvironment)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(viper.GetString("VIPER_CONFIG_PATH"))

	setupServer()
	setupCookie()

	if err := viper.MergeInConfig(); err != nil {
		wrapped := errors.Wrap(err, errs.ErrReadConfig)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	log.Info().Msg("Viper initialized")
	return nil
}
