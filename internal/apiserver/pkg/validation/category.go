// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"strings"

	genericvalidation "github.com/clin211/miniblog-v2/pkg/validation"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
)

// ValidateCategoryRules 定义分类相关的校验规则
func (v *Validator) ValidateCategoryRules() genericvalidation.Rules {
	// 分类ID校验函数
	validateCategoryID := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("category ID cannot be empty")
			}
			return nil
		}
	}

	// 分类名称校验函数
	validateName := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var nameStr string
			var hasName bool

			switch v := value.(type) {
			case string:
				nameStr = v
				hasName = v != ""
			case *string:
				if v != nil {
					nameStr = *v
					hasName = *v != ""
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("name field type error")
			}

			if hasName {
				if len(strings.TrimSpace(nameStr)) == 0 {
					return errno.ErrInvalidArgument.WithMessage("category name cannot be empty or whitespace only")
				}
				if len(nameStr) > 50 {
					return errno.ErrInvalidArgument.WithMessage("category name is too long")
				}
			}
			return nil
		}
	}

	// 分类描述校验函数
	validateDescription := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var descStr string
			var hasDesc bool

			switch v := value.(type) {
			case string:
				descStr = v
				hasDesc = v != ""
			case *string:
				if v != nil {
					descStr = *v
					hasDesc = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("description field type error")
			}

			if hasDesc && len(descStr) > 200 {
				return errno.ErrInvalidArgument.WithMessage("category description is too long")
			}
			return nil
		}
	}

	// 分类图标校验函数
	validateIcon := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var iconStr string
			var hasIcon bool

			switch v := value.(type) {
			case string:
				iconStr = v
				hasIcon = v != ""
			case *string:
				if v != nil {
					iconStr = *v
					hasIcon = *v != ""
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("icon field type error")
			}

			if hasIcon {
				if len(strings.TrimSpace(iconStr)) == 0 {
					return errno.ErrInvalidArgument.WithMessage("category icon cannot be empty or whitespace only")
				}
				if len(iconStr) > 255 {
					return errno.ErrInvalidArgument.WithMessage("category icon is too long")
				}
				// 可以添加更多图标格式验证，比如只允许特定的图标名称
				// 这里简单验证不包含特殊字符
				if strings.ContainsAny(iconStr, "<>\"'&") {
					return errno.ErrInvalidArgument.WithMessage("category icon contains invalid characters")
				}
			}
			return nil
		}
	}

	// 分类主题校验函数
	validateTheme := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var themeStr string
			var hasTheme bool

			switch v := value.(type) {
			case string:
				themeStr = v
				hasTheme = v != ""
			case *string:
				if v != nil {
					themeStr = *v
					hasTheme = *v != ""
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("theme field type error")
			}

			if hasTheme {
				if len(strings.TrimSpace(themeStr)) == 0 {
					return errno.ErrInvalidArgument.WithMessage("category theme cannot be empty or whitespace only")
				}
				if len(themeStr) > 50 {
					return errno.ErrInvalidArgument.WithMessage("category theme is too long")
				}
				// 验证主题名称格式（只允许字母、数字、连字符）
				validThemes := []string{
					"cyan", "emerald", "violet", "orange", "red", "pink", "lime", "sky",
					"stone", "zinc", "yellow", "fuchsia", "neutral", "gray", "blue",
					"purple", "green", "amber", "rose", "indigo", "teal", "slate",
				}
				validTheme := false
				for _, valid := range validThemes {
					if themeStr == valid {
						validTheme = true
						break
					}
				}
				if !validTheme {
					return errno.ErrInvalidArgument.WithMessage("invalid category theme")
				}
			}
			return nil
		}
	}

	// 父分类ID校验函数
	validateParentID := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：int32 和 *int32
			var parentID int32
			var hasParentID bool

			switch v := value.(type) {
			case int32:
				parentID = v
				hasParentID = true
			case *int32:
				if v != nil {
					parentID = *v
					hasParentID = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("parent ID field type error")
			}

			if hasParentID && parentID < 0 {
				return errno.ErrInvalidArgument.WithMessage("parent ID cannot be negative")
			}
			return nil
		}
	}

	// 排序值校验函数
	validateSortOrder := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：int32 和 *int32
			var sortOrder int32
			var hasSortOrder bool

			switch v := value.(type) {
			case int32:
				sortOrder = v
				hasSortOrder = true
			case *int32:
				if v != nil {
					sortOrder = *v
					hasSortOrder = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("sort order field type error")
			}

			if hasSortOrder && sortOrder < 0 {
				return errno.ErrInvalidArgument.WithMessage("sort order cannot be negative")
			}
			return nil
		}
	}

	// 激活状态校验函数
	validateIsActive := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：v1.IsActive 和 *v1.IsActive
			switch val := value.(type) {
			case v1.IsActive:
				// 验证 enum 值是否有效
				if val != v1.IsActive_IS_ACTIVE_ACTIVE && val != v1.IsActive_IS_ACTIVE_DISABLED {
					return errno.ErrInvalidArgument.WithMessage("invalid IsActive enum value")
				}
				return nil
			case *v1.IsActive:
				if val != nil {
					// 验证 enum 值是否有效
					if *val != v1.IsActive_IS_ACTIVE_ACTIVE && *val != v1.IsActive_IS_ACTIVE_DISABLED {
						return errno.ErrInvalidArgument.WithMessage("invalid IsActive enum value")
					}
				}
				return nil
			default:
				return errno.ErrInvalidArgument.WithMessage("IsActive field must be of 0 or 1")
			}
		}
	}

	// 定义各字段的校验逻辑
	return genericvalidation.Rules{
		// 基本字段校验
		"CategoryID":  validateCategoryID(),
		"Name":        validateName(),
		"Description": validateDescription(),
		"Icon":        validateIcon(),
		"Theme":       validateTheme(),
		"ParentID":    validateParentID(),
		"SortOrder":   validateSortOrder(),
		"IsActive":    validateIsActive(),
	}
}

// ValidateCreateCategoryRequest 校验 CreateCategoryRequest 结构体的有效性
func (v *Validator) ValidateCreateCategoryRequest(ctx context.Context, rq *v1.CreateCategoryRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateCategoryRules())
}

// ValidateGetCategoryRequest 校验 GetCategoryRequest 结构体的有效性
func (v *Validator) ValidateGetCategoryRequest(ctx context.Context, rq *v1.GetCategoryRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateCategoryRules())
}

// ValidateDeleteCategoryRequest 校验 DeleteCategoryRequest 结构体的有效性
func (v *Validator) ValidateDeleteCategoryRequest(ctx context.Context, rq *v1.DeleteCategoryRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateCategoryRules())
}

// ValidateUpdateCategoryRequest 校验 UpdateCategoryRequest 结构体的有效性
func (v *Validator) ValidateUpdateCategoryRequest(ctx context.Context, rq *v1.UpdateCategoryRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateCategoryRules())
}

// ValidateListCategoryRequest 校验 ListCategoryRequest 结构体的有效性
func (v *Validator) ValidateListCategoryRequest(ctx context.Context, rq *v1.ListCategoryRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateCategoryRules())
}
