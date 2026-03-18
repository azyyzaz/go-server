package user

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserDuplicated = errors.New("user email duplicated")
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	List(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id string) error
}

type inMemoryRepository struct {
	mu   sync.RWMutex
	data map[string]User
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{data: make(map[string]User)}
}

func (r *inMemoryRepository) Create(_ context.Context, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.data {
		if existing.Email == user.Email {
			return User{}, ErrUserDuplicated
		}
	}
	r.data[user.ID] = user
	return user, nil
}

func (r *inMemoryRepository) GetByID(_ context.Context, id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.data[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryRepository) List(_ context.Context) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]User, 0, len(r.data))
	for _, u := range r.data {
		users = append(users, u)
	}
	return users, nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrUserNotFound
	}
	delete(r.data, id)
	return nil
}
