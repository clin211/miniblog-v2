// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package post

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	"github.com/clin211/miniblog-v2/internal/apiserver/pkg/conversion"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	"github.com/clin211/miniblog-v2/internal/pkg/known"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/clin211/miniblog-v2/pkg/where"
)

var (
	// 对象池：复用 map 对象减少内存分配
	categoryMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[int32]*model.CategoryM, 16)
		},
	}

	postTagsMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[string][]*model.TagM, 64)
		},
	}

	categoryIDsPool = sync.Pool{
		New: func() interface{} {
			return make([]int32, 0, 16)
		},
	}

	postIDsPool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 64)
		},
	}

	tagIDsPool = sync.Pool{
		New: func() interface{} {
			return make([]int32, 0, 32)
		},
	}
)

// Redis 缓存管理器
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache 创建 Redis 缓存管理器
func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

// 缓存键前缀
const (
	categoryKeyPrefix = "miniblog:category:"
	tagKeyPrefix      = "miniblog:tag:"
)

// getBatchCategories 批量获取缓存的分类
func (rc *RedisCache) getBatchCategories(ctx context.Context, ids []int32) (map[int32]*model.CategoryM, []int32) {
	if len(ids) == 0 {
		return make(map[int32]*model.CategoryM), nil
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = fmt.Sprintf("%s%d", categoryKeyPrefix, id)
	}

	vals, err := rc.client.MGet(ctx, keys...).Result()
	if err != nil {
		log.W(ctx).Errorw("Failed to batch get cached categories", "error", err)
		return make(map[int32]*model.CategoryM), ids
	}

	cached := make(map[int32]*model.CategoryM)
	uncached := make([]int32, 0)

	for i, val := range vals {
		if val == nil {
			uncached = append(uncached, ids[i])
			continue
		}

		var category model.CategoryM
		if err := json.Unmarshal([]byte(val.(string)), &category); err != nil {
			log.W(ctx).Errorw("Failed to unmarshal cached category", "error", err, "categoryID", ids[i])
			uncached = append(uncached, ids[i])
			continue
		}

		cached[ids[i]] = &category
	}

	return cached, uncached
}

