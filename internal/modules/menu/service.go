package menu

import (
	"context"
	"strings"
	"time"

	"go-server/internal/errcode"
)

type Service interface {
	ListMenus(ctx context.Context) ([]MenuTreeNode, error)
	GetMenu(ctx context.Context, id uint) (MenuResponse, error)
	CreateMenu(ctx context.Context, req CreateMenuRequest) (MenuResponse, error)
	UpdateMenu(ctx context.Context, id uint, req UpdateMenuRequest) (MenuResponse, error)
	DeleteMenu(ctx context.Context, id uint) error
	UpdateMenuSorts(ctx context.Context, req UpdateMenuSortsRequest) error
	ListCurrentUserMenus(ctx context.Context, userID uint) ([]MenuTreeNode, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListMenus(ctx context.Context) ([]MenuTreeNode, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(items, false), nil
}

func (s *service) GetMenu(ctx context.Context, id uint) (MenuResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrMenuNotFound {
			return MenuResponse{}, errcode.ErrMenuNotFound.AsError()
		}
		return MenuResponse{}, err
	}
	return toResponse(item), nil
}

func (s *service) CreateMenu(ctx context.Context, req CreateMenuRequest) (MenuResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return MenuResponse{}, err
	}
	if err := validateParent(items, 0, req.ParentID); err != nil {
		return MenuResponse{}, err
	}
	if hasSiblingName(items, 0, req.ParentID, req.Name) {
		return MenuResponse{}, errcode.ErrMenuNameExists.AsError()
	}

	item := Menu{
		ParentID:   req.ParentID,
		Name:       strings.TrimSpace(req.Name),
		Type:       req.Type,
		Path:       strings.TrimSpace(req.Path),
		Component:  strings.TrimSpace(req.Component),
		Permission: strings.TrimSpace(req.Permission),
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     req.Status,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	created, err := s.repo.Create(ctx, item)
	if err != nil {
		return MenuResponse{}, err
	}
	return toResponse(created), nil
}

func (s *service) UpdateMenu(ctx context.Context, id uint, req UpdateMenuRequest) (MenuResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrMenuNotFound {
			return MenuResponse{}, errcode.ErrMenuNotFound.AsError()
		}
		return MenuResponse{}, err
	}

	items, err := s.repo.List(ctx)
	if err != nil {
		return MenuResponse{}, err
	}
	if err := validateParent(items, id, req.ParentID); err != nil {
		return MenuResponse{}, err
	}
	if hasSiblingName(items, id, req.ParentID, req.Name) {
		return MenuResponse{}, errcode.ErrMenuNameExists.AsError()
	}

	item.ParentID = req.ParentID
	item.Name = strings.TrimSpace(req.Name)
	item.Type = req.Type
	item.Path = strings.TrimSpace(req.Path)
	item.Component = strings.TrimSpace(req.Component)
	item.Permission = strings.TrimSpace(req.Permission)
	item.Sort = req.Sort
	item.Visible = req.Visible
	item.Status = req.Status
	item.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.Update(ctx, item)
	if err != nil {
		if err == ErrMenuNotFound {
			return MenuResponse{}, errcode.ErrMenuNotFound.AsError()
		}
		return MenuResponse{}, err
	}
	return toResponse(updated), nil
}

func (s *service) DeleteMenu(ctx context.Context, id uint) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if err == ErrMenuNotFound {
			return errcode.ErrMenuNotFound.AsError()
		}
		return err
	}

	total, err := s.repo.CountChildren(ctx, id)
	if err != nil {
		return err
	}
	if total > 0 {
		return errcode.ErrMenuHasChildren.AsError()
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if err == ErrMenuNotFound {
			return errcode.ErrMenuNotFound.AsError()
		}
		return err
	}
	return nil
}

func (s *service) UpdateMenuSorts(ctx context.Context, req UpdateMenuSortsRequest) error {
	if err := s.repo.UpdateSorts(ctx, req.Items); err != nil {
		if err == ErrMenuNotFound {
			return errcode.ErrMenuNotFound.AsError()
		}
		return err
	}
	return nil
}

func (s *service) ListCurrentUserMenus(ctx context.Context, userID uint) ([]MenuTreeNode, error) {
	items, err := s.repo.ListCurrentUserMenus(ctx, userID)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(items, true), nil
}

func validateParent(items []Menu, currentID uint, parentID *uint) error {
	if parentID == nil {
		return nil
	}
	if *parentID == currentID {
		return errcode.ErrMenuParentInvalid.AsError()
	}

	nodes := make(map[uint]Menu, len(items))
	for _, item := range items {
		nodes[item.ID] = item
	}
	parent, ok := nodes[*parentID]
	if !ok {
		return errcode.ErrMenuParentInvalid.AsError()
	}

	for parent.ParentID != nil {
		if *parent.ParentID == currentID {
			return errcode.ErrMenuParentInvalid.AsError()
		}
		next, ok := nodes[*parent.ParentID]
		if !ok {
			break
		}
		parent = next
	}
	return nil
}

func hasSiblingName(items []Menu, currentID uint, parentID *uint, name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, item := range items {
		if item.ID == currentID {
			continue
		}
		if strings.ToLower(strings.TrimSpace(item.Name)) != name {
			continue
		}
		if sameParent(item.ParentID, parentID) {
			return true
		}
	}
	return false
}

func sameParent(a, b *uint) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func buildMenuTree(items []Menu, excludeButtons bool) []MenuTreeNode {
	nodes := make(map[uint]*MenuTreeNode, len(items))
	roots := make([]MenuTreeNode, 0)

	for _, item := range items {
		if excludeButtons && item.Type == "button" {
			continue
		}
		nodes[item.ID] = &MenuTreeNode{
			ID:         item.ID,
			ParentID:   item.ParentID,
			Name:       item.Name,
			Type:       item.Type,
			Path:       item.Path,
			Component:  item.Component,
			Permission: item.Permission,
			Sort:       item.Sort,
			Visible:    item.Visible,
			Status:     item.Status,
		}
	}

	for _, item := range items {
		node, ok := nodes[item.ID]
		if !ok {
			continue
		}
		if item.ParentID != nil {
			if parent, ok := nodes[*item.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
				continue
			}
		}
		roots = append(roots, *node)
	}

	return roots
}
