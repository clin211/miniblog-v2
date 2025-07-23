// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"net"
	"net/url"

	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	genericvalidation "github.com/onexstack/onexstack/pkg/validation"

	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
)

// ValidateUserRules 定义用户相关的校验规则
func (v *Validator) ValidateUserRules() genericvalidation.Rules {
	// 通用的密码校验函数
	validatePassword := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			return isValidPassword(value.(string))
		}
	}

	// 年龄校验函数
	validateAge := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			age := value.(int32)
			if age < 0 || age > 150 {
				return errno.ErrInvalidArgument.WithMessage("age must be between 0 and 150")
			}
			return nil
		}
	}

	// 头像URL校验函数
	validateAvatar := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			avatar := value.(string)
			if avatar != "" {
				if _, err := url.ParseRequestURI(avatar); err != nil {
					return errno.ErrInvalidArgument.WithMessage("avatar must be a valid URL")
				}
			}
			return nil
		}
	}

	// 性别校验函数
	validateGender := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			gender := value.(v1.Gender)
			// 只允许已定义的枚举值
			switch gender {
			case v1.Gender_GENDER_UNSPECIFIED, v1.Gender_GENDER_MALE, v1.Gender_GENDER_FEMALE, v1.Gender_GENDER_OTHER:
				return nil
			default:
				return errno.ErrInvalidArgument.WithMessage("invalid gender value")
			}
		}
	}

	// 状态校验函数
	validateStatus := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			status := value.(int32)
			// 状态只能是0（禁用）或1（正常）
			if status != 0 && status != 1 {
				return errno.ErrInvalidArgument.WithMessage("status must be 0 (disabled) or 1 (active)")
			}
			return nil
		}
	}

	// IP地址校验函数 TODO: 需要优化，应该是在服务端自动获取IP地址，而不是在客户端获取
	validateIP := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			ipStr := value.(string)
			if ipStr != "" {
				if net.ParseIP(ipStr) == nil {
					return errno.ErrInvalidArgument.WithMessage("invalid IP address format")
				}
			}
			return nil
		}
	}

	// 注册来源校验函数
	validateRegisterSource := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			source := value.(v1.RegisterSource)
			// 只允许已定义的枚举值
			switch source {
			case v1.RegisterSource_REGISTER_SOURCE_UNSPECIFIED, v1.RegisterSource_REGISTER_SOURCE_WEB,
				v1.RegisterSource_REGISTER_SOURCE_APP, v1.RegisterSource_REGISTER_SOURCE_WECHAT,
				v1.RegisterSource_REGISTER_SOURCE_QQ, v1.RegisterSource_REGISTER_SOURCE_GITHUB,
				v1.RegisterSource_REGISTER_SOURCE_GOOGLE:
				return nil
			default:
				return errno.ErrInvalidArgument.WithMessage("invalid register source value")
			}
		}
	}

	// 微信OpenID校验函数
	validateWechatOpenID := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			openID := value.(string)
			// 微信OpenID格式校验，如果提供了OpenID则进行基本格式检查
			if openID != "" && len(openID) < 10 {
				return errno.ErrInvalidArgument.WithMessage("wechat openID format is invalid")
			}
			return nil
		}
	}

	// 定义各字段的校验逻辑，通过一个 map 实现模块化和简化
	return genericvalidation.Rules{
		// 密码相关校验
		"Password":    validatePassword(),
		"OldPassword": validatePassword(),
		"NewPassword": validatePassword(),

		// 基本信息校验
		"UserID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("userID cannot be empty")
			}
			return nil
		},
		"Username": func(value any) error {
			if !isValidUsername(value.(string)) {
				return errno.ErrUsernameInvalid
			}
			return nil
		},
		"Email": func(value any) error {
			return isValidEmail(value.(string))
		},
		"Phone": func(value any) error {
			return isValidPhone(value.(string))
		},

		// 新增字段校验
		"Age":    validateAge(),
		"Avatar": validateAvatar(),
		"Gender": validateGender(),
		"Status": validateStatus(),
		"IsRisk": func(value any) error {
			// 布尔值无需特殊校验
			return nil
		},
		"RegisterSource": validateRegisterSource(),
		"RegisterIP":     validateIP(),
		"WechatOpenID":   validateWechatOpenID(),

		// 分页参数校验
		"Limit": func(value any) error {
			// 允许 limit 为 0（使用默认值）或正数，只有负数时才报错
			if value.(int64) < 0 {
				return errno.ErrInvalidArgument.WithMessage("limit cannot be negative")
			}
			return nil
		},
		"Offset": func(value any) error {
			// offset 可以为 0 或正数，只有负数时才报错
			if value.(int64) < 0 {
				return errno.ErrInvalidArgument.WithMessage("offset cannot be negative")
			}
			return nil
		},
	}
}

// ValidateLogin 校验修改密码请求.
func (v *Validator) ValidateLoginRequest(ctx context.Context, rq *v1.LoginRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateChangePasswordRequest 校验 ChangePasswordRequest 结构体的有效性.
func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, rq *v1.ChangePasswordRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateCreateUserRequest 校验 CreateUserRequest 结构体的有效性.
func (v *Validator) ValidateCreateUserRequest(ctx context.Context, rq *v1.CreateUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateUpdateUserRequest 校验更新用户请求.
func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *v1.UpdateUserRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateSelectedFields(rq, v.ValidateUserRules(), "UserID")
}

// ValidateDeleteUserRequest 校验 DeleteUserRequest 结构体的有效性.
func (v *Validator) ValidateDeleteUserRequest(ctx context.Context, rq *v1.DeleteUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateGetUserRequest 校验 GetUserRequest 结构体的有效性.
func (v *Validator) ValidateGetUserRequest(ctx context.Context, rq *v1.GetUserRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateListUserRequest 校验 ListUserRequest 结构体的有效性.
func (v *Validator) ValidateListUserRequest(ctx context.Context, rq *v1.ListUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}
