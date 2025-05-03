package mocks

type Movie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	PreviewURL string `json:"preview_url"`
	Duration   string `json:"duration,omitempty"`
}

type Collection map[int]Movie

type Collections map[string]Collection

var BestOfAllTime = Collection{
	0: {ID: 1, Title: "Бойцовский клуб", PreviewURL: "/img/0.webp"},
	1: {ID: 2, Title: "Матрица", PreviewURL: "/img/7.webp"},
	2: {ID: 3, Title: "Форрест Гамп", PreviewURL: "/img/2.webp"},
	3: {ID: 4, Title: "Крестный отец", PreviewURL: "/img/3.webp"},
	4: {ID: 5, Title: "Интерстеллар", PreviewURL: "/img/4.webp"},
	5: {ID: 6, Title: "Криминальное чтиво ", PreviewURL: "/img/5.webp"},
	6: {ID: 7, Title: "Побег из Шоушенка", PreviewURL: "/img/6.webp"},
	7: {ID: 8, Title: "Тёмный рыцарь", PreviewURL: "/img/1.webp"},
	8: {ID: 9, Title: "Зелёная миля", PreviewURL: "/img/8.webp"},
	9: {ID: 10, Title: "Одержимость", PreviewURL: "/img/9.webp"},
}

var OskarNominees = Collection{
	0: {ID: 11, Title: "Оппенгеймер", PreviewURL: "/img/11.webp"},
	1: {ID: 12, Title: "Звёздные войны: Эпизод 4 – Новая надежда", PreviewURL: "/img/12.webp"},
	2: {ID: 13, Title: "Рокки", PreviewURL: "/img/13.webp"},
	3: {ID: 14, Title: "Джокер", PreviewURL: "/img/14.webp"},
	4: {ID: 15, Title: "Игра в имитацию", PreviewURL: "/img/15.webp"},
	5: {ID: 16, Title: "Начало", PreviewURL: "/img/16.webp"},
	6: {ID: 17, Title: "Назад в будущее", PreviewURL: "/img/17.webp"},
	7: {ID: 18, Title: "Гладиатор", PreviewURL: "/img/18.webp"},
	8: {ID: 19, Title: "Титаник", PreviewURL: "/img/19.webp"},
	9: {ID: 20, Title: "Ford против Ferrari", PreviewURL: "/img/10.webp"},
}

var MainPageCollections = Collections{
	"Лучшие за всё время": BestOfAllTime,
	"Номинанты на оскар":  OskarNominees,
	"promo":               Promo,
}

var Promo = Collection{
	0: {ID: 21, Title: "Пророк. История Александра Пушкина", PreviewURL: "/img/promo_prorok.webp", Duration: "1ч 53м"},
	1: {ID: 22, Title: "Батя", PreviewURL: "/img/promo_batya.webp", Duration: "1ч 16м"},
	2: {ID: 23, Title: "Финист. первый Богатырь", PreviewURL: "/img/promo_bogatyr.webp", Duration: "1ч 52м"},
}
