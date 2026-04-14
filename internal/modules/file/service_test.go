package file

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"go-server/internal/config"
)

func TestServiceUploadListAndDelete(t *testing.T) {
	baseDir := t.TempDir()
	storage := newLocalStorage(config.FileConfig{
		Local: config.LocalFileConfig{
			BaseDir:    baseDir,
			BaseURL:    "/uploads",
			DateLayout: "2006/01/02",
		},
	})
	svc := NewService(
		NewInMemoryRepository(),
		storage,
		NewValidator(10, []string{".png", ".pdf"}),
		NewValidator(2, []string{".png"}),
	)

	fileHeader := makeFileHeader(t, "hello.png", []byte("fake png content"))
	uploaderID := uint(7)

	uploaded, err := svc.Upload(context.Background(), &uploaderID, fileHeader, UploadRequest{Category: "image"})
	if err != nil {
		t.Fatalf("upload file: %v", err)
	}
	if uploaded.ID == 0 {
		t.Fatalf("expected persisted file id")
	}
	if uploaded.Category != "image" {
		t.Fatalf("unexpected category: %s", uploaded.Category)
	}
	if uploaded.URL == "" {
		t.Fatalf("expected file url")
	}

	listed, err := svc.List(context.Background(), ListFilesQuery{Page: 1, PageSize: 10, Category: "image"})
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if listed.Total != 1 || len(listed.Items) != 1 {
		t.Fatalf("unexpected list result: %+v", listed)
	}

	relativePath := filepath.Join(baseDir, filepath.FromSlash(listed.Items[0].URL[len("/uploads/"):]))
	if _, err := os.Stat(relativePath); err != nil {
		t.Fatalf("stored file missing: %v", err)
	}

	if err := svc.Delete(context.Background(), uploaded.ID); err != nil {
		t.Fatalf("delete file: %v", err)
	}
	if _, err := os.Stat(relativePath); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, got: %v", err)
	}
}

func TestServiceUploadAvatarRejectsNonImage(t *testing.T) {
	baseDir := t.TempDir()
	storage := newLocalStorage(config.FileConfig{
		Local: config.LocalFileConfig{
			BaseDir:    baseDir,
			BaseURL:    "/uploads",
			DateLayout: "2006/01/02",
		},
	})
	svc := NewService(
		NewInMemoryRepository(),
		storage,
		NewValidator(10, []string{".png", ".pdf"}),
		NewValidator(2, []string{".png"}),
	)

	fileHeader := makeFileHeader(t, "not-allowed.pdf", []byte("%PDF-1.7"))
	if _, err := svc.UploadAvatar(context.Background(), 1, fileHeader); err == nil {
		t.Fatalf("expected avatar upload to fail")
	}
}

func makeFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(len(body.Bytes()) + 1024)); err != nil {
		t.Fatalf("parse multipart form: %v", err)
	}

	fileHeader := req.MultipartForm.File["file"][0]
	fileHeader.Size = int64(len(content))
	return fileHeader
}
