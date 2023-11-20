package dinfra

import di "github.com/AlgerDu/go-di/src"

type (
	Config interface {
		BindStruct(path string, dst any) error
	}

	ConfigChangeListener interface {
		OnChange() error
	}
)

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
