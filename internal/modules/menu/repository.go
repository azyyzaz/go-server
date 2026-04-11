package menu

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
)

var ErrMenuNotFound = errors.New("menu not found")

type Repository interface {
	Create(ctx context.Context, item Menu) (Menu, error)
	GetByID(ctx context.Context, id uint) (Menu, error)
	List(ctx context.Context) ([]Menu, error)
	Update(ctx context.Context, item Menu) (Menu, error)
	Delete(ctx context.Context, id uint) error
	CountChildren(ctx context.Context, id uint) (int64, error)
	UpdateSorts(ctx context.Context, items []MenuSortItem) error
	ListCurrentUserMenus(ctx context.Context, userID uint) ([]Menu, error)
}

type inMemoryRepository struct {
	mu      sync.RWMutex
	data    map[uint]Menu
	counter uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{data: make(map[uint]Menu)}
}

func (r *inMemoryRepository) nextID() uint {
	return uint(atomic.AddUint64(&r.counter, 1))
}

func (r *inMemoryRepository) Create(_ context.Context, item Menu) (Menu, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = r.nextID()
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) GetByID(_ context.Context, id uint) (Menu, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.data[id]
	if !ok {
		return Menu{}, ErrMenuNotFound
	}
	return item, nil
}

func (r *inMemoryRepository) List(_ context.Context) ([]Menu, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]Menu, 0, len(r.data))
	for _, item := range r.data {
		items = append(items, item)
	}
	slices.SortFunc(items, func(a, b Menu) int {
		if a.Sort != b.Sort {
			return a.Sort - b.Sort
		}
		return int(a.ID) - int(b.ID)
	})
	return items, nil
}

func (r *inMemoryRepository) Update(_ context.Context, item Menu) (Menu, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[item.ID]; !ok {
		return Menu{}, ErrMenuNotFound
	}
	r.data[item.ID] = item
	return item, nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrMenuNotFound
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

func (r *inMemoryRepository) UpdateSorts(_ context.Context, items []MenuSortItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		menu, ok := r.data[item.ID]
		if !ok {
			return ErrMenuNotFound
		}
		menu.Sort = item.Sort
		r.data[item.ID] = menu
	}
	return nil
}

func (r *inMemoryRepository) ListCurrentUserMenus(ctx context.Context, _ uint) ([]Menu, error) {
	return r.List(ctx)
}
