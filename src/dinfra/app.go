package dinfra

type (
	// 应用程序
	App interface {
		// 启动
		Run() error
	}
)
