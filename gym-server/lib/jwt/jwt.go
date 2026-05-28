package jwt

import (
	"health/internal/domain/models"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	secret      = "mainGymUnique"
	authSecret  = "authGymUnique"
	resetSecret = "resetGymUnique"
)

func NewToken(user models.AuthUser, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":            user.ID,
		"department_id": user.DepartmentID,
		"source":        user.Source,
		"exp":           time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (int64, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return 0, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, err
	}

	return int64(claims["id"].(float64)), int64(claims["department_id"].(float64)), nil
}

func NewAuthToken(authUser models.AuthUser, duration time.Duration) (string, error) {
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    authUser.Email,
		"passhash": string(authUser.PassHash),
		"source":   authUser.Source,
		"exp":      time.Now().Add(duration).Unix(),
	})

	authTokenString, err := authToken.SignedString([]byte(authSecret))
	if err != nil {
		return "", err
	}

	return authTokenString, nil
}

func ParseAuthToken(tokenString string) (*models.AuthUser, error) {
	authToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(authSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !authToken.Valid {
		return nil, err
	}

	claims, ok := authToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return &models.AuthUser{
		Email:    claims["email"].(string),
		PassHash: []byte(claims["passhash"].(string)),
		Source:   claims["source"].(string),
	}, nil
}

func NewResetToken(email, source string, duration time.Duration) (string, error) {
	resetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"source": source,
		"exp":    time.Now().Add(duration).Unix(),
	})

	resetTokenString, err := resetToken.SignedString([]byte(resetSecret))
	if err != nil {
		return "", err
	}

	return resetTokenString, nil
}

func ParseResetToken(tokenString string) (*models.AuthUser, error) {
	authToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(resetSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !authToken.Valid {
		return nil, err
	}

	claims, ok := authToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return &models.AuthUser{
		Email:  claims["email"].(string),
		Source: claims["source"].(string),
	}, nil
}
