package file

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-server/internal/config"
	"go-server/internal/errcode"
)

type localStorage struct {
	baseDir    string
	baseURL    string
	dateLayout string
}

func newLocalStorage(cfg config.FileConfig) Storage {
	baseDir := strings.TrimSpace(cfg.Local.BaseDir)
	if baseDir == "" {
		baseDir = "uploads"
	}
	baseURL := strings.TrimSpace(cfg.Local.BaseURL)
	if baseURL == "" {
		baseURL = "/uploads"
	}
	dateLayout := strings.TrimSpace(cfg.Local.DateLayout)
	if dateLayout == "" {
		dateLayout = "2006/01/02"
	}
	return &localStorage{
		baseDir:    baseDir,
		baseURL:    baseURL,
		dateLayout: dateLayout,
	}
}

func (s *localStorage) Save(_ context.Context, fileHeader *multipart.FileHeader, category string) (StoredFile, error) {
	category = normalizeCategory(category, "default")
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	objectName := newObjectName(ext)

	dateDir := time.Now().UTC().Format(s.dateLayout)
	relativeDir := filepath.Join(category, filepath.FromSlash(dateDir))
	fullDir := filepath.Join(s.baseDir, relativeDir)
	if err := os.MkdirAll(fullDir, 0o755); err != nil {
		return StoredFile{}, errcode.ErrInternalError.AsError()
	}

	relativePath := filepath.Join(relativeDir, objectName)
	fullPath := filepath.Join(s.baseDir, relativePath)

	src, err := fileHeader.Open()
	if err != nil {
		return StoredFile{}, errcode.ErrInternalError.AsError()
	}
	defer func() { _ = src.Close() }()

	dst, err := os.Create(fullPath)
	if err != nil {
		return StoredFile{}, errcode.ErrInternalError.AsError()
	}
	defer func() { _ = dst.Close() }()

	if _, err := io.Copy(dst, src); err != nil {
		return StoredFile{}, errcode.ErrInternalError.AsError()
	}

	return StoredFile{
		ObjectName: objectName,
		Ext:        ext,
		Size:       fileHeader.Size,
		Path:       filepath.ToSlash(relativePath),
		URL:        localPublicURL(s.baseURL, relativePath),
	}, nil
}

func (s *localStorage) Delete(_ context.Context, item File) error {
	if strings.TrimSpace(item.Path) == "" {
		return nil
	}
	target := filepath.Join(s.baseDir, filepath.FromSlash(item.Path))
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return errcode.ErrInternalError.AsError()
	}
	return nil
}

func (s *localStorage) Name() string {
	return "local"
}
