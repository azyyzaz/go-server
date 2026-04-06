package user

import (
	"context"
	"strings"
	"time"

	appcasbin "go-server/internal/casbin"
	"go-server/internal/errcode"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
	GetUser(ctx context.Context, id uint) (UserResponse, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	ListUsers(ctx context.Context) ([]UserResponse, error)
	ListUsersPage(ctx context.Context, q ListUsersQuery) (UserPageResult, error)
	UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
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

	// 角色分配
	if len(req.RoleIDs) > 0 {
		if err := s.repo.SetUserRoles(ctx, created.ID, req.RoleIDs); err != nil {
			return UserResponse{}, err
		}
	}

	// 重新查一次，带上 Roles
	resp, err := s.GetUser(ctx, created.ID)
	if err != nil {
		return UserResponse{}, err
	}

	// 同步写入 Casbin：用户名 → 角色 code
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

	// 同步 Casbin：先删旧绑定，再写新角色
	if enforcer := appcasbin.Get(); enforcer != nil {
		_, _ = enforcer.DeleteRolesForUser(resp.Username)
		for _, r := range resp.Roles {
			_, _ = enforcer.AddRoleForUser(resp.Username, r.Code)
		}
	}

	return resp, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err == ErrUserNotFound {
		return errcode.ErrUserNotFound.AsError()
	}
	return err
}
