package config

import (
	"fmt"

	"github.com/someview/trade-base/util"
	"github.com/spf13/viper"
)

// ConfigPath 返回指定配置文件的完整路径
// fileName: 配置文件名
// 返回值: 配置文件的完整路径
func ConfigPath(fileName string) string {
	path := util.AppPath()

	return path + "/" + fileName
}

// Configer 是一个接口，用于检查配置的有效性
type Configer interface {
	// CheckValidity 检查配置的有效性
	// 返回值: 如果配置有效返回nil，否则返回错误
	CheckValidity() error
}

// NewConfigRooter 创建一个新的Viper配置解析器
// configPath: 配置文件的路径
// 返回值: Viper配置解析器和可能的错误
func NewConfigRooter(configPath string) (rooter *viper.Viper, err error) {
	rooter = viper.New()
	rooter.SetConfigFile(configPath)
	err = rooter.ReadInConfig()
	return
}

// GetConfig 从Viper配置解析器中获取指定键的配置并解析到结构体
// rooter: Viper配置解析器
// key: 配置键
// t: 配置结构体实例
// 返回值: 解析后的配置结构体和可能的错误
func GetConfig[T Configer](rooter *viper.Viper, key string, t T) (T, error) {
	// 直接解析 key 到结构体
	if err := rooter.UnmarshalKey(key, t); err != nil {
		return t, fmt.Errorf("解析配置失败: %w", err)
	}

	// 检查配置有效性
	if err := t.CheckValidity(); err != nil {
		return t, fmt.Errorf("配置校验失败: %w", err)
	}
	return t, nil
}
