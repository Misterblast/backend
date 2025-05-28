package repo

import (
	"fmt"
	"strings"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

func (r *questionRepository) UpsertAndSyncAnswers(questionID int32, answers []questionEntity.SetAnswer) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Error("[Repo][UpsertAndSyncAnswers] Error beginning transaction: ", err)
		return app.NewAppError(500, "failed to begin transaction")
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM answers WHERE question_id = $1`, questionID)
	if err != nil {
		log.Error("[Repo][UpsertAndSyncAnswers] Error deleting old answers: ", err)
		return app.NewAppError(500, "failed to delete old answers")
	}

	log.Debug("[Repo][UpsertAndSyncAnswers] Deleted old answers for question ID: ", questionID)

	var insertValues []string
	var insertArgs []interface{}

	insertIdx := 1

	if len(answers) == 0 {
		return tx.Commit()
	}

	for _, ans := range answers {
		insertValues = append(insertValues, fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d)`, insertIdx, insertIdx+1, insertIdx+2, insertIdx+3, insertIdx+4))
		insertArgs = append(insertArgs, questionID, ans.Code, ans.Content, ans.ImgURL, ans.IsAnswer)
		insertIdx += 5
	}

	insertQuery := `INSERT INTO answers (question_id, code, content, img_url, is_answer) VALUES ` + strings.Join(insertValues, ", ")

	log.Debug("[Repo][UpsertAndSyncAnswers] Insert Query: ", insertQuery)

	if _, err := tx.Exec(insertQuery, insertArgs...); err != nil {
		log.Error("[Repo][UpsertAndSyncAnswers] Error inserting new answers: ", err)
		return app.NewAppError(500, "failed to insert new answers")
	}

	// if len(upsertValues) > 0 {
	// 	upsertQuery := `
	// 		INSERT INTO answers (id, question_id, code, content, img_url, is_answer) VALUES
	// 	` + strings.Join(upsertValues, ", ") + `
	// 		ON CONFLICT (id) DO UPDATE SET
	// 			question_id = EXCLUDED.question_id,
	// 			code = EXCLUDED.code,
	// 			content = EXCLUDED.content,
	// 			img_url = EXCLUDED.img_url,
	// 			is_answer = EXCLUDED.is_answer
	// 	`
	// 	if _, err := tx.Exec(upsertQuery, upsertArgs...); err != nil {
	// 		log.Error("[Repo][UpsertAndSyncAnswers] Error upserting answers: ", err)
	// 		return app.NewAppError(500, "failed to upsert answers")
	// 	}
	// }
	return tx.Commit()

}
