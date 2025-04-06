package repository

import (
	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
)

type UserRepositoryInterface interface {
}

type UserRepository struct {
	// username --> hashedPass
	users synccredmap.SyncCredentialsMap
	cfg   *config.Cookie
}

func NewUserRepository() UserRepositoryInterface {
	res := &UserRepository{
		users: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}
