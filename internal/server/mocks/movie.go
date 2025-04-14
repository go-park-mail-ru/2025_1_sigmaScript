package mocks

type ReviewUserDataJSON struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar,omitempty"`
}

type ReviewJSON struct {
	ID         int                `json:"id"`
	Score      int                `json:"score"`
	ReviewText string             `json:"review_text"`
	CreatedAt  string             `json:"created_at"`
	User       ReviewUserDataJSON `json:"user"`
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
	1: {
		ID:              1,
		Name:            "Бойцовский клуб",
		About:           `Сотрудник страховой компании страдает хронической бессонницей и отчаянно пытается вырваться из мучительно скучной жизни. Однажды в очередной командировке он встречает некоего Тайлера Дёрдена — харизматического торговца мылом с извращенной философией. Тайлер уверен, что самосовершенствование — удел слабых, а единственное, ради чего стоит жить, — саморазрушение. \nПроходит немного времени, и вот уже новые друзья лупят друг друга почем зря на стоянке перед баром, и очищающий мордобой доставляет им высшее блаженство. Приобщая других мужчин к простым радостям физической жестокости, они основывают тайный Бойцовский клуб, который начинает пользоваться невероятной популярностью.`,
		Poster:          "/static/img/0.webp",
		ReleaseYear:     1999,
		Country:         "США",
		Budget:          63000000,
		BoxOfficeUS:     37030102,
		BoxOfficeGlobal: 100853753,
		Rating:          8.2,
		Duration:        "2ч 19м",
		Genres: []GenreJSON{
			{ID: 1, Name: "триллер"},
			{ID: 2, Name: "драма"},
		},
		Staff: []PersonJSON{
			{ID: 7, FullName: "Брэд Питт", Photo: "/static/img/brad_pitt.webp"},
			{ID: 12, FullName: "Эдвард Нортон"},
			{ID: 13, FullName: "Хелена Бонем Картер"},
			{ID: 14, FullName: "Марвин Ли Эдей"},
			{ID: 15, FullName: "Джаред Лето"},
			{ID: 16, FullName: "Зэк Гренье"},
			{ID: 17, FullName: "Холт Маккэллани"},
			{ID: 18, FullName: "Эйон Бэйли"},
			{ID: 19, FullName: "Ричмонд Аркетт"},
			{ID: 20, FullName: "Дэвид Эндрюс"},
		},
		Reviews: []ReviewJSON{
			{
				ID:         1,
				User:       ReviewUserDataJSON{Login: "KinoKritik77"},
				ReviewText: "Абсолютный шедевр! Фильм, который заставляет задуматься о современном обществе, консьюмеризме и поиске себя. Потрясающая игра актеров и неожиданный финал.",
				Score:      10,
				CreatedAt:  "2023-10-15 22:01:16.716787654 +0300 MSK m=+35.697200482",
			},
			{
				ID:         2,
				User:       ReviewUserDataJSON{Login: "Alice_F"},
				ReviewText: "Сначала показался странным и жестоким, но потом поняла глубину. Финал просто взрывает мозг! Пересматривала несколько раз.",
				Score:      9,
				CreatedAt:  "2024-01-20 12:11:00.716787654 +0300 MSK m=+35.697200482",
			},
			{
				ID:         3,
				User:       ReviewUserDataJSON{Login: "Sergey_N"},
				ReviewText: "Не мое. Слишком много неоправданного насилия и псевдофилософии. Пытается быть глубоким, но выглядит претенциозно. Финал предсказуем, если внимательно смотреть.",
				Score:      5,
				CreatedAt:  "2023-11-01 02:27:15.716787654 +0300 MSK m=+35.697200482",
			},
			{
				ID:         4,
				User:       ReviewUserDataJSON{Login: "Tyler_Fan99"},
				ReviewText: "Лучший фильм ЭВЕР! Нортон и Питт на высоте. Идея анархии и разрушения системы - то, что нужно! Первое правило - никому не рассказывать!",
				Score:      10,
				CreatedAt:  "2024-03-08 14:56:19.716787654 +0300 MSK m=+35.697200482",
			},
			{
				ID:         5,
				User:       ReviewUserDataJSON{Login: "RegularViewer"},
				ReviewText: "Интересный фильм с неожиданным поворотом. Хорошая сатира на общество потребления, но местами затянуто. Стоит посмотреть хотя бы раз.",
				Score:      7,
				CreatedAt:  "2023-12-25 16:46:21.610787654 +0300 MSK m=+35.697200482",
			},
		},
	},
	// другие фильмы
	2: {
		ID:              1,
		Name:            "Матрица",
		About:           `Жизнь Томаса Андерсона разделена на две части: днём он — самый обычный офисный работник, получающий нагоняи от начальства, а ночью превращается в хакера по имени Нео, и нет места в сети, куда он бы не смог проникнуть. Но однажды всё меняется. Томас узнаёт ужасающую правду о реальности.`,
		Poster:          "/static/img/7.webp",
		ReleaseYear:     1999,
		Country:         "США, Австралия",
		Budget:          63000000,
		BoxOfficeUS:     171479930,
		BoxOfficeGlobal: 463517383,
		Rating:          8.5,
		Duration:        "2ч 16м",
		Genres: []GenreJSON{
			{ID: 3, Name: "фантастика"},
			{ID: 4, Name: "боевик"},
		},
		Staff: []PersonJSON{
			{ID: 11, FullName: "Киану Ривз", Photo: "https://i.pinimg.com/originals/a3/70/0b/a3700bdf15fcceabf740e1f347dbb5a2.jpg"},
		},
		Reviews: []ReviewJSON{
			{
				ID:         1,
				User:       ReviewUserDataJSON{Login: "KinoKritik77"},
				ReviewText: "Абсолютный шедевр! Потрясающая игра актеров и неожиданный финал.",
				Score:      10,
				CreatedAt:  "2023-12-01 18:11:53.516787654 +0300 MSK m=+35.697200482",
			},
		},
	},
}
