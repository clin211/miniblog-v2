// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import "github.com/clin211/miniblog-v2/internal/pkg/log"

// Logger 是一个实现了 Logger 接口的日志记录器
// 它使用 log 包来记录带有附加上下文的错误信息
type Logger struct{}

// NewLogger 创建并返回一个新的 Logger 实例
func NewLogger() *Logger {
	return &Logger{}
}

// Error 使用 log 包记录带有提供上下文的错误信息
func (l *Logger) Error(err error, msg string, kvs ...any) {
	log.Errorw(msg, append(kvs, "err", err)...)
}
