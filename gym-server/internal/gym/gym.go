package gym

import (
	"log/slog"
	"time"

	grpcgym "health/internal/gym/grpc"
)

type Gym struct {
	GRPCSrv *grpcgym.Gym
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *Gym {
	grpcGym := grpcgym.New(log, grpcPort)

	return &Gym{
		GRPCSrv: grpcGym,
	}
}
