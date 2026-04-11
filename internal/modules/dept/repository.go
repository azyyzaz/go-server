package dept

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
)

var ErrDeptNotFound = errors.New("dept not found")

type Repository interface {
	Create(ctx context.Context, item Dept) (Dept, error)
	GetByID(ctx context.Context, id uint) (Dept, error)
	List(ctx context.Context) ([]Dept, error)
	Update(ctx context.Context, item Dept) (Dept, error)
	Delete(ctx context.Context, id uint) error
	CountChildren(ctx context.Context, id uint) (int64, error)
	CountUsers(ctx context.Context, id uint) (int64, error)
	ListDeptUsersPage(ctx context.Context, id uint, q ListDeptUsersQuery) ([]DeptUserResponse, int64, error)
}

type inMemoryRepository struct {
	mu      sync.RWMutex
	data    map[uint]Dept
	counter uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{data: make(map[uint]Dept)}
}

func (r *inMemoryRepository) nextID() uint {
	return uint(atomic.AddUint64(&r.counter, 1))
}

func (r *inMemoryRepository) Create(_ context.Context, item Dept) (Dept, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = r.nextID()
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) GetByID(_ context.Context, id uint) (Dept, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.data[id]
	if !ok {
		return Dept{}, ErrDeptNotFound
	}
	return item, nil
}

func (r *inMemoryRepository) List(_ context.Context) ([]Dept, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]Dept, 0, len(r.data))
	for _, item := range r.data {
		items = append(items, item)
	}
	slices.SortFunc(items, func(a, b Dept) int {
		if a.Sort != b.Sort {
			return a.Sort - b.Sort
		}
		return int(a.ID) - int(b.ID)
	})
	return items, nil
}

func (r *inMemoryRepository) Update(_ context.Context, item Dept) (Dept, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[item.ID]; !ok {
		return Dept{}, ErrDeptNotFound
	}
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return ErrDeptNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *inMemoryRepository) CountChildren(_ context.Context, id uint) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var total int64
	for _, item := range r.data {
		if item.ParentID != nil && *item.ParentID == id {
			total++
		}
	}
	return total, nil
}

func (r *inMemoryRepository) CountUsers(_ context.Context, _ uint) (int64, error) {
	return 0, nil
}

func (r *inMemoryRepository) ListDeptUsersPage(_ context.Context, id uint, _ ListDeptUsersQuery) ([]DeptUserResponse, int64, error) {
	if _, err := r.GetByID(context.Background(), id); err != nil {
		return nil, 0, err
	}
	return []DeptUserResponse{}, 0, nil
}
