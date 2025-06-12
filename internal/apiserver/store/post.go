// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
)

// PostStore 定义了 post 模块在 store 层所实现的方法
type PostStore interface {
	genericstore.IStore[model.PostM]

	PostExpansion
}

// PostExpansion 定义了帖子操作的附加方法
type PostExpansion interface{}

// postStore 是 PostStore 接口的实现
type postStore struct {
	*genericstore.Store[model.PostM]
}

// 确保 postStore 实现了 PostStore 接口
var _ PostStore = (*postStore)(nil)

// newPostStore 创建 postStore 的实例.
func newPostStore(store *datastore) *postStore {
	return &postStore{
		Store: genericstore.NewStore[model.PostM](store, genericstore.NewLogger()),
	}
}
