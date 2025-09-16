package base

type IModule interface {

	// Run 启动模块，在初始化完成后调用
	Start() error

	// Stop 停止模块，在关闭前调用，用于优雅关闭
	Stop(reason string) error

	// GetName 获取模块名称，用于日志和调试
	GetName() string
}

type IStrategy interface {
	IModule
}

type Event interface {
	EventType() string
}
