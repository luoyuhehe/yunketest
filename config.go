package thinkgo

import (
	"fmt"
	"gitee.com/sahara-go/thinkgo/config"
	"gitee.com/sahara-go/thinkgo/log"
	"gitee.com/sahara-go/thinkgo/utils"
	"github.com/pkg/errors"
)

const (
	defaultAppConfigFile = "./config/app.yaml"
)

const (
	DevAppMode  string = "dev"
	TestAppMode string = "test"
	ProdAppMode string = "prod"
)

var appConfig = loadAppConfig()
var AppConfig = &AppSetting{}

// loadAppConfig 加载应用配置
func loadAppConfig() config.Config {
	c, err := config.NewConfig(config.FileProvider{File: defaultAppConfigFile})
	if err != nil {
		panic(fmt.Errorf("读取配置文件%s失败：%s", defaultAppConfigFile, err))
	}

	if err := c.Unmarshal(AppConfig); err != nil {
		panic(err)
	}
	log.Debug("加载配置")
	log.Debug(AppConfig)

	if err := checkAppConfig(AppConfig); err != nil {
		panic(err)
	}
	return c
}

func checkAppConfig(c *AppSetting) (err error) {
	if c.AppName == "" {
		return errors.New("缺少app_name配置")
	}

	if c.AppMode == "" || !utils.InStringArray(c.AppMode, DevAppMode, TestAppMode, ProdAppMode) {
		c.AppMode = DevAppMode
	}

	if c.Session.Enable {
		if c.Session.Store == "" {
			return errors.New("session配置错误")
		}
		if c.Session.SessionIdKey == "" {
			c.Session.SessionIdKey = "thinkgo_sessionid"
		}
	}

	return nil
}

type AppSetting struct {
	AppName         string      `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	BaseRouterGroup string      `mapstructure:"base_router_group" json:"base_router_group" yaml:"base_router_group"`
	EnvPrefix       string      `mapstructure:"env_profix" json:"env_profix" yaml:"env_profix"`
	AppMode         string      `mapstructure:"app_mode" json:"app_mode" yaml:"app_mode"`
	AutoLogResp     string      `mapstructure:"auto_log_resp" json:"auto_log_resp" yaml:"auto_log_resp"`
	Rbac            RbacSetting `mapstructure:"rbac" json:"rbac" yaml:"rbac"`

	Database DatabaseSetting `json:"database" yaml:"database"`
	Redis    RedisSetting    `json:"redis" yaml:"redis"`
	Server   ServerSetting   `json:"server" yaml:"server"`
	Captcha  CaptchaSetting  `json:"captcha" yaml:"captcha"`
	Log      LogSetting      `json:"log" yaml:"log"`
	SSL      SslSetting      `json:"ssl" yaml:"ssl"`
	Session  SessionSetting  `json:"session" yaml:"session"`
}

func (a *AppSetting) IsProdMode() bool {
	return AppConfig.AppMode == ProdAppMode
}

type ServerSetting struct {
	IsHttps      bool   `mapstructure:"is_https" json:"is_https" yaml:"is_https"`
	Env          string `json:"env" yaml:"env"`
	Addr         string `json:"addr" yaml:"addr"`
	Port         int    `json:"port" yaml:"port"`
	ReadTimeOut  int64  `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeOut int64  `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
}

type DatabaseSetting struct {
	Type  string       `json:"type" yaml:"type"`
	Mysql MysqlSetting `json:"mysql" yaml:"mysql"`
}

type MysqlSetting struct {
	Username        string `json:"username" yaml:"username"`
	Password        string `json:"password" yaml:"password"`
	Host            string `json:"host" yaml:"host"`
	Port            string `json:"port" yaml:"port"`
	DBName          string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	Charset         string `json:"charset" yaml:"charset"`
	Loc             string `json:"loc" yaml:"loc"`
	IdleNum         int    `mapstructure:"idle_num" json:"idle_num" yaml:"idle_num"`
	PoolNum         int    `mapstructure:"pool_num" json:"pool_num" yaml:"pool_num"`
	MaxLifeSecond   int    `mapstructure:"max_life_second" json:"max_life_second" yaml:"max_life_second"`
	LogMode         bool   `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	MultiStatements bool   `mapstructure:"multi_statements" json:"multi_statements" yaml:"multi_statements"`
	ParseTime       bool   `mapstructure:"parse_time" json:"parse_time" yaml:"parse_time"`
}

