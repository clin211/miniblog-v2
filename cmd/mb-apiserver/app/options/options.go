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
	genericoptions "github.com/clin211/miniblog-v2/pkg/options"
	"github.com/clin211/miniblog-v2/pkg/strings"
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
	// TLSOptions 包含 TLS 配置选项.
	TLSOptions *genericoptions.TLSOptions `json:"tls" mapstructure:"tls"`
	// HTTPOptions 包含 HTTP 配置选项.
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`
	// GRPCOptions 包含 gRPC 配置选项.
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
	// MySQLOptions 包含 MySQL 配置选项.
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	// MongoOptions 包含 MongoDB 选项
	MongoOptions *genericoptions.MongoOptions `json:"mongodb" mapstructure:"mongo"`
	// RedisOptions 包含 Redis 配置选项
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	// UploadOptions 包含文件上传配置选项
	UploadOptions *genericoptions.UploadOptions `json:"upload" mapstructure:"upload"`
}

// NewServerOptions 创建带有默认值的 ServerOptions 实例
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:    apiserver.GRPCGatewayServerMode,
		JWTKey:        "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:    2 * time.Hour,
		TLSOptions:    genericoptions.NewTLSOptions(),
		HTTPOptions:   genericoptions.NewHTTPOptions(),
		GRPCOptions:   genericoptions.NewGRPCOptions(),
		MySQLOptions:  genericoptions.NewMySQLOptions(),
		MongoOptions:  genericoptions.NewMongoOptions(),
		RedisOptions:  genericoptions.NewRedisOptions(),
		UploadOptions: genericoptions.NewUploadOptions(),
	}
	opts.HTTPOptions.Addr = ":5555"
	opts.GRPCOptions.Addr = ":6666"
	return opts
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")

	// 绑定 JWT Token 的过期时间选项到命令行标志
	// 参数名称 `--expiration`，默认值为 o.Expiration
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "JWT expiration")
	o.TLSOptions.AddFlags(fs)
	o.HTTPOptions.AddFlags(fs)
	o.GRPCOptions.AddFlags(fs)
	o.MySQLOptions.AddFlags(fs)
	o.MongoOptions.AddFlags(fs)
	o.RedisOptions.AddFlags(fs)
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

	// 校验子选项
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.MongoOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)

	// 如果是 gRPC 或 gRPC-Gateway 模式，校验 gRPC 配置
	if strings.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}

	// 合并所有错误并返回
	return utilerrors.NewAggregate(errs)
}

// Config 基于 ServerOptions 创建新的 apiserver.Config。
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:    o.ServerMode,
		JWTKey:        o.JWTKey,
		TLSOptions:    o.TLSOptions,
		Expiration:    o.Expiration,
		HTTPOptions:   o.HTTPOptions,
		GRPCOptions:   o.GRPCOptions,
		MySQLOptions:  o.MySQLOptions,
		MongoOptions:  o.MongoOptions,
		RedisOptions:  o.RedisOptions,
		UploadOptions: o.UploadOptions,
	}, nil
}
