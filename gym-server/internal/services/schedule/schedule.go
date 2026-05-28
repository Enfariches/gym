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
	UpdateSchedule(ctx context.Context, schedule_id int64, updateFields map[string]any) (*models.Schedule, error)
	Schedule(ctx context.Context, schedule_id int64) (*models.Schedule, error)
	DeleteSchedule(ctx context.Context, schedule_id int64) (*models.Schedule, error)
	ListSchedule(ctx context.Context, admin_id int64) ([]*models.Schedule, error)
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
	//s.jobScheduler.UpdateScheduler(schedules)

	return schedules, nil
}

func (s *Schedule) GetSchedule(ctx context.Context, schedule_id int64) (*models.Schedule, error) {
	const op = "schedule.GetSchedule"
	log := s.log.With("op", op)

	schedule, err := s.scheduleManager.Schedule(ctx, schedule_id)
	if err != nil {
		log.Error("failed to get schedule", sl.Err(err))
		return nil, err
	}

	return schedule, nil
}

func (s *Schedule) UpdateSchedule(ctx context.Context, schedule_id int64, updateFields map[string]any) (*models.Schedule, error) {
	const op = "schedule.UpdateSchedule"
	log := s.log.With("op", op)

	schedule, err := s.scheduleManager.UpdateSchedule(ctx, schedule_id, updateFields)
	if err != nil {
		log.Error("failed to update schedule", sl.Err(err))
		return nil, err
	}

	// Обновление задачи планировщика
	//s.jobScheduler.UpdateScheduler([]*models.Schedule{schedule})

	return schedule, nil
}

func (s *Schedule) DeleteSchedule(ctx context.Context, schedule_id int64) error {
	const op = "schedule.DeleteSchedule"
	log := s.log.With("op", op)

	_, err := s.scheduleManager.DeleteSchedule(ctx, schedule_id)
	if err != nil {
		log.Error("failed to delete schedule", sl.Err(err))
		return err
	}

	// Удаление задачи планировщика
	//s.jobScheduler.UpdateScheduler([]*models.Schedule{schedule})

	return nil
}

func (s *Schedule) ListSchedule(ctx context.Context, admin_id int64) ([]*models.Schedule, error) {
	const op = "schedule.ListSchedules"
	log := s.log.With("op", op)

	schedules, err := s.scheduleManager.ListSchedule(ctx, admin_id)
	if err != nil {
		log.Error("failed to list schedules", sl.Err(err))
		return nil, err
	}

	return schedules, nil
}
