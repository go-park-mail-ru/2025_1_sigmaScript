package mocks

type Genre struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

var Genres = map[int]Genre{
  1: {ID: 1, Name: "Sci-Fi"},
  2: {ID: 2, Name: "Drama"},
}
