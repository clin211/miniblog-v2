// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package id

// Package id 提供了两种 ID 生成方案：基于 Sony 的分布式 ID 生成器 Sonyflake 和基于编码的字符串 ID 生成器

// NewCode 根据数字 ID 生成一个唯一的字符串编码
// 参数:
//   - id: 输入的数字 ID，通常由 Sonyflake.Id() 生成
//   - options: 配置函数列表，用于自定义编码生成过程
//
// 返回值:
//   - 一个字符串类型的编码
//
// 算法说明:
//  1. 先对输入 ID 进行扩大并添加盐值，增加安全性
//  2. 使用扩散算法让每一位数字相互影响，增加编码复杂度
//  3. 使用混淆算法（排列盒）重新排列编码，进一步提高安全性
//
// 使用示例:
//
//	id := sf.Id(context.Background())
//	code := id.NewCode(id, id.WithCodeL(10))
func NewCode(id uint64, options ...func(*CodeOptions)) string {
	ops := getCodeOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	// 扩大并添加盐值
	id = id*uint64(ops.n1) + ops.salt

	code := make([]rune, 0, ops.l) // 预分配容量以提高性能
	slIdx := make([]byte, ops.l)

	charLen := len(ops.chars)
	charLenUI := uint64(charLen)

	// 扩散
	for i := range ops.l {
		slIdx[i] = byte(id % charLenUI)                          // 获取每个数字
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(charLen) // 让个位数影响其他位
		id /= charLenUI                                          // 右移
	}

	// 混淆（https://en.wikipedia.org/wiki/Permutation_box）
	for i := range ops.l {
		idx := (byte(i) * byte(ops.n2)) % byte(ops.l)
		code = append(code, ops.chars[slIdx[idx]])
	}
	return string(code)
}
