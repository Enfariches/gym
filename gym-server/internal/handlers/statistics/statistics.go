package statistics

import (
	"context"
	"fmt"
	"health/internal/domain/models"
	pb "health/protogen/v1/statistics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StatisticsService interface {
	CreateStatistics(ctx context.Context, statProgress string, percentView, media_id int64) (error)
	GetEmployeeStatistics(ctx context.Context, employee_id, media_id int64) (*models.Statistics, error)
}

type StatisticsServerManagmentApi struct {
	pb.UnimplementedStatisticsServiceServer
	statistics StatisticsService
}

func RegisterGRPCServer(gRPC *grpc.Server, statistics StatisticsService) {
	pb.RegisterStatisticsServiceServer(gRPC, &StatisticsServerManagmentApi{statistics: statistics})
}

func (s *StatisticsServerManagmentApi) CreateStatistics(ctx context.Context, r *pb.CreateStatisticsRequest) (*emptypb.Empty, error) {

	if r.MediaId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "media_id is required")
	}

	if r.PercentageView < 0 || r.PercentageView > 100 {
		return nil, status.Error(codes.InvalidArgument, "percentage_view must be between 0 and 100")
	}

	if r.PercentageView > 0 && r.PercentageView < 100 && r.Progress != pb.MediaProgress_INCOMPLETE {
		return nil, status.Error(codes.InvalidArgument, "if 0 < percentage_view < 100 progress cant be completed or skipped")
	}

	statProgress, err := progressFromProto(r.Progress)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid media progress: %v", err))
	}

	err = s.statistics.CreateStatistics(ctx, statProgress, r.PercentageView, r.MediaId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create statistics: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (s *StatisticsServerManagmentApi) GetEmployeeStatistics(ctx context.Context, r *pb.GetEmployeeStatisticsRequest) (*pb.Statistic, error) {
	if r.EmployeeId == 0 {
        return nil, status.Errorf(codes.InvalidArgument, "employee_id is required")
    }

	if r.MediaId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "media_id is required")
	}

    stats, err := s.statistics.GetEmployeeStatistics(ctx, r.EmployeeId, r.MediaId)
    if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get employee statistics: %v", err))
    }

	pbProgress, err := progressToProto(stats.Progress)
	if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to parse media progress to pb: %v", err))
    }

    pbStats := &pb.Statistic{
        Id: stats.ID,
		Progress: pbProgress,
		PercentageView: stats.PercentageView,
		EmployeeName: stats.EmployeeName,
		EmployeeSurname: stats.EmployeeSurname,
		MediaTitle: stats.MediaTitle,
		CreatedAt: timestamppb.New(stats.CreatedAt),
    }

    return pbStats, nil
}

func progressFromProto(pbProgress pb.MediaProgress) (string, error) {
	switch pbProgress {
	case pb.MediaProgress_SKIPPED:
		return "skipped", nil
	case pb.MediaProgress_INCOMPLETE:
		return "incomplete", nil
	case pb.MediaProgress_COMPLETED:
		return "complete", nil
	default:
		return "", fmt.Errorf("unknown media progress: %v", pbProgress)
	}
}

func progressToProto(progress string) (pb.MediaProgress, error) {
	switch progress {
	case "skipped":
		return pb.MediaProgress_SKIPPED, nil
	case "incomplete":
		return pb.MediaProgress_INCOMPLETE, nil
	case "complete":
		return pb.MediaProgress_COMPLETED, nil
	default:
		return 0, fmt.Errorf("unknown parse media progress to pb: %v", progress)
	}
}
