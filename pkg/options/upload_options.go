// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package options

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// UploadOptions 定义上传相关配置。
type UploadOptions struct {
	Provider     string           `json:"provider" mapstructure:"provider"`
	MaxSize      string           `json:"maxSize" mapstructure:"maxSize"` // e.g. "20MB"
	AllowedMIMEs []string         `json:"allowedMIMEs" mapstructure:"allowedMIMEs"`
	Deduplicate  bool             `json:"deduplicate" mapstructure:"deduplicate"`
	KeyTemplate  string           `json:"keyTemplate" mapstructure:"keyTemplate"`
	Local        *LocalOptions    `json:"local" mapstructure:"local"`
	AliOSS       *AliOSSOptions   `json:"alioss" mapstructure:"alioss"`
	Multipart    *MultipartConfig `json:"multipart" mapstructure:"multipart"`
}

type LocalOptions struct {
	BaseDir   string `json:"baseDir" mapstructure:"baseDir"`
	BasePath  string `json:"basePath" mapstructure:"basePath"`
	BaseURL   string `json:"baseURL" mapstructure:"baseURL"`
	MkdirPerm uint32 `json:"mkdirPerm" mapstructure:"mkdirPerm"`
}

type AliOSSOptions struct {
	Endpoint        string `json:"endpoint" mapstructure:"endpoint"`
	Bucket          string `json:"bucket" mapstructure:"bucket"`
	AccessKeyID     string `json:"accessKeyID" mapstructure:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret" mapstructure:"accessKeySecret"`
	SecurityToken   string `json:"securityToken" mapstructure:"securityToken"`
	ACL             string `json:"acl" mapstructure:"acl"`
	PathStyleAccess bool   `json:"pathStyleAccess" mapstructure:"pathStyleAccess"`
	UploadTimeout   string `json:"uploadTimeout" mapstructure:"uploadTimeout"`
	CallbackURL     string `json:"callbackURL" mapstructure:"callbackURL"`
}

type MultipartConfig struct {
	Enabled                 bool   `json:"enabled" mapstructure:"enabled"`
	MinSize                 string `json:"minSize" mapstructure:"minSize"`
	MaxParts                int32  `json:"maxParts" mapstructure:"maxParts"`
	PartSizeMin             string `json:"partSizeMin" mapstructure:"partSize.min"`
	PartSizeMax             string `json:"partSizeMax" mapstructure:"partSize.max"`
	PartSizeDefault         string `json:"partSizeDefault" mapstructure:"partSize.default"`
	ConcurrencyLimitPerUser int32  `json:"concurrencyLimitPerUser" mapstructure:"concurrencyLimitPerUser"`
	UploadIDTTL             string `json:"uploadIdTTL" mapstructure:"uploadIdTTL"`
	TempDir                 string `json:"tempDir" mapstructure:"tempDir"`
	Checksum                string `json:"checksum" mapstructure:"checksum"`
	PresignEnabled          bool   `json:"presignEnabled" mapstructure:"presign.enabled"`
	PresignExpires          string `json:"presignExpires" mapstructure:"presign.expires"`
	PresignMode             string `json:"presignMode" mapstructure:"presign.mode"`
}

// NewUploadOptions 返回带默认值的 UploadOptions。
func NewUploadOptions() *UploadOptions {
	return &UploadOptions{
		Provider:     "local",
		MaxSize:      "20MB",
		AllowedMIMEs: []string{"image/jpeg", "image/png", "image/gif", "image/webp", "application/pdf"},
		Deduplicate:  true,
		KeyTemplate:  "{date:2006/01/02}/{sha256:16}{ext}",
		Local: &LocalOptions{
			BaseDir:   "/data/miniblog/uploads",
			BasePath:  "/static/uploads",
			BaseURL:   "http://localhost:5555/static/uploads",
			MkdirPerm: 0755,
		},
		AliOSS: &AliOSSOptions{
			ACL:             "private",
			PathStyleAccess: false,
			UploadTimeout:   "30s",
		},
		Multipart: &MultipartConfig{
			Enabled:                 true,
			MinSize:                 "8MB",
			MaxParts:                10000,
			PartSizeMin:             "5MB",
			PartSizeMax:             "128MB",
			PartSizeDefault:         "8MB",
			ConcurrencyLimitPerUser: 2,
			UploadIDTTL:             "24h",
			TempDir:                 "/data/miniblog/uploads/.multipart",
			Checksum:                "sha256",
			PresignEnabled:          true,
			PresignExpires:          "15m",
			PresignMode:             "direct",
		},
	}
}

// ParseSize 解析如 20MB/8MB/1GB 等表示，返回字节数。
func ParseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, nil
	}
	re := regexp.MustCompile(`^([0-9]+)(B|KB|KIB|MB|MIB|GB|GIB)$`)
	m := re.FindStringSubmatch(s)
	if len(m) != 3 {
		// 尝试纯数字（字节）
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return n, nil
		}
		return 0, fmt.Errorf("invalid size: %s", s)
	}
	n, _ := strconv.ParseInt(m[1], 10, 64)
	unit := m[2]
	switch unit {
	case "B":
		return n, nil
	case "KB", "KIB":
		return n * 1024, nil
	case "MB", "MIB":
		return n * 1024 * 1024, nil
	case "GB", "GIB":
		return n * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("invalid size unit: %s", unit)
	}
}

// ParseDuration 包装 time.ParseDuration，允许空值。
func ParseDuration(s string) (time.Duration, error) {
	if strings.TrimSpace(s) == "" {
		return 0, nil
	}
	return time.ParseDuration(s)
}
