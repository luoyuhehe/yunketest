package thinkgo

import (
	"fmt"
	"gitee.com/sahara-go/thinkgo/db"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"sync"
)

var DB *gorm.DB
var initDbOnce sync.Once

// 初始化数据库
func createDatabaseConnPool() {
	initDbOnce.Do(func() {
		var err error
		switch AppConfig.Database.Type {
		case "mysql":
			setting := AppConfig.Database.Mysql
			DB, err = db.NewMysqlDB(db.MysqlSetting{
				Host:            setting.Host,
				Port:            setting.Port,
				Username:        setting.Username,
				Password:        setting.Password,
				DBName:          setting.DBName,
				Charset:         setting.Charset,
				PoolNum:         setting.PoolNum,
				IdleNum:         setting.IdleNum,
				LogMod:          setting.LogMode,
				Loc:             setting.Loc,
				MaxLifeSecond:   setting.MaxLifeSecond,
				MultiStatements: setting.MultiStatements,
				ParseTime:       setting.ParseTime,
			})
		default:
			panic(errors.New(fmt.Sprintf("不支持此数据库类型：%s", AppConfig.Database.Type)))
		}

		if err != nil {
			panic(err)
		}

		// 设置回调函数
		db.SetCreateCallback(DB, db.UpdateTimeStampForCreateCallback)
		db.SetUpdateCallback(DB, db.UpdateTimeStampForUpdateCallback)
		db.SetDeleteCallback(DB, db.DeleteCallback)
	})
}
