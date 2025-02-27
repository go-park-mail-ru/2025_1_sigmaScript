package defaults

import (
  "net/http"
  "time"
)

// server constants
const (
  Address         = "localhost"
  Port            = 8080
  ReadTimeout     = time.Second * 5
  WriteTimeout    = time.Second * 5
  ShutdownTimeout = time.Second * 30
  IdleTimeout     = time.Second * 60
)

// cookie constants
const (
  SessionName   = "session_id"
  SessionLength = 32
  HTTPOnly      = true
  Secure        = false
  SameSite      = http.SameSiteStrictMode
  Path          = "/"
  ResetMaxAge   = -1
)
