package main

import (
	"health/internal/config"
	"health/internal/gym"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envDev  = "envDev"
	envProd = "envProd"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg

	log := setupLogger(envDev)
	log.Info("Starting server...")

	app := gym.New(log, cfg.GRPC.Port, cfg.Storage, cfg.TokenTTL)
	go app.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.GRPCSrv.Stop()

	// Приложение
	// grpc-сервер
	// Тесты
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
