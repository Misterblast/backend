package entity

type Task struct {
	ID          int32
	Title       string
	Description string
	Content     string
	AttachedURL string
	UpdatedAt   int64
	CreatedAt   int64
	DeletedAt   int64
}
