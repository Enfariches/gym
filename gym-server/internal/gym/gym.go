package gym

import (
	"log/slog"
	"time"

	"health/internal/config"
	grpcgym "health/internal/gym/grpc"
	"health/internal/services/admin"
	"health/internal/services/auth"
	"health/internal/services/schedule"
	"health/internal/storage/postgres"
	"health/lib/gcron"
)

type Gym struct {
	GRPCSrv *grpcgym.Gym
}

func New(log *slog.Logger, grpcPort int, storagePath string,
	smtpConfig config.SMTPConfig, tokenTTL, authTokenTLL time.Duration) *Gym {

	storage, err := postgres.New(log, storagePath)
	if err != nil {
		log.Error("failed to create storage", slog.String("path", storagePath))
	}

	jobScheduler, err := gcron.NewJobScheduler()
	if err != nil {
		log.Error("failed to create job scheduler")
	}

	authService := auth.New(log, storage, storage, smtpConfig, tokenTTL, authTokenTLL)
	adminService := admin.New(log, storage)
	scheduleService := schedule.New(log, storage, jobScheduler)

	grpcGym := grpcgym.New(log, authService, adminService, scheduleService, grpcPort)

	return &Gym{
		GRPCSrv: grpcGym,
	}
}
