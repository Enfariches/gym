package gym

import (
	"log/slog"
	"time"

	"health/internal/config"
	grpcgym "health/internal/gym/grpc"
	httpgym "health/internal/gym/http"
	"health/internal/services/admin"
	"health/internal/services/auth"
	"health/internal/services/media"
	"health/internal/services/schedule"
	"health/internal/storage/minio"
	"health/internal/storage/postgres"
	"health/lib/gcron"
	"health/lib/logger/sl"
)

type Gym struct {
	GRPCSrv *grpcgym.GymGRPC
	HTTPSrv *httpgym.GymHTTP
}

func New(log *slog.Logger, grpcPort, httpPort int, storagePath string, minioConfig config.MinioConfig,
	smtpConfig config.SMTPConfig, tokenTTL, authTokenTLL time.Duration) *Gym {

	pgStorage, err := postgres.New(log, storagePath)
	if err != nil {
		log.Error("failed to create storage", slog.String("path", storagePath))
		return nil
	}

	minioStorage, err := minio.NewMinioClient(log, minioConfig)
	if err != nil {
		log.Error("failed to create minio storage", sl.Err(err))
		return nil
	}

	jobScheduler, err := gcron.NewJobScheduler()
	if err != nil {
		log.Error("failed to create job scheduler")
		return nil
	}

	authService := auth.New(log, pgStorage, smtpConfig, tokenTTL, authTokenTLL)
	adminService := admin.New(log, pgStorage)
	scheduleService := schedule.New(log, pgStorage, jobScheduler)
	mediaService := media.New(log, pgStorage, minioStorage)

	grpcGym := grpcgym.New(log, authService, adminService, scheduleService, mediaService, grpcPort)
	httpGym := httpgym.New(log, pgStorage, minioStorage, httpPort)

	return &Gym{
		GRPCSrv: grpcGym,
		HTTPSrv: httpGym,
	}
}
