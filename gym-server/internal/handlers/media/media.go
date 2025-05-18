package media

import (
	"context"
	"health/internal/domain/models"
	pb "health/protogen/v1/media"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MediaService interface {
	GetMedia(ctx context.Context, media_id int64, expiryDuration time.Duration) (*models.Media, error)
}

type MediaServerManagmentApi struct {
	pb.UnimplementedMediaServiceServer
	media MediaService
}

func RegisterGRPCServer(gRPC *grpc.Server, media MediaService) {
	pb.RegisterMediaServiceServer(gRPC, &MediaServerManagmentApi{media: media})
}

func (s *MediaServerManagmentApi) GetMedia(ctx context.Context, r *pb.GetMediaRequest) (*pb.Media, error) {
	if r.DepartmentId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "department_id is required")
	}
	if r.MediaId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "media_id is required")
	}

	expiryDuration := r.Expiry.AsDuration()

	media, err := s.media.GetMedia(ctx, r.MediaId, expiryDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get media: %v", err)
	}

	return &pb.Media{
		Id:            media.ID,
		PressignedUrl: media.PressignedUrl,
		DepartmentId: media.DepartmentID,
		CreatedAt:     media.CreatedAt.Format(time.RFC3339),
	}, nil
}
