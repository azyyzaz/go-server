package dashboard

import (
	"context"
	"time"
)

type Repository interface {
	CountUsers(ctx context.Context) (total int64, active int64, err error)
	CountOperationLogs(ctx context.Context) (int64, error)
	CountLoginLogs(ctx context.Context) (int64, error)
	CountOperationLogsSince(ctx context.Context, since time.Time) (int64, error)
	CountOnlineUsers(ctx context.Context) (int64, error)
}

type inMemoryRepository struct{}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{}
}

func (r *inMemoryRepository) CountUsers(context.Context) (int64, int64, error) {
	return 0, 0, nil
}

func (r *inMemoryRepository) CountOperationLogs(context.Context) (int64, error) {
	return 0, nil
}

func (r *inMemoryRepository) CountLoginLogs(context.Context) (int64, error) {
	return 0, nil
}

func (r *inMemoryRepository) CountOperationLogsSince(context.Context, time.Time) (int64, error) {
	return 0, nil
}

func (r *inMemoryRepository) CountOnlineUsers(context.Context) (int64, error) {
	return 0, nil
}
