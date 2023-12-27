package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/AlgerDu/go-cli/v1"
)

type (
	PackCommandFlags struct {
		Name         string   `flag:"usage:应用名称，未设置时构建所有应用"`
		Output       string   `flag:"usage:输出路径，绝对路径或者相对于 dream-project.json 的路径"`
		TargetOSs    []string `flag:"name:os,usage:目标操作系统"`
		Type         string   `flag:"usage:打包方式，当前仅支持 docker"`
		BuildVersion string   `flag:"require,usage:使用的 build 版本号"`
	}

	PackCommand struct {
		*cli.BaseCommand
	}

	// App 编译时需要的一些信息
	AppPackContext struct {
		App        *ProjectApp
		TargetOS   string // 目标操作系统，多个时，会编译多次
		Output     string // 输出路径
		OutputName string // 输出文件名称
		Type       string // 打包方式
		Version    string
	}
)

var (
	AppPackType_Docker = "docker"

	defaultPackCommandFlags = &PackCommandFlags{
		Output: "./bin/packs",
		Type:   AppPackType_Docker,
	}

	packCommand = &PackCommand{
		BaseCommand: &cli.BaseCommand{
			Meta: &cli.CommandMeta{
				Name:  "pack",
				Usage: "打包某个或者全部的应用",
			},
			DefaultFlags: defaultPackCommandFlags,
		},
	}
)

func (cmd *PackCommand) Action(c *cli.Context) error {
	f := c.Value.(*PackCommandFlags)
	project := CliContextExts_GetProject(c)
	targetOSs := f.TargetOSs

	toPackApps := []*ProjectApp{}

	if f.Name != "" {
		app, err := project.FindApp(f.Name)
		if err != nil {
			return err
		}

		toPackApps = append(toPackApps, app)
	} else {
		toPackApps = project.Apps
	}

	if f.Type == AppPackType_Docker {
		targetOSs = append(targetOSs, TargetOS_Linux)
	}

	if len(targetOSs) == 0 {
		targetOSs = append(targetOSs, runtime.GOOS)
	}

	for _, app := range toPackApps {
		for _, targetOS := range targetOSs {
			packCtx := &AppPackContext{
				App:      app,
				TargetOS: targetOS,
				Output:   f.Output,
				Type:     f.Type,
				Version:  fmt.Sprintf("%s.%s", project.Version, f.BuildVersion),
			}

			if err := packApp(c, packCtx); err != nil {
				return err
			}
		}
	}
	return nil
}

// 打包 App ，过程中可能会修改上下文 packCtx
func packApp(cliCtx *cli.Context, packCtx *AppPackContext) error {
	if packCtx.Type != AppPackType_Docker {
		return fmt.Errorf("unsupport pack type %s", packCtx.Type)
	}

	return packAppToDocker(cliCtx, packCtx)
}

// 将应用打包为 Docker
func packAppToDocker(cliCtx *cli.Context, packCtx *AppPackContext) error {

	if packCtx.TargetOS != TargetOS_Linux {
		return nil
	}

	buildCtx := &AppBulidContext{
		App:      packCtx.App,
		TargetOS: packCtx.TargetOS,
		Output:   path.Join(packCtx.Output, "cache"),
	}
	err := buildApp(cliCtx, buildCtx)
	if err != nil {
		return err
	}

	srcDockerfile := path.Join(cliCtx.WorkDir, packCtx.App.Docker.Dockerfile)
	dstDockerfile := path.Join(cliCtx.WorkDir, buildCtx.Output, "Dockerfile")
	CopyFile(srcDockerfile, dstDockerfile)

	extEnvs := []string{}
	args := []string{
		"build",
		"-t",
		fmt.Sprintf("%s:%s", packCtx.App.Docker.ImageName, packCtx.Version),
		".",
	}

	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), extEnvs...)
	cmd.Dir = buildCtx.Output

	fmt.Println(args)

	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOut))
		return err
	}

	return nil
}
