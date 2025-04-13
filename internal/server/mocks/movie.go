package mocks

type UserJSON struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar,omitempty"`
}

type ReviewJSON struct {
	ID         int      `json:"id"`
	Score      int      `json:"score"`
	ReviewText string   `json:"review_text"`
	CreatedAt  string   `json:"created_at"`
	User       UserJSON `json:"user"`
}

type GenreJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// объявлено в staff_person.go, берем оттуда или копию типа здесь оставляем
// type PersonJSON struct {
// 	ID         int    `json:"id"`
// 	FullName   string `json:"full_name"`
// 	EnFullName string `json:"en_full_name,omitempty"`
// 	Photo      string `json:"photo"`
// 	About      string `json:"about"`
// 	Sex        string `json:"sex,omitempty"`
// 	Growth     string `json:"growth,omitempty"`
// 	Birthday   string `json:"birthday,omitempty"`
// 	Death      string `json:"death,omitempty"`
// }

type MovieJSON struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	OriginalName    string       `json:"original_name,omitempty"`
	About           string       `json:"about,omitempty"`
	Poster          string       `json:"poster,omitempty"`
	ReleaseYear     int          `json:"release_year,omitempty"`
	Country         string       `json:"country,omitempty"`
	Slogan          string       `json:"slogan,omitempty"`
	Director        string       `json:"director,omitempty"`
	Budget          int64        `json:"budget,omitempty"`
	BoxOfficeUS     int64        `json:"box_office_us,omitempty"`
	BoxOfficeGlobal int64        `json:"box_office_global,omitempty"`
	BoxOfficeRussia int64        `json:"box_office_russia,omitempty"`
	PremierRussia   string       `json:"premier_russia,omitempty"`
	PremierGlobal   string       `json:"premier_global,omitempty"`
	Rating          float64      `json:"rating,omitempty"`
	Duration        string       `json:"duration,omitempty"`
	Genres          []GenreJSON  `json:"genres,omitempty"`
	Staff           []PersonJSON `json:"staff,omitempty"`
	Reviews         []ReviewJSON `json:"reviews,omitempty"`
}

type Movies map[int]MovieJSON

var ExistingMovies = Movies{
	0: {
		ID:              0,
		Name:            "Бойцовский клуб",
		About:           `Сотрудник страховой компании страдает хронической бессонницей и отчаянно пытается вырваться из мучительно скучной жизни. Однажды в очередной командировке он встречает некоего Тайлера Дёрдена — харизматического торговца мылом с извращенной философией. Тайлер уверен, что самосовершенствование — удел слабых, а единственное, ради чего стоит жить, — саморазрушение. \nПроходит немного времени, и вот уже новые друзья лупят друг друга почем зря на стоянке перед баром, и очищающий мордобой доставляет им высшее блаженство. Приобщая других мужчин к простым радостям физической жестокости, они основывают тайный Бойцовский клуб, который начинает пользоваться невероятной популярностью.`,
		Poster:          "/static/img/0.webp",
		ReleaseYear:     1999,
		Country:         "США",
		Budget:          63000000,
		BoxOfficeUS:     37030102,
		BoxOfficeGlobal: 100853753,
		Rating:          8.8,
		Duration:        "2ч 19м",
		Genres: []GenreJSON{
			{ID: 1, Name: "триллер"},
			{ID: 2, Name: "драма"},
		},
		Staff: []PersonJSON{
			{ID: 1, FullName: "Брэд Питт", Photo: "/static/img/brad_pitt.webp"},
			{ID: 2, FullName: "Эдвард Нортон"},
			{ID: 3, FullName: "Хелена Бонем Картер"},
			{ID: 4, FullName: "Мит Лоуф"},
			{ID: 5, FullName: "Джаред Лето"},
			{ID: 6, FullName: "Зэк Гренье"},
			{ID: 7, FullName: "Холт Маккэллани"},
			{ID: 8, FullName: "Эйон Бэйли"},
			{ID: 9, FullName: "Ричмонд Аркетт"},
			{ID: 10, FullName: "Дэвид Эндрюс"},
		},
		Reviews: []ReviewJSON{
			{
				ID:         1,
				User:       UserJSON{Login: "KinoKritik77"},
				ReviewText: "Абсолютный шедевр! Фильм, который заставляет задуматься о современном обществе, консьюмеризме и поиске себя. Потрясающая игра актеров и неожиданный финал.",
				Score:      10,
				CreatedAt:  "15.10.2023",
			},
			{
				ID:         2,
				User:       UserJSON{Login: "Alice_F"},
				ReviewText: "Сначала показался странным и жестоким, но потом поняла глубину. Финал просто взрывает мозг! Пересматривала несколько раз.",
				Score:      9,
				CreatedAt:  "20.01.2024",
			},
			{
				ID:         3,
				User:       UserJSON{Login: "Sergey_N"},
				ReviewText: "Не мое. Слишком много неоправданного насилия и псевдофилософии. Пытается быть глубоким, но выглядит претенциозно. Финал предсказуем, если внимательно смотреть.",
				Score:      5,
				CreatedAt:  "01.11.2023",
			},
			{
				ID:         4,
				User:       UserJSON{Login: "Tyler_Fan99"},
				ReviewText: "Лучший фильм ЭВЕР! Нортон и Питт на высоте. Идея анархии и разрушения системы - то, что нужно! Первое правило - никому не рассказывать!",
				Score:      10,
				CreatedAt:  "08.03.2024",
			},
			{
				ID:         5,
				User:       UserJSON{Login: "RegularViewer"},
				ReviewText: "Интересный фильм с неожиданным поворотом. Хорошая сатира на общество потребления, но местами затянуто. Стоит посмотреть хотя бы раз.",
				Score:      7,
				CreatedAt:  "25.12.2023",
			},
		},
	},
	// другие фильмы
}
