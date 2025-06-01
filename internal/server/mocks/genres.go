package mocks

//go:generate easyjson -all genres.go
type Genre struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies,omitempty"`
}
