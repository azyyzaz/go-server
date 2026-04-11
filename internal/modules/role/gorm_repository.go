package role

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "1062") ||
		strings.Contains(strings.ToLower(err.Error()), "unique constraint")
}

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, role Role) (Role, error) {
	if err := r.db.WithContext(ctx).Create(&role).Error; err != nil {
		if isDuplicateError(err) {
			return Role{}, ErrRoleDuplicated
		}
		return Role{}, err
	}
	return role, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (Role, error) {
	var role Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if err == gorm.ErrRecordNotFound {
		return Role{}, ErrRoleNotFound
	}
	return role, err
}

func (r *gormRepository) ListPage(ctx context.Context, q ListRolesQuery) ([]Role, int64, error) {
	var (
		items []Role
		total int64
	)

	db := r.db.WithContext(ctx).Model(&Role{})
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.Code != "" {
		db = db.Where("code LIKE ?", "%"+q.Code+"%")
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.PageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *gormRepository) Update(ctx context.Context, role Role) (Role, error) {
	if err := r.db.WithContext(ctx).Save(&role).Error; err != nil {
		if isDuplicateError(err) {
			return Role{}, ErrRoleDuplicated
		}
		return Role{}, err
	}
	return role, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Role{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRoleNotFound
	}
	return nil
}

func (r *gormRepository) CountUsers(ctx context.Context, roleID uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("user_roles").Where("role_id = ?", roleID).Count(&total).Error
	return total, err
}

func (r *gormRepository) ReplaceRoleMenus(ctx context.Context, roleID uint, menuIDs []uint) error {
	role := Role{ID: roleID}
	menus := make([]Menu, 0, len(menuIDs))
	for _, id := range menuIDs {
		menus = append(menus, Menu{ID: id})
	}
	return r.db.WithContext(ctx).Model(&role).Association("Menus").Replace(menus)
}

func (r *gormRepository) ListMenus(ctx context.Context) ([]Menu, error) {
	var menus []Menu
	err := r.db.WithContext(ctx).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *gormRepository) GetRoleMenuIDs(ctx context.Context, roleID uint) ([]uint, error) {
	if _, err := r.GetByID(ctx, roleID); err != nil {
		return nil, err
	}

	var ids []uint
	if err := r.db.WithContext(ctx).Table("role_menus").Where("role_id = ?", roleID).Pluck("menu_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *gormRepository) ListRoleUsersPage(ctx context.Context, roleID uint, q ListRoleUsersQuery) ([]RoleUserResponse, int64, error) {
	if _, err := r.GetByID(ctx, roleID); err != nil {
		return nil, 0, err
	}

	base := r.db.WithContext(ctx).Table("users AS u").
		Joins("JOIN user_roles ur ON ur.user_id = u.id").
		Where("ur.role_id = ?", roleID)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []RoleUserResponse
	offset := (q.Page - 1) * q.PageSize
	err := base.Select("u.id, u.username, u.name, u.email, u.phone, u.status").
		Order("u.created_at DESC").
		Offset(offset).
		Limit(q.PageSize).
		Scan(&items).Error
	return items, total, err
}
