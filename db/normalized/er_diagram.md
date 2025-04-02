# ER Diagram

```mermaid
erDiagram
    user {
        id
        login
        hashed_password
        avatar
        birth_date
        created_at
        updated_at
    }
    collection {
        id
        name
        slug
        created_at
        updated_at
    }
    person {
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
    genre {
        id
        name
        created_at
        updated_at
    }
    country {
        id
        name
        code
        flag
    }
    movie {
        id
        name
        about
        poster
        card
        release_year
        country_id  // FK to country
        slogan
        director
        budget
        box_office_us
        box_office_global
        boxo_office_russia
        premiere_russia
        premiere_global
        rating
        duration
        created_at
        updated_at
    }
    collection_movie {
        collection_id // FK to collection
        movie_id    // FK to movie
        // Composite PK (collection_id, movie_id)
    }
    movie_staff {
        staff_id   // FK to person
        movie_id   // FK to movie
        role
        // Composite PK (staff_id, movie_id)
    }
    movie_genre {
        genre_id   // FK to genre
        movie_id   // FK to movie
        // Composite PK (genre_id, movie_id)
    }
    review {
        id
        user_id    // FK to user
        movie_id   // FK to movie
        review_text
        score
        created_at
        updated_at
    }
    like {
        id
        user_id    // FK to user
        review_id  // FK to review
        is_valid
        created_at
        updated_at
    }
    dislike {
        id
        user_id    // FK to user
        review_id  // FK to review
        is_valid
        created_at
        updated_at
    }

    user                ||--o{ review           : "пишет"
    movie               ||--o{ review           : "имеет"
    user                ||--o{ like             : "ставит"
    review              ||--o{ like             : "получает"
    user                ||--o{ dislike          : "ставит"
    review              ||--o{ dislike          : "получает"
    country             ||--o{ movie            : "страна_производства" // Связь Один-ко-Многим (одна страна -> много фильмов)

    collection          ||--o{ collection_movie : "содержит" // Связь Один-ко-Многим (одна коллекция -> много записей в collection_movie)
    movie               ||--o{ collection_movie : "содержится_в" // Связь Один-ко-Многим (один фильм -> много записей в collection_movie)

    person              ||--o{ movie_staff      : "участвует_в" // Связь Один-ко-Многим (один человек -> много ролей в фильмах)
    movie               ||--o{ movie_staff      : "имеет_состав" // Связь Один-ко-Многим (один фильм -> много участников)

    genre               ||--o{ movie_genre      : "определяет" // Связь Один-ко-Многим (один жанр -> много записей в movie_genre)
    movie               ||--o{ movie_genre      : "принадлежит_к" // Связь Один-ко-Многим (один фильм -> много записей в movie_genre)

```
