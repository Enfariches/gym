package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"health/internal/domain/models"
	"health/internal/storage"
	"health/lib/logger/sl"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
	pq "github.com/lib/pq"
)

type Storage struct {
	db  *goqu.Database
	log *slog.Logger
}

func New(log *slog.Logger, databaseURL string) (*Storage, error) {
	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	db := goqu.New("postgres", dbConn)

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) SaveUser(ctx context.Context, authUser *models.AuthUser) error {
	const op = "postgres.SaveUser"

	log := s.log.With(slog.String("op", op), slog.String("email", authUser.Email))

	query, args, _ := goqu.Insert(authUser.Source).
		Cols("email", "passhash").
		Vals(goqu.Vals{authUser.Email, authUser.PassHash}).
		ToSQL()

	_, err := s.db.Exec(query, args...)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return HandleDBError(err)
	}
	return nil
}

func (s *Storage) CheckUser(ctx context.Context, authUser *models.AuthUser) (bool) {
	query, args, _ := goqu.From(authUser.Source).
		Select("id").
		Where(goqu.Ex{"email": authUser.Email}).
		Limit(1).
		ToSQL()
	
	var id int

	err := s.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		if HandleDBError(err) == storage.ErrUserNotFound {
			return false // пользователь не найден
		}
		return false
	}

	return true
}

func (s *Storage) User(ctx context.Context, email, source string) (models.AuthUser, error) {
	query, args, _ := goqu.From(source).
		Select("id", "email", "passhash").
		Where(goqu.Ex{"email": email}).
		Limit(1).
		ToSQL()

	var authUser models.AuthUser
	err := s.db.QueryRow(query, args...).Scan(&authUser.Id, &authUser.Email, &authUser.PassHash)
	if err != nil {
		if HandleDBError(err) == storage.ErrUserNotFound {
			return models.AuthUser{}, storage.ErrUserNotFound
		}
		return models.AuthUser{}, err
	}
	authUser.Source = source

	return authUser, nil
}

func HandleDBError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return storage.ErrUserExists
		case "23503":
			return errors.New("invalid reference (foreign key error)")
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrUserNotFound
	}

	return err
}
