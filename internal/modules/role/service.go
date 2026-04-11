package role

import (
	"context"
	"strings"
	"time"

	appcasbin "go-server/internal/casbin"
	"go-server/internal/errcode"
)

type Service interface {
	ListRolesPage(ctx context.Context, q ListRolesQuery) (RolePageResult, error)
	CreateRole(ctx context.Context, req CreateRoleRequest) (RoleResponse, error)
	GetRole(ctx context.Context, id uint) (RoleResponse, error)
	UpdateRole(ctx context.Context, id uint, req UpdateRoleRequest) (RoleResponse, error)
	DeleteRole(ctx context.Context, id uint) error
	GetRoleMenus(ctx context.Context, id uint) (RoleMenusResponse, error)
	UpdateRoleMenus(ctx context.Context, id uint, req AssignMenusRequest) error
	GetRoleAPIs(ctx context.Context, id uint) (RoleAPIsResponse, error)
	UpdateRoleAPIs(ctx context.Context, id uint, req AssignAPIsRequest) error
	ListRoleUsers(ctx context.Context, id uint, q ListRoleUsersQuery) (RoleUserPageResult, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListRolesPage(ctx context.Context, q ListRolesQuery) (RolePageResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	items, total, err := s.repo.ListPage(ctx, q)
	if err != nil {
		return RolePageResult{}, err
	}

	resp := make([]RoleResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toResponse(item))
	}
	return RolePageResult{
		Items:    resp,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func (s *service) CreateRole(ctx context.Context, req CreateRoleRequest) (RoleResponse, error) {
	role := Role{
		Name:      strings.TrimSpace(req.Name),
		Code:      strings.TrimSpace(req.Code),
		Remark:    strings.TrimSpace(req.Remark),
		Status:    req.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	created, err := s.repo.Create(ctx, role)
	if err != nil {
		if err == ErrRoleDuplicated {
			return RoleResponse{}, errcode.ErrRoleCodeExists.AsError()
		}
		return RoleResponse{}, err
	}
	return toResponse(created), nil
}

func (s *service) GetRole(ctx context.Context, id uint) (RoleResponse, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrRoleNotFound {
			return RoleResponse{}, errcode.ErrRoleNotFound.AsError()
		}
		return RoleResponse{}, err
	}
	return toResponse(role), nil
}

func (s *service) UpdateRole(ctx context.Context, id uint, req UpdateRoleRequest) (RoleResponse, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrRoleNotFound {
			return RoleResponse{}, errcode.ErrRoleNotFound.AsError()
		}
		return RoleResponse{}, err
	}
	if strings.TrimSpace(req.Code) != "" && strings.TrimSpace(req.Code) != role.Code {
		return RoleResponse{}, errcode.ErrRoleCodeImmutable.AsError()
	}

	role.Name = strings.TrimSpace(req.Name)
	role.Remark = strings.TrimSpace(req.Remark)
	role.Status = req.Status
	role.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.Update(ctx, role)
	if err != nil {
		if err == ErrRoleDuplicated {
			return RoleResponse{}, errcode.ErrRoleCodeExists.AsError()
		}
		return RoleResponse{}, err
	}
	return toResponse(updated), nil
}

func (s *service) DeleteRole(ctx context.Context, id uint) error {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrRoleNotFound {
			return errcode.ErrRoleNotFound.AsError()
		}
		return err
	}

	total, err := s.repo.CountUsers(ctx, id)
	if err != nil {
		return err
	}
	if total > 0 {
		return errcode.ErrRoleHasUsers.AsError()
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if err == ErrRoleNotFound {
			return errcode.ErrRoleNotFound.AsError()
		}
		return err
	}

	if enforcer := appcasbin.Get(); enforcer != nil {
		_, _ = enforcer.RemoveFilteredPolicy(0, role.Code)
		_, _ = enforcer.RemoveFilteredGroupingPolicy(1, role.Code)
	}
	return nil
}

func (s *service) GetRoleMenus(ctx context.Context, id uint) (RoleMenusResponse, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if err == ErrRoleNotFound {
			return RoleMenusResponse{}, errcode.ErrRoleNotFound.AsError()
		}
		return RoleMenusResponse{}, err
	}

	menus, err := s.repo.ListMenus(ctx)
	if err != nil {
		return RoleMenusResponse{}, err
	}
	checkedIDs, err := s.repo.GetRoleMenuIDs(ctx, id)
	if err != nil {
		if err == ErrRoleNotFound {
			return RoleMenusResponse{}, errcode.ErrRoleNotFound.AsError()
		}
		return RoleMenusResponse{}, err
	}
	return RoleMenusResponse{
		CheckedIDs: checkedIDs,
		Menus:      buildMenuTree(menus),
	}, nil
}

