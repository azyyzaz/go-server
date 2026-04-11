package role

import (
	"context"
	"testing"

	appcasbin "go-server/internal/casbin"
	"go-server/internal/response"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

func newTestEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()

	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`)
	if err != nil {
		t.Fatalf("new model: %v", err)
	}
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("new enforcer: %v", err)
	}
	return e
}

func TestCreateRoleRejectsDuplicateCode(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	if _, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Admin", Code: "admin", Status: 1}); err != nil {
		t.Fatalf("create role: %v", err)
	}

	_, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Admin 2", Code: "admin", Status: 1})
	if err == nil {
		t.Fatal("expected duplicate code error")
	}

	appErr, ok := err.(*response.AppError)
	if !ok || appErr.Code != "ROLE_CODE_EXISTS" {
		t.Fatalf("expected ROLE_CODE_EXISTS, got %v", err)
	}
}

func TestUpdateRoleRejectsCodeChange(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	created, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Admin", Code: "admin", Status: 1})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	_, err = svc.UpdateRole(context.Background(), created.ID, UpdateRoleRequest{Name: "Admin Updated", Code: "super-admin", Status: 1})
	if err == nil {
		t.Fatal("expected immutable code error")
	}
}

func TestDeleteRoleReturnsNotFound(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	if err := svc.DeleteRole(context.Background(), 42); err == nil {
		t.Fatal("expected not found error")
	}
}

func TestUpdateRoleMenusReplacesMenuIDs(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	created, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Editor", Code: "editor", Status: 1})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	if err := svc.UpdateRoleMenus(context.Background(), created.ID, AssignMenusRequest{MenuIDs: []uint{1, 2, 3}}); err != nil {
		t.Fatalf("update menus: %v", err)
	}
	if err := svc.UpdateRoleMenus(context.Background(), created.ID, AssignMenusRequest{MenuIDs: []uint{2}}); err != nil {
		t.Fatalf("replace menus: %v", err)
	}

	result, err := svc.GetRoleMenus(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get role menus: %v", err)
	}
	if len(result.CheckedIDs) != 1 || result.CheckedIDs[0] != 2 {
		t.Fatalf("expected checked ids [2], got %#v", result.CheckedIDs)
	}
}

func TestUpdateRoleAPIsReplacesPolicies(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)
	enforcer := newTestEnforcer(t)
	appcasbin.Set(enforcer)
	defer appcasbin.Set(nil)

	created, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Manager", Code: "manager", Status: 1})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	if err := svc.UpdateRoleAPIs(context.Background(), created.ID, AssignAPIsRequest{
		Permissions: []APIPermission{{Path: "/api/v1/system/users", Method: "GET"}},
	}); err != nil {
		t.Fatalf("first update apis: %v", err)
	}
	if err := svc.UpdateRoleAPIs(context.Background(), created.ID, AssignAPIsRequest{
		Permissions: []APIPermission{{Path: "/api/v1/system/roles", Method: "POST"}},
	}); err != nil {
		t.Fatalf("second update apis: %v", err)
	}

	policies, err := enforcer.GetFilteredPolicy(0, "manager")
	if err != nil {
		t.Fatalf("get policies: %v", err)
	}
	if len(policies) != 1 || policies[0][1] != "/api/v1/system/roles" || policies[0][2] != "POST" {
		t.Fatalf("unexpected policies: %#v", policies)
	}
}
