package casbin

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/persist/file-adapter"
	"gorm.io/gorm"
)

type Rule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:16;not null;index"`
	V0    string `gorm:"size:255;index"`
	V1    string `gorm:"size:255;index"`
	V2    string `gorm:"size:255;index"`
	V3    string `gorm:"size:255"`
	V4    string `gorm:"size:255"`
	V5    string `gorm:"size:255"`
}

func (Rule) TableName() string {
	return "casbin_rule"
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(db *gorm.DB) *Adapter {
	return &Adapter{db: db}
}

func (a *Adapter) LoadPolicy(m model.Model) error {
	var rules []Rule
	if err := a.db.Order("id ASC").Find(&rules).Error; err != nil {
		return err
	}
	for _, rule := range rules {
		if err := persist.LoadPolicyLine(rule.line(), m); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) SavePolicy(m model.Model) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Rule{}).Error; err != nil {
			return err
		}

		var rules []Rule
		for ptype, ast := range m["p"] {
			for _, policy := range ast.Policy {
				rules = append(rules, newRule(ptype, policy))
			}
		}
		for ptype, ast := range m["g"] {
			for _, policy := range ast.Policy {
				rules = append(rules, newRule(ptype, policy))
			}
		}

		if len(rules) == 0 {
			return nil
		}
		return tx.Create(&rules).Error
	})
}

func (a *Adapter) AddPolicy(_ string, ptype string, rule []string) error {
	item := newRule(ptype, rule)
	return a.db.Create(&item).Error
}

func (a *Adapter) RemovePolicy(_ string, ptype string, rule []string) error {
	return a.applyFilter(a.db, ptype, 0, rule...).Delete(&Rule{}).Error
}

func (a *Adapter) RemoveFilteredPolicy(_ string, ptype string, fieldIndex int, fieldValues ...string) error {
	return a.applyFilter(a.db, ptype, fieldIndex, fieldValues...).Delete(&Rule{}).Error
}

func (a *Adapter) applyFilter(db *gorm.DB, ptype string, fieldIndex int, fieldValues ...string) *gorm.DB {
	q := db.Where("ptype = ?", ptype)
	columns := []string{"v0", "v1", "v2", "v3", "v4", "v5"}
	for idx, value := range fieldValues {
		colIdx := fieldIndex + idx
		if colIdx >= len(columns) {
			break
		}
		if value == "" {
			continue
		}
		q = q.Where(fmt.Sprintf("%s = ?", columns[colIdx]), value)
	}
	return q
}

func newRule(ptype string, values []string) Rule {
	rule := Rule{Ptype: ptype}
	fields := []*string{&rule.V0, &rule.V1, &rule.V2, &rule.V3, &rule.V4, &rule.V5}
	for idx, value := range values {
		if idx >= len(fields) {
			break
		}
		*fields[idx] = value
	}
	return rule
}

func (r Rule) line() string {
	parts := []string{r.Ptype, r.V0, r.V1, r.V2, r.V3, r.V4, r.V5}
	last := len(parts) - 1
	for last > 0 && parts[last] == "" {
		last--
	}
	return strings.Join(parts[:last+1], ", ")
}

var (
	once     sync.Once
	enforcer *casbin.Enforcer
)

func Init(modelPath, policyPath string, db *gorm.DB) (*casbin.Enforcer, error) {
	var err error
	once.Do(func() {
		m, e := model.NewModelFromFile(modelPath)
		if e != nil {
			err = e
			return
		}

		var adapter any
		if db != nil {
			adapter = NewAdapter(db)
		} else {
			adapter = fileadapter.NewAdapter(policyPath)
		}

		enforcer, e = casbin.NewEnforcer(m, adapter)
		if e != nil {
			err = e
			return
		}
		enforcer.EnableAutoSave(true)
	})
	return enforcer, err
}

func Get() *casbin.Enforcer {
	return enforcer
}

func Set(e *casbin.Enforcer) {
	enforcer = e
}

func CleanupPoliciesByRole(ctx context.Context, db *gorm.DB, roleCode string) error {
	return db.WithContext(ctx).
		Where("ptype = ? AND v0 = ?", "p", roleCode).
		Or("ptype = ? AND v1 = ?", "g", roleCode).
		Delete(&Rule{}).Error
}
