// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package category

import (
	"context"
	"sort"
	"time"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	"github.com/clin211/miniblog-v2/internal/apiserver/pkg/conversion"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/clin211/miniblog-v2/pkg/copier"
	"github.com/clin211/miniblog-v2/pkg/where"
)

type CategoryBiz interface {
	Create(ctx context.Context, rq *v1.CreateCategoryRequest) (*v1.CreateCategoryResponse, error)
	Update(ctx context.Context, rq *v1.UpdateCategoryRequest) (*v1.UpdateCategoryResponse, error)
	Delete(ctx context.Context, rq *v1.DeleteCategoryRequest) (*v1.DeleteCategoryResponse, error)
	Get(ctx context.Context, rq *v1.GetCategoryRequest) (*v1.GetCategoryResponse, error)
	List(ctx context.Context, rq *v1.ListCategoryRequest) (*v1.ListCategoryResponse, error)
	AppList(ctx context.Context, rq *v1.ListCategoryRequest) (*v1.ListCategoryResponse, error)
}

type categoryBiz struct {
	store store.IStore
}

// 确保 categoryBiz 实现了 CategoryBiz 接口
var _ CategoryBiz = (*categoryBiz)(nil)

// 创建一个 CategoryBiz 的实例
func New(store store.IStore) *categoryBiz {
	return &categoryBiz{store: store}
}

// Create 实现 CategoryBiz 接口中的 Create 方法.
func (b *categoryBiz) Create(ctx context.Context, rq *v1.CreateCategoryRequest) (*v1.CreateCategoryResponse, error) {
	var categoryM model.CategoryM
	_ = copier.Copy(&categoryM, rq)

	// 手动设置创建时间
	now := time.Now()
	categoryM.CreatedAt = &now
	categoryM.UpdatedAt = &now

	if err := b.store.Category().Create(ctx, &categoryM); err != nil {
		return nil, err
	}
	return &v1.CreateCategoryResponse{CategoryID: categoryM.CategoryID}, nil
}

// Update 实现 CategoryBiz 接口中的 Update 方法.
func (b *categoryBiz) Update(ctx context.Context, rq *v1.UpdateCategoryRequest) (*v1.UpdateCategoryResponse, error) {
	whr := where.F("category_id", rq.GetCategoryID())
	categoryM, err := b.store.Category().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	_ = copier.Copy(&categoryM, rq)

	// 手动设置更新时间
	now := time.Now()
	categoryM.UpdatedAt = &now

	if err := b.store.Category().Update(ctx, categoryM); err != nil {
		return nil, err
	}
	return &v1.UpdateCategoryResponse{}, nil
}

// Delete 实现 CategoryBiz 接口中的 Delete 方法.
func (b *categoryBiz) Delete(ctx context.Context, rq *v1.DeleteCategoryRequest) (*v1.DeleteCategoryResponse, error) {
	whr := where.F("category_id", rq.GetCategoryID())
	if err := b.store.Category().Delete(ctx, whr); err != nil {
		return nil, err
	}
	return &v1.DeleteCategoryResponse{}, nil
}

// Get 实现 CategoryBiz 接口中的 Get 方法.
func (b *categoryBiz) Get(ctx context.Context, rq *v1.GetCategoryRequest) (*v1.GetCategoryResponse, error) {
	whr := where.F("category_id", rq.GetCategoryID())
	categoryM, err := b.store.Category().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &v1.GetCategoryResponse{Category: conversion.CategoryModelToCategoryV1(categoryM)}, nil
}

// List 实现 CategoryBiz 接口中的 List 方法.
func (b *categoryBiz) List(ctx context.Context, rq *v1.ListCategoryRequest) (*v1.ListCategoryResponse, error) {
	// 读取全量列表
	categoryList, err := b.store.Category().ListAllWithCache(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ListCategoryResponse{
		Total:      int32(len(categoryList)),
		Categories: buildHierarchicalCategories(categoryList),
	}, nil
}

// List 实现 CategoryBiz 接口中的 List 方法.
func (b *categoryBiz) AppList(ctx context.Context, rq *v1.ListCategoryRequest) (*v1.ListCategoryResponse, error) {
	// 读取 active 列表
	categoryList, err := b.store.Category().ListActiveWithCache(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ListCategoryResponse{
		Total:      int32(len(categoryList)),
		Categories: buildHierarchicalCategories(categoryList),
	}, nil
}

// sortFunc 对分类列表按 SortOrder 排序
func sortFunc(categories []*model.CategoryM) {
	sort.Slice(categories, func(i, j int) bool {
		order1 := int32(0)
		order2 := int32(0)
		if categories[i].SortOrder != nil {
			order1 = *categories[i].SortOrder
		}
		if categories[j].SortOrder != nil {
			order2 = *categories[j].SortOrder
		}
		return order1 < order2
	})
}

// buildHierarchicalCategories 构建分层的分类结构
func buildHierarchicalCategories(categoryList []*model.CategoryM) []*v1.ListCategoryResponse_Categories {
	// 构建父子关系映射
	childrenMap := make(map[int32][]*model.CategoryM)
	rootCategoriesMap := make(map[int32]*model.CategoryM)

	// 一次遍历构建父子关系映射和根分类切片
	for _, category := range categoryList {
		if category.ParentID == nil || *category.ParentID == 0 {
			rootCategoriesMap[category.ID] = category
		} else {
			parentID := *category.ParentID
			childrenMap[parentID] = append(childrenMap[parentID], category)
		}
	}

	// 将 map 转换为切片并排序
	rootCategories := make([]*model.CategoryM, 0, len(rootCategoriesMap))
	for _, category := range rootCategoriesMap {
		rootCategories = append(rootCategories, category)
	}
	sortFunc(rootCategories)

	// 预分配响应切片容量，减少内存重新分配
	responseCategories := make([]*v1.ListCategoryResponse_Categories, 0, len(rootCategories))

	// 构建响应结构
	for _, rootCategory := range rootCategories {
		children := childrenMap[rootCategory.ID]
		sortFunc(children)

		// 预分配子分类切片容量
		childrenV1 := make([]*v1.Category, 0, len(children))
		for _, child := range children {
			childrenV1 = append(childrenV1, conversion.CategoryModelToCategoryV1(child))
		}

		// 使用结构体字面量初始化，减少零值赋值
		categoryWithChildren := &v1.ListCategoryResponse_Categories{
			Id:          rootCategory.ID,
			CategoryID:  rootCategory.CategoryID,
			Name:        rootCategory.Name,
			Icon:        rootCategory.Icon,
			Theme:       rootCategory.Theme,
			Description: rootCategory.Description,
			Children:    childrenV1,
		}

		// 设置可选字段
		if rootCategory.ParentID != nil {
			categoryWithChildren.ParentID = *rootCategory.ParentID
		}
		if rootCategory.SortOrder != nil {
			categoryWithChildren.SortOrder = *rootCategory.SortOrder
		}
		if rootCategory.IsActive != nil {
			if *rootCategory.IsActive == 1 {
				categoryWithChildren.IsActive = v1.IsActive_IS_ACTIVE_ACTIVE
			} else {
				categoryWithChildren.IsActive = v1.IsActive_IS_ACTIVE_DISABLED
			}
		}
		if rootCategory.CreatedAt != nil {
			categoryWithChildren.CreatedAt = rootCategory.CreatedAt.Unix()
		}
		if rootCategory.UpdatedAt != nil {
			categoryWithChildren.UpdatedAt = rootCategory.UpdatedAt.Unix()
		}

		responseCategories = append(responseCategories, categoryWithChildren)
	}

	return responseCategories
}
