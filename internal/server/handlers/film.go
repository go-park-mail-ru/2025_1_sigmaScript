package handlers

import (
  "encoding/json"
  "net/http"
  "strconv"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/gorilla/mux"
)

func GetFilm(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    http.Error(w, "Invalid film ID", http.StatusBadRequest)
    return
  }

  film, exists := mocks.Films[id]
  if !exists {
    http.Error(w, "Film not found", http.StatusNotFound)
    return
  }

  if err = json.NewEncoder(w).Encode(film); err != nil {
    http.Error(w, "Encode error", http.StatusInternalServerError)
    return
  }
}
