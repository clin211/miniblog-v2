package http

import (
	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

// CreatePost 创建博客帖子.
func (h *Handler) CreatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Create)
}

// UpdatePost 更新博客帖子.
func (h *Handler) UpdatePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Update)
}

// DeletePost 删除博客帖子.
func (h *Handler) DeletePost(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.PostV1().Delete)
}

// GetPost 获取博客帖子.
func (h *Handler) GetPost(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.PostV1().Get)
}

// ListPosts 列出用户的所有博客帖子.
func (h *Handler) ListPost(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.PostV1().List)
}
