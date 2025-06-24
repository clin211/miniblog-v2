// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
)

// TagStore 定义了 tag 模块在 store 层所实现的方法
type TagStore interface {
	genericstore.IStore[model.TagM]
}

// tagStore 是 TagStore 接口的实现
type tagStore struct {
	*genericstore.Store[model.TagM]
}

// 确保 tagStore 实现了 TagStore 接口
var _ TagStore = (*tagStore)(nil)

// newTagStore 创建 tagStore 的实例
func newTagStore(store *datastore) *tagStore {
	return &tagStore{
		Store: genericstore.NewStore[model.TagM](store, genericstore.NewLogger()),
	}
}
