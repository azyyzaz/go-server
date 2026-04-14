package user

import (
	"context"
	"strings"

	"go-server/internal/modules/role"

	"gorm.io/gorm"
)

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "1062")
}

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, u User) (User, error) {
	result := r.db.WithContext(ctx).Create(&u)
	if result.Error != nil {
		if isDuplicateError(result.Error) {
			return User{}, ErrUserDuplicated
		}
		return User{}, result.Error
	}
	return u, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (User, error) {
	var u User
	err := r.db.WithContext(ctx).Preload("Roles").First(&u, id).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, ErrUserNotFound
	}
	return u, err
}

func (r *gormRepository) GetByUsername(ctx context.Context, username string) (User, error) {
	var u User
	err := r.db.WithContext(ctx).Preload("Roles").Where("username = ?", username).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, ErrUserNotFound
	}
	return u, err
}

func (r *gormRepository) List(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).Preload("Roles").Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *gormRepository) ListPage(ctx context.Context, q ListUsersQuery) ([]User, int64, error) {
	var (
		users []User
		total int64
	)

	db := r.db.WithContext(ctx).Model(&User{})

	if q.Username != "" {
		db = db.Where("username LIKE ?", "%"+q.Username+"%")
	}
	if q.Name != "" {
		db = db.Where("name LIKE ?", "%"+q.Name+"%")
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.PageSize
	err := db.Preload("Roles").Order("created_at DESC").Offset(offset).Limit(q.PageSize).Find(&users).Error
	return users, total, err
}

// SetUserRoles 替换用户的全部角色（先清后写）
func (r *gormRepository) SetUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	u := User{ID: userID}
	// 构造只含 ID 的 role slice，避免触发 GORM 的 upsert
	roles := make([]role.Role, 0, len(roleIDs))
	for _, id := range roleIDs {
		roles = append(roles, role.Role{ID: id})
	}
	return r.db.WithContext(ctx).Model(&u).Association("Roles").Replace(roles)
}

func (r *gormRepository) Update(ctx context.Context, u User) (User, error) {
	result := r.db.WithContext(ctx).Save(&u)
	if result.Error != nil {
		if isDuplicateError(result.Error) {
			return User{}, ErrUserDuplicated
		}
		return User{}, result.Error
	}
	return u, nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *gormRepository) DeleteBatch(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&User{}, ids).Error
}

func (r *gormRepository) UpdateStatus(ctx context.Context, id uint, status int8) error {
	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *gormRepository) UpdatePassword(ctx context.Context, id uint, hashedPwd string) error {
	result := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", hashedPwd)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
