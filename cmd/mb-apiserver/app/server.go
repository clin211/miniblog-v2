package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 指定命令的名字，该名字会出现在帮助信息中
		Use: "mb-apiserver",
		// 命令的简短描述
		Short: "一个小型博客展示了开发全功能 Go 项目的最佳实践",
		// 命令的详细描述
		Long: `迷你博客系统，展示Go项目最佳实践：
简洁架构、标准目录结构、JWT认证、Casbin授权、多种测试、Makefile管理、
常用包集成(gorm/gin/cobra/zap等)、Web功能(中间件/跨域/优雅关停)、
多服务器支持(HTTP/HTTPS/gRPC)、RESTful API、OpenAPI文档、快速部署。`,
		// 命令出错时，不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("hello world")
			return nil
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./miniblog param1 param2
		Args: cobra.NoArgs,
	}

	return cmd
}
