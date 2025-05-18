package minio

import (
	"context"
	"fmt"
	"health/internal/config"
	"health/lib/ctxkey"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	log        *slog.Logger
	client     *minio.Client
	bucketName string
}

func NewMinioClient(log *slog.Logger, cfg config.MinioConfig) (*Storage, error) {
	// Инициализация клиента
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	ctx := context.Background()
	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	// Создание bucket если не существует
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{
			Region: cfg.Location,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Storage{log: log, client: client, bucketName: cfg.BucketName}, nil
}

func (ms *Storage) SaveMediaMinio(ctx context.Context, media multipart.File, media_id, department_id, contentType string, size int64) (int64, error) {
	const op = "minio.UploadFile"

	admin_id := ctx.Value(ctxkey.UserKey).(int64)
	admin_id_str := fmt.Sprintf("%d", admin_id)

	log := ms.log.With("op", op)
	log.Info("uploading file")

	m := map[string]string{"department_id": department_id, "admin_id": admin_id_str}
	info, err := ms.client.PutObject(ctx, ms.bucketName, media_id, media, size, minio.PutObjectOptions{ContentType: contentType, UserMetadata: m})
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Info("uploaded file", slog.Int64("size", info.Size))
	return info.Size, nil
}

func (ms *Storage) CreatePressignedUrl(ctx context.Context, media_id int64, expriyDuration time.Duration) (string, error) {
	const op = "minio.CreatePresignedUrl"

	media_id_str := fmt.Sprintf("%d", media_id)

	presignedURL, err := ms.client.PresignedGetObject(
		ctx,
		ms.bucketName,
		media_id_str,
		expriyDuration,
		nil, // дополнительные параметры запроса
	)
	//externalURL := strings.Replace(presignedURL.String(), "http://minio:9000", "http://localhost:9000", 1)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return presignedURL.String(), nil
}
