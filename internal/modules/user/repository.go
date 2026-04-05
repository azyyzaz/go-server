package user

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserDuplicated = errors.New("user duplicated")
)

type Repository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id uint) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	List(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id uint) error
}

type inMemoryRepository struct {
	mu      sync.RWMutex
	data    map[uint]User
	counter uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{data: make(map[uint]User)}
}

func (r *inMemoryRepository) nextID() uint {
	return uint(atomic.AddUint64(&r.counter, 1))
}

func (r *inMemoryRepository) Create(_ context.Context, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.data {
		if existing.Username == user.Username || existing.Email == user.Email {
			return User{}, ErrUserDuplicated
		}
	}
	user.ID = r.nextID()
	r.data[user.ID] = user
	return user, nil
}

func (r *inMemoryRepository) GetByID(_ context.Context, id uint) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.data[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryRepository) GetByUsername(_ context.Context, username string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.data {
		if u.Username == username {
			return u, nil
		}
	}
	return User{}, ErrUserNotFound
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

func (r *inMemoryRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrUserNotFound
	}
	delete(r.data, id)
	return nil
}
