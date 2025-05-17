package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"health/internal/domain/models"
	"health/internal/storage"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
	pq "github.com/lib/pq"
)

type Scanner interface {
	Scan(dest ...any) error
}

type Storage struct {
	db *goqu.Database
}

func New(log *slog.Logger, databaseURL string) (*Storage, error) {
	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	db := goqu.New("postgres", dbConn)

	return &Storage{db: db}, nil
}

func HandleDBError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return storage.ErrUserExists
		default:
			return fmt.Errorf("unknown error: %v", pqErr.Error())
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrUserNotFound
	}

	return err
}

func scanSchedule(s Scanner, opts ...func(s *models.Schedule)) (*models.Schedule, error) {
	schedule := &models.Schedule{}

	err := s.Scan(&schedule.ID, &schedule.CronExpression, &schedule.IsActive, &schedule.MediaID, &schedule.AdminID, &schedule.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	for _, opt := range opts {
		opt(schedule)
	}

	return schedule, nil
}

func scanAdmin(s Scanner, opts ...func(s *models.Admin)) (*models.Admin, error) {
	admin := &models.Admin{}

	err := s.Scan(&admin.ID, &admin.Name, &admin.Surname, &admin.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	for _, opt := range opts {
		opt(admin)
	}

	return admin, nil
}

func scanMedia(s Scanner, opts...func(s *models.Media)) (*models.Media, error) {
	media := &models.Media{}

    err := s.Scan(&media.ID, &media.AdminID, &media.DepartmentID, &media.CreatedAt)
    if err != nil {
        return nil, fmt.Errorf("failed to scan row: %w", err)
    }

    for _, opt := range opts {
        opt(media)
    }

    return media, nil
}
