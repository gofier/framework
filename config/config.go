package config

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var v *viper.Viper

func Init(configPath string) {
	v = viper.New()
	v.SetConfigType("env")
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.SetEnvPrefix("gofier")
	v.AutomaticEnv()
}

// Env 获取变量,如果不存在则使用默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return Get(envName, defaultValue[0])
	}
	return Get(envName)
}

// Add 设置配置
// func Add(name string, configure map[string]interface{}) {
// 	v.Set(name, configure)
// }

func Add(name string, configure interface{}) {
	v.Set(name, configure)
}

// Get 获取配置,如果不存在则使用默认值
func Get(path string, defaultValue ...interface{}) interface{} {
	if !v.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return v.Get(path)
}

// GetInterface 获取任意类型配置
func GetInterface(path string, defaultValue ...interface{}) interface{} {
	var value interface{}
	if len(defaultValue) > 0 {
		value = Get(path, defaultValue[0])
	} else {
		value = Get(path)
	}
	return value
}

// GetString 获取string类型配置
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(GetInterface(path, defaultValue...))
}

// GetInt 获取int类型配置
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(GetInterface(path, defaultValue...))
}

// GetInt64 获取int64类型配置
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(GetInterface(path, defaultValue...))
}

// GetUInt 获取uint类型配置
func GetUInt(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(GetInterface(path, defaultValue...))
}

// GetBoolean 获取boolean类型配置
func GetBoolean(path string, defaultValue ...interface{}) bool {
	var value interface{}
	if len(defaultValue) > 0 {
		value = Get(path, defaultValue[0])
	} else {
		value = Get(path)
	}
	return cast.ToBool(value)
}
