// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package system

import (
	"time"

	"github.com/clin211/miniblog-v2/internal/pkg/core"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

// Healthz 服务健康检查.
func (h *Handler) Healthz(c *gin.Context) {
	core.WriteResponse(c, v1.HealthzResponse{
		Status:    v1.ServiceStatus_Healthy,
		Timestamp: time.Now().Format(time.DateTime),
	}, nil)
}
