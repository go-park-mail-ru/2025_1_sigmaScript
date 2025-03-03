package mocks

type Film struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

type Collection map[int]Film

var BestOfAllTime = Collection{
	0: {ID: 1, Title: "Бойцовский клуб", Path: "/img/0.webp"},
	1: {ID: 2, Title: "Тёмный рыцарь", Path: "/img/1.webp"},
	2: {ID: 3, Title: "Форрест Гамп", Path: "/img/2.webp"},
	3: {ID: 4, Title: "Крестный отец", Path: "/img/3.webp"},
	4: {ID: 5, Title: "Интерстеллар", Path: "/img/4.webp"},
	5: {ID: 6, Title: "Криминальное чтиво ", Path: "/img/5.webp"},
	6: {ID: 7, Title: "Побег из Шоушенка", Path: "/img/6.webp"},
	7: {ID: 8, Title: "Матрица", Path: "/img/7.webp"},
	8: {ID: 9, Title: "Зелёная миля", Path: "/img/8.webp"},
	9: {ID: 10, Title: "Одержимость", Path: "/img/9.webp"},
}

var OskarNominees = Collection{
	0: {ID: 10, Title: "Ford против Ferrari", Path: "/img/10.webp"},
	1: {ID: 11, Title: "Оппенгеймер", Path: "/img/11.webp"},
	2: {ID: 12, Title: "Звёздные войны: Эпизод 4 – Новая надежда", Path: "/img/12.webp"},
	3: {ID: 13, Title: "Рокки", Path: "/img/13.webp"},
	4: {ID: 14, Title: "Джокер", Path: "/img/14.webp"},
	5: {ID: 15, Title: "Игра в имитацию ", Path: "/img/15.webp"},
	6: {ID: 16, Title: "Начало", Path: "/img/16.webp"},
	7: {ID: 17, Title: "Назад в будущее", Path: "/img/17.webp"},
	8: {ID: 18, Title: "Гладиатор", Path: "/img/18.webp"},
	9: {ID: 19, Title: "Титаник", Path: "/img/19.webp"},
}

var Collections = map[string]Collection{
	"OskarNominees": OskarNominees,
	"BestOfAllTime": BestOfAllTime,
}
