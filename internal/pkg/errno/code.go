// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package errno

import (
	"net/http"
)

// errorsx 预定义标准的错误.
var (
	// OK 代表请求成功.
	OK = &ErrorX{Code: http.StatusOK, Message: ""}

	// OKWithMsg 代表请求成功.
	OKWithMsg = &ErrorX{Code: http.StatusOK, Message: "success"} //nolint:errname

	// ErrInternal 表示所有未知的服务器端错误.
	ErrInternal = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError", Message: "Internal server error."}

	// ErrNotFound 表示资源未找到.
	ErrNotFound = &ErrorX{Code: http.StatusNotFound, Reason: "NotFound", Message: "Resource not found."}

	// ErrBind 表示请求体绑定错误.
	ErrBind = &ErrorX{Code: http.StatusBadRequest, Reason: "BindError", Message: "Error occurred while binding the request body to the struct."}

	// ErrInvalidArgument 表示参数验证失败.
	ErrInvalidArgument = &ErrorX{Code: http.StatusBadRequest, Reason: "InvalidArgument", Message: "Argument verification failed."}

	// ErrUnauthenticated 表示认证失败.
	ErrUnauthenticated = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated", Message: "Unauthenticated."}

	// ErrPermissionDenied 表示请求没有权限.
	ErrPermissionDenied = &ErrorX{Code: http.StatusForbidden, Reason: "PermissionDenied", Message: "Permission denied. Access to the requested resource is forbidden."}

	// ErrOperationFailed 表示操作失败.
	ErrOperationFailed = &ErrorX{Code: http.StatusConflict, Reason: "OperationFailed", Message: "The requested operation has failed. Please try again later."}

	// ErrPageNotFound 表示页面未找到.
	ErrPageNotFound = &ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PageNotFound", Message: "Page not found."}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.SignToken", Message: "Error occurred while signing the JSON web token."}

	// ErrTokenInvalid 表示 JWT Token 格式无效.
	ErrTokenInvalid = &ErrorX{Code: http.StatusUnauthorized, Reason: "Unauthenticated.TokenInvalid", Message: "Token was invalid."}

	// ErrDBRead 表示数据库读取失败.
	ErrDBRead = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.DBRead", Message: "Database read failure."}

	// ErrDBWrite 表示数据库写入失败.
	ErrDBWrite = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.DBWrite", Message: "Database write failure."}

	// ErrAddRole 表示在添加角色时发生错误.
	ErrAddRole = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.AddRole", Message: "Error occurred while adding the role."}

	// ErrRemoveRole 表示在删除角色时发生错误.
	ErrRemoveRole = &ErrorX{Code: http.StatusInternalServerError, Reason: "InternalError.RemoveRole", Message: "Error occurred while removing the role."}
)
