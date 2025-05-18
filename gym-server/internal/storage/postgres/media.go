package postgres

import (
	"context"
	"fmt"
	"health/internal/domain/models"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) SaveMediaPostgres(admin_id int64) (int64, int64, error) {
	const op = "postgres.SaveMediaPostgres"

	departmentId, err := s.getDepIdByAdminID(admin_id)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, err)
	}

	query, args, _ := goqu.Insert("mediafiles").
		Cols("admin_id", "department_id").
		Vals(goqu.Vals{admin_id, departmentId}).
		Returning("id").
		ToSQL()

	var media_id int64

	err = s.db.QueryRow(query, args...).Scan(&media_id)
	if err != nil {
		return 0, 0, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return media_id, departmentId, nil
}

func (s *Storage) Media(ctx context.Context, media_id int64) (*models.Media, error) {
	const op = "postgres.Media"

	query, args, _ := goqu.
		From("mediafiles").
		Select("id", "admin_id", "department_id", "created_at").
		Where(goqu.C("id").Eq(media_id)).
		ToSQL()

	row := s.db.QueryRow(query, args...)
	media, err := scanMedia(row)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return media, nil
}
