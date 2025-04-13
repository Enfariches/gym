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

	f := func(rAdmin models.Admin, duration time.Duration) {
		claims := wantMapClaims(rAdmin, duration)
		tokenString, err := NewToken(rAdmin, time.Hour)

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
	f(models.Admin{
		ID:       "1",
		Email:    "test@example.ru",
		PassHash: wantPassHash,
	}, timeDuration)
}

func TestNewToken_InvalidToken(t *testing.T) {
	timeDuration := time.Duration(time.Nanosecond)

	f := func(rAdmin models.Admin, duration time.Duration) {
		claims := wantMapClaims(rAdmin, duration)
		tokenString, err := NewToken(rAdmin, duration)

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
	f(models.Admin{
		ID:       "2",
		Email:    "admin@yandex.ru",
		PassHash: wantPassHash,
	}, timeDuration)
}

func wantMapClaims(admin models.Admin, duration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"uid":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(duration).Unix(),
	}
}
