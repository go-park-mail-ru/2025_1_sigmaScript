import json
import psycopg2
from psycopg2 import sql
from psycopg2.extras import execute_values
from datetime import datetime

# --- НАСТРОЙКИ ПОДКЛЮЧЕНИЯ К БД ---
DB_NAME = "filmlk"
DB_USER = "filmlk_user"
DB_PASSWORD = "filmlk_password"
DB_HOST = "localhost" # или адрес вашего сервера PostgreSQL
DB_PORT = "5433"      # стандартный порт PostgreSQL

# --- Путь к файлу с данными ---
SEED_FILE_PATH = "seed.txt"


def get_db_connection():
    """Устанавливает соединение с базой данных PostgreSQL."""
    try:
        conn = psycopg2.connect(
            dbname=DB_NAME,
            user=DB_USER,
            password=DB_PASSWORD,
            host=DB_HOST,
            port=DB_PORT
        )
        return conn
    except psycopg2.Error as e:
        print(f"Ошибка подключения к PostgreSQL: {e}")
        return None

def format_date(date_str):
    """Конвертирует строку с датой в объект date или None."""
    if not date_str:
        return None
    try:
        return datetime.strptime(date_str, '%Y-%m-%d').date()
    except ValueError:
        try:
            year = int(date_str)
            return datetime(year, 1, 1).date()
        except ValueError:
            print(f"Некорректный формат даты: {date_str}")
            return None

def map_sex(sex_str):
    """Маппинг строкового представления пола в значение ENUM sex_type."""
    if sex_str == "Мужской":
        return "Мужчина"
    elif sex_str == "Женский":
        return "Женщина"
    return ""

def map_career(career_str):
    """Маппинг карьеры к значениям ENUM career_type."""
    career_str_lower = career_str.lower() if career_str else ""
    if career_str_lower in ["актер", "актёр", "актриса"]:
        return "Актёр"
    elif career_str_lower == "продюсер":
        return "Продюсер"
    elif career_str_lower == "режиссёр": # Используем ё
        return "Режиссёр"
    elif career_str_lower == "сценарист":
        return "Сценарист"
    # Если карьера не распознана, можно вернуть None или значение по умолчанию
    # print(f"Неизвестная карьера: {career_str}, будет использовано null или default")
    return "Актёр" # По умолчанию, если не определено, или None если это предпочтительнее

