// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package options

import (
	"errors"
	"fmt"
	"time"

	"github.com/clin211/miniblog-v2/internal/apiserver"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
	stringsutil "github.com/onexstack/onexstack/pkg/util/strings"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// 定义支持的服务器模式集合
var availableServerModes = sets.New(
	apiserver.GinServerMode,
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
)

type ServerOptions struct {
	// ServerMode 定义服务器模式 gRPC、Gin HTTP、HTTP Reverse Proxy
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`
	// JWTKey 定义 JWT 秘钥
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key"`
	// Expiration 定义 JWT Token 过期时间
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`
	// GRPCOptions 包含 gRPC 配置选项.
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
}

// NewServerOptions 创建带有默认值的 ServerOptions 实例
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:  "grpc-gateway",
		JWTKey:      "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:  time.Hour * 2,
		GRPCOptions: genericoptions.NewGRPCOptions(),
	}
	opts.GRPCOptions.Addr = ":6666"
	return opts
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")

	// 绑定 JWT Token 的过期时间选项到命令行标志
	// 参数名称 `--expiration`，默认值为 o.Expiration
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "JWT expiration")
	o.GRPCOptions.AddFlags(fs)
}

// Validate 检验 ServerOptions 中的选项是否合法
func (o *ServerOptions) Validate() error {
	errs := []error{}

	// 校验 ServerMode 是否有效
	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: %s", availableServerModes.UnsortedList()))
	}

	// 校验 JWTKey 长度
	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("JWTKey must be at least 6 characters long"))
	}

	// 如果是 gRPC 或 gRPC-Gateway 模式，校验 gRPC 配置
	if stringsutil.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}

	// 合并所有错误并返回
	return utilerrors.NewAggregate(errs)
}

// Config 基于 ServerOptions 创建新的 apiserver.Config。
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:  o.ServerMode,
		JWTKey:      o.JWTKey,
		Expiration:  o.Expiration,
		GRPCOptions: o.GRPCOptions,
	}, nil
}
