package extgookit

import (
	gookit "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type (
	ConfigOptions struct {
		Files     []string
		EnvPrefix string
	}

	Config struct {
		*gookit.Config

		options *ConfigOptions
	}
)

func New(
	options *ConfigOptions,
) (*Config, error) {

	core := createGookitConfig()

	err := core.LoadFiles(options.Files...)
	if err != nil {
		return nil, err
	}

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
