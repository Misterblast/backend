package repo

import (
	"database/sql"
	"fmt"

	"github.com/ghulammuzz/misterblast/internal/lesson/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type LessonRepository interface {
	Add(lesson entity.Lesson) error
	Delete(id int32) error
	List() ([]entity.Lesson, error)
	Exists(lesson string) (bool, error)
}

type lessonRepository struct {
	db *sql.DB
}

func NewLessonRepository(db *sql.DB) LessonRepository {
	return &lessonRepository{db: db}
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
	query := `INSERT INTO lessons (name) VALUES ($1)`
	_, err := r.db.Exec(query, lesson.Name)
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

func (r *lessonRepository) List() ([]entity.Lesson, error) {
	query := `SELECT id, name FROM lessons order by id`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Error("[Repo][ListLessons] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch lessons")
	}
	defer rows.Close()

	var lessons []entity.Lesson

	for rows.Next() {
		var lesson entity.Lesson
		if err := rows.Scan(&lesson.ID, &lesson.Name); err != nil {
			log.Error("[Repo][ListLessons] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan lesson")
		}
		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListLessons] Error Iterating Rows: ", err)
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return lessons, nil
}
