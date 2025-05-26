package repo

import (
	"fmt"
	"strings"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/lib/pq"
)

func (r *questionRepository) UpsertAndSyncAnswers(questionID int32, answers []questionEntity.SetAnswer) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Error("[Repo][UpsertAndSyncAnswers] Error beginning transaction: ", err)
		return app.NewAppError(500, "failed to begin transaction")
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id FROM answers WHERE question_id = $1", questionID)
	if err != nil {
		log.Error("[Repo][UpsertAndSyncAnswers] Error fetching existing answers: ", err)
		return app.NewAppError(500, "failed to fetch existing answers")
	}
	defer rows.Close()

	existingIDs := make(map[int32]bool)
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			log.Error("[Repo][UpsertAndSyncAnswers] Error scanning answer id: ", err)
			return app.NewAppError(500, "failed to scan answer id")
		}
		existingIDs[id] = true
	}

	inputIDs := make(map[int32]bool)
	for _, ans := range answers {
		if ans.ID > 0 {
			inputIDs[ans.ID] = true
		}
	}

	var idsToDelete []int32
	for id := range existingIDs {
		if !inputIDs[id] {
			idsToDelete = append(idsToDelete, id)
		}
	}

	if len(idsToDelete) > 0 {
		_, err := tx.Exec(`DELETE FROM answers WHERE id = ANY($1)`, pq.Array(idsToDelete))
		if err != nil {
			log.Error("[Repo][UpsertAndSyncAnswers] Error deleting old answers: ", err)
			return app.NewAppError(500, "failed to delete old answers")
		}
	}

	var insertValues []string
	var insertArgs []interface{}
	var upsertValues []string
	var upsertArgs []interface{}

	insertIdx := 1
	upsertIdx := 1

	for _, ans := range answers {
		if ans.ID == 0 {
			insertValues = append(insertValues, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", insertIdx, insertIdx+1, insertIdx+2, insertIdx+3, insertIdx+4))
			insertArgs = append(insertArgs, questionID, ans.Code, ans.Content, ans.ImgURL, ans.IsAnswer)
			insertIdx += 5
		} else {
			upsertValues = append(upsertValues, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)", upsertIdx, upsertIdx+1, upsertIdx+2, upsertIdx+3, upsertIdx+4, upsertIdx+5))
			upsertArgs = append(upsertArgs, ans.ID, questionID, ans.Code, ans.Content, ans.ImgURL, ans.IsAnswer)
			upsertIdx += 6
		}
	}

	if len(insertValues) > 0 {
		insertQuery := `INSERT INTO answers (question_id, code, content, img_url, is_answer) VALUES ` + strings.Join(insertValues, ", ")
		if _, err := tx.Exec(insertQuery, insertArgs...); err != nil {
			log.Error("[Repo][UpsertAndSyncAnswers] Error inserting new answers: ", err)
			return app.NewAppError(500, "failed to insert new answers")
		}
	}

	if len(upsertValues) > 0 {
		upsertQuery := `
			INSERT INTO answers (id, question_id, code, content, img_url, is_answer) VALUES
		` + strings.Join(upsertValues, ", ") + `
			ON CONFLICT (id) DO UPDATE SET
				question_id = EXCLUDED.question_id,
				code = EXCLUDED.code,
				content = EXCLUDED.content,
				img_url = EXCLUDED.img_url,
				is_answer = EXCLUDED.is_answer
		`
		if _, err := tx.Exec(upsertQuery, upsertArgs...); err != nil {
			log.Error("[Repo][UpsertAndSyncAnswers] Error upserting answers: ", err)
			return app.NewAppError(500, "failed to upsert answers")
		}
	}

	return tx.Commit()
}
