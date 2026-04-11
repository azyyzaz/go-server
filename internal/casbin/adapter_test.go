package casbin

import (
	"testing"

	"github.com/casbin/casbin/v2/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newAdapterTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&Rule{}); err != nil {
		t.Fatalf("migrate casbin rules: %v", err)
	}
	return db
}

func TestAdapterLoadPolicy(t *testing.T) {
	db := newAdapterTestDB(t)
	if err := db.Create(&[]Rule{
		newRule("p", []string{"admin", "/api/v1/system/roles", "GET"}),
		newRule("g", []string{"alice", "admin"}),
	}).Error; err != nil {
		t.Fatalf("seed rules: %v", err)
	}

	adapter := NewAdapter(db)
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

	if err := adapter.LoadPolicy(m); err != nil {
		t.Fatalf("load policy: %v", err)
	}

	if got := m["p"]["p"].Policy; len(got) != 1 || got[0][0] != "admin" || got[0][1] != "/api/v1/system/roles" {
		t.Fatalf("unexpected p policies: %#v", got)
	}
	if got := m["g"]["g"].Policy; len(got) != 1 || got[0][0] != "alice" || got[0][1] != "admin" {
		t.Fatalf("unexpected g policies: %#v", got)
	}
}

func TestAdapterAddAndRemoveFilteredPolicy(t *testing.T) {
	db := newAdapterTestDB(t)
	adapter := NewAdapter(db)

	if err := adapter.AddPolicy("p", "p", []string{"admin", "/api/v1/system/roles", "GET"}); err != nil {
		t.Fatalf("add policy: %v", err)
	}
	if err := adapter.AddPolicy("g", "g", []string{"alice", "admin"}); err != nil {
		t.Fatalf("add grouping: %v", err)
	}
	if err := adapter.RemoveFilteredPolicy("p", "p", 0, "admin"); err != nil {
		t.Fatalf("remove policy: %v", err)
	}
	if err := adapter.RemoveFilteredPolicy("g", "g", 1, "admin"); err != nil {
		t.Fatalf("remove grouping: %v", err)
	}

	var count int64
	if err := db.Model(&Rule{}).Count(&count).Error; err != nil {
		t.Fatalf("count rules: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 rules, got %d", count)
	}
}
