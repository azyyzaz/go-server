package file

import (
	"context"
	"mime/multipart"
	"net/url"
	"path"
	"strings"
	"time"

	"go-server/internal/config"
	"go-server/internal/errcode"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client   *minio.Client
	bucket   string
	baseURL  string
	useSSL   bool
	endpoint string
}

func newMinIOStorage(cfg config.FileConfig) (Storage, error) {
	if !cfg.MinIO.Enabled || strings.TrimSpace(cfg.MinIO.Endpoint) == "" || strings.TrimSpace(cfg.MinIO.Bucket) == "" {
		return nil, errcode.ErrFileStorageUnsupported.AsError()
	}

	client, err := minio.New(strings.TrimSpace(cfg.MinIO.Endpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(strings.TrimSpace(cfg.MinIO.AccessKey), strings.TrimSpace(cfg.MinIO.SecretKey), ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, errcode.ErrFileStorageUnsupported.AsError()
	}

	return &minioStorage{
		client:   client,
		bucket:   strings.TrimSpace(cfg.MinIO.Bucket),
		baseURL:  strings.TrimSpace(cfg.MinIO.BaseURL),
		useSSL:   cfg.MinIO.UseSSL,
		endpoint: strings.TrimSpace(cfg.MinIO.Endpoint),
	}, nil
}

func (s *minioStorage) Save(ctx context.Context, fileHeader *multipart.FileHeader, category string) (StoredFile, error) {
	category = normalizeCategory(category, "default")
	ext := strings.ToLower(path.Ext(fileHeader.Filename))
	objectName := newObjectName(ext)
	objectPath := path.Join(category, time.Now().UTC().Format("2006/01/02"), objectName)

	src, err := fileHeader.Open()
	if err != nil {
		return StoredFile{}, errcode.ErrInternalError.AsError()
	}
	defer func() { _ = src.Close() }()

	_, err = s.client.PutObject(ctx, s.bucket, objectPath, src, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return StoredFile{}, errcode.ErrInternalError.AsError()
	}

	return StoredFile{
		ObjectName: objectName,
		Ext:        ext,
		Size:       fileHeader.Size,
		Bucket:     s.bucket,
		Path:       objectPath,
		URL:        s.publicURL(objectPath),
	}, nil
}

func (s *minioStorage) Delete(ctx context.Context, item File) error {
	if strings.TrimSpace(item.Path) == "" {
		return nil
	}
	if err := s.client.RemoveObject(ctx, s.bucket, item.Path, minio.RemoveObjectOptions{}); err != nil {
		return errcode.ErrInternalError.AsError()
	}
	return nil
}

func (s *minioStorage) Name() string {
	return "minio"
}

func (s *minioStorage) publicURL(objectPath string) string {
	if s.baseURL != "" {
		base := strings.TrimRight(s.baseURL, "/")
		return base + "/" + strings.TrimLeft(objectPath, "/")
	}

	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return (&url.URL{
		Scheme: scheme,
		Host:   s.endpoint,
		Path:   path.Join(s.bucket, objectPath),
	}).String()
}
