package jwt

import (
	"health/internal/domain/models"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "gymnastic"

func NewToken(admin models.Admin, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
		"uid":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
