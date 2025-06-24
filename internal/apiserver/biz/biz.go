// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package biz

import (
	"github.com/google/wire"

	postv1 "github.com/clin211/miniblog-v2/internal/apiserver/biz/v1/post"
	tagv1 "github.com/clin211/miniblog-v2/internal/apiserver/biz/v1/tag"
	userv1 "github.com/clin211/miniblog-v2/internal/apiserver/biz/v1/user"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	"github.com/clin211/miniblog-v2/pkg/auth"
	// Post V2 版本（未实现，仅展示用）
	// postv2 "github.com/clin211/miniblog-v2/internal/apiserver/biz/v2/post".
)

// ProviderSet 是一个 Wire 的 Provider 集合，用于声明依赖注入的规则.Add commentMore actions
// 包含 NewBiz 构造函数，用于生成 biz 实例.
// wire.Bind 用于将接口 IBiz 与具体实现 *biz 绑定，
// 这样依赖 IBiz 的地方会自动注入 *biz 实例.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

// IBiz 定义了业务层需要实现的方法.
type IBiz interface {
	// 获取用户业务接口.
	UserV1() userv1.UserBiz
	// 获取帖子业务接口.
	PostV1() postv1.PostBiz
	// 获取标签业务接口.
	TagV1() tagv1.TagBiz
	// 获取帖子业务接口（V2版本）.
	// PostV2() post.PostBiz
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	store store.IStore
	authz *auth.Authz
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(store store.IStore, authz *auth.Authz) *biz {
	return &biz{store: store, authz: authz}
}

// UserV1 返回一个实现了 UserBiz 接口的实例.
func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store, b.authz)
}

// PostV1 返回一个实现了 PostBiz 接口的实例.
func (b *biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}

// TagV1 返回一个实现了 TagBiz 接口的实例.
func (b *biz) TagV1() tagv1.TagBiz {
	return tagv1.New(b.store)
}
