package gym

import (
	"log/slog"
	"time"

	"health/internal/config"
	grpcgym "health/internal/gym/grpc"
	"health/internal/services/auth"
	"health/internal/storage/postgres"
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

	authService := auth.New(log, storage, storage, storage, smtpConfig, tokenTTL, authTokenTLL)

	grpcGym := grpcgym.New(log, authService, grpcPort)

	return &Gym{
		GRPCSrv: grpcGym,
	}
}
