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
	BuildCommandFlags struct {
		Name      string   `flag:"usage:应用名称，为空时构建所有应用"`
		Output    string   `flag:"usage:输出路径，绝对路径或者相对于 dream-project.json 的路径"`
		TargetOSs []string `flag:"name:os,usage:输出路径，绝对路径或者相对于 dream-project.json 的路径"`
	}

	BuildCommand struct {
		*cli.BaseCommand
	}

	// App 编译时需要的一些信息
	AppBulidContext struct {
		App        *ProjectApp
		TargetOS   string // 目标操作系统，多个时，会编译多次
		Output     string // 输出路径
		OutputName string // 输出文件名称
	}
)

var (
	defaultHelloCommandFlags = &BuildCommandFlags{
		Output: "./bin",
	}

	buildCommand = &BuildCommand{
		BaseCommand: &cli.BaseCommand{
			Meta: &cli.CommandMeta{
				Name:  "build",
				Usage: "编译全部或者某个应用",
			},
			DefaultFlags: defaultHelloCommandFlags,
		},
	}
)

func (cmd *BuildCommand) Action(c *cli.Context) error {
	f := c.Value.(*BuildCommandFlags)
	project := CliContextExts_GetProject(c)
	toBuildTargetOSs := f.TargetOSs

	toBuildApps := []*ProjectApp{}

	if f.Name != "" {
		app, err := project.FindApp(f.Name)
		if err != nil {
			return err
		}

		toBuildApps = append(toBuildApps, app)
	} else {
		toBuildApps = project.Apps
	}

	if len(toBuildTargetOSs) == 0 {
		toBuildTargetOSs = append(toBuildTargetOSs, runtime.GOOS)
	}

	for _, app := range toBuildApps {
		for _, targetOS := range toBuildTargetOSs {
			buildCtx := &AppBulidContext{
				App:      app,
				TargetOS: targetOS,
				Output:   f.Output,
			}

			if err := buildApp(c, buildCtx); err != nil {
				return err
			}
		}
	}
	return nil
}

// 编译 App ，过程中可能会修改上下文 buildCtx
func buildApp(cliCtx *cli.Context, buildCtx *AppBulidContext) error {
	extEnvs := []string{}
	if buildCtx.OutputName == "" {
		if buildCtx.App.Type == PAT_Tool {
			buildCtx.OutputName = buildCtx.App.Name
		} else {
			buildCtx.OutputName = "main"
		}
	}

	if buildCtx.TargetOS == TargetOS_Windows {
		buildCtx.OutputName = fmt.Sprintf("%s.exe", buildCtx.OutputName)
	}

	buildCtx.Output = path.Join(buildCtx.Output, buildCtx.App.Name)

	if buildCtx.TargetOS != runtime.GOOS {
		buildCtx.Output = fmt.Sprintf("%s/%s", buildCtx.Output, buildCtx.TargetOS)
		extEnvs = append(extEnvs, fmt.Sprintf("GOOS=%s", buildCtx.TargetOS))
	}

	args := []string{
		"build",
		"-o",
		path.Join(cliCtx.WorkDir, buildCtx.Output, buildCtx.OutputName),
		path.Join(cliCtx.WorkDir, buildCtx.App.Src),
	}

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), extEnvs...)
	cmd.Dir = cliCtx.WorkDir

	fmt.Println(args)

	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(cmdOut))
		return err
	}

	if buildCtx.App.ConfigFile != "" {
		example := path.Join(cliCtx.WorkDir, buildCtx.App.Src, fmt.Sprintf("%s.example", buildCtx.App.ConfigFile))
		dst := path.Join(cliCtx.WorkDir, buildCtx.Output, buildCtx.App.ConfigFile)

		err = copyConfigFile(example, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

// 向目标目录复制配置文件，如果文件存在，则跳过
func copyConfigFile(src, dst string) error {

	if PathExists(dst) {
		return nil
	}

	return CopyFile(src, dst)
}
