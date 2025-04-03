package entity

type TaskAssignment struct {
	Id      int64  `json:"id" db:"id"`
	TaskId  int64  `json:"task_id" db:"task_id"`
	UserId  int64  `json:"user_id" db:"user_id"`
	Content string `json:"content" db:"content"`
	Score   int32  `json:"score" db:"score"`
}
