package models

import "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"

type SearchResponseJSON struct {
	MovieCollection []mocks.Movie      `json:"movie_collection,omitempty"`
	Actors          []mocks.PersonJSON `json:"actors,omitempty"`
}

type SearchRequestJSON struct {
	SearchString string `json:"search"`
}
