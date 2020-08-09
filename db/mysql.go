package db

import (
	"bytes"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
	"time"
)

type MysqlSetting struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	Charset         string
	PoolNum         int //数据库的最大打开连接数，gorm包当前默认值为2
	IdleNum         int
	Loc             string
	MaxLifeSecond   int
	MultiStatements bool
	ParseTime       bool
	LogMod          bool
}

// NewMysqlDB returns *gorm.DB instance.
func NewMysqlDB(setting MysqlSetting) (*gorm.DB, error) {
	if setting.Username == "" {
		return nil, errors.New("username配置不能为空")
	}
	if setting.Password == "" {
		return nil, errors.New("password配置不能为空")
	}
	if setting.Host == "" {
		setting.Host = "localhost"
	}
	if setting.DBName == "" {
		return nil, errors.New("dbname配置不能为空")
	}
	if setting.Charset == "" {
		setting.Charset = "utf8"
	}
	if setting.IdleNum == 0 {
		//gorm包当前默认值为2
		setting.PoolNum = 2
	}

	var buf bytes.Buffer
	buf.WriteString(setting.Username)
	buf.WriteString(":")
	buf.WriteString(setting.Password)
	buf.WriteString("@tcp(")
	buf.WriteString(setting.Host)
	buf.WriteString(":")
	buf.WriteString(setting.Port)
	buf.WriteString(")/")
	buf.WriteString(setting.DBName)
	buf.WriteString("?charset=")
	buf.WriteString(setting.Charset)
	buf.WriteString("&parseTime=" + strconv.FormatBool(setting.ParseTime))
	buf.WriteString("&multiStatements=" + strconv.FormatBool(setting.MultiStatements))
	if setting.Loc == "" {
		buf.WriteString("&loc=Local")
	} else {
		buf.WriteString("&loc=" + url.QueryEscape(setting.Loc))
	}

	db, err := gorm.Open("mysql", buf.String())
	if err != nil {
		return nil, errors.WithMessage(err, "mysql数据库连接失败")
	}
	db.LogMode(setting.LogMod)

	if setting.MaxLifeSecond > 0 {
		db.DB().SetConnMaxLifetime(time.Duration(setting.MaxLifeSecond) * time.Second)
	}
	db.DB().SetMaxIdleConns(setting.IdleNum)
	db.DB().SetMaxOpenConns(setting.PoolNum)

	return db, nil
}
