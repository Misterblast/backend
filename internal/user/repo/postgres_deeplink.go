package repo

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"

	emailEntity "github.com/ghulammuzz/misterblast/internal/email/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
)

func (r *userRepository) SetDeeplink(userID int32, token string, expiresAt int64) error {

	_, err := r.DB.Exec("INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiresAt)
	if err != nil {
		return app.NewAppError(500, "Gagal Set Deeplink")
	}
	return nil

}

func (r *userRepository) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", app.NewAppError(500, "Gagal Generate")
	}
	return hex.EncodeToString(bytes), nil
}

func (r *userRepository) GetDeeplink(token string) (emailEntity.DeeplinkResponse, error) {
	var resp emailEntity.DeeplinkResponse

	err := r.DB.QueryRow("SELECT token, user_id, expires_at FROM password_reset_tokens WHERE token=$1 LIMIT 1", token).Scan(&resp.Token, &resp.UserID, &resp.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, app.NewAppError(404, "Deeplink not found")
		}
		return resp, app.NewAppError(500, "Failed to retrieve Deeplink")
	}
	return resp, nil
}
