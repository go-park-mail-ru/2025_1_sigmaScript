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
    backdrop TEXT DEFAULT '#'
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

INSERT INTO movie (name, poster) VALUES
('Бойцовский клуб', '/img/0.webp'),
('Матрица', '/img/7.webp'),
('Форрест Гамп', '/img/2.webp'),
('Крестный отец', '/img/3.webp'),
('Интерстеллар', '/img/4.webp'),
('Криминальное чтиво', '/img/5.webp'),
('Побег из Шоушенка', '/img/6.webp'),
('Тёмный рыцарь', '/img/1.webp'),
('Зелёная миля', '/img/8.webp'),
('Одержимость', '/img/9.webp'),
('Оппенгеймер', '/img/11.webp'),
('Звёздные войны: Эпизод 4 – Новая надежда', '/img/12.webp'),
('Рокки', '/img/13.webp'),
('Джокер', '/img/14.webp'),
('Игра в имитацию', '/img/15.webp'),
('Начало', '/img/16.webp'),
('Назад в будущее', '/img/17.webp'),
('Гладиатор', '/img/18.webp'),
('Титаник', '/img/19.webp'),
('Ford против Ferrari', '/img/10.webp'),
('Пророк. История Александра Пушкина', '/img/21.webp'),
('Батя', '/img/23.webp'),
('Финист. первый Богатырь', '/img/22.webp'),
('Матрица 2: Перезагрузка', '/img/7.webp');

-- movi id = 25
INSERT INTO movie (name, original_name, release_year, poster, duration, rating) VALUES
(
    'Легенда об Очи',
    'The Legend of Ochi',
    '2025-05-16T00:00:00.000000Z',
    'https://www.kino-teatr.ru/movie/poster/185445/232809.jpg',
    '1ч 35м',
    NULL
);

UPDATE movie SET promo_poster = '/img/promo_prorok.webp' WHERE name = 'Пророк. История Александра Пушкина'
AND poster = '/img/21.webp';
UPDATE movie SET promo_poster = '/img/promo_batya.webp' WHERE name = 'Батя'
AND poster = '/img/23.webp';
UPDATE movie SET promo_poster = '/img/promo_bogatyr.webp' WHERE name = 'Финист. первый Богатырь'
AND poster = '/img/22.webp';


-- movie id = 26
INSERT INTO movie (
    name, original_name, about, short_description, poster, release_year, slogan,
    budget, box_office_us, box_office_global, box_office_russia, premier_russia, premier_global,
    rating, rating_kp, rating_imdb, duration, logo, backdrop 
) VALUES
(
    'Человек-паук: Через вселенные',
    'Spider-Man: Into the Spider-Verse',
    'Мы всё знаем о Питере Паркере. Он спас город, влюбился, а потом спасал город снова и снова… Но все это – в нашем измерении. А что если в результате работы гигантского коллайдера откроется окно из одного измерения в другое? Найдется ли в нем свой Человек-паук? И как он будет выглядеть? Приготовьтесь к тому, что в разных вселенных могут быть разные Люди-пауки и однажды им придется собраться вместе...',
    'Пауки из разных измерений объединяются перед общей угрозой. Изобретательный кинокомикс с «Оскаром» за анимацию',
    'https://image.openmoviedb.com/kinopoisk-images/1900788/64417bd3-838b-4910-a9f4-c278c509d568/x1000',
    '2018-12-12T00:00:00.000000Z',
    'Более одного носит маску',
    90000000,
    190241310,
    375582637,
    6298173,
    '2018-12-13 03:00:00',
    '2018-12-12 03:00:00',
    9.1,
    8.18,
    8.4,
    '1ч 57м',
    'https://image.openmoviedb.com/tmdb-images/w500/bJ3VpPP3VkJM9H8GfRK8wSvTwPy.png',
    'https://avatars.mds.yandex.net/get-ott/1531675/2a00000178cb2efc9762720b53a3b57633ef/2016x1134'
);

INSERT INTO collection (name, is_main_collection) VALUES
('Лучшие за всё время', TRUE),
('Номинанты на оскар', TRUE),
('promo', TRUE);

