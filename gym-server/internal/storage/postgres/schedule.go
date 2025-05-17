package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"health/internal/domain/models"
	ctxkey "health/lib/ctxkey"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) CreateSchedules(ctx context.Context, schedules []*models.Schedule) ([]*models.Schedule, error) {
	const op = "postgres.CreateSchedules"

	admin_id := ctx.Value(ctxkey.UserKey).(int64)
	records := make([]goqu.Record, 0, len(schedules))

	for _, sc := range schedules {
		media_id, err := s.getMediaIDByAdminID(admin_id, sc.MediaID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		records = append(records, goqu.Record{
			"cron_expression": sc.CronExpression,
			"is_active":       sc.IsActive,
			"media_id":        media_id,
			"admin_id":        admin_id,
		})
	}

	query, args, _ := goqu.
		Insert("schedules").
		Rows(records).
		Returning("id", "cron_expression", "is_active", "media_id", "admin_id", "created_at").
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	resultSchedules := make([]*models.Schedule, 0, len(schedules))

	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		resultSchedules = append(resultSchedules, s)
	}

	return resultSchedules, nil
}

func (s *Storage) Schedule(ctx context.Context, schedule_id int64) (*models.Schedule, error) {
	const op = "postgres.Schedule"

	query, args, _ := goqu.
		From("schedules").
		Select("id", "cron_expression", "is_active", "media_id", "admin_id", "created_at").
		Where(goqu.Ex{"id": schedule_id}).
		Limit(1).
		ToSQL()
	row := s.db.QueryRow(query, args...)

	schedule, err := scanSchedule(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Schedule{}, fmt.Errorf("%s: %w", op, errors.New("schedule not found"))
		}
		return &models.Schedule{}, fmt.Errorf("%s: %w", op, err)
	}

	return schedule, nil
}

func (s *Storage) UpdateSchedule(ctx context.Context, schedule_id int64, updateFields map[string]any) (*models.Schedule, error) {
	const op = "postgres.UpdateSchedule"

	if mediaId, ok := updateFields["media_id"]; ok {
		admin_id := ctx.Value(ctxkey.UserKey).(int64)

		_, err := s.getMediaIDByAdminID(admin_id, mediaId.(int64))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	query, args, err := goqu.
		Update("schedules").
		Set(updateFields).
		Where(goqu.C("id").Eq(schedule_id)).
		Returning("id", "cron_expression", "is_active", "media_id", "admin_id", "created_at").
		Limit(1).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	row := s.db.QueryRow(query, args...)

	schedule, err := scanSchedule(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Schedule{}, fmt.Errorf("%s: %w", op, errors.New("schedule not found"))
		}
		return &models.Schedule{}, fmt.Errorf("%s: %w", op, err)
	}

	return schedule, nil
}

func (s *Storage) DeleteSchedule(ctx context.Context, schedule_id int64) (*models.Schedule, error) {
	const op = "postgres.DeleteSchedule"

	query, args, _ := goqu.
		Delete("schedules").
		Where(goqu.C("id").Eq(schedule_id)).
		Returning("id", "cron_expression", "is_active", "media_id", "admin_id", "created_at").
		Limit(1).
		ToSQL()

	row := s.db.QueryRow(query, args...)

	// Деактивируем удаленное расписание, чтобы прекратить его работу в планировщике.
	schedule, err := scanSchedule(row, func(s *models.Schedule) { s.IsActive = false })
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Schedule{}, fmt.Errorf("%s: %w", op, errors.New("schedule not found"))
		}
		return &models.Schedule{}, fmt.Errorf("%s: %w", op, err)
	}
	return schedule, nil
}

func (s *Storage) ListSchedule(ctx context.Context, admin_id int64) ([]*models.Schedule, error) {
	const op = "postgres.ListSchedules"

	query, args, _ := goqu.
		From("schedules").
		Select("id", "cron_expression", "is_active", "media_id", "admin_id", "created_at").
		Where(goqu.C("admin_id").Eq(admin_id)).
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	resultSchedules := make([]*models.Schedule, 0)

	for rows.Next() {
		s, err := scanSchedule(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		resultSchedules = append(resultSchedules, s)
	}

	return resultSchedules, nil
}
