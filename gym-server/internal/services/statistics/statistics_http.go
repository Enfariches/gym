package statistics

import (
	"context"
	"health/internal/domain/models"
	"health/lib/ctxkey"
	"health/lib/logger/sl"
	"health/lib/pdf"
	"log/slog"
	"net/http"
)

type PgStatsProvider interface {
	ListDepartmentStatistics(ctx context.Context, department_id int64) ([]*models.Statistics, error)

	GetAdminInfo(ctx context.Context, admin_id int64) (string, string, error)

	GetDepNameByDepID(department_id int64) (string, error)
}

func ExportStatistics(log *slog.Logger, pgProvider PgStatsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "mediaHttp.ExportStatistics"

		log = log.With(slog.String("op", op))
		ctx := r.Context()

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=gymnastics_report.pdf")

		department_id := ctx.Value(ctxkey.DepartmentKey).(int64)
		admin_id := ctx.Value(ctxkey.UserKey).(int64)

		departmentName, err := pgProvider.GetDepNameByDepID(department_id)
		if err != nil {
			log.Error("failed to get department name", sl.Err(err))
			http.Error(w, "failed to get department name", http.StatusInternalServerError)
			return
		}

		adminName, adminPosition, err := pgProvider.GetAdminInfo(r.Context(), admin_id)
		if err != nil {
			log.Error("failed to get admin info", sl.Err(err))
			http.Error(w, "failed to get admin info", http.StatusInternalServerError)
			return
		}

		stats, err := pgProvider.ListDepartmentStatistics(ctx, department_id)

		pdfFile := pdf.ExportToPDF(stats, adminName, adminPosition, departmentName)

		err = pdfFile.Output(w)
		if err != nil {
			log.Error("failed gen pdf", sl.Err(err))
			http.Error(w, "failed gen pdf", http.StatusInternalServerError)
			return
		}
	}
}
