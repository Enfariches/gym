package gymgrpc

import (
	"fmt"
	authgrpc "health/internal/handlers/auth"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Gym struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *Gym {
	gRPCServer := grpc.NewServer()

	authgrpc.RegisterGRPCServer(gRPCServer)

	return &Gym{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (g *Gym) Run() error {
	const op = "gymgrpc.Run"

	log := g.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := g.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (g *Gym) Stop() {
	const op = "gymgrpc.Stop"

	g.log.With(slog.String("op", op)).Info("Gracefuly stopping grpc server", slog.Int("port", g.port))

	g.gRPCServer.GracefulStop()
}
