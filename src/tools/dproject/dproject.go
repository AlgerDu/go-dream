package main

type (
	ProjectAppType string // 支持的应用类型

	// 应用 docker 相关的一些配置
	ProjectAppDockerSettings struct {
		Dockerfile string // 应用 dockerfile 文件的路径
		ImageName  string // 应用构建的 docker 镜像的名称
	}

	// 项目配置的应用，一个源码仓库下可以包含应用
	ProjectApp struct {
		Name       string                   // 名称
		Type       ProjectAppType           // 应用类型
		Src        string                   // 应用根路径（ main 文件的路径）
		ConfigFile string                   // 配置文件的名称，需要以 example 的形式放置在 src 目录下
		Docker     ProjectAppDockerSettings // 应用打包为 docker 时的一些配置
	}
)

const (
	PAT_Server ProjectAppType = "server" // 一般服务
	PAT_Tool   ProjectAppType = "tool"   // 工具，一般是 cli 工具

	TargetOS_Windows = "windows"
	TargetOS_Linux   = "linux"
)
