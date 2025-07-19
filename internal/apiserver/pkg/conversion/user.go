// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package conversion

import (
	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	"github.com/clin211/miniblog-v2/pkg/copier"

	appv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1/app"
)

// UserModelToUserV1 将模型层的 UserM（用户模型对象）转换为 Protobuf 层的 User（v1 用户对象）.
func UserModelToUserV1(userModel *model.UserM) *appv1.User {
	if userModel == nil {
		return nil
	}

	var protoUser appv1.User
	_ = copier.CopyWithConverters(&protoUser, userModel)
	return &protoUser
}

// UserV1ToUserModel 将 Protobuf 层的 User（v1 用户对象）转换为模型层的 UserM（用户模型对象）.
func UserV1ToUserModel(protoUser *appv1.User) *model.UserM {
	if protoUser == nil {
		return nil
	}

	var userModel model.UserM
	_ = copier.CopyWithConverters(&userModel, protoUser)
	return &userModel
}
