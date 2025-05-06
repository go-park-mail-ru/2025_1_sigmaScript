package mocks

type Genre struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies,omitempty"`
}
