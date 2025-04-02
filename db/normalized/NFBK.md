Relation user:

{id} -> login, hashed_password, avatar, birth_date, created_at, updated_at

---

Relation collection:

{id} -> name, slug, created_at, updated_at

---

Relation person:

{id} -> full_name, en_full_name, photo, about, sex, growth, birthday, death, age, birth_place, death_place, created_at, updated_at

---

Relation genre:

{id} -> name, created_at, updated_at

---

Relation country:

{id} -> name, code, flag

---

Relation movie:

{id} -> name, about, poster, card, release_year, country, slogan, director, budget, box_office_us, box_office_global, boxo_office_russia, premiere_russia, premiere_global, rating, duration, created_at, updated_at

---

Relation collection_movie:

{collection_id, movie_id} -> {}

---

Relation movie_staff:

{staff_id, movie_id} -> role

---

Relation movie_genre:

{genre_id, movie_id} -> {}

---

Relation review:

{id} -> user_id, movie_id, review_text, score, created_at, updated_at

---

Relation like:

{id} -> user_id, review_id, is_valid, created_at, updated_at

---

Relation dislike:

{id} -> user_id, review_id, is_valid, created_at, updated_at
