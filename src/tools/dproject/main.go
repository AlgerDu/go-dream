package main

import (
	"fmt"
	"os"

	"github.com/AlgerDu/go-cli/v1"
)

func main() {

	builder := cli.NewBuilder("dproject").
		SetUsage("简易 golang 脚手架").
		SetPipeline(&cli.PipelineSettings{
			BeferCheckGlobalFlags: []cli.PipelineAction{
				PipelineAction_LoadProject,
			},
		}).
		AddCommand(buildCommand).
		AddCommand(packCommand)

	app := builder.Build()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Print(err)
	}
}
