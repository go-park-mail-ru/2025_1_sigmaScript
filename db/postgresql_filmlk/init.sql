DROP TABLE IF EXISTS "review" CASCADE;
DROP TABLE IF EXISTS "watch_provider" CASCADE;
DROP TABLE IF EXISTS "user_person_favorite" CASCADE;
DROP TABLE IF EXISTS "user_movie_favorite" CASCADE;
DROP TABLE IF EXISTS "career_person" CASCADE;
DROP TABLE IF EXISTS "movie_country" CASCADE;
DROP TABLE IF EXISTS "person_genre" CASCADE;
DROP TABLE IF EXISTS "movie_genre" CASCADE;
DROP TABLE IF EXISTS "movie_staff" CASCADE;
DROP TABLE IF EXISTS "collection_movie" CASCADE;
DROP TABLE IF EXISTS "country" CASCADE;
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

CREATE TABLE "collection" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CONSTRAINT collection_namechk CHECK (char_length(name) <= 255),
    is_main_collection BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


DROP TYPE IF EXISTS career_type;
CREATE TYPE career_type AS ENUM ('Актёр', 'Продюсер', 'Режиссёр', 'Сценарист');

CREATE TABLE "career" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    career career_type DEFAULT 'Актёр' NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE "genre" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "country" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    flag TEXT DEFAULT '/static/flags/flag_default_picture.webp',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE "collection_movie" (
    collection_id INTEGER REFERENCES collection(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (collection_id, movie_id)
);

CREATE TABLE "movie_staff" (
    staff_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    role career_type DEFAULT 'Актёр' NOT NULL,
    PRIMARY KEY (staff_id, movie_id)
);

CREATE TABLE "movie_genre" (
    genre_id INTEGER REFERENCES genre(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (genre_id, movie_id)
);

CREATE TABLE "career_person" (
    career_id INTEGER REFERENCES career(id) ON DELETE CASCADE,
    person_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    PRIMARY KEY (career_id, person_id)
);

CREATE TABLE "person_genre" (
    person_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genre(id) ON DELETE CASCADE,
    PRIMARY KEY (genre_id, person_id)
);

CREATE TABLE "movie_country" (
    country_id INTEGER REFERENCES country(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    PRIMARY KEY (country_id, movie_id)
);

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

CREATE TABLE "user_person_favorite" (
    person_id INTEGER REFERENCES person(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, person_id)
);

CREATE TABLE "user_movie_favorite" (
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, movie_id)
);

CREATE TABLE "watch_provider" (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    movie_id INTEGER REFERENCES movie(id) ON DELETE CASCADE,
	name TEXT NOT NULL CONSTRAINT provider_namechk CHECK (char_length(name) <= 255),
	logo TEXT DEFAULT '/static/svg/play.svg',
	watch_url TEXT DEFAULT '#'
);


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

CREATE TRIGGER update_country_updated_at
BEFORE UPDATE ON "country"
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



-- some default data

-- creating collections
INSERT INTO collection (name, is_main_collection) VALUES
('Лучшие за всё время', TRUE),
('Номинанты на оскар', TRUE),
('promo', TRUE);


-- creating career roles
INSERT INTO career (career) VALUES
('Актёр'),
('Продюсер'),
('Режиссёр'),
('Сценарист');

-- creating genres
INSERT INTO genre (name) values
('триллер'), -- 1
('драма'), -- 2
('фантастика'), -- 3
('боевик'), -- 4
('комедия'), -- 5
('мелодрама'), -- 6
('история'), -- 7
('военный'), -- 8
('криминал'), -- 9
('приключения'), -- 10
('фэнтези'), -- 11
('музыка'), -- 12
('биография'), -- 13
('спорт'), -- 14
('детектив'), -- 15
('мюзикл'), -- 16
('мультфильм'), -- 17
('семейный'); -- 18

-- creating countries
INSERT INTO country (name) values
('Канада'),
('США');




INSERT INTO person (full_name, en_full_name, photo) VALUES
('Леонардо Ди Каприо', 'Leonardo DiCaprio', 'https://avatars.mds.yandex.net/get-entity_search/2310675/1130394491/S600xU_2x'),
('Морган Фримен', 'Morgan Freeman', 'https://avatars.mds.yandex.net/get-entity_search/2057552/1132084397/S600xU_2x'),
('Джонни Депп', 'Johnny Depp', '/static/avatars/avatar_default_picture.svg'),
('Том Круз', 'Tom Cruise', '/static/avatars/avatar_default_picture.svg'),
('Уилл Смит', 'Will Smith', '/static/avatars/avatar_default_picture.svg'),
('Мэтт Дэймон', 'Matt Damon', 'https://avatars.mds.yandex.net/get-entity_search/1245892/935872902/S600xU_2x'),
('Киану Ривз', '', ''),
('Боб Персичетти', '', ''),
('Питер Рэмзи', '', ''),
('Родни Ротман', '', '');

-- peron id 15-17
INSERT INTO person (full_name, en_full_name, photo) VALUES
('Шамеик Мур', 'Shameik Moore', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_2378404.jpg'),
('Джейк Джонсон', 'Jake Johnson', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_1089330.jpg'),
('Хейли Стайнфелд', 'Hailee Steinfeld', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_1478559.jpg');



UPDATE person SET (full_name, en_full_name, photo, about, sex, growth, birthday) =
(
    'Киану Ривз',
    'Keanu Reeves',
    'https://i.pinimg.com/originals/a3/70/0b/a3700bdf15fcceabf740e1f347dbb5a2.jpg',
    '\nКиану Чарльз Ривз — канадский актёр, кинорежиссёр, кинопродюсер и музыкант.\nНаиболее известен своими ролями в киносериях «Матрица», «Билл и Тед», «Джон Уик», а также в фильмах «На гребне волны», «Скорость», «Адвокат дьявола», «Константин: Повелитель тьмы».\nОбладатель звезды на Голливудской «Аллее славы».',
    'Мужчина',
    '186',
    '1964-09-2'
) where full_name = 'Киану Ривз';


-- Start Transaction
BEGIN;

-- Optional: Ensure career types exist (usually handled by ENUM/initial setup)
-- INSERT INTO career (career) VALUES ('Актёр') ON CONFLICT (career) DO NOTHING;
-- INSERT INTO career (career) VALUES ('Продюсер') ON CONFLICT (career) DO NOTHING;
-- INSERT INTO career (career) VALUES ('Режиссёр') ON CONFLICT (career) DO NOTHING;
-- INSERT INTO career (career) VALUES ('Сценарист') ON CONFLICT (career) DO NOTHING;

-- ==========================================================================
-- Movie: 361 - Бойцовский клуб (Fight Club)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 3]
INSERT INTO genre (name) VALUES ('триллер') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('драма') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('криминал') ON CONFLICT (name) DO NOTHING;

-- Insert Persons (Staff) if they don't exist [cite: 3, 4]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(25774, 'Эдвард Нортон', 'Edward Norton', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_25774.jpg'),
(25584, 'Брэд Питт', 'Brad Pitt', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_25584.jpg'),
(25775, 'Хелена Бонем Картер', 'Helena Bonham Carter', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_25775.jpg'),
(14127, 'Мит Лоаф', 'Meat Loaf', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_14127.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 1, 2, 3, 4]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    361,                                                                           -- ID [cite: 1]
    'Бойцовский клуб',                                                           -- Name [cite: 1]
    'Fight Club',                                                                  -- OriginalName [cite: 1]
    E'Сотрудник страховой компании страдает хронической бессонницей и отчаянно пытается вырваться из мучительно скучной жизни. Однажды в очередной командировке он встречает некоего Тайлера Дёрдена — харизматического торговца мылом с извращенной философией. Тайлер уверен, что самосовершенствование — удел слабых, а единственное, ради чего стоит жить, — саморазрушение. Проходит немного времени, и \nвот уже новые друзья лупят друг друга почем зря на стоянке перед баром, и очищающий мордобой доставляет им высшее блаженство. Приобщая других мужчин к простым радостям физической жестокости, они основывают тайный Бойцовский клуб, который начинает пользоваться невероятной популярностью.', -- About [cite: 1, 2]
    'https://image.openmoviedb.com/kinopoisk-images/4716873/85b585ea-410f-4d1c-aaa5-8d242756c2a4/orig', -- Poster [cite: 2]
    NULL,                                                                          -- PromoURL is empty [cite: 2]
    '1999-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 2]
    '«Интриги. Хаос. Мыло»',                                                       -- Slogan [cite: 2, 3]
    'Дэвид Финчер',                                                                -- Director [cite: 3]
    'США',                                                                         -- Country [cite: 2]
    63000000,                                                                      -- Budget [cite: 3]
    37030102,                                                                      -- BoxOfficeUS [cite: 3]
    100853753,                                                                     -- BoxOfficeGlobal [cite: 3]
    334590,                                                                        -- BoxOfficeRussia [cite: 3]
    '2000-01-13'::DATE,                                                            -- PremierRussia [cite: 3]
    '1999-09-10'::DATE,                                                            -- PremierGlobal [cite: 3]
    7.5,                                                                           -- Rating [cite: 3]
    '139',                                                                         -- Duration [cite: 3]
    8.673,                                                                         -- RatingKP [cite: 4]
    8.8,                                                                           -- RatingIMDB [cite: 4]
    'Страховой работник разрушает рутину своей благополучной жизни. Культовая драма по книге Чака Паланика', -- ShortDescription [cite: 2]
    'https://image.openmoviedb.com/tmdb-images/original/y9RSpK5PpMYEkfdCRofBp09KpW9.png', -- Logo [cite: 4]
    'https://image.openmoviedb.com/tmdb-images/original/hZkgoQYus5vegHoetLkCJzb17zJ.jpg', -- Backdrop [cite: 4]
    FALSE                                                                          -- Watchability is nil [cite: 4]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 3, 4]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(25774, 361, 'Актёр'),
(25584, 361, 'Актёр'),
(25775, 361, 'Актёр'),
(14127, 361, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 3]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'триллер'), 361),
((SELECT id from genre where name = 'драма'), 361),
((SELECT id from genre where name = 'криминал'), 361)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 400787 - Матрица (Matrix) - Note: This seems to be data for a TV Series 'Matrix (1993)' not the movie.
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 9]
INSERT INTO genre (name) VALUES ('фэнтези') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('боевик') ON CONFLICT (name) DO NOTHING;
-- 'триллер', 'драма' already handled

-- Insert Persons (Staff) if they don't exist [cite: 9]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(18386, 'Ник Манкузо', 'Nick Mancuso', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_18386.jpg'),
(1768, 'Филлип Джарретт', 'Phillip Jarrett', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_1768.jpg'),
(6226, 'Кэрри-Энн Мосс', 'Carrie-Anne Moss', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_6226.jpg'),
(45338, 'Джон Вернон', 'John Vernon', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_45338.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 4, 5, 6, 7, 8, 9, 10]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    400787,                                                                        -- ID [cite: 4]
    'Матрица',                                                                     -- Name [cite: 4]
    'Matrix',                                                                      -- OriginalName [cite: 4]
    E'В центре сюжета — жестокий наемный убийца Стивен Матрица. Однажды с ним происходит то же, что и с его жертвами, — он попадает в загробный мир. В преисподней Стивену дают второй шанс — сделать выбор: отправиться в ад или искупить свою вину. Заработать отсрочку он может только помогая людям на Земле. Киллер соглашается и приходит в сознание в больнице, выйдя из состояния клинической смерти. Но став своеобразным ангелом-хранителем, он часто использует в работе свои старые профессиональные методы.', -- About [cite: 4, 5, 6, 7, 8]
    'https://image.openmoviedb.com/kinopoisk-images/1629390/55c4ed3b-4ddf-4e33-9607-de343320ee96/orig', -- Poster [cite: 8]
    NULL,                                                                          -- PromoURL is empty [cite: 8]
    '1993-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 8]
    NULL,                                                                          -- Slogan is empty [cite: 8]
    'Хорхе Монтеси',                                                               -- Director [cite: 8]
    'Канада',                                                                      -- Country [cite: 8]
    0,                                                                             -- Budget [cite: 8]
    0,                                                                             -- BoxOfficeUS [cite: 8]
    0,                                                                             -- BoxOfficeGlobal [cite: 8]
    0,                                                                             -- BoxOfficeRussia [cite: 8]
    NULL,                                                                          -- PremierRussia is empty [cite: 8]
    '1993-03-01'::DATE,                                                            -- PremierGlobal [cite: 8]
    7.5,                                                                           -- Rating [cite: 8, 9]
    '43',                                                                          -- Duration [cite: 9]
    6.054,                                                                         -- RatingKP [cite: 9]
    7.2,                                                                           -- RatingIMDB [cite: 9]
    NULL,                                                                          -- ShortDescription is empty [cite: 8]
    NULL,                                                                          -- Logo is empty [cite: 9, 10]
    NULL,                                                                          -- Backdrop is empty [cite: 10]
    FALSE                                                                          -- Watchability is nil [cite: 9]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 9]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(18386, 400787, 'Актёр'),
(1768, 400787, 'Актёр'),
(6226, 400787, 'Актёр'),
(45338, 400787, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 9]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'фэнтези'), 400787),
((SELECT id from genre where name = 'боевик'), 400787),
((SELECT id from genre where name = 'триллер'), 400787),
((SELECT id from genre where name = 'драма'), 400787)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 448 - Форрест Гамп (Forrest Gump)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 13]
INSERT INTO genre (name) VALUES ('комедия') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('мелодрама') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('история') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('военный') ON CONFLICT (name) DO NOTHING;
-- 'драма' already handled

