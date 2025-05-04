package mocks

type ReviewUserDataJSON struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar,omitempty"`
}

type ReviewJSON struct {
	ID         int                `json:"id"`
	Score      float64            `json:"score"`
	ReviewText string             `json:"review_text"`
	CreatedAt  string             `json:"created_at"`
	User       ReviewUserDataJSON `json:"user"`
}

type NewReviewDataJSON struct {
	ReviewText string  `json:"review_text,omitempty"`
	Score      float64 `json:"score"`
}

type GenreJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MovieJSON struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	OriginalName    string       `json:"original_name,omitempty"`
	About           string       `json:"about,omitempty"`
	Poster          string       `json:"poster,omitempty"`
	ReleaseYear     string       `json:"release_year,omitempty"`
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
	Genres          string       `json:"genres,omitempty"`
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
		ReleaseYear:     "1999",
		Country:         "США",
		Budget:          63000000,
		BoxOfficeUS:     37030102,
		BoxOfficeGlobal: 100853753,
		Rating:          8.2,
		Duration:        "2ч 19м",
		Genres:          "thriller, " + "drama",
		Staff: []PersonJSON{
			{ID: 7, FullName: "Брэд Питт", Photo: "/static/img/brad_pitt.webp"},
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
	2: {
		ID:              2,
		Name:            "Матрица",
		About:           `Жизнь Томаса Андерсона разделена на две части: днём он — самый обычный офисный работник, получающий нагоняи от начальства, а ночью превращается в хакера по имени Нео, и нет места в сети, куда он бы не смог проникнуть. Но однажды всё меняется. Томас узнаёт ужасающую правду о реальности.`,
		Poster:          "/static/img/7.webp",
		ReleaseYear:     "1999",
		Country:         "США, Австралия",
		Budget:          63000000,
		BoxOfficeUS:     171479930,
		BoxOfficeGlobal: 463517383,
		Rating:          8.5,
		Duration:        "2ч 16м",
		Genres:          "fantastic, " + "action movie",
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

	3: {
		ID:   3,
		Name: "Форрест Гамп",
		About: `Сидя на автобусной остановке, Форрест Гамп — не очень умный, но добрый и открытый парень — рассказывает случайным встречным историю своей необыкновенной жизни.

С самого малолетства парень страдал от заболевания ног, соседские мальчишки дразнили его, но в один прекрасный день Форрест открыл в себе невероятные способности к бегу. Подруга детства Дженни всегда его поддерживала и защищала, но вскоре дороги их разошлись.`,
		Poster:          "/static/img/2.webp",
		ReleaseYear:     "1994",
		Country:         "США",
		Budget:          55000000,
		BoxOfficeUS:     329694499,
		BoxOfficeGlobal: 677387716,
		Rating:          8.9,
		Duration:        "2ч 22м",
		Genres:          "драма, " + "комедия, " + "мелодрама, " + "история, " + "военный",
		Staff: []PersonJSON{
			{ID: 3, FullName: "Том Хэнкс", EnFullName: "Tom Hanks", Photo: "https://avatars.mds.yandex.net/get-entity_search/2005770/833182325/S600xU_2x"},
		},
		Reviews: []ReviewJSON{},
	},

	4: {
		ID:   4,
		Name: "Крестный отец",
		About: `Криминальная сага, повествующая о нью-йоркской сицилийской мафиозной семье Корлеоне. Фильм охватывает период 1945-1955 годов.

Глава семьи, Дон Вито Корлеоне, выдаёт замуж свою дочь. В это время со Второй мировой войны возвращается его любимый сын Майкл. Майкл, герой войны, гордость семьи, не выражает желания заняться жестоким семейным бизнесом. Дон Корлеоне ведёт дела по старым правилам, но наступают иные времена, и появляются люди, желающие изменить сложившиеся порядки. На Дона Корлеоне совершается покушение.`,
		Poster:          "/static/img/3.webp",
		ReleaseYear:     "1972",
		Country:         "США",
		Budget:          6000000,
		BoxOfficeUS:     133698921,
		BoxOfficeGlobal: 243862778,
		Rating:          8.7,
		Duration:        "2ч 55м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 9, Name: "криминал"},
		}`,
		Staff: []PersonJSON{
			{ID: 21, FullName: "Аль Пачино", EnFullName: "Al Pacino", Photo: "https://avatars.mds.yandex.net/get-entity_search/516537/726983645/SUx182_2x"},
		},
		Reviews: []ReviewJSON{},
	},

	5: {
		ID:              5,
		Name:            "Интерстеллар",
		About:           `Когда засуха, пыльные бури и вымирание растений приводят человечество к продовольственному кризису, коллектив исследователей и учёных отправляется сквозь червоточину (которая предположительно соединяет области пространства-времени через большое расстояние) в путешествие, чтобы превзойти прежние ограничения для космических путешествий человека и найти планету с подходящими для человечества условиями.`,
		Poster:          "/static/img/4.webp",
		ReleaseYear:     "2014",
		Country:         "США, Великобритания, Канада",
		Budget:          165000000,
		BoxOfficeUS:     192445017,
		BoxOfficeGlobal: 736546575,
		Rating:          8.7,
		Duration:        "2ч 49м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 3, Name: "фантастика"},
			{ID: 10, Name: "приключения"},
		}`,
		Staff: []PersonJSON{
			{ID: 22, FullName: "Мэттью Макконахи", EnFullName: "Matthew David McConaughey", Photo: "https://avatars.mds.yandex.net/get-entity_search/1987348/921017284/S600xU_2x"},
			{ID: 10, FullName: "Мэтт Дэймон", EnFullName: "Matt Damon", Photo: "https://avatars.mds.yandex.net/get-entity_search/1245892/935872902/S600xU_2x"},
		},
		Reviews: []ReviewJSON{},
	},

	6: {
		ID:   6,
		Name: "Криминальное чтиво",
		About: `Двое бандитов Винсент Вега и Джулс Винфилд ведут философские беседы в перерывах между разборками и решением проблем с должниками криминального босса Марселласа Уоллеса.

В первой истории Винсент проводит незабываемый вечер с женой Марселласа Мией. Во второй Марселлас покупает боксёра Бутча Кулиджа, чтобы тот сдал бой. В третьей истории Винсент и Джулс по нелепой случайности попадают в неприятности.`,
		Poster:          "/static/img/5.webp",
		ReleaseYear:     "1994",
		Country:         "США",
		Budget:          8000000,
		BoxOfficeUS:     107928762,
		BoxOfficeGlobal: 213928762,
		Rating:          8.7,
		Duration:        "2ч 34м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 9, Name: "криминал"},
		}`,
		Staff: []PersonJSON{
			{ID: 6, FullName: "Сэмюэл Л. Джексон", EnFullName: "Samuel L. Jackson", Photo: "https://avatars.mds.yandex.net/get-entity_search/98180/952678918/S600xU_2x"},
			{ID: 23, FullName: "Брюс Уиллис", EnFullName: "Bruce Willis", Photo: "https://avatars.mds.yandex.net/get-entity_search/1734588/930928060/S600xU_2x"},
			{ID: 24, FullName: "Джон Траволта", EnFullName: "John Joseph Travolta", Photo: "https://avatars.mds.yandex.net/get-entity_search/2300207/1131324291/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	7: {
		ID:              7,
		Name:            "Побег из Шоушенка",
		About:           `Бухгалтер Энди Дюфрейн обвинён в убийстве собственной жены и её любовника. Оказавшись в тюрьме под названием Шоушенк, он сталкивается с жестокостью и беззаконием, царящими по обе стороны решётки. Каждый, кто попадает в эти стены, становится их рабом до конца жизни. Но Энди, обладающий живым умом и доброй душой, находит подход как к заключённым, так и к охранникам, добиваясь их особого к себе расположения.`,
		Poster:          "/static/img/5.webp",
		ReleaseYear:     "1994",
		Country:         "США",
		Budget:          25000000,
		BoxOfficeUS:     28341469,
		BoxOfficeGlobal: 28418687,
		Rating:          9.1,
		Duration:        "2ч 22м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
		}`,
		Staff: []PersonJSON{
			{ID: 2, FullName: "Морган Фримен", EnFullName: "Morgan Freeman", Photo: "https://avatars.mds.yandex.net/get-entity_search/2057552/1132084397/S600xU_2x"},
			{ID: 25, FullName: "Тим Роббинс", EnFullName: "Tim Robbins", Photo: "https://avatars.mds.yandex.net/get-entity_search/60958/784029965/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	8: {
		ID:              8,
		Name:            "Тёмный рыцарь",
		About:           `Бэтмен поднимает ставки в войне с криминалом. С помощью лейтенанта Джима Гордона и прокурора Харви Дента он намерен очистить улицы Готэма от преступности. Сотрудничество оказывается эффективным, но скоро они обнаружат себя посреди хаоса, развязанного восходящим криминальным гением, известным напуганным горожанам под именем Джокер.`,
		Poster:          "/static/img/1.webp",
		ReleaseYear:     "2008",
		Country:         "США, Великобритания",
		Budget:          185000000,
		BoxOfficeUS:     533345358,
		BoxOfficeGlobal: 1003045358,
		Rating:          8.5,
		Duration:        "2ч 32м",
		Genres: `[]GenreJSON{
			{ID: 3, Name: "фантастика"},
			{ID: 4, Name: "боевик"},
			{ID: 1, Name: "триллер"},
			{ID: 9, Name: "криминал"},
			{ID: 2, Name: "драма"},
		}`,
		Staff: []PersonJSON{
			{ID: 26, FullName: "Кристиан Бэйл", EnFullName: "Christian Charles Philip Bale", Photo: "https://avatars.mds.yandex.net/get-entity_search/141104/728910606/S600xU_2x", About: "Информация по этому человеку не указана"},
			{ID: 2, FullName: "Морган Фримен", EnFullName: "Morgan Freeman", Photo: "https://avatars.mds.yandex.net/get-entity_search/2057552/1132084397/S600xU_2x"},
			{ID: 27, FullName: "Хит Леджер", EnFullName: "Heathcliff Andrew Ledger", Photo: "https://avatars.mds.yandex.net/get-entity_search/935097/952763126/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	9: {
		ID:              9,
		Name:            "Зелёная миля",
		About:           `Пол Эджкомб — начальник блока смертников в тюрьме «Холодная гора», каждый из узников которого однажды проходит «зеленую милю» по пути к месту казни. Пол повидал много заключённых и надзирателей за время работы. Однако гигант Джон Коффи, обвинённый в страшном преступлении, стал одним из самых необычных обитателей блока.`,
		Poster:          "/static/img/8.webp",
		ReleaseYear:     "1999",
		Country:         "США",
		Budget:          60000000,
		BoxOfficeUS:     136801374,
		BoxOfficeGlobal: 286801374,
		Rating:          9.1,
		Duration:        "3ч 9м",
		Genres: `[]GenreJSON{
			{ID: 9, Name: "криминал"},
			{ID: 2, Name: "драма"},
			{ID: 11, Name: "фэнтези"},
		}`,
		Staff: []PersonJSON{
			{ID: 28, FullName: "Майкл Кларк Дункан", EnFullName: "Michael Clarke Duncan", Photo: "https://avatars.mds.yandex.net/get-entity_search/5578840/826293398/S600xU_2x", About: "Информация по этому человеку не указана"},
			{ID: 3, FullName: "Том Хэнкс", EnFullName: "Tom Hanks", Photo: "https://avatars.mds.yandex.net/get-entity_search/2005770/833182325/S600xU_2x"},
			{ID: 29, FullName: "Дэвид Морс", EnFullName: "David Morse", Photo: "https://avatars.mds.yandex.net/get-entity_search/10105370/1132643699/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	10: {
		ID:              10,
		Name:            "Одержимость",
		About:           `Эндрю мечтает стать великим. Казалось бы, вот-вот его мечта осуществится. Юношу замечает настоящий гений, дирижер лучшего в стране оркестра. Желание Эндрю добиться успеха быстро становится одержимостью, а безжалостный наставник продолжает подталкивать его все дальше и дальше – за пределы человеческих возможностей. Кто выйдет победителем из этой схватки?`,
		Poster:          "/static/img/9.webp",
		ReleaseYear:     "2013",
		Country:         "США",
		Budget:          3300000,
		BoxOfficeUS:     13092000,
		BoxOfficeGlobal: 48982041,
		Rating:          8.4,
		Duration:        "1ч 46м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 12, Name: "музыка"},
		}`,
		Staff: []PersonJSON{
			{ID: 30, FullName: "Майлз Теллер", EnFullName: "Miles Teller", Photo: "https://avatars.mds.yandex.net/get-entity_search/95107/918423255/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	11: {
		ID:              11,
		Name:            "Оппенгеймер",
		About:           `История жизни американского физика-теоретика Роберта Оппенгеймера, который во времена Второй мировой войны руководил Манхэттенским проектом — секретными разработками ядерного оружия.`,
		Poster:          "/static/img/11.webp",
		ReleaseYear:     "2023",
		Country:         "США",
		Budget:          100000000,
		BoxOfficeUS:     330078895,
		BoxOfficeGlobal: 975811333,
		Rating:          8.1,
		Duration:        "3ч",
		Genres: `[]GenreJSON{
			{ID: 13, Name: "биография"},
			{ID: 2, Name: "драма"},
			{ID: 7, Name: "история"},
		}`,
		Staff: []PersonJSON{
			{ID: 31, FullName: "Роберт Дауни мл.", EnFullName: "Robert Downey Jr.", Photo: "https://avatars.mds.yandex.net/get-entity_search/1579191/725945074/S600xU_2x", About: "Информация по этому человеку не указана"},
			{ID: 32, FullName: "Киллиан Мерфи", EnFullName: "Cillian Murphy", Photo: "https://avatars.mds.yandex.net/get-entity_search/2273637/1132284892/S600xU_2x", About: "Информация по этому человеку не указана"},
			{ID: 33, FullName: "Эмили Блант", EnFullName: "Emily Olivia Laura Blunt", Photo: "https://avatars.mds.yandex.net/get-entity_search/5508932/1131749885/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	12: {
		ID:              12,
		Name:            "Звёздные войны: Эпизод 4 – Новая надежда",
		About:           `Татуин. Планета-пустыня. Уже постаревший рыцарь Джедай Оби Ван Кеноби спасает молодого Люка Скайуокера, когда тот пытается отыскать пропавшего дроида. С этого момента Люк осознает свое истинное назначение: он один из рыцарей Джедай. В то время как гражданская война охватила галактику, а войска повстанцев ведут бои против сил злого Императора, к Люку и Оби Вану присоединяется отчаянный пилот-наемник Хан Соло, и в сопровождении двух дроидов, R2D2 и C-3PO, этот необычный отряд отправляется на поиски предводителя повстанцев – принцессы Леи. Героям предстоит отчаянная схватка с устрашающим Дартом Вейдером – правой рукой Императора и его секретным оружием – «Звездой Смерти».`,
		Poster:          "/static/img/12.webp",
		ReleaseYear:     "1977",
		Country:         "США",
		Budget:          11000000,
		BoxOfficeUS:     307263857,
		BoxOfficeGlobal: 503015849,
		Rating:          8.1,
		Duration:        "2ч 1м",
		Genres: `[]GenreJSON{
			{ID: 3, Name: "фантастика"},
			{ID: 4, Name: "боевик"},
			{ID: 11, Name: "фэнтези"},
			{ID: 10, Name: "приключения"},
		}`,
		Staff: []PersonJSON{
			{ID: 34, FullName: "Марк Хэмилл", EnFullName: "Mark Hamill", Photo: "https://avatars.mds.yandex.net/get-entity_search/1634327/991431229/S600xU_2x", About: "Информация по этому человеку не указана"},
			{ID: 35, FullName: "Харрисон Форд", EnFullName: "Harrison Ford", Photo: "https://avatars.mds.yandex.net/get-entity_search/4740766/953026487/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	13: {
		ID:              13,
		Name:            "Рокки",
		About:           `Филадельфия. Рокки Бальбоа — молодой боксёр, который живёт в захудалой квартирке и еле сводит концы с концами, занимаясь выбиванием долгов для своего босса Тони Гаццо и периодически участвуя в боях. Каждый его унылый день похож на предыдущий, и особо радужных перспектив не наблюдается. Но однажды удача наконец улыбается парню, когда ему поступает неожиданное предложение выступить против действующего чемпиона Аполло Крида.`,
		Poster:          "/static/img/13.webp",
		ReleaseYear:     "1976",
		Country:         "США",
		Budget:          1100000,
		BoxOfficeUS:     117235147,
		BoxOfficeGlobal: 225000000,
		Rating:          8.0,
		Duration:        "2ч",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 14, Name: "спорт"},
		}`,
		Staff: []PersonJSON{
			{ID: 36, FullName: "Сильвестр Сталлоне", EnFullName: "Sylvester Stallone", Photo: "https://avatars.mds.yandex.net/get-entity_search/5449393/1132278994/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	14: {
		ID:              14,
		Name:            "Джокер",
		About:           `Готэм, начало 1980-х годов. Комик Артур Флек живет с больной матерью, которая с детства учит его «ходить с улыбкой». Пытаясь нести в мир хорошее и дарить людям радость, Артур сталкивается с человеческой жестокостью и постепенно приходит к выводу, что этот мир получит от него не добрую улыбку, а ухмылку злодея Джокера.`,
		Poster:          "/static/img/14.webp",
		ReleaseYear:     "2019",
		Country:         "США, Канада, Австралия",
		Budget:          55000000,
		BoxOfficeUS:     335451311,
		BoxOfficeGlobal: 1078751311,
		Rating:          8.0,
		Duration:        "2ч 2м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 9, Name: "криминал"},
			{ID: 1, Name: "триллер"},
		}`,
		Staff: []PersonJSON{
			{ID: 37, FullName: "Хоакин Феникс", EnFullName: "Joaquin Phoenix", Photo: "https://i.pinimg.com/736x/1f/25/37/1f2537bd8057c6ee1115a5ab23b486b4.jpg", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	15: {
		ID:              15,
		Name:            "Игра в имитацию",
		About:           `Английский математик и логик Алан Тьюринг пытается взломать код немецкой шифровальной машины Enigma во время Второй мировой войны.`,
		Poster:          "/static/img/15.webp",
		ReleaseYear:     "2014",
		Country:         "США",
		Budget:          14000000,
		BoxOfficeUS:     91125683,
		BoxOfficeGlobal: 233555708,
		Rating:          7.8,
		Duration:        "1ч 54м",
		Genres: `[]GenreJSON{
			{ID: 13, Name: "биография"},
			{ID: 8, Name: "военный"},
			{ID: 2, Name: "драма"},
			{ID: 7, Name: "история"},
		}`,
		Staff: []PersonJSON{
			{ID: 38, FullName: "Бенедикт Камбербэтч", EnFullName: "Benedict Cumberbatch", Photo: "https://avatars.mds.yandex.net/get-entity_search/2162801/1132456058/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	16: {
		ID:   16,
		Name: "Начало",
		About: `Кобб – талантливый вор, лучший из лучших в опасном искусстве извлечения: он крадет ценные секреты из глубин подсознания во время сна, когда человеческий разум наиболее уязвим. Редкие способности Кобба сделали его ценным игроком в привычном к предательству мире промышленного шпионажа, но они же превратили его в извечного беглеца и лишили всего, что он когда-либо любил.

И вот у Кобба появляется шанс исправить ошибки. Его последнее дело может вернуть все назад, но для этого ему нужно совершить невозможное – инициацию. Вместо идеальной кражи Кобб и его команда спецов должны будут провернуть обратное. Теперь их задача – не украсть идею, а внедрить ее. Если у них получится, это и станет идеальным преступлением.

Но никакое планирование или мастерство не могут подготовить команду к встрече с опасным противником, который, кажется, предугадывает каждый их ход. Врагом, увидеть которого мог бы лишь Кобб.`,
		Poster:          "/static/img/16.webp",
		ReleaseYear:     "2010",
		Country:         "США, Великобритания",
		Budget:          160000000,
		BoxOfficeUS:     292576195,
		BoxOfficeGlobal: 828322032,
		Rating:          8.7,
		Duration:        "2ч 28м",
		Genres: `[]GenreJSON{
			{ID: 3, Name: "фантастика"},
			{ID: 4, Name: "боевик"},
			{ID: 1, Name: "триллер"},
			{ID: 2, Name: "драма"},
			{ID: 15, Name: "детектив"},
		}`,
		Staff: []PersonJSON{
			{ID: 1, FullName: "Леонардо Ди Каприо", EnFullName: "Leonardo DiCaprio", Photo: "https://avatars.mds.yandex.net/get-entity_search/2310675/1130394491/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	17: {
		ID:              17,
		Name:            "Назад в будущее",
		About:           `Подросток Марти с помощью машины времени, сооружённой его другом-профессором доком Брауном, попадает из 80-х в далекие 50-е. Там он встречается со своими будущими родителями, ещё подростками, и другом-профессором, совсем молодым.`,
		Poster:          "/static/img/17.webp",
		ReleaseYear:     "1985",
		Country:         "США",
		Budget:          19000000,
		BoxOfficeUS:     210609762,
		BoxOfficeGlobal: 381109762,
		Rating:          8.6,
		Duration:        "1ч 56м",
		Genres: `[]GenreJSON{
			{ID: 3, Name: "фантастика"},
			{ID: 5, Name: "комедия"},
			{ID: 10, Name: "приключения"},
		}`,
		Staff: []PersonJSON{
			{ID: 39, FullName: "Майкл Дж. Фокс", EnFullName: "Michael J. Fox", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1704946/9fe88f52-8452-4d38-b698-536b1af91439/600x900", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	18: {
		ID:              18,
		Name:            "Гладиатор",
		About:           `Римская империя. Бесстрашного и благородного генерала Максимуса боготворят солдаты, а старый император Марк Аврелий безгранично доверяет ему и относится как к сыну. Однако опытный воин, готовый сразиться с любым противником в честном бою, оказывается бессильным перед коварными придворными интригами. Коммод, сын Марка Аврелия, убивает отца, который планировал сделать преемником не его, а Максимуса, и захватывает власть. Решив избавиться от опасного соперника, который к тому же отказывается присягнуть ему на верность, Коммод отдаёт приказ убить Максимуса и всю его семью. Чудом выжив, но не сумев спасти близких, Максимус попадает в плен к работорговцу, который продаёт его организатору гладиаторских боёв Проксимо. Так легендарный полководец становится гладиатором. Но вскоре ему представится шанс встретиться со своим смертельным врагом лицом к лицу.`,
		Poster:          "/static/img/18.webp",
		ReleaseYear:     "2000",
		Country:         "США, Великобритания, Мальта, Марокко",
		Budget:          103000000,
		BoxOfficeUS:     187705427,
		BoxOfficeGlobal: 460583960,
		Rating:          8.6,
		Duration:        "2ч 35м",
		Genres: `[]GenreJSON{
			{ID: 7, Name: "история"},
			{ID: 4, Name: "боевик"},
			{ID: 2, Name: "драма"},
		}`,
		Staff: []PersonJSON{
			{ID: 37, FullName: "Хоакин Феникс", EnFullName: "Joaquin Phoenix", Photo: "https://i.pinimg.com/736x/1f/25/37/1f2537bd8057c6ee1115a5ab23b486b4.jpg", About: "Информация по этому человеку не указана"},
			{ID: 8, FullName: "Рассел Кроу", EnFullName: "Russell Crowe", Photo: "https://avatars.mds.yandex.net/get-entity_search/478647/809836058/S600xU_2x"},
		},
		Reviews: []ReviewJSON{},
	},

	19: {
		ID:              19,
		Name:            "Титаник",
		About:           `Апрель 1912 года. В первом и последнем плавании шикарного «Титаника» встречаются двое. Пассажир нижней палубы Джек выиграл билет в карты, а богатая наследница Роза отправляется в Америку, чтобы выйти замуж по расчёту. Чувства молодых людей только успевают расцвести, и даже не классовые различия создадут испытания влюблённым, а айсберг, вставший на пути считавшегося непотопляемым лайнера.`,
		Poster:          "/static/img/19.webp",
		ReleaseYear:     "1997",
		Country:         "США, Мексика",
		Budget:          200000000,
		BoxOfficeUS:     674354882,
		BoxOfficeGlobal: 2264805579,
		Rating:          8.4,
		Duration:        "3ч 14м",
		Genres: `[]GenreJSON{
			{ID: 6, Name: "мелодрама"},
			{ID: 7, Name: "история"},
			{ID: 1, Name: "триллер"},
			{ID: 2, Name: "драма"},
		}`,
		Staff: []PersonJSON{
			{ID: 1, FullName: "Леонардо Ди Каприо", EnFullName: "Leonardo DiCaprio", Photo: "https://avatars.mds.yandex.net/get-entity_search/2310675/1130394491/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	20: {
		ID:              20,
		Name:            "Ford против Ferrari",
		About:           `В начале 1960-х Генри Форд II принимает решение улучшить имидж компании и сменить курс на производство более модных автомобилей. После неудавшейся попытки купить практически банкрота Ferrari американцы решают бросить вызов итальянским конкурентам на трассе и выиграть престижную гонку 24 часа Ле-Мана. Чтобы создать подходящую машину, компания нанимает автоконструктора Кэррола Шэлби, а тот отказывается работать без выдающегося, но, как считается, трудного в общении гонщика Кена Майлза. Вместе они принимаются за разработку впоследствии знаменитого спорткара Ford GT40.`,
		Poster:          "/static/img/10.webp",
		ReleaseYear:     "2019",
		Country:         "США",
		Budget:          97600000,
		BoxOfficeUS:     117624357,
		BoxOfficeGlobal: 225508210,
		Rating:          8.2,
		Duration:        "2ч 32м",
		Genres: `[]GenreJSON{
			{ID: 14, Name: "спорт"},
			{ID: 13, Name: "биография"},
			{ID: 2, Name: "драма"},
			{ID: 4, Name: "боевик"},
		}`,
		Staff: []PersonJSON{
			{ID: 26, FullName: "Кристиан Бэйл", EnFullName: "Christian Charles Philip Bale", Photo: "https://avatars.mds.yandex.net/get-entity_search/141104/728910606/S600xU_2x", About: "Информация по этому человеку не указана"},
			{ID: 10, FullName: "Мэтт Дэймон", EnFullName: "Matt Damon", Photo: "https://avatars.mds.yandex.net/get-entity_search/1245892/935872902/S600xU_2x", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	21: {
		ID:              21,
		Name:            "Пророк. История Александра Пушкина",
		About:           `Пушкин молод, дерзок и любим публикой: он — звезда любого бала, поклонники носят его на руках, а девушки мечтают о его внимании. Но ни высокие покровители, ни верные друзья, ни известность не могут отвести от него злой рок — ссылки, дуэли, безденежье. Пока однажды он не встречает ту самую, единственную, кто осветит его жизнь и станет его судьбой.`,
		Poster:          "/static/img/21.webp",
		ReleaseYear:     "2024",
		Country:         "Россия",
		Budget:          9700000,
		BoxOfficeRussia: 19134822,
		// BoxOfficeGlobal: 19134822,
		Rating:   7.1,
		Duration: "1ч 53м",
		Genres: `[]GenreJSON{
			{ID: 2, Name: "драма"},
			{ID: 16, Name: "мюзикл"},
		}`,
		Staff: []PersonJSON{
			{ID: 40, FullName: "Юра Борисов", EnFullName: "Ura Borisov", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1898899/dc5079bf-795c-4953-a8c3-e3a5a025a1ff/600x900"},
		},
		Reviews: []ReviewJSON{},
	},

	22: {
		ID:              22,
		Name:            "Финист. первый Богатырь",
		About:           `В начале 1960-х Генри Форд II принимает решение улучшить имидж компании и сменить курс на производство более модных автомобилей. После неудавшейся попытки купить практически банкрота Ferrari американцы решают бросить вызов итальянским конкурентам на трассе и выиграть престижную гонку 24 часа Ле-Мана. Чтобы создать подходящую машину, компания нанимает автоконструктора Кэррола Шэлби, а тот отказывается работать без выдающегося, но, как считается, трудного в общении гонщика Кена Майлза. Вместе они принимаются за разработку впоследствии знаменитого спорткара Ford GT40.`,
		Poster:          "/static/img/22.webp",
		ReleaseYear:     "2024",
		Country:         "Россия",
		Budget:          15000000,
		BoxOfficeRussia: 31892047,
		// BoxOfficeGlobal: 31892047,
		Rating:   7.4,
		Duration: "2ч 32м",
		Genres: `[]GenreJSON{
			{ID: 10, Name: "приключения"},
			{ID: 11, Name: "фэнтези"},
		}`,
		Staff: []PersonJSON{
			{ID: 41, FullName: "Кирилл Зайцев", EnFullName: "Kirill Zaitzev", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1777765/5833d6a1-c63f-4d86-be47-b476746f4d1c/600x900", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},

	23: {
		ID:          23,
		Name:        "Батя",
		About:       `В начале 1960-х Генри Форд II принимает решение улучшить имидж компании и сменить курс на производство более модных автомобилей. После неудавшейся попытки купить практически банкрота Ferrari американцы решают бросить вызов итальянским конкурентам на трассе и выиграть престижную гонку 24 часа Ле-Мана. Чтобы создать подходящую машину, компания нанимает автоконструктора Кэррола Шэлби, а тот отказывается работать без выдающегося, но, как считается, трудного в общении гонщика Кена Майлза. Вместе они принимаются за разработку впоследствии знаменитого спорткара Ford GT40.`,
		Poster:      "/static/img/23.webp",
		ReleaseYear: "2020",
		Country:     "Россия",
		// Budget:          0,
		BoxOfficeRussia: 9404757,
		// BoxOfficeGlobal: 9 404 757,
		Rating:   7.8,
		Duration: "1ч 16м",
		Genres: `[]GenreJSON{
			{ID: 5, Name: "комедия"},
			{ID: 6, Name: "мелодрама"},
		}`,
		Staff: []PersonJSON{
			{ID: 42, FullName: "Владимир Вдовиченков", EnFullName: "Vladimir Vdovichenkov", Photo: "https://avatars.mds.yandex.net/get-kinopoisk-image/1777765/34a084da-0eb2-4fbd-af20-99cedf91de2a/600x900", About: "Информация по этому человеку не указана"},
		},
		Reviews: []ReviewJSON{},
	},
}

var ExistingGenres = map[string]GenreJSON{
	"thriller":     {ID: 1, Name: "триллер"},
	"drama":        {ID: 2, Name: "драма"},
	"fantastic":    {ID: 3, Name: "фантастика"},
	"action movie": {ID: 4, Name: "боевик"},
	"comedy":       {ID: 5, Name: "комедия"},
	"melodram":     {ID: 6, Name: "мелодрама"},
	"history":      {ID: 7, Name: "история"},
	"military":     {ID: 8, Name: "военный"},
	"criminal":     {ID: 9, Name: "криминал"},
	"advanture":    {ID: 10, Name: "приключения"},
	"fantasy":      {ID: 11, Name: "фэнтези"},
	"music":        {ID: 12, Name: "музыка"},
	"biography":    {ID: 13, Name: "биография"},
	"sport":        {ID: 14, Name: "спорт"},
	"detective":    {ID: 15, Name: "детектив"},
	"musicle":      {ID: 16, Name: "мюзикл"},
}
