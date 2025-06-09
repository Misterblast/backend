package repo

func (r *quizRepository) GetAvgTotal(userID int) (int, float64, error) {
	quizQuery := `SELECT COUNT(*), COALESCE(AVG(grade), 0) FROM quiz_submissions WHERE user_id = $1`
	var quizCount int
	var avgQuiz float64
	err := r.db.QueryRow(quizQuery, userID).Scan(&quizCount, &avgQuiz)
	if err != nil {
		return 0, 0, err
	}

	return quizCount, avgQuiz, nil
}
