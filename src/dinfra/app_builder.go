package dinfra

import di "github.com/AlgerDu/go-di/src"

type (
	// 依赖注入配置整个程序
	ServiceConfigFunc func(di.ServiceCollector)

	// App 构建
	AppBuilder interface {
		// 依赖注入
		ConfigService(configFuncs ...ServiceConfigFunc)

		// 构建一个全新的应用
		Build() (App, error)
	}
)
