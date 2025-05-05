package models

import (
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
)

type User struct {
	ID             string `json:"id,omitempty"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"`
	Avatar         string `json:"avatar"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type Profile struct {
	Username        string             `json:"username"`
	Avatar          string             `json:"avatar"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at,omitempty"`
	MovieCollection []mocks.Movie      `json:"movie_collection,omitempty"`
	Actors          []mocks.PersonJSON `json:"actors,omitempty"`
	Reviews         []mocks.ReviewJSON `json:"reviews,omitempty"`
}
