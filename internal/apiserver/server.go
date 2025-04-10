// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package apiserver

import (
	"time"

	"github.com/spf13/viper"

	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

// Config 配置结构体，用于存储应用相关的配置。
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项。
type Config struct {
	ServerMode string
	JWTKey     string
	Expiration time.Duration
}

// UnionServer 定义一个联合服务器。 根据 ServerMode 决定要启动的服务器类型。
type UnionServer struct {
	cfg *Config
}

// NewUnionServer 根据配置创建联合服务器。
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	return &UnionServer{cfg: cfg}, nil
}

// Run 运行应用。
func (s *UnionServer) Run() error {
	log.Infow("ServerMode from ServerOptions: %s\n", s.cfg.JWTKey)
	log.Infow("ServerMode from Viper: %s\n\n", viper.GetString("jwt-key"))

	select {}
}
