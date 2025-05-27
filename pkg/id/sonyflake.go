// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

// Package id 提供了两种 ID 生成方案：基于 Sony 的分布式 ID 生成器 Sonyflake 和基于编码的字符串 ID 生成器
package id

import (
	"context"
	"fmt"
	"time"

	"github.com/sony/sonyflake"
)

// Sonyflake 结构体封装了 Sony 的分布式 ID 生成器
// 它可以生成分布式环境下的唯一 ID，可用于数据库主键等场景
type Sonyflake struct {
	ops   SonyflakeOptions     // Sonyflake 的配置选项
	sf    *sonyflake.Sonyflake // 内部的 Sonyflake 实例
	Error error                // 初始化过程中的错误信息
}

// NewSonyflake 创建一个新的 Sonyflake 实例
// 参数:
//   - options: 配置函数列表，用于自定义 Sonyflake 的行为
//
// 返回值:
//   - 一个配置好的 Sonyflake 实例
//
// 使用示例:
//
//	sf := id.NewSonyflake(id.WithSonyflakeMachineId(1))
//	id := sf.Id(context.Background())
func NewSonyflake(options ...func(*SonyflakeOptions)) *Sonyflake {
	ops := getSonyflakeOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	sf := &Sonyflake{
		ops: *ops,
	}
	st := sonyflake.Settings{
		StartTime: ops.startTime,
	}
	if ops.machineId > 0 {
		st.MachineID = func() (uint16, error) {
			return ops.machineId, nil
		}
	}
	ins := sonyflake.NewSonyflake(st)
	if ins == nil {
		sf.Error = fmt.Errorf("create snoyflake failed")
	}
	_, err := ins.NextID()
	if err != nil {
		sf.Error = fmt.Errorf("invalid start time")
	}
	sf.sf = ins
	return sf
}

// Id 生成一个新的唯一 ID
// 参数:
//   - ctx: 上下文对象，可用于取消操作
//
// 返回值:
//   - 一个 uint64 类型的唯一 ID
//
// 备注:
//   - 如果生成失败会进行指数退避重试
func (s *Sonyflake) Id(ctx context.Context) (id uint64) {
	if s.Error != nil {
		return
	}
	var err error
	id, err = s.sf.NextID()
	if err == nil {
		return
	}

	sleep := 1
	for {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		id, err = s.sf.NextID()
		if err == nil {
			return
		}
		sleep *= 2
	}
}
