// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	"github.com/clin211/miniblog-v2/internal/pkg/known"
	"github.com/clin211/miniblog-v2/pkg/ipwho"
)

// RequestIDInterceptor 是一个 gRPC 拦截器，用于设置请求 ID 和客户端 IP 信息.
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 获取详细的 IP 信息
		client := ipwho.NewClient()
		details, err := client.GetHostIPDetail(ctx)
		if err != nil || details.IP == "" { // 获取客户端 IP 信息
			clientIP := extractClientIP(ctx)
			// 如果获取详细信息失败，仍然使用原始 IP
			ctx = contextx.WithClientIP(ctx, clientIP)
		} else {
			// 将详细的 IP 信息保存到 context 中
			ctx = contextx.WithClientIP(ctx, details.IP)

		}

		var requestID string
		md, _ := metadata.FromIncomingContext(ctx)

		// 从请求中获取请求 ID
		if requestIDs := md[known.XRequestID]; len(requestIDs) > 0 {
			requestID = requestIDs[0]
		}

		// 如果没有请求 ID，则生成一个新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
			md.Append(known.XRequestID, requestID)
		}

		// 将元数据设置为新的 incoming context
		ctx = metadata.NewIncomingContext(ctx, md)

		// 将请求 ID 设置到响应的 Header Metadata 中
		// grpc.SetHeader 会在 gRPC 方法响应中添加元数据（Metadata），
		// 此处将包含请求 ID 的 Metadata 设置到 Header 中。
		// 注意：grpc.SetHeader 仅设置数据，它不会立即发送给客户端。
		// Header Metadata 会在 RPC 响应返回时一并发送。
		_ = grpc.SetHeader(ctx, md)

		// 将请求 ID 添加到 ctx 中，使用已经包含 IP 信息的 context
		ctx = contextx.WithRequestID(ctx, requestID)

		// 继续处理请求
		res, err := handler(ctx, req)
		// 错误处理，附加请求 ID
		if err != nil {
			return res, errno.FromError(err).WithRequestID(requestID)
		}

		return res, nil
	}
}

// extractClientIP 从 gRPC context 中提取客户端 IP 地址，处理更复杂的情况
func extractClientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}

	// 处理不同类型的地址
	switch addr := p.Addr.(type) {
	case *net.TCPAddr:
		return addr.IP.String()
	case *net.UDPAddr:
		return addr.IP.String()
	default:
		// 尝试解析地址字符串
		addrStr := addr.String()
		if host, _, err := net.SplitHostPort(addrStr); err == nil {
			return host
		}
		// 如果解析失败，尝试直接提取 IP（可能格式不标准）
		if strings.Contains(addrStr, ":") {
			parts := strings.Split(addrStr, ":")
			if len(parts) > 0 {
				return parts[0]
			}
		}
		return addrStr
	}
}
