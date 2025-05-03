package mocks

// PersonJSON delivery layer staff person info
type PersonJSON struct {
	ID              int        `json:"id"`
	FullName        string     `json:"full_name"`
	EnFullName      string     `json:"en_full_name,omitempty"`
	Photo           string     `json:"photo"`
	About           string     `json:"about"`
	Sex             string     `json:"sex,omitempty"`
	Growth          string     `json:"growth,omitempty"`
	Birthday        string     `json:"birthday,omitempty"`
	Death           string     `json:"death,omitempty"`
	Career          string     `json:"career,omitempty"`
	Genres          string     `json:"genres,omitempty"`
	TotalFilms      string     `json:"total_films,omitempty"`
	MovieCollection Collection `json:"movie_collection,omitempty"`
}

type Persons map[int]PersonJSON

var ExistingActors = Persons{
	1:  {ID: 1, FullName: "Леонардо Ди Каприо", EnFullName: "Leonardo DiCaprio", Photo: "https://avatars.mds.yandex.net/get-entity_search/2310675/1130394491/S600xU_2x", About: "Информация по этому человеку не указана"},
	2:  {ID: 2, FullName: "Морган Фримен", EnFullName: "Morgan Freeman", Photo: "https://avatars.mds.yandex.net/get-entity_search/2057552/1132084397/S600xU_2x", About: "Информация по этому человеку не указана"},
	3:  {ID: 3, FullName: "Том Хэнкс", EnFullName: "Tom Hanks", Photo: "https://avatars.mds.yandex.net/get-entity_search/2005770/833182325/S600xU_2x", About: "Информация по этому человеку не указана"},
	4:  {ID: 4, FullName: "Джонни Депп", EnFullName: "Johnny Depp", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	5:  {ID: 5, FullName: "Том Круз", EnFullName: "Tom Cruise", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	6:  {ID: 6, FullName: "Сэмюэл Л. Джексон", EnFullName: "Samuel L. Jackson", Photo: "https://avatars.mds.yandex.net/get-entity_search/98180/952678918/S600xU_2x", About: "Информация по этому человеку не указана"},
	7:  {ID: 7, FullName: "Брэд Питт", EnFullName: "Brad Pitt", Photo: "/static/img/brad_pitt.webp", About: "Информация по этому человеку не указана"},
	8:  {ID: 8, FullName: "Рассел Кроу", EnFullName: "Russell Crowe", Photo: "https://avatars.mds.yandex.net/get-entity_search/478647/809836058/S600xU_2x", About: "Информация по этому человеку не указана", Birthday: "2010-04-10 13:39:11.041078099 +0300 MSK m=+0.000049449"},
	9:  {ID: 9, FullName: "Уилл Смит", EnFullName: "Will Smith", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	10: {ID: 10, FullName: "Мэтт Дэймон", EnFullName: "Matt Damon", Photo: "https://avatars.mds.yandex.net/get-entity_search/1245892/935872902/S600xU_2x", About: "Информация по этому человеку не указана"},
	11: {
		ID:         11,
		FullName:   "Киану Ривз",
		EnFullName: "Keanu Reeves",
		Photo:      "https://i.pinimg.com/originals/a3/70/0b/a3700bdf15fcceabf740e1f347dbb5a2.jpg",
		Career:     "Актёр, Продюсер, Режиссёр",
		Growth:     "186",
		Sex:        "Мужчина",
		Birthday:   "1964-09-2 0:0:0.041078099 +0300 MSK m=+0.000049449",
		Genres:     "Боевик, Фантастика, Триллер",
		TotalFilms: "274, 1981-2025",
		About: `
Киану Чарльз Ривз — канадский актёр, кинорежиссёр, кинопродюсер и музыкант.
Наиболее известен своими ролями в киносериях «Матрица», «Билл и Тед», «Джон Уик», а также в фильмах «На гребне волны», «Скорость», «Адвокат дьявола», «Константин: Повелитель тьмы».
Обладатель звезды на Голливудской «Аллее славы».`,
	},
	12: {ID: 12, FullName: "Эдвард Нортон", EnFullName: "Edward Norton", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	13: {ID: 13, FullName: "Хелена Бонем Картер", EnFullName: "Helena Bonham Carter", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	14: {ID: 14, FullName: "Мит Лоуф", EnFullName: "Meat Loaf", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	15: {ID: 15, FullName: "Джаред Лето", EnFullName: "Jared Leto", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	16: {ID: 16, FullName: "Зэк Гренье", EnFullName: "Zack Grenier", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	17: {ID: 17, FullName: "Холт Маккэллани", EnFullName: "Holt McCallany", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	18: {ID: 18, FullName: "Эйон Бэйли", EnFullName: "Eion Bailey", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	19: {ID: 19, FullName: "Ричмонд Аркетт", EnFullName: "Richmond Arquette", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	20: {ID: 20, FullName: "Дэвид Эндрюс", EnFullName: "David Andrews", Photo: "/static/avatars/avatar_default_picture.svg", About: "Информация по этому человеку не указана"},
	21: {ID: 21, FullName: "Аль Пачино", EnFullName: "Al Pacino", Photo: "https://avatars.mds.yandex.net/get-entity_search/516537/726983645/SUx182_2x", About: "Информация по этому человеку не указана"},
	22: {ID: 22, FullName: "Мэттью Макконахи", EnFullName: "Matthew David McConaughey", Photo: "https://avatars.mds.yandex.net/get-entity_search/1987348/921017284/S600xU_2x", About: "Информация по этому человеку не указана"},
	23: {ID: 23, FullName: "Брюс Уиллис", EnFullName: "Bruce Willis", Photo: "https://avatars.mds.yandex.net/get-entity_search/1734588/930928060/S600xU_2x", About: "Информация по этому человеку не указана"},
	24: {ID: 24, FullName: "Джон Траволта", EnFullName: "John Joseph Travolta", Photo: "https://avatars.mds.yandex.net/get-entity_search/2300207/1131324291/S600xU_2x", About: "Информация по этому человеку не указана"},
	25: {ID: 25, FullName: "Тим Роббинс", EnFullName: "Tim Robbins", Photo: "https://avatars.mds.yandex.net/get-entity_search/60958/784029965/S600xU_2x", About: "Информация по этому человеку не указана"},
	26: {ID: 26, FullName: "Кристиан Бэйл", EnFullName: "Christian Charles Philip Bale", Photo: "https://avatars.mds.yandex.net/get-entity_search/141104/728910606/S600xU_2x", About: "Информация по этому человеку не указана"},
	27: {ID: 27, FullName: "Хит Леджер", EnFullName: "Heathcliff Andrew Ledger", Photo: "https://avatars.mds.yandex.net/get-entity_search/935097/952763126/S600xU_2x", About: "Информация по этому человеку не указана"},
	28: {ID: 28, FullName: "Майкл Кларк Дункан", EnFullName: "Michael Clarke Duncan", Photo: "https://avatars.mds.yandex.net/get-entity_search/5578840/826293398/S600xU_2x", About: "Информация по этому человеку не указана"},
	29: {ID: 29, FullName: "Дэвид Морс", EnFullName: "David Morse", Photo: "https://avatars.mds.yandex.net/get-entity_search/10105370/1132643699/S600xU_2x", About: "Информация по этому человеку не указана"},
	30: {ID: 30, FullName: "Майлз Теллер", EnFullName: "Miles Teller", Photo: "https://avatars.mds.yandex.net/get-entity_search/95107/918423255/S600xU_2x", About: "Информация по этому человеку не указана"},
	31: {ID: 31, FullName: "Роберт Дауни мл.", EnFullName: "Robert Downey Jr.", Photo: "https://avatars.mds.yandex.net/get-entity_search/1579191/725945074/S600xU_2x", About: "Информация по этому человеку не указана"},
	32: {ID: 32, FullName: "Киллиан Мерфи", EnFullName: "Cillian Murphy", Photo: "https://avatars.mds.yandex.net/get-entity_search/2273637/1132284892/S600xU_2x", About: "Информация по этому человеку не указана"},
	33: {ID: 33, FullName: "Эмили Блант", EnFullName: "Emily Olivia Laura Blunt", Photo: "https://avatars.mds.yandex.net/get-entity_search/5508932/1131749885/S600xU_2x", About: "Информация по этому человеку не указана"},
	34: {ID: 34, FullName: "Марк Хэмилл", EnFullName: "Mark Hamill", Photo: "https://avatars.mds.yandex.net/get-entity_search/1634327/991431229/S600xU_2x", About: "Информация по этому человеку не указана"},
	35: {ID: 35, FullName: "Харрисон Форд", EnFullName: "Harrison Ford", Photo: "https://avatars.mds.yandex.net/get-entity_search/4740766/953026487/S600xU_2x", About: "Информация по этому человеку не указана"},
	36: {ID: 36, FullName: "Сильвестр Сталлоне", EnFullName: "Sylvester Stallone", Photo: "https://avatars.mds.yandex.net/get-entity_search/5449393/1132278994/S600xU_2x", About: "Информация по этому человеку не указана"},
	37: {ID: 37, FullName: "Хоакин Феникс", EnFullName: "Joaquin Phoenix", Photo: "https://i.pinimg.com/736x/1f/25/37/1f2537bd8057c6ee1115a5ab23b486b4.jpg", About: "Информация по этому человеку не указана"},
	38: {ID: 38, FullName: "Бенедикт Камбербэтч", EnFullName: "Benedict Cumberbatch", Photo: "https://avatars.mds.yandex.net/get-entity_search/2162801/1132456058/S600xU_2x", About: "Информация по этому человеку не указана"},
	39: {ID: 39, FullName: "Майкл Дж. Фокс", EnFullName: "Michael J. Fox", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1704946/9fe88f52-8452-4d38-b698-536b1af91439/600x900", About: "Информация по этому человеку не указана"},
	40: {ID: 40, FullName: "Юра Борисов", EnFullName: "Ura Borisov", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1898899/dc5079bf-795c-4953-a8c3-e3a5a025a1ff/600x900", About: "Информация по этому человеку не указана"},
	41: {ID: 41, FullName: "Кирилл Зайцев", EnFullName: "Kirill Zaitzev", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1777765/5833d6a1-c63f-4d86-be47-b476746f4d1c/600x900", About: "Информация по этому человеку не указана"},
	42: {ID: 42, FullName: "Владимир Вдовиченков", EnFullName: "Vladimir Vdovichenkov", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1777765/34a084da-0eb2-4fbd-af20-99cedf91de2a/600x900", About: "Информация по этому человеку не указана"},
}
