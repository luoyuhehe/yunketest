package rbac

import (
	"gitee.com/sahara-go/thinkgo/model"
	rbac_adapter "gitee.com/sahara-go/thinkgo/rbac/adapter"
	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
)

type Rbac struct {
	enforcer *casbin.Enforcer
}

func NewRbac(modelFile string, db *gorm.DB) (*Rbac, error) {
	adapter, err := rbac_adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(modelFile, adapter)
	if err != nil {
		return nil, err
	}

	return &Rbac{
		enforcer: enforcer,
	}, nil
}

// LoadByDB
func (r *Rbac) LoadPolicy() error {
	// Load the policy from DB.
	err := r.enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	return nil
}

// Save the policy back to DB
func (r *Rbac) Save() error {
	return r.enforcer.SavePolicy()
}

// Check check the permission
func (r *Rbac) Check(rvals interface{}) (bool, error) {
	return r.enforcer.Enforce(rvals)
}

// AddFunction 添加自定义规则函数
func (r *Rbac) AddFunction(name string, f func(args ...interface{}) (interface{}, error)) {
	r.enforcer.AddFunction(name, f)
}

// @title AddPolicy
// @description 添加权限
// @auth luoyu
// @param
// @return bool
func (r *Rbac) AddPolicy(cm model.CasbinModel) (bool, error) {
	return r.enforcer.AddPolicy(cm.AuthorityId, cm.Path, cm.Method)
}

// @title RemovePolicy
// @description 删除权限
// @auth luoyu
// @param
// @return bool
func (r *Rbac) RemovePolicy(cm model.CasbinModel) (bool, error) {
	return r.enforcer.RemovePolicy(cm.AuthorityId, cm.Path, cm.Method)
}