def load_data(conn, data):
    """Загружает данные в таблицы PostgreSQL."""
    cursor = conn.cursor()

    inserted_genres_map = {}
    inserted_careers_map = {}

    print("Заполнение таблицы 'career'...")
    defined_careers = ['Актёр', 'Продюсер', 'Режиссёр', 'Сценарист']
    for career_name_enum in defined_careers:
        try:
            cursor.execute("""
                INSERT INTO career (career) VALUES (%s)
                ON CONFLICT (career) DO NOTHING
                RETURNING id;
            """, (career_name_enum,))
            res = cursor.fetchone()
            if res:
                inserted_careers_map[career_name_enum] = res[0]
            else:
                cursor.execute("SELECT id FROM career WHERE career = %s;", (career_name_enum,))
                existing_career_id = cursor.fetchone()
                if existing_career_id:
                    inserted_careers_map[career_name_enum] = existing_career_id[0]
        except psycopg2.Error as e:
            print(f"Ошибка при добавлении/получении карьеры '{career_name_enum}': {e}")
            conn.rollback()
            return
    print(f"ID карьер: {inserted_careers_map}")

    for movie_data in data:
        json_movie_id = movie_data.get('id')
        if json_movie_id is None:
            print(f"Ошибка: Фильм '{movie_data.get('name')}' не имеет 'id' в JSON. Пропуск фильма.")
            continue
        print(f"\nОбработка фильма: {movie_data.get('name')} (JSON ID: {json_movie_id})")

        movie_genres_ids = []
        if movie_data.get('genres'):
            genre_names = [genre.strip() for genre in movie_data['genres'].split(',')]
            for genre_name in genre_names:
                if not genre_name: continue
                if genre_name not in inserted_genres_map:
                    try:
                        cursor.execute("""
                            INSERT INTO genre (name) VALUES (%s)
                            ON CONFLICT (name) DO NOTHING
                            RETURNING id;
                        """, (genre_name,))
                        genre_id_tuple = cursor.fetchone()
                        if genre_id_tuple:
                            inserted_genres_map[genre_name] = genre_id_tuple[0]
                        else:
                            cursor.execute("SELECT id FROM genre WHERE name = %s;", (genre_name,))
                            existing_genre_id = cursor.fetchone()
                            if existing_genre_id:
                                inserted_genres_map[genre_name] = existing_genre_id[0]
                    except psycopg2.Error as e:
                        print(f"Ошибка при добавлении жанра '{genre_name}': {e}")
                        conn.rollback()
                        continue
                if genre_name in inserted_genres_map:
                     movie_genres_ids.append(inserted_genres_map[genre_name])

        staff_db_ids_with_roles = []
        for staff_member in movie_data.get('staff', []):
            json_person_id = staff_member.get('id')
            if json_person_id is None:
                print(f"Ошибка: Персона '{staff_member.get('fullName')}' не имеет 'id' в JSON. Пропуск персоны.")
                continue

            person_full_name = staff_member.get('fullName')
            person_en_full_name = staff_member.get('enFullName')
            person_photo = staff_member.get('photo', '/static/avatars/avatar_default_picture.svg')
            person_about = staff_member.get('about') if staff_member.get('about') else 'Информация по этому человеку не указана'
            person_sex = map_sex(staff_member.get('sex'))
            person_growth = staff_member.get('growth')
            person_birthday = format_date(staff_member.get('birthday'))
            person_death = format_date(staff_member.get('death'))
            
            try:
                cursor.execute("""
                    INSERT INTO person (id, full_name, en_full_name, photo, about, sex, growth, birthday, death, updated_at)
                    OVERRIDING SYSTEM VALUE VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, CURRENT_TIMESTAMP)
                    ON CONFLICT (id) DO UPDATE SET
                        full_name = EXCLUDED.full_name,
                        en_full_name = EXCLUDED.en_full_name,
                        photo = EXCLUDED.photo,
                        about = EXCLUDED.about,
                        sex = EXCLUDED.sex,
                        growth = EXCLUDED.growth,
                        birthday = EXCLUDED.birthday,
                        death = EXCLUDED.death,
                        updated_at = CURRENT_TIMESTAMP
                    RETURNING id;
                """, (
                    json_person_id, person_full_name, person_en_full_name, person_photo,
                    person_about, person_sex, person_growth, person_birthday, person_death
                ))
                db_person_id = cursor.fetchone()[0] # ID будет равен json_person_id
                print(f"Добавлена/Обновлена персона '{person_full_name}' с ID: {db_person_id}")

                person_career_str = staff_member.get('career')
                if person_career_str:
                    person_career_enum = map_career(person_career_str)
                    if person_career_enum in inserted_careers_map:
                        career_db_id = inserted_careers_map[person_career_enum]
                        cursor.execute("""
                            INSERT INTO career_person (person_id, career_id)
                            VALUES (%s, %s)
                            ON CONFLICT (person_id, career_id) DO NOTHING;
                        """, (db_person_id, career_db_id))
                    else:
                        print(f"Карьера '{person_career_enum}' не найдена в карте ID для персоны {person_full_name}")
                
                role_enum_value = map_career(staff_member.get('career', 'Актёр'))
                staff_db_ids_with_roles.append((db_person_id, role_enum_value))

            except psycopg2.Error as e:
                print(f"Ошибка при добавлении/обновлении персоны '{person_full_name}' (JSON ID: {json_person_id}): {e}")
                conn.rollback()
                continue
        
        db_movie_id = None
        try:
            movie_name = movie_data.get('name')
            movie_original_name = movie_data.get('originalName')
            movie_about = movie_data.get('about', 'Информация по этому фильму не указана')
            movie_poster = movie_data.get('poster', '/static/movies/poster_default_picture.webp')
            movie_logo = movie_data.get('logo', '#')
            movie_backdrop = movie_data.get('backdrop', '#')
            release_year_val = movie_data.get('releaseYear')
            release_date_val = None
            if release_year_val:
                try: release_date_val = datetime(int(release_year_val), 1, 1).date()
                except ValueError: print(f"Некорректный год релиза: {release_year_val}")
            
            movie_country = movie_data.get('country')
            movie_slogan = movie_data.get('slogan')
            movie_director = movie_data.get('director')
            movie_budget = movie_data.get('budget')
            movie_box_office_us = movie_data.get('boxOfficeUS')
            movie_box_office_global = movie_data.get('boxOfficeGlobal')
            movie_box_office_russia = movie_data.get('boxOfficeRussia')
            movie_premier_russia = format_date(movie_data.get('premierRussia'))
            movie_premier_global = format_date(movie_data.get('premierGlobal'))
            movie_rating_val = movie_data.get('rating', 5.00) # default rating
            movie_rating_kp = movie_data.get('ratingKp')
            movie_rating_imdb = movie_data.get('ratingImdb')
            movie_duration = movie_data.get('duration')
            movie_short_description = movie_data.get('shortDescription')
            movie_trailer_url = movie_data.get('trailerUrl')
            movie_watchability_bool = True if movie_data.get('watchability') else False

            cursor.execute("""
                INSERT INTO movie (
                    id, name, original_name, about, poster, logo, backdrop, release_year, country, slogan,
                    director, budget, box_office_us, box_office_global, box_office_russia,
                    premier_russia, premier_global, rating, rating_kp, rating_imdb, duration,
                    short_description, trailerUrl, watchability, updated_at
                )
                OVERRIDING SYSTEM VALUE VALUES (
                    %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, CURRENT_TIMESTAMP
                )
                ON CONFLICT (id) DO UPDATE SET
                    name = EXCLUDED.name, original_name = EXCLUDED.original_name, about = EXCLUDED.about,
                    poster = EXCLUDED.poster, logo = EXCLUDED.logo, backdrop = EXCLUDED.backdrop,
                    release_year = EXCLUDED.release_year, country = EXCLUDED.country, slogan = EXCLUDED.slogan,
                    director = EXCLUDED.director, budget = EXCLUDED.budget, box_office_us = EXCLUDED.box_office_us,
                    box_office_global = EXCLUDED.box_office_global, box_office_russia = EXCLUDED.box_office_russia,
                    premier_russia = EXCLUDED.premier_russia, premier_global = EXCLUDED.premier_global,
                    rating = EXCLUDED.rating, rating_kp = EXCLUDED.rating_kp, rating_imdb = EXCLUDED.rating_imdb,
                    duration = EXCLUDED.duration, short_description = EXCLUDED.short_description,
                    trailerUrl = EXCLUDED.trailerUrl, watchability = EXCLUDED.watchability, updated_at = CURRENT_TIMESTAMP
                RETURNING id;
            """, (
                json_movie_id, movie_name, movie_original_name, movie_about, movie_poster, movie_logo, movie_backdrop,
                release_date_val, movie_country, movie_slogan, movie_director, movie_budget, movie_box_office_us,
                movie_box_office_global, movie_box_office_russia, movie_premier_russia, movie_premier_global,
                movie_rating_val, movie_rating_kp, movie_rating_imdb, movie_duration, movie_short_description,
                movie_trailer_url, movie_watchability_bool
            ))
            db_movie_id = cursor.fetchone()[0] # ID будет равен json_movie_id
            print(f"Добавлен/Обновлен фильм '{movie_name}' с DB ID: {db_movie_id}")
        except psycopg2.Error as e:
            print(f"Ошибка при добавлении/обновлении фильма '{movie_data.get('name')}' (JSON ID: {json_movie_id}): {e}")
            conn.rollback()
            continue

        if db_movie_id:
            if movie_genres_ids:
                try:
                    execute_values(cursor, """
                        INSERT INTO movie_genre (movie_id, genre_id) VALUES %s
                        ON CONFLICT (movie_id, genre_id) DO NOTHING;
                    """, [(db_movie_id, genre_id) for genre_id in movie_genres_ids])
                except psycopg2.Error as e:
                    print(f"Ошибка при связывании жанров с фильмом ID {db_movie_id}: {e}")
                    conn.rollback()

            if staff_db_ids_with_roles:
                staff_movie_link_data = [(person_id, db_movie_id, role) for person_id, role in staff_db_ids_with_roles]
                try:
                    execute_values(cursor, """
                        INSERT INTO movie_staff (staff_id, movie_id, role) VALUES %s
                        ON CONFLICT (staff_id, movie_id) DO UPDATE SET role = EXCLUDED.role;
                    """, staff_movie_link_data)
                except psycopg2.Error as e:
                    print(f"Ошибка при связывании состава с фильмом ID {db_movie_id}: {e}")
                    conn.rollback()

            if movie_data.get('watchability'):
                providers_data = []
                for provider in movie_data['watchability']:
                    providers_data.append((
                        db_movie_id, provider.get('name'),
                        provider.get('logoUrl', '/static/svg/play.svg'),
                        provider.get('movieUrl', '#')
                    ))
                if providers_data:
                    try:
                        # Для watch_provider нет явного ON CONFLICT в схеме, зависит от бизнес-логики
                        # Если нужно обновлять, добавьте ON CONFLICT (movie_id, name) DO UPDATE SET ...
                        # Здесь просто вставляем, возможны дубликаты если провайдеры для одного фильма повторяются в JSON
                        # или при повторном запуске. Удалим старые перед вставкой новых для этого фильма.
                        cursor.execute("DELETE FROM watch_provider WHERE movie_id = %s;", (db_movie_id,))
                        execute_values(cursor, """
                            INSERT INTO watch_provider (movie_id, name, logo, watch_url) VALUES %s
                        """, providers_data)
                        print(f"  Добавлены провайдеры для фильма DB ID {db_movie_id}")
                    except psycopg2.Error as e:
                        print(f"Ошибка при добавлении провайдеров для фильма {db_movie_id}: {e}")
                        conn.rollback()
        
        # Reviews (без изменений, так как user_id не из JSON фильма)
        if movie_data.get('reviews'):
            # ... (логика для отзывов остается прежней)
            pass


    conn.commit()
    cursor.close()
    print("\nЗагрузка данных успешно завершена.")


