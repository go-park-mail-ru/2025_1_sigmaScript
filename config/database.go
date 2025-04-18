package config

import (
	"context"

	"github.com/rs/zerolog/log"
)

type ContextPgDBKey struct{}

type ConfigPgDB struct {
	Listener  Listener  `yaml:"listener"`
	Databases Databases `yaml:"databases"`
}

// Listener contains listener port
type Listener struct {
	Port string `yaml:"port"`
}

// Databases contains database configuration
type Databases struct {
	Postgres     Postgres            `yaml:"postgres"`
	LocalStorage LocalAvatarsStorage `yaml:"localStorage"`
}

// Postgres contains postgres configuration
type Postgres struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Name            string `yaml:"name"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
}

// LocalAvatarsStorage contains local avatar storage paths
type LocalAvatarsStorage struct {
	UserAvatarsFullPath     string `yaml:"userAvatarsFullPath"`
	UserAvatarsRelativePath string `yaml:"userAvatarsRelativePath"`
}

func FromPgDatabaseContext(ctx context.Context) *ConfigPgDB {
	pgDatabaseConfig, ok := ctx.Value(ContextPgDBKey{}).(ConfigPgDB)
	if !ok {
		log.Error().Msg("cant convert from context")
		return nil
	}
	return &pgDatabaseConfig
}

func WrapPgDatabaseContext(ctx context.Context, data interface{}) context.Context {
	return context.WithValue(ctx, ContextPgDBKey{}, data)
}
