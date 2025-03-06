package repo

import (
	"database/sql"

	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

func (r *userRepository) AdminActivation(adminID int32) error {
	query := `UPDATE users SET is_verified=true WHERE id=$1`
	_, err := r.DB.Exec(query, adminID)
	if err != nil {
		log.Error("[Repo][userRepo.AdminActivation] Error Exec: ", err)
		return app.NewAppError(500, "failed to update user activation status")
	}
	return nil
}

func (r *userRepository) Auth(id int32) (userEntity.UserAuth, error) {
	query := `SELECT id, name, email, COALESCE(img_url, ''), is_admin, is_verified  FROM users WHERE id=$1`
	var user userEntity.UserAuth
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.ImgUrl, &user.IsAdmin, &user.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, app.NewAppError(404, "user not found")
		}
		return userEntity.UserAuth{}, app.NewAppError(500, err.Error())
	}
	return user, nil
}
