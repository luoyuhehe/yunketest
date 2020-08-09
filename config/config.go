package config

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"time"
)

// 配置接口
type Config interface {
	Unmarshal(rawVal interface{}) error
	WatchConfig()
	OnConfigChange(callback func())
	Get(key string) interface{}
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetIntSlice(key string) []int
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	IsSet(key string) bool
	Set(key string, value interface{})
	AllSettings() map[string]interface{}
	Sub(key string) Config
	UnmarshalKey(key string, rawVal interface{}) error
}

type config struct {
	*viper.Viper
}

// NewConfig 构建配置实例
func NewConfig(provider interface{}) (Config, error) {
	//加载配置文件
	v := viper.New()
	// 配置存储在本地文件
	if p, ok := provider.(FileProvider); ok {
		v.SetConfigFile(p.File)
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				return nil, errors.New("config file not found")
			}

			return nil, err
		}
	}

	// 配置存储在etcd
	if p, ok := provider.(EtcdProvider); ok {
		err := v.AddSecureRemoteProvider(p.GetType(), p.EndPoint, p.Path, p.SecretKey)
		if err != nil {
			return nil, err
		}
		v.SetConfigType(p.ConfigType) // 因为在字节流中没有文件扩展名，所以这里需要设置下类型。支持的扩展名有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
		if err := v.ReadRemoteConfig(); err != nil {
			return nil, err
		}
	}

	// 配置存储在etcd
	if p, ok := provider.(ConsulProvider); ok {
		err := v.AddSecureRemoteProvider(p.GetType(), p.Addr, p.Key, p.SecretKey)
		if err != nil {
			return nil, err
		}
		v.SetConfigType(p.ConfigType) // 因为在字节流中没有文件扩展名，所以这里需要设置下类型。支持的扩展名有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
		if err := v.ReadRemoteConfig(); err != nil {
			return nil, err
		}
	}

	return &config{v}, nil
}

// Unmarshal 配置解析到结构体
func (c *config) Unmarshal(rawVal interface{}) error {
	if err := c.Viper.Unmarshal(rawVal); err != nil {
		return errors.New(fmt.Sprintf("config file parse err:%s", err))
	}

	return nil
}

// UnmarshalKey 配置解析到结构体
func (c *config) UnmarshalKey(key string, rawVal interface{}) error {
	if err := c.Viper.UnmarshalKey(key, rawVal); err != nil {
		return err
	}

	return nil
}

// OnConfigChange 回调函数
func (c *config) OnConfigChange(callback func()) {
	c.Viper.OnConfigChange(func(e fsnotify.Event) {
		callback()
	})
}

func (c *config) Sub(key string) Config {
	v := c.Viper.Sub(key)
	return &config{v}
}

//////////////////////  Viper 原有的方法  /////////////////////////////////////

func (c *config) Get(key string) interface{} {
	return c.Viper.Get(key)
}

func (c *config) GetString(key string) string {
	return c.Viper.GetString(key)
}

func (c *config) GetStringMap(key string) map[string]interface{} {
	return c.GetStringMap(key)
}

func (c *config) GetStringMapString(key string) map[string]string {
	return c.Viper.GetStringMapString(key)
}

func (c *config) GetStringSlice(key string) []string {
	return c.Viper.GetStringSlice(key)
}

func (c *config) GetBool(key string) bool {
	return c.Viper.GetBool(key)
}

func (c *config) GetInt(key string) int {
	return c.Viper.GetInt(key)
}

func (c *config) GetIntSlice(key string) []int {
	return c.Viper.GetIntSlice(key)
}

func (c *config) GetTime(key string) time.Time {
	return c.Viper.GetTime(key)
}

func (c *config) GetDuration(key string) time.Duration {
	return c.Viper.GetDuration(key)
}

func (c *config) IsSet(key string) bool {
	return c.Viper.IsSet(key)
}

func (c *config) Set(key string, value interface{}) {
	c.Viper.Set(key, value)
}

func (c *config) AllSettings() map[string]interface{} {
	return c.Viper.AllSettings()
}
