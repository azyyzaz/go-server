package menu

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

func (r *gormRepository) Create(ctx context.Context, item Menu) (Menu, error) {
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return Menu{}, err
	}
	return item, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (Menu, error) {
	var item Menu
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err == gorm.ErrRecordNotFound {
		return Menu{}, ErrMenuNotFound
	}
	return item, err
}

func (r *gormRepository) List(ctx context.Context) ([]Menu, error) {
	var items []Menu
	err := r.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&items).Error
	return items, err
}

func (r *gormRepository) Update(ctx context.Context, item Menu) (Menu, error) {
	if err := r.db.WithContext(ctx).Save(&item).Error; err != nil {
		return Menu{}, err
	}
	return item, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Menu{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMenuNotFound
	}
	return nil
}

func (r *gormRepository) CountChildren(ctx context.Context, id uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&Menu{}).Where("parent_id = ?", id).Count(&total).Error
	return total, err
}

func (r *gormRepository) UpdateSorts(ctx context.Context, items []MenuSortItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			result := tx.Model(&Menu{}).Where("id = ?", item.ID).Update("sort", item.Sort)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrMenuNotFound
			}
		}
		return nil
	})
}

func (r *gormRepository) ListCurrentUserMenus(ctx context.Context, userID uint) ([]Menu, error) {
	var items []Menu
	err := r.db.WithContext(ctx).
		Table("menus AS m").
		Distinct("m.id, m.parent_id, m.name, m.type, m.path, m.component, m.permission, m.sort, m.visible, m.status, m.created_at, m.updated_at").
		Joins("JOIN role_menus rm ON rm.menu_id = m.id").
		Joins("JOIN user_roles ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ? AND m.status = ? AND m.visible = ?", userID, 1, 1).
		Order("m.sort ASC, m.id ASC").
		Scan(&items).Error
	return items, err
}
