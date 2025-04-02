# ER Diagram

```mermaid
erDiagram
    user {
        int id PK
        text login
        text hashed_password
        text avatar
        date birth_date
        timestamptz created_at
        timestamptz updated_at
    }

    collection {
        int id PK
        text name
        text slug
        timestamptz created_at
        timestamptz updated_at
    }

    person {
        int id PK
        text full_name
        text en_full_name
        text photo
        text about
        text sex
        text growth
        date birthday
        date death
        interval age
        text birth_place
        text death_place
        timestamptz created_at
        timestamptz updated_at
    }

    genre {
        int id PK
        text name
        timestamptz created_at
        timestamptz updated_at
    }

    country {
        int id PK
        text name
        text code
        text flag
    }

    movie {
        int id PK
        text name
        text about
        text poster
        text card
        date release_year
        int country FK
        text slogan
        text director
        int budget
        text box_office_us
        text box_office_global
        text boxo_office_russia
        date premiere_russia
        date premiere_global
        numeric rating
        text duration
        timestamptz created_at
        timestamptz updated_at
    }

    collection_movie {
        int collection_id FK PK
        int movie_id FK PK
    }

    movie_staff {
        int staff_id FK
        int movie_id FK
        text role
        PK staff_id, movie_id
    }

    movie_genre {
        int genre_id FK
        int movie_id FK
        PK genre_id, movie_id
    }

    review {
        int id PK
        int user_id FK
        int movie_id FK
        text review_text
        numeric score
        timestamptz created_at
        timestamptz updated_at
    }

    like {
        int id PK
        int user_id FK
        int review_id FK
        boolean is_valid
        timestamptz created_at
        timestamptz updated_at
    }

    dislike {
        int id PK
        int user_id FK
        int review_id FK
        boolean is_valid
        timestamptz created_at
        timestamptz updated_at
    }

    user ||--o{ review : "writes"
    movie ||--o{ review : "reviews"
    review ||--o{ like : "has"
    review ||--o{ dislike : "has"
    user ||--o{ like : "creates"
    user ||--o{ dislike : "creates"
    movie ||--o{ collection_movie : "belongs to"
    collection ||--o{ collection_movie : "contains"
    movie ||--o{ movie_staff : "has"
    person ||--o{ movie_staff : "works in"
    movie ||--o{ movie_genre : "has"
    genre ||--o{ movie_genre : "is"
    country ||--o{ movie : "produced in"
```
