// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package device

import (
	"context"

	"github.com/clin211/miniblog-v2/internal/apiserver/store"
)

// DeviceBiz 定义了 device 模块在 biz 层所实现的方法.
type DeviceBiz interface {
	Create(ctx context.Context, device *store.DeviceM) error
	Update(ctx context.Context, device *store.DeviceM) error
	Delete(ctx context.Context, deviceID string) error
	Get(ctx context.Context, deviceID string) (*store.DeviceM, error)
	List(ctx context.Context, limit, offset int) ([]*store.DeviceM, int64, error)
}

// deviceBiz 是 DeviceBiz 接口的实现.
type deviceBiz struct {
	store store.IStore
}

// 确保 deviceBiz 实现了 DeviceBiz 接口.
var _ DeviceBiz = (*deviceBiz)(nil)

// NewDeviceBiz 创建一个 DeviceBiz 的实例.
func NewDeviceBiz(store store.IStore) *deviceBiz {
	return &deviceBiz{store: store}
}

// Create 创建一个新的设备记录.
func (b *deviceBiz) Create(ctx context.Context, device *store.DeviceM) error {
	return b.store.Device().Create(ctx, device)
}

// Update 更新一个设备记录.
func (b *deviceBiz) Update(ctx context.Context, device *store.DeviceM) error {
	return b.store.Device().Update(ctx, device)
}

// Delete 删除一个设备记录.
func (b *deviceBiz) Delete(ctx context.Context, deviceID string) error {
	return b.store.Device().Delete(ctx, deviceID)
}

// Get 根据设备ID获取设备记录.
func (b *deviceBiz) Get(ctx context.Context, deviceID string) (*store.DeviceM, error) {
	return b.store.Device().Get(ctx, deviceID)
}

// List 获取设备列表.
func (b *deviceBiz) List(ctx context.Context, limit, offset int) ([]*store.DeviceM, int64, error) {
	return b.store.Device().List(ctx, limit, offset)
}
