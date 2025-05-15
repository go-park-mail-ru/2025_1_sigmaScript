# ДЗ 2 Администрирование СУБД

Все данные по ДЗ2 дисциплины СУБД находятся в директории `internal/db/postgresql_filmlk`

## 1.	Безопасность сервера СУБД

### Работа с БД через сервисную учетную запись
<p>Был настроен скрипт создания сервисную учетную запись service_user, который создает пользователя
с привилегиями для работы с публичной схемой и модификации таблиц БД. 
</p>

<p>Скрипт создания находится в директории:</p>

```
internal/db/postgresql_filmlk/scripts/service_user.sql
```

### Защита от SQL Injections
  <p>Защита от sql инъекций производится в коде приложения:
 
 1) валидация входных данных производится в слое delivery. Пример:

```go
// ...

if err = jsonutil.ReadJSON(r, &newReviewDataJSON); err != nil {
  logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
  jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(), errs.ErrBadPayload)
  return
}

if newReviewDataJSON.Score < 1 || newReviewDataJSON.Score > 10 {
  logger.Error().Err(errors.New(errs.ErrBadPayload)).Msg(fmt.Sprintf("bad score of new review: %f", newReviewDataJSON.Score))
  jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload,
    errs.ErrBadPayload)
  return
}

// ...
```

 2) Экранирование спецсимволов производится с помощью функции ValidateInputTextData
 из пакета escapingutil в слое delivery. Пример:

```go
package escapingutil

import (
	"errors"
	"html"
	"strings"
)

const (
	DEFAULT_TEXT_MAX_LENGTH = int64(500)
)

var (
	ErrorMaxLength            = errors.New("text exeeds the allowed maximum length")
	ErrorEmptyOrInvalidString = errors.New("text string is null or contains only invalid symbols")
)

func ValidateInputTextData(textData string, textMaxLength ...int64) (string, error) {
	maxLength := DEFAULT_TEXT_MAX_LENGTH
	if len(textMaxLength) > 0 && textMaxLength[0] > 0 {
		maxLength = textMaxLength[0]
	}
	if int64(len(textData)) > maxLength {
		return "", ErrorMaxLength
	}
	trimmedDataString := html.EscapeString(strings.TrimSpace(textData))

	if len(trimmedDataString) < 1 {
		return "", ErrorEmptyOrInvalidString
	}
	return trimmedDataString, nil
}
```

3) Использование Prepared Statement производится в слое repository. Пример:

```go
//...

func (r *MoviePostgresRepository) CreateNewMovieReviewInRepo(
	ctx context.Context,
	userID string,
	movieID string,
	newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error) {
	logger := log.Ctx(ctx)

	var newReviewID sql.NullInt64
	execRow, err := r.pgdb.Prepare(insertNewReviewQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo").Error())
		return nil, errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo")
	}
	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	errExec := execRow.QueryRow(
		userID,
		movieID,
		newReview.ReviewText,
		newReview.Score,
	).Scan(&newReviewID)
	if errExec != nil {
		errPg := fmt.Errorf("postgres: error while creating review - %w", errExec)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		sqlErr, ok := errExec.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return nil, errors.New(errs.ErrAlreadyExists)
		}
		return nil, errors.New(errs.ErrSomethingWentWrong)
	}

	execRowRating, err := r.pgdb.Prepare(updateMovieRatingQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo").Error())
		return nil, errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo")
	}
	defer func() {
		if clErr := execRowRating.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRowRating.Exec(movieID)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while updating movie rating - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return nil, errors.New(errs.ErrAlreadyExists)
		}
		return nil, errors.New(errs.ErrSomethingWentWrong)
	}

	if !newReviewID.Valid {
		errPg := fmt.Errorf("postgres: error while updating movie rating - %s", "got not valid review id")
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		return nil, errPg
	}

	logger.Info().Msgf("successfully updated movie rating by movie id: %s", movieID)
	return &mocks.NewReviewDataJSON{ID: int(newReviewID.Int64), Score: newReview.Score}, nil
}

// ...
```
</p>

## 2.	Настройка параметров сервера и клиента
  <p>Был создан файл с конфигурацией постгрес postgresql.conf . Он находится в директории:

```
  internal/db/postgresql_filmlk/postgresql.conf
```
  Listen_adressess '*', потому что в запуск производится в docker.
  Конфиг удовлетворяет отзывчивости системы.
  </p>


## 3.	Таймауты
  <p>
  statement_timeout и lock_timeout прописаны в файле конфигурации. Выбраны оптимальные значения с точки зрения колическва запросов и времени их выполнения.
  </p>


## 4.	pg_stat_statements
  <p> Прописаны в internal/db/postgresql_filmlk/postgresql.conf:

```
# - PG_STAT_STATEMENTS -
pg_stat_statements.max = 10000
pg_stat_statements.track = all
pg_stat_statements.save = on
```
  </p>


## 5.	auto_explain
  <p> Прописаны в internal/db/postgresql_filmlk/postgresql.conf: 

```
# - AUTO_EXPLAIN -
auto_explain.log_min_duration = 150
auto_explain.log_analyze = true
auto_explain.log_buffers = true
auto_explain.log_verbose = true
auto_explain.log_timing = true
```
  </p>


## 6.	Логгирование медленных запросов в формате, который можно распарсить PGBadger’ом
  <p>Посмотреть логи с помощью pgBadger можно, если зайти в bash контейнера

```bash
docker exec -it filmlk_db pgbadger /var/lib/postgresql/data/log/postgresql-*.log -o output.html 
docker cp filmlk_db:/output.html ./output.html
```

  Условия логгирования медленных запросов прописаны в internal/db/postgresql_filmlk/postgresql.conf.
  Например, выражение вида:

```sql
select id, login from "user" where id = 1 OR 1=1 UNION ALL select 1, (CASE WHEN (1=1) THEN PG_SLEEP(100) ELSE NULL END)::TEXT ;
```

  выдаст ошибку по времени выполнения, и запрос отправится в лог (выполнение запроса больше 150 мСек). Также в лог отправляются все ошибки.
  </p>

<p>PS: Данная конструкция запроса приведена только как пример и была иполнена вручную из терминала,
  подобные sql инъекции сервис не пропускает.</p>

