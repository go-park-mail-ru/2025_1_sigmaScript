package models

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	PreviewURL  string  `json:"preview_url"`
	Duration    string  `json:"duration,omitempty"`
	ReleaseDate string  `json:"release_date,omitempty"` // add
	Rating      float64 `json:"rating,omitempty"`
}

type Collection map[int]Movie

type PersonJSON struct {
	ID              int        `json:"id"`
	FullName        string     `json:"full_name"`
	EnFullName      string     `json:"en_full_name,omitempty"`
	Photo           string     `json:"photo"`
	About           string     `json:"about,omitempty"`
	Sex             string     `json:"sex,omitempty"`
	Growth          string     `json:"growth,omitempty"`
	Birthday        string     `json:"birthday,omitempty"`
	Death           string     `json:"death,omitempty"`
	Career          string     `json:"career,omitempty"`
	Genres          string     `json:"genres,omitempty"`
	TotalFilms      string     `json:"total_films,omitempty"`
	MovieCollection Collection `json:"movie_collection,omitempty"`
}

type User struct {
	ID             string `json:"id,omitempty"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"`
	Avatar         string `json:"avatar"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type Profile struct {
	Username        string     `json:"username"`
	Avatar          string     `json:"avatar"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at,omitempty"`
	MovieCollection []Movie      `json:"movie_collection,omitempty"`
	Actors          []PersonJSON `json:"actors,omitempty"`
	Reviews         []ReviewJSON `json:"reviews,omitempty"`
}

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

type RegisterData struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
