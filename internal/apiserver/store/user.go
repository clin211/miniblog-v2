// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
)

// UserStore 定义了 user 模块在 store 层所实现的方法
type UserStore interface {
	genericstore.IStore[model.UserM]
}

// userStore 是 UserStore 接口的实现
type userStore struct {
	*genericstore.Store[model.UserM]
}

// 确保 userStore 实现了 UserStore 接口
var _ UserStore = (*userStore)(nil)

// newUserStore 创建 userStore 的实例
func newUserStore(store *datastore) *userStore {
	return &userStore{
		Store: genericstore.NewStore[model.UserM](store, genericstore.NewLogger()),
	}
}
