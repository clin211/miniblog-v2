// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package conversion

import (
	"github.com/clin211/miniblog-v2/pkg/copier"

	"github.com/clin211/miniblog-v2/internal/apiserver/model"
	appv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1/app"
)

// TagModelToTagV1 将模型层的 TagM（标签模型对象）转换为 Protobuf 层的 Tag（v1 标签对象）.
func TagModelToTagV1(tagModel *model.TagM) *appv1.Tag {
	var protoTag appv1.Tag
	_ = copier.CopyWithConverters(&protoTag, tagModel)
	return &protoTag
}

// TagV1ToTagModel 将 Protobuf 层的 Tag（v1 标签对象）转换为模型层的 TagM（标签模型对象）.
func TagV1ToTagModel(protoTag *appv1.Tag) *model.TagM {
	var tagModel model.TagM
	_ = copier.CopyWithConverters(&tagModel, protoTag)
	return &tagModel
}
