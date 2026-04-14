package profile

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-server/internal/errcode"
	"go-server/internal/modules/audit"
	"go-server/internal/modules/user"
)

type Service interface {
	GetProfile(ctx context.Context, userID uint) (ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (ProfileResponse, error)
	ChangePassword(ctx context.Context, userID uint, req ChangePasswordRequest) error
	UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (AvatarUploadResponse, error)
	ListMyLoginLogs(ctx context.Context, userID uint, q audit.ListLoginLogsQuery) (LoginLogPageResult, error)
}

type service struct {
	userSvc   user.Service
	auditSvc  audit.Service
	uploadDir string
}

func NewService(userSvc user.Service, auditSvc audit.Service, uploadDir string) Service {
	if strings.TrimSpace(uploadDir) == "" {
		uploadDir = filepath.Join("uploads", "avatars")
	}
	return &service{userSvc: userSvc, auditSvc: auditSvc, uploadDir: uploadDir}
}

func (s *service) GetProfile(ctx context.Context, userID uint) (ProfileResponse, error) {
	return s.userSvc.GetUser(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (ProfileResponse, error) {
	return s.userSvc.UpdateProfile(ctx, userID, user.UpdateProfileRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})
}

func (s *service) ChangePassword(ctx context.Context, userID uint, req ChangePasswordRequest) error {
	return s.userSvc.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
}

func (s *service) UploadAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (AvatarUploadResponse, error) {
	if file == nil || file.Size <= 0 || file.Size > 2*1024*1024 {
		return AvatarUploadResponse{}, errcode.ErrInvalidParam.AsError()
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
	default:
		return AvatarUploadResponse{}, errcode.ErrInvalidParam.AsError()
	}

	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return AvatarUploadResponse{}, errcode.ErrInternalError.AsError()
	}

	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	dstPath := filepath.Join(s.uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return AvatarUploadResponse{}, errcode.ErrInternalError.AsError()
	}
	defer func() { _ = src.Close() }()

	dst, err := os.Create(dstPath)
	if err != nil {
		return AvatarUploadResponse{}, errcode.ErrInternalError.AsError()
	}
	defer func() { _ = dst.Close() }()

	if _, err := io.Copy(dst, src); err != nil {
		return AvatarUploadResponse{}, errcode.ErrInternalError.AsError()
	}

	avatar := "/" + filepath.ToSlash(filepath.Join("uploads", "avatars", filename))
	if _, err := s.userSvc.UpdateAvatar(ctx, userID, avatar); err != nil {
		_ = os.Remove(dstPath)
		return AvatarUploadResponse{}, err
	}

	return AvatarUploadResponse{Avatar: avatar}, nil
}

func (s *service) ListMyLoginLogs(ctx context.Context, userID uint, q audit.ListLoginLogsQuery) (LoginLogPageResult, error) {
	if s.auditSvc == nil {
		if q.Page <= 0 {
			q.Page = 1
		}
		if q.PageSize <= 0 {
			q.PageSize = 10
		}
		return LoginLogPageResult{
			Items:    []audit.LoginLogResponse{},
			Total:    0,
			Page:     q.Page,
			PageSize: q.PageSize,
		}, nil
	}

	currentUser, err := s.userSvc.GetUser(ctx, userID)
	if err != nil {
		return LoginLogPageResult{}, err
	}
	q.Username = currentUser.Username
	return s.auditSvc.ListLoginLogs(ctx, q)
}
