// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package system

import (
	"github.com/clin211/miniblog-v2/internal/pkg/core"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	"github.com/gin-gonic/gin"
)

// UploadFile 单文件上传。
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}
	fh, err := file.Open()
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}
	defer fh.Close()

	var mime string
	if file.Header != nil {
		mime = file.Header.Get("Content-Type")
	}
	obj, err := h.upl.Upload(c.Request.Context(), file.Filename, mime, fh)
	if err != nil {
		core.WriteResponse(c, nil, err)
		return
	}
	core.WriteResponse(c, &v1.UploadedObject{
		Provider: obj.Provider,
		Key:      obj.Key,
		Url:      obj.URL,
		Size:     obj.Size,
		Mime:     obj.MIME,
		Hash:     obj.Hash,
		Metadata: obj.Metadata,
	}, nil)
}

// InitMultipart 初始化分片上传（占位）。
func (h *Handler) InitMultipart(c *gin.Context) {
	core.WriteResponse(c, &v1.InitMultipartResponse{}, nil)
}

// PresignParts 预签名（占位）。
func (h *Handler) PresignParts(c *gin.Context) {
	core.WriteResponse(c, &v1.PresignPartsResponse{}, nil)
}

// UploadPart 上传分片（占位）。
func (h *Handler) UploadPart(c *gin.Context) { core.WriteResponse(c, &v1.UploadPartResponse{}, nil) }

// ListParts 查询分片（占位）。
func (h *Handler) ListParts(c *gin.Context) { core.WriteResponse(c, &v1.ListPartsResponse{}, nil) }

// CompleteMultipart 完成分片（占位）。
func (h *Handler) CompleteMultipart(c *gin.Context) {
	core.WriteResponse(c, &v1.CompleteMultipartResponse{}, nil)
}

// AbortMultipart 中止分片（占位）。
func (h *Handler) AbortMultipart(c *gin.Context) {
	core.WriteResponse(c, &v1.AbortMultipartResponse{}, nil)
}
