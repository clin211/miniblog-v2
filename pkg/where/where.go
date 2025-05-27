// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

// Package where 提供了一个简洁的 GORM 查询条件构建器，用于构建数据库查询条件
package where

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// defaultLimit 定义分页的默认限制数量
	defaultLimit = -1
)

// Tenant 表示一个租户，包含一个键和一个用于获取其值的函数
type Tenant struct {
	Key       string                           // 与租户关联的键
	ValueFunc func(ctx context.Context) string // 基于上下文获取租户值的函数
}

// Where 定义了一个接口，用于修改 GORM 数据库查询
type Where interface {
	Where(db *gorm.DB) *gorm.DB
}

// Query 表示带参数的数据库查询
// 它包含查询条件和任何相关的参数
type Query struct {
	// Query 保存要在 GORM 查询中使用的条件
	// 可以是字符串、映射、结构体或 GORM 的 Where 子句支持的其他类型
	Query interface{}

	// Args 保存将传递给查询条件的参数
	// 这些值将用于替换查询中的占位符
	Args []interface{}
}

// Option 定义了一个修改 Options 的函数类型
type Option func(*Options)

// Options 保存 GORM 的 Where 查询条件的选项
type Options struct {
	// Offset 定义分页的起始点
	// +可选
	Offset int `json:"offset"`
	// Limit 定义要返回的最大结果数
	// +可选
	Limit int `json:"limit"`
	// Filters 包含用于过滤记录的键值对
	Filters map[any]any
	// Clauses 包含要附加到查询的自定义子句
	Clauses []clause.Expression
	// Queries 包含要执行的查询列表
	Queries []Query
}

// tenant 保存已注册的租户实例
var registeredTenant Tenant

// WithOffset 使用给定的偏移值初始化 Options 中的 Offset 字段
func WithOffset(offset int64) Option {
	return func(whr *Options) {
		if offset < 0 {
			offset = 0
		}
		whr.Offset = int(offset)
	}
}

// WithLimit 使用给定的限制值初始化 Options 中的 Limit 字段
func WithLimit(limit int64) Option {
	return func(whr *Options) {
		if limit <= 0 {
			limit = defaultLimit
		}
		whr.Limit = int(limit)
	}
}

// WithPage 是一个语法糖函数，用于将页码和页面大小转换为 Options 中的 limit 和 offset
// 此函数通常用于业务逻辑中以便于分页
func WithPage(page int, pageSize int) Option {
	return func(whr *Options) {
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = defaultLimit
		}

		whr.Offset = (page - 1) * pageSize
		whr.Limit = pageSize
	}
}

// WithFilter 使用给定的过滤条件初始化 Options 中的 Filters 字段
func WithFilter(filter map[any]any) Option {
	return func(whr *Options) {
		whr.Filters = filter
	}
}

// WithClauses 将子句追加到 Options 中的 Clauses 字段
func WithClauses(conds ...clause.Expression) Option {
	return func(whr *Options) {
		whr.Clauses = append(whr.Clauses, conds...)
	}
}

// WithQuery 创建一个 Option，将带参数的查询条件添加到 Options 结构体中
// 查询参数可以是字符串、映射、结构体或 GORM 的 Where 子句支持的任何其他类型
// args 参数包含将替换查询字符串中占位符的值
func WithQuery(query interface{}, args ...interface{}) Option {
	return func(whr *Options) {
		whr.Queries = append(whr.Queries, Query{Query: query, Args: args})
	}
}

// NewWhere 构造一个新的 Options 对象，应用给定的 where 选项
func NewWhere(opts ...Option) *Options {
	whr := &Options{
		Offset:  0,
		Limit:   defaultLimit,
		Filters: map[any]any{},
		Clauses: make([]clause.Expression, 0),
	}

	for _, opt := range opts {
		opt(whr) // 将每个 Option 应用到 opts
	}

	return whr
}

// O 设置查询的偏移量
func (whr *Options) O(offset int) *Options {
	if offset < 0 {
		offset = 0
	}
	whr.Offset = offset
	return whr
}

// L 设置查询的限制数量
func (whr *Options) L(limit int) *Options {
	if limit <= 0 {
		limit = defaultLimit // 确保 defaultLimit 在其他地方定义
	}
	whr.Limit = limit
	return whr
}

// P 根据页码和页面大小设置分页
func (whr *Options) P(page int, pageSize int) *Options {
	if page < 1 {
		page = 1 // 确保页码至少为 1
	}
	if pageSize <= 0 {
		pageSize = defaultLimit // 确保 defaultLimit 在其他地方定义
	}
	whr.Offset = (page - 1) * pageSize
	whr.Limit = pageSize
	return whr
}

// C 向查询添加条件
func (whr *Options) C(conds ...clause.Expression) *Options {
	whr.Clauses = append(whr.Clauses, conds...)
	return whr
}

// Q 向 Options 结构体添加带参数的查询条件并返回修改后的 Options
// 此方法将一个新的 Query 实例追加到 Queries 切片中
func (whr *Options) Q(query interface{}, args ...interface{}) *Options {
	whr.Queries = append(whr.Queries, Query{Query: query, Args: args})
	return whr
}

// T 使用提供的上下文获取与已注册租户关联的值
func (whr *Options) T(ctx context.Context) *Options {
	if registeredTenant.Key != "" && registeredTenant.ValueFunc != nil {
		whr.F(registeredTenant.Key, registeredTenant.ValueFunc(ctx))
	}
	return whr
}

// F 向查询添加过滤条件
func (whr *Options) F(kvs ...any) *Options {
	if len(kvs)%2 != 0 {
		// 处理错误：键值对数量不均
		return whr
	}

	for i := 0; i < len(kvs); i += 2 {
		key := kvs[i]
		value := kvs[i+1]
		whr.Filters[key] = value
	}

	return whr
}

// Where 将过滤条件和子句应用到给定的 gorm.DB 实例
func (whr *Options) Where(db *gorm.DB) *gorm.DB {
	for _, query := range whr.Queries {
		conds := db.Statement.BuildCondition(query.Query, query.Args...)
		whr.Clauses = append(whr.Clauses, conds...)
	}
	return db.Where(whr.Filters).Clauses(whr.Clauses...).Offset(whr.Offset).Limit(whr.Limit)
}

// O 是一个便捷函数，用于创建带偏移量的新 Options
func O(offset int) *Options {
	return NewWhere().O(offset)
}

// L 是一个便捷函数，用于创建带限制数量的新 Options
func L(limit int) *Options {
	return NewWhere().L(limit)
}

// P 是一个便捷函数，用于创建带页码和页面大小的新 Options
func P(page int, pageSize int) *Options {
	return NewWhere().P(page, pageSize)
}

// C 是一个便捷函数，用于创建带条件的新 Options
func C(conds ...clause.Expression) *Options {
	return NewWhere().C(conds...)
}

// T 是一个便捷函数，用于创建带租户的新 Options
func T(ctx context.Context) *Options {
	return NewWhere().F(registeredTenant.Key, registeredTenant.ValueFunc(ctx))
}

// F 是一个便捷函数，用于创建带过滤条件的新 Options
func F(kvs ...any) *Options {
	return NewWhere().F(kvs...)
}

// RegisterTenant 注册一个新的租户，使用指定的键和值函数
func RegisterTenant(key string, valueFunc func(context.Context) string) {
	registeredTenant = Tenant{
		Key:       key,
		ValueFunc: valueFunc,
	}
}
