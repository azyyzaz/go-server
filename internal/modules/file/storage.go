package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-server/internal/config"
	"go-server/internal/errcode"

	"github.com/google/uuid"
)

type Storage interface {
	Save(ctx context.Context, fileHeader *multipart.FileHeader, category string) (StoredFile, error)
	Delete(ctx context.Context, item File) error
	Name() string
}

type StoredFile struct {
	ObjectName string
	Ext        string
	MIMEType   string
	Size       int64
	Bucket     string
	Path       string
	URL        string
}

type Validator struct {
	maxSizeBytes int64
	allowedExts  map[string]struct{}
}

func NewStorage(cfg config.FileConfig) (Storage, error) {
	switch normalizeStorage(cfg.Storage) {
	case "local":
		return newLocalStorage(cfg), nil
	case "minio":
		return newMinIOStorage(cfg)
	default:
		return nil, errcode.ErrFileStorageUnsupported.AsError()
	}
}

func NewValidator(maxSizeMB int64, allowedExts []string) Validator {
	if maxSizeMB <= 0 {
		maxSizeMB = 10
	}
	set := make(map[string]struct{}, len(allowedExts))
	for _, ext := range allowedExts {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		set[ext] = struct{}{}
	}
	return Validator{
		maxSizeBytes: maxSizeMB * 1024 * 1024,
		allowedExts:  set,
	}
}

func (v Validator) Validate(fileHeader *multipart.FileHeader) (PreparedFile, error) {
	if fileHeader == nil || fileHeader.Size <= 0 {
		return PreparedFile{}, errcode.ErrInvalidParam.AsError()
	}
	if v.maxSizeBytes > 0 && fileHeader.Size > v.maxSizeBytes {
		return PreparedFile{}, errcode.ErrFileTooLarge.AsError()
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if len(v.allowedExts) > 0 {
		if _, ok := v.allowedExts[ext]; !ok {
			return PreparedFile{}, errcode.ErrFileTypeNotAllowed.AsError()
		}
	}

	src, err := fileHeader.Open()
	if err != nil {
		return PreparedFile{}, errcode.ErrInternalError.AsError()
	}
	defer func() { _ = src.Close() }()

	head := make([]byte, 512)
	n, err := src.Read(head)
	if err != nil && err != io.EOF {
		return PreparedFile{}, errcode.ErrInternalError.AsError()
	}

	return PreparedFile{
		Ext:      ext,
		MIMEType: http.DetectContentType(head[:n]),
	}, nil
}

type PreparedFile struct {
	Ext      string
	MIMEType string
}

func normalizeStorage(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "local"
	}
	return raw
}

func normalizeCategory(raw string, fallback string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return fallback
	}
	return raw
}

func newObjectName(ext string) string {
	return strings.ReplaceAll(uuid.NewString(), "-", "") + ext
}

func parseTimeRange(startRaw, endRaw string) (*time.Time, *time.Time, error) {
	parse := func(raw string) (*time.Time, error) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return nil, nil
		}
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return nil, ErrFileInvalidTimeRange
		}
		utc := t.UTC()
		return &utc, nil
	}

	start, err := parse(startRaw)
	if err != nil {
		return nil, nil, err
	}
	end, err := parse(endRaw)
	if err != nil {
		return nil, nil, err
	}
	if start != nil && end != nil && start.After(*end) {
		return nil, nil, ErrFileInvalidTimeRange
	}
	return start, end, nil
}

func localPublicURL(baseURL, relativePath string) string {
	baseURL = "/" + strings.Trim(strings.TrimSpace(baseURL), "/")
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(filepath.ToSlash(relativePath), "/"))
}
