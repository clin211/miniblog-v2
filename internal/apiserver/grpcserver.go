// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package apiserver

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	handler "github.com/clin211/miniblog-v2/internal/apiserver/handler/grpc"
	mw "github.com/clin211/miniblog-v2/internal/pkg/middleware/grpc"
	apiv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/clin211/miniblog-v2/pkg/server"
)

// grpcServer 定义一个 gRPC 服务器.
type grpcServer struct {
	srv server.Server
	// stop 为优雅关停函数.
	stop func(context.Context)
}

// 确保 *grpcServer 实现了 server.Server 接口.
var _ server.Server = (*grpcServer)(nil)

// NewGRPCServerOr 创建并初始化 gRPC 或者 gRPC +  gRPC-Gateway 服务器.
// 在 Go 项目开发中，NewGRPCServerOr 这个函数命名中的 Or 一般用来表示“或者”的含义，
// 通常暗示该函数会在两种或多种选择中选择一种可能性。具体的含义需要结合函数的实现
// 或上下文来理解。以下是一些可能的解释：
//  1. 提供多种构建方式的选择
//  2. 处理默认值或回退逻辑
//  3. 表达灵活选项
func (c *ServerConfig) NewGRPCServerOr() (server.Server, error) {
	// 配置 gRPC 服务器选项，包括拦截器链
	serverOptions := []grpc.ServerOption{
		// 注意拦截器顺序！
		grpc.ChainUnaryInterceptor(
			// 请求 ID 拦截器
			mw.RequestIDInterceptor(),
			// Bypass 拦截器，通过所有请求的认证
			mw.AuthnBypasswInterceptor(),
			// 请求默认值设置拦截器
			mw.DefaulterInterceptor(),
		),
	}

	// 创建 gRPC 服务器
	grpcsrv, err := server.NewGRPCServer(
		c.cfg.GRPCOptions,
		serverOptions,
		func(s grpc.ServiceRegistrar) {
			apiv1.RegisterMiniBlogServer(s, handler.NewHandler(c.biz))
		},
	)
	if err != nil {
		return nil, err
	}

	if c.cfg.ServerMode == GRPCServerMode {
		return &grpcServer{
			srv: grpcsrv,
			stop: func(ctx context.Context) {
				grpcsrv.GracefulStop(ctx)
			},
		}, nil
	}

	// 先启动 gRPC 服务器，因为 HTTP 服务器依赖 gRPC 服务器.
	go grpcsrv.RunOrDie()

	httpsrv, err := server.NewGRPCGatewayServer(
		c.cfg.HTTPOptions,
		c.cfg.GRPCOptions,
		func(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return apiv1.RegisterMiniBlogHandler(context.Background(), mux, conn)
		},
	)
	if err != nil {
		return nil, err
	}

	return &grpcServer{
		srv: httpsrv,
		stop: func(ctx context.Context) {
			grpcsrv.GracefulStop(ctx)
			httpsrv.GracefulStop(ctx)
		},
	}, nil
}

// RunOrDie 启动 gRPC 服务器或 HTTP 反向代理服务器，异常时退出.
func (s *grpcServer) RunOrDie() {
	s.srv.RunOrDie()
}

// GracefulStop 优雅停止 HTTP 和 gRPC 服务器.
func (s *grpcServer) GracefulStop(ctx context.Context) {
	s.stop(ctx)
}
