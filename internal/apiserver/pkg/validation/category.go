// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"strings"

	genericvalidation "github.com/onexstack/onexstack/pkg/validation"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
)

// ValidateCategoryRules 定义分类相关的校验规则
func (v *Validator) ValidateCategoryRules() genericvalidation.Rules {
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
			// 处理两种可能的类型：bool 和 *bool
			switch value.(type) {
			case bool, *bool:
				return nil
			default:
				return errno.ErrInvalidArgument.WithMessage("is active field type error")
			}
		}
	}

	// 定义各字段的校验逻辑
	return genericvalidation.Rules{
		// 基本字段校验
		"ID": func(value any) error {
			id := value.(int32)
			if id <= 0 {
				return errno.ErrInvalidArgument.WithMessage("category ID must be positive")
			}
			return nil
		},
		"Name":        validateName(),
		"Description": validateDescription(),
		"ParentID":    validateParentID(),
		"SortOrder":   validateSortOrder(),
		"IsActive":    validateIsActive(),
	}
}

// ValidateCreateCategoryRequest 校验 CreateCategoryRequest 结构体的有效性
func (v *Validator) ValidateCreateCategoryRequest(ctx context.Context, rq *v1.CreateCategoryRequest) error {

	rules := v.ValidateCategoryRules()
	if rq.Name != "" {
		if err := rules["Name"](rq.Name); err != nil {
			return err
		}
	}
	if rq.Description != nil {
		if err := rules["Description"](rq.Description); err != nil {
			return err
		}
	}
	if rq.ParentID != nil {
		if err := rules["ParentID"](rq.ParentID); err != nil {
			return err
		}
	}
	if rq.SortOrder != nil {
		if err := rules["SortOrder"](rq.SortOrder); err != nil {
			return err
		}
	}
	if rq.IsActive != nil {
		if err := rules["IsActive"](rq.IsActive); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateCategoryRequest 校验 UpdateCategoryRequest 结构体的有效性
func (v *Validator) ValidateUpdateCategoryRequest(ctx context.Context, rq *v1.UpdateCategoryRequest) error {
	// 校验必填的分类ID
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("category ID must be positive")
	}

	// 校验可选字段
	rules := v.ValidateCategoryRules()
	if rq.Name != nil {
		if err := rules["Name"](*rq.Name); err != nil {
			return err
		}
	}
	if rq.Description != nil {
		if err := rules["Description"](rq.Description); err != nil {
			return err
		}
	}
	if rq.ParentID != nil {
		if err := rules["ParentID"](rq.ParentID); err != nil {
			return err
		}
	}
	if rq.SortOrder != nil {
		if err := rules["SortOrder"](rq.SortOrder); err != nil {
			return err
		}
	}
	if rq.IsActive != nil {
		if err := rules["IsActive"](rq.IsActive); err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteCategoryRequest 校验 DeleteCategoryRequest 结构体的有效性
func (v *Validator) ValidateDeleteCategoryRequest(ctx context.Context, rq *v1.DeleteCategoryRequest) error {
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("category ID must be positive")
	}
	return nil
}

// ValidateGetCategoryRequest 校验 GetCategoryRequest 结构体的有效性
func (v *Validator) ValidateGetCategoryRequest(ctx context.Context, rq *v1.GetCategoryRequest) error {
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("category ID must be positive")
	}
	return nil
}

// ValidateListCategoryRequest 校验 ListCategoryRequest 结构体的有效性
func (v *Validator) ValidateListCategoryRequest(ctx context.Context, rq *v1.ListCategoryRequest) error {
	// 校验可选的父分类ID过滤参数
	if rq.ParentID != nil {
		if *rq.ParentID < 0 {
			return errno.ErrInvalidArgument.WithMessage("parent ID filter cannot be negative")
		}
	}

	// 校验可选的激活状态过滤参数（bool类型无需特殊校验）
	return nil
}
