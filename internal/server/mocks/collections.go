package mocks

type Film struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

type Collection map[int]Film

var BestOfAllTime = Collection{
	0: {ID: 1, Title: "Бойцовский клуб", Path: "/img/0.png"},
	1: {ID: 2, Title: "Тёмный рыцарь", Path: "/img/1.png"},
	2: {ID: 3, Title: "Форрест Гамп", Path: "/img/2.png"},
	3: {ID: 4, Title: "Крестный отец", Path: "/img/3.png"},
	4: {ID: 5, Title: "Интерстеллар", Path: "/img/4.png"},
	5: {ID: 6, Title: "Криминальное чтиво ", Path: "/img/5.png"},
	6: {ID: 7, Title: "Побег из Шоушенка", Path: "/img/6.png"},
	7: {ID: 8, Title: "Матрица", Path: "/img/7.png"},
	8: {ID: 9, Title: "Зелёная миля", Path: "/img/8.png"},
	9: {ID: 10, Title: "Одержимость", Path: "/img/9.png"},
}

var OskarNominees = Collection{
	0: {ID: 10, Title: "Ford против Ferrari", Path: "/img/10.png"},
	1: {ID: 11, Title: "Оппенгеймер", Path: "/img/11.png"},
	2: {ID: 12, Title: "Звёздные войны: Эпизод 4 – Новая надежда", Path: "/img/12.png"},
	3: {ID: 13, Title: "Рокки", Path: "/img/13.png"},
	4: {ID: 14, Title: "Джокер", Path: "/img/14.png"},
	5: {ID: 15, Title: "Игра в имитацию ", Path: "/img/15.png"},
	6: {ID: 16, Title: "Начало", Path: "/img/16.png"},
	7: {ID: 17, Title: "Назад в будущее", Path: "/img/17.png"},
	8: {ID: 18, Title: "Гладиатор", Path: "/img/18.png"},
	9: {ID: 19, Title: "Титаник", Path: "/img/19.png"},
}

var Collections = map[string]Collection{
	"BestOfAllTime": BestOfAllTime,
	"OskarNominees": OskarNominees,
}
