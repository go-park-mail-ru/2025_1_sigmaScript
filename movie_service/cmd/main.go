package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/internal/movie"
)

func main() {
	a, err := movie.New(false)
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}
