package user

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	appcasbin "go-server/internal/casbin"
	"go-server/internal/errcode"

	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
	GetUser(ctx context.Context, id uint) (UserResponse, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	ListUsers(ctx context.Context) ([]UserResponse, error)
	ListUsersPage(ctx context.Context, q ListUsersQuery) (UserPageResult, error)
	UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (UserResponse, error)
	UpdateProfile(ctx context.Context, id uint, req UpdateProfileRequest) (UserResponse, error)
	UpdateAvatar(ctx context.Context, id uint, avatar string) (UserResponse, error)
	ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error
	DeleteUser(ctx context.Context, id uint) error
	DeleteUserBatch(ctx context.Context, ids []uint) error
	UpdateUserStatus(ctx context.Context, id uint, status int8) (UserResponse, error)
	ResetPassword(ctx context.Context, id uint, newPassword string) error
	ExportUsers(ctx context.Context, q ListUsersQuery) ([]byte, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, errcode.ErrInternalError.AsError()
	}

	u := User{
		Username:  strings.TrimSpace(req.Username),
		Password:  string(hashed),
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:     strings.TrimSpace(req.Phone),
		DeptID:    req.DeptID,
		Status:    1,
		CreatedAt: time.Now().UTC(),
	}

	created, err := s.repo.Create(ctx, u)
	if err != nil {
		if err == ErrUserDuplicated {
			return UserResponse{}, errcode.ErrUserEmailExists.AsError()
		}
		return UserResponse{}, err
	}

	if len(req.RoleIDs) > 0 {
		if err := s.repo.SetUserRoles(ctx, created.ID, req.RoleIDs); err != nil {
			return UserResponse{}, err
		}
	}

	resp, err := s.GetUser(ctx, created.ID)
	if err != nil {
		return UserResponse{}, err
	}

	if enforcer := appcasbin.Get(); enforcer != nil {
		for _, r := range resp.Roles {
			_, _ = enforcer.AddRoleForUser(resp.Username, r.Code)
		}
	}

	return resp, nil
}

func (s *service) GetUser(ctx context.Context, id uint) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}
	return toResponse(u), nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *service) ListUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toResponse(u))
	}
	return resp, nil
}

func (s *service) ListUsersPage(ctx context.Context, q ListUsersQuery) (UserPageResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	users, total, err := s.repo.ListPage(ctx, q)
	if err != nil {
		return UserPageResult{}, err
	}

	items := make([]UserResponse, 0, len(users))
	for _, u := range users {
		items = append(items, toResponse(u))
	}
	return UserPageResult{
		Items:    items,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}

	u.Name = strings.TrimSpace(req.Name)
	u.Email = strings.ToLower(strings.TrimSpace(req.Email))
	u.Phone = strings.TrimSpace(req.Phone)
	u.DeptID = req.DeptID

	if _, err := s.repo.Update(ctx, u); err != nil {
		return UserResponse{}, err
	}

	if err := s.repo.SetUserRoles(ctx, id, req.RoleIDs); err != nil {
		return UserResponse{}, err
	}

	resp, err := s.GetUser(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}

	if enforcer := appcasbin.Get(); enforcer != nil {
		_, _ = enforcer.DeleteRolesForUser(resp.Username)
		for _, r := range resp.Roles {
			_, _ = enforcer.AddRoleForUser(resp.Username, r.Code)
		}
	}

	return resp, nil
}

func (s *service) UpdateProfile(ctx context.Context, id uint, req UpdateProfileRequest) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}

	u.Name = strings.TrimSpace(req.Name)
	u.Email = strings.ToLower(strings.TrimSpace(req.Email))
	u.Phone = strings.TrimSpace(req.Phone)

	if _, err := s.repo.Update(ctx, u); err != nil {
		if err == ErrUserDuplicated {
			return UserResponse{}, errcode.ErrUserEmailExists.AsError()
		}
		return UserResponse{}, err
	}

	return s.GetUser(ctx, id)
}

func (s *service) UpdateAvatar(ctx context.Context, id uint, avatar string) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}

	u.Avatar = strings.TrimSpace(avatar)
	if _, err := s.repo.Update(ctx, u); err != nil {
		return UserResponse{}, err
	}

	return s.GetUser(ctx, id)
}

func (s *service) ChangePassword(ctx context.Context, id uint, oldPassword, newPassword string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return errcode.ErrUserNotFound.AsError()
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword)); err != nil {
		return errcode.ErrInvalidCredentials.AsError()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternalError.AsError()
	}

	return s.repo.UpdatePassword(ctx, id, string(hashed))
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err == ErrUserNotFound {
		return errcode.ErrUserNotFound.AsError()
	}
	return err
}

func (s *service) DeleteUserBatch(ctx context.Context, ids []uint) error {
	return s.repo.DeleteBatch(ctx, ids)
}

func (s *service) UpdateUserStatus(ctx context.Context, id uint, status int8) (UserResponse, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}
	return s.GetUser(ctx, id)
}

func (s *service) ResetPassword(ctx context.Context, id uint, newPassword string) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if err == ErrUserNotFound {
			return errcode.ErrUserNotFound.AsError()
		}
		return err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternalError.AsError()
	}
	return s.repo.UpdatePassword(ctx, id, string(hashed))
}

func (s *service) ExportUsers(ctx context.Context, q ListUsersQuery) ([]byte, error) {
	q.Page = 1
	q.PageSize = 10000
	users, _, err := s.repo.ListPage(ctx, q)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	defer func() { _ = file.Close() }()

	const sheetName = "Users"
	defaultSheet := file.GetSheetName(file.GetActiveSheetIndex())
	if defaultSheet != sheetName {
		file.SetSheetName(defaultSheet, sheetName)
	}

	headers := []string{"ID", "用户名", "姓名", "邮箱", "手机", "状态", "创建时间"}
	for idx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(idx+1, 1)
		if err := file.SetCellValue(sheetName, cell, header); err != nil {
			return nil, errcode.ErrInternalError.AsError()
		}
	}

	headerStyle, err := file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	if err != nil {
		return nil, errcode.ErrInternalError.AsError()
	}
	if err := file.SetCellStyle(sheetName, "A1", "G1", headerStyle); err != nil {
		return nil, errcode.ErrInternalError.AsError()
	}

	for idx, u := range users {
		row := idx + 2
		status := "启用"
		if u.Status == 0 {
			status = "禁用"
		}

		values := []string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Username,
			u.Name,
			u.Email,
			u.Phone,
			status,
			u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for col, value := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			if err := file.SetCellValue(sheetName, cell, value); err != nil {
				return nil, errcode.ErrInternalError.AsError()
			}
		}
	}

	widths := map[string]float64{
		"A": 10,
		"B": 18,
		"C": 18,
		"D": 28,
		"E": 18,
		"F": 12,
		"G": 22,
	}
	for column, width := range widths {
		if err := file.SetColWidth(sheetName, column, column, width); err != nil {
			return nil, errcode.ErrInternalError.AsError()
		}
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return nil, errcode.ErrInternalError.AsError()
	}
	return buf.Bytes(), nil
}
