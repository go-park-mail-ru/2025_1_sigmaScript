package ds

const (
  SuccessfulRegister = "Successfully registered"
  SuccessfulLogin    = "Successfully logged in"
  SuccessfulLogout   = "Successfully logged out"
)

type Response struct {
  Message string `json:"message"`
}
