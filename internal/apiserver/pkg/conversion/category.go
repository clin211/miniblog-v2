// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package conversion

import (
	"github.com/clin211/miniblog-v2/pkg/copier"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	v1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
)

// CategoryModelToCategoryV1 将模型层的 CategoryM 转换为 Protobuf 层的 Category
func CategoryModelToCategoryV1(categoryModel *model.CategoryM) *v1.Category {
	var protoCategory v1.Category
	_ = copier.CopyWithConverters(&protoCategory, categoryModel)

	// 手动处理 IsActive 字段的转换
	if categoryModel.IsActive != nil {
		if *categoryModel.IsActive == 1 {
			protoCategory.IsActive = v1.IsActive_IS_ACTIVE_ACTIVE
		} else {
			protoCategory.IsActive = v1.IsActive_IS_ACTIVE_DISABLED
		}
	} else {
		// 如果数据库中为 null，默认设置为激活状态
		protoCategory.IsActive = v1.IsActive_IS_ACTIVE_ACTIVE
	}

	return &protoCategory
}

// CategoryV1ToCategoryModel 将 Protobuf 层的 Category 转换为模型层的 CategoryM
func CategoryV1ToCategoryModel(protoCategory *v1.Category) *model.CategoryM {
	var categoryModel model.CategoryM
	_ = copier.CopyWithConverters(&categoryModel, protoCategory)

	// 手动处理 IsActive 字段的转换
	var isActiveValue int32
	if protoCategory.IsActive == v1.IsActive_IS_ACTIVE_ACTIVE {
		isActiveValue = 1
	} else {
		isActiveValue = 0
	}
	categoryModel.IsActive = &isActiveValue

	return &categoryModel
}
