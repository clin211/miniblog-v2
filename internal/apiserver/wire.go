// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"

	"github.com/clin211/miniblog-v2/internal/apiserver/biz"
	"github.com/clin211/miniblog-v2/internal/apiserver/pkg/validation"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	ginmw "github.com/clin211/miniblog-v2/internal/pkg/middleware/gin"
	"github.com/clin211/miniblog-v2/pkg/auth"
	"github.com/clin211/miniblog-v2/pkg/server"
)

func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB, // 提供数据库实例
		ProvideMongoDB,
		ProvideRedis,
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginmw.UserRetriever), new(*UserRetriever)),
		),
		auth.ProviderSet,
	)
	return nil, nil
}
