// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// defaultHomeDir 指定 miniblog 服务的默认配置文件目录
	defaultHomeDir = ".miniblog"

	// defaultConfigName 指定 miniblog 服务的默认配置文件名
	defaultConfigName = "mb-apiserver.yaml"
)

func onInitialize() {
	if configFile != "" {
		// 从命令行选项指定的配置文件中读取
		viper.SetConfigFile(configFile)
	} else {
		// 使用默认配置文件路径和名称
		for _, dir := range searchDirs() {
			// 将 dir 目录加入到配置文件搜索路径
			viper.AddConfigPath(dir)
		}

		// 设置配置文件格式为 YAML
		viper.SetConfigType("yaml")
		// 设置配置文件名称
		viper.SetConfigName(defaultConfigName)
	}

	// 读取环境变量并设置前缀
	setupEnvironmentVariables()

	// 读取配置文件。如果指定了配置文件名，则使用指定的配置文件，否则在注册的搜索路径中搜索配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read viper configuration file, err: %v", err)
	}

	// 打印当前使用的配置文件，方便调试
	log.Printf("Using config file: %s", viper.ConfigFileUsed())
}

// setupEnvironmentVariables 设置环境变量前缀
func setupEnvironmentVariables() {
	// 允许 viper 自动匹配环境变量
	viper.AutomaticEnv()

	// 设置环境变量前缀
	viper.SetEnvPrefix("MINIBLOG")

	// 设置环境变量分隔符
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}

// searchDirs 返回默认的配置文件搜索目录
func searchDirs() []string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()

	// 如果获取用户主目录失败，则打印错误信息并退出程序
	cobra.CheckErr(err)

	return []string{filepath.Join(homeDir, defaultHomeDir), "."}
}

// filePath 获取默认配置文件的完整路径.
func filePath() string {
	home, err := os.UserHomeDir()
	// 如果不能获取用户主目录，则记录错误并返回空路径
	cobra.CheckErr(err)
	return filepath.Join(home, defaultHomeDir, defaultConfigName)
}
