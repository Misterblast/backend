package repo

import (
	"database/sql"
	"fmt"

	classEntity "github.com/ghulammuzz/misterblast/internal/class/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type ClassRepository interface {
	Add(class classEntity.SetClass) error
	Delete(id int32) error
	List() ([]classEntity.Class, error)
	Exists(class string) (bool, error)
}
type classRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) ClassRepository {
	return &classRepository{db: db}
}

func (c *classRepository) Exists(class string) (bool, error) {
	query := `SELECT 1 FROM classes WHERE name = $1`
	var exists bool
	err := c.db.QueryRow(query, class).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Error("[Repo][ExistsClass] Error QueryRow: ", err)
		return false, app.NewAppError(500, "failed to check if class exists")
	}
	return exists, nil
}

func (c *classRepository) Add(class classEntity.SetClass) error {
	if err := class.Validate(); err != nil {
		log.Error("[Repo][AddClass] Error Validate: ", err)
		return app.NewAppError(400, "validation failed")
	}

	query := `INSERT INTO classes (name) VALUES ($1)`
	_, err := c.db.Exec(query, class.Name)
	if err != nil {
		log.Error("[Repo][AddClass] Error Exec: ", err)
		return app.NewAppError(500, "failed to insert class")
	}

	return nil
}

func (c *classRepository) Delete(id int32) error {
	query := `DELETE FROM classes WHERE id = $1`
	result, err := c.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteClass] Error Exec: ", err)
		return app.NewAppError(500, "failed to delete class")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("[Repo][DeleteClass] Error RowsAffected: ", err)
		return app.NewAppError(500, "failed to check rows affected")
	}
	if rowsAffected == 0 {
		return app.ErrNotFound
	}

	return nil
}

func (c *classRepository) List() ([]classEntity.Class, error) {
	query := `SELECT id, name FROM classes`
	rows, err := c.db.Query(query)
	if err != nil {
		log.Error("[Repo][ListClass] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch classes")
	}
	defer rows.Close()

	var classes []classEntity.Class

	for rows.Next() {
		var class classEntity.Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			log.Error("[Repo][ListClass] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan class")
		}
		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListClass] Error Iterating Rows: ", err)
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return classes, nil
}