-- Insert Persons (Staff) if they don't exist [cite: 13]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(9144, 'Том Хэнкс', 'Tom Hanks', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_9144.jpg'),
(8887, 'Робин Райт', 'Robin Wright', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_8887.jpg'),
(13477, 'Салли Филд', 'Sally Field', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_13477.jpg'),
(3100, 'Гэри Синиз', 'Gary Sinise', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_3100.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 10, 11, 12, 13, 14]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    448,                                                                           -- ID [cite: 10]
    'Форрест Гамп',                                                                -- Name [cite: 10]
    'Forrest Gump',                                                                -- OriginalName [cite: 10]
    E'Сидя на\u00a0автобусной остановке, Форрест Гамп\u00a0—\u00a0не очень умный, но\u00a0добрый и\u00a0открытый парень\u00a0—\u00a0рассказывает случайным встречным историю своей необыкновенной жизни.\n\nС самого малолетства парень страдал от\u00a0заболевания ног, соседские мальчишки дразнили его, но\u00a0в один прекрасный день Форрест открыл в\u00a0себе невероятные способности к\u00a0бегу. Подруга детства Дженни всегда его\u00a0поддерживала и\u00a0защищала, но\u00a0вскоре дороги их\u00a0разошлись.', -- About [cite: 10, 11]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/3560b757-9b95-45ec-af8c-623972370f9d/orig', -- Poster [cite: 12]
    NULL,                                                                          -- PromoURL is empty [cite: 12]
    '1994-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 12]
    'Мир уже никогда не будет прежним, после того как вы увидите его глазами Форреста Гампа', -- Slogan [cite: 12]
    'Роберт Земекис',                                                              -- Director [cite: 12]
    'США',                                                                         -- Country [cite: 12]
    55000000,                                                                      -- Budget [cite: 12]
    329694499,                                                                     -- BoxOfficeUS [cite: 12]
    677387716,                                                                     -- BoxOfficeGlobal [cite: 12]
    84460,                                                                         -- BoxOfficeRussia [cite: 12]
    '2020-02-13'::DATE,                                                            -- PremierRussia [cite: 12, 13]
    '1994-06-23'::DATE,                                                            -- PremierGlobal [cite: 13]
    7.5,                                                                           -- Rating [cite: 13]
    '142',                                                                         -- Duration [cite: 13]
    8.921,                                                                         -- RatingKP [cite: 13]
    8.8,                                                                           -- RatingIMDB [cite: 13]
    'Полувековая история США глазами чудака из Алабамы. Абсолютная классика Роберта Земекиса с Томом Хэнксом', -- ShortDescription [cite: 11, 12]
    'https://avatars.mds.yandex.net/get-ott/200035/2a00000170ed554ce17a2db2b2cfdc134a6c/orig', -- Logo [cite: 13, 14]
    'https://image.openmoviedb.com/kinopoisk-ott-images/200035/2a0000016127256d3d2a4f6bf76eac33c45f/orig', -- Backdrop [cite: 14]
    FALSE                                                                          -- Watchability is nil [cite: 13]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 13]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(9144, 448, 'Актёр'),
(8887, 448, 'Актёр'),
(13477, 448, 'Актёр'),
(3100, 448, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 13]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'драма'), 448),
((SELECT id from genre where name = 'комедия'), 448),
((SELECT id from genre where name = 'мелодрама'), 448),
((SELECT id from genre where name = 'история'), 448),
((SELECT id from genre where name = 'военный'), 448)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 325 - Крестный отец (The Godfather)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 21]
-- 'драма', 'криминал' already handled

