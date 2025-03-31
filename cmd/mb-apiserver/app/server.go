package app

import (
	"encoding/json"
	"fmt"

	"github.com/clin211/miniblog-v2/cmd/mb-apiserver/app/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string // 配置文件路径

func NewMiniBlogCommand() *cobra.Command {
	// 创建默认的应用命令行选项
	opts := options.NewServerOptions()

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
			// 将 viper 中的配置解析到选项 opts 变量中.
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}

			// 对命令行选项值进行校验.
			if err := opts.Validate(); err != nil {
				return err
			}

			fmt.Printf("ServerMode from ServerOptions: %s\n", opts.JWTKey)
			fmt.Printf("ServerMode from Viper: %s\n\n", viper.GetString("jwt-key"))

			jsonData, _ := json.MarshalIndent(opts, "", "  ")
			fmt.Println(string(jsonData))

			return nil
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数。例如：./miniblog param1 param2
		Args: cobra.NoArgs,
	}
	// 初始化配置函数，在每个命令运行时调用
	cobra.OnInitialize(onInitialize)

	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the miniblog configuration file.")

	// 将 ServerOptions 中的选项绑定到命令标志
	opts.AddFlags(cmd.PersistentFlags())

	return cmd
}
