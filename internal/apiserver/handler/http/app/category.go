// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package app

import (
	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

func (h *Handler) ListCategories(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.CategoryV1().List, h.val.ValidateListCategoryRequest)
}

func (h *Handler) GetCategory(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.CategoryV1().Get, h.val.ValidateGetCategoryRequest)
}
