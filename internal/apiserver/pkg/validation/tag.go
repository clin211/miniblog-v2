// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"strings"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	genericvalidation "github.com/clin211/miniblog-v2/pkg/validation"
)

func (v *Validator) ValidateTagRules() genericvalidation.Rules {
	return genericvalidation.Rules{
		"TagID": func(id any) error {
			if id.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("tagID cannot be empty")
			}
			return nil
		},
		"Name": func(value any) error {
			name := value.(string)
			if name == "" {
				return errno.ErrInvalidArgument.WithMessage("name cannot be empty")
			}
			if len(strings.TrimSpace(name)) == 0 {
				return errno.ErrInvalidArgument.WithMessage("name cannot be whitespace only")
			}
			if len(name) > 50 {
				return errno.ErrInvalidArgument.WithMessage("name cannot exceed 50 characters")
			}
			return nil
		},
		"Color": func(value any) error {
			color := value.(string)
			if color == "" {
				return errno.ErrInvalidArgument.WithMessage("color cannot be empty")
			}

			color = strings.TrimSpace(color)
			if len(color) == 0 {
				return errno.ErrInvalidArgument.WithMessage("color cannot be whitespace only")
			}

			return nil
		},
		// 分页参数校验
		"Limit": func(value any) error {
			if value.(int64) < 0 {
				return errno.ErrInvalidArgument.WithMessage("limit cannot be negative")
			}
			return nil
		},
		"Offset": func(value any) error {
			if value.(int64) < 0 {
				return errno.ErrInvalidArgument.WithMessage("offset cannot be negative")
			}
			return nil
		},
	}
}

// ValidateCreateTagRequest 校验 CreateTagRequest 结构体的有效性.
func (v *Validator) ValidateCreateTagRequest(ctx context.Context, rq *v1.CreateTagRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateUpdateTagRequest 校验 UpdateTagRequest 结构体的有效性.
func (v *Validator) ValidateUpdateTagRequest(ctx context.Context, rq *v1.UpdateTagRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateDeleteTagRequest 校验 DeleteTagRequest 结构体的有效性.
func (v *Validator) ValidateDeleteTagRequest(ctx context.Context, rq *v1.DeleteTagRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateGetTagRequest 校验 GetTagRequest 结构体的有效性.
func (v *Validator) ValidateGetTagRequest(ctx context.Context, rq *v1.GetTagRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateListTagRequest 校验 ListTagRequest 结构体的有效性.
func (v *Validator) ValidateListTagRequest(ctx context.Context, rq *v1.ListTagRequest) error {
	rules := v.ValidateTagRules()
	if err := rules["Offset"](rq.GetOffset()); err != nil {
		return err
	}
	if err := rules["Limit"](rq.GetLimit()); err != nil {
		return err
	}
	return nil
}
