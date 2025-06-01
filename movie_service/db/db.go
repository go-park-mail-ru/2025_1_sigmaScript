package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/config"
	"github.com/rs/zerolog/log"
)

const (
	DbMaxPings = 3
)

// SetupDatabase connects to Postgres and returns instance of sql.DB
func SetupDatabase(ctx context.Context, cancel context.CancelFunc) (*sql.DB, error) {
	cfg := config.FromDatabaseContext(ctx)
	fmt.Println("ctxVals: ", cfg)
	defer cancel()
	log.Info().Msg("Trying to connect to Postgres database")

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			DB, err := connectToPostgresDB(cfg)
			if err == nil {
				log.Info().Msg("Postgres database connection opened successfully")
				return DB, nil
			}

			log.Error().Err(fmt.Errorf("failed to connect to database. Error: %v Retrying", err)).Msg("setup_db_error")
			time.Sleep(3 * time.Second)
		}
	}
}

func connectToPostgresDB(cfg *config.Database) (*sql.DB, error) {
	if cfg == nil {
		errConf := fmt.Errorf("error reading DB config: %v", cfg)
		log.Error().Err(errConf).Msg(errConf.Error())
		return nil, errConf
	}
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	DB, err := sql.Open("postgres", connectionString)
	if err != nil {
		errMsg := fmt.Errorf("error while opening DB: %w", err)
		log.Error().Err(errMsg).Msg("connect_db_error")

		return nil, errMsg
	}

	DB.SetMaxOpenConns(cfg.MaxOpenConns)
	DB.SetMaxIdleConns(cfg.MaxIdleConns)
	DB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	DB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	time.Sleep(1 * time.Second)

	for i := range DbMaxPings {
		err = DB.Ping()
		if err == nil {
			break
		}

		errMsg := fmt.Errorf("ping №%d:error while pinging DB: %w", i+1, err)
		log.Error().Err(errMsg).Msg("ping_db_error")
		if i == DbMaxPings-1 {
			return nil, errMsg
		}

		time.Sleep(2 * time.Second)
	}

	log.Info().Msg("Database pinged successfully")

	return DB, nil
}
