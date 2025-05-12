package schedule

import (
	"context"
	"health/internal/domain/models"
	"health/lib/ctxkey"
	pb "health/protogen/v1/schedule"

	"github.com/mennanov/fmutils"
	"github.com/robfig/cron/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ScheduleService interface {
	CreateSchedules(ctx context.Context, schedules []*models.Schedule) ([]*models.Schedule, error)
	GetSchedule(ctx context.Context, schedule_id int64) (*models.Schedule, error)
	UpdateSchedule(ctx context.Context, schedule_id int64, updateFields map[string]any) (*models.Schedule, error)
	DeleteSchedule(ctx context.Context, schedule_id int64) error
	ListSchedule(ctx context.Context, admin_id int64) ([]*models.Schedule, error)
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

func (s *ScheduleServerManagmentApi) GetSchedule(ctx context.Context, r *pb.GetScheduleRequest) (*pb.Schedule, error) {
	if r.ScheduleId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "schedule_id is required")
	}

	schedule, err := s.schedule.GetSchedule(ctx, r.ScheduleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get schedule: %v", err)
	}

	return &pb.Schedule{
		Id:             schedule.ID,
		CronExpression: schedule.CronExpression,
		IsActive:       schedule.IsActive,
		VideoId:        schedule.VideoID,
		AdminId:        schedule.AdminID,
		CreatedAt:      timestamppb.New(schedule.CreatedAt),
	}, nil
}

func (s *ScheduleServerManagmentApi) UpdateSchedule(ctx context.Context, r *pb.UpdateScheduleRequest) (*pb.Schedule, error) {
	if r.Schedule.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "schedule_id is required")
	}
	schedule_id := r.Schedule.Id

	r.FieldMask.Normalize()

	if !r.FieldMask.IsValid(r.Schedule) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	fmutils.Filter(r.GetSchedule(), r.FieldMask.GetPaths())

	updateFields, err := applyFieldMask(r.Schedule, r.FieldMask)
	if err != nil {
		return nil, err
	}

	updatedSchedule, err := s.schedule.UpdateSchedule(ctx, schedule_id, updateFields)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update schedule: %v", err)
	}

	return &pb.Schedule{
		Id:             updatedSchedule.ID,
		CronExpression: updatedSchedule.CronExpression,
		IsActive:       updatedSchedule.IsActive,
		VideoId:        updatedSchedule.VideoID,
		AdminId:        updatedSchedule.AdminID,
		CreatedAt:      timestamppb.New(updatedSchedule.CreatedAt),
	}, nil
}

func (s *ScheduleServerManagmentApi) DeleteSchedule(ctx context.Context, r *pb.DeleteScheduleRequest) (*emptypb.Empty, error) {
	if r.ScheduleId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "schedule_id is required")
	}

	err := s.schedule.DeleteSchedule(ctx, r.ScheduleId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete schedule: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ScheduleServerManagmentApi) ListSchedule(ctx context.Context, _ *emptypb.Empty) (*pb.ListScheduleResponse, error) {
	admin_id := ctx.Value(ctxkey.UserKey).(int64)

	schedules, err := s.schedule.ListSchedule(ctx, admin_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list schedules: %v", err)
	}

	schedulesPb := makePbSchedules(schedules)

	return &pb.ListScheduleResponse{
		Schedules: schedulesPb,
	}, nil

}

func applyFieldMask(req *pb.Schedule, mask *fieldmaskpb.FieldMask) (map[string]any, error) {
	updateMap := make(map[string]any)

	for _, path := range mask.GetPaths() {
		switch path {
		case "cron_expression":
			_, err := cron.ParseStandard(req.CronExpression)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid cron_expression: %v", req.CronExpression)
			}
			updateMap[path] = req.CronExpression
		case "is_active":
			updateMap[path] = req.IsActive
		case "video_id":
			updateMap[path] = req.VideoId
		default:
			return nil, status.Errorf(codes.InvalidArgument, "incorrect fields: %s", path)
		}
	}

	for _, value := range updateMap {
		if value == "" {
			return nil, status.Error(codes.InvalidArgument, "field value is required")
		}
	}
	return updateMap, nil
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
			AdminId:        schedule.AdminID,
			CreatedAt:      timestamppb.New(schedule.CreatedAt),
		})
	}

	return responseSchedules
}
