package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/config/defaults"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/errors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	MaxFindingEnvDepth = 100

	KinolkAvatarsFolder = "KINOLK_AVATARS_FOLDER"
)

type Config struct {
	Listener Listener `yaml:"listener" mapstructure:"listener"`
	Database Database `yaml:"database" mapstructure:"database"`
	Storage  Storage  `yaml:"storage" mapstructure:"storage"`
}

type Listener struct {
	Port string `yaml:"port" mapstructure:"port"`
}

type Database struct {
	Host            string        `yaml:"host" mapstructure:"host"`
	Port            int           `yaml:"port" mapstructure:"port"`
	User            string        `yaml:"user" mapstructure:"user"`
	Password        string        `yaml:"password" mapstructure:"password"`
	Name            string        `yaml:"dbname" mapstructure:"dbname"`
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
}

type Storage struct {
	UserAvatarsFullPath   string `yaml:"full_path" mapstructure:"full_path"`
	UserAvatarsStaticPath string `yaml:"relative_path" mapstructure:"relative_path"`
}

func New() (*Config, error) {
	log.Info().Msg("Initializing config")

	if err := setupViper(); err != nil {
		log.Error().Err(errors.Wrap(err, errs.ErrInitializeConfig)).Msg(errors.Wrap(err, errs.ErrInitializeConfig).Error())
		return nil, errors.Wrap(err, errs.ErrInitializeConfig)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error().Err(errors.Wrap(err, errs.ErrUnmarshalConfig)).Msg(errors.Wrap(err, errs.ErrUnmarshalConfig).Error())
		return nil, errors.Wrap(err, errs.ErrUnmarshalConfig)
	}

	log.Info().Msg("Config initialized")
	return &cfg, nil
}

func setupListener() {
	viper.SetDefault("listener.port", defaults.ListenerPort)
}

func setupDatabase() {
	viper.SetDefault("database.host", defaults.DBHost)
	viper.SetDefault("database.port", defaults.DBPort)
	viper.SetDefault("database.user", defaults.DBUser)
	viper.SetDefault("database.password", defaults.DBPassword)
	viper.SetDefault("database.dbname", defaults.DBName)
	viper.SetDefault("database.max_open_conns", defaults.DBMaxOpenConns)
	viper.SetDefault("database.max_idle_conns", defaults.DBMaxIdleConns)
	viper.SetDefault("database.conn_max_lifetime", defaults.DBConnMaxLifetime)
	viper.SetDefault("database.conn_max_idle_time", defaults.DBConnMaxIdleTime)
}

func setupStorage() {
	viper.SetDefault("storage.full_path", viper.GetString(KinolkAvatarsFolder))
	viper.SetDefault("storage.relative_path", defaults.StorageRelativePath)
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
	viper.AddConfigPath(viper.GetString("VIPER_MOVIE_CONFIG_PATH"))

	setupListener()
	setupDatabase()
	setupStorage()

	if err := viper.MergeInConfig(); err != nil {
		wrapped := errors.Wrap(err, errs.ErrReadConfig)
		log.Error().Err(wrapped).Msg(wrapped.Error())
		return wrapped
	}

	log.Info().Msg("Viper initialized")
	return nil
}
