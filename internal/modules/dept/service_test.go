package dept

import (
	"context"
	"testing"
)

func TestCreateDeptRejectsDuplicateSiblingName(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	if _, err := svc.CreateDept(context.Background(), CreateDeptRequest{Name: "研发部", Status: 1}); err != nil {
		t.Fatalf("create dept: %v", err)
	}
	if _, err := svc.CreateDept(context.Background(), CreateDeptRequest{Name: "研发部", Status: 1}); err == nil {
		t.Fatal("expected duplicate name error")
	}
}

func TestDeleteDeptRejectsWhenHasChildren(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	parent, _ := svc.CreateDept(context.Background(), CreateDeptRequest{Name: "总部", Status: 1})
	if _, err := svc.CreateDept(context.Background(), CreateDeptRequest{ParentID: &parent.ID, Name: "财务部", Status: 1}); err != nil {
		t.Fatalf("create child: %v", err)
	}

	if err := svc.DeleteDept(context.Background(), parent.ID); err == nil {
		t.Fatal("expected child conflict error")
	}
}
