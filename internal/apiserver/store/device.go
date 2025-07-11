// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

// DeviceStore 定义了 device 模块在 store 层所实现的方法，演示 MongoDB 的使用.
type DeviceStore interface {
	Create(ctx context.Context, device *DeviceM) error
	Update(ctx context.Context, device *DeviceM) error
	Delete(ctx context.Context, deviceID string) error
	Get(ctx context.Context, deviceID string) (*DeviceM, error)
	List(ctx context.Context, limit, offset int) ([]*DeviceM, int64, error)
}

// DeviceM 定义设备模型，用于演示 MongoDB 的使用.
// 完全扁平化的结构，可以存储任意数据，无需登录
type DeviceM struct {
	ID        string                 `bson:"_id,omitempty" json:"id,omitempty"`
	Data      map[string]interface{} `bson:"data" json:"data"` // 扁平化存储，直接展开到根级别
	CreatedAt int64                  `bson:"created_at" json:"created_at"`
	UpdatedAt int64                  `bson:"updated_at" json:"updated_at"`
}

// deviceStore 是 DeviceStore 接口的实现，演示 MongoDB 的使用.
type deviceStore struct {
	store *datastore
}

// 确保 deviceStore 实现了 DeviceStore 接口.
var _ DeviceStore = (*deviceStore)(nil)

// newDeviceStore 创建 deviceStore 的实例.
func newDeviceStore(store *datastore) *deviceStore {
	return &deviceStore{store: store}
}

// getCollection 获取设备集合
func (s *deviceStore) getCollection() *mongo.Collection {
	return s.store.mongo.Database("miniblog_v2").Collection("devices")
}

// Create 创建新设备记录
func (s *deviceStore) Create(ctx context.Context, device *DeviceM) error {
	collection := s.getCollection()

	// 如果没有提供ID，生成一个新的ObjectID
	if device.ID == "" {
		device.ID = primitive.NewObjectID().Hex()
	}

	// 设置时间戳
	now := time.Now().Unix()
	device.CreatedAt = now
	device.UpdatedAt = now

	// 确保Data字段不为nil
	if device.Data == nil {
		device.Data = make(map[string]interface{})
	}

	_, err := collection.InsertOne(ctx, device)
	if err != nil {
		log.Errorw("Failed to insert device into MongoDB", "err", err, "device_id", device.ID)
		return err
	}

	log.Infow("Device created successfully", "device_id", device.ID)

	return nil
}

// Update 更新设备记录
func (s *deviceStore) Update(ctx context.Context, device *DeviceM) error {
	collection := s.getCollection()

	// 更新时间戳
	device.UpdatedAt = time.Now().Unix()

	// 构建更新文档，包含所有数据字段
	updateDoc := bson.M{"updated_at": device.UpdatedAt}
	for key, value := range device.Data {
		updateDoc[key] = value
	}

	filter := bson.M{"_id": device.ID}
	update := bson.M{"$set": updateDoc}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorw("Failed to update device in MongoDB", "err", err, "device_id", device.ID)
		return err
	}

	if result.MatchedCount == 0 {
		log.Warnw("Device not found for update", "device_id", device.ID)
		return errors.New("device not found")
	}

	log.Infow("Device updated successfully", "device_id", device.ID, "matched", result.MatchedCount)
	return nil
}

// Delete 删除设备记录
func (s *deviceStore) Delete(ctx context.Context, deviceID string) error {
	collection := s.getCollection()

	filter := bson.M{"_id": deviceID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Errorw("Failed to delete device from MongoDB", "err", err, "device_id", deviceID)
		return err
	}

	if result.DeletedCount == 0 {
		log.Warnw("Device not found for deletion", "device_id", deviceID)
		return errors.New("device not found")
	}

	log.Infow("Device deleted successfully", "device_id", deviceID, "deleted_count", result.DeletedCount)
	return nil
}

// Get 根据设备ID获取设备记录
func (s *deviceStore) Get(ctx context.Context, deviceID string) (*DeviceM, error) {
	collection := s.getCollection()

	filter := bson.M{"_id": deviceID}
	var device DeviceM

	err := collection.FindOne(ctx, filter).Decode(&device)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Warnw("Device not found", "device_id", deviceID)
			return nil, errors.New("device not found")
		}
		log.Errorw("Failed to get device from MongoDB", "err", err, "device_id", deviceID)
		return nil, err
	}

	log.Infow("Device retrieved successfully", "device_id", deviceID)
	return &device, nil
}

// List 获取设备列表（无需用户验证）
func (s *deviceStore) List(ctx context.Context, limit, offset int) ([]*DeviceM, int64, error) {
	collection := s.getCollection()

	// 获取所有设备，无用户限制
	filter := bson.M{}

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Errorw("Failed to count devices in MongoDB", "err", err)
		return nil, 0, err
	}

	// 构建查询选项
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // 按创建时间降序排列
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}

	// 执行查询
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Errorw("Failed to list devices from MongoDB", "err", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// 解析结果
	var devices []*DeviceM
	for cursor.Next(ctx) {
		var device DeviceM
		if err := cursor.Decode(&device); err != nil {
			log.Errorw("Failed to decode device from MongoDB", "err", err)
			continue
		}
		devices = append(devices, &device)
	}

	if err := cursor.Err(); err != nil {
		log.Errorw("Cursor error while listing devices", "err", err)
		return nil, 0, err
	}

	log.Infow("Devices listed successfully", "count", len(devices), "total", total)
	return devices, total, nil
}
