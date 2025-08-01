// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package post

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	"github.com/clin211/miniblog-v2/internal/apiserver/pkg/conversion"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/clin211/miniblog-v2/pkg/where"
)

// PostBiz 定义处理帖子请求所需的方法.
type PostBiz interface {
	Create(ctx context.Context, rq *v1.CreatePostRequest) (*v1.CreatePostResponse, error)
	Update(ctx context.Context, rq *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *v1.DeletePostRequest) (*v1.DeletePostResponse, error)
	Get(ctx context.Context, rq *v1.GetPostRequest) (*v1.GetPostResponse, error)
	List(ctx context.Context, rq *v1.ListPostRequest) (*v1.ListPostResponse, error)

	PostExpansion
}

// PostExpansion 定义额外的帖子操作方法.
type PostExpansion interface {
	AppList(ctx context.Context, rq *v1.ListPostRequest) (*v1.ListPostResponse, error)
}

// postBiz 是 PostBiz 接口的实现.
type postBiz struct {
	store store.IStore
}

// 确保 postBiz 实现了 PostBiz 接口.
var _ PostBiz = (*postBiz)(nil)

// New 创建 postBiz 的实例.
func New(store store.IStore) *postBiz {
	return &postBiz{store: store}
}

// Create 实现 PostBiz 接口中的 Create 方法.
func (b *postBiz) Create(ctx context.Context, rq *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, rq)
	postM.UserID = contextx.UserID(ctx)

	// 手动设置创建时间
	now := time.Now()
	postM.CreatedAt = &now
	postM.UpdatedAt = &now

	// 使用事务确保创建文章和标签关联的原子性
	err := b.store.TX(ctx, func(txCtx context.Context) error {
		// 创建文章
		if err := b.store.Post().Create(txCtx, &postM); err != nil {
			log.W(ctx).Errorw("create post failed", "error", err)
			return err
		}

		// 创建文章标签关联
		for _, tagID := range rq.GetTags() {
			postTagM := model.PostTagM{
				PostID:    postM.PostID,
				TagID:     tagID,
				CreatedAt: &now,
				UpdatedAt: &now,
			}
			if err := b.store.PostTag().Create(txCtx, &postTagM); err != nil {
				log.W(ctx).Errorw("create post tag failed", "error", err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.W(ctx).Errorw("create post failed", "error", err)
		return nil, err
	}

	return &v1.CreatePostResponse{PostID: postM.PostID}, nil
}

// Update 实现 PostBiz 接口中的 Update 方法.
func (b *postBiz) Update(ctx context.Context, rq *v1.UpdatePostRequest) (*v1.UpdatePostResponse, error) {
	whr := where.T(ctx).F("post_id", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	// 使用事务确保更新文章和标签关联的原子性
	err = b.store.TX(ctx, func(txCtx context.Context) error {
		// 更新文章基本信息
		if rq.Title != nil {
			postM.Title = rq.GetTitle()
		}

		if rq.Content != nil {
			content := rq.GetContent()
			postM.Content = &content
		}

		if rq.Cover != nil {
			cover := rq.GetCover()
			postM.Cover = &cover
		}

		if rq.Summary != nil {
			summary := rq.GetSummary()
			postM.Summary = &summary
		}

		if rq.CategoryID != nil {
			categoryID := rq.GetCategoryID()
			postM.CategoryID = &categoryID
		}

		if rq.PostType != nil {
			postType := int32(rq.GetPostType())
			postM.PostType = &postType
		}

		if rq.OriginalAuthor != nil {
			originalAuthor := rq.GetOriginalAuthor()
			postM.OriginalAuthor = &originalAuthor
		}

		if rq.OriginalSource != nil {
			originalSource := rq.GetOriginalSource()
			postM.OriginalSource = &originalSource
		}

		if rq.OriginalAuthorIntro != nil {
			originalAuthorIntro := rq.GetOriginalAuthorIntro()
			postM.OriginalAuthorIntro = &originalAuthorIntro
		}

		if rq.Position != nil {
			position := rq.GetPosition()
			postM.Position = &position
		}

		if rq.Status != nil {
			status := int32(rq.GetStatus())
			postM.Status = &status
		}

		// 手动设置更新时间
		now := time.Now()
		postM.UpdatedAt = &now

		// 更新文章信息
		if err := b.store.Post().Update(txCtx, postM); err != nil {
			return err
		}

		// 如果提供了标签，则更新标签关联
		if len(rq.GetTags()) > 0 {
			// 删除现有的标签关联
			postTagWhr := where.T(txCtx).F("post_id", rq.GetPostID())
			if err := b.store.PostTag().Delete(txCtx, postTagWhr); err != nil {
				return err
			}

			// 创建新的标签关联
			for _, tagID := range rq.GetTags() {
				postTagM := model.PostTagM{
					PostID: postM.PostID,
					TagID:  tagID,
				}
				if err := b.store.PostTag().Create(txCtx, &postTagM); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &v1.UpdatePostResponse{}, nil
}

// Delete 实现 PostBiz 接口中的 Delete 方法.
func (b *postBiz) Delete(ctx context.Context, rq *v1.DeletePostRequest) (*v1.DeletePostResponse, error) {
	whr := where.T(ctx).F("post_id", rq.GetPostIDs())
	if err := b.store.Post().Delete(ctx, whr); err != nil {
		return nil, err
	}

	return &v1.DeletePostResponse{}, nil
}

// Get 实现 PostBiz 接口中的 Get 方法.
func (b *postBiz) Get(ctx context.Context, rq *v1.GetPostRequest) (*v1.GetPostResponse, error) {
	whr := where.T(ctx).F("post_id", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &v1.GetPostResponse{Post: conversion.PostModelToPostV1(postM)}, nil
}

// List 实现 PostBiz 接口中的 List 方法.
func (b *postBiz) List(ctx context.Context, rq *v1.ListPostRequest) (*v1.ListPostResponse, error) {
	whr := where.T(ctx).P(int(rq.GetOffset()), int(rq.GetLimit()))
	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*v1.Post, 0, len(postList))
	for _, post := range postList {
		converted := conversion.PostModelToPostV1(post)
		posts = append(posts, converted)
	}

	return &v1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}

func (b *postBiz) AppList(ctx context.Context, rq *v1.ListPostRequest) (*v1.ListPostResponse, error) {
	var whr *where.Options
	// 检查 categoryID 是否提供
	if rq.GetCategoryID() != "" {
		whr = where.F("category_id", *rq.CategoryID).P(int(rq.GetOffset()), int(rq.GetLimit()))
	} else {
		whr = where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	}

	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*v1.Post, 0, len(postList))
	for _, post := range postList {
		converted := conversion.PostModelToPostV1(post)
		posts = append(posts, converted)
	}

	return &v1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
