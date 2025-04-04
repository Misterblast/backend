package entity

type Submission struct {
	Id        int64  `db:"id"`
	TaskId    int64  `db:"task_id"`
	UserId    int64  `db:"user_id"`
	Answer    string `db:"answer"`
	Score     int32  `db:"score"`
	Feedback  string `db:"feedback"`
	ScoredBy  int64  `db:"scored_by"`
	ScoredAt  int64  `db:"scored_at"`
	CreatedAt int64  `db:"created_at"`
}
