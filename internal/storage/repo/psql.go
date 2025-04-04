package repo

import (
	"database/sql"
	"fmt"
	entity2 "github.com/ghulammuzz/misterblast/internal/storage/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/lib/pq"
	"net/http"
	"strings"
)

type StorageRepository interface {
	InsertAttachments(attachments []entity2.Attachment) ([]int64, error)
	DeleteAttachments(attachments []int64)
}

type StorageRepositoryImpl struct {
	db *sql.DB
}

func (r *StorageRepositoryImpl) DeleteAttachments(attachments []int64) {
	query := "DELETE FROM attachments WHERE id IN($1)"
	_, err := r.db.Exec(query, pq.Array(attachments))
	if err != nil {
		log.Error("[Repo.Attachment.DELETE] failed to delete attachments, cause: %s", err.Error())
	}
}

func (r *StorageRepositoryImpl) InsertAttachments(attachments []entity2.Attachment) ([]int64, error) {
	var attachmentIds []int64
	query := "INSERT INTO attachments (type, url) VALUES "
	var values []interface{}
	var placeholders []string

	for i, attachment := range attachments {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		values = append(values, attachment.Type, attachment.Url)
	}

	query += strings.Join(placeholders, ", ") + " RETURNING id"

	rows, err := r.db.Query(query, values...)
	if err != nil {
		log.Error("[Repo.Attachment.INSERT] failed to insert attachments, cause: %s", err.Error())
		return []int64{}, app.NewAppError(http.StatusInternalServerError, "failed to insert attachments")
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Error("[Repo.Attachment.INSERT] failed to scan attachment ID, cause: %s", err.Error())
			return []int64{}, app.NewAppError(http.StatusInternalServerError, "failed to retrieve attachment IDs")
		}
		attachmentIds = append(attachmentIds, id)
	}
	return attachmentIds, nil
}

func NewStorageRepository(db *sql.DB) StorageRepository {
	return &StorageRepositoryImpl{db: db}
}
