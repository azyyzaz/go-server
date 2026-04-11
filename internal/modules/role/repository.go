package role

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleDuplicated = errors.New("role duplicated")
)

type Repository interface {
	Create(ctx context.Context, role Role) (Role, error)
	GetByID(ctx context.Context, id uint) (Role, error)
	ListPage(ctx context.Context, q ListRolesQuery) ([]Role, int64, error)
	Update(ctx context.Context, role Role) (Role, error)
	Delete(ctx context.Context, id uint) error
	CountUsers(ctx context.Context, roleID uint) (int64, error)
	ReplaceRoleMenus(ctx context.Context, roleID uint, menuIDs []uint) error
	ListMenus(ctx context.Context) ([]Menu, error)
	GetRoleMenuIDs(ctx context.Context, roleID uint) ([]uint, error)
	ListRoleUsersPage(ctx context.Context, roleID uint, q ListRoleUsersQuery) ([]RoleUserResponse, int64, error)
}

type inMemoryRepository struct {
	mu        sync.RWMutex
	roles     map[uint]Role
	roleMenus map[uint][]uint
	menus     map[uint]Menu
	counter   uint64
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		roles:     make(map[uint]Role),
		roleMenus: make(map[uint][]uint),
		menus:     make(map[uint]Menu),
	}
}

func (r *inMemoryRepository) nextID() uint {
	return uint(atomic.AddUint64(&r.counter, 1))
}

func (r *inMemoryRepository) Create(_ context.Context, role Role) (Role, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.roles {
		if strings.EqualFold(existing.Code, role.Code) {
			return Role{}, ErrRoleDuplicated
		}
	}
	role.ID = r.nextID()
	r.roles[role.ID] = role
	return role, nil
}

func (r *inMemoryRepository) GetByID(_ context.Context, id uint) (Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, ok := r.roles[id]
	if !ok {
		return Role{}, ErrRoleNotFound
	}
	return role, nil
}

func (r *inMemoryRepository) ListPage(_ context.Context, q ListRolesQuery) ([]Role, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var roles []Role
	for _, item := range r.roles {
		if q.Name != "" && !strings.Contains(strings.ToLower(item.Name), strings.ToLower(q.Name)) {
			continue
		}
		if q.Code != "" && !strings.Contains(strings.ToLower(item.Code), strings.ToLower(q.Code)) {
			continue
		}
		if q.Status != nil && item.Status != *q.Status {
			continue
		}
		roles = append(roles, item)
	}

	total := int64(len(roles))
	start := (q.Page - 1) * q.PageSize
	if start >= len(roles) {
		return []Role{}, total, nil
	}
	end := start + q.PageSize
	if end > len(roles) {
		end = len(roles)
	}
	return roles[start:end], total, nil
}

func (r *inMemoryRepository) Update(_ context.Context, role Role) (Role, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.roles[role.ID]; !ok {
		return Role{}, ErrRoleNotFound
	}
	for id, existing := range r.roles {
		if id != role.ID && strings.EqualFold(existing.Code, role.Code) {
			return Role{}, ErrRoleDuplicated
		}
	}
	r.roles[role.ID] = role
	return role, nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.roles[id]; !ok {
		return ErrRoleNotFound
	}
	delete(r.roles, id)
	delete(r.roleMenus, id)
	return nil
}

func (r *inMemoryRepository) CountUsers(_ context.Context, _ uint) (int64, error) {
	return 0, nil
}

func (r *inMemoryRepository) ReplaceRoleMenus(_ context.Context, roleID uint, menuIDs []uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.roles[roleID]; !ok {
		return ErrRoleNotFound
	}
	r.roleMenus[roleID] = append([]uint(nil), menuIDs...)
	return nil
}

func (r *inMemoryRepository) ListMenus(_ context.Context) ([]Menu, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	menus := make([]Menu, 0, len(r.menus))
	for _, menu := range r.menus {
		menus = append(menus, menu)
	}
	return menus, nil
}

func (r *inMemoryRepository) GetRoleMenuIDs(_ context.Context, roleID uint) ([]uint, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.roles[roleID]; !ok {
		return nil, ErrRoleNotFound
	}
	return append([]uint(nil), r.roleMenus[roleID]...), nil
}

func (r *inMemoryRepository) ListRoleUsersPage(_ context.Context, roleID uint, _ ListRoleUsersQuery) ([]RoleUserResponse, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.roles[roleID]; !ok {
		return nil, 0, ErrRoleNotFound
	}
	return []RoleUserResponse{}, 0, nil
}
