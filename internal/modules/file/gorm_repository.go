package file

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

func (r *gormRepository) Create(ctx context.Context, item File) (File, error) {
	err := r.db.WithContext(ctx).Create(&item).Error
	return item, err
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (File, error) {
	var item File
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err == gorm.ErrRecordNotFound {
		return File{}, ErrFileRepoNotFound
	}
	return item, err
}

func (r *gormRepository) List(ctx context.Context, q ListFilesQuery, start, end *time.Time) ([]File, int64, error) {
	var (
		items []File
		total int64
	)

	db := r.db.WithContext(ctx).Model(&File{})
	if q.Keyword != "" {
		keyword := "%" + strings.TrimSpace(q.Keyword) + "%"
		db = db.Where("original_name LIKE ? OR object_name LIKE ?", keyword, keyword)
	}
	if q.Category != "" {
		db = db.Where("category = ?", strings.TrimSpace(q.Category))
	}
	if q.Storage != "" {
		db = db.Where("storage = ?", strings.TrimSpace(q.Storage))
	}
	if q.UploaderID != nil {
		db = db.Where("uploader_id = ?", *q.UploaderID)
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

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&File{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrFileRepoNotFound
	}
	return nil
}
