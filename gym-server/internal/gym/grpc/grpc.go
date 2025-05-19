package gymgrpc

import (
	"context"
	"fmt"
	admingrpc "health/internal/handlers/admin"
	authgrpc "health/internal/handlers/auth"
	employeegrpc "health/internal/handlers/employee"
	mediagrpc "health/internal/handlers/media"
	schedulegrpc "health/internal/handlers/schedule"
	statisticsgrpc "health/internal/handlers/statistics"
	"health/lib/jwt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	"google.golang.org/grpc"
)

type GymGRPC struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.AuthService, adminService admingrpc.AdminService,
	scheduleService schedulegrpc.ScheduleService, employeeService employeegrpc.EmployeeService,
	mediaService mediagrpc.MediaService, statisticsService statisticsgrpc.StatisticsService, port int) *GymGRPC {

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			jwt.JWTServerInterceptor,
		),
	)

	authgrpc.RegisterGRPCServer(gRPCServer, authService)
	admingrpc.RegisterGRPCServer(gRPCServer, adminService)
	employeegrpc.RegisterGRPCServer(gRPCServer, employeeService)
	schedulegrpc.RegisterGRPCServer(gRPCServer, scheduleService)
	mediagrpc.RegisterGRPCServer(gRPCServer, mediaService)
	statisticsgrpc.RegisterGRPCServer(gRPCServer, statisticsService)

	return &GymGRPC{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (g *GymGRPC) Run() error {
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

func (g *GymGRPC) Stop() {
	const op = "gymgrpc.Stop"

	g.log.With(slog.String("op", op)).Info("Gracefuly stopping grpc server", slog.Int("port", g.port))

	g.gRPCServer.GracefulStop()
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
