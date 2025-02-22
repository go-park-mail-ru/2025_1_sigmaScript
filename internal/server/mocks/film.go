package mocks

type Film struct {
  ID       int    `json:"id"`
  Title    string `json:"title"`
  GenreID  int    `json:"genre_id"`
  ActorIDs []int  `json:"actor_ids"`
}

var Films = map[int]Film{
  1: {ID: 1, Title: "Inception", GenreID: 1, ActorIDs: []int{1}},
  2: {ID: 2, Title: "Fight Club", GenreID: 2, ActorIDs: []int{2}},
}
