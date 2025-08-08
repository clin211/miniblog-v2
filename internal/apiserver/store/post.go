// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"context"
	"strconv"
	"time"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
	"github.com/clin211/miniblog-v2/pkg/where"
	"github.com/redis/go-redis/v9"
)

// PostStore 定义了 post 模块在 store 层所实现的方法
type PostStore interface {
	genericstore.IStore[model.PostM]

	// ListApp 返回应用层列表（无 Count），按 id desc 排序
	ListApp(ctx context.Context, opts *where.Options) ([]*model.PostM, error)
	// CountApp 返回应用层列表的总数，带短 TTL 缓存
	CountApp(ctx context.Context, opts *where.Options) (int64, error)
}

// postStore 是 PostStore 接口的实现
type postStore struct {
	*genericstore.Store[model.PostM]
	ds *datastore
}

// 确保 postStore 实现了 PostStore 接口
var _ PostStore = (*postStore)(nil)

// newPostStore 创建 postStore 的实例.
func newPostStore(store *datastore) *postStore {
	return &postStore{
		Store: genericstore.NewStore[model.PostM](store, genericstore.NewLogger()),
		ds:    store,
	}
}

// ListApp 返回应用层列表（无 Count），按 id desc 排序
func (s *postStore) ListApp(ctx context.Context, opts *where.Options) ([]*model.PostM, error) {
	var posts []*model.PostM
	if err := s.ds.DB(ctx, opts).Order("id desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// CountApp 返回应用层列表的总数，带短 TTL 缓存
func (s *postStore) CountApp(ctx context.Context, opts *where.Options) (int64, error) {
	// 构造 cache key（仅受 status、category_id 影响）
	var status, category string
	if v, ok := opts.Filters["status"]; ok {
		status = strconv.FormatInt(int64(v.(int32)), 10)
	}
	if v, ok := opts.Filters["category_id"]; ok {
		category = strconv.FormatInt(int64(v.(int32)), 10)
	}
	key := "miniblog:count:posts:status:" + status + ":category:" + category

	rdb := s.ds.Redis(ctx)
	if rdb != nil {
		if _, ok := interface{}(rdb).(*redis.Client); ok {
			if val, err := rdb.Get(ctx, key).Result(); err == nil {
				if n, err2 := strconv.ParseInt(val, 10, 64); err2 == nil {
					return n, nil
				}
			}
		}
	}

	var n int64
	if err := s.ds.DB(ctx, opts).Model(&model.PostM{}).Count(&n).Error; err != nil {
		return 0, err
	}
	if rdb != nil {
		_ = rdb.Set(ctx, key, strconv.FormatInt(n, 10), 30*time.Second).Err()
	}
	return n, nil
}
