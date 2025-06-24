// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package http

import (
	"github.com/gin-gonic/gin"

	"github.com/clin211/miniblog-v2/internal/pkg/core"
)

// CreateTag 创建新标签.
func (h *Handler) CreateTag(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.TagV1().Create, h.val.ValidateCreateTagRequest)
}

// UpdateTag 更新标签信息.
func (h *Handler) UpdateTag(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.TagV1().Update, h.val.ValidateUpdateTagRequest)
}

// DeleteTag 删除标签.
func (h *Handler) DeleteTag(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.TagV1().Delete, h.val.ValidateDeleteTagRequest)
}

// GetTag 获取标签信息.
func (h *Handler) GetTag(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.TagV1().Get, h.val.ValidateGetTagRequest)
}

// ListTag 列出标签信息.
func (h *Handler) ListTag(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.TagV1().List, h.val.ValidateListTagRequest)
}
