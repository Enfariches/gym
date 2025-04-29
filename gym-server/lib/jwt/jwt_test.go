package jwt

import (
	"health/internal/domain/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	wantPassHash, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
)

func TestNewToken_ValidToken(t *testing.T) {
	timeDuration := time.Duration(time.Minute * 1)

	f := func(reqUser models.AuthUser, duration time.Duration) {
		claims := wantMapClaims(reqUser, duration)
		tokenString, err := NewToken(reqUser, time.Hour)

		assert.NoError(t, err, "Failed to generate token")

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Секретный ключ из NewToken()
			return []byte("gymnastic"), nil
		},
		)

		assert.NoError(t, err, "Failed to parse token")
		assert.True(t, token.Valid, "Token is invalid")
	}

	// Корректный случай
	f(models.AuthUser{
		Id:       1,
		Email:    "test@example.ru",
		PassHash: wantPassHash,
	}, timeDuration)
}

func TestNewToken_InvalidToken(t *testing.T) {
	timeDuration := time.Duration(time.Nanosecond)

	f := func(reqUser models.AuthUser, duration time.Duration) {
		claims := wantMapClaims(reqUser, duration)
		tokenString, err := NewToken(reqUser, duration)

		assert.NoError(t, err, "Failed to generate token")

		//Дожидаемся просрочки токена
		time.Sleep(time.Second)

		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Секретный ключ из NewToken()
			return []byte("gymnastic"), nil
		})

		assert.Error(t, err, "Expected error for expired token")
		assert.Contains(t, err.Error(), "token is expired", "Error should indicate expiration")
	}

	// Корректный случай
	f(models.AuthUser{
		Id:       2,
		Email:    "admin@yandex.ru",
		PassHash: wantPassHash,
	}, timeDuration)
}

func wantMapClaims(user models.AuthUser, duration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"uid":   user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
	}
}
