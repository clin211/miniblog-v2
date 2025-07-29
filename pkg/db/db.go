// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package db

import (
	"github.com/google/wire"
	redis "github.com/redis/go-redis/v9"
)

// ProviderSet is db providers.
var ProviderSet = wire.NewSet(
	NewMySQL,
	NewRedis,
	// 正确绑定接口和实现
	wire.Bind(
		new(redis.UniversalClient),
		new(*redis.Client),
	),
)
