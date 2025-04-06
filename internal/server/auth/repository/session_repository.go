package repository

import (
	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
)

type SessionRepositoryInterface interface {
}

type SessionRepository struct {
	// sessionID --> username
	sessions synccredmap.SyncCredentialsMap
	cfg      *config.Cookie
}

func NewSessionRepository() SessionRepositoryInterface {
	res := &SessionRepository{
		sessions: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}
