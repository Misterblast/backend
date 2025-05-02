package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	setEntity "github.com/ghulammuzz/misterblast/internal/set/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/redis/go-redis/v9"
)

type SetRepository interface {
	Add(class setEntity.SetSet) error
	Delete(id int32) error
	List(ctx context.Context, filter map[string]string) ([]setEntity.ListSet, error)
}

type setRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewSetRepository(db *sql.DB, redis *redis.Client) SetRepository {
	return &setRepository{db, redis}
}

func (c *setRepository) Add(class setEntity.SetSet) error {

	query := `INSERT INTO sets (name, lesson_id, class_id, is_quiz) VALUES ($1, $2, $3, $4)`
	_, err := c.db.Exec(query, class.Name, class.LessonID, class.ClassID, class.IsQuiz)
	if err != nil {
		log.Error("[Repo][AddSet] Error Exec: ", err)
		return app.NewAppError(500, "failed to insert class")
	}

	return nil
}

func (c *setRepository) Delete(id int32) error {
	query := `DELETE FROM sets WHERE id = $1`
	result, err := c.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteSet] Error Exec: ", err)
		return app.NewAppError(500, "failed to delete class")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("[Repo][DeleteSet] Error RowsAffected: ", err)
		return app.NewAppError(500, "failed to check rows affected")
	}
	if rowsAffected == 0 {
		return app.ErrNotFound
	}

	return nil
}

func (r *setRepository) List(ctx context.Context, filter map[string]string) ([]setEntity.ListSet, error) {
	redisKey := "set:list"
	for k, v := range filter {
		redisKey += fmt.Sprintf(":%s=%s", k, v)
	}
	if r.redis != nil {
		cachedVal, err := cache.Get(ctx, redisKey, r.redis)
		if err == nil && cachedVal != "" {
			var cachedSets []setEntity.ListSet
			if err := json.Unmarshal([]byte(cachedVal), &cachedSets); err == nil {
				return cachedSets, nil
			}
		}
	}

	query := `SELECT s.id, s.name, l.name AS lesson, c.name AS class, s.is_quiz FROM sets s
		JOIN lessons l ON s.lesson_id = l.id
		JOIN classes c ON s.class_id = c.id WHERE 1=1`
	args := []interface{}{}
	argCounter := 1

	if lesson, ok := filter["lesson"]; ok {
		query += fmt.Sprintf(" AND l.name = $%d", argCounter)
		args = append(args, lesson)
		argCounter++
	}
	if class, ok := filter["class"]; ok {
		query += fmt.Sprintf(" AND c.name = $%d", argCounter)
		args = append(args, class)
		argCounter++
	}
	if isQuiz, ok := filter["is_quiz"]; ok {
		query += fmt.Sprintf(" AND s.is_quiz = $%d", argCounter)
		args = append(args, isQuiz == "true")
		argCounter++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("[Repo][ListSets] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch sets")
	}
	defer rows.Close()

	var sets []setEntity.ListSet
	for rows.Next() {
		var set setEntity.ListSet
		if err := rows.Scan(&set.ID, &set.Name, &set.Lesson, &set.Class, &set.IsQuiz); err != nil {
			log.Error("[Repo][ListSets] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan set")
		}
		sets = append(sets, set)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListSets] Error Iterating Rows: ", err)
		return nil, app.NewAppError(500, "error iterating rows")
	}

	if r.redis != nil {
		serialized, err := json.Marshal(sets)
		if err == nil {
			_ = cache.Set(ctx, redisKey, string(serialized), r.redis, cache.ExpBlazing)
		}
	}

	return sets, nil
}
