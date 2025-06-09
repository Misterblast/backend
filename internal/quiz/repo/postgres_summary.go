package repo

import (
	"fmt"
	"strconv"

	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

func (r *quizRepository) GetAvgTotal(userID int, filter map[string]string) (int, float64, error) {
	baseQuery := `
		SELECT COUNT(qs.*), COALESCE(AVG(qs.grade), 0)
		FROM quiz_submissions qs
		JOIN sets s ON qs.set_id = s.id
		WHERE qs.user_id = $1
	`

	args := []interface{}{userID}
	argIdx := 2
	// log.Debug("[QuizRepo][GetAvgTotal] user_id: %d", userID)
	// log.Debug("[QuizRepo][GetAvgTotal] filter: %v", filter)

	if lessonIDStr, ok := filter["lesson_id"]; ok && lessonIDStr != "" {
		lessonID, err := strconv.Atoi(lessonIDStr)
		// log.Debug("[QuizRepo][GetAvgTotal] lesson_id: %s", lessonIDStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid lesson_id: %v", err)
		}
		baseQuery += fmt.Sprintf(" AND s.lesson_id = $%d", argIdx)
		args = append(args, lessonID)
		argIdx++
	}

	var count int
	var avg float64
	err := r.db.QueryRow(baseQuery, args...).Scan(&count, &avg)
	if err != nil {
		log.Error("[QuizRepo][GetAvgTotal] Error executing query: %v", err)
		return 0, 0, fmt.Errorf("[QuizRepo][GetAvgTotal] Error executing query: %v", err)
	}

	return count, avg, nil
}
