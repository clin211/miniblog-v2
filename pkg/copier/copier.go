// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package copier

import (
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TypeConverters 定义时间类型转换器，用于 copier 的深度拷贝。
func TypeConverters() []copier.TypeConverter {
	return []copier.TypeConverter{
		// time.Time -> *timestamppb.Timestamp
		{
			SrcType: time.Time{},
			DstType: &timestamppb.Timestamp{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(time.Time)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				return timestamppb.New(s), nil
			},
		},
		// *time.Time -> *timestamppb.Timestamp
		{
			SrcType: &time.Time{},
			DstType: &timestamppb.Timestamp{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*time.Time)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s == nil {
					return nil, nil
				}
				return timestamppb.New(*s), nil
			},
		},
		// *timestamppb.Timestamp -> time.Time
		{
			SrcType: &timestamppb.Timestamp{},
			DstType: time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*timestamppb.Timestamp)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s == nil {
					return time.Time{}, nil
				}
				return s.AsTime(), nil
			},
		},
		// *timestamppb.Timestamp -> *time.Time
		{
			SrcType: &timestamppb.Timestamp{},
			DstType: &time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*timestamppb.Timestamp)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s == nil {
					return nil, nil
				}
				t := s.AsTime()
				return &t, nil
			},
		},
		// *time.Time -> int64 (Unix 时间戳)
		{
			SrcType: &time.Time{},
			DstType: int64(0),
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(*time.Time)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s == nil {
					return int64(0), nil
				}
				return s.Unix(), nil
			},
		},
		// time.Time -> int64 (Unix 时间戳)
		{
			SrcType: time.Time{},
			DstType: int64(0),
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(time.Time)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				return s.Unix(), nil
			},
		},
		// int64 -> *time.Time (从 Unix 时间戳)
		{
			SrcType: int64(0),
			DstType: &time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(int64)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				if s == 0 {
					return nil, nil
				}
				t := time.Unix(s, 0)
				return &t, nil
			},
		},
		// int64 -> time.Time (从 Unix 时间戳)
		{
			SrcType: int64(0),
			DstType: time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(int64)
				if !ok {
					return nil, errors.New("source type not matching")
				}
				return time.Unix(s, 0), nil
			},
		},
	}
}

func CopyWithConverters(to any, from any) error {
	return copier.CopyWithOption(to, from, copier.Option{IgnoreEmpty: true, DeepCopy: true, Converters: TypeConverters()})
}

func Copy(to any, from any) error {
	return copier.Copy(to, from)
}
