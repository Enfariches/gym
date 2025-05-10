package main

import (
	"health/internal/config"
	"health/internal/gym"
	"health/lib/logger/handlers/slogpretty"
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

	app := gym.New(log, cfg.GRPC.Port, cfg.Storage, cfg.SMTP, cfg.TokenTTL, cfg.AuthTokenTTL)
	go app.GRPCSrv.Run()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.GRPCSrv.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
