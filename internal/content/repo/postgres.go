package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	contentEntity "github.com/ghulammuzz/misterblast/internal/content/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/redis/go-redis/v9"
)

type ContentRepository interface {
	Add(content contentEntity.Content, lang string) error
	List(filter map[string]string, ctx context.Context) (*response.PaginateResponse, error)
	Delete(id int32) error
	Detail(ctx context.Context, id int32) (contentEntity.Content, error)
	Edit(id int32, content contentEntity.Content) error
}

type contentRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewContentRepository(db *sql.DB, redis *redis.Client) ContentRepository {
	return &contentRepository{db, redis}
}

func (c *contentRepository) Add(content contentEntity.Content, lang string) error {
	query := `INSERT INTO content (title, description, img_url, site_url, lang) VALUES ($1, $2, $3, $4, $5)`
	_, err := c.db.Exec(query, content.Title, content.Desc, content.ImgURL, content.SiteURL, lang)
	return err
}

func (c *contentRepository) List(filter map[string]string, ctx context.Context) (*response.PaginateResponse, error) {
	var contents []contentEntity.Content

	redisKey := "cache:contents:list"
	if lang, ok := filter["lang"]; ok {
		redisKey += fmt.Sprintf(":lang=%s", lang)
	}
	if p, ok := filter["page"]; ok {
		redisKey += fmt.Sprintf(":page=%s", p)
	}
	if l, ok := filter["limit"]; ok {
		redisKey += fmt.Sprintf(":limit=%s", l)
	}

	if c.redis != nil {
		cached, err := cache.Get(ctx, redisKey, c.redis)
		if err != nil && err != redis.Nil {
			log.Warn("[ContentRepository.List] Redis error: ", err)
		}
		if err == nil && cached != "" {
			var cachedResp response.PaginateResponse
			if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
				return &cachedResp, nil
			} else {
				log.Warn("[ContentRepository.List] Failed to unmarshal Redis cache: ", err)
			}
		}
	}

	query := `SELECT id, title, description, img_url, site_url, lang FROM content`
	countQuery := `SELECT COUNT(*) FROM content`
	var args []interface{}
	var conditions []string
	argCounter := 1

	if l, ok := filter["lang"]; ok && (l == "id" || l == "en") {
		conditions = append(conditions, fmt.Sprintf("lang = $%d", argCounter))
		args = append(args, l)
		argCounter++
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	page := 1
	limit := 10
	if p, ok := filter["page"]; ok {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if l, ok := filter["limit"]; ok {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	offset := (page - 1) * limit
	query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("[ContentRepository.List] Error executing query: ", err)
		return nil, app.NewAppError(500, "failed to execute content list query")
	}
	defer rows.Close()

	for rows.Next() {
		var cont contentEntity.Content
		if err := rows.Scan(&cont.ID, &cont.Title, &cont.Desc, &cont.ImgURL, &cont.SiteURL, &cont.Lang); err != nil {
			log.Error("[ContentRepository.List] Error scanning row: ", err)
			return nil, app.NewAppError(500, "failed to scan content row")
		}
		contents = append(contents, cont)
	}

	var total int64
	if err := c.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Error("[ContentRepository.List] Error counting total records: ", err)
		return nil, app.NewAppError(500, "failed to count total records")
	}
	rs := &response.PaginateResponse{
		Total: total,
		Limit: limit,
		Page:  page,
		Data:  contents,
	}

	if c.redis != nil {
		if dataJSON, err := json.Marshal(rs); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), c.redis, cache.ExpSecond)
		} else {
			log.Warn("[ContentRepository.List] Failed to marshal data for Redis: ", err)
		}
	}

	return rs, nil
}

func (c *contentRepository) Delete(id int32) error {
	query := `DELETE FROM content WHERE id = $1`
	res, err := c.db.Exec(query, id)
	if err != nil {
		log.Error("[ContentRepository.Delete]", "Error executing delete query", err)
		return app.NewAppError(500, "failed to delete content")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error("[ContentRepository.Delete]", "Error getting rows affected", err)
		return app.NewAppError(500, "failed to get rows affected for delete operation")
	}

	if rowsAffected == 0 {
		log.Warn("[ContentRepository.Delete]", "No rows affected for delete operation", "id", id)
		return app.NewAppError(404, "content not found or already deleted")
	}

	return nil
}

func (c *contentRepository) Detail(ctx context.Context, id int32) (contentEntity.Content, error) {
	query := `SELECT id, title, description, img_url, site_url, lang FROM content WHERE id = $1`
	var cont contentEntity.Content
	err := c.db.QueryRowContext(ctx, query, id).Scan(&cont.ID, &cont.Title, &cont.Desc, &cont.ImgURL, &cont.SiteURL, &cont.Lang)
	if err != nil {
		if err == sql.ErrNoRows {
			return cont, app.NewAppError(404, "content not found")
		}
		return cont, app.NewAppError(500, err.Error())
	}
	return cont, nil
}

func (c *contentRepository) Edit(id int32, content contentEntity.Content) error {
	query := `UPDATE content SET title = $1, description = $2, img_url = $3, site_url = $4, lang = $5 WHERE id = $6`
	_, err := c.db.Exec(query, content.Title, content.Desc, content.ImgURL, content.SiteURL, content.Lang, id)
	if err != nil {
		log.Error("[ContentRepository.Edit]", "Error updating content", err)
		return fmt.Errorf("failed to update content: %w", err)
	}
	return err
}
