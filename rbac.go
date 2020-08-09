package thinkgo

import (
	"gitee.com/luoyusnnu/thinkgo/rbac"
	"github.com/pkg/errors"
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
