package repo

import (
	"database/sql"
	"fmt"

	emailEntity "github.com/ghulammuzz/misterblast/internal/email/entity"
	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Add(user userEntity.Register, IsVerified bool) error
	Check(user userEntity.UserLogin) (*userEntity.UserJWT, error)
	Exists(id int32) (bool, error)
	List(filter map[string]string, page, limit int) (*response.PaginateResponse, error)
	Detail(id int32) (userEntity.DetailUser, error)
	Edit(id int32, user userEntity.EditUser) error
	Delete(id int32) error
	Auth(id int32) (userEntity.UserAuth, error)
	AdminActivation(adminID int32) error
	GetIDByEmail(email string) (int32, error)
	EditPassword(id int32, newPass string) error
	SetDeeplink(userID int32, token string, expiresAt int64) error
	GetDeeplink(token string) (emailEntity.DeeplinkResponse, error)
	GenerateToken() (string, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Exists(id int32) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`
	err := r.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, app.NewAppError(500, "failed to check if user exists")
	}
	return exists, nil
}

func (r *userRepository) Add(user userEntity.Register, IsVerified bool) error {

	query := `INSERT INTO users (name, email, password, img_url, is_verified) VALUES ($1, $2, $3, $4, $5)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query, user.Name, user.Email, hashedPassword, nil, IsVerified)
	return err
}

func (r *userRepository) Check(user userEntity.UserLogin) (*userEntity.UserJWT, error) {
	userResult := userEntity.UserJWT{}
	query := "SELECT id, email, password, is_admin, is_verified FROM users WHERE email=$1"
	err := r.DB.QueryRow(query, user.Email).Scan(&userResult.ID, &userResult.Email, &userResult.Password, &userResult.IsAdmin, &userResult.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewAppError(404, "user not found")
		}
		return nil, app.NewAppError(500, "failed to get user data")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(user.Password)); err != nil {
		return nil, app.NewAppError(400, "wrong password")
	}

	return &userResult, nil
}

func (r *userRepository) List(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	baseQuery := `FROM users WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Handle pencarian (search)
	searchClause := ""
	if search, exists := filter["search"]; exists && search != "" {
		searchClause = ` AND (LOWER(name) LIKE LOWER($` + fmt.Sprintf("%d", argCount) + `) OR LOWER(email) LIKE LOWER($` + fmt.Sprintf("%d", argCount+1) + `))`
		args = append(args, "%"+search+"%", "%"+search+"%")
		argCount += 2
	}

	// Query untuk menghitung total user sebelum pagination
	countQuery := `SELECT COUNT(*) ` + baseQuery + searchClause
	var total int64
	if err := r.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, app.NewAppError(500, "failed to count users")
	}

	// Query untuk mengambil user dengan pagination
	query := `SELECT id, name, email, COALESCE(img_url, '') ` + baseQuery + searchClause + ` ORDER BY id LIMIT $` + fmt.Sprintf("%d", argCount) + ` OFFSET $` + fmt.Sprintf("%d", argCount+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, app.NewAppError(500, err.Error())
	}
	defer rows.Close()

	var users []userEntity.ListUser
	for rows.Next() {
		var user userEntity.ListUser
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.ImgUrl); err != nil {
			return nil, app.NewAppError(500, err.Error())
		}
		users = append(users, user)
	}

	response := &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  users,
	}

	return response, nil
}

func (r *userRepository) Detail(id int32) (userEntity.DetailUser, error) {
	query := `SELECT id, name, email, COALESCE(img_url, '') FROM users WHERE id=$1`
	var user userEntity.DetailUser
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.ImgUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, app.NewAppError(404, "question not found")
		}
		return userEntity.DetailUser{}, app.NewAppError(500, err.Error())
	}
	return user, nil
}

func (r *userRepository) Edit(id int32, user userEntity.EditUser) error {
	query := `UPDATE users SET name=$1, email=$2, password=$3, img_url=$4, updated_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$5`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query, user.Name, user.Email, hashedPassword, user.ImgUrl, id)
	return err
}

func (r *userRepository) Delete(id int32) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteUser] Error Exec: ", err)
		return app.NewAppError(500, "failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("[Repo][DeleteUser] Error RowsAffected: ", err)
		return app.NewAppError(500, "failed to check rows affected")
	}
	if rowsAffected == 0 {
		return app.ErrNotFound
	}

	return nil
}

func (r *userRepository) GetIDByEmail(email string) (int32, error) {
	var id int32
	query := `SELECT id FROM users WHERE email=$1`
	err := r.DB.QueryRow(query, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, app.NewAppError(404, "user not found")
		}
		return 0, app.NewAppError(500, "failed to get user ID")
	}
	return id, nil
}

func (r *userRepository) EditPassword(id int32, newPass string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "UPDATE users SET password = $1, updated_at = EXTRACT(EPOCH from now()) WHERE id = $2"
	_, err = r.DB.Exec(query, hashedPassword, id)
	if err != nil {
		return app.NewAppError(500, "failed to update password")
	}
	return nil
}
