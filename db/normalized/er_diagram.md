# ER Diagram

```mermaid
erDiagram
    user {
        INTEGER id PK
        TEXT login
        TEXT hashed_password
        TEXT avatar
        DATE birth_date
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    collection {
        INTEGER id PK
        TEXT name
        TEXT slug
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    person {
        INTEGER id PK
        TEXT full_name
        TEXT en_full_name
        TEXT photo
        TEXT about
        TEXT sex
        TEXT growth
        DATE birthday
        DATE death
        INTERVAL age
        TEXT birth_place
        TEXT death_place
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    genre {
        INTEGER id PK
        TEXT name
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    country {
        INTEGER id PK
        TEXT name
        TEXT code
        TEXT flag
    }
    movie {
        INTEGER id PK
        TEXT name
        TEXT about
        TEXT poster
        TEXT card
        DATE release_year
        INTEGER country_id FK
        TEXT slogan
        TEXT director
        INTEGER budget
        TEXT box_office_us
        TEXT box_office_global
        TEXT boxo_office_russia
        DATE premiere_russia
        DATE premiere_global
        NUMERIC rating
        TEXT duration
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    collection_movie {
        INTEGER collection_id PK, FK
        INTEGER movie_id PK, FK
    }
    movie_staff {
        INTEGER staff_id PK, FK
        INTEGER movie_id PK, FK
        TEXT role
    }
    movie_genre {
        INTEGER genre_id PK, FK
        INTEGER movie_id PK, FK
    }
    review {
        INTEGER id PK
        INTEGER user_id FK
        INTEGER movie_id FK
        TEXT review_text
        NUMERIC score
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    like {
        INTEGER id PK
        INTEGER user_id FK
        INTEGER review_id FK
        BOOLEAN is_valid
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
    dislike {
        INTEGER id PK
        INTEGER user_id FK
        INTEGER review_id FK
        BOOLEAN is_valid
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    user                ||--o{ review           : "writes"
    movie               ||--o{ review           : "has"
    user                ||--o{ like             : "gives"
    review              ||--o{ like             : "receives"
    user                ||--o{ dislike          : "gives"
    review              ||--o{ dislike          : "receives"
    country             ||--o{ movie            : "originates_from"

    collection          ||--o{ collection_movie : "groups"
    movie               ||--o{ collection_movie : "is_grouped_in"

    person              ||--o{ movie_staff      : "is_staff_for"
    movie               ||--o{ movie_staff      : "has_staff"

    genre               ||--o{ movie_genre      : "categorizes"
    movie               ||--o{ movie_genre      : "has_genre"

```
