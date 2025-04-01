Relation User:

{ID} -> Login, HashedPassword, Avatar, BirthDate, CreatedAt, UpdatedAt

---

Relation Collection:

{ID} -> Name, Slug, CreatedAt, UpdatedAt

---

Relation Person:

{ID} -> FullName, EnFullName, Photo, About, Sex, Growth, Birthday, Death, Age, Birthplace, Deathplace, CreatedAt, UpdatedAt

---

Relation Genre:

{GenreID} -> Name, CreatedAt, UpdatedAt

---

Relation Movie:

{ID} -> Name, About, Poster, Card, ReleaseYear, Country, Slogan, Director, Budget, BoxOfficeUS, BoxOfficeGlobal, BoxOfficeRussia, PremiereRussia, PremiereGlobal, Rating, Duration, CreatedAt, UpdatedAt

---

Relation Collection_Movie:

{CollectionID, MovieID} -> {}

---

Relation Movie_Staff:

{StaffID, MovieID} -> Role

---

Relation Movie_Genre:

{GenreID, MovieID} -> {}

---

Relation Review:

{ReviewID} -> UserID, MovieID, ReviewText, Score, CreatedAt, UpdatedAt

---

Relation Like:

{LikeID} -> UserID, ReviewID, IsValid, CreatedAt, UpdatedAt

---

Relation Dislike:

{DislikeID} -> UserID, ReviewID, IsValid, CreatedAt, UpdatedAt