// setBatchCategories 批量设置缓存的分类
func (rc *RedisCache) setBatchCategories(ctx context.Context, categories []*model.CategoryM) {
	if len(categories) == 0 {
		return
	}

	pipe := rc.client.Pipeline()
	for _, category := range categories {
		key := fmt.Sprintf("%s%d", categoryKeyPrefix, category.ID)
		data, err := json.Marshal(category)
		if err != nil {
			log.W(ctx).Errorw("Failed to marshal category for cache", "error", err, "categoryID", category.ID)
			continue
		}
		pipe.Set(ctx, key, data, rc.ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.W(ctx).Errorw("Failed to batch set cached categories", "error", err)
	}
}

// getBatchTags 批量获取缓存的标签
func (rc *RedisCache) getBatchTags(ctx context.Context, ids []int32) (map[int32]*model.TagM, []int32) {
	if len(ids) == 0 {
		return make(map[int32]*model.TagM), nil
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = fmt.Sprintf("%s%d", tagKeyPrefix, id)
	}

	vals, err := rc.client.MGet(ctx, keys...).Result()
	if err != nil {
		log.W(ctx).Errorw("Failed to batch get cached tags", "error", err)
		return make(map[int32]*model.TagM), ids
	}

	cached := make(map[int32]*model.TagM)
	uncached := make([]int32, 0)

	for i, val := range vals {
		if val == nil {
			uncached = append(uncached, ids[i])
			continue
		}

		var tag model.TagM
		if err := json.Unmarshal([]byte(val.(string)), &tag); err != nil {
			log.W(ctx).Errorw("Failed to unmarshal cached tag", "error", err, "tagID", ids[i])
			uncached = append(uncached, ids[i])
			continue
		}

		cached[ids[i]] = &tag
	}

	return cached, uncached
}

// setBatchTags 批量设置缓存的标签
func (rc *RedisCache) setBatchTags(ctx context.Context, tags []*model.TagM) {
	if len(tags) == 0 {
		return
	}

	pipe := rc.client.Pipeline()
	for _, tag := range tags {
		key := fmt.Sprintf("%s%d", tagKeyPrefix, tag.ID)
		data, err := json.Marshal(tag)
		if err != nil {
			log.W(ctx).Errorw("Failed to marshal tag for cache", "error", err, "tagID", tag.ID)
			continue
		}
		pipe.Set(ctx, key, data, rc.ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.W(ctx).Errorw("Failed to batch set cached tags", "error", err)
	}
}

// resetPool 重置切片池对象
func resetSlice[T any](s []T) []T {
	return s[:0]
}

// resetMap 重置 map 池对象
func resetCategoryMap(m map[int32]*model.CategoryM) {
	for k := range m {
		delete(m, k)
	}
}

func resetPostTagsMap(m map[string][]*model.TagM) {
	for k := range m {
		delete(m, k)
	}
}

// ================================
// 关联数据加载器接口和实现
// ================================

// RelationLoader 定义关联数据加载器接口
type RelationLoader interface {
	// Load 加载关联数据，返回加载函数
	Load(ctx context.Context, posts []*model.PostM) func() error
}

// ================================
// 分类数据加载器
// ================================

// HighPerformanceCategoryLoader 高性能分类数据加载器
type HighPerformanceCategoryLoader struct {
	store         store.IStore
	categoriesMap map[int32]*model.CategoryM
	mu            *sync.RWMutex
	cache         *RedisCache
}

// NewHighPerformanceCategoryLoaderWithCache 创建带缓存的高性能分类加载器
func NewHighPerformanceCategoryLoaderWithCache(store store.IStore, categoriesMap map[int32]*model.CategoryM, mu *sync.RWMutex, cache *RedisCache) *HighPerformanceCategoryLoader {
	return &HighPerformanceCategoryLoader{
		store:         store,
		categoriesMap: categoriesMap,
		mu:            mu,
		cache:         cache,
	}
}

// Load 实现 RelationLoader 接口 - 高性能版本（使用 Redis 缓存）
func (hcl *HighPerformanceCategoryLoader) Load(ctx context.Context, posts []*model.PostM) func() error {
	// 从对象池获取切片，减少内存分配
	categoryIDs := categoryIDsPool.Get().([]int32)
	categoryIDs = resetSlice(categoryIDs)

	defer func() {
		categoryIDsPool.Put(categoryIDs)
	}()

	// 收集需要查询的分类ID
	categoryIDSet := make(map[int32]bool, len(posts)/2) // 假设平均每2篇文章一个分类
	for _, post := range posts {
		if post.CategoryID != nil && !categoryIDSet[*post.CategoryID] {
			categoryIDs = append(categoryIDs, *post.CategoryID)
			categoryIDSet[*post.CategoryID] = true
		}
	}

	if len(categoryIDs) == 0 {
		return func() error { return nil }
	}

	return func() error {
		var cachedCategories map[int32]*model.CategoryM
		var uncachedIDs []int32

		// 如果有 Redis 缓存，使用批量缓存查询
		if hcl.cache != nil {
			cachedCategories, uncachedIDs = hcl.cache.getBatchCategories(ctx, categoryIDs)
		} else {
			// 没有缓存时，直接查询所有ID
			cachedCategories = make(map[int32]*model.CategoryM)
			uncachedIDs = categoryIDs
		}

		// 先设置缓存中的数据
		if len(cachedCategories) > 0 {
			hcl.mu.Lock()
			for id, category := range cachedCategories {
				hcl.categoriesMap[id] = category
			}
			hcl.mu.Unlock()
		}

		// 查询未缓存的分类
		if len(uncachedIDs) > 0 {
			whr := where.F("id", uncachedIDs)
			_, categories, err := hcl.store.Category().List(ctx, whr)
			if err != nil {
				log.W(ctx).Errorw("Failed to load categories", "error", err, "categoryIDs", uncachedIDs)
				return err
			}

			hcl.mu.Lock()
			for _, category := range categories {
				hcl.categoriesMap[category.ID] = category
			}
			hcl.mu.Unlock()

			// 批量更新 Redis 缓存
			if hcl.cache != nil {
				hcl.cache.setBatchCategories(ctx, categories)
			}
		}

		return nil
	}
}

// ================================
// 标签数据加载器
// ================================

// HighPerformanceTagLoader 高性能标签数据加载器
type HighPerformanceTagLoader struct {
	store       store.IStore
	postTagsMap map[string][]*model.TagM
	mu          *sync.RWMutex
	cache       *RedisCache
}

// NewHighPerformanceTagLoaderWithCache 创建带缓存的高性能标签加载器
func NewHighPerformanceTagLoaderWithCache(store store.IStore, postTagsMap map[string][]*model.TagM, mu *sync.RWMutex, cache *RedisCache) *HighPerformanceTagLoader {
	return &HighPerformanceTagLoader{
		store:       store,
		postTagsMap: postTagsMap,
		mu:          mu,
		cache:       cache,
	}
}

// Load 实现 RelationLoader 接口 - 高性能版本
func (htl *HighPerformanceTagLoader) Load(ctx context.Context, posts []*model.PostM) func() error {
	// 从对象池获取切片，减少内存分配
	postIDs := postIDsPool.Get().([]string)
	postIDs = resetSlice(postIDs)

	defer func() {
		postIDsPool.Put(postIDs)
	}()

	// 预分配容量
	if cap(postIDs) < len(posts) {
		postIDs = make([]string, 0, len(posts))
	}

	for _, post := range posts {
		postIDs = append(postIDs, post.PostID)
	}

	if len(postIDs) == 0 {
		return func() error { return nil }
	}

	return func() error {
		return htl.loadTagRelationsOptimized(ctx, postIDs)
	}
}

// loadTagRelationsOptimized 高性能标签关联数据加载实现（使用 Redis 缓存）
func (htl *HighPerformanceTagLoader) loadTagRelationsOptimized(ctx context.Context, postIDs []string) error {
	// 查询 post_tag 关联 - 使用预编译语句提升性能
	postTagWhr := where.F("post_id", postIDs)
	var allPostTags []*model.PostTagM
	db := htl.store.DB(ctx, postTagWhr)

	// 预分配切片容量，假设平均每篇文章3个标签
	allPostTags = make([]*model.PostTagM, 0, len(postIDs)*3)
	err := db.Find(&allPostTags).Error
	if err != nil {
		log.W(ctx).Errorw("Failed to query post_tag relations", "error", err, "postIDs", postIDs)
		return err
	}

	if len(allPostTags) == 0 {
		return nil
	}

	// 从对象池获取标签ID切片
	allTagIDs := tagIDsPool.Get().([]int32)
	allTagIDs = resetSlice(allTagIDs)
	defer func() {
		tagIDsPool.Put(allTagIDs)
	}()

	// 按文章分组并收集所有标签ID - 预分配容量优化
	postTagsGrouped := make(map[string][]int32, len(postIDs))
	tagIDSet := make(map[int32]bool, len(allPostTags))

	for _, pt := range allPostTags {
		postTagsGrouped[pt.PostID] = append(postTagsGrouped[pt.PostID], pt.TagID)

		if !tagIDSet[pt.TagID] {
			allTagIDs = append(allTagIDs, pt.TagID)
			tagIDSet[pt.TagID] = true
		}
	}

	if len(allTagIDs) == 0 {
		return nil
	}

	var cachedTags map[int32]*model.TagM
	var uncachedTagIDs []int32

	// 如果有 Redis 缓存，使用批量缓存查询
	if htl.cache != nil {
		cachedTags, uncachedTagIDs = htl.cache.getBatchTags(ctx, allTagIDs)
	} else {
		// 没有缓存时，直接查询所有ID
		cachedTags = make(map[int32]*model.TagM)
		uncachedTagIDs = allTagIDs
	}

	// 建立标签ID到标签对象的映射
	tagsMap := make(map[int32]*model.TagM, len(cachedTags)+len(uncachedTagIDs))

	// 先添加缓存的标签
	for id, tag := range cachedTags {
		tagsMap[id] = tag
	}

	// 批量查询未缓存的标签详情
	if len(uncachedTagIDs) > 0 {
		tagWhr := where.F("id", uncachedTagIDs)
		_, allTags, err := htl.store.Tag().List(ctx, tagWhr)
		if err != nil {
			log.W(ctx).Errorw("Failed to query tags", "error", err, "tagIDs", uncachedTagIDs)
			return err
		}

		for _, tag := range allTags {
			tagsMap[tag.ID] = tag
		}

		// 批量更新 Redis 缓存
		if htl.cache != nil {
			htl.cache.setBatchTags(ctx, allTags)
		}
	}

	// 为每篇文章组装标签列表 - 减少锁的持有时间
	postTagsResult := make(map[string][]*model.TagM, len(postTagsGrouped))
	for postID, tagIDs := range postTagsGrouped {
		tags := make([]*model.TagM, 0, len(tagIDs))
		for _, tagID := range tagIDs {
			if tag, exists := tagsMap[tagID]; exists {
				tags = append(tags, tag)
			}
		}
		postTagsResult[postID] = tags
	}

	// 一次性更新结果映射
	htl.mu.Lock()
	for postID, tags := range postTagsResult {
		htl.postTagsMap[postID] = tags
	}
	htl.mu.Unlock()

	return nil
}

// ================================
// 关联数据加载协调器
// ================================

// RelationLoadCoordinator 关联数据加载协调器
type RelationLoadCoordinator struct {
	loaders []RelationLoader
}

// NewRelationLoadCoordinator 创建关联数据加载协调器
func NewRelationLoadCoordinator(loaders ...RelationLoader) *RelationLoadCoordinator {
	return &RelationLoadCoordinator{
		loaders: loaders,
	}
}

// LoadConcurrently 并发加载所有关联数据
func (rlc *RelationLoadCoordinator) LoadConcurrently(ctx context.Context, posts []*model.PostM) error {
	if len(posts) == 0 {
		return nil
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 为每个加载器创建并发任务
	for _, loader := range rlc.loaders {
		loadFunc := loader.Load(egCtx, posts)
		eg.Go(loadFunc)
	}

	if err := eg.Wait(); err != nil {
		log.W(ctx).Errorw("Failed to load post relations", "error", err)
		return err
	}

	return nil
}

// ================================
// 文章关联数据建造者
// ================================

// PostWithRelationsBuilder 文章关联数据建造者
type PostWithRelationsBuilder struct {
	categoriesMap map[int32]*model.CategoryM
	postTagsMap   map[string][]*model.TagM
	mu            *sync.RWMutex
}

// NewPostWithRelationsBuilder 创建文章关联数据建造者
func NewPostWithRelationsBuilder(categoriesMap map[int32]*model.CategoryM, postTagsMap map[string][]*model.TagM, mu *sync.RWMutex) *PostWithRelationsBuilder {
	return &PostWithRelationsBuilder{
		categoriesMap: categoriesMap,
		postTagsMap:   postTagsMap,
		mu:            mu,
	}
}

// BuildPosts 构建带关联数据的文章列表
func (builder *PostWithRelationsBuilder) BuildPosts(posts []*model.PostM) []*v1.Post {
	results := make([]*v1.Post, len(posts)) // 精确分配容量

	// 批量获取锁，减少锁竞争
	builder.mu.RLock()
	defer builder.mu.RUnlock()

	for i, post := range posts {
		var category *model.CategoryM
		if post.CategoryID != nil {
			category = builder.categoriesMap[*post.CategoryID]
		}

		tags := builder.postTagsMap[post.PostID]
		protoPost := conversion.PostModelToPostV1WithRelations(post, category, tags)
		results[i] = protoPost
	}

	return results
}

// BuildPostsConcurrently 并发构建带关联数据的文章列表 - 适用于大量数据
func (builder *PostWithRelationsBuilder) BuildPostsConcurrently(posts []*model.PostM) []*v1.Post {
	results := make([]*v1.Post, len(posts))

	// 对于小数据量，直接使用串行处理
	if len(posts) < 50 {
		return builder.BuildPosts(posts)
	}

	// 大数据量使用并发处理
	const batchSize = 50
	var wg sync.WaitGroup

	for i := 0; i < len(posts); i += batchSize {
		end := i + batchSize
		if end > len(posts) {
			end = len(posts)
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			builder.mu.RLock()
			defer builder.mu.RUnlock()

			for j := start; j < end; j++ {
				post := posts[j]
				var category *model.CategoryM
				if post.CategoryID != nil {
					category = builder.categoriesMap[*post.CategoryID]
				}

				tags := builder.postTagsMap[post.PostID]
				protoPost := conversion.PostModelToPostV1WithRelations(post, category, tags)
				results[j] = protoPost
			}
		}(i, end)
	}

	wg.Wait()
	return results
}

// ================================
// 便捷工厂方法
// ================================

// LoadPostsWithRelations 高性能版本：批量加载文章及其关联数据
// 包含 Redis 缓存、对象池、并发优化和性能监控
func LoadPostsWithRelations(ctx context.Context, store store.IStore, posts []*model.PostM) ([]*v1.Post, error) {
	if len(posts) == 0 {
		return []*v1.Post{}, nil
	}

	// 性能监控
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		log.W(ctx).Infow("LoadPostsWithRelations performance",
			"posts_count", len(posts),
			"duration_ms", duration.Milliseconds(),
			"avg_ms_per_post", float64(duration.Milliseconds())/float64(len(posts)))
	}()

	// 从对象池获取 map 对象，减少内存分配
	categoriesMap := categoryMapPool.Get().(map[int32]*model.CategoryM)
	postTagsMap := postTagsMapPool.Get().(map[string][]*model.TagM)

	// 确保使用后归还到池中
	defer func() {
		resetCategoryMap(categoriesMap)
		resetPostTagsMap(postTagsMap)
		categoryMapPool.Put(categoriesMap)
		postTagsMapPool.Put(postTagsMap)
	}()

	var mu sync.RWMutex

	// 创建 Redis 缓存管理器
	cache := NewRedisCache(store.Redis(ctx), 12*time.Hour)

	// 创建带 Redis 缓存的高性能关联数据加载器
	categoryLoader := NewHighPerformanceCategoryLoaderWithCache(store, categoriesMap, &mu, cache)
	tagLoader := NewHighPerformanceTagLoaderWithCache(store, postTagsMap, &mu, cache)

	// 创建协调器并发加载所有关联数据
	coordinator := NewRelationLoadCoordinator(categoryLoader, tagLoader)
	if err := coordinator.LoadConcurrently(ctx, posts); err != nil {
		return nil, err
	}

	// 使用建造者模式构建最终结果
	builder := NewPostWithRelationsBuilder(categoriesMap, postTagsMap, &mu)

	// 根据数据量选择构建策略
	if len(posts) > 100 {
		return builder.BuildPostsConcurrently(posts), nil
	}
	return builder.BuildPosts(posts), nil
}
