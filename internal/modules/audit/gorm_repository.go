package audit

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateOperationLog(ctx context.Context, item OperationLog) error {
	return r.db.WithContext(ctx).Create(&item).Error
}

func (r *gormRepository) CreateLoginLog(ctx context.Context, item LoginLog) error {
	return r.db.WithContext(ctx).Create(&item).Error
}

func (r *gormRepository) ListOperationLogs(ctx context.Context, q ListOperationLogsQuery, start, end *time.Time) ([]OperationLog, int64, error) {
	var (
		items []OperationLog
		total int64
	)

	db := r.db.WithContext(ctx).Model(&OperationLog{})
	if q.Username != "" {
		db = db.Where("username LIKE ?", "%"+strings.TrimSpace(q.Username)+"%")
	}
	if q.Method != "" {
		db = db.Where("method = ?", strings.ToUpper(strings.TrimSpace(q.Method)))
	}
	if q.Path != "" {
		db = db.Where("path LIKE ?", "%"+strings.TrimSpace(q.Path)+"%")
	}
	if q.Status != nil {
		db = db.Where("status_code = ?", *q.Status)
	}
	if start != nil {
		db = db.Where("created_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("created_at <= ?", *end)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.PageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(q.PageSize).Find(&items).Error
	return items, total, err
}

func (r *gormRepository) ListLoginLogs(ctx context.Context, q ListLoginLogsQuery, start, end *time.Time) ([]LoginLog, int64, error) {
	var (
		items []LoginLog
		total int64
	)

	db := r.db.WithContext(ctx).Model(&LoginLog{})
	if q.Username != "" {
		db = db.Where("username LIKE ?", "%"+strings.TrimSpace(q.Username)+"%")
	}
	if q.IP != "" {
		db = db.Where("ip LIKE ?", "%"+strings.TrimSpace(q.IP)+"%")
	}
	if q.Success != nil {
		db = db.Where("success = ?", *q.Success)
	}
	if start != nil {
		db = db.Where("created_at >= ?", *start)
	}
	if end != nil {
		db = db.Where("created_at <= ?", *end)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.PageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(q.PageSize).Find(&items).Error
	return items, total, err
}

func (r *gormRepository) DeleteOperationLogsBefore(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&OperationLog{})
	return result.RowsAffected, result.Error
}

func (r *gormRepository) DeleteLoginLogsBefore(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&LoginLog{})
	return result.RowsAffected, result.Error
}