func (s *service) UpdateRoleMenus(ctx context.Context, id uint, req AssignMenusRequest) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if err == ErrRoleNotFound {
			return errcode.ErrRoleNotFound.AsError()
		}
		return err
	}
	return s.repo.ReplaceRoleMenus(ctx, id, req.MenuIDs)
}

func (s *service) GetRoleAPIs(ctx context.Context, id uint) (RoleAPIsResponse, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrRoleNotFound {
			return RoleAPIsResponse{}, errcode.ErrRoleNotFound.AsError()
		}
		return RoleAPIsResponse{}, err
	}

	resp := RoleAPIsResponse{Permissions: []APIPermission{}}
	enforcer := appcasbin.Get()
	if enforcer == nil {
		return resp, nil
	}

	policies, err := enforcer.GetFilteredPolicy(0, role.Code)
	if err != nil {
		return RoleAPIsResponse{}, err
	}
	for _, policy := range policies {
		if len(policy) < 3 {
			continue
		}
		resp.Permissions = append(resp.Permissions, APIPermission{
			Path:   policy[1],
			Method: policy[2],
		})
	}
	return resp, nil
}

func (s *service) UpdateRoleAPIs(ctx context.Context, id uint, req AssignAPIsRequest) error {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrRoleNotFound {
			return errcode.ErrRoleNotFound.AsError()
		}
		return err
	}

	enforcer := appcasbin.Get()
	if enforcer == nil {
		return nil
	}

	if _, err := enforcer.RemoveFilteredPolicy(0, role.Code); err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(req.Permissions))
	for _, permission := range req.Permissions {
		item := normalizePermission(permission.Path, permission.Method)
		key := item.Path + "#" + item.Method
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		if _, err := enforcer.AddPermissionForUser(role.Code, item.Path, item.Method); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ListRoleUsers(ctx context.Context, id uint, q ListRoleUsersQuery) (RoleUserPageResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	items, total, err := s.repo.ListRoleUsersPage(ctx, id, q)
	if err != nil {
		if err == ErrRoleNotFound {
			return RoleUserPageResult{}, errcode.ErrRoleNotFound.AsError()
		}
		return RoleUserPageResult{}, err
	}
	return RoleUserPageResult{
		Items:    items,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func buildMenuTree(menus []Menu) []MenuTreeNode {
	nodes := make(map[uint]*MenuTreeNode, len(menus))
	roots := make([]MenuTreeNode, 0)

	for _, menu := range menus {
		nodes[menu.ID] = &MenuTreeNode{
			ID:         menu.ID,
			ParentID:   menu.ParentID,
			Name:       menu.Name,
			Type:       menu.Type,
			Path:       menu.Path,
			Component:  menu.Component,
			Permission: menu.Permission,
			Sort:       menu.Sort,
			Visible:    menu.Visible,
			Status:     menu.Status,
		}
	}

	for _, menu := range menus {
		node := nodes[menu.ID]
		if menu.ParentID != nil {
			if parent, ok := nodes[*menu.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
				continue
			}
		}
		roots = append(roots, *node)
	}

	return roots
}
