// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package http

import (
	"github.com/gin-gonic/gin"

	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

// Login 用户登录并返回 JWT Token.
func (h *Handler) Login(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Login, h.val.ValidateLoginRequest)
}

// RefreshToken 刷新 JWT Token.
func (h *Handler) RefreshToken(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().RefreshToken)
}

// ChangeUserPassword 修改用户密码.
func (h *Handler) ChangePassword(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().ChangePassword, h.val.ValidateChangePasswordRequest)
}

// CreateUser 创建新用户.
func (h *Handler) CreateUser(c *gin.Context) {
	log.Infow("CreateUser handler is called", "method", "CreateUser")
	core.HandleJSONRequest(c, h.biz.UserV1().Create, h.val.ValidateCreateUserRequest)
}

// UpdateUser 更新用户信息.
func (h *Handler) UpdateUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Update, h.val.ValidateUpdateUserRequest)
}

// DeleteUser 删除用户.
func (h *Handler) DeleteUser(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.UserV1().Delete, h.val.ValidateDeleteUserRequest)
}

// GetUser 获取用户信息.
func (h *Handler) GetUser(c *gin.Context) {
	core.HandleUriRequest(c, h.biz.UserV1().Get, h.val.ValidateGetUserRequest)
}

// ListUser 列出用户信息.
func (h *Handler) ListUser(c *gin.Context) {
	core.HandleQueryRequest(c, h.biz.UserV1().List, h.val.ValidateListUserRequest)
}
