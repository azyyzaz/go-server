package user

import (
	"context"
	"sort"
	"strings"
	"time"

	"go-server/internal/errcode"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
	GetUser(ctx context.Context, id uint) (UserResponse, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	ListUsers(ctx context.Context) ([]UserResponse, error)
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

	user := User{
		Username:  strings.TrimSpace(req.Username),
		Password:  string(hashed),
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		CreatedAt: time.Now().UTC(),
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		if err == ErrUserDuplicated {
			return UserResponse{}, errcode.ErrUserEmailExists.AsError()
		}
		return UserResponse{}, err
	}
	return toResponse(created), nil
}

func (s *service) GetUser(ctx context.Context, id uint) (UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}
	return toResponse(user), nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *service) ListUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.After(users[j].CreatedAt)
	})

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toResponse(u))
	}
	return resp, nil
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return errcode.ErrUserNotFound.AsError()
		}
		return err
	}
	return nil
}
