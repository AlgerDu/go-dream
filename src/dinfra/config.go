package dinfra

import di "github.com/AlgerDu/go-di/src"

type (
	OnChangeFunc func() error

	// 配置
	Config interface {
		// 将配置绑定到结构体
		BindStruct(path string, dst any) error
	}

	// 各个服务的可选参数
	Options interface {
		// option 的使用者可以通过注册 listener 来响应 config 的运行时变更
		OnChange(listener OnChangeFunc)
	}
)

// 将 config path 对应的数据绑定给 options 然后注入 IoC 容器
func DI_ConfigOptions[insType any](
	services di.ServiceCollector,
	path string,
	ins insType,
	opts ...func(insType),
) {
	di.AddScope(services, func(config Config) insType {
		config.BindStruct(path, ins)
		if len(opts) > 0 {
			for _, opt := range opts {
				opt(ins)
			}
		}
		return ins
	})
}
