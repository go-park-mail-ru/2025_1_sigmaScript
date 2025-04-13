package models

import "time"

type User struct {
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
