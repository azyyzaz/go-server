package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"go-server/internal/errcode"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error)
	GetUser(ctx context.Context, id string) (UserResponse, error)
	ListUsers(ctx context.Context) ([]UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user := User{
		ID:        newID(),
		Name:      name,
		Email:     email,
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

func (s *service) GetUser(ctx context.Context, id string) (UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return UserResponse{}, errcode.ErrUserNotFound.AsError()
		}
		return UserResponse{}, err
	}
	return toResponse(user), nil
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

func (s *service) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			return errcode.ErrUserNotFound.AsError()
		}
		return err
	}
	return nil
}

func newID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().UTC().Format("20060102150405.000000")
	}
	return hex.EncodeToString(buf)
}
