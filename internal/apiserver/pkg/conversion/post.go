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

// PostModelToPostV1 将模型层的 PostM（博客模型对象）转换为 Protobuf 层的 Post（v1 博客对象）.
func PostModelToPostV1(postModel *model.PostM) *v1.Post {
	var protoPost v1.Post
	_ = copier.CopyWithConverters(&protoPost, postModel)
	return &protoPost
}

// PostV1ToPostModel 将 Protobuf 层的 Post（v1 博客对象）转换为模型层的 PostM（博客模型对象）.
func PostV1ToPostModel(protoPost *v1.Post) *model.PostM {
	var postModel model.PostM
	_ = copier.CopyWithConverters(&postModel, protoPost)
	return &postModel
}

// PostModelToPostV1WithRelations 将文章模型和关联数据转换为 Protobuf 层的 Post 对象
func PostModelToPostV1WithRelations(postModel *model.PostM, category *model.CategoryM, tags []*model.TagM) *v1.Post {
	if postModel == nil {
		return nil
	}

	// 首先转换基础的文章信息
	var protoPost v1.Post
	_ = copier.CopyWithConverters(&protoPost, postModel)

	// 转换分类信息
	if category != nil {
		protoPost.Category = CategoryModelToCategoryV1(category)
	}

	// 转换标签信息
	if len(tags) > 0 {
		protoPost.Tags = make([]*v1.Tag, 0, len(tags))
		for _, tag := range tags {
			if tag != nil {
				protoPost.Tags = append(protoPost.Tags, TagModelToTagV1(tag))
			}
		}
	}

	return &protoPost
}
