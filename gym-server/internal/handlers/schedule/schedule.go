package schedule

import (
	"context"
	"health/internal/domain/models"
	pb "health/protogen/v1/schedule"

	"github.com/robfig/cron/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ScheduleService interface {
	CreateSchedules(ctx context.Context, schedules []*models.Schedule) ([]*models.Schedule, error)
}

type ScheduleServerManagmentApi struct {
	pb.UnimplementedScheduleServiceServer
	schedule ScheduleService
}

func RegisterGRPCServer(gRPC *grpc.Server, schedule ScheduleService) {
	pb.RegisterScheduleServiceServer(gRPC, &ScheduleServerManagmentApi{schedule: schedule})
}

func (s *ScheduleServerManagmentApi) CreateSchedules(ctx context.Context, r *pb.CreateSchedulesRequest) (*pb.CreateSchedulesResponse, error) {
	if len(r.Schedules) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "schedules is required")
	}

	modelsSchedules, err := makeModelsSchedules(r.Schedules)
	if err != nil {
		return nil, err
	}

	savedSchedule, err := s.schedule.CreateSchedules(ctx, modelsSchedules)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save schedule: %v", err)
	}

	pbSchedules := makePbSchedules(savedSchedule)

	return &pb.CreateSchedulesResponse{
		Schedules: pbSchedules,
	}, nil
}

func makeModelsSchedules(pbSchedules []*pb.Schedule) ([]*models.Schedule, error) {
	modelsSchedules := make([]*models.Schedule, 0, len(pbSchedules))

	for _, s := range pbSchedules {
		if s.CronExpression == "" {
			return nil, status.Errorf(codes.InvalidArgument, "cron_expression is required")
		}

		if s.VideoId == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "video_id is required")
		}

		_, err := cron.ParseStandard(s.CronExpression)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cron_expression: %v", s.CronExpression)
		}

		modelsSchedules = append(modelsSchedules, &models.Schedule{
			CronExpression: s.CronExpression,
			IsActive:       s.IsActive,
			VideoID:        s.VideoId,
		})
	}

	return modelsSchedules, nil
}

func makePbSchedules(savedSchedule []*models.Schedule) []*pb.Schedule {
	responseSchedules := make([]*pb.Schedule, 0, len(savedSchedule))

	for _, schedule := range savedSchedule {
		responseSchedules = append(responseSchedules, &pb.Schedule{
			Id:             schedule.ID,
			CronExpression: schedule.CronExpression,
			IsActive:       schedule.IsActive,
			VideoId:        schedule.VideoID,
			CreatedAt:      timestamppb.New(schedule.CreatedAt),
		})
	}

	return responseSchedules
}
