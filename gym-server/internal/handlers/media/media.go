package media

import (
	"context"
	"health/internal/domain/models"
	"health/lib/ctxkey"
	pb "health/protogen/v1/media"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MediaService interface {
	GetMedia(ctx context.Context, media_id int64, expiryDuration time.Duration) (*models.Media, error)
	ListMedia(ctx context.Context, admin_id int64) ([]*models.Media, error)
	DeleteMedia(ctx context.Context, media_id int64) error
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
		Title:         media.Title,
		PressignedUrl: media.PressignedUrl,
		AdminId:       media.AdminID,
		DepartmentId:  media.DepartmentID,
		CreatedAt:     media.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MediaServerManagmentApi) ListMedia(ctx context.Context, _ *emptypb.Empty) (*pb.ListMediaResponse, error) {
	admin_id := ctx.Value(ctxkey.UserKey).(int64)

	media, err := s.media.ListMedia(ctx, admin_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list media: %v", err)
	}

	mediasPb := makePbMedias(media)

	return &pb.ListMediaResponse{
		Medias: mediasPb,
	}, nil
}

func (s *MediaServerManagmentApi) DeleteMedia(ctx context.Context, r *pb.DeleteMediaRequest) (*emptypb.Empty, error) {
	if r.MediaId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "media_id is required")
	}

	err := s.media.DeleteMedia(ctx, r.MediaId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete media: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func makePbMedias(savedMedia []*models.Media) []*pb.Media {
	responseMedias := make([]*pb.Media, 0, len(savedMedia))

	for _, m := range savedMedia {
		responseMedias = append(responseMedias, &pb.Media{
			Id:    m.ID,
			Title: m.Title,
			// Выставляем ссылку на временный объект, который мы храним в БД
			// Ссылка всегда пустая, так как мы не храним в БД временные ссылки для передачи клиенту
			PressignedUrl: "",
			AdminId:       m.AdminID,
			DepartmentId:  m.DepartmentID,
			CreatedAt:     m.CreatedAt.Format(time.RFC3339),
		})
	}

	return responseMedias
}