type RedisSetting struct {
	Host        string `json:"host" yaml:"host"`
	Password    string `json:"password" yaml:"password"`
	DB          int    `json:"db" yaml:"db"`
	PoolNum     int    `mapstructure:"pool_num" json:"pool_num" yaml:"pool_num"`
	IdleTimeout int    `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
}

type CaptchaSetting struct {
	Store string                   `mapstructure:"store" json:"store" yaml:"store"`
	Redis CaptchaRedisStoreSetting `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type CaptchaRedisStoreSetting struct {
	Host        string `json:"host" yaml:"host"`
	Password    string `json:"password" yaml:"password"`
	DB          int    `json:"db" yaml:"db"`
	PoolNum     int    `mapstructure:"pool_num" json:"pool_num" yaml:"pool_num"`
	IdleTimeout int    `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
	KeyPrefix   string `mapstructure:"key_prefix" json:"key_prefix" yaml:"key_prefix"`
	Expire      int    `mapstructure:"expire" json:"expire" yaml:"expire"`
}

type LogSetting struct {
	File    FileLogSetting    `json:"file" yaml:"file"`
	Console ConsoleLogSetting `json:"console" yaml:"console"`
	Kafka   KafkaLogSetting   `json:"kafka" yaml:"kafka"`
}

type FileLogSetting struct {
	Enable     bool   `json:"enable" yaml:"enable"`
	Prefix     string `json:"prefix" yaml:"prefix"`
	Path       string `json:"path" yaml:"path"`
	Level      string `json:"level" yaml:"level"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

type ConsoleLogSetting struct {
	Enable bool   `json:"enable" yaml:"enable"`
	Prefix string `json:"prefix" yaml:"prefix"`
	Level  string `json:"level" yaml:"level"`
}

type KafkaLogSetting struct {
	Enable bool     `json:"enable" yaml:"enable"`
	Prefix string   `json:"prefix" yaml:"prefix"`
	Level  string   `json:"level" yaml:"level"`
	Addr   []string `json:"addr" yaml:"addr"`
	Topic  string   `json:"topic" yaml:"topic"`
}

type SessionSetting struct {
	Enable       bool              `json:"enable" yaml:"enable"`
	SessionIdKey string            `mapstructure:"sessionid_key" json:"sessionid_key" yaml:"sessionid_key"`
	Store        string            `json:"store" yaml:"store"`
	Redis        RedisStoreSetting `json:"redis" yaml:"redis"`
}

type RedisStoreSetting struct {
	Enable      bool   `json:"enable" yaml:"enable"`
	Host        string `json:"host" yaml:"host"`
	MaxIdleConn int    `mapstructure:"max_idle_conn" json:"max_idle_conn" yaml:"max_idle_conn"`
	Password    string `json:"password" yaml:"password"`
	Key         string `json:"key" yaml:"key"`
}

// GetKey redis key
func (rs *RedisStoreSetting) GetKey() string {
	if rs.Key == "" {
		return "thinkgo-session"
	}
	return rs.Key
}

type SslSetting struct {
	KeyFile  string `mapstructure:"key_file" json:"key_file" yaml:"key_file"`
	CertFile string `mapstructure:"cert_file" json:"cert_file" yaml:"cert_file"`
}

type RbacSetting struct {
	Enable    bool   `json:"enable" yaml:"enable"`
	ModelFile string `mapstructure:"model_file" json:"model_file" yaml:"model_file"`
}
