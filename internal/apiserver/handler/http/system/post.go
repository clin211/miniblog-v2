// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package app

import (
	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

// CreatePost 创建博客帖子.
func (h *Handler) CreatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Create, h.val.ValidateCreatePostRequest)
}

// UpdatePost 更新博客帖子.
func (h *Handler) UpdatePost(c *gin.Context) {
	core.HandleJSONWithURIRequest(c, h.biz.PostV1().Update, h.val.ValidateUpdatePostRequest)
}

// DeletePost 删除博客帖子.
func (h *Handler) DeletePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Delete, h.val.ValidateDeletePostRequest)
}

// GetPost 获取博客帖子.
func (h *Handler) GetPost(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.PostV1().Get, h.val.ValidateGetPostRequest)
}

// ListPosts 列出用户的所有博客帖子.
func (h *Handler) ListPost(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.PostV1().List, h.val.ValidateListPostRequest)
}
