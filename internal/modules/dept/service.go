package dept

import (
	"context"
	"strings"
	"time"

	"go-server/internal/errcode"
)

type Service interface {
	ListDepts(ctx context.Context) ([]DeptTreeNode, error)
	GetDept(ctx context.Context, id uint) (DeptResponse, error)
	CreateDept(ctx context.Context, req CreateDeptRequest) (DeptResponse, error)
	UpdateDept(ctx context.Context, id uint, req UpdateDeptRequest) (DeptResponse, error)
	DeleteDept(ctx context.Context, id uint) error
	ListDeptUsers(ctx context.Context, id uint, q ListDeptUsersQuery) (DeptUserPageResult, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListDepts(ctx context.Context) ([]DeptTreeNode, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return buildDeptTree(items), nil
}

func (s *service) GetDept(ctx context.Context, id uint) (DeptResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrDeptNotFound {
			return DeptResponse{}, errcode.ErrDeptNotFound.AsError()
		}
		return DeptResponse{}, err
	}
	return toResponse(item), nil
}

func (s *service) CreateDept(ctx context.Context, req CreateDeptRequest) (DeptResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return DeptResponse{}, err
	}
	if err := validateParent(items, 0, req.ParentID); err != nil {
		return DeptResponse{}, err
	}
	if hasSiblingName(items, 0, req.ParentID, req.Name) {
		return DeptResponse{}, errcode.ErrDeptNameExists.AsError()
	}

	item := Dept{
		ParentID:  req.ParentID,
		Name:      strings.TrimSpace(req.Name),
		Leader:    strings.TrimSpace(req.Leader),
		Phone:     strings.TrimSpace(req.Phone),
		Email:     strings.TrimSpace(req.Email),
		Sort:      req.Sort,
		Status:    req.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	created, err := s.repo.Create(ctx, item)
	if err != nil {
		return DeptResponse{}, err
	}
	return toResponse(created), nil
}

func (s *service) UpdateDept(ctx context.Context, id uint, req UpdateDeptRequest) (DeptResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == ErrDeptNotFound {
			return DeptResponse{}, errcode.ErrDeptNotFound.AsError()
		}
		return DeptResponse{}, err
	}

	items, err := s.repo.List(ctx)
	if err != nil {
		return DeptResponse{}, err
	}
	if err := validateParent(items, id, req.ParentID); err != nil {
		return DeptResponse{}, err
	}
	if hasSiblingName(items, id, req.ParentID, req.Name) {
		return DeptResponse{}, errcode.ErrDeptNameExists.AsError()
	}

	item.ParentID = req.ParentID
	item.Name = strings.TrimSpace(req.Name)
	item.Leader = strings.TrimSpace(req.Leader)
	item.Phone = strings.TrimSpace(req.Phone)
	item.Email = strings.TrimSpace(req.Email)
	item.Sort = req.Sort
	item.Status = req.Status
	item.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.Update(ctx, item)
	if err != nil {
		if err == ErrDeptNotFound {
			return DeptResponse{}, errcode.ErrDeptNotFound.AsError()
		}
		return DeptResponse{}, err
	}
	return toResponse(updated), nil
}

func (s *service) DeleteDept(ctx context.Context, id uint) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if err == ErrDeptNotFound {
			return errcode.ErrDeptNotFound.AsError()
		}
		return err
	}

	children, err := s.repo.CountChildren(ctx, id)
	if err != nil {
		return err
	}
	if children > 0 {
		return errcode.ErrDeptHasChildren.AsError()
	}

	users, err := s.repo.CountUsers(ctx, id)
	if err != nil {
		return err
	}
	if users > 0 {
		return errcode.ErrDeptHasUsers.AsError()
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if err == ErrDeptNotFound {
			return errcode.ErrDeptNotFound.AsError()
		}
		return err
	}
	return nil
}

func (s *service) ListDeptUsers(ctx context.Context, id uint, q ListDeptUsersQuery) (DeptUserPageResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}

	items, total, err := s.repo.ListDeptUsersPage(ctx, id, q)
	if err != nil {
		if err == ErrDeptNotFound {
			return DeptUserPageResult{}, errcode.ErrDeptNotFound.AsError()
		}
		return DeptUserPageResult{}, err
	}
	return DeptUserPageResult{
		Items:    items,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

func validateParent(items []Dept, currentID uint, parentID *uint) error {
	if parentID == nil {
		return nil
	}
	if *parentID == currentID {
		return errcode.ErrDeptParentInvalid.AsError()
	}

	nodes := make(map[uint]Dept, len(items))
	for _, item := range items {
		nodes[item.ID] = item
	}
	parent, ok := nodes[*parentID]
	if !ok {
		return errcode.ErrDeptParentInvalid.AsError()
	}

	for parent.ParentID != nil {
		if *parent.ParentID == currentID {
			return errcode.ErrDeptParentInvalid.AsError()
		}
		next, ok := nodes[*parent.ParentID]
		if !ok {
			break
		}
		parent = next
	}
	return nil
}

func hasSiblingName(items []Dept, currentID uint, parentID *uint, name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
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

func buildDeptTree(items []Dept) []DeptTreeNode {
	nodes := make(map[uint]*DeptTreeNode, len(items))
	roots := make([]DeptTreeNode, 0)

	for _, item := range items {
		nodes[item.ID] = &DeptTreeNode{
			ID:       item.ID,
			ParentID: item.ParentID,
			Name:     item.Name,
			Leader:   item.Leader,
			Phone:    item.Phone,
			Email:    item.Email,
			Sort:     item.Sort,
			Status:   item.Status,
		}
	}

	for _, item := range items {
		node := nodes[item.ID]
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
