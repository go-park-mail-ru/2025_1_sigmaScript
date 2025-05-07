package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/user"
)

func main() {
	a, err := user.New(false)
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}
