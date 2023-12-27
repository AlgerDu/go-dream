package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/AlgerDu/go-cli/v1"
)

const (
	DefaultProjectJsonName = "dream-project.json" // tool 默认使用的项目配置文件的名称
)

var (
	ErrAppNotExist = errors.New("app not exist")
)

type (
	// 整个项目的配置
	Project struct {
		ToolVersion string        // 使用的 dproject 工具的版本
		Version     string        // 项目的版本
		Apps        []*ProjectApp // 项目包含的应用
	}
)

// 通过 app 名称查找配置中的 app
func (project *Project) FindApp(name string) (*ProjectApp, error) {
	for _, app := range project.Apps {
		if app.Name == name {
			return app, nil
		}
	}

	return nil, ErrAppNotExist
}

// 加载运行目录下的 project 配置文件
func PipelineAction_LoadProject(c *cli.Context) error {
	fileName := path.Join(c.WorkDir, DefaultProjectJsonName)
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	project := &Project{}
	if err = json.Unmarshal(fileContents, project); err != nil {
		return err
	}

	CliContextExts_SetProject(c, project)

	return nil
}

func CliContextExts_SetProject(c *cli.Context, project *Project) {
	c.Items["project"] = project
}

func CliContextExts_GetProject(c *cli.Context) *Project {
	return c.Items["project"].(*Project)
}
