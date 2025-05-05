package gymgrpc

import (
	"context"
	"fmt"
	admingrpc "health/internal/handlers/admin"
	authgrpc "health/internal/handlers/auth"
	"health/lib/jwt"
	ctxkey "health/lib/ctxkey"
	"log/slog"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var PublicMethods = map[string]bool{
	"/auth.AuthService/Login":                true,
	"/auth.AuthService/Register":             true,
	"/auth.AuthService/ChangePassword":       true,
	"/auth.AuthService/VerifyChangePassword": true,
	"/auth.AuthService/VerifyRegister":       true,
}

type Gym struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.AuthService, adminService admingrpc.AdminService, port int) *Gym {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(JWTServerInterceptor),
	)

	authgrpc.RegisterGRPCServer(gRPCServer, authService)
	admingrpc.RegisterGRPCServer(gRPCServer, adminService)

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

func JWTServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if PublicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	authHeader := md.Get("Authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

	userId, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	ctx = context.WithValue(ctx, ctxkey.UserKey, userId)

	return handler(ctx, req)
}
