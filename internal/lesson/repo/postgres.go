package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	"github.com/ghulammuzz/misterblast/internal/lesson/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/redis/go-redis/v9"
)

type LessonRepository interface {
	Add(lesson entity.Lesson) error
	Delete(id int32) error
	List(ctx context.Context) ([]entity.Lesson, error)
	Exists(lesson string) (bool, error)
}

type lessonRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewLessonRepository(db *sql.DB, redis *redis.Client) LessonRepository {
	return &lessonRepository{db, redis}
}

func (r *lessonRepository) Exists(lesson string) (bool, error) {
	query := `SELECT 1 FROM lessons WHERE name = $1`
	var exists bool
	err := r.db.QueryRow(query, lesson).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Error("[Repo][ExistsLesson] Error: ", err)
		return false, app.NewAppError(500, "failed to check if lesson exists")
	}
	return true, nil
}

func (r *lessonRepository) Add(lesson entity.Lesson) error {
	query := `INSERT INTO lessons (name, code) VALUES ($1)`
	_, err := r.db.Exec(query, lesson.Name, lesson.Code)
	if err != nil {
		log.Error("[Repo][AddLesson] Error: ", err)
		return app.NewAppError(500, "failed to add lesson")
	}
	return nil
}

func (r *lessonRepository) Delete(id int32) error {
	query := `DELETE FROM lessons WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteLesson] Error: ", err)
		return app.NewAppError(500, "failed to delete lesson")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("[Repo][DeleteLesson] Error RowsAffected: ", err)
		return app.NewAppError(500, "failed to check rows affected")
	}
	if rowsAffected == 0 {
		return app.ErrNotFound
	}

	return nil
}

func (r *lessonRepository) List(ctx context.Context) ([]entity.Lesson, error) {
	var lessons []entity.Lesson

	redisKey := "lessons"

	if r.redis != nil {
		cached, err := cache.Get(ctx, redisKey, r.redis)
		if err == nil && cached != "" {
			if err := json.Unmarshal([]byte(cached), &lessons); err == nil {
				return lessons, nil
			}
		}
	}

	query := `SELECT id, name, code FROM lessons order by id`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Error("[Repo][ListLessons] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch lessons")
	}
	defer rows.Close()

	for rows.Next() {
		var lesson entity.Lesson
		if err := rows.Scan(&lesson.ID, &lesson.Name, &lesson.Code); err != nil {
			log.Error("[Repo][ListLessons] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan lesson")
		}
		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListLessons] Error Iterating Rows: ", err)
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	if r.redis != nil {
		if dataJSON, err := json.Marshal(lessons); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis, cache.ExpBlazing)
		}
	}

	return lessons, nil
}
