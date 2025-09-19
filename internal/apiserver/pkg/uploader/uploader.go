// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package uploader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	opt "github.com/clin211/miniblog-v2/pkg/options"
)

// UploadedObject 为上传成功后的对象信息。
type UploadedObject struct {
	Provider string
	Key      string
	URL      string
	Size     int64
	MIME     string
	Hash     string
	Metadata map[string]string
}

// Uploader 定义上传抽象
type Uploader interface {
	Upload(ctx context.Context, filename string, mimeType string, r io.Reader) (*UploadedObject, error)
}

// NewLocalOrFromConfig 根据配置返回 Uploader（默认 local）
func NewLocalOrFromConfig(cfg *opt.UploadOptions) Uploader {
	if cfg == nil || strings.ToLower(cfg.Provider) == "local" {
		if cfg == nil {
			cfg = opt.NewUploadOptions()
		}
		return &localUploader{cfg: cfg}
	}
	// 仅占位：当前仅实现本地
	return &localUploader{cfg: cfg}
}

type localUploader struct {
	cfg *opt.UploadOptions
}

func (l *localUploader) Upload(ctx context.Context, filename string, mimeType string, r io.Reader) (*UploadedObject, error) {
	// 解析扩展名
	ext := filepath.Ext(filename)
	if ext == "" {
		if exts, _ := mime.ExtensionsByType(mimeType); len(exts) > 0 {
			ext = exts[0]
		}
	}

	// 读取并计算哈希
	hasher := sha256.New()
	tmpPath := filepath.Join(l.cfg.Local.BaseDir, ".tmp", fmt.Sprintf("%d.tmp", time.Now().UnixNano()))
	if err := os.MkdirAll(filepath.Dir(tmpPath), os.FileMode(l.cfg.Local.MkdirPerm)); err != nil {
		return nil, err
	}
	f, err := os.Create(tmpPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	n, err := io.Copy(io.MultiWriter(f, hasher), r)
	if err != nil {
		_ = os.Remove(tmpPath)
		return nil, err
	}

	sum := hasher.Sum(nil)
	sumHex := hex.EncodeToString(sum)

	// 生成对象键（简化版：日期/哈希/扩展名）
	datePrefix := time.Now().Format("2006/01/02")
	key := fmt.Sprintf("%s/%s%s", datePrefix, sumHex[:16], ext)
	absPath := filepath.Join(l.cfg.Local.BaseDir, key)

	// 确保目录存在并移动临时文件
	if err := os.MkdirAll(filepath.Dir(absPath), os.FileMode(l.cfg.Local.MkdirPerm)); err != nil {
		_ = os.Remove(tmpPath)
		return nil, err
	}
	if err := os.Rename(tmpPath, absPath); err != nil {
		_ = os.Remove(tmpPath)
		return nil, err
	}

	// 不做任何图片转换或衍生文件生成
	var publicKey = key
	var publicSize = n
	var publicMIME = mimeType

	// 构造 URL：
	// - 对外访问优先使用 Local.BaseURL（包含协议/域名/路径），满足不同环境公网前缀不同的需求；
	// - 若未配置 BaseURL，则回退为 Local.BasePath（仅路径），由前端拼接域名或同源访问。
	var publicURL string
	if l.cfg != nil && l.cfg.Local != nil && strings.TrimSpace(l.cfg.Local.BaseURL) != "" {
		if uu, err := url.Parse(l.cfg.Local.BaseURL); err == nil {
			uu.Path = strings.TrimRight(uu.Path, "/") + "/" + publicKey
			publicURL = uu.String()
		}
	}
	if publicURL == "" {
		var pathBase string
		if l.cfg != nil && l.cfg.Local != nil {
			if strings.TrimSpace(l.cfg.Local.BasePath) != "" {
				pathBase = l.cfg.Local.BasePath
			} else if strings.TrimSpace(l.cfg.Local.BaseURL) != "" { // 兼容：若只给了 BaseURL，则取其 path
				if uu, err := url.Parse(l.cfg.Local.BaseURL); err == nil {
					pathBase = uu.Path
				}
			}
		}
		if pathBase == "" {
			pathBase = "/static/uploads"
		}
		u, _ := url.Parse(pathBase)
		u.Path = strings.TrimRight(u.Path, "/") + "/" + publicKey
		publicURL = u.String()
	}

	return &UploadedObject{
		Provider: "local",
		Key:      publicKey,
		URL:      publicURL,
		Size:     publicSize,
		MIME:     publicMIME,
		Hash:     sumHex, // sha256
		Metadata: map[string]string{},
	}, nil
}
