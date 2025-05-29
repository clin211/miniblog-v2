// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

// Authorizer 用于定义授权接口的实现.
type Authorizer interface {
	Authorize(subject, object, action string) (bool, error)
}

// AuthzInterceptor 是一个 gRPC 拦截器，用于进行请求授权.
func AuthzInterceptor(authorizer Authorizer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		subject := contextx.UserID(ctx) // 获取用户ID
		object := info.FullMethod       // 获取请求资源
		action := "CALL"                // 默认操作

		// 记录授权上下文信息
		log.Debugw("Build authorize context", "subject", subject, "object", object, "action", action)

		// 调用授权接口进行验证
		if allowed, err := authorizer.Authorize(subject, object, action); err != nil || !allowed {
			return nil, errno.ErrPermissionDenied.WithMessage(
				"access denied: subject=%s, object=%s, action=%s, reason=%v",
				subject,
				object,
				action,
				err,
			)
		}

		// 继续处理请求
		return handler(ctx, req)
	}
}
