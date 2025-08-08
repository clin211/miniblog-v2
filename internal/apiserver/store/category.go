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
	"golang.org/x/sync/singleflight"
)

// CategoryStore 定义了 category 模块在 store 层所实现的方法
type CategoryStore interface {
	genericstore.IStore[model.CategoryM]

	// BatchGetByIDsWithCache 批量获取分类并带短 TTL 缓存
	BatchGetByIDsWithCache(ctx context.Context, ids []int32) (map[int32]*model.CategoryM, error)

	// ListAllWithCache 返回全量分类列表并带缓存
	ListAllWithCache(ctx context.Context) ([]*model.CategoryM, error)
	// ListActiveWithCache 返回 is_active=1 的分类列表并带缓存
	ListActiveWithCache(ctx context.Context) ([]*model.CategoryM, error)
}

// categoryStore 是 CategoryStore 接口的实现
type categoryStore struct {
	*genericstore.Store[model.CategoryM]
	ds *datastore
}

// 确保 categoryStore 实现了 CategoryStore 接口
var _ CategoryStore = (*categoryStore)(nil)

// newCategoryStore 创建 categoryStore 的实例
func newCategoryStore(store *datastore) *categoryStore {
	return &categoryStore{
		Store: genericstore.NewStore[model.CategoryM](store, genericstore.NewLogger()),
		ds:    store,
	}
}

const (
	cacheKeyCategoryListAll    = "miniblog:category:list:all"
	cacheKeyCategoryListActive = "miniblog:category:list:active"
	cacheTTLCategoryList       = 12 * time.Hour
)

var listGroup singleflight.Group

func (s *categoryStore) invalidateListCache(ctx context.Context) {
	rdb := s.ds.Redis(ctx)
	if rdb == nil {
		return
	}
	_ = rdb.Del(ctx, cacheKeyCategoryListAll, cacheKeyCategoryListActive).Err()
}

// Create 覆盖通用 Create，在成功后失效列表缓存
func (s *categoryStore) Create(ctx context.Context, data *model.CategoryM) error {
	if err := s.Store.Create(ctx, data); err != nil {
		return err
	}
	s.invalidateListCache(ctx)
	return nil
}

// Update 覆盖通用 Update，在成功后失效列表缓存及单条缓存
func (s *categoryStore) Update(ctx context.Context, data *model.CategoryM) error {
	if err := s.Store.Update(ctx, data); err != nil {
		return err
	}
	rdb := s.ds.Redis(ctx)
	if rdb != nil && data != nil && data.ID != 0 {
		_ = rdb.Del(ctx, fmt.Sprintf("miniblog:category:%d", data.ID)).Err()
	}
	s.invalidateListCache(ctx)
	return nil
}

// Delete 覆盖通用 Delete，在成功后失效列表缓存
func (s *categoryStore) Delete(ctx context.Context, opts *where.Options) error {
	if err := s.Store.Delete(ctx, opts); err != nil {
		return err
	}
	s.invalidateListCache(ctx)
	return nil
}

// BatchGetByIDsWithCache 批量获取分类并带短 TTL 缓存
func (s *categoryStore) BatchGetByIDsWithCache(ctx context.Context, ids []int32) (map[int32]*model.CategoryM, error) {
	result := make(map[int32]*model.CategoryM)
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
			keys[i] = fmt.Sprintf("miniblog:category:%d", id)
		}
		if vals, err := rdb.MGet(ctx, keys...).Result(); err == nil {
			for i, v := range vals {
				if v == nil {
					uncached = append(uncached, uniq[i])
					continue
				}
				var cat model.CategoryM
				if err := json.Unmarshal([]byte(v.(string)), &cat); err != nil {
					uncached = append(uncached, uniq[i])
					continue
				}
				result[uniq[i]] = &cat
			}
		} else {
			// 缓存不可用则全部回源
			uncached = uniq
		}
	} else {
		uncached = uniq
	}

	if len(uncached) > 0 {
		var list []*model.CategoryM
		if err := s.ds.DB(ctx, where.F("id", uncached)).Find(&list).Error; err != nil {
			return nil, err
		}
		for _, c := range list {
			result[c.ID] = c
		}
		if rdb != nil && len(list) > 0 {
			pipe := rdb.Pipeline()
			for _, c := range list {
				key := fmt.Sprintf("miniblog:category:%d", c.ID)
				if data, err := json.Marshal(c); err == nil {
					pipe.Set(ctx, key, data, 12*time.Hour)
				}
			}
			_, _ = pipe.Exec(ctx)
		}
	}

	return result, nil
}

// ListAllWithCache 返回全量分类列表并带缓存
func (s *categoryStore) ListAllWithCache(ctx context.Context) ([]*model.CategoryM, error) {
	rdb := s.ds.Redis(ctx)
	if rdb != nil {
		if bs, err := rdb.Get(ctx, cacheKeyCategoryListAll).Bytes(); err == nil && len(bs) > 0 {
			var list []*model.CategoryM
			if jsonErr := json.Unmarshal(bs, &list); jsonErr == nil {
				return list, nil
			}
		}
	}

	v, err, _ := listGroup.Do(cacheKeyCategoryListAll, func() (any, error) {
		whr := where.NewWhere()
		_, categoryList, err := s.Store.List(ctx, whr)
		if err != nil {
			return nil, err
		}
		if rdb != nil {
			if data, mErr := json.Marshal(categoryList); mErr == nil {
				_ = rdb.Set(ctx, cacheKeyCategoryListAll, data, cacheTTLCategoryList).Err()
			}
		}
		return categoryList, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]*model.CategoryM), nil
}

// ListActiveWithCache 返回 is_active=1 的分类列表并带缓存
func (s *categoryStore) ListActiveWithCache(ctx context.Context) ([]*model.CategoryM, error) {
	rdb := s.ds.Redis(ctx)
	if rdb != nil {
		if bs, err := rdb.Get(ctx, cacheKeyCategoryListActive).Bytes(); err == nil && len(bs) > 0 {
			var list []*model.CategoryM
			if jsonErr := json.Unmarshal(bs, &list); jsonErr == nil {
				return list, nil
			}
		}
	}

	v, err, _ := listGroup.Do(cacheKeyCategoryListActive, func() (any, error) {
		whr := where.F("is_active", 1)
		_, categoryList, err := s.Store.List(ctx, whr)
		if err != nil {
			return nil, err
		}
		if rdb != nil {
			if data, mErr := json.Marshal(categoryList); mErr == nil {
				_ = rdb.Set(ctx, cacheKeyCategoryListActive, data, cacheTTLCategoryList).Err()
			}
		}
		return categoryList, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]*model.CategoryM), nil
}
