package file

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Repository interface {
	Create(ctx context.Context, item File) (File, error)
	GetByID(ctx context.Context, id uint) (File, error)
	List(ctx context.Context, q ListFilesQuery, start, end *time.Time) ([]File, int64, error)
	Delete(ctx context.Context, id uint) error
}

var (
	ErrFileRepoNotFound     = errors.New("file not found")
	ErrFileInvalidTimeRange = errors.New("invalid file time range")
)

type inMemoryRepository struct {
	mu     sync.RWMutex
	items  []File
	nextID uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{}
}

func (r *inMemoryRepository) Create(_ context.Context, item File) (File, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = uint(atomic.AddUint64(&r.nextID, 1))
	r.items = append(r.items, item)
	return item, nil
}

func (r *inMemoryRepository) GetByID(_ context.Context, id uint) (File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}
	return File{}, ErrFileRepoNotFound
}

func (r *inMemoryRepository) List(_ context.Context, q ListFilesQuery, start, end *time.Time) ([]File, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]File, 0, len(r.items))
	for _, item := range r.items {
		if q.Keyword != "" {
			keyword := strings.ToLower(strings.TrimSpace(q.Keyword))
			if !strings.Contains(strings.ToLower(item.OriginalName), keyword) &&
				!strings.Contains(strings.ToLower(item.ObjectName), keyword) {
				continue
			}
		}
		if q.Category != "" && !strings.EqualFold(item.Category, q.Category) {
			continue
		}
		if q.Storage != "" && !strings.EqualFold(item.Storage, q.Storage) {
			continue
		}
		if q.UploaderID != nil {
			if item.UploaderID == nil || *item.UploaderID != *q.UploaderID {
				continue
			}
		}
		if start != nil && item.CreatedAt.Before(*start) {
			continue
		}
		if end != nil && item.CreatedAt.After(*end) {
			continue
		}
		filtered = append(filtered, item)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := int64(len(filtered))
	offset := (q.Page - 1) * q.PageSize
	if offset >= len(filtered) {
		return []File{}, total, nil
	}
	limit := offset + q.PageSize
	if limit > len(filtered) {
		limit = len(filtered)
	}
	return filtered[offset:limit], total, nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for idx, item := range r.items {
		if item.ID == id {
			r.items = append(r.items[:idx], r.items[idx+1:]...)
			return nil
		}
	}
	return ErrFileRepoNotFound
}
