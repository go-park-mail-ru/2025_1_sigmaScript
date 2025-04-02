# ER Diagram

```mermaid
erDiagram
    user {
        int id PK
        text login
        text hashed_password
        text avatar
        date birth_date
        timestamp created_at
        timestamp updated_at
    }
    collection {
        int id PK
        text name
        text slug
        timestamp created_at
        timestamp updated_at
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
        timestamp created_at
        timestamp updated_at
    }
    genre {
        int id PK
        text name
        timestamp created_at
        timestamp updated_at
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
        timestamp created_at
        timestamp updated_at
    }
    collection_movie {
        int collection_id FK
        int movie_id FK
        PRIMARY KEY (collection_id, movie_id)
    }
    movie_staff {
        int staff_id FK
        int movie_id FK
        text role
        PRIMARY KEY (staff_id, movie_id)
    }
    movie_genre {
        int genre_id FK
        int movie_id FK
        PRIMARY KEY (genre_id, movie_id)
    }
    review {
        int id PK
        int user_id FK
        int movie_id FK
        text review_text
        numeric score
        timestamp created_at
        timestamp updated_at
    }
    like {
        int id PK
        int user_id FK
        int review_id FK
        boolean is_valid
        timestamp created_at
        timestamp updated_at
    }
    dislike {
        int id PK
        int user_id FK
        int review_id FK
        boolean is_valid
        timestamp created_at
        timestamp updated_at
    }

    user ||--o{ review : "writes"
    user ||--o{ like : "gives"
    user ||--o{ dislike : "gives"
    collection ||--o{ collection_movie : "contains"
    movie ||--o{ collection_movie : "is in"
    movie ||--o{ movie_staff : "has"
    person ||--o{ movie_staff : "works on"
    movie ||--o{ movie_genre : "has"
    genre ||--o{ movie_genre : "is"
    movie ||--o{ review : "is reviewed"
    movie ||--o{ country : "is from"
    review ||--o{ like : "has"
    review ||--o{ dislike : "has"
```
