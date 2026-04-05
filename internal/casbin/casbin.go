package casbin

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist/file-adapter"
)

var (
	once     sync.Once
	enforcer *casbin.Enforcer
)

func Init(modelPath, policyPath string) (*casbin.Enforcer, error) {
	var err error
	once.Do(func() {
		m, e := model.NewModelFromFile(modelPath)
		if e != nil {
			err = e
			return
		}
		a := fileadapter.NewAdapter(policyPath)
		enforcer, e = casbin.NewEnforcer(m, a)
		if e != nil {
			err = e
		}
	})
	return enforcer, err
}

func Get() *casbin.Enforcer {
	return enforcer
}
