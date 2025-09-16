package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfigPath = "./config.test.yaml"

// 定义一个简单的配置结构体用于测试
type TestConfig struct {
	Key string `mapstructure:"key"`
}

func (t *TestConfig) CheckValidity() error {
	if t.Key == "" {
		return fmt.Errorf("key is empty")
	}
	return nil
}

func TestNewConfigRooter(t *testing.T) {
	// 测试 NewConfigRooter 函数
	configPath := testConfigPath
	rooter, err := NewConfigRooter(configPath)
	assert.NoError(t, err)
	assert.NotNil(t, rooter)
}

func TestGetConfig(t *testing.T) {
	// 测试 GetConfig 函数
	configPath := testConfigPath
	rooter, err := NewConfigRooter(configPath)
	assert.NoError(t, err)

	// 测试配置键不存在的情况
	var testConfig TestConfig
	_, err = GetConfig(rooter, "test_key", &testConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "解析配置失败")

	// 测试正常情况
	testConfig.Key = "value"
	config, err := GetConfig(rooter, "server_config", &testConfig)
	assert.NoError(t, err)
	assert.Equal(t, "value", config.Key)

	// todo mock testConfig的函数
	//testConfig.Key = ""
	//_, err = GetConfig(rooter, "server_config", &testConfig)
	//assert.Error(t, err)
	//assert.Equal(t, "配置校验失败: key is empty", err.Error())
}

func TestNestedConfig(t *testing.T) {
	var conf Config
	rooter, err := NewConfigRooter("config.dev.yaml")
	if err != nil {
		t.Log("加载配置文件发生错误")
		return
	}
	if err := rooter.UnmarshalKey(TestConfigKey, &conf); err != nil {
		t.Log("解析策略部分发生错误", err)
	}
	t.Log(conf)
}
