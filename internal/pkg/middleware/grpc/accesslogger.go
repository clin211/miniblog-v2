// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package grpc

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

const (
	// MaxRequestLogSize 请求日志记录的最大大小 (10KB)
	MaxRequestLogSize = 10 * 1024
	// MaxResponseLogSize 响应日志记录的最大大小 (10KB)
	MaxResponseLogSize = 10 * 1024
	// DataTooLargeMessage 数据过大时的提示信息
	DataTooLargeMessage = "Data size is too large to log"
)

// AccessLogger 是一个 gRPC 拦截器，用于记录访问日志
func AccessLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		// 获取客户端IP
		clientIP := getClientIP(ctx)

		// 获取请求ID
		requestID := contextx.RequestID(ctx)

		// 获取用户ID（可能从metadata中获取）
		userID := getUserID(ctx)

		// 序列化请求参数
		reqBody := serializeRequest(req)

		// 记录访问开始
		accessLog(ctx, "access_start", 0, info.FullMethod, clientIP, userID, requestID, reqBody, "", 0)

		// 执行请求处理
		resp, err := handler(ctx, req)

		// 计算执行时间
		duration := time.Since(start)

		// 序列化响应
		respBody := serializeResponse(resp)

		// 获取状态码（从错误中推断）
		status := getStatusCode(err)

		// 记录访问结束
		accessLog(ctx, "access_end", duration, info.FullMethod, clientIP, userID, requestID, reqBody, respBody, status)

		return resp, err
	}
}

// getClientIP 从 context 中获取客户端 IP 地址
func getClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return "unknown"
}

// getUserID 从 context 或 metadata 中获取用户ID
func getUserID(ctx context.Context) string {
	// 尝试从 metadata 中获取用户ID
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userIDs := md.Get("user-id"); len(userIDs) > 0 {
			return userIDs[0]
		}
		if userIDs := md.Get("x-user-id"); len(userIDs) > 0 {
			return userIDs[0]
		}
	}

	// 也可以尝试从 context 中获取（如果有存储的话）
	// 这里可以根据具体的认证方式来扩展
	return "unknown"
}

// serializeRequest 序列化请求参数
func serializeRequest(req any) string {
	if req == nil {
		return ""
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "failed to serialize request"
	}

	if len(data) > MaxRequestLogSize {
		return DataTooLargeMessage
	}

	return string(data)
}

// serializeResponse 序列化响应数据
func serializeResponse(resp any) string {
	if resp == nil {
		return ""
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return "failed to serialize response"
	}

	if len(data) > MaxResponseLogSize {
		return DataTooLargeMessage
	}

	return string(data)
}

// getStatusCode 从错误中获取状态码
func getStatusCode(err error) int {
	if err == nil {
		return 0 // gRPC OK
	}

	// 这里可以根据具体的错误类型来返回对应的状态码
	// 例如使用 status.FromError(err) 来获取 gRPC 状态码
	return -1 // 表示有错误，具体状态码可以进一步处理
}

// accessLog 记录访问日志
func accessLog(ctx context.Context, accessType string, dur time.Duration, method, ip, userID, requestID, reqBody, respBody string, status int) {
	log.Infow("gRPC AccessLog",
		"type", accessType,
		"ip", ip,
		"userID", userID,
		"requestID", requestID,
		"method", method,
		"request", reqBody,
		"response", respBody,
		"time(ms)", int64(dur/time.Millisecond),
		"status", status,
	)
}
