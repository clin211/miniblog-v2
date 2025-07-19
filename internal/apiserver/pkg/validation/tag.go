// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"strings"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	appv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1/app"
	genericvalidation "github.com/onexstack/onexstack/pkg/validation"
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
			if len(name) > 20 {
				return errno.ErrInvalidArgument.WithMessage("name cannot exceed 20 characters")
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
	}
}

// ValidateCreateTagRequest 校验 CreateTagRequest 结构体的有效性.
func (v *Validator) ValidateCreateTagRequest(ctx context.Context, rq *appv1.CreateTagRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateUpdateTagRequest 校验 UpdateTagRequest 结构体的有效性.
func (v *Validator) ValidateUpdateTagRequest(ctx context.Context, rq *appv1.UpdateTagRequest) error {
	// 先校验 ID 字段
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("tag ID must be positive")
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateDeleteTagRequest 校验 DeleteTagRequest 结构体的有效性.
func (v *Validator) ValidateDeleteTagRequest(ctx context.Context, rq *appv1.DeleteTagRequest) error {
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("tag ID must be positive")
	}
	return nil
}

// ValidateGetTagRequest 校验 GetTagRequest 结构体的有效性.
func (v *Validator) ValidateGetTagRequest(ctx context.Context, rq *appv1.GetTagRequest) error {
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("tag ID must be positive")
	}
	return nil
}

// ValidateListTagRequest 校验 ListTagRequest 结构体的有效性.
func (v *Validator) ValidateListTagRequest(ctx context.Context, rq *appv1.ListTagRequest) error {
	// 校验可选的名称过滤参数
	if name := rq.GetName(); name != "" {
		if len(strings.TrimSpace(name)) == 0 {
			return errno.ErrInvalidArgument.WithMessage("name filter cannot be whitespace only")
		}
		if len(name) > 20 {
			return errno.ErrInvalidArgument.WithMessage("name filter cannot exceed 20 characters")
		}
	}
	return nil
}
