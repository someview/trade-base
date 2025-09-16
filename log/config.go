package log

import (
	"path/filepath"
	"time"

	"github.com/someview/trade-base/config"
)

const LogConfigKey = "log_config"

// Config 日志配置
type Config struct {
	LogLevel               LogLevel      `mapstructure:"log_level"`
	LogPath                string        `mapstructure:"log_path"`
	MaxSize                int64         `mapstructure:"rotate_size_mb"`
	BufferSize             int           `mapstructure:"buffer_size"`
	FlushInterval          time.Duration `mapstructure:"flush_interval"`
	CloseWithoutCompressed bool          `mapstructure:"close_without_compressed"`
	CompressTime           string        `mapstructure:"compress_time"` // 新增：每日定时打印时间（HH:MM）
}

func NewDefaultConfig() Config {
	return Config{
		LogLevel:               DebugLevel,
		LogPath:                filepath.Join(LogPath() + "/logs/"),
		MaxSize:                512 * 1024 * 1024,
		BufferSize:             16384,
		FlushInterval:          time.Second,
		CloseWithoutCompressed: false,
		CompressTime:           "00:00", // 默认每天0点打印日志
	}
}

func (c *Config) CheckValidity() error {
	return nil
}

var _ config.Configer = (*Config)(nil)
