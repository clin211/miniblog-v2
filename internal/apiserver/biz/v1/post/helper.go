// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package post

import (
	"context"
	"sync"

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

// 删除 Redis 具体实现，交由 store 层统一缓存

// 删除具体缓存读写逻辑

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
}

// NewHighPerformanceCategoryLoaderWithCache 创建带缓存的高性能分类加载器
func NewHighPerformanceCategoryLoader(store store.IStore, categoriesMap map[int32]*model.CategoryM, mu *sync.RWMutex) *HighPerformanceCategoryLoader {
	return &HighPerformanceCategoryLoader{store: store, categoriesMap: categoriesMap, mu: mu}
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
		// 使用 store 层批量缓存 API
		cached, err := hcl.store.Category().BatchGetByIDsWithCache(ctx, categoryIDs)
		if err != nil {
			return err
		}
		hcl.mu.Lock()
		for id, c := range cached {
			hcl.categoriesMap[id] = c
		}
		hcl.mu.Unlock()
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
}

// NewHighPerformanceTagLoaderWithCache 创建带缓存的高性能标签加载器
func NewHighPerformanceTagLoader(store store.IStore, postTagsMap map[string][]*model.TagM, mu *sync.RWMutex) *HighPerformanceTagLoader {
	return &HighPerformanceTagLoader{store: store, postTagsMap: postTagsMap, mu: mu}
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

	// 由 store 层提供批量缓存 API
	cachedTags, err := htl.store.Tag().BatchGetByIDsWithCache(ctx, allTagIDs)
	if err != nil {
		return err
	}
	// 建立标签ID到标签对象的映射
	tagsMap := make(map[int32]*model.TagM, len(cachedTags))
	for id, tag := range cachedTags {
		tagsMap[id] = tag
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

	// 创建加载器（缓存由 store 层负责）
	categoryLoader := NewHighPerformanceCategoryLoader(store, categoriesMap, &mu)
	tagLoader := NewHighPerformanceTagLoader(store, postTagsMap, &mu)

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

// LoadSinglePostWithRelations 加载单篇文章及其关联数据。
// 相比批量装载路径，单篇装载会避免不必要的并发与对象池开销，减少额外分配与同步成本。
func LoadSinglePostWithRelations(ctx context.Context, store store.IStore, post *model.PostM) (*v1.Post, error) {
	if post == nil {
		return nil, nil
	}

	var category *model.CategoryM
	var tags []*model.TagM

	// 分类：优先使用缓存，其次数据库。
	if post.CategoryID != nil {
		categoryID := *post.CategoryID
		cachedMap, err := store.Category().BatchGetByIDsWithCache(ctx, []int32{categoryID})
		if err != nil {
			return nil, err
		}
		category = cachedMap[categoryID]
	}

	// 标签列表：先查 post_tag，再批量获取标签详情（优先缓存）。
	postTagWhr := where.F("post_id", post.PostID)
	_, postTags, err := store.PostTag().List(ctx, postTagWhr)
	if err != nil {
		return nil, err
	}

	if len(postTags) > 0 {
		tagIDs := make([]int32, 0, len(postTags))
		for _, pt := range postTags {
			tagIDs = append(tagIDs, pt.TagID)
		}
		cachedTags, err := store.Tag().BatchGetByIDsWithCache(ctx, tagIDs)
		if err != nil {
			return nil, err
		}
		tags = make([]*model.TagM, 0, len(tagIDs))
		for _, id := range tagIDs {
			if t, ok := cachedTags[id]; ok {
				tags = append(tags, t)
			}
		}
	}

	protoPost := conversion.PostModelToPostV1WithRelations(post, category, tags)
	return protoPost, nil
}
