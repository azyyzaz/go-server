package file

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"go-server/internal/errcode"
)

type Service interface {
	Upload(ctx context.Context, uploaderID *uint, fileHeader *multipart.FileHeader, req UploadRequest) (FileResponse, error)
	UploadAvatar(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader) (FileResponse, error)
	List(ctx context.Context, q ListFilesQuery) (FilePageResult, error)
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo            Repository
	storage         Storage
	validator       Validator
	avatarValidator Validator
}

func NewService(repo Repository, storage Storage, validator Validator, avatarValidator Validator) Service {
	return &service{
		repo:            repo,
		storage:         storage,
		validator:       validator,
		avatarValidator: avatarValidator,
	}
}

func (s *service) Upload(ctx context.Context, uploaderID *uint, fileHeader *multipart.FileHeader, req UploadRequest) (FileResponse, error) {
	return s.save(ctx, uploaderID, fileHeader, normalizeCategory(req.Category, inferCategory(fileHeader)))
}

func (s *service) UploadAvatar(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader) (FileResponse, error) {
	if fileHeader == nil {
		return FileResponse{}, errcode.ErrInvalidParam.AsError()
	}

	prepared, err := s.avatarValidator.Validate(fileHeader)
	if err != nil {
		return FileResponse{}, err
	}
	if !strings.HasPrefix(prepared.MIMEType, "image/") {
		return FileResponse{}, errcode.ErrFileTypeNotAllowed.AsError()
	}
	id := uploaderID
	return s.savePrepared(ctx, &id, fileHeader, "avatar", prepared)
}

func (s *service) List(ctx context.Context, q ListFilesQuery) (FilePageResult, error) {
	normalizeListQuery(&q)
	start, end, err := parseTimeRange(q.StartTime, q.EndTime)
	if err != nil {
		return FilePageResult{}, err
	}
	items, total, err := s.repo.List(ctx, q, start, end)
	if err != nil {
		return FilePageResult{}, err
	}

	result := make([]FileResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toResponse(item))
	}
	return FilePageResult{Items: result, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrFileRepoNotFound {
			return errcode.ErrFileNotFound.AsError()
		}
		return err
	}
	if err := s.storage.Delete(ctx, item); err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if err == ErrFileRepoNotFound {
			return errcode.ErrFileNotFound.AsError()
		}
		return err
	}
	return nil
}

func (s *service) save(ctx context.Context, uploaderID *uint, fileHeader *multipart.FileHeader, category string) (FileResponse, error) {
	prepared, err := s.validator.Validate(fileHeader)
	if err != nil {
		return FileResponse{}, err
	}
	return s.savePrepared(ctx, uploaderID, fileHeader, category, prepared)
}

func (s *service) savePrepared(ctx context.Context, uploaderID *uint, fileHeader *multipart.FileHeader, category string, prepared PreparedFile) (FileResponse, error) {
	if s.storage == nil {
		return FileResponse{}, errcode.ErrFileStorageUnsupported.AsError()
	}

	stored, err := s.storage.Save(ctx, fileHeader, category)
	if err != nil {
		return FileResponse{}, err
	}
	if stored.Ext == "" {
		stored.Ext = prepared.Ext
	}
	if stored.MIMEType == "" {
		stored.MIMEType = prepared.MIMEType
	}

	now := time.Now().UTC()
	item := File{
		UploaderID:   uploaderID,
		Storage:      s.storage.Name(),
		Category:     category,
		OriginalName: strings.TrimSpace(fileHeader.Filename),
		ObjectName:   stored.ObjectName,
		Ext:          strings.ToLower(filepath.Ext(fileHeader.Filename)),
		MIMEType:     stored.MIMEType,
		Size:         fileHeader.Size,
		Bucket:       stored.Bucket,
		Path:         stored.Path,
		URL:          stored.URL,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.repo.Create(ctx, item)
	if err != nil {
		_ = s.storage.Delete(ctx, item)
		return FileResponse{}, err
	}
	return toResponse(created), nil
}

func normalizeListQuery(q *ListFilesQuery) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
	q.Keyword = strings.TrimSpace(q.Keyword)
	q.Category = strings.TrimSpace(q.Category)
	q.Storage = strings.TrimSpace(q.Storage)
}

func inferCategory(fileHeader *multipart.FileHeader) string {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt":
		return "document"
	default:
		return "default"
	}
}
