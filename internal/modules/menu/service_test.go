package menu

import (
	"context"
	"testing"
)

func TestCreateMenuRejectsDuplicateSiblingName(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	if _, err := svc.CreateMenu(context.Background(), CreateMenuRequest{
		Name: "System", Type: "directory", Visible: 1, Status: 1,
	}); err != nil {
		t.Fatalf("create menu: %v", err)
	}

	if _, err := svc.CreateMenu(context.Background(), CreateMenuRequest{
		Name: "System", Type: "directory", Visible: 1, Status: 1,
	}); err == nil {
		t.Fatal("expected duplicate name error")
	}
}

func TestDeleteMenuRejectsWhenHasChildren(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	parent, err := svc.CreateMenu(context.Background(), CreateMenuRequest{
		Name: "System", Type: "directory", Visible: 1, Status: 1,
	})
	if err != nil {
		t.Fatalf("create parent: %v", err)
	}
	if _, err := svc.CreateMenu(context.Background(), CreateMenuRequest{
		ParentID: &parent.ID, Name: "Users", Type: "menu", Visible: 1, Status: 1,
	}); err != nil {
		t.Fatalf("create child: %v", err)
	}

	if err := svc.DeleteMenu(context.Background(), parent.ID); err == nil {
		t.Fatal("expected child conflict error")
	}
}

func TestUpdateMenuSorts(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	first, _ := svc.CreateMenu(context.Background(), CreateMenuRequest{Name: "A", Type: "menu", Visible: 1, Status: 1})
	second, _ := svc.CreateMenu(context.Background(), CreateMenuRequest{Name: "B", Type: "menu", Visible: 1, Status: 1})

	if err := svc.UpdateMenuSorts(context.Background(), UpdateMenuSortsRequest{
		Items: []MenuSortItem{{ID: first.ID, Sort: 20}, {ID: second.ID, Sort: 10}},
	}); err != nil {
		t.Fatalf("update sorts: %v", err)
	}

	items, err := svc.ListMenus(context.Background())
	if err != nil {
		t.Fatalf("list menus: %v", err)
	}
	if len(items) != 2 || items[0].ID != second.ID {
		t.Fatalf("unexpected menu order: %#v", items)
	}
}
