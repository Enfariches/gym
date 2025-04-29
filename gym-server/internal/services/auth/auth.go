package auth

import (
	"context"
	"errors"
	"fmt"
	"health/internal/config"
	"health/internal/domain/models"
	"health/internal/storage"
	"health/lib/jwt"
	"health/lib/logger/sl"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userChecker  UserChecker
	userProvider UserProvider
	smtpConfig   config.SMTPConfig
	tokenTTL     time.Duration
	authTokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, authUser *models.AuthUser) error
}

type UserChecker interface {
	CheckUser(ctx context.Context, authUser *models.AuthUser) error
}

type UserProvider interface {
	User(ctx context.Context, email, source string) (models.AuthUser, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordIsIncorrect = errors.New("password is incorrect")
)

// New возвращает инстанс Auth сервиса
func New(log *slog.Logger, userSaver UserSaver, userChecker UserChecker, userProvider UserProvider,
	smtpConfig config.SMTPConfig, tokenTTL, authTokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userChecker:  userChecker,
		userProvider: userProvider,
		smtpConfig:   smtpConfig,
		tokenTTL:     tokenTTL,
		authTokenTTL: authTokenTTL,
	}
}

// Register регистрирует нового пользователя в системе и возвращает ошибку, если она имеется.
//
// Если пользователь уже существует в системе вернется ошибка
func (a *Auth) RegisterNewUser(ctx context.Context, email, password, source string) (string, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	authUser := models.AuthUser{
		Email:    email,
		PassHash: passHash,
		Source:   source,
	}

	if err := a.userChecker.CheckUser(ctx, &authUser); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn(storage.ErrUserExists.Error(), sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		log.Warn("failed to check user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	authToken, err := jwt.NewAuthToken(authUser, a.authTokenTTL)
	if err != nil {
		log.Error("failed to generate auth token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// err = smtp.SendMail(a.smtpConfig, authToken, email)
	// if err != nil {
	// 	log.Error("failed to send email", sl.Err(err))
	// 	return "", fmt.Errorf("%s: %w", op, err)
	// }

	return authToken, nil
}

// Login проверяет наличия пользователя в системе, если он есть возвращается токен.
//
// Если пользователь есть в системе, но пароль не правильный, то вернется ошибка.
// Если пользователя нет в системе, то вернется ошибка.
func (a *Auth) Login(ctx context.Context, email, password, source string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	user, err := a.userProvider.User(ctx, email, source)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn(storage.ErrUserNotFound.Error(), sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("failed to compare password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrPasswordIsIncorrect)
	}

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate JWT token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) VerifyRegister(ctx context.Context, authToken string) error {
	const op = "auth.VerifyRegister"

	log := a.log.With(slog.String("op", op))

	authUser, err := jwt.ParseAuthToken(authToken)
	if err != nil {
		log.Error("failed to parse token or token has expired", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log = log.With(slog.String("email", authUser.Email), slog.String("source", authUser.Source))

	err = a.userSaver.SaveUser(ctx, authUser)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn(storage.ErrUserExists.Error(), sl.Err(err))

			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		log.Error("failed to save user in system", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("registered new user")
	return nil
}

func (a *Auth) ResetPassword(ctx context.Context, email, source string) (string, error) {
	const op = "auth.ResetPassword"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	if err := a.userChecker.CheckUser(ctx, &models.AuthUser{Email: email, Source: source}); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn(storage.ErrUserExists.Error(), sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	resetToken, err := jwt.NewResetToken(email, source, a.authTokenTTL)
	if err != nil {
		log.Error("failed to generate reset token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// err = smtp.SendMail(a.smtpConfig, resetToken, email)
	// if err != nil {
	// 	log.Error("failed to send email", sl.Err(err))
	// 	return "", fmt.Errorf("%s: %w", op, err)
	// }

	return resetToken, nil
}

func (a *Auth) ChangePassword(ctx context.Context, resetToken, newPassword string) error {
	return nil
}
