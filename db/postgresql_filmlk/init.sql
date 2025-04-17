DROP TABLE IF EXISTS "review" CASCADE;
DROP TABLE IF EXISTS "movie_genre" CASCADE;
DROP TABLE IF EXISTS "movie_staff" CASCADE;
DROP TABLE IF EXISTS "collection_movie" CASCADE;
DROP TABLE IF EXISTS "country" CASCADE;
DROP TABLE IF EXISTS "genre" CASCADE;
DROP TABLE IF EXISTS "person" CASCADE;
DROP TABLE IF EXISTS "collection" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "movie" CASCADE;

CREATE TABLE "user" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login TEXT UNIQUE NOT NULL CONSTRAINT loginchk CHECK (char_length(login) <= 255),
    hashed_password TEXT NOT NULL,
    avatar TEXT DEFAULT '/static/avatars/avatar_default_picture.svg',
    birth_date DATE DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "collection" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CONSTRAINT collection_namechk CHECK (char_length(name) <= 255),
    slug TEXT NOT NULL CONSTRAINT collection_slugchk CHECK (char_length(slug) <= 255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "person" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name TEXT NOT NULL,
    en_full_name TEXT DEFAULT NULL,
    photo TEXT DEFAULT '/static/avatars/avatar_default_picture.svg',
    about TEXT DEFAULT 'Информация по этому человеку не указана',
    sex TEXT DEFAULT 'secret',
    growth TEXT DEFAULT NULL,
    birthday DATE DEFAULT NULL,
    death DATE DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "genre" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "country" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    code TEXT NOT NULL,
    flag TEXT DEFAULT '/static/flags/flag_default_picture.webp',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "movie" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CONSTRAINT movie_namechk CHECK (char_length(name) <= 255),
    about TEXT DEFAULT 'Информация по этому фильму не указана',
    poster TEXT DEFAULT '/static/movies/poster_default_picture.webp',
    card TEXT DEFAULT '/static/movies/card_default_picture.webp',
    release_year DATE NOT NULL,
    country INTEGER REFERENCES country(id),
    slogan TEXT DEFAULT NULL,
    director TEXT DEFAULT NULL,
    budget DECIMAL DEFAULT 0,
    box_office_us DECIMAL DEFAULT 0,
    box_office_global DECIMAL DEFAULT 0,
    box_office_russia DECIMAL DEFAULT 0,
    premiere_russia DATE DEFAULT NULL,
    premiere_global DATE DEFAULT NULL,
    rating NUMERIC(4,2) CHECK (rating <= 10.00) DEFAULT 5.00,
    duration TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "collection_movie" (
    collection_id INTEGER REFERENCES collection(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, movie_id)
);

CREATE TABLE "movie_staff" (
    staff_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    role TEXT DEFAULT 'actor',
    PRIMARY KEY (staff_id, movie_id)
);

CREATE TABLE "movie_genre" (
    genre_id INTEGER REFERENCES genre(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (genre_id, movie_id)
);

CREATE TABLE "review" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    review_text TEXT NOT NULL,
    score NUMERIC(4,2) CHECK (score <= 10.00) DEFAULT 5.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
