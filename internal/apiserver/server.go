// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package apiserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	handler "github.com/clin211/miniblog-v2/internal/apiserver/handler/grpc"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
	apiv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"
	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// Config 配置结构体，用于存储应用相关的配置.
// 不用 viper.Get，是因为这种方式能更加清晰的知道应用提供了哪些配置项.
type Config struct {
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	HTTPOptions *genericoptions.HTTPOptions
	GRPCOptions *genericoptions.GRPCOptions
}

// UnionServer 定义一个联合服务器。 根据 ServerMode 决定要启动的服务器类型。
type UnionServer struct {
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

// NewUnionServer 根据配置创建联合服务器。
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	lis, err := net.Listen("tcp", cfg.GRPCOptions.Addr)
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
		return nil, err
	}

	// 创建 GRPC Server 实例
	grpcsrv := grpc.NewServer()
	apiv1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())
	reflection.Register(grpcsrv)

	return &UnionServer{cfg: cfg, srv: grpcsrv, lis: lis}, nil
}

// Run 运行应用。
func (s *UnionServer) Run() error {
	// 打印一条日志，用来提示 GRPC 服务已经起来，方便排障
	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)
	// nolint: errcheck
	go s.srv.Serve(s.lis)

	//nolint: staticcheck
	dialOptions := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(s.cfg.GRPCOptions.Addr, dialOptions...)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// 设置序列化 protobuf 数据时，枚举类型的字段以数字格式输出.
			// 否则，默认会以字符串格式输出，跟枚举类型定义不一致，带来理解成本.
			UseEnumNumbers: true,
		},
	}))
	if err := apiv1.RegisterMiniBlogHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}

	log.Infow("Start to listening the incoming requests", "protocol", "http", "addr", s.cfg.HTTPOptions.Addr)
	httpsrv := &http.Server{
		Addr:    s.cfg.HTTPOptions.Addr,
		Handler: gwmux,
	}

	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
