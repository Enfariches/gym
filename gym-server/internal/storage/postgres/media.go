package postgres

import (
	"context"
	"fmt"
	"health/internal/domain/models"
	ctxkey "health/lib/ctxkey"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) SaveMediaPostgres(admin_id, departmentId int64, mediaTitle string) (int64, error) {
	const op = "postgres.SaveMediaPostgres"

	query, args, _ := goqu.
		Insert("mediafiles").
		Cols("title", "admin_id", "department_id").
		Vals(goqu.Vals{mediaTitle, admin_id, departmentId}).
		Returning("id").
		ToSQL()

	var media_id int64

	err := s.db.QueryRow(query, args...).Scan(&media_id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return media_id, nil
}

func (s *Storage) Media(ctx context.Context, media_id int64) (*models.Media, error) {
	const op = "postgres.Media"

	query, args, _ := goqu.
		From("mediafiles").
		Select("id", "title", "admin_id", "department_id", "created_at").
		Where(goqu.C("id").Eq(media_id)).
		ToSQL()

	row := s.db.QueryRow(query, args...)
	media, err := scanMedia(row)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return media, nil
}

func (s *Storage) ListMedia(ctx context.Context, admin_id int64) ([]*models.Media, error) {
	const op = "postgres.ListMedia"

	query, args, _ := goqu.
		From("mediafiles").
		Select("id", "title", "admin_id", "department_id", "created_at").
		Where(goqu.C("admin_id").Eq(admin_id)).
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, HandleDBError(err))
	}
	defer rows.Close()

	var medias []*models.Media
	for rows.Next() {
		mediaItem, err := scanMedia(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		medias = append(medias, mediaItem)
	}
	return medias, nil
}

func (s *Storage) DeleteMedia(ctx context.Context, media_id int64) error {
	const op = "postgres.DeleteMedia"

	admin_id := ctx.Value(ctxkey.UserKey).(int64)

	if _, err := s.getMediaIDByAdminID(admin_id, media_id); err != nil {
		return fmt.Errorf("%s: %s", op, "admin does not have such a media_id")
	}

	query, args, _ := goqu.
		Delete("mediafiles").
		Where(goqu.C("id").Eq(media_id)).
		ToSQL()

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleDBError(err))
	}

	return nil
}
