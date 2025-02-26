package models

type RegisterData struct {
  Username         string `json:"username"`
  Email            string `json:"email"`
  Password         string `json:"password"`
  RepeatedPassword string `json:"repeated_password"`
}

type LoginData struct {
  Username string `json:"username"`
  Password string `json:"password"`
}
