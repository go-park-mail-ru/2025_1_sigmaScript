package mocks

import "time"

// PersonJSON delivery layer staff person info
type PersonJSON struct {
	ID         int    `json:"id"`
	FullName   string `json:"full_name"`
	EnFullName string `json:"en_full_name,omitempty"`
	Photo      string `json:"photo"`
	About      string `json:"about"`
	Sex        string `json:"sex,omitempty"`
	Growth     string `json:"growth,omitempty"`
	Birthday   string `json:"birthday,omitempty"`
	Death      string `json:"death,omitempty"`
}

type Persons map[int]PersonJSON

var ExistingActors = Persons{
	1:  {ID: 1, FullName: "Леонардо Ди Каприо", EnFullName: "Leonardo DiCaprio", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	2:  {ID: 2, FullName: "Морган Фримен", EnFullName: "Morgan Freeman", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	3:  {ID: 3, FullName: "Том Хэнкс", EnFullName: "Tom Hanks", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	4:  {ID: 4, FullName: "Джонни Депп", EnFullName: "Johnny Depp", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	5:  {ID: 5, FullName: "Том Круз", EnFullName: "Tom Cruise", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	6:  {ID: 6, FullName: "Сэмюэл Л. Джексон", EnFullName: "Samuel L. Jackson", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	7:  {ID: 7, FullName: "Брэд Питт", EnFullName: "Brad Pitt", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	8:  {ID: 8, FullName: "Рассел Кроу", EnFullName: "Russell Crowe", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана", Birthday: time.Now().String()},
	9:  {ID: 9, FullName: "Уилл Смит", EnFullName: "Will Smith", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	10: {ID: 10, FullName: "Мэтт Деймон", EnFullName: "Matt Damon", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
}
