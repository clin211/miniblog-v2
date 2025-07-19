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

// PostModelToPostV1 将模型层的 PostM（博客模型对象）转换为 Protobuf 层的 Post（v1 博客对象）.
func PostModelToPostV1(postModel *model.PostM) *appv1.Post {
	var protoPost appv1.Post
	_ = copier.CopyWithConverters(&protoPost, postModel)
	return &protoPost
}

// PostV1ToPostModel 将 Protobuf 层的 Post（v1 博客对象）转换为模型层的 PostM（博客模型对象）.
func PostV1ToPostModel(protoPost *appv1.Post) *model.PostM {
	var postModel model.PostM
	_ = copier.CopyWithConverters(&postModel, protoPost)
	return &postModel
}
