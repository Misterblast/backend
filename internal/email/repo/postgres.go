package repo

import (
	"database/sql"
	"time"

	"github.com/ghulammuzz/misterblast/pkg/app"
)

type EmailRepository interface {
	SetOTP(adminID int32, otp string, expiresAt int64) error
	GetOTP(adminID int32) (string, int64, error)
}

type emailRepository struct {
	DB *sql.DB
}

func NewEmailRepository(db *sql.DB) EmailRepository {
	return &emailRepository{db}
}

func (r *emailRepository) SetOTP(adminID int32, otp string, expiresAt int64) error {
	if expiresAt <= time.Now().Unix() {
		return app.NewAppError(400, "Waktu kedaluwarsa tidak valid")
	}
	query := `
        INSERT INTO user_otps (admin_id, otp_code, expires_at) 
        VALUES ($1, $2, $3)
        ON CONFLICT (admin_id) 
        DO UPDATE SET otp_code = EXCLUDED.otp_code, expires_at = EXCLUDED.expires_at;
    `
	_, err := r.DB.Exec(query, adminID, otp, expiresAt)
	if err != nil {
		return app.NewAppError(500, "Gagal menyimpan OTP")
	}
	return nil
}

func (r *emailRepository) GetOTP(adminID int32) (string, int64, error) {
	var dbOtp string
	var expiresAt int64
	query := `SELECT otp_code, expires_at FROM user_otps WHERE admin_id=$1 LIMIT 1`
	err := r.DB.QueryRow(query, adminID).Scan(&dbOtp, &expiresAt)
	if err != nil {
		return "", 0, app.NewAppError(404, "OTP tidak ditemukan atau sudah kadaluarsa")
	}
	return dbOtp, expiresAt, nil
}