-- Insert Persons (Staff) if they don't exist [cite: 21]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(33037, 'Марлон Брандо', 'Marlon Brando', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_33037.jpg'),
(26240, 'Аль Пачино', 'Al Pacino', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_26240.jpg'),
(7088, 'Джеймс Каан', 'James Caan', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_7088.jpg'),
(6273, 'Роберт Дювалл', 'Robert Duvall', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_6273.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 14, 15, 16, 17, 18, 19, 20, 21, 22]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    325,                                                                           -- ID [cite: 14]
    'Крестный отец',                                                               -- Name [cite: 14]
    'The Godfather',                                                               -- OriginalName [cite: 14]
    E'Криминальная сага, повествующая о нью-йоркской сицилийской мафиозной семье Корлеоне. Фильм охватывает период 1945-1955 годов.\n\nГлава семьи, Дон Вито Корлеоне, выдаёт замуж свою дочь. В это время со Второй мировой войны возвращается его любимый сын Майкл. Майкл, герой войны, гордость семьи, не выражает желания заняться жестоким семейным бизнесом. Дон Корлеоне ведёт дела по старым правилам, но наступают иные времена, и появляются люди, желающие изменить сложившиеся порядки. На Дона Корлеоне совершается покушение.', -- About [cite: 14, 15, 16, 17, 18, 19]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/c11652e8-653b-47c1-8e72-1552399a775b/orig', -- Poster [cite: 20]
    NULL,                                                                          -- PromoURL is empty [cite: 20]
    '1972-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 20]
    'Настоящая сила не может быть дана, она может быть взята...',                  -- Slogan [cite: 20]
    'Фрэнсис Форд Коппола',                                                        -- Director [cite: 20]
    'США',                                                                         -- Country [cite: 20]
    6000000,                                                                       -- Budget [cite: 20]
    133698921,                                                                     -- BoxOfficeUS [cite: 20]
    243862778,                                                                     -- BoxOfficeGlobal [cite: 20]
    151566,                                                                        -- BoxOfficeRussia [cite: 20]
    '1992-02-03'::DATE,                                                            -- PremierRussia [cite: 20]
    '1972-03-14'::DATE,                                                            -- PremierGlobal [cite: 20, 21]
    7.5,                                                                           -- Rating [cite: 21]
    '175',                                                                         -- Duration [cite: 21]
    8.705,                                                                         -- RatingKP [cite: 21]
    9.2,                                                                           -- RatingIMDB [cite: 21]
    'В семье крупного нью-йоркского мафиози наметился кризис. Революция в гангстерском кино и начало большого эпоса', -- ShortDescription [cite: 19, 20]
    'https://avatars.mds.yandex.net/get-ott/2419418/2a00000170ed4fb3680bbf2eea1b08e2565d/orig', -- Logo [cite: 21]
    'https://image.openmoviedb.com/kinopoisk-ott-images/224348/2a0000016127199c4c157b9e5c245f697def/orig', -- Backdrop [cite: 21, 22]
    FALSE                                                                          -- Watchability is nil [cite: 21]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 21]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(33037, 325, 'Актёр'),
(26240, 325, 'Актёр'),
(7088, 325, 'Актёр'),
(6273, 325, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 21]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'драма'), 325),
((SELECT id from genre where name = 'криминал'), 325)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 258687 - Интерстеллар (Interstellar)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 24]
INSERT INTO genre (name) VALUES ('фантастика') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('приключения') ON CONFLICT (name) DO NOTHING;
-- 'драма' already handled

-- Insert Persons (Staff) if they don't exist [cite: 24]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(797, 'Мэттью Макконахи', 'Matthew McConaughey', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_797.jpg'),
(38703, 'Энн Хэтэуэй', 'Anne Hathaway', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_38703.jpg'),
(1111242, 'Джессика Честейн', 'Jessica Chastain', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1111242.jpg'),
(2007922, 'Маккензи Фой', 'Mackenzie Foy', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_2007922.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 22, 23, 24]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    258687,                                                                        -- ID [cite: 22]
    'Интерстеллар',                                                                -- Name [cite: 22]
    'Interstellar',                                                                -- OriginalName [cite: 22]
    'Когда засуха, пыльные бури и вымирание растений приводят человечество к продовольственному кризису, коллектив исследователей и учёных отправляется сквозь червоточину (которая предположительно соединяет области пространства-времени через большое расстояние) в путешествие, чтобы превзойти прежние ограничения для космических путешествий человека и найти планету с подходящими для человечества условиями.', -- About [cite: 22]
    'https://image.openmoviedb.com/kinopoisk-images/1600647/430042eb-ee69-4818-aed0-a312400a26bf/orig', -- Poster [cite: 23]
    NULL,                                                                          -- PromoURL is empty [cite: 23]
    '2014-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 23]
    'Следующий шаг человечества станет величайшим',                                 -- Slogan [cite: 23]
    'Кристофер Нолан',                                                             -- Director [cite: 23]
    'США',                                                                         -- Country [cite: 23]
    165000000,                                                                     -- Budget [cite: 23]
    192445017,                                                                     -- BoxOfficeUS [cite: 23]
    736546575,                                                                     -- BoxOfficeGlobal [cite: 23]
    26192066,                                                                      -- BoxOfficeRussia [cite: 23]
    '2014-11-06'::DATE,                                                            -- PremierRussia [cite: 23]
    '2014-10-26'::DATE,                                                            -- PremierGlobal [cite: 23]
    7.5,                                                                           -- Rating [cite: 23, 24]
    '169',                                                                         -- Duration [cite: 24]
    8.656,                                                                         -- RatingKP [cite: 24]
    8.7,                                                                           -- RatingIMDB [cite: 24]
    'Фантастический эпос про задыхающуюся Землю, космические полеты и парадоксы времени. «Оскар» за спецэффекты', -- ShortDescription [cite: 22, 23]
    'https://image.openmoviedb.com/tmdb-images/original/8YGHe69tPNbSaMbUwSo3AHkELKJ.png', -- Logo [cite: 24]
    'https://image.openmoviedb.com/kinopoisk-ott-images/1531675/2a000001862c40fa5c5394735529ee7e7188/orig', -- Backdrop [cite: 24]
    FALSE                                                                          -- Watchability is nil [cite: 24]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 24]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(797, 258687, 'Актёр'),
(38703, 258687, 'Актёр'),
(1111242, 258687, 'Актёр'),
(2007922, 258687, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 24]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'фантастика'), 258687),
((SELECT id from genre where name = 'драма'), 258687),
((SELECT id from genre where name = 'приключения'), 258687)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 342 - Криминальное чтиво (Pulp Fiction)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 28]
-- 'криминал', 'драма' already handled

-- Insert Persons (Staff) if they don't exist [cite: 28, 29]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(6479, 'Джон Траволта', 'John Travolta', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_6479.jpg'),
(7164, 'Сэмюэл Л. Джексон', 'Samuel L. Jackson', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_7164.jpg'),
(110, 'Брюс Уиллис', 'Bruce Willis', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_110.jpg'),
(29595, 'Ума Турман', 'Uma Thurman', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_29595.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 24, 25, 26, 27, 28, 29]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    342,                                                                           -- ID [cite: 24, 25]
    'Криминальное чтиво',                                                          -- Name [cite: 25]
    'Pulp Fiction',                                                                -- OriginalName [cite: 25]
    E'Двое бандитов Винсент Вега и\u00a0Джулс Винфилд ведут философские беседы в\u00a0перерывах между разборками и\u00a0решением проблем с\u00a0должниками криминального босса Марселласа Уоллеса.\nВ первой истории Винсент проводит незабываемый вечер с\u00a0женой Марселласа Мией. Во\u00a0второй Марселлас покупает боксёра Бутча Кулиджа, чтобы тот\u00a0сдал бой. В\u00a0третьей истории Винсент и\u00a0Джулс по\u00a0нелепой случайности попадают в\u00a0неприятности.', -- About [cite: 25, 26]
    'https://image.openmoviedb.com/kinopoisk-images/1900788/87b5659d-a159-4224-9bff-d5a5d109a53b/orig', -- Poster [cite: 27]
    NULL,                                                                          -- PromoURL is empty [cite: 27]
    '1994-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 27]
    'Just because you are a character doesn''t mean you have character',           -- Slogan (Note: Escaped single quote) [cite: 27]
    'Квентин Тарантино',                                                           -- Director [cite: 27]
    'США',                                                                         -- Country [cite: 27]
    8000000,                                                                       -- Budget [cite: 27]
    107928762,                                                                     -- BoxOfficeUS [cite: 27]
    213928762,                                                                     -- BoxOfficeGlobal [cite: 27]
    83843,                                                                         -- BoxOfficeRussia [cite: 27]
    '1995-09-29'::DATE,                                                            -- PremierRussia [cite: 27]
    '1994-05-21'::DATE,                                                            -- PremierGlobal [cite: 27, 28]
    7.5,                                                                           -- Rating [cite: 28]
    '154',                                                                         -- Duration [cite: 28]
    8.653,                                                                         -- RatingKP [cite: 29]
    8.9,                                                                           -- RatingIMDB [cite: 29]
    'Несколько связанных историй из жизни бандитов. Шедевр Квентина Тарантино, который изменил мировое кино', -- ShortDescription [cite: 26, 27]
    'https://imagetmdb.com/t/p/original/inuYhCBbTof4gw7f9Ized0SQQuW.png',           -- Logo [cite: 29]
    'https://image.openmoviedb.com/kinopoisk-ott-images/224348/2a0000016224ef525aabbc8bfad13709e4a8/orig', -- Backdrop [cite: 29]
    FALSE                                                                          -- Watchability is nil [cite: 29]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 28, 29]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(6479, 342, 'Актёр'),
(7164, 342, 'Актёр'),
(110, 342, 'Актёр'),
(29595, 342, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 28]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'криминал'), 342),
((SELECT id from genre where name = 'драма'), 342)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 326 - Побег из Шоушенка (The Shawshank Redemption)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 35]
-- 'драма' already handled

