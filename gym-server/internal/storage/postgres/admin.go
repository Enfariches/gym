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

	department_id := ctx.Value(ctxkey.DepartmentKey).(int64)

	query, args, _ := goqu.
		From("admins").
		Select("id", "name", "surname", "email").
		Where(goqu.Ex{"id": admin_id}).
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	departmentName, err := s.GetDepNameByDepID(department_id)
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
	department_id := ctx.Value(ctxkey.DepartmentKey).(int64)

	if deptName, ok := updateFields["department"]; ok {
		err := s.updateDepNameByDepID(department_id, deptName.(string))
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

	departmentName, err := s.GetDepNameByDepID(department_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department by admin id: %w", err)
	}

	updatedAdmin, err := scanAdmin(row, func(s *models.Admin) { s.Department = departmentName })
	if err != nil {
		return &models.Admin{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedAdmin, nil
}

func (s *Storage) GetDepNameByDepID(department_id int64) (string, error) {
	const op = "postgres.GetDepNameByDepID"

	query, args, _ := goqu.
		From("departments").
		Select("name").
		Where(goqu.C("id").Eq(department_id)).
		ToSQL()

	var departmentName string

	err := s.db.QueryRow(query, args...).Scan(&departmentName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return departmentName, nil
}

func (s *Storage) updateDepNameByDepID(department_id int64, deptName string) error {
	const op = "postgres.updateDepNameByDepID"

	query, args, _ := goqu.
		Update("departments").
		Set(goqu.Record{"name": deptName}).
		Where(goqu.C("id").Eq(department_id)).
		ToSQL()

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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
