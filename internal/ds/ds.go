package ds

//go:generate easyjson -all ds.go
type Response struct {
	Message string `json:"message"`
}

type User struct {
	Username string `json:"username"`
}
