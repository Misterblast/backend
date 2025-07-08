// repo/author_repository.go
package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	"github.com/ghulammuzz/misterblast/internal/content/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/redis/go-redis/v9"
)

type AuthorRepository interface {
	Add(author entity.Author) error
	Update(author entity.Author) error
	Delete(id int32) error
	Get(ctx context.Context, id int32) (*entity.Author, error)
	List(ctx context.Context) ([]entity.Author, error)
	Exists(name string) (bool, error)
}

type authorRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewAuthorRepository(db *sql.DB, redis *redis.Client) AuthorRepository {
	return &authorRepository{db, redis}
}

func (r *authorRepository) Add(author entity.Author) error {
	query := `INSERT INTO authors (name, img_url, description) VALUES ($1, $2, $3)`
	if _, err := r.db.Exec(query, author.Name, author.ImgURL, author.Description); err != nil {
		log.Error("[Repo][AddAuthor] ", err.Error())
		return app.NewAppError(500, "failed to add author")
	}
	// invalidate cache
	if r.redis != nil {
		_ = r.redis.Del(context.Background(), "authors").Err()
	}
	return nil
}

func (r *authorRepository) Update(author entity.Author) error {
	query := `UPDATE authors SET name=$1, img_url=$2, description=$3 WHERE id=$4`
	result, err := r.db.Exec(query, author.Name, author.ImgURL, author.Description, author.ID)
	if err != nil {
		log.Error("[Repo][UpdateAuthor] ", err.Error())
		return app.NewAppError(500, "failed to update author")
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return app.ErrNotFound
	}
	if r.redis != nil {
		_ = r.redis.Del(context.Background(), "authors").Err()
	}
	return nil
}

func (r *authorRepository) Delete(id int32) error {
	query := `DELETE FROM authors WHERE id=$1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteAuthor] ", err.Error())
		return app.NewAppError(500, "failed to delete author")
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return app.ErrNotFound
	}
	if r.redis != nil {
		_ = r.redis.Del(context.Background(), "authors").Err()
	}
	return nil
}

func (r *authorRepository) Get(ctx context.Context, id int32) (*entity.Author, error) {
	var a entity.Author
	query := `SELECT id, name, img_url, description FROM authors WHERE id=$1`
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Name, &a.ImgURL, &a.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, app.ErrNotFound
		}
		log.Error("[Repo][GetAuthor] ", err.Error())
		return nil, app.NewAppError(500, "failed to fetch author")
	}
	return &a, nil
}

func (r *authorRepository) List(ctx context.Context) ([]entity.Author, error) {
	var authors []entity.Author
	redisKey := "cache:authors:list"

	if r.redis != nil {
		raw, err := cache.Get(ctx, redisKey, r.redis)
		if err != nil && err != redis.Nil {
			log.Warn("[Repo][ListAuthors] Redis error: ", err)
		}
		if err == nil && raw != "" {
			if err := json.Unmarshal([]byte(raw), &authors); err == nil {
				return authors, nil
			} else {
				log.Warn("[Repo][ListAuthors] Failed to unmarshal Redis cache: ", err)
			}
		}
	}

	rows, err := r.db.QueryContext(ctx, `SELECT id, name, img_url, description FROM authors ORDER BY id`)
	if err != nil {
		log.Error("[Repo][ListAuthors] Query error: ", err)
		return nil, app.NewAppError(500, "failed to list authors")
	}
	defer rows.Close()

	for rows.Next() {
		var a entity.Author
		if err := rows.Scan(&a.ID, &a.Name, &a.ImgURL, &a.Description); err != nil {
			log.Error("[Repo][ListAuthors] Scan error: ", err)
			return nil, app.NewAppError(500, "failed to scan author")
		}
		authors = append(authors, a)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListAuthors] Row iteration error: ", err)
		return nil, app.NewAppError(500, "error iterating rows")
	}

	if r.redis != nil {
		if b, err := json.Marshal(authors); err == nil {
			_ = cache.Set(ctx, redisKey, string(b), r.redis, cache.ExpInstant)
		} else {
			log.Warn("[Repo][ListAuthors] Failed to marshal data for cache: ", err)
		}
	}

	return authors, nil
}

func (r *authorRepository) Exists(name string) (bool, error) {
	query := `SELECT 1 FROM authors WHERE name = $1`
	var exists int
	if err := r.db.QueryRow(query, name).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Error("[Repo][ExistsAuthor] ", err.Error())
		return false, app.NewAppError(500, "failed to check author existence")
	}
	return true, nil
}
