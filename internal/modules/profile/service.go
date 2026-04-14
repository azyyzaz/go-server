package profile

import (
	"context"
	"mime/multipart"

	"go-server/internal/errcode"
	"go-server/internal/modules/audit"
	filemodule "go-server/internal/modules/file"
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
	userSvc  user.Service
	auditSvc audit.Service
	fileSvc  filemodule.Service
}

func NewService(userSvc user.Service, auditSvc audit.Service, fileSvc filemodule.Service) Service {
	return &service{userSvc: userSvc, auditSvc: auditSvc, fileSvc: fileSvc}
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
	if file == nil {
		return AvatarUploadResponse{}, errcode.ErrInvalidParam.AsError()
	}

	uploaded, err := s.fileSvc.UploadAvatar(ctx, userID, file)
	if err != nil {
		return AvatarUploadResponse{}, err
	}
	if _, err := s.userSvc.UpdateAvatar(ctx, userID, uploaded.URL); err != nil {
		return AvatarUploadResponse{}, err
	}

	return AvatarUploadResponse{Avatar: uploaded.URL}, nil
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
