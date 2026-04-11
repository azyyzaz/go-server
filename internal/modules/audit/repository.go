package audit

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
	CreateOperationLog(ctx context.Context, item OperationLog) error
	CreateLoginLog(ctx context.Context, item LoginLog) error
	ListOperationLogs(ctx context.Context, q ListOperationLogsQuery, start, end *time.Time) ([]OperationLog, int64, error)
	ListLoginLogs(ctx context.Context, q ListLoginLogsQuery, start, end *time.Time) ([]LoginLog, int64, error)
	DeleteOperationLogsBefore(ctx context.Context, before time.Time) (int64, error)
	DeleteLoginLogsBefore(ctx context.Context, before time.Time) (int64, error)
}

var ErrAuditInvalidTimeRange = errors.New("invalid audit time range")

type inMemoryRepository struct {
	mu              sync.RWMutex
	operationLogs   []OperationLog
	loginLogs       []LoginLog
	operationIDSeed uint64
	loginIDSeed     uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{}
}

func (r *inMemoryRepository) CreateOperationLog(_ context.Context, item OperationLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = uint(atomic.AddUint64(&r.operationIDSeed, 1))
	r.operationLogs = append(r.operationLogs, item)
	return nil
}

func (r *inMemoryRepository) CreateLoginLog(_ context.Context, item LoginLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item.ID = uint(atomic.AddUint64(&r.loginIDSeed, 1))
	r.loginLogs = append(r.loginLogs, item)
	return nil
}

func (r *inMemoryRepository) ListOperationLogs(_ context.Context, q ListOperationLogsQuery, start, end *time.Time) ([]OperationLog, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]OperationLog, 0, len(r.operationLogs))
	for _, item := range r.operationLogs {
		if q.Username != "" && !strings.Contains(strings.ToLower(item.Username), strings.ToLower(q.Username)) {
			continue
		}
		if q.Method != "" && !strings.EqualFold(item.Method, q.Method) {
			continue
		}
		if q.Path != "" && !strings.Contains(strings.ToLower(item.Path), strings.ToLower(q.Path)) {
			continue
		}
		if q.Status != nil && item.StatusCode != *q.Status {
			continue
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
	startIndex := (q.Page - 1) * q.PageSize
	if startIndex >= len(filtered) {
		return []OperationLog{}, total, nil
	}
	endIndex := startIndex + q.PageSize
	if endIndex > len(filtered) {
		endIndex = len(filtered)
	}
	return filtered[startIndex:endIndex], total, nil
}

func (r *inMemoryRepository) ListLoginLogs(_ context.Context, q ListLoginLogsQuery, start, end *time.Time) ([]LoginLog, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]LoginLog, 0, len(r.loginLogs))
	for _, item := range r.loginLogs {
		if q.Username != "" && !strings.Contains(strings.ToLower(item.Username), strings.ToLower(q.Username)) {
			continue
		}
		if q.IP != "" && !strings.Contains(item.IP, q.IP) {
			continue
		}
		if q.Success != nil && item.Success != *q.Success {
			continue
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
	startIndex := (q.Page - 1) * q.PageSize
	if startIndex >= len(filtered) {
		return []LoginLog{}, total, nil
	}
	endIndex := startIndex + q.PageSize
	if endIndex > len(filtered) {
		endIndex = len(filtered)
	}
	return filtered[startIndex:endIndex], total, nil
}

func (r *inMemoryRepository) DeleteOperationLogsBefore(_ context.Context, before time.Time) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := r.operationLogs[:0]
	var deleted int64
	for _, item := range r.operationLogs {
		if item.CreatedAt.Before(before) {
			deleted++
			continue
		}
		filtered = append(filtered, item)
	}
	r.operationLogs = filtered
	return deleted, nil
}

func (r *inMemoryRepository) DeleteLoginLogsBefore(_ context.Context, before time.Time) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := r.loginLogs[:0]
	var deleted int64
	for _, item := range r.loginLogs {
		if item.CreatedAt.Before(before) {
			deleted++
			continue
		}
		filtered = append(filtered, item)
	}
	r.loginLogs = filtered
	return deleted, nil
}
