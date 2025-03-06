package jwt

import (
	"os"
	"time"

	"github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userResult entity.UserJWT) (string, error) {
	claims := jwt.MapClaims{
		"apps":     "misterblast-core",
		"email":    userResult.Email,
		"user_id":  userResult.ID,
		"is_admin": userResult.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