def main():
    """Главная функция для запуска процесса."""
    conn = get_db_connection()
    if not conn:
        return

    try:
        with open(SEED_FILE_PATH, 'r', encoding='utf-8') as f:
            # Ожидаем, что файл содержит JSON-массив объектов
            try:
                data = json.load(f)
            except json.JSONDecodeError as e:
                print(f"Ошибка декодирования JSON из файла {SEED_FILE_PATH}: {e}")
                print("Убедитесь, что файл содержит валидный JSON-массив (объекты в '[...]' и разделены запятыми).")
                # Попытка исправить распространенную ошибку (отсутствие обрамляющих скобок)
                f.seek(0)
                content = f.read()
                try:
                    data = json.loads(f"[{content}]") # Обертываем в скобки
                    print("Данные успешно прочитаны после добавления квадратных скобок.")
                except json.JSONDecodeError as e_retry:
                    print(f"Повторная попытка декодирования JSON также не удалась: {e_retry}")
                    return


        if isinstance(data, list): # Убедимся, что это список фильмов
            load_data(conn, data)





        else:
            print("Ошибка: данные в файле не являются JSON-массивом.")

    except FileNotFoundError:
        print(f"Файл {SEED_FILE_PATH} не найден.")
    except Exception as e:
        print(f"Произошла непредвиденная ошибка: {e}")
        if conn: # Если ошибка произошла после подключения, откатываем транзакцию
             conn.rollback()
    finally:
        if conn:
            conn.close()
            print("Соединение с PostgreSQL закрыто.")

if __name__ == "__main__":
    main()
    
    

