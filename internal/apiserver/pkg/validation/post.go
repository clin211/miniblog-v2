// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"net/url"
	"strings"

	genericvalidation "github.com/onexstack/onexstack/pkg/validation"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	apiv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
)

// ValidatePostRules 定义文章相关的校验规则
func (v *Validator) ValidatePostRules() genericvalidation.Rules {
	// 文章类型校验函数
	validatePostType := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：PostType 和 *PostType
			var postType apiv1.PostType
			var hasPostType bool

			switch v := value.(type) {
			case apiv1.PostType:
				postType = v
				hasPostType = true
			case *apiv1.PostType:
				if v != nil {
					postType = *v
					hasPostType = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("post type field type error")
			}

			if hasPostType {
				// 只允许已定义的枚举值
				switch postType {
				case apiv1.PostType_POST_TYPE_UNSPECIFIED, apiv1.PostType_POST_TYPE_ORIGINAL,
					apiv1.PostType_POST_TYPE_REPOST, apiv1.PostType_POST_TYPE_CONTRIBUTION:
					return nil
				default:
					return errno.ErrInvalidArgument.WithMessage("invalid post type value")
				}
			}
			return nil
		}
	}

	// 文章状态校验函数
	validatePostStatus := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：PostStatus 和 *PostStatus
			var status apiv1.PostStatus
			var hasStatus bool

			switch v := value.(type) {
			case apiv1.PostStatus:
				status = v
				hasStatus = true
			case *apiv1.PostStatus:
				if v != nil {
					status = *v
					hasStatus = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("post status field type error")
			}

			if hasStatus {
				// 只允许已定义的枚举值
				switch status {
				case apiv1.PostStatus_POST_STATUS_UNSPECIFIED, apiv1.PostStatus_POST_STATUS_DRAFT,
					apiv1.PostStatus_POST_STATUS_PUBLISHED, apiv1.PostStatus_POST_STATUS_ARCHIVED:
					return nil
				default:
					return errno.ErrInvalidArgument.WithMessage("invalid post status value")
				}
			}
			return nil
		}
	}

	// 封面URL校验函数
	validateCover := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var coverStr string
			var hasCover bool

			switch v := value.(type) {
			case string:
				coverStr = v
				hasCover = v != ""
			case *string:
				if v != nil {
					coverStr = *v
					hasCover = *v != ""
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("cover field type error")
			}

			if hasCover {
				if _, err := url.ParseRequestURI(coverStr); err != nil {
					return errno.ErrInvalidArgument.WithMessage("cover must be a valid URL")
				}
			}
			return nil
		}
	}

	// 原文链接校验函数
	validateOriginalSource := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var sourceStr string
			var hasSource bool

			switch v := value.(type) {
			case string:
				sourceStr = v
				hasSource = v != ""
			case *string:
				if v != nil {
					sourceStr = *v
					hasSource = *v != ""
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("original source field type error")
			}

			if hasSource {
				if _, err := url.ParseRequestURI(sourceStr); err != nil {
					return errno.ErrInvalidArgument.WithMessage("original source must be a valid URL")
				}
			}
			return nil
		}
	}

	// 文章摘要校验函数
	validateSummary := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var summaryStr string
			var hasSummary bool

			switch v := value.(type) {
			case string:
				summaryStr = v
				hasSummary = v != ""
			case *string:
				if v != nil {
					summaryStr = *v
					hasSummary = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("summary field type error")
			}

			if hasSummary && len(summaryStr) > 500 {
				return errno.ErrInvalidArgument.WithMessage("summary cannot exceed 500 characters")
			}
			return nil
		}
	}

	// 原作者姓名校验函数
	validateOriginalAuthor := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var authorStr string
			var hasAuthor bool

			switch v := value.(type) {
			case string:
				authorStr = v
				hasAuthor = v != ""
			case *string:
				if v != nil {
					authorStr = *v
					hasAuthor = *v != ""
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("original author field type error")
			}

			if hasAuthor {
				if len(strings.TrimSpace(authorStr)) == 0 {
					return errno.ErrInvalidArgument.WithMessage("original author cannot be empty or whitespace only")
				}
				if len(authorStr) > 100 {
					return errno.ErrInvalidArgument.WithMessage("original author name cannot exceed 100 characters")
				}
			}
			return nil
		}
	}

	// 原作者简介校验函数
	validateOriginalAuthorIntro := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：string 和 *string
			var introStr string
			var hasIntro bool

			switch v := value.(type) {
			case string:
				introStr = v
				hasIntro = v != ""
			case *string:
				if v != nil {
					introStr = *v
					hasIntro = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("original author intro field type error")
			}

			if hasIntro && len(introStr) > 200 {
				return errno.ErrInvalidArgument.WithMessage("original author intro cannot exceed 200 characters")
			}
			return nil
		}
	}

	// 分类ID校验函数
	validateCategoryID := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：int32 和 *int32
			var categoryID int32
			var hasCategoryID bool

			switch v := value.(type) {
			case int32:
				categoryID = v
				hasCategoryID = true
			case *int32:
				if v != nil {
					categoryID = *v
					hasCategoryID = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("category ID field type error")
			}

			if hasCategoryID && categoryID < 0 {
				return errno.ErrInvalidArgument.WithMessage("category ID cannot be negative")
			}
			return nil
		}
	}

	// 排序位置校验函数
	validatePosition := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			// 处理两种可能的类型：int32 和 *int32
			var position int32
			var hasPosition bool

			switch v := value.(type) {
			case int32:
				position = v
				hasPosition = true
			case *int32:
				if v != nil {
					position = *v
					hasPosition = true
				}
			default:
				return errno.ErrInvalidArgument.WithMessage("position field type error")
			}

			if hasPosition && position < 0 {
				return errno.ErrInvalidArgument.WithMessage("position cannot be negative")
			}
			return nil
		}
	}

	// 定义各字段的校验逻辑，通过一个 map 实现模块化和简化
	return genericvalidation.Rules{
		// 基本字段校验
		"PostID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("postID cannot be empty")
			}
			return nil
		},
		"Title": func(value any) error {
			title := value.(string)
			if title == "" {
				return errno.ErrInvalidArgument.WithMessage("title cannot be empty")
			}
			if len(strings.TrimSpace(title)) == 0 {
				return errno.ErrInvalidArgument.WithMessage("title cannot be whitespace only")
			}
			if len(title) > 200 {
				return errno.ErrInvalidArgument.WithMessage("title cannot exceed 200 characters")
			}
			return nil
		},
		"Content": func(value any) error {
			content := value.(string)
			if content == "" {
				return errno.ErrInvalidArgument.WithMessage("content cannot be empty")
			}
			if len(strings.TrimSpace(content)) == 0 {
				return errno.ErrInvalidArgument.WithMessage("content cannot be whitespace only")
			}
			return nil
		},
		"UserID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("userID cannot be empty")
			}
			return nil
		},

		// 可选字段校验
		"Cover":               validateCover(),
		"Summary":             validateSummary(),
		"CategoryID":          validateCategoryID(),
		"OriginalAuthor":      validateOriginalAuthor(),
		"OriginalSource":      validateOriginalSource(),
		"OriginalAuthorIntro": validateOriginalAuthorIntro(),
		"Position":            validatePosition(),

		// 枚举字段校验
		"PostType": validatePostType(),
		"Status":   validatePostStatus(),

		// 计数字段校验
		"ViewCount": func(value any) error {
			if value.(int32) < 0 {
				return errno.ErrInvalidArgument.WithMessage("view count cannot be negative")
			}
			return nil
		},
		"LikeCount": func(value any) error {
			if value.(int32) < 0 {
				return errno.ErrInvalidArgument.WithMessage("like count cannot be negative")
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

// ValidateCreatePostRequest 校验 CreatePostRequest 结构体的有效性
func (v *Validator) ValidateCreatePostRequest(ctx context.Context, rq *apiv1.CreatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateUpdatePostRequest 校验 UpdatePostRequest 结构体的有效性
func (v *Validator) ValidateUpdatePostRequest(ctx context.Context, rq *apiv1.UpdatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateDeletePostRequest 校验 DeletePostRequest 结构体的有效性
func (v *Validator) ValidateDeletePostRequest(ctx context.Context, rq *apiv1.DeletePostRequest) error {
	if len(rq.GetPostIDs()) == 0 {
		return errno.ErrInvalidArgument.WithMessage("postIDs cannot be empty")
	}
	for _, postID := range rq.GetPostIDs() {
		if postID == "" {
			return errno.ErrInvalidArgument.WithMessage("postID in the list cannot be empty")
		}
	}
	return nil
}

// ValidateGetPostRequest 校验 GetPostRequest 结构体的有效性
func (v *Validator) ValidateGetPostRequest(ctx context.Context, rq *apiv1.GetPostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateListPostRequest 校验 ListPostRequest 结构体的有效性
func (v *Validator) ValidateListPostRequest(ctx context.Context, rq *apiv1.ListPostRequest) error {
	// 校验分页参数
	rules := v.ValidatePostRules()
	if err := rules["Offset"](rq.GetOffset()); err != nil {
		return err
	}
	if err := rules["Limit"](rq.GetLimit()); err != nil {
		return err
	}

	// 校验可选的标题过滤参数
	if title := rq.GetTitle(); title != "" {
		if len(strings.TrimSpace(title)) == 0 {
			return errno.ErrInvalidArgument.WithMessage("title filter cannot be whitespace only")
		}
		if len(title) > 200 {
			return errno.ErrInvalidArgument.WithMessage("title filter cannot exceed 200 characters")
		}
	}

	return nil
}
