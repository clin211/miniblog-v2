// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
)

// PostTagStore 定义了 post_tag 模块在 store 层所实现的方法
type PostTagStore interface {
	genericstore.IStore[model.PostTagM]
}

// postTagStore 是 PostTagStore 接口的实现
type postTagStore struct {
	*genericstore.Store[model.PostTagM]
}

// 确保 postTagStore 实现了 PostTagStore 接口
var _ PostTagStore = (*postTagStore)(nil)

// newPostTagStore 创建 postTagStore 的实例
func newPostTagStore(store *datastore) *postTagStore {
	return &postTagStore{
		Store: genericstore.NewStore[model.PostTagM](store, genericstore.NewLogger()),
	}
}
