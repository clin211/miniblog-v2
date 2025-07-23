// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package app

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/clin211/miniblog-v2/internal/pkg/core"
	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

// CreateDevice 创建新设备 - 演示MongoDB插入操作（无需登录）
func (h *Handler) CreateDevice(c *gin.Context) {
	log.Infow("MongoDB Create operation demo - open access")

	// 接收任意扁平化数据
	var deviceData any
	if err := c.ShouldBindJSON(&deviceData); err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	if err := h.biz.Device().Create(c, deviceData); err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, nil, nil)
}

// UpdateDevice 更新设备 - 演示MongoDB更新操作（无需登录）
func (h *Handler) UpdateDevice(c *gin.Context) {
	log.Infow("MongoDB Update operation demo - open access")

	deviceID := c.Param("id")
	if deviceID == "" {
		core.WriteResponse(c, nil, errors.New("device id is required"))
		return
	}

	// 接收更新数据
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	// if err := h.biz.Device().Update(c, device); err != nil {
	// 	core.WriteResponse(c, nil, err)
	// 	return
	// }

	core.WriteResponse(c, nil, nil)
}

// DeleteDevice 删除设备 - 演示MongoDB删除操作（无需登录）
func (h *Handler) DeleteDevice(c *gin.Context) {
	log.Infow("MongoDB Delete operation demo - open access")

	deviceID := c.Param("id")
	if deviceID == "" {
		core.WriteResponse(c, nil, errors.New("device id is required"))
		return
	}

	if err := h.biz.Device().Delete(c, deviceID); err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, gin.H{"message": "Device deleted successfully"}, nil)
}

// GetDevice 获取单个设备 - 演示MongoDB查询单个文档（无需登录）
func (h *Handler) GetDevice(c *gin.Context) {
	log.Infow("MongoDB Get operation demo - open access")

	deviceID := c.Param("id")
	if deviceID == "" {
		core.WriteResponse(c, nil, errors.New("device id is required"))
		return
	}

	device, err := h.biz.Device().Get(c, deviceID)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	core.WriteResponse(c, device, nil)
}

// ListDevices 获取设备列表 - 演示MongoDB分页查询（无需登录）
func (h *Handler) ListDevices(c *gin.Context) {
	log.Infow("MongoDB List operation demo - open access")

	// 解析分页参数
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	devices, total, err := h.biz.Device().List(c, limit, offset)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}

	response := gin.H{
		"devices": devices,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	}
	core.WriteResponse(c, response, nil)
}
