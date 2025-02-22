package handlers

import (
  "encoding/json"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
)

func GetGenres(w http.ResponseWriter, r *http.Request) {
  if err := json.NewEncoder(w).Encode(mocks.Genres); err != nil {
    http.Error(w, "Encode error", http.StatusInternalServerError)
    return
  }
}
