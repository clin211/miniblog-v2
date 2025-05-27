package post

import (
	"context"

	"github.com/jinzhu/copier"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	"github.com/clin211/miniblog-v2/internal/apiserver/pkg/conversion"
	"github.com/clin211/miniblog-v2/internal/apiserver/store"
	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	apiv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/clin211/miniblog-v2/pkg/where"
)

// PostBiz 定义处理帖子请求所需的方法.
type PostBiz interface {
	Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error)
	Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error)
	List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error)

	PostExpansion
}

// PostExpansion 定义额外的帖子操作方法.
type PostExpansion interface{}

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
func (b *postBiz) Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, rq)
	postM.UserID = contextx.UserID(ctx)

	if err := b.store.Post().Create(ctx, &postM); err != nil {
		return nil, err
	}

	return &apiv1.CreatePostResponse{PostID: postM.PostID}, nil
}

// Update 实现 PostBiz 接口中的 Update 方法.
func (b *postBiz) Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	if rq.Title != nil {
		postM.Title = rq.GetTitle()
	}

	if rq.Content != nil {
		postM.Content = rq.GetContent()
	}

	if err := b.store.Post().Update(ctx, postM); err != nil {
		return nil, err
	}

	return &apiv1.UpdatePostResponse{}, nil
}

// Delete 实现 PostBiz 接口中的 Delete 方法.
func (b *postBiz) Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostIDs())
	if err := b.store.Post().Delete(ctx, whr); err != nil {
		return nil, err
	}

	return &apiv1.DeletePostResponse{}, nil
}

// Get 实现 PostBiz 接口中的 Get 方法.
func (b *postBiz) Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &apiv1.GetPostResponse{Post: conversion.PostModelToPostV1(postM)}, nil
}

// List 实现 PostBiz 接口中的 List 方法.
func (b *postBiz) List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	whr := where.T(ctx).P(int(rq.GetOffset()), int(rq.GetLimit()))
	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*apiv1.Post, 0, len(postList))
	for _, post := range postList {
		converted := conversion.PostModelToPostV1(post)
		posts = append(posts, converted)
	}

	return &apiv1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
