package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"health/internal/domain/models"
	ctxkey "health/lib/ctxkey"

	"github.com/doug-martin/goqu/v9"
)

func (s *Storage) CreateSchedules(ctx context.Context, schedules []*models.Schedule) ([]*models.Schedule, error) {
	const op = "postgres.CreateSchedules"

	admin_id := ctx.Value(ctxkey.UserKey).(int64)
	records := makeRecords(schedules, admin_id)

	query, args, _ := goqu.
		Insert("schedules").
		Rows(records).
		Returning("id", "cron_expression", "is_active", "video_id", "admin_id", "created_at").
		ToSQL()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	resultSchedules := make([]*models.Schedule, 0, len(schedules))

	for rows.Next() {
		s, err := readScheduleRows(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		resultSchedules = append(resultSchedules, s)
	}

	return resultSchedules, nil
}

func readScheduleRows(row *sql.Rows) (*models.Schedule, error) {
	schedule := &models.Schedule{}
	err := row.Scan(&schedule.ID, &schedule.CronExpression, &schedule.IsActive, &schedule.VideoID, &schedule.AdminID, &schedule.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return schedule, nil
}

func makeRecords(schedules []*models.Schedule, admin_id int64) []goqu.Record {
	records := make([]goqu.Record, 0, len(schedules))

	for _, s := range schedules {
		records = append(records, goqu.Record{
			"cron_expression": s.CronExpression,
			"is_active":       s.IsActive,
			"video_id":        s.VideoID,
			"admin_id":        admin_id,
		})
	}

	return records
}
