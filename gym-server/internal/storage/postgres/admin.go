package postgres

import (
	"context"
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
		Select("id", "name", "surname", "email").
		Where(goqu.Ex{"id": admin_id}).
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	departmentName, err := s.GetDepNameByAdminID(admin_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department by admin id: %w", err)
	}

	admin, err := scanAdmin(row, func(s *models.Admin) { s.Department = departmentName })
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

	admin_id := ctx.Value(ctxkey.UserKey).(int64)

	if deptName, ok := updateFields["department"]; ok {
		err := s.updateDepNameByAdminID(admin_id, deptName.(string))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		delete(updateFields, "department")
	}

	if len(updateFields) == 0 {
		admin, err := s.Admin(ctx, admin_id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return admin, nil
	}

	query, args, _ := goqu.
		Update("admins").
		Set(updateFields).
		Where(goqu.C("id").Eq(admin_id)).
		Returning("id", "name", "surname", "email").
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	departmentName, err := s.GetDepNameByAdminID(admin_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department by admin id: %w", err)
	}

	updatedAdmin, err := scanAdmin(row, func(s *models.Admin) { s.Department = departmentName })
	if err != nil {
		return &models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedAdmin, nil
}

func (s *Storage) GetDepNameByAdminID(admin_id int64) (string, error) {
	const op = "postgres.GetDepNameByAdminID"

	query, args, _ := goqu.
		From("admins").
		Select(goqu.I("d.name")).
		Join(
			goqu.T("departments").As("d"),
			goqu.On(goqu.I("admins.department_id").Eq(goqu.I("d.id"))),
		).
		Where(goqu.I("admins.id").Eq(admin_id)).
		ToSQL()

	var departmentName string

	err := s.db.QueryRow(query, args...).Scan(&departmentName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return departmentName, nil
}

func (s *Storage) updateDepNameByAdminID(admin_id int64, deptName string) error {
	const op = "postgres.updateDepNameByAdminID"

	query, args, _ := goqu.
		Update("departments").
		Set(goqu.Record{"name": deptName}).
		Where(goqu.C("id").Eq(
			goqu.L("(SELECT department_id FROM admins WHERE id = ?)", admin_id),
		)).
		ToSQL()

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) getDepIdByAdminID(admin_id int64) (int64, error) {
	const op = "postgres.getDepIdByAdminID"

	query, args, _ := goqu.
		From("admins").
		Select("department_id").
		Where(goqu.C("id").Eq(admin_id)).
		Limit(1).
		ToSQL()

	var departmentID int64

	err := s.db.QueryRow(query, args...).Scan(&departmentID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return departmentID, nil
}

func (s *Storage) getMediaIDByAdminID(admin_id, media_id int64) (int64, error) {
	const op = "postgres.getMediaIDByAdminID"

	query, args, _ := goqu.
		From("mediafiles").
		Select("id").
		Where(goqu.And(goqu.C("admin_id").Eq(admin_id), goqu.C("id").Eq(media_id))).
		Limit(1).
		ToSQL()

	var mediaID int64

	err := s.db.QueryRow(query, args...).Scan(&mediaID)
	if err != nil {
		return 0, fmt.Errorf("%s: %s", op, "mediafiles not found")
	}

	return mediaID, nil
}
