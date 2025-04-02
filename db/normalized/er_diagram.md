# ER Diagram

```mermaid
erDiagram
    USER {
        id
        login
        hashed_password
        avatar
        birth_date
        created_at
        updated_at
    }

    COLLECTION {
        id
        name
        slug
        created_at
        updated_at
    }

    PERSON {
        id
        full_name
        en_full_name
        photo
        about
        sex
        growth
        birthday
        death
        age
        birth_place
        death_place
        created_at
        updated_at
    }
    
    GENRE {
        id
        name
        created_at
        updated_at
    }
    
    COUNTRY {
        id
        name
        code
        flag
    }
    
    MOVIE {
        id
        name
        about
        poster
        card
        release_year
        slogan
        director
        budget
        box_office_us
        box_office_global
        box_office_russia
        premiere_russia
        premiere_global
        rating
        duration
        created_at
        updated_at
    }
    
    COLLECTION_MOVIE {
        collection_id
        movie_id
    }
    
    MOVIE_GENRE {
        genre_id
        movie_id
    }
    
    MOVIE_STAFF {
        staff_id
        movie_id
        role
    }
    
    REVIEW {
        id
        user_id
        movie_id
        review_text
        score
        created_at
        updated_at
    }
    
    LIKE {
        id
        user_id
        review_id
        is_valid
        created_at
        updated_at
    }
    
    DISLIKE {
        id
        user_id
        review_id
        is_valid
        created_at
        updated_at
    }

    USER ||--o{ REVIEW : writes
    USER ||--o{ LIKE : gives
    USER ||--o{ DISLIKE : gives
    REVIEW ||--o{ LIKE : has
    REVIEW ||--o{ DISLIKE : has
    REVIEW }o--|| MOVIE : belongs_to
    COLLECTION ||--o{ COLLECTION_MOVIE : contains
    COLLECTION_MOVIE }o--|| MOVIE : includes
    MOVIE ||--o{ MOVIE_GENRE : categorized_as
    GENRE ||--o{ MOVIE_GENRE : belongs_to
    MOVIE ||--o{ MOVIE_STAFF : has
    PERSON ||--o{ MOVIE_STAFF : works_in
    MOVIE }o--|| COUNTRY : produced_in
```
