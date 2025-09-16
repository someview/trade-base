package global

import (
	"github.com/someview/trade-base/config"
	"github.com/someview/trade-base/util"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	IsProxy         bool   `mapstructure:"is_proxy"`
	ProxyUrl        string `mapstructure:"proxy_url"`
	ServerId        string `mapstructure:"server_id"`
	EnableSafeCheck bool   `mapstructure:"enable_safe_check"` // 是否开启安全检测，//是否进行安全验证，如果为false,则ExpireTime、ValidIpList设置无效
	Env             string `mapstructure:"env"`
}

func (s *ServerConfig) CheckValidity() error {
	return nil
}

var (
	IsProxy         = false
	ProxyUrl        string
	ServerId        string
	EnableSafeCheck = false
	localIps        []string
	env             = "test"
)

func init() {
	var err error
	localIps, err = util.GetAllLocalIPs()
	if err != nil {
		panic(err)
	}
}

func GetLocalIPs() []string {
	return localIps
}

func InitServerConfig(rooter *viper.Viper) error {
	conf, err := config.GetConfig(rooter, "server_config", &ServerConfig{})
	if err != nil {
		return err
	}
	IsProxy = conf.IsProxy
	ProxyUrl = conf.ProxyUrl
	ServerId = conf.ServerId
	EnableSafeCheck = conf.EnableSafeCheck
	if conf.Env != "" {
		env = conf.Env
	}
	return nil
}

func IsTestEnv() bool {
	return env == "test"
}

func IsProdEnv() bool {
	return env == "prod"
}
