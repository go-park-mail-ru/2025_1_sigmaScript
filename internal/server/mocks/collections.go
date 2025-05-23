package mocks

type Film struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	PreviewURL string `json:"preview_url"`
}

type Collection map[int]Film

type Collections map[string]Collection

var BestOfAllTime = Collection{
	0: {ID: 0, Title: "Бойцовский клуб", PreviewURL: "/img/0.webp"},
	1: {ID: 1, Title: "Тёмный рыцарь", PreviewURL: "/img/1.webp"},
	2: {ID: 2, Title: "Форрест Гамп", PreviewURL: "/img/2.webp"},
	3: {ID: 3, Title: "Крестный отец", PreviewURL: "/img/3.webp"},
	4: {ID: 4, Title: "Интерстеллар", PreviewURL: "/img/4.webp"},
	5: {ID: 5, Title: "Криминальное чтиво ", PreviewURL: "/img/5.webp"},
	6: {ID: 6, Title: "Побег из Шоушенка", PreviewURL: "/img/6.webp"},
	7: {ID: 7, Title: "Матрица", PreviewURL: "/img/7.webp"},
	8: {ID: 8, Title: "Зелёная миля", PreviewURL: "/img/8.webp"},
	9: {ID: 9, Title: "Одержимость", PreviewURL: "/img/9.webp"},
}

var OskarNominees = Collection{
	0: {ID: 10, Title: "Ford против Ferrari", PreviewURL: "/img/10.webp"},
	1: {ID: 11, Title: "Оппенгеймер", PreviewURL: "/img/11.webp"},
	2: {ID: 12, Title: "Звёздные войны: Эпизод 4 – Новая надежда", PreviewURL: "/img/12.webp"},
	3: {ID: 13, Title: "Рокки", PreviewURL: "/img/13.webp"},
	4: {ID: 14, Title: "Джокер", PreviewURL: "/img/14.webp"},
	5: {ID: 15, Title: "Игра в имитацию ", PreviewURL: "/img/15.webp"},
	6: {ID: 16, Title: "Начало", PreviewURL: "/img/16.webp"},
	7: {ID: 17, Title: "Назад в будущее", PreviewURL: "/img/17.webp"},
	8: {ID: 18, Title: "Гладиатор", PreviewURL: "/img/18.webp"},
	9: {ID: 19, Title: "Титаник", PreviewURL: "/img/19.webp"},
}

var MainPageCollections = Collections{
	"Лучшие за всё время": BestOfAllTime,
	"Номинанты на оскар":  OskarNominees,
}
