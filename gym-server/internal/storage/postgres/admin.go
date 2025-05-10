package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"health/internal/domain/models"
	"health/internal/storage"
	ctxkey "health/lib/ctxkey"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) Admin(ctx context.Context, admin_id int64) (*models.Admin, error) {
	const op = "postgres.Admin"

	query, args, _ := goqu.
		From("admins").
		Select(&models.Admin{}).
		Where(goqu.Ex{"id": admin_id}).
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)
	admin, err := readAdminRow(row)
	if err != nil {
		if HandleDBError(err) == storage.ErrUserNotFound {
			return &models.Admin{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return &models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	return admin, nil
}

func (s *Storage) UpdateAdmin(ctx context.Context, updateFields map[string]any) (*models.Admin, error) {
	const op = "postgres.UpdateAdmin"

	admin_id := ctx.Value(ctxkey.UserKey)

	query, args, _ := goqu.
		Update("admins").
		Set(updateFields).
		Where(goqu.C("id").Eq(admin_id)).
		Returning("departament", "email", "id", "name", "surname").
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	updatedAdmin, err := readAdminRow(row)
	if err != nil {
		return &models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedAdmin, nil
}

func readAdminRow(row *sql.Row) (*models.Admin, error) {
	var admin models.Admin

	err := row.Scan(&admin.Departament, &admin.Email, &admin.ID, &admin.Name, &admin.Surname)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &admin, nil
}
