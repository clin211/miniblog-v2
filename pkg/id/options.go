// Package id 提供了两种 ID 生成方案：基于 Sony 的分布式 ID 生成器 Sonyflake 和基于编码的字符串 ID 生成器
package id

import "time"

// CodeOptions 定义了编码 ID 生成器的配置选项
type CodeOptions struct {
	chars []rune // 用于生成编码的字符集
	n1    int    // 用于扩散算法的参数，须与字符集长度互质
	n2    int    // 用于混淆算法的参数，须与编码长度互质
	l     int    // 生成的编码长度
	salt  uint64 // 随机盐值，用于增强安全性
}

// WithCodeChars 设置自定义字符集
// 参数:
//   - arr: 自定义的字符集
func WithCodeChars(arr []rune) func(*CodeOptions) {
	return func(options *CodeOptions) {
		if len(arr) > 0 {
			getCodeOptionsOrSetDefault(options).chars = arr
		}
	}
}

// WithCodeN1 设置 n1 参数
// 参数:
//   - n: 扩散算法的参数值，建议与字符集长度互质
func WithCodeN1(n int) func(*CodeOptions) {
	return func(options *CodeOptions) {
		getCodeOptionsOrSetDefault(options).n1 = n
	}
}

// WithCodeN2 设置 n2 参数
// 参数:
//   - n: 混淆算法的参数值，建议与编码长度互质
func WithCodeN2(n int) func(*CodeOptions) {
	return func(options *CodeOptions) {
		getCodeOptionsOrSetDefault(options).n2 = n
	}
}

// WithCodeL 设置编码长度
// 参数:
//   - l: 生成的编码长度
func WithCodeL(l int) func(*CodeOptions) {
	return func(options *CodeOptions) {
		if l > 0 {
			getCodeOptionsOrSetDefault(options).l = l
		}
	}
}

// WithCodeSalt 设置随机盐值
// 参数:
//   - salt: 用于增强安全性的随机数
func WithCodeSalt(salt uint64) func(*CodeOptions) {
	return func(options *CodeOptions) {
		if salt > 0 {
			getCodeOptionsOrSetDefault(options).salt = salt
		}
	}
}

// getCodeOptionsOrSetDefault 获取 CodeOptions，如果为空则创建默认配置
// 参数:
//   - options: 可选的 CodeOptions 实例
//
// 返回值:
//   - 配置好的 CodeOptions 实例
func getCodeOptionsOrSetDefault(options *CodeOptions) *CodeOptions {
	if options == nil {
		return &CodeOptions{
			// 基础字符集，移除了 0,1,I,O,U,Z
			chars: []rune{
				'2', '3', '4', '5', '6',
				'7', '8', '9', 'A', 'B',
				'C', 'D', 'E', 'F', 'G',
				'H', 'J', 'K', 'L', 'M',
				'N', 'P', 'Q', 'R', 'S',
				'T', 'V', 'W', 'X', 'Y',
			},
			// n1 / len(chars)=30 互质
			n1: 17,
			// n2 / l 互质
			n2: 5,
			// 编码长度
			l: 8,
			// 随机数
			salt: 123567369,
		}
	}
	return options
}

// SonyflakeOptions 定义了 Sonyflake ID 生成器的配置选项
type SonyflakeOptions struct {
	machineId uint16    // 机器 ID，用于分布式环境中标识不同机器
	startTime time.Time // 时间起点，Sonyflake 从该时间点开始计算时间差
}

// WithSonyflakeMachineId 设置机器 ID
// 参数:
//   - id: 机器 ID，用于在分布式环境中唯一标识服务器
func WithSonyflakeMachineId(id uint16) func(*SonyflakeOptions) {
	return func(options *SonyflakeOptions) {
		if id > 0 {
			getSonyflakeOptionsOrSetDefault(options).machineId = id
		}
	}
}

// WithSonyflakeStartTime 设置时间起点
// 参数:
//   - startTime: Sonyflake 从该时间点开始计算时间差
func WithSonyflakeStartTime(startTime time.Time) func(*SonyflakeOptions) {
	return func(options *SonyflakeOptions) {
		if !startTime.IsZero() {
			getSonyflakeOptionsOrSetDefault(options).startTime = startTime
		}
	}
}

// getSonyflakeOptionsOrSetDefault 获取 SonyflakeOptions，如果为空则创建默认配置
// 参数:
//   - options: 可选的 SonyflakeOptions 实例
//
// 返回值:
//   - 配置好的 SonyflakeOptions 实例
func getSonyflakeOptionsOrSetDefault(options *SonyflakeOptions) *SonyflakeOptions {
	if options == nil {
		return &SonyflakeOptions{
			machineId: 1,
			startTime: time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC),
		}
	}
	return options
}
