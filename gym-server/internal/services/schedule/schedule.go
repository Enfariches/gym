package schedule

import (
	"context"
	"health/internal/domain/models"
	"health/lib/gcron"
	"health/lib/logger/sl"
	"log/slog"
)

type Schedule struct {
	log             *slog.Logger
	scheduleManager ScheduleManager
	jobScheduler    *gcron.JobScheduler
}

type ScheduleManager interface {
	CreateSchedules(ctx context.Context, schedules []*models.Schedule) ([]*models.Schedule, error)
}

func New(log *slog.Logger, scheduleManager ScheduleManager, jobScheduler *gcron.JobScheduler) *Schedule {
	jobScheduler.Start()

	return &Schedule{
		log:             log,
		scheduleManager: scheduleManager,
		jobScheduler:    jobScheduler,
	}
}

func (s *Schedule) CreateSchedules(ctx context.Context, schedules []*models.Schedule) ([]*models.Schedule, error) {
	const op = "schedule.CreateSchedule"
	log := s.log.With("op", op)

	schedules, err := s.scheduleManager.CreateSchedules(ctx, schedules)
	if err != nil {
		log.Error("failed to create schedule", sl.Err(err))
		return nil, err
	}

	// Обновление задачи планировщика
	s.jobScheduler.UpdateScheduler(schedules)

	return schedules, nil
}
