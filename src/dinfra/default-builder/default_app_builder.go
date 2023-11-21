package defaultbuilder

import (
	di "github.com/AlgerDu/go-di/src"
	"github.com/AlgerDu/go-dream/src/dinfra"
	extgookit "github.com/AlgerDu/go-dream/src/dinfra/ext-gookit"
	extlogrus "github.com/AlgerDu/go-dream/src/dinfra/ext-logrus"
)

type (
	// 默认实现的 dinfra.AppBuilder
	DefaultAppBuilder struct {
		container          *di.Container
		serviceConfigFuncs []dinfra.ServiceConfigFunc
	}
)

func New() *DefaultAppBuilder {

	builder := &DefaultAppBuilder{
		container:          di.New(),
		serviceConfigFuncs: []dinfra.ServiceConfigFunc{},
	}

	return builder
}

func (builder *DefaultAppBuilder) ConfigService(configFuncs ...dinfra.ServiceConfigFunc) *DefaultAppBuilder {
	builder.serviceConfigFuncs = append(builder.serviceConfigFuncs, configFuncs...)
	return builder
}

func (builder *DefaultAppBuilder) Build() (dinfra.App, error) {

	for _, configFunc := range builder.serviceConfigFuncs {
		configFunc(builder.container)
	}

	return di.GetService[dinfra.App](builder.container)
}

// 使用默认的 config ；可以自行注入不同的实现
func (builder *DefaultAppBuilder) UseConfig(options *extgookit.ConfigOptions) *DefaultAppBuilder {
	di.AddInstance(builder.container, options)
	di.AddSingleton(builder.container, extgookit.New)
	return builder
}

// 使用默认的 logger ；可以自行注入不同的实现
func (builder *DefaultAppBuilder) UseLogger(options *extgookit.ConfigOptions) *DefaultAppBuilder {
	dinfra.DI_ConfigOptions(builder.container, "logger", extlogrus.NewDefaultOptions())
	di.AddSingleton(builder.container, extlogrus.New)
	return builder
}
