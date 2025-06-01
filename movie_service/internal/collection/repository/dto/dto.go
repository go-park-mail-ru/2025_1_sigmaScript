package dto

type MoviePostgresql struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	OriginalName    string  `json:"original_name,omitempty"`
	About           string  `json:"about,omitempty"`
	Poster          string  `json:"poster,omitempty"`
	ReleaseYear     int     `json:"release_year,omitempty"`
	Country         string  `json:"country,omitempty"`
	Slogan          string  `json:"slogan,omitempty"`
	Director        string  `json:"director,omitempty"`
	Budget          int64   `json:"budget,omitempty"`
	BoxOfficeUS     int64   `json:"box_office_us,omitempty"`
	BoxOfficeGlobal int64   `json:"box_office_global,omitempty"`
	BoxOfficeRussia int64   `json:"box_office_russia,omitempty"`
	PremierRussia   string  `json:"premier_russia,omitempty"`
	PremierGlobal   string  `json:"premier_global,omitempty"`
	Rating          float64 `json:"rating,omitempty"`
	Duration        string  `json:"duration,omitempty"`
	Created_at      string  `json:"created_at,omitempty"`
	Updated_at      string  `json:"updated_at,omitempty"`
}
