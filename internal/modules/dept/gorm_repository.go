package dept

import (
	"context"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, item Dept) (Dept, error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return Dept{}, err
	}
	return item, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (Dept, error) {
	var item Dept
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err == gorm.ErrRecordNotFound {
		return Dept{}, ErrDeptNotFound
	}
	return item, err
}

func (r *gormRepository) List(ctx context.Context) ([]Dept, error) {
	var items []Dept
	err := r.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *gormRepository) Update(ctx context.Context, item Dept) (Dept, error) {
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return Dept{}, err
	}
	return item, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Dept{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrDeptNotFound
	}
	return nil
}

func (r *gormRepository) CountChildren(ctx context.Context, id uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&Dept{}).Where("parent_id = ?", id).Count(&total).Error
	return total, err
}

func (r *gormRepository) CountUsers(ctx context.Context, id uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("users").Where("dept_id = ?", id).Count(&total).Error
	return total, err
}

func (r *gormRepository) ListDeptUsersPage(ctx context.Context, id uint, q ListDeptUsersQuery) ([]DeptUserResponse, int64, error) {
	if _, err := r.GetByID(ctx, id); err != nil {
		return nil, 0, err
	}

	base := r.db.WithContext(ctx).Table("users").Where("dept_id = ?", id)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []DeptUserResponse
	offset := (q.Page - 1) * q.PageSize
	err := base.Select("id, username, name, email, phone, status").
		Order("created_at DESC").
		Offset(offset).
		Limit(q.PageSize).
		Scan(&items).Error
	return items, total, err
}
