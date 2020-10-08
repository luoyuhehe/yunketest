package thinkgo

import (
	"github.com/pkg/errors"
	"github.com/sahara-go/thinkgo/rbac"
	"sync"
)

var Rbac *rbac.Rbac
var initRbacOnce sync.Once

// 初始化数据库
func createDefaultRbac() {
	initRbacOnce.Do(func() {
		var err error
		Rbac, err = rbac.NewRbac(AppConfig.Rbac.ModelFile, DB)
		if err != nil {
			panic(errors.WithMessage(err, "权限控制初始化失败"))
		}
		err = Rbac.LoadPolicy()
		if err != nil {
			panic(errors.WithMessage(err, "加载权限失败"))
		}
	})
}
