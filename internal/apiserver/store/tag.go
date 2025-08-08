// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	genericstore "github.com/clin211/miniblog-v2/pkg/store"
	"github.com/clin211/miniblog-v2/pkg/where"
)

// TagStore 定义了 tag 模块在 store 层所实现的方法
type TagStore interface {
	genericstore.IStore[model.TagM]

	// BatchGetByIDsWithCache 批量获取标签并带短 TTL 缓存
	BatchGetByIDsWithCache(ctx context.Context, ids []int32) (map[int32]*model.TagM, error)
}

// tagStore 是 TagStore 接口的实现
type tagStore struct {
	*genericstore.Store[model.TagM]
	ds *datastore
}

// 确保 tagStore 实现了 TagStore 接口
var _ TagStore = (*tagStore)(nil)

// newTagStore 创建 tagStore 的实例
func newTagStore(store *datastore) *tagStore {
	return &tagStore{
		Store: genericstore.NewStore[model.TagM](store, genericstore.NewLogger()),
		ds:    store,
	}
}

// BatchGetByIDsWithCache 批量获取标签并带短 TTL 缓存
func (s *tagStore) BatchGetByIDsWithCache(ctx context.Context, ids []int32) (map[int32]*model.TagM, error) {
	result := make(map[int32]*model.TagM)
	if len(ids) == 0 {
		return result, nil
	}

	// 去重
	uniq := make([]int32, 0, len(ids))
	seen := make(map[int32]bool, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			uniq = append(uniq, id)
		}
	}

	rdb := s.ds.Redis(ctx)
	uncached := make([]int32, 0, len(uniq))
	if rdb != nil {
		keys := make([]string, len(uniq))
		for i, id := range uniq {
			keys[i] = fmt.Sprintf("miniblog:tag:%d", id)
		}
		if vals, err := rdb.MGet(ctx, keys...).Result(); err == nil {
			for i, v := range vals {
				if v == nil {
					uncached = append(uncached, uniq[i])
					continue
				}
				var tag model.TagM
				if err := json.Unmarshal([]byte(v.(string)), &tag); err != nil {
					uncached = append(uncached, uniq[i])
					continue
				}
				result[uniq[i]] = &tag
			}
		} else {
			uncached = uniq
		}
	} else {
		uncached = uniq
	}

	if len(uncached) > 0 {
		var list []*model.TagM
		if err := s.ds.DB(ctx, where.F("id", uncached)).Find(&list).Error; err != nil {
			return nil, err
		}
		for _, t := range list {
			result[t.ID] = t
		}
		if rdb != nil && len(list) > 0 {
			pipe := rdb.Pipeline()
			for _, t := range list {
				key := fmt.Sprintf("miniblog:tag:%d", t.ID)
				if data, err := json.Marshal(t); err == nil {
					pipe.Set(ctx, key, data, 12*time.Hour)
				}
			}
			_, _ = pipe.Exec(ctx)
		}
	}

	return result, nil
}