-- Insert Persons (Staff) if they don't exist [cite: 35, 36]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(7987, 'Тим Роббинс', 'Tim Robbins', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_7987.jpg'),
(6750, 'Морган Фриман', 'Morgan Freeman', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_6750.jpg'),
(23481, 'Боб Гантон', 'Bob Gunton', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_23481.jpg'),
(24267, 'Уильям Сэдлер', 'William Sadler', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_24267.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 29, 30, 31, 32, 33, 34, 35, 36]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    326,                                                                           -- ID [cite: 29]
    'Побег из Шоушенка',                                                           -- Name [cite: 29]
    'The Shawshank Redemption',                                                    -- OriginalName [cite: 29]
    E'Бухгалтер Энди Дюфрейн обвинён в убийстве собственной жены и её любовника. Оказавшись в тюрьме под названием Шоушенк, он сталкивается с жестокостью и беззаконием, царящими по обе стороны решётки. Каждый, кто попадает в эти стены, становится их рабом до конца жизни. Но Энди, обладающий живым умом и доброй душой, находит подход как к заключённым, так и к охранникам, добиваясь их особого к себе расположения.', -- About [cite: 29, 30, 31, 32, 33]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/0b76b2a2-d1c7-4f04-a284-80ff7bb709a4/orig', -- Poster [cite: 34]
    NULL,                                                                          -- PromoURL is empty [cite: 34]
    '1994-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 34]
    E'Страх - это кандалы.\nНадежда - это свобода',                                 -- Slogan [cite: 34, 35]
    'Фрэнк Дарабонт',                                                              -- Director [cite: 35]
    'США',                                                                         -- Country [cite: 34]
    25000000,                                                                      -- Budget [cite: 35]
    28341469,                                                                      -- BoxOfficeUS [cite: 35]
    28418687,                                                                      -- BoxOfficeGlobal [cite: 35]
    87432,                                                                         -- BoxOfficeRussia [cite: 35]
    '2019-10-24'::DATE,                                                            -- PremierRussia [cite: 35]
    '1994-09-10'::DATE,                                                            -- PremierGlobal [cite: 35]
    7.5,                                                                           -- Rating [cite: 35]
    '142',                                                                         -- Duration [cite: 35]
    9.109,                                                                         -- RatingKP [cite: 36]
    9.3,                                                                           -- RatingIMDB [cite: 36]
    'Несправедливо осужденный банкир готовит побег из тюрьмы. Тим Роббинс в выдающейся экранизации Стивена Кинга', -- ShortDescription [cite: 33, 34]
    'https://image.openmoviedb.com/tmdb-images/original/2FBRJDL06YPtoPTrwU4rTARLs76.png', -- Logo [cite: 36]
    'https://image.openmoviedb.com/kinopoisk-ott-images/1672343/2a0000016b03d1f5365474a90d26998e2a9f/orig', -- Backdrop [cite: 36]
    FALSE                                                                          -- Watchability is nil [cite: 36]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 35, 36]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(7987, 326, 'Актёр'),
(6750, 326, 'Актёр'),
(23481, 326, 'Актёр'),
(24267, 326, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 35]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'драма'), 326)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 111543 - Темный рыцарь (The Dark Knight)
-- ==========================================================================

-- Insert Genres if they don't exist [cite: 40]
-- 'фантастика', 'боевик', 'триллер', 'криминал', 'драма' already handled

