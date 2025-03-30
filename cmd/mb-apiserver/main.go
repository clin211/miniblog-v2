// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package main

import (
	"fmt"
	"os"

	"github.com/clin211/miniblog-v2/cmd/mb-apiserver/app"

	// 导入 automaxprocs 包，可以在程序启动时自动设置 GOMAXPROCS 配置，
	// 使其与 Linux 容器的 CPU 配额相匹配。
	// 这避免了在容器中运行时，因默认 GOMAXPROCS 值不合适导致的性能问题，
	// 确保 Go 程序能够充分利用可用的 CPU 资源，避免 CPU 浪费。
	_ "go.uber.org/automaxprocs"
)

func main() {
	fmt.Println("hello world")
	// 创建 miniblog 程序命令行
	command := app.NewMiniBlogCommand()

	// 执行命令并处理错误
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		// 如果发生错误，则退出程序
		// 返回退出码，可以使其它进程能够检测到错误
		os.Exit(1)
	}
}
