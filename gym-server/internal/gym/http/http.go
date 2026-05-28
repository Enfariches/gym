package gymhttp

import (
	"context"
	"fmt"
	"health/internal/services/media"
	"health/internal/services/statistics"
	"health/internal/storage/minio"
	"health/internal/storage/postgres"
	"health/lib/jwt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type GymHTTP struct {
	log          *slog.Logger
	httpServer   *http.Server
	pgStorage    *postgres.Storage
	minioStorage *minio.Storage
	router       *chi.Mux
	port         int
}

func New(log *slog.Logger, pgStorage *postgres.Storage, minioStorage *minio.Storage, port int) *GymHTTP {
	router := chi.NewRouter()

	return &GymHTTP{
		log:          log,
		httpServer:   &http.Server{Addr: ":3000", Handler: router},
		pgStorage:    pgStorage,
		minioStorage: minioStorage,
		router:       router,
		port:         port,
	}
}

func (g *GymHTTP) Run() error {
	const op = "gymhttp.Run"

	log := g.log.With(slog.String("op", op))

	g.router.Use(middleware.Logger)
	g.router.Use(middleware.Recoverer)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions, http.MethodPatch},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})
	g.router.Use(c.Handler)

	g.router.Group(func(r chi.Router) {
		r.Use(jwt.JWTMiddleware)
		r.Post("/api/upload", media.UploadMedia(log, g.pgStorage, g.minioStorage))
		r.Get("/api/export", statistics.ExportStatistics(log, g.pgStorage))
	})

	log.Info("http server is running", slog.String("addr", g.httpServer.Addr))

	if err := g.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (g *GymHTTP) Stop() {
	const op = "gymhttp.Stop"

	g.log.With(slog.String("op", op)).Info("Gracefully stopping http server", slog.Int("port", g.port))

	g.httpServer.Shutdown(context.Background())
}
