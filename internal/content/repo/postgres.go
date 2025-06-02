package repo

// import (
// 	"context"
// 	"database/sql"

// 	contentEntity "github.com/ghulammuzz/misterblast/internal/content/entity"
// 	"github.com/redis/go-redis/v9"
// )

// type ContentRepository interface {
// 	// Questions
// 	Add(content contentEntity.Content, lang string) error
// 	List(ctx context.Context) ([]contentEntity.Content, error)
// 	Delete(id int32) error
// 	Detail(ctx context.Context, id int32) (contentEntity.Content, error)
// 	Exists(setID int32, number int) (bool, error)
// 	Edit(id int32, question contentEntity.Content) error
// }

// type contentRepository struct {
// 	db    *sql.DB
// 	redis *redis.Client
// }

// func NewContentRepository(db *sql.DB, redis *redis.Client) ContentRepository {
// 	return &contentRepository{db, redis}
// }
