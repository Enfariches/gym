package suite

import (
	"context"
	"health/internal/config"
	"health/internal/storage/postgres"
	authpb "health/protogen/v1/auth"
	mediapb "health/protogen/v1/media"
	schedulepb "health/protogen/v1/schedule"
	statspb "health/protogen/v1/statistics"
	userpb "health/protogen/v1/users"
	"log/slog"
	"net"
	"os"
	"strconv"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T // Потребуется для вызова методов *testing.T внутри Suite
	Cfg        *config.Config
	Storage    *postgres.Storage // Конфигурация приложения

	AuthClient       authpb.AuthServiceClient
	AdminClient      userpb.AdminServiceClient
	ScheduleClient   schedulepb.ScheduleServiceClient
	MediaClient      mediapb.MediaServiceClient
	EmployeeClient   userpb.EmployeeServiceClient
	StatisticsClient statspb.StatisticsServiceClient
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
//
// TODO: for pipeline tests we need to wait for app is ready
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("../config/dev.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	lg := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	pg, err := postgres.New(lg, "postgres://postgres:root@localhost:5436/health?sslmode=disable")
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials())) // Используем insecure-коннект для тестов
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:       t,
		Cfg:     cfg,
		Storage: pg,

		AuthClient:       authpb.NewAuthServiceClient(cc),
		AdminClient:      userpb.NewAdminServiceClient(cc),
		ScheduleClient:   schedulepb.NewScheduleServiceClient(cc),
		MediaClient:      mediapb.NewMediaServiceClient(cc),
		EmployeeClient:   userpb.NewEmployeeServiceClient(cc),
		StatisticsClient: statspb.NewStatisticsServiceClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
