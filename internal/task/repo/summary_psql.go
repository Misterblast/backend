package repo

func (r *TaskRepositoryImpl) GetAvgTotal(userID int32) (int, float64, error) {
	taskQuery := `SELECT COUNT(*), COALESCE(AVG(score), 0) FROM task_submissions WHERE user_id = $1`
	var taskCount int
	var avgTask float64
	err := r.db.QueryRow(taskQuery, userID).Scan(&taskCount, &avgTask)
	if err != nil {
		return 0, 0, err
	}

	return taskCount, avgTask, nil
}
