// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package system

import (
	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

// CreateCategory 创建新分类
func (h *Handler) CreateCategory(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.CategoryV1().Create, h.val.ValidateCreateCategoryRequest)
}

// GetCategory 获取分类信息
func (h *Handler) GetCategory(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.CategoryV1().Get, h.val.ValidateGetCategoryRequest)
}

// UpdateCategory 更新分类信息
func (h *Handler) UpdateCategory(c *gin.Context) {
	core.HandleJSONWithURIRequest(c, h.biz.CategoryV1().Update, h.val.ValidateUpdateCategoryRequest)
}

// DeleteCategory 删除分类
func (h *Handler) DeleteCategory(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.CategoryV1().Delete, h.val.ValidateDeleteCategoryRequest)
}

// ListCategory 列出所有分类
func (h *Handler) ListCategory(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.CategoryV1().List, h.val.ValidateListCategoryRequest)
}
