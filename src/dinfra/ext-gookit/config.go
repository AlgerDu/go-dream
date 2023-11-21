package extgookit

import (
	"github.com/AlgerDu/go-dream/src/dinfra"
	gookit "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type (
	ConfigOptions struct {
		Files     []string
		EnvPrefix string
		Envs      []string
	}

	Config struct {
		*gookit.Config

		options *ConfigOptions
	}
)

func New(
	logger dinfra.Logger,
	options *ConfigOptions,
) (*Config, error) {

	core := createGookitConfig()

	err := core.LoadFiles(options.Files...)
	if err != nil {
		return nil, err
	}

	logger.WithField(dinfra.LogField_Source, "extgookit.New").Info("create gookit config")

	return &Config{
		Config:  core,
		options: options,
	}, nil
}

func createGookitConfig() *gookit.Config {
	cfg := gookit.NewWithOptions("core", gookit.ParseEnv)
	cfg.AddDriver(yaml.Driver)
	return cfg
}
