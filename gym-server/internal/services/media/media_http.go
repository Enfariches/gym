package media

import (
	"context"
	"health/lib/ctxkey"
	"health/lib/logger/sl"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type MediaSaverPostgres interface {
	SaveMediaPostgres(admin_id, department_id int64) (int64, error)
}

type MediaSaverMinio interface {
	SaveMediaMinio(ctx context.Context, media multipart.File, media_id, department_id, contentType string, size int64) (int64, error)
}

// Максимальный размер файла 200Мб
const maxFileSize = 200 << 20

// HTTP-функция для загрузки видео на сервер
func UploadMedia(log *slog.Logger, mediaSaverPg MediaSaverPostgres, mediaSaverMinio MediaSaverMinio) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "mediaHttp.UploadMedia"

		log = log.With(slog.String("op", op))
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

		media, fileHeader, err := r.FormFile("mediafile")
		if err != nil {
			if strings.Contains(err.Error(), "request body too large") {
				log.Error("mediafile > 200MB", sl.Err(err))
				http.Error(w, "mediafile too large (max 200MB)", http.StatusBadRequest)
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := fileHeader.Header.Get("Content-Type")

		if contentType != "video/mp4" {
			log.Error("unsupported mediafile type")
			http.Error(w, "unsupported mediafile type", http.StatusBadRequest)
		}

		size, err := getFileSize(media)
		if err != nil {
			log.Error("failed to get mediafile size", sl.Err(err))
			http.Error(w, "failed to get mediafile size", http.StatusInternalServerError)
			return
		}

		defer media.Close()

		ctx := r.Context()
		admin_id := ctx.Value(ctxkey.UserKey).(int64)
		department_id := ctx.Value(ctxkey.DepartmentKey).(int64)

		media_id, err := mediaSaverPg.SaveMediaPostgres(admin_id, department_id)
		if err != nil {
			log.Error("failed to save mediafile to postgres", sl.Err(err))
			http.Error(w, "failed to save mediafile to postgres", http.StatusInternalServerError)
			return
		}

		mediaId := strconv.FormatInt(media_id, 10)
		departmentId := strconv.FormatInt(department_id, 10)

		mediaSize, err := mediaSaverMinio.SaveMediaMinio(ctx, media, mediaId, departmentId, contentType, size)
		if err != nil {
			log.Error("failed to save mediafile to minio", sl.Err(err))
			http.Error(w, "failed to save mediafile to minio", http.StatusInternalServerError)
			return
		}

		log.Info("mediafile uploaded successfully",
			slog.String("filename", fileHeader.Filename),
			slog.Int64("size", mediaSize),
		)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mediafile uploaded successfully"))

	}
}

func getFileSize(file multipart.File) (int64, error) {
	// Перемещаем указатель в конец
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Возвращаем указатель в начало
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	return size, nil
}
