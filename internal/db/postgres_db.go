package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"

	"github.com/rs/zerolog/log"
)

const (
	DB_MAX_PINGS = 3
)

// SetupDatabase connects to Postgres and returns instance of sql.DB
func SetupDatabase(ctx context.Context, cancel context.CancelFunc) (*sql.DB, error) {
	ctxVals := config.FromPgDatabaseContext(ctx)
	defer cancel()
	log.Info().Msg("Trying to connect to Postgres database")

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			DB, err := connectToPostgresDB(ctxVals)
			if err == nil {
				log.Info().Msg("Postgres database connection opened successfully")
				return DB, nil
			}

			log.Error().Err(fmt.Errorf("failed to connect to database. Error: %v Retrying", err)).Msg("setup_db_error")
			time.Sleep(3 * time.Second)
		}
	}
}

func connectToPostgresDB(cfg *config.ConfigPgDB) (*sql.DB, error) {
	if cfg == nil {
		errConf := fmt.Errorf("error reading DB config: %v", cfg)
		log.Error().Err(errConf).Msg(errConf.Error())
		return nil, errConf
	}
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Databases.Postgres.Host,
		cfg.Databases.Postgres.Port,
		cfg.Databases.Postgres.User,
		cfg.Databases.Postgres.Password,
		cfg.Databases.Postgres.Name,
	)

	DB, err := sql.Open("postgres", connectionString)
	if err != nil {
		errMsg := fmt.Errorf("error while opening DB: %w", err)
		log.Error().Err(errMsg).Msg("connect_db_error")

		return nil, errMsg
	}

	DB.SetMaxOpenConns(cfg.Databases.Postgres.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.Databases.Postgres.MaxIdleConns)
	DB.SetConnMaxLifetime(time.Duration(cfg.Databases.Postgres.ConnMaxLifetime) * time.Minute)
	DB.SetConnMaxIdleTime(time.Duration(cfg.Databases.Postgres.ConnMaxIdleTime) * time.Minute)

	time.Sleep(1 * time.Second)

	for i := range DB_MAX_PINGS {
		err = DB.Ping()
		if err == nil {
			break
		}

		errMsg := fmt.Errorf("ping №%d:error while pinging DB: %w", i+1, err)
		log.Error().Err(errMsg).Msg("ping_db_error")
		if i == DB_MAX_PINGS-1 {
			return nil, errMsg
		}

		time.Sleep(2 * time.Second)
	}

	log.Info().Msg("Database pinged successfully")

	return DB, nil
}
