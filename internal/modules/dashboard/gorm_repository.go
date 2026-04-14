package dashboard

import (
	"context"
	"time"

	"go-server/internal/modules/audit"
	"go-server/internal/modules/user"

	rdb "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const onlineUserKeyPrefix = "dashboard:online:user:"

type gormRepository struct {
	db    *gorm.DB
	redis *rdb.Client
}

func NewGORMRepository(db *gorm.DB, redis *rdb.Client) Repository {
	return &gormRepository{db: db, redis: redis}
}

func (r *gormRepository) CountUsers(ctx context.Context) (int64, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&total).Error; err != nil {
		return 0, 0, err
	}

	var active int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Where("status = ?", 1).Count(&active).Error; err != nil {
		return 0, 0, err
	}
	return total, active, nil
}

func (r *gormRepository) CountOperationLogs(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&audit.OperationLog{}).Count(&total).Error
	return total, err
}

func (r *gormRepository) CountLoginLogs(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&audit.LoginLog{}).Count(&total).Error
	return total, err
}

func (r *gormRepository) CountOperationLogsSince(ctx context.Context, since time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&audit.OperationLog{}).Where("created_at >= ?", since).Count(&total).Error
	return total, err
}

func (r *gormRepository) CountOnlineUsers(ctx context.Context) (int64, error) {
	if r.redis == nil {
		return 0, nil
	}

	var (
		cursor uint64
		total  int64
	)
	for {
		keys, next, err := r.redis.Scan(ctx, cursor, onlineUserKeyPrefix+"*", 100).Result()
		if err != nil {
			return 0, err
		}
		total += int64(len(keys))
		cursor = next
		if cursor == 0 {
			return total, nil
		}
	}
}
