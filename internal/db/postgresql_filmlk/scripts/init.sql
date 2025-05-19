DROP TABLE IF EXISTS "review" CASCADE;
DROP TABLE IF EXISTS "user_person_favorite" CASCADE;
DROP TABLE IF EXISTS "user_movie_favorite" CASCADE;
DROP TABLE IF EXISTS "career_person" CASCADE;
DROP TABLE IF EXISTS "person_genre" CASCADE;
DROP TABLE IF EXISTS "movie_genre" CASCADE;
DROP TABLE IF EXISTS "movie_staff" CASCADE;
DROP TABLE IF EXISTS "collection_movie" CASCADE;
DROP TABLE IF EXISTS "genre" CASCADE;
DROP TABLE IF EXISTS "career" CASCADE;
DROP TABLE IF EXISTS "person" CASCADE;
DROP TABLE IF EXISTS "collection" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS "movie" CASCADE;


DROP TYPE IF EXISTS sex_type;
CREATE TYPE sex_type AS ENUM ('Мужчина', 'Женщина', 'secret', '');

CREATE TABLE "user" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login TEXT UNIQUE NOT NULL CONSTRAINT loginchk CHECK (char_length(login) <= 255),
    hashed_password TEXT NOT NULL,
    avatar TEXT DEFAULT '/static/avatars/avatar_default_picture.svg',
    birth_date DATE DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_id ON "user"(id);

CREATE TABLE "collection" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CONSTRAINT collection_namechk CHECK (char_length(name) <= 255),
    is_main_collection BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_collection_id ON collection(id);

DROP TYPE IF EXISTS career_type;
CREATE TYPE career_type AS ENUM ('Актёр', 'Продюсер', 'Режиссёр', 'Сценарист');

CREATE TABLE "career" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    career career_type DEFAULT 'Актёр' NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_career_id ON career(id);

CREATE TABLE "person" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name TEXT NOT NULL,
    en_full_name TEXT DEFAULT NULL,
    photo TEXT DEFAULT '/static/avatars/avatar_default_picture.svg',
    about TEXT DEFAULT 'Информация по этому человеку не указана',
    sex sex_type DEFAULT '',
    growth TEXT DEFAULT NULL,
    birthday DATE DEFAULT NULL,
    death DATE DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_person_id ON person(id);

CREATE TABLE "genre" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_genre_id ON genre(id);

CREATE TABLE "movie" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CONSTRAINT movie_namechk CHECK (char_length(name) <= 255),
    original_name TEXT CONSTRAINT movie_orig_namechk CHECK (char_length(original_name) <= 255) DEFAULT NULL,
    about TEXT DEFAULT 'Информация по этому фильму не указана',
    poster TEXT DEFAULT '/static/movies/poster_default_picture.webp',
    promo_poster TEXT DEFAULT '/static/movies/poster_default_picture.webp',
    release_year TIMESTAMPTZ DEFAULT NULL,
    slogan TEXT DEFAULT NULL,
    director TEXT DEFAULT NULL,
    country TEXT DEFAULT NULL,
    budget DECIMAL DEFAULT 0,
    box_office_us DECIMAL DEFAULT 0,
    box_office_global DECIMAL DEFAULT 0,
    box_office_russia DECIMAL DEFAULT 0,
    premier_russia DATE DEFAULT NULL,
    premier_global DATE DEFAULT NULL,
    rating DECIMAL(4,2) CHECK (rating <= 10.00) DEFAULT 5.00,
    duration TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    rating_kp DECIMAL(4,2) CHECK (rating <= 10.00) DEFAULT NULL,
    rating_imdb DECIMAL(4,2) CHECK (rating <= 10.00) DEFAULT NULL,
    short_description TEXT DEFAULT NULL,
    logo TEXT DEFAULT '#',
    backdrop TEXT DEFAULT '#',
    watchability BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_movie_id ON movie(id);

CREATE TABLE "collection_movie" (
    collection_id INTEGER REFERENCES collection(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, movie_id)
);
CREATE INDEX idx_collection_movie_movie_id ON collection_movie(movie_id);
CREATE INDEX idx_collection_movie_collection_id ON collection_movie(collection_id);

CREATE TABLE "movie_staff" (
    staff_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    role career_type DEFAULT 'Актёр' NOT NULL,
    PRIMARY KEY (staff_id, movie_id)
);
CREATE INDEX idx_movie_staff_movie_id ON movie_staff(movie_id);
CREATE INDEX idx_movie_staff_staff_id ON movie_staff(staff_id);

CREATE TABLE "movie_genre" (
    genre_id INTEGER REFERENCES genre(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (genre_id, movie_id)
);
CREATE INDEX idx_movie_genre_movie_id ON movie_genre(movie_id);
CREATE INDEX idx_movie_genre_genre_id ON movie_genre(genre_id);

CREATE TABLE "career_person" (
    career_id INTEGER REFERENCES career(id) ON DELETE CASCADE,
    person_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    PRIMARY KEY (career_id, person_id)
);
CREATE INDEX idx_career_person_career_id ON career_person(career_id);
CREATE INDEX idx_career_person_person_id ON career_person(person_id);

CREATE TABLE "person_genre" (
    person_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genre(id) ON DELETE CASCADE,
    PRIMARY KEY (genre_id, person_id)
);
CREATE INDEX idx_person_genre_genre_id ON person_genre(genre_id);
CREATE INDEX idx_person_genre_person_id ON person_genre(person_id);


CREATE TABLE "review" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    review_text TEXT DEFAULT NULL,
    score DECIMAL(4,2) CHECK (score <= 10.00) DEFAULT 5.00,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_user_movie_id UNIQUE (user_id, movie_id)
);
CREATE INDEX idx_review_id ON review(id);
CREATE INDEX idx_review_user_id ON review(user_id);
CREATE INDEX idx_review_movie_id ON review(movie_id);


CREATE TABLE "user_person_favorite" (
    person_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, person_id)
);
CREATE INDEX idx_user_person_favorite_user_id ON user_person_favorite(user_id);
CREATE INDEX idx_user_person_favorite_person_id ON user_person_favorite(person_id);

CREATE TABLE "user_movie_favorite" (
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, movie_id)
);
CREATE INDEX idx_user_movie_favorite_user_id ON user_movie_favorite(user_id);
CREATE INDEX idx_user_movie_favorite_movie_id ON user_movie_favorite(movie_id);


-- add update triggers to tables
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_review_updated_at
BEFORE UPDATE ON "review"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_genre_updated_at
BEFORE UPDATE ON "genre"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_person_updated_at
BEFORE UPDATE ON "person"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_collection_updated_at
BEFORE UPDATE ON "collection"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_movie_updated_at
BEFORE UPDATE ON "movie"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_career_updated_at
BEFORE UPDATE ON "career"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
