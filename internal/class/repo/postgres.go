package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	classEntity "github.com/ghulammuzz/misterblast/internal/class/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/redis/go-redis/v9"
)

type ClassRepository interface {
	Add(class classEntity.SetClass) error
	Delete(id int32) error
	List(ctx context.Context) ([]classEntity.Class, error)
	Exists(class string) (bool, error)
}
type classRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewClassRepository(db *sql.DB, redis *redis.Client) ClassRepository {
	return &classRepository{db, redis}
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

func (c *classRepository) List(ctx context.Context) ([]classEntity.Class, error) {
	var classes []classEntity.Class
	redisKey := "cache:classes:list"

	if c.redis != nil {
		cached, err := cache.Get(ctx, redisKey, c.redis)
		if err != nil && err != redis.Nil {
			log.Warn("[Repo][ListClass] Redis error: ", err)
		}
		if err == nil && cached != "" {
			if err := json.Unmarshal([]byte(cached), &classes); err == nil {
				return classes, nil
			} else {
				log.Warn("[Repo][ListClass] Failed to unmarshal Redis cache: ", err)
			}
		}
	}

	query := `SELECT id, name FROM classes`
	rows, err := c.db.Query(query)
	if err != nil {
		log.Error("[Repo][ListClass] Error executing query: ", err)
		return nil, app.NewAppError(500, "failed to fetch classes")
	}
	defer rows.Close()

	for rows.Next() {
		var class classEntity.Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			log.Error("[Repo][ListClass] Error scanning row: ", err)
			return nil, app.NewAppError(500, "failed to scan class")
		}
		classes = append(classes, class)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListClass] Error iterating rows: ", err)
		return nil, app.NewAppError(500, "error iterating rows")
	}

	if c.redis != nil {
		if dataJSON, err := json.Marshal(classes); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), c.redis, cache.ExpBlazing)
		} else {
			log.Warn("[Repo][ListClass] Failed to marshal data for Redis: ", err)
		}
	}

	return classes, nil
}