-- Insert Persons (Staff) if they don't exist [cite: 40]
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(21495, 'Кристиан Бэйл', 'Christian Bale', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_21495.jpg'),
(1183, 'Хит Леджер', 'Heath Ledger', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1183.jpg'),
(6752, 'Аарон Экхарт', 'Aaron Eckhart', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_6752.jpg'),
(10384, 'Мэгги Джилленхол', 'Maggie Gyllenhaal', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10384.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie [cite: 36, 37, 38, 39, 40]
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    111543,                                                                        -- ID [cite: 36]
    'Темный рыцарь',                                                               -- Name [cite: 36]
    'The Dark Knight',                                                             -- OriginalName [cite: 36]
    E'Бэтмен поднимает ставки в войне с криминалом. С помощью лейтенанта Джима Гордона и прокурора Харви Дента он намерен очистить улицы Готэма от преступности. Сотрудничество оказывается эффективным, но скоро они обнаружат себя посреди хаоса, развязанного восходящим криминальным гением, известным напуганным горожанам под именем Джокер.', -- About [cite: 36, 37, 38]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/0fa5bf50-d5ad-446f-a599-b26d070c8b99/orig', -- Poster [cite: 39]
    NULL,                                                                          -- PromoURL is empty [cite: 39]
    '2008-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 39]
    'Добро пожаловать в мир Хаоса!',                                               -- Slogan [cite: 39]
    'Кристофер Нолан',                                                             -- Director [cite: 39]
    'США',                                                                         -- Country [cite: 39]
    185000000,                                                                     -- Budget [cite: 39]
    533345358,                                                                     -- BoxOfficeUS [cite: 39]
    1003045358,                                                                    -- BoxOfficeGlobal [cite: 39]
    8589100,                                                                       -- BoxOfficeRussia [cite: 39]
    '2008-08-14'::DATE,                                                            -- PremierRussia [cite: 39]
    '2008-07-14'::DATE,                                                            -- PremierGlobal [cite: 39]
    7.5,                                                                           -- Rating [cite: 39, 40]
    '152',                                                                         -- Duration [cite: 40]
    8.53,                                                                          -- RatingKP [cite: 40]
    9.0,                                                                           -- RatingIMDB [cite: 40]
    'У Бэтмена появляется новый враг — философ-террорист Джокер. Кинокомикс, который вывел жанр на новый уровень', -- ShortDescription [cite: 38, 39]
    'https://avatars.mds.yandex.net/get-ott/224348/2a00000176f159505eff31a41fe3e4ccf723/orig', -- Logo [cite: 40]
    'https://image.openmoviedb.com/kinopoisk-ott-images/374297/2a000001670d5f5d21a5ed5475c6a524aa8a/orig', -- Backdrop [cite: 40]
    FALSE                                                                          -- Watchability is nil [cite: 40]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie [cite: 40]
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(21495, 111543, 'Актёр'),
(1183, 111543, 'Актёр'),
(6752, 111543, 'Актёр'),
(10384, 111543, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie [cite: 40]
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'фантастика'), 111543),
((SELECT id from genre where name = 'боевик'), 111543),
((SELECT id from genre where name = 'триллер'), 111543),
((SELECT id from genre where name = 'криминал'), 111543),
((SELECT id from genre where name = 'драма'), 111543)
ON CONFLICT (genre_id, movie_id) DO NOTHING;


-- Commit Transaction
COMMIT;


-- Start Transaction
BEGIN;

-- ==========================================================================
-- Movie: 1108577 - Зеленая книга (Green Book)
-- ==========================================================================

-- Insert Genres if they don't exist
INSERT INTO genre (name) VALUES ('биография') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('комедия') ON CONFLICT (name) DO NOTHING;
INSERT INTO genre (name) VALUES ('драма') ON CONFLICT (name) DO NOTHING;

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(10779, 'Вигго Мортенсен', 'Viggo Mortensen', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10779.jpg'),
(542248, 'Махершала Али', 'Mahershala Ali', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_542248.jpg'),
(27320, 'Линда Карделлини', 'Linda Cardellini', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_27320.jpg'),
(1147285, 'Себастьян Манискалко', 'Sebastian Maniscalco', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1147285.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    1108577,                                                                       -- ID [cite: 1]
    'Зеленая книга',                                                               -- Name [cite: 1]
    'Green Book',                                                                  -- OriginalName [cite: 1]
    E'1960-е годы. После закрытия нью-йоркского ночного клуба на\u00a0ремонт вышибала Тони по\u00a0прозвищу Болтун ищет подработку на\u00a0пару месяцев. Как\u00a0раз в\u00a0это время Дон\u00a0Ширли\u00a0—\u00a0утонченный светский лев, богатый и\u00a0талантливый чернокожий музыкант, исполняющий классическую музыку\u00a0—\u00a0собирается в\u00a0турне по\u00a0южным штатам, где\u00a0ещё сильны расистские убеждения и\u00a0царит сегрегация. Он\u00a0нанимает Тони в\u00a0качестве водителя, телохранителя и\u00a0человека, способного решать текущие проблемы. У\u00a0этих двоих так\u00a0мало общего, и\u00a0эта поездка навсегда изменит жизнь обоих.', -- About [cite: 1, 2]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/4b27e219-a8a5-4d85-9874-57d6016e0837/orig', -- Poster [cite: 2]
    NULL,                                                                          -- PromoURL is empty [cite: 2]
    '2018-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 2]
    'Inspired by a True Friendship',                                               -- Slogan [cite: 2]
    'Питер Фаррелли',                                                              -- Director [cite: 2]
    'США',                                                                         -- Country [cite: 2]
    23000000,                                                                      -- Budget [cite: 2]
    85080171,                                                                      -- BoxOfficeUS [cite: 2]
    321752656,                                                                     -- BoxOfficeGlobal [cite: 3]
    9319859,                                                                       -- BoxOfficeRussia [cite: 3]
    '2019-01-24'::DATE,                                                            -- PremierRussia [cite: 3]
    '2018-09-11'::DATE,                                                            -- PremierGlobal [cite: 3]
    7.5,                                                                           -- Rating [cite: 3]
    '130',                                                                         -- Duration [cite: 3]
    8.517,                                                                         -- RatingKP [cite: 4]
    8.2,                                                                           -- RatingIMDB [cite: 4]
    'Путешествие итальянца-вышибалы и чернокожего пианиста по Америке 1960-х. «Оскар» за лучший фильм', -- ShortDescription [cite: 2]
    'https://image.openmoviedb.com/tmdb-images/original/teZDEtsuxhmeyuEyv6Ww5TrdHJi.png', -- Logo [cite: 4]
    'https://image.openmoviedb.com/kinopoisk-ott-images/224348/2a00000169233a1e496732daa6f91dc207b9/orig', -- Backdrop [cite: 4]
    FALSE                                                                          -- Watchability is nil [cite: 4]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(10779, 1108577, 'Актёр'),
(542248, 1108577, 'Актёр'),
(27320, 1108577, 'Актёр'),
(1147285, 1108577, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'биография'), 1108577),
((SELECT id from genre where name = 'комедия'), 1108577),
((SELECT id from genre where name = 'драма'), 1108577)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 725190 - Одержимость (Whiplash)
-- ==========================================================================

-- Insert Genres if they don't exist
INSERT INTO genre (name) VALUES ('музыка') ON CONFLICT (name) DO NOTHING;
-- 'драма' already handled

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(1669945, 'Майлз Теллер', 'Miles Teller', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1669945.jpg'),
(8552, 'Дж.К. Симмонс', 'J.K. Simmons', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_8552.jpg'),
(23944, 'Пол Райзер', 'Paul Reiser', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_23944.jpg'),
(1207946, 'Мелисса Бенойст', 'Melissa Benoist', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1207946.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    725190,                                                                        -- ID [cite: 4]
    'Одержимость',                                                                 -- Name [cite: 4]
    'Whiplash',                                                                    -- OriginalName [cite: 4]
    E'Эндрю мечтает стать великим. Казалось бы, вот-вот его мечта осуществится. Юношу замечает настоящий гений, дирижер лучшего в стране оркестра. Желание Эндрю добиться успеха быстро становится одержимостью, а безжалостный наставник продолжает подталкивать его все дальше и дальше – за пределы человеческих возможностей. Кто выйдет победителем из этой схватки?', -- About [cite: 5, 6, 7]
    'https://image.openmoviedb.com/kinopoisk-images/6201401/16af46be-bcfe-461e-af54-ff17b905b82e/orig', -- Poster [cite: 8]
    NULL,                                                                          -- PromoURL is empty [cite: 8]
    '2013-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 8]
    'The road to greatness can take you to the edge',                              -- Slogan [cite: 8]
    'Дэмьен Шазелл',                                                               -- Director [cite: 8]
    'США',                                                                         -- Country [cite: 8]
    3300000,                                                                       -- Budget [cite: 8]
    13092000,                                                                      -- BoxOfficeUS [cite: 8]
    48982041,                                                                      -- BoxOfficeGlobal [cite: 8]
    189987,                                                                        -- BoxOfficeRussia [cite: 8]
    '2014-10-23'::DATE,                                                            -- PremierRussia [cite: 8]
    '2014-01-16'::DATE,                                                            -- PremierGlobal [cite: 9]
    7.5,                                                                           -- Rating [cite: 9]
    '106',                                                                         -- Duration [cite: 9]
    8.362,                                                                         -- RatingKP [cite: 10]
    8.5,                                                                           -- RatingIMDB [cite: 10]
    'Юный барабанщик на тернистом пути к величию. Остросюжетная драма Дэмьена Шазелла, отмеченная тремя «Оскарами»', -- ShortDescription [cite: 8]
    NULL,                                                                          -- Logo is empty [cite: 10]
    'https://image.openmoviedb.com/kinopoisk-ott-images/1534341/2a0000018737410431c24b72c399351bdeeb/orig', -- Backdrop [cite: 10]
    FALSE                                                                          -- Watchability is nil [cite: 10]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(1669945, 725190, 'Актёр'),
(8552, 725190, 'Актёр'),
(23944, 725190, 'Актёр'),
(1207946, 725190, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'драма'), 725190),
((SELECT id from genre where name = 'музыка'), 725190)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 6462 - Рокки (Rocky)
-- ==========================================================================

-- Insert Genres if they don't exist
INSERT INTO genre (name) VALUES ('спорт') ON CONFLICT (name) DO NOTHING;
-- 'драма' already handled

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(8816, 'Сильвестр Сталлоне', 'Sylvester Stallone', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_8816.jpg'),
(33823, 'Талия Шайр', 'Talia Shire', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_33823.jpg'),
(25025, 'Берт Янг', 'Burt Young', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_25025.jpg'),
(7639, 'Карл Уэзерс', 'Carl Weathers', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_7639.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    6462,                                                                          -- ID [cite: 10]
    'Рокки',                                                                       -- Name [cite: 10]
    'Rocky',                                                                       -- OriginalName [cite: 10]
    E'Филадельфия. Рокки Бальбоа\u00a0—\u00a0молодой боксёр, который живёт в\u00a0захудалой квартирке и\u00a0еле сводит концы с\u00a0концами, занимаясь выбиванием долгов для\u00a0своего босса Тони Гаццо и\u00a0периодически участвуя в\u00a0боях. Каждый его\u00a0унылый день похож на\u00a0предыдущий, и\u00a0особо радужных перспектив не\u00a0наблюдается. Но\u00a0однажды удача наконец улыбается парню, когда ему\u00a0поступает неожиданное предложение выступить против действующего чемпиона Аполло Крида.', -- About [cite: 11, 12]
    'https://image.openmoviedb.com/kinopoisk-images/6201401/2ec74380-dda9-4f1c-a228-1475061a5152/orig', -- Poster [cite: 13]
    NULL,                                                                          -- PromoURL is empty [cite: 13]
    '1976-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 13]
    'His whole life was a million-to-one shot',                                    -- Slogan [cite: 13]
    'Джон Г. Эвилдсен',                                                            -- Director [cite: 14]
    'США',                                                                         -- Country [cite: 13]
    1100000,                                                                       -- Budget [cite: 14]
    117235147,                                                                     -- BoxOfficeUS [cite: 14]
    225000000,                                                                     -- BoxOfficeGlobal [cite: 14]
    87432,                                                                         -- BoxOfficeRussia [cite: 14]
    NULL,                                                                          -- PremierRussia is empty [cite: 14]
    '1976-11-20'::DATE,                                                            -- PremierGlobal [cite: 14]
    7.5,                                                                           -- Rating [cite: 14]
    '120',                                                                         -- Duration [cite: 14]
    8.02,                                                                          -- RatingKP [cite: 15]
    8.1,                                                                           -- RatingIMDB [cite: 15]
    'Вышибала-боксер выходит на большой бой, чтобы стать легендой. Рождение звезды Сильвестра Сталлоне', -- ShortDescription [cite: 13]
    'https://avatars.mds.yandex.net/get-ott/1534341/2a0000017e7ce2ea0f17a5db3974724355a7/orig', -- Logo [cite: 15]
    'https://image.openmoviedb.com/kinopoisk-ott-images/374297/2a0000017e7cd46446d2a92ed0e27e0c7f66/orig', -- Backdrop [cite: 15]
    FALSE                                                                          -- Watchability is nil [cite: 15]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(8816, 6462, 'Актёр'),
(33823, 6462, 'Актёр'),
(25025, 6462, 'Актёр'),
(7639, 6462, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'драма'), 6462),
((SELECT id from genre where name = 'спорт'), 6462)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 635772 - Игра в имитацию (The Imitation Game)
-- ==========================================================================

-- Insert Genres if they don't exist
-- 'биография', 'военный', 'драма', 'история' already handled

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(34549, 'Бенедикт Камбербэтч', 'Benedict Cumberbatch', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_34549.jpg'),
(24302, 'Кира Найтли', 'Keira Knightley', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_24302.jpg'),
(397499, 'Мэттью Гуд', 'Matthew Goode', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_397499.jpg'),
(722361, 'Рори Киннер', 'Rory Kinnear', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_722361.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    635772,                                                                        -- ID [cite: 15]
    'Игра в имитацию',                                                             -- Name [cite: 15]
    'The Imitation Game',                                                          -- OriginalName [cite: 15]
    'Английский математик и логик Алан Тьюринг пытается взломать код немецкой шифровальной машины Enigma во время Второй мировой войны.', -- About [cite: 15]
    'https://image.openmoviedb.com/kinopoisk-images/10703859/8ac3b989-2556-42bb-9d11-913d42bacc63/orig', -- Poster [cite: 16]
    NULL,                                                                          -- PromoURL is empty [cite: 16]
    '2014-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 16]
    'Основано на невероятной, но реальной истории',                                -- Slogan [cite: 16]
    'Мортен Тильдум',                                                              -- Director [cite: 16]
    'Великобритания',                                                              -- Country [cite: 16]
    14000000,                                                                      -- Budget [cite: 16]
    91125683,                                                                      -- BoxOfficeUS [cite: 16]
    233555708,                                                                     -- BoxOfficeGlobal [cite: 16]
    3171684,                                                                       -- BoxOfficeRussia [cite: 16]
    '2015-02-05'::DATE,                                                            -- PremierRussia [cite: 16]
    '2014-08-29'::DATE,                                                            -- PremierGlobal [cite: 16]
    7.5,                                                                           -- Rating [cite: 17]
    '114',                                                                         -- Duration [cite: 17]
    7.752,                                                                         -- RatingKP [cite: 17]
    8.0,                                                                           -- RatingIMDB [cite: 17]
    'Математики должны изобрести дешифратор, чтобы остановить войну. Оскароносная драма с Бенедиктом Камбербэтчем', -- ShortDescription [cite: 16]
    'https://image.openmoviedb.com/tmdb-images/original/knWKYhJhmrJowHyaeVxJsTys3wN.png', -- Logo [cite: 17]
    'https://image.openmoviedb.com/kinopoisk-ott-images/13051577/2a000001900c6d14d65d6157d7290ad66e92/orig', -- Backdrop [cite: 17]
    FALSE                                                                          -- Watchability is nil [cite: 17]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(34549, 635772, 'Актёр'),
(24302, 635772, 'Актёр'),
(397499, 635772, 'Актёр'),
(722361, 635772, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'биография'), 635772),
((SELECT id from genre where name = 'военный'), 635772),
((SELECT id from genre where name = 'драма'), 635772),
((SELECT id from genre where name = 'история'), 635772)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 447301 - Начало (Inception)
-- ==========================================================================

-- Insert Genres if they don't exist
INSERT INTO genre (name) VALUES ('детектив') ON CONFLICT (name) DO NOTHING;
-- 'фантастика', 'боевик', 'триллер', 'драма' already handled

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(37859, 'Леонардо ДиКаприо', 'Leonardo DiCaprio', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_37859.jpg'),
(9867, 'Джозеф Гордон-Левитт', 'Joseph Gordon-Levitt', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_9867.jpg'),
(43503, 'Эллиот Пейдж', 'Elliot Page', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_43503.jpg'),
(39984, 'Том Харди', 'Tom Hardy', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_39984.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    447301,                                                                        -- ID [cite: 18]
    'Начало',                                                                      -- Name [cite: 18]
    'Inception',                                                                   -- OriginalName [cite: 18]
    E'Кобб – талантливый вор, лучший из лучших в опасном искусстве извлечения: он крадет ценные секреты из глубин подсознания во время сна, когда человеческий разум наиболее уязвим. Редкие способности Кобба сделали его ценным игроком в привычном к предательству мире промышленного шпионажа, но они же превратили его в извечного беглеца и лишили всего, что он когда-либо любил. \n\nИ вот у Кобба появляется шанс исправить ошибки. Его последнее дело может вернуть все назад, но для этого ему нужно совершить невозможное – инициацию. Вместо идеальной кражи Кобб и его команда спецов должны будут провернуть обратное. Теперь их задача – не украсть идею, а внедрить ее. Если у них получится, это и станет идеальным преступлением. \n\nНо никакое планирование или мастерство не могут подготовить команду к встрече с опасным противником, который, кажется, предугадывает каждый их ход. Врагом, увидеть которого мог бы лишь Кобб.', -- About [cite: 18, 19, 20, 21, 22, 23, 24]
    'https://image.openmoviedb.com/kinopoisk-images/1629390/8ab9a119-dd74-44f0-baec-0629797483d7/orig', -- Poster [cite: 25]
    NULL,                                                                          -- PromoURL is empty [cite: 25]
    '2010-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 25]
    'Твой разум - место преступления',                                             -- Slogan [cite: 25]
    'Кристофер Нолан',                                                             -- Director [cite: 25]
    'США',                                                                         -- Country [cite: 25]
    160000000,                                                                     -- Budget [cite: 25]
    292576195,                                                                     -- BoxOfficeUS [cite: 25]
    828322032,                                                                     -- BoxOfficeGlobal [cite: 25]
    21691531,                                                                      -- BoxOfficeRussia [cite: 25]
    '2010-07-22'::DATE,                                                            -- PremierRussia [cite: 25]
    '2010-07-08'::DATE,                                                            -- PremierGlobal [cite: 25]
    7.5,                                                                           -- Rating [cite: 26]
    '148',                                                                         -- Duration [cite: 26]
    8.667,                                                                         -- RatingKP [cite: 26]
    8.8,                                                                           -- RatingIMDB [cite: 26]
    'Профессиональные воры внедряются в сон наследника огромной империи. Фантастический боевик Кристофера Нолана', -- ShortDescription [cite: 25]
    'https://avatars.mds.yandex.net/get-ott/200035/2a00000178c5fc5e63481655114331b766a3/orig', -- Logo [cite: 26]
    'https://image.openmoviedb.com/kinopoisk-ott-images/2439731/2a00000178c5fc1db03b03dceda170357652/orig', -- Backdrop [cite: 26]
    FALSE                                                                          -- Watchability is nil [cite: 26]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(37859, 447301, 'Актёр'),
(9867, 447301, 'Актёр'),
(43503, 447301, 'Актёр'),
(39984, 447301, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'фантастика'), 447301),
((SELECT id from genre where name = 'боевик'), 447301),
((SELECT id from genre where name = 'триллер'), 447301),
((SELECT id from genre where name = 'драма'), 447301),
((SELECT id from genre where name = 'детектив'), 447301)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 476 - Назад в будущее (Back to the Future)
-- ==========================================================================

-- Insert Genres if they don't exist
-- 'фантастика', 'комедия', 'приключения' already handled

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(181, 'Майкл Дж. Фокс', 'Michael J. Fox', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_181.jpg'),
(3514, 'Кристофер Ллойд', 'Christopher Lloyd', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_3514.jpg'),
(77429, 'Лиа Томпсон', 'Lea Thompson', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_77429.jpg'),
(3558, 'Криспин Гловер', 'Crispin Glover', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_3558.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    476,                                                                           -- ID [cite: 27]
    'Назад в будущее',                                                             -- Name [cite: 27]
    'Back to the Future',                                                          -- OriginalName [cite: 27]
    E'Подросток Марти с помощью машины времени, сооружённой его другом-профессором доком Брауном, попадает из 80-х в далекие 50-е. Там он встречается со своими будущими родителями, ещё подростками, и другом-профессором, совсем молодым.', -- About [cite: 27, 28]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/73cf2ed0-fd52-47a2-9e26-74104360786a/orig', -- Poster [cite: 29]
    NULL,                                                                          -- PromoURL is empty [cite: 29]
    '1985-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 29]
    E'Семнадцатилетний Марти МакФлай пришел вчера домой пораньше. На 30 лет раньше', -- Slogan [cite: 29, 30]
    'Роберт Земекис',                                                              -- Director [cite: 30]
    'США',                                                                         -- Country [cite: 30]
    19000000,                                                                      -- Budget [cite: 30]
    210609762,                                                                     -- BoxOfficeUS [cite: 30]
    381109762,                                                                     -- BoxOfficeGlobal [cite: 30]
    0,                                                                             -- BoxOfficeRussia [cite: 30]
    '2013-06-27'::DATE,                                                            -- PremierRussia [cite: 30]
    '1985-07-03'::DATE,                                                            -- PremierGlobal [cite: 30]
    7.5,                                                                           -- Rating [cite: 30]
    '116',                                                                         -- Duration [cite: 30]
    8.647,                                                                         -- RatingKP [cite: 31]
    8.5,                                                                           -- RatingIMDB [cite: 31]
    'Безумный ученый и 17-летний оболтус тестируют машину времени, наводя шороху в 1950-х. Классика кинофантастики', -- ShortDescription [cite: 29]
    'https://avatars.mds.yandex.net/get-ott/223007/2a0000016eadd2590232b1f0f0ea7b27b0b0/orig', -- Logo [cite: 31]
    'https://image.openmoviedb.com/kinopoisk-ott-images/223007/2a0000016128aa9721be7e9d22f4bb2be635/orig', -- Backdrop [cite: 31]
    FALSE                                                                          -- Watchability is nil [cite: 31]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(181, 476, 'Актёр'),
(3514, 476, 'Актёр'),
(77429, 476, 'Актёр'),
(3558, 476, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'фантастика'), 476),
((SELECT id from genre where name = 'комедия'), 476),
((SELECT id from genre where name = 'приключения'), 476)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 474 - Гладиатор (Gladiator)
-- ==========================================================================

-- Insert Genres if they don't exist
-- 'история', 'боевик', 'драма', 'приключения' already handled

-- Insert Persons (Staff) if they don't exist
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(10019, 'Рассел Кроу', 'Russell Crowe', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10019.jpg'),
(10020, 'Хоакин Феникс', 'Joaquin Phoenix', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10020.jpg'),
(7988, 'Конни Нильсен', 'Connie Nielsen', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_7988.jpg'),
(10021, 'Оливер Рид', 'Oliver Reed', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10021.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    474,                                                                           -- ID [cite: 31]
    'Гладиатор',                                                                   -- Name [cite: 31]
    'Gladiator',                                                                   -- OriginalName [cite: 31]
    E'В великой Римской империи не было военачальника, равного генералу Максимусу. Непобедимые легионы, которыми командовал этот благородный воин, боготворили его и могли последовать за ним даже в ад. Но случилось так, что отважный Максимус, готовый сразиться с любым противником в честном бою, оказался бессилен против вероломных придворных интриг. Генерала предали и приговорили к смерти. Чудом избежав гибели, Максимус становится гладиатором. Быстро снискав себе славу в кровавых поединках, он оказывается в знаменитом римском Колизее, на арене которого он встретится в смертельной схватке со своим заклятым врагом...', -- About [cite: 32, 33, 34, 35, 36]
    'https://image.openmoviedb.com/kinopoisk-images/1599028/7c3460dc-344d-433f-8220-f18d86c8397d/orig', -- Poster [cite: 37]
    NULL,                                                                          -- PromoURL is empty [cite: 37]
    '2000-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 37]
    E'«Генерал, Ставший Рабом. Раб, Ставший Гладиатором»',                           -- Slogan [cite: 37, 38]
    'Ридли Скотт',                                                                 -- Director [cite: 38]
    'США',                                                                         -- Country [cite: 37]
    103000000,                                                                     -- Budget [cite: 38]
    187705427,                                                                     -- BoxOfficeUS [cite: 38]
    460583960,                                                                     -- BoxOfficeGlobal [cite: 38]
    1280000,                                                                       -- BoxOfficeRussia [cite: 38]
    '2000-05-18'::DATE,                                                            -- PremierRussia [cite: 38]
    '2000-05-01'::DATE,                                                            -- PremierGlobal [cite: 38]
    7.5,                                                                           -- Rating [cite: 38]
    '155',                                                                         -- Duration [cite: 38]
    8.589,                                                                         -- RatingKP [cite: 39]
    8.5,                                                                           -- RatingIMDB [cite: 39]
    'Отважный генерал, ставший рабом, мстит империи. Культовая историческая драма Ридли Скотта с пятью «Оскарами»', -- ShortDescription [cite: 37]
    'https://image.openmoviedb.com/tmdb-images/original/cesHtsrDcrjFhSsGKrckKhTLTyp.png', -- Logo [cite: 39]
    'https://image.openmoviedb.com/tmdb-images/original/Ar7QuJ7sJEiC0oP3I8fKBKIQD9u.jpg', -- Backdrop [cite: 39]
    FALSE                                                                          -- Watchability is nil [cite: 39]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(10019, 474, 'Актёр'),
(10020, 474, 'Актёр'),
(7988, 474, 'Актёр'),
(10021, 474, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'история'), 474),
((SELECT id from genre where name = 'боевик'), 474),
((SELECT id from genre where name = 'драма'), 474),
((SELECT id from genre where name = 'приключения'), 474)
ON CONFLICT (genre_id, movie_id) DO NOTHING;

-- ==========================================================================
-- Movie: 2213 - Титаник (Titanic)
-- ==========================================================================

-- Insert Genres if they don't exist
-- 'мелодрама', 'история', 'триллер', 'драма' already handled

-- Insert Persons (Staff) if they don't exist
-- Note: Person 37859 (Leonardo DiCaprio) already inserted for movie 447301
INSERT INTO person (id, full_name, en_full_name, photo) OVERRIDING SYSTEM VALUE VALUES
(37859, 'Леонардо ДиКаприо', 'Leonardo DiCaprio', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_37859.jpg'),
(21709, 'Кейт Уинслет', 'Kate Winslet', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_21709.jpg'),
(45019, 'Билли Зейн', 'Billy Zane', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_45019.jpg'),
(379, 'Кэти Бейтс', 'Kathy Bates', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_379.jpg')
ON CONFLICT (id) DO NOTHING;

-- Insert Movie
INSERT INTO movie (
    id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, watchability
)
OVERRIDING SYSTEM VALUE VALUES (
    2213,                                                                          -- ID [cite: 39]
    'Титаник',                                                                     -- Name [cite: 39]
    'Titanic',                                                                     -- OriginalName [cite: 39]
    E'Апрель 1912 года. В\u00a0первом и\u00a0последнем плавании шикарного «Титаника» встречаются двое. Пассажир нижней палубы Джек выиграл билет в\u00a0карты, а\u00a0богатая наследница Роза отправляется в\u00a0Америку, чтобы выйти замуж по\u00a0расчёту. Чувства молодых людей только успевают расцвести, и\u00a0даже не\u00a0классовые различия создадут испытания влюблённым, а\u00a0айсберг, вставший на\u00a0пути считавшегося непотопляемым лайнера.', -- About [cite: 40, 41]
    'https://image.openmoviedb.com/kinopoisk-images/10592371/7f0e6761-4635-46ad-b804-59d5cf1ae85c/orig', -- Poster [cite: 42]
    NULL,                                                                          -- PromoURL is empty [cite: 42]
    '1997-01-01 00:00:00+00'::TIMESTAMPTZ,                                         -- ReleaseYear [cite: 42]
    'Ничто на Земле не сможет разлучить их',                                       -- Slogan [cite: 42]
    'Джеймс Кэмерон',                                                              -- Director [cite: 42]
    'США',                                                                         -- Country [cite: 42]
    200000000,                                                                     -- Budget [cite: 42]
    674292608,                                                                     -- BoxOfficeUS [cite: 42]
    2264743305,                                                                    -- BoxOfficeGlobal [cite: 42]
    18400000,                                                                      -- BoxOfficeRussia [cite: 42]
    '1998-02-20'::DATE,                                                            -- PremierRussia [cite: 42]
    '1997-11-01'::DATE,                                                            -- PremierGlobal [cite: 42]
    7.5,                                                                           -- Rating [cite: 43]
    '194',                                                                         -- Duration [cite: 43]
    8.386,                                                                         -- RatingKP [cite: 43]
    7.9,                                                                           -- RatingIMDB [cite: 43]
    'Запретная любовь на фоне гибели легендарного лайнера. Великий фильм-катастрофа — в отреставрированной версии', -- ShortDescription [cite: 42]
    'https://image.openmoviedb.com/tmdb-images/original/pfaZhIXrJQmfsMTxrC3o7kTHBuD.png', -- Logo [cite: 43]
    'https://image.openmoviedb.com/kinopoisk-ott-images/200035/2a00000161288210c88a7eb838405d44fd9d/orig', -- Backdrop [cite: 43]
    FALSE                                                                          -- Watchability is nil [cite: 43]
) ON CONFLICT (id) DO NOTHING;

-- Link Staff to Movie
INSERT INTO movie_staff (staff_id, movie_id, role) VALUES
(37859, 2213, 'Актёр'),
(21709, 2213, 'Актёр'),
(45019, 2213, 'Актёр'),
(379, 2213, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Link Genres to Movie
INSERT INTO movie_genre (genre_id, movie_id) VALUES
((SELECT id from genre where name = 'мелодрама'), 2213),
((SELECT id from genre where name = 'история'), 2213),
((SELECT id from genre where name = 'триллер'), 2213),
((SELECT id from genre where name = 'драма'), 2213)
ON CONFLICT (genre_id, movie_id) DO NOTHING;


-- Commit Transaction
COMMIT;


-- Start Transaction
BEGIN;

-- Inserting data into the movie table
INSERT INTO movie (id, name, original_name, about, poster, promo_poster, release_year, slogan, director, country, budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global, rating, duration, rating_kp, rating_imdb, short_description, logo, backdrop, created_at, updated_at)
OVERRIDING SYSTEM VALUE VALUES
(920265, 'Человек-паук: Через вселенные', 'Spider-Man: Into the Spider-Verse', 'В центре сюжета уникальной и инновационной в визуальном плане картины подросток из Нью-Йорка Майлз Моралес, который живет в мире безграничных возможностей вселенных Человека-паука, где костюм супергероя носит не только он.', 'https://image.openmoviedb.com/kinopoisk-images/1900788/64417bd3-838b-4910-a9f4-c278c509d568/orig', '/static/img/promo_4k_spider-man.webp', '2018-01-01', 'Более одного носит маску', 'Боб Персичетти', 'Канада', 90000000, 190241310, 375582637, 6298173, '2018-12-13', '2018-12-12', 7.5, '117', 8.178, 8.4, 'Пауки из разных измерений объединяются перед общей угрозой. Изобретательный кинокомикс с «Оскаром» за анимацию', 'https://image.openmoviedb.com/tmdb-images/original/bJ3VpPP3VkJM9H8GfRK8wSvTwPy.png', 'https://image.openmoviedb.com/tmdb-images/original/8mnXR9rey5uQ08rZAvzojKWbDQS.jpg', DEFAULT, DEFAULT),
(841081, 'Ла-Ла Ленд', 'La La Land', 'Это история любви старлетки, которая между прослушиваниями подает кофе состоявшимся кинозвездам, и фанатичного джазового музыканта, вынужденного подрабатывать в заштатных барах. Но пришедший к влюбленным успех начинает подтачивать их отношения.', 'https://image.openmoviedb.com/kinopoisk-images/10835644/7e786437-eada-4d23-baea-f3a5ebf57e06/orig', '/static/img/promo_4k_lalaland.webp', '2016-01-01', 'Бесстрашным мечтателям посвящается…', 'Дэмьен Шазелл', 'США', 30000000, 151101803, 447407695, 5913961, '2017-01-12', '2016-08-31', 7.5, '128', 8.029, 8.0, 'Миа и Себастьян выбирают между личным счастьем и амбициями. Трагикомичный мюзикл о компромиссе в жизни артиста', 'https://image.openmoviedb.com/tmdb-images/original/97DYpdDmiOpfn9l9I2dGxDH01t1.svg', 'https://image.openmoviedb.com/kinopoisk-ott-images/224348/2a00000161286a724010933d8658b746baa2/orig', DEFAULT, DEFAULT),
(854, 'Такси 2', 'Taxi 2', 'Во Францию прибывает министр обороны Японии. Цель его визита - ознакомиться с французским опытом борьбы с терроризмом и подписать «контракт века» о взаимном сотрудничестве.\n\nНеожиданно во время показательных выступлений французской полиции министра обороны похищают якудза, желающая сорвать заключение наиважнейшего контракта. Даниэль и Эмильен отправляются на поиски высокого гостя. В дело вступает уже хорошо знакомое нам такси.', 'https://image.openmoviedb.com/kinopoisk-images/1629390/a89ba16c-271c-4fb9-b1a5-806449bc3d62/orig', '/static/img/promo_4k_taxi.webp', '2000-01-01', 'Le 29 mars, il passe la seconde', 'Жерар Кравчик', 'Франция', 12000000, 626164, 60726164, 580000, '2000-12-21', '2000-03-25', 7.5, '88', 7.735, 6.5, 'Жажда скорости Даниэля поможет спасти японского министра от якудза. Сиквел комедийного проекта Люка Бессона', 'https://image.openmoviedb.com/tmdb-images/original/fOmbtIqbmbsAhAfkm1LR0kM7phb.png', 'https://image.openmoviedb.com/kinopoisk-ott-images/224348/2a0000016b2c54c6a65e64b1585f126b5a49/orig', DEFAULT, DEFAULT),
(1048334, 'Джокер', 'Joker', 'Готэм, начало 1980-х годов. Комик Артур Флек живет с больной матерью, которая с детства учит его «ходить с улыбкой». Пытаясь нести в мир хорошее и дарить людям радость, Артур сталкивается с человеческой жестокостью и постепенно приходит к выводу, что этот мир получит от него не добрую улыбку, а ухмылку злодея Джокера.', 'https://image.openmoviedb.com/kinopoisk-images/1946459/84934543-5991-4c93-97eb-beb6186a3ad7/orig', '/static/img/promo_4k_joker.webp', '2019-01-01', '«Сделай счастливое лицо»', 'Тодд Филлипс', 'США', 55000000, 335451311, 1078751311, 31418225, '2019-10-03', '2019-08-31', 7.5, '122', 7.996, 8.3, 'Как неудачливый комик стал самым опасным человеком в Готэме. Бенефис Хоакина Феникса и «Оскар» за саундтрек', 'https://image.openmoviedb.com/tmdb-images/original/tDzUDdxHk4yNcZskimqHQbD1JZ1.png', 'https://image.openmoviedb.com/tmdb-images/original/rlay2M5QYvi6igbGcFjq8jxeusY.jpg', DEFAULT, DEFAULT)
ON CONFLICT (id) DO NOTHING;

-- Inserting data into the person table
INSERT INTO person (id, full_name, en_full_name, photo, created_at, updated_at)
OVERRIDING SYSTEM VALUE VALUES
(2378404, 'Шамеик Мур', 'Shameik Moore', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_2378404.jpg', DEFAULT, DEFAULT),
(1089330, 'Джейк Джонсон', 'Jake Johnson', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_1089330.jpg', DEFAULT, DEFAULT),
(1478559, 'Хейли Стайнфелд', 'Hailee Steinfeld', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_1478559.jpg', DEFAULT, DEFAULT),
(542248, 'Махершала Али', 'Mahershala Ali', 'https://st.kp.yandex.net/images/actor_iphone/iphone360_542248.jpg', DEFAULT, DEFAULT),
(10143, 'Райан Гослинг', 'Ryan Gosling', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10143.jpg', DEFAULT, DEFAULT),
(1130955, 'Эмма Стоун', 'Emma Stone', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1130955.jpg', DEFAULT, DEFAULT),
(1074124, 'Джон Ледженд', 'John Legend', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_1074124.jpg', DEFAULT, DEFAULT),
(8552, 'Дж.К. Симмонс', 'J.K. Simmons', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_8552.jpg', DEFAULT, DEFAULT),
(55289, 'Сами Насери', 'Samy Naceri', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_55289.jpg', DEFAULT, DEFAULT),
(2365, 'Фредерик Дифенталь', 'Frédéric Diefenthal', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_2365.jpg', DEFAULT, DEFAULT),
(32545, 'Марион Котийяр', 'Marion Cotillard', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_32545.jpg', DEFAULT, DEFAULT),
(55290, 'Эмма Виклунд', 'Emma Wiklund', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_55290.jpg', DEFAULT, DEFAULT),
(10020, 'Хоакин Феникс', 'Joaquin Phoenix', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_10020.jpg', DEFAULT, DEFAULT),
(718, 'Роберт Де Ниро', 'Robert De Niro', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_718.jpg', DEFAULT, DEFAULT),
(3394604, 'Зази Битц', 'Zazie Beetz', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_3394604.jpg', DEFAULT, DEFAULT),
(40447, 'Фрэнсис Конрой', 'Frances Conroy', 'https://image.openmoviedb.com/kinopoisk-st-images//actor_iphone/iphone360_40447.jpg', DEFAULT, DEFAULT)
ON CONFLICT (id) DO NOTHING;


-- Inserting data into the movie_staff table
INSERT INTO movie_staff (staff_id, movie_id, role)
VALUES
(2378404, 920265, 'Актёр'),
(1089330, 920265, 'Актёр'),
(1478559, 920265, 'Актёр'),
(542248, 920265, 'Актёр'),
(10143, 841081, 'Актёр'),
(1130955, 841081, 'Актёр'),
(1074124, 841081, 'Актёр'),
(8552, 841081, 'Актёр'),
(55289, 854, 'Актёр'),
(2365, 854, 'Актёр'),
(32545, 854, 'Актёр'),
(55290, 854, 'Актёр'),
(10020, 1048334, 'Актёр'),
(718, 1048334, 'Актёр'),
(3394604, 1048334, 'Актёр'),
(40447, 1048334, 'Актёр')
ON CONFLICT (staff_id, movie_id) DO NOTHING;

-- Inserting data into the movie_genre table
INSERT INTO movie_genre (movie_id, genre_id)
VALUES
(920265, (SELECT id FROM genre WHERE name = 'мультфильм')),
(920265, (SELECT id FROM genre WHERE name = 'фантастика')),
(920265, (SELECT id FROM genre WHERE name = 'фэнтези')),
(920265, (SELECT id FROM genre WHERE name = 'боевик')),
(920265, (SELECT id FROM genre WHERE name = 'комедия')),
(920265, (SELECT id FROM genre WHERE name = 'приключения')),
(920265, (SELECT id FROM genre WHERE name = 'семейный')),
(841081, (SELECT id FROM genre WHERE name = 'мюзикл')),
(841081, (SELECT id FROM genre WHERE name = 'драма')),
(841081, (SELECT id FROM genre WHERE name = 'мелодрама')),
(841081, (SELECT id FROM genre WHERE name = 'комедия')),
(841081, (SELECT id FROM genre WHERE name = 'музыка')),
(854, (SELECT id FROM genre WHERE name = 'боевик')),
(854, (SELECT id FROM genre WHERE name = 'комедия')),
(854, (SELECT id FROM genre WHERE name = 'криминал')),
(1048334, (SELECT id FROM genre WHERE name = 'драма')),
(1048334, (SELECT id FROM genre WHERE name = 'криминал')),
(1048334, (SELECT id FROM genre WHERE name = 'триллер'))
ON CONFLICT (movie_id, genre_id) DO NOTHING;

COMMIT;

-- add collections
INSERT INTO collection_movie (collection_id, movie_id) VALUES
(1, 361),
(1, 400787),
(1, 448),
(1, 325),
(1, 258687),
(1, 342),
(1, 326),
(1, 111543),


(2, 1108577),
(2, 725190),
(2, 6462),
(2, 635772),
(2, 447301),
(2, 476),
(2, 474),
(2, 2213),

(3, 920265),
(3, 841081),
(3, 854),
(3, 1048334);



-- create incoming releses
INSERT INTO movie (id, name, original_name, release_year, poster, duration, rating) OVERRIDING SYSTEM VALUE VALUES
(
    25,
    'Легенда об Очи',
    'The Legend of Ochi',
    '2025-05-16T00:00:00.000000Z',
    'https://www.kino-teatr.ru/movie/poster/185445/232809.jpg',
    '1ч 35м',
    NULL
);

-- movie id = 25
UPDATE movie set (name, original_name, about, poster, promo_poster, release_year, slogan,
director, country, premier_russia, premier_global, duration, 
rating_kp, rating_imdb, short_description, logo, backdrop) =
(
    'Легенда об Очи',
    'The Legend of Ochi',
    'В отдаленной деревне на острове Карпатия застенчивую девочку воспитывают в страхе перед неуловимым видом животных, известным как очи. Но когда она обнаруживает, что раненый детеныш Очи остался дома, она убегает, чтобы вернуть его домой.',
    'https://www.kino-teatr.ru/movie/poster/185445/232809.jpg',
    NULL,
    '2025-05-08T00:00:00.000000Z',
    'Там есть что-то еще.',
    'Исайя Саксон',
    'США, Финляндия, Великобритания',
    '2025-05-08T00:00:00.000000Z',
    '2025-01-26T00:00:00.000000Z',
    '1ч 35мин',
    6.1,
    6.3,
    'Девочка спасает маленькое лесное чудище и меняет мир',
    '/static/img/legend_of_ochi_name.webp',
    'https://platform.polygon.com/wp-content/uploads/sites/2/2025/01/https___cdn.sanity.io_images_xq1bjtf4_production_68467d7a0d4c6fa2936f45ec7c0405573bd00daf-2000x1125-1.jpg?quality=90&strip=all&crop=7.8125%2C0%2C84.375%2C100&w=2400'
)
where name = 'Легенда об Очи';

-- genres приключения семейный фэнтези
-- Inserting data into the movie_genre table
INSERT INTO movie_genre (movie_id, genre_id)
VALUES
((SELECT id FROM movie WHERE name = 'Легенда об Очи'), (SELECT id FROM genre WHERE name = 'фэнтези')),
((SELECT id FROM movie WHERE name = 'Легенда об Очи'), (SELECT id FROM genre WHERE name = 'приключения')),
((SELECT id FROM movie WHERE name = 'Легенда об Очи'), (SELECT id FROM genre WHERE name = 'семейный'))
ON CONFLICT (movie_id, genre_id) DO NOTHING;


------------
-- test user
INSERT into "user" (login, hashed_password)
VALUES ('KinoLooker', '123456');

-- add reviews to test user
INSERT INTO review (user_id, movie_id, review_text, score) VALUES
(1, 361, 'Отличный фильм!', 9.0),
(1, 400787, 'Отличный фильм!', 8.0),
(1, 448, 'Отличный фильм!', 10.0);

-- add favorites to test user
insert into user_person_favorite (person_id, user_id) values (1089330, 1), (2378404, 1);
insert into user_movie_favorite  (movie_id, user_id) values (361, 1), (400787, 1), (448, 1);