INSERT INTO collection_movie (collection_id, movie_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(1, 6),
(1, 7),
(1, 8),
(1, 9),
(1, 10),

(2, 11),
(2, 12),
(2, 13),
(2, 14),
(2, 15),
(2, 16),
(2, 17),
(2, 18),
(2, 19),
(2, 20),

(3, 21),
(3, 22),
(3, 23);

UPDATE movie SET duration = '1ч 53м' WHERE id = 21;
UPDATE movie SET duration = '1ч 16м' WHERE id = 22;
UPDATE movie SET duration = '1ч 52м' WHERE id = 23;

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


-- creating movies
UPDATE movie SET (name, original_name, about, poster, release_year, Budget, box_office_us, box_office_global, rating, duration) =
(
    'Матрица',
    'Matrix',
    'Жизнь Томаса Андерсона разделена на две части: днём он — самый обычный офисный работник, получающий нагоняи от начальства, а ночью превращается в хакера по имени Нео, и нет места в сети, куда он бы не смог проникнуть. Но однажды всё меняется. Томас узнаёт ужасающую правду о реальности.',
    '/static/img/7.webp',
    '1999-1-1',
    63000000,
    171479930,
    463517383,
    8.5,
    '2ч 16м'
)
where id = 2;

INSERT INTO person (full_name, en_full_name, photo) VALUES
('Леонардо Ди Каприо', 'Leonardo DiCaprio', 'https://avatars.mds.yandex.net/get-entity_search/2310675/1130394491/S600xU_2x'),
('Морган Фримен', 'Morgan Freeman', 'https://avatars.mds.yandex.net/get-entity_search/2057552/1132084397/S600xU_2x'),
( 'Том Хэнкс', 'Tom Hanks', 'https://avatars.mds.yandex.net/get-entity_search/2005770/833182325/S600xU_2x'),
('Джонни Депп', 'Johnny Depp', '/static/avatars/avatar_default_picture.svg'),
('Том Круз', 'Tom Cruise', '/static/avatars/avatar_default_picture.svg'),
('Сэмюэл Л. Джексон', 'Samuel L. Jackson', 'https://avatars.mds.yandex.net/get-entity_search/98180/952678918/S600xU_2x'),
('Брэд Питт', 'Brad Pitt', '/static/img/brad_pitt.webp'),
('Рассел Кроу', 'Russell Crowe', 'https://avatars.mds.yandex.net/get-entity_search/478647/809836058/S600xU_2x'),
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
) where id = 11;


-- inserting movie genres
INSERT INTO movie_genre (movie_id, genre_id) VALUES
(2, 4), -- matrixix
(2, 3),
(24, 4), -- matrix 2
(24, 3),
(24, 2),
(26, 17),
(26, 3),
(26, 11),
(26, 4),
(26, 5),
(26, 10),
(26, 18);


-- inserting movie staff
INSERT INTO movie_staff (movie_id, staff_id) VALUES
(1, 7), -- brad pitt
(2, 11), -- keanu reeves
(24, 11), -- keanu reeves
(26, 12),
(26, 13),
(26, 14),
(26, 15),
(26, 16),
(26, 17);

-- inserting person genres
INSERT INTO person_genre (person_id, genre_id) VALUES
(11, 4),
(11, 3),
(11, 2);

-- creating careers for persons
insert into career_person (career_id, person_id) values (1, 11), (2, 11), (3, 11), (3, 12), (3, 13), (3, 14);


-- creating countries for movies
INSERT INTO movie_country (movie_id, country_id) VALUES
(26, 1),
(26, 2);

-- creating watch provirders for movie 
INSERT INTO watch_provider (movie_id, name, logo, watch_url) VALUES
(
    26,
    'Триколор Кино и ТВ',
    'https://avatars.mds.yandex.net/get-ott/239697/947e777c-2f73-4cbc-b09d-6bfa3966ba13/orig',
    'https://kino.tricolor.tv/watch/chelovek-pauk-cherez-vselennye-2018/'
),
(
    26,
    'Кинопоиск HD',
    'https://play-lh.googleusercontent.com/5czw6iycA8YhjI653GQdwnnmu8NNzEMXV32gZKoVCYZV6PQUAv_YV0uJ2PU1E-Jm9PE=w480-h960-rw',
    'https://hd.kinopoisk.ru/film/4d924361a6c32b09aeee7d1a63f9c3bf?content_tab=overview'
);

------------
-- test user
INSERT into "user" (login, hashed_password)
VALUES ('KinoLooker', '123456');

-- add reviews to test user
INSERT INTO review (user_id, movie_id, review_text, score) VALUES
(1, 1, 'Отличный фильм!', 9.0),
(1, 2, 'Отличный фильм!', 8.0),
(1, 3, 'Отличный фильм!', 10.0);

-- add favorites to test user
insert into user_person_favorite (person_id, user_id) values (1, 1), (2, 1);
insert into user_movie_favorite  (movie_id, user_id) values (4, 1), (3, 1), (26, 1);

