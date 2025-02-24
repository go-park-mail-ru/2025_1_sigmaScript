package models

type RegisterData struct {
  Username string `json:"username"`
  Email    string `json:"email"`
  Password string `json:"password"`
}

type LoginData struct {
  Username string `json:"username"`
  Password string `json:"password"`
}
