package postgres

import (
	"context"
	"fmt"
	"health/internal/domain/models"
	"health/internal/storage"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) SaveUser(ctx context.Context, authUser *models.AuthUser) error {
	const op = "postgres.SaveUser"

	var (
		query string
		args  []any
		err   error
	)

	if authUser.Source == "employees" {
		query, args, _ = goqu.
			Insert(authUser.Source).
			Cols("email", "passhash").
			Vals(goqu.Vals{authUser.Email, authUser.PassHash}).
			ToSQL()
	} else {
		cteQuery := goqu.
			Insert("departments").
			Rows(goqu.Record{"name": ""}).
			Returning("id")

		query, args, _ = goqu.
			Insert(authUser.Source).
			Cols("email", "passhash", "department_id").
			Vals(goqu.Vals{authUser.Email, authUser.PassHash, goqu.L("(SELECT id FROM new_department)")}).
			With("new_department", cteQuery).
			ToSQL()
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return nil
}

func (s *Storage) CheckUser(ctx context.Context, authUser *models.AuthUser) error {
	const op = "postgres.CheckUser"

	query, args, _ := goqu.
		From(authUser.Source).
		Select("id").
		Where(goqu.Ex{"email": authUser.Email}).
		Limit(1).
		ToSQL()

	var id int

	err := s.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		if HandleDBError(err) == storage.ErrUserNotFound {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
}

func (s *Storage) User(ctx context.Context, email, source string) (*models.AuthUser, error) {
	const op = "postgres.User"

	query, args, _ := goqu.
		From(source).
		Select("id", "email", "department_id", "passhash").
		Where(goqu.Ex{"email": email}).
		Limit(1).
		ToSQL()

	var authUser models.AuthUser

	err := s.db.QueryRow(query, args...).Scan(&authUser.ID, &authUser.Email, &authUser.DepartmentID, &authUser.PassHash)
	if err != nil {
		if HandleDBError(err) == storage.ErrUserNotFound {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	authUser.Source = source

	return &authUser, nil
}

func (s *Storage) UpdateUserPassword(ctx context.Context, authUser *models.AuthUser) error {
	const op = "postgres.UpdatePassword"

	query, args, _ := goqu.
		Update(authUser.Source).
		Set(goqu.Record{"passhash": authUser.PassHash}).
		Where(goqu.Ex{"email": authUser.Email}).
		ToSQL()

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return nil
}
