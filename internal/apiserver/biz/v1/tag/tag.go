// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package tag

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	"github.com/clin211/miniblog-v2/internal/apiserver/pkg/conversion"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	appv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1/app"
	"github.com/clin211/miniblog-v2/pkg/where"
)

// TagBiz 定义处理标签请求所需的方法.
type TagBiz interface {
	Create(ctx context.Context, rq *appv1.CreateTagRequest) (*appv1.CreateTagResponse, error)
	Update(ctx context.Context, rq *appv1.UpdateTagRequest) (*appv1.UpdateTagResponse, error)
	Delete(ctx context.Context, rq *appv1.DeleteTagRequest) (*appv1.DeleteTagResponse, error)
	Get(ctx context.Context, rq *appv1.GetTagRequest) (*appv1.GetTagResponse, error)
	List(ctx context.Context, rq *appv1.ListTagRequest) (*appv1.ListTagResponse, error)

	TagExpansion
}

// TagExpansion 定义额外的标签操作方法.
type TagExpansion interface{}

// tagBiz 是 TagBiz 接口的实现.
type tagBiz struct {
	store store.IStore
}

// 确保 tagBiz 实现了 TagBiz 接口.
var _ TagBiz = (*tagBiz)(nil)

// New 创建 tagBiz 的实例.
func New(store store.IStore) *tagBiz {
	return &tagBiz{store: store}
}

// Create 实现 TagBiz 接口中的 Create 方法.
func (b *tagBiz) Create(ctx context.Context, rq *appv1.CreateTagRequest) (*appv1.CreateTagResponse, error) {
	var tagM model.TagM
	_ = copier.Copy(&tagM, rq)

	// 手动设置创建时间
	now := time.Now()
	tagM.CreatedAt = &now
	tagM.UpdatedAt = &now

	if err := b.store.Tag().Create(ctx, &tagM); err != nil {
		return nil, err
	}

	return &appv1.CreateTagResponse{Id: tagM.ID}, nil
}

// Update 实现 TagBiz 接口中的 Update 方法.
func (b *tagBiz) Update(ctx context.Context, rq *appv1.UpdateTagRequest) (*appv1.UpdateTagResponse, error) {
	whr := where.F("id", rq.GetId())
	tagM, err := b.store.Tag().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	if rq.Name != nil {
		tagM.Name = rq.GetName()
	}

	if rq.Color != nil {
		color := rq.GetColor()
		tagM.Color = &color
	}

	// 手动设置更新时间
	now := time.Now()
	tagM.UpdatedAt = &now

	if err := b.store.Tag().Update(ctx, tagM); err != nil {
		return nil, err
	}

	return &appv1.UpdateTagResponse{}, nil
}

// Delete 实现 TagBiz 接口中的 Delete 方法.
func (b *tagBiz) Delete(ctx context.Context, rq *appv1.DeleteTagRequest) (*appv1.DeleteTagResponse, error) {
	whr := where.F("id", rq.GetId())
	if err := b.store.Tag().Delete(ctx, whr); err != nil {
		return nil, err
	}

	return &appv1.DeleteTagResponse{}, nil
}

// Get 实现 TagBiz 接口中的 Get 方法.
func (b *tagBiz) Get(ctx context.Context, rq *appv1.GetTagRequest) (*appv1.GetTagResponse, error) {
	whr := where.F("id", rq.GetId())
	tagM, err := b.store.Tag().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &appv1.GetTagResponse{Tag: conversion.TagModelToTagV1(tagM)}, nil
}

// List 实现 TagBiz 接口中的 List 方法.
func (b *tagBiz) List(ctx context.Context, rq *appv1.ListTagRequest) (*appv1.ListTagResponse, error) {
	whr := where.NewWhere()

	// 如果有名称过滤参数，添加到查询条件中
	if name := rq.GetName(); name != "" {
		whr = whr.F("name", name)
	}

	_, tagList, err := b.store.Tag().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	tags := make([]*appv1.Tag, 0, len(tagList))
	for _, tag := range tagList {
		converted := conversion.TagModelToTagV1(tag)
		tags = append(tags, converted)
	}

	return &appv1.ListTagResponse{Tags: tags}, nil
}
