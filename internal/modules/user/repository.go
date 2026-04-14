package user

import (
	"context"
	"errors"
	"strings"
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
	ListPage(ctx context.Context, q ListUsersQuery) ([]User, int64, error)
	SetUserRoles(ctx context.Context, userID uint, roleIDs []uint) error
	Update(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id uint) error
	DeleteBatch(ctx context.Context, ids []uint) error
	UpdateStatus(ctx context.Context, id uint, status int8) error
	UpdatePassword(ctx context.Context, id uint, hashedPwd string) error
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

func (r *inMemoryRepository) ListPage(_ context.Context, q ListUsersQuery) ([]User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []User
	for _, u := range r.data {
		if q.Username != "" && !strings.Contains(u.Username, q.Username) {
			continue
		}
		if q.Name != "" && !strings.Contains(u.Name, q.Name) {
			continue
		}
		if q.Status != nil && u.Status != *q.Status {
			continue
		}
		filtered = append(filtered, u)
	}
	total := int64(len(filtered))
	start := (q.Page - 1) * q.PageSize
	if start >= len(filtered) {
		return []User{}, total, nil
	}
	end := start + q.PageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (r *inMemoryRepository) SetUserRoles(_ context.Context, _ uint, _ []uint) error {
	return nil // 内存实现不支持角色持久化
}

func (r *inMemoryRepository) Update(_ context.Context, u User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[u.ID]; !ok {
		return User{}, ErrUserNotFound
	}
	for id, existing := range r.data {
		if id == u.ID {
			continue
		}
		if existing.Username == u.Username || existing.Email == u.Email {
			return User{}, ErrUserDuplicated
		}
	}
	r.data[u.ID] = u
	return u, nil
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

func (r *inMemoryRepository) DeleteBatch(_ context.Context, ids []uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		delete(r.data, id)
	}
	return nil
}

func (r *inMemoryRepository) UpdateStatus(_ context.Context, id uint, status int8) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.data[id]
	if !ok {
		return ErrUserNotFound
	}
	u.Status = status
	r.data[id] = u
	return nil
}

func (r *inMemoryRepository) UpdatePassword(_ context.Context, id uint, hashedPwd string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.data[id]
	if !ok {
		return ErrUserNotFound
	}
	u.Password = hashedPwd
	r.data[id] = u
	return nil
}
