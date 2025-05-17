package media

import (
	"context"
	"health/internal/domain/models"
	"health/internal/storage/minio"
	"health/lib/logger/sl"
	"log/slog"
	"time"
)

type Media struct {
	log           *slog.Logger
	mediaProvider MediaProvider
	minioStorage  *minio.Storage
}

type MediaProvider interface {
	Media(ctx context.Context, media_id int64) (*models.Media, error)
}

func New(log *slog.Logger, mediaProvider MediaProvider, minioStorage *minio.Storage) *Media {
	return &Media{
		log:           log,
		mediaProvider: mediaProvider,
		minioStorage:  minioStorage,
	}
}

func (m *Media) GetMedia(ctx context.Context, media_id int64, expiryDuration time.Duration) (*models.Media, error) {
	const op = "media.GetMedia"
	log := m.log.With("op", op)

	media, err := m.mediaProvider.Media(ctx, media_id)
	if err != nil {
		log.Error("failed to get media", sl.Err(err))
		return nil, err
	}

	pressignedUrl, err := m.minioStorage.CreatePressignedUrl(ctx, media.ID, expiryDuration)
	if err != nil {
		log.Error("failed to create presigned url", sl.Err(err))
        return nil, err
	}

	return &models.Media{
		ID:            media.ID,
		PressignedUrl: pressignedUrl,
		AdminID:       media.AdminID,
		DepartmentID:  media.DepartmentID,
		CreatedAt:     media.CreatedAt,
	}, nil
}
