package handlers

import (
  "encoding/json"
  "net/http"
  "strconv"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/gorilla/mux"
)

func GetActor(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    http.Error(w, "Invalid actor ID", http.StatusBadRequest)
    return
  }

  actor, exists := mocks.Actors[id]
  if !exists {
    http.Error(w, "Actor not found", http.StatusNotFound)
    return
  }

  if err = json.NewEncoder(w).Encode(actor); err != nil {
    http.Error(w, "Encode error", http.StatusInternalServerError)
    return
  }
}
