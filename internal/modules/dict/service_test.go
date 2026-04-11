package dict

import (
	"context"
	"testing"
)

func TestCreateTypeRejectsDuplicateCode(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo, nil)

	if _, err := svc.CreateType(context.Background(), CreateDictTypeRequest{Name: "状态", Code: "sys_status", Status: 1}); err != nil {
		t.Fatalf("create type: %v", err)
	}
	if _, err := svc.CreateType(context.Background(), CreateDictTypeRequest{Name: "状态2", Code: "sys_status", Status: 1}); err == nil {
		t.Fatal("expected duplicate code error")
	}
}

func TestLookupByTypeCodeReturnsOnlyActiveItems(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo, nil)

	typ, _ := svc.CreateType(context.Background(), CreateDictTypeRequest{Name: "状态", Code: "sys_status", Status: 1})
	if _, err := svc.CreateData(context.Background(), CreateDictDataRequest{TypeID: typ.ID, Label: "启用", Value: "1", Status: 1}); err != nil {
		t.Fatalf("create active data: %v", err)
	}
	if _, err := svc.CreateData(context.Background(), CreateDictDataRequest{TypeID: typ.ID, Label: "禁用", Value: "0", Status: 0}); err != nil {
		t.Fatalf("create inactive data: %v", err)
	}

	items, err := svc.LookupByTypeCode(context.Background(), "sys_status")
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	if len(items) != 1 || items[0].Value != "1" {
		t.Fatalf("unexpected lookup result: %#v", items)
	}
}
