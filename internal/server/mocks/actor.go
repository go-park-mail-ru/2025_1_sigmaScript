package mocks

type Actor struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

var Actors = map[int]Actor{
  1: {ID: 1, Name: "Leonardo DiCaprio"},
  2: {ID: 2, Name: "Brad Pitt"},
}
