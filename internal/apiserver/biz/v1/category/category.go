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
	return &v1.CreateCategoryResponse{Id: categoryM.ID}, nil
}

// Update 实现 CategoryBiz 接口中的 Update 方法.
func (b *categoryBiz) Update(ctx context.Context, rq *v1.UpdateCategoryRequest) (*v1.UpdateCategoryResponse, error) {
	whr := where.F("id", rq.GetId())
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
	whr := where.F("id", rq.GetId())
	if err := b.store.Category().Delete(ctx, whr); err != nil {
		return nil, err
	}
	return &v1.DeleteCategoryResponse{}, nil
}

// Get 实现 CategoryBiz 接口中的 Get 方法.
func (b *categoryBiz) Get(ctx context.Context, rq *v1.GetCategoryRequest) (*v1.GetCategoryResponse, error) {
	whr := where.F("id", rq.GetId())
	categoryM, err := b.store.Category().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &v1.GetCategoryResponse{Category: conversion.CategoryModelToCategoryV1(categoryM)}, nil
}

// List 实现 CategoryBiz 接口中的 List 方法.
func (b *categoryBiz) List(ctx context.Context, rq *v1.ListCategoryRequest) (*v1.ListCategoryResponse, error) {
	whr := where.NewWhere()

	// 获取所有分类数据
	_, categoryList, err := b.store.Category().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	// 构建父子关系映射
	childrenMap := make(map[int32][]*model.CategoryM)

	// 一次遍历构建父子关系映射和根分类切片
	rootCategoriesMap := make(map[int32]*model.CategoryM)
	for _, category := range categoryList {
		if category.ParentID == nil || *category.ParentID == 0 {
			rootCategoriesMap[category.ID] = category
		} else {
			parentID := *category.ParentID
			childrenMap[parentID] = append(childrenMap[parentID], category)
		}
	}

	// 使用 sort.Slice 替代冒泡排序，时间复杂度从 O(n²) 降低到 O(n log n)
	sortFunc := func(categories []*model.CategoryM) {
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
			Name:        rootCategory.Name,
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
			categoryWithChildren.IsActive = *rootCategory.IsActive
		}
		if rootCategory.CreatedAt != nil {
			categoryWithChildren.CreatedAt = rootCategory.CreatedAt.Unix()
		}
		if rootCategory.UpdatedAt != nil {
			categoryWithChildren.UpdatedAt = rootCategory.UpdatedAt.Unix()
		}

		responseCategories = append(responseCategories, categoryWithChildren)
	}

	return &v1.ListCategoryResponse{
		Total:      int32(len(categoryList)),
		Categories: responseCategories,
	}, nil
}
