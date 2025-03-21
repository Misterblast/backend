package entity

type TaskAttachment struct {
	Id     int32  `json:"-" db:"id"`
	TaskId int32  `json:"-" db:"task_id"`
	Type   string `json:"type" db:"type"`
	Url    string `json:"attachment_url" db:"url"`
}

type Task struct {
	ID          int32  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Content     string `db:"content"`
	Author      int32  `db:"author"`
	UpdatedBy   int32  `db:"updated_by"`
	UpdatedAt   int64  `db:"updated_at"`
	CreatedAt   int64  `db:"created_at"`
	DeletedAt   int64  `db:"deleted_at"`
	Attachments []TaskAttachment
}
