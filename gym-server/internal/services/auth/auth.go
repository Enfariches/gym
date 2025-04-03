package auth

import (
	"context"
	"errors"
	"fmt"
	"health/internal/domain/models"
	"health/internal/storage"
	"health/lib/jwt"
	"health/lib/logger/sl"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	admSaver    AdminSaver
	admProvider AdminProvider
	tokenTTL    time.Duration
}

type AdminSaver interface {
	SaveAdmin(ctx context.Context, email string, passHash []byte) (err error)
}

type AdminProvider interface {
	Admin(ctx context.Context, email string) (models.Admin, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// New возвращает инстанс Auth сервиса
func New(log *slog.Logger, admSaver AdminSaver, admProvider AdminProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		admSaver:    admSaver,
		admProvider: admProvider,
		tokenTTL:    tokenTTL,
	}
}

// Register регистрирует нового админа в системе и возвращает ошибку, если она имеется.
// 
// Если админ существует в системе вернется ошибка
func (a *Auth) Register(ctx context.Context, email, password string) error {
	const op = "auth.Register"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash password", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.admSaver.SaveAdmin(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrAdminExists) {
			a.log.Warn(storage.ErrAdminExists.Error(), sl.Err(err))

			return fmt.Errorf("%s: %w", op, storage.ErrAdminExists)
		}
		log.Error("failed to save admin in system", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("registered new admin")
	return nil
}

// Login проверяет наличия админа в системе, если он есть возвращается токен.
//
// Если админ есть в системе, но пароль не правильный, то вернется ошибка.
// Если админа нет в системе, то вернется ошибка.
func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	admin, err := a.admProvider.Admin(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrAdminNotFound) {
			a.log.Warn(storage.ErrAdminNotFound.Error(), sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get admin", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(admin.PassHash, []byte(password)); err != nil {
		log.Error("failed to compare password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(admin, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate JWT token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
