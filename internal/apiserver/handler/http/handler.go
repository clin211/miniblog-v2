// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package http

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/biz"
)

// Handler 处理博客模块的请求.
type Handler struct {
	biz biz.IBiz
}

// NewHandler 创建新的 Handler 实例.
func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
