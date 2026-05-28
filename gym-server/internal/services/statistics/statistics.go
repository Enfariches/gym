package statistics

import (
	"context"
	"health/internal/domain/models"
	"health/lib/logger/sl"
	"log/slog"
)

type Statistics struct {
	log                *slog.Logger
	statisticsProvider StatisticsProvider
}

type StatisticsProvider interface {
	CreateStatistics(ctx context.Context, statProgress string, percentView, media_id int64) error
	GetEmployeeStatistics(ctx context.Context, employee_id, media_id int64) (*models.Statistics, error)

	ListMediaStatistics(ctx context.Context, media_id int64) ([]*models.Statistics, error)
	ListEmployeeStatistics(ctx context.Context, employee_id int64) ([]*models.Statistics, error)
	ListDepartmentStatistics(ctx context.Context, department_id int64) ([]*models.Statistics, error)
}

func New(log *slog.Logger, statisticsProvider StatisticsProvider) *Statistics {
	return &Statistics{
		log:                log,
		statisticsProvider: statisticsProvider,
	}
}

func (s *Statistics) CreateStatistics(ctx context.Context, statProgress string, percentView, media_id int64) error {
	const op = "statistics.CreateStatistics"
	log := s.log.With("op", op)

	err := s.statisticsProvider.CreateStatistics(ctx, statProgress, percentView, media_id)
	if err != nil {
		log.Error("failed to create statistics", sl.Err(err))
		return err
	}

	return nil
}

func (s *Statistics) GetEmployeeStatistics(ctx context.Context, employee_id, media_id int64) (*models.Statistics, error) {
	const op = "statistics.GetEmployeeStatistics"
	log := s.log.With("op", op)

	stats, err := s.statisticsProvider.GetEmployeeStatistics(ctx, employee_id, media_id)
	if err != nil {
		log.Error("failed to get employee statistics", sl.Err(err))
		return nil, err
	}

	return stats, nil
}

func (s *Statistics) ListMediaStatistics(ctx context.Context, media_id int64) ([]*models.Statistics, error) {
	const op = "statistics.ListMediaStatistics"
	log := s.log.With("op", op)

	stats, err := s.statisticsProvider.ListMediaStatistics(ctx, media_id)
	if err != nil {
		log.Error("failed to list media statistics", sl.Err(err))
		return nil, err
	}

	return stats, nil
}

func (s *Statistics) ListEmployeeStatistics(ctx context.Context, employee_id int64) ([]*models.Statistics, error) {
	const op = "statistics.ListEmployeeStatistics"
	log := s.log.With("op", op)

	stats, err := s.statisticsProvider.ListEmployeeStatistics(ctx, employee_id)
	if err != nil {
		log.Error("failed to list employee statistics", sl.Err(err))
		return nil, err
	}

	return stats, nil
}

func (s *Statistics) ListDepartmentStatistics(ctx context.Context, department_id int64) ([]*models.Statistics, error) {
	const op = "statistics.ListDepartmentStatistics"
	log := s.log.With("op", op)

	stats, err := s.statisticsProvider.ListDepartmentStatistics(ctx, department_id)
	if err != nil {
		log.Error("failed to list department statistics", sl.Err(err))
		return nil, err
	}

	return stats, nil
}
