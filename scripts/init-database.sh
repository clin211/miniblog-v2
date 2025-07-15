#!/bin/bash

# Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/clin211/miniblog-v2.git.

# 数据库初始化脚本
# 用于自动化 MySQL 数据库的初始化过程

set -e

echo "🚀 MiniBlog 数据库初始化开始..."

# 检查是否在项目根目录
if [ ! -f "configs/miniblog.sql" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    echo "当前目录：$(pwd)"
    exit 1
fi

# 检查 Docker 是否运行
if ! docker ps >/dev/null 2>&1; then
    echo "❌ 错误：Docker 服务未运行，请先启动 Docker"
    exit 1
fi

# 检查 MySQL 容器是否存在
echo "🔍 检查 MySQL 容器状态..."
if ! docker ps | grep -q "miniblog-mysql"; then
    echo "❌ 错误：MySQL 容器未运行"
    echo "请先启动数据库服务："
    echo "  cd deployment/docker/dev"
    echo "  docker compose -f docker-compose.env.yml up -d"
    exit 1
fi

# 等待 MySQL 完全启动
echo "⏳ 等待 MySQL 服务完全启动..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if docker exec miniblog-mysql mysql -u root -p123456 -e "SELECT 1" >/dev/null 2>&1; then
        echo "✅ MySQL 服务已准备就绪"
        break
    fi

    attempt=$((attempt + 1))
    echo "等待中... ($attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ 错误：MySQL 服务启动超时"
    echo "请检查 MySQL 容器日志："
    echo "  docker logs miniblog-mysql"
    exit 1
fi

# 检查数据库是否已存在表
echo "🔍 检查数据库状态..."
table_count=$(docker exec miniblog-mysql mysql -u root -p123456 -e "USE miniblog_v2; SHOW TABLES;" 2>/dev/null | wc -l)

if [ "$table_count" -gt 1 ]; then
    echo "⚠️  数据库表已存在，是否重新初始化？(y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "✅ 跳过数据库初始化"
        exit 0
    fi
fi

# 导入数据库表结构
echo "📊 导入数据库表结构..."
if docker exec -i miniblog-mysql mysql -u root -p123456 miniblog_v2 <configs/miniblog.sql; then
    echo "✅ 数据库表结构导入成功"
else
    echo "❌ 错误：数据库表结构导入失败"
    exit 1
fi

# 验证数据库表
echo "🔍 验证数据库表..."
echo "已创建的表："
docker exec miniblog-mysql mysql -u root -p123456 -e "USE miniblog_v2; SHOW TABLES;" 2>/dev/null | grep -v "Tables_in_miniblog_v2" | sed 's/^/  - /'

echo ""
echo "🎉 数据库初始化完成！"
echo ""
echo "📋 下一步："
echo "1. 启动应用程序："
echo "   go run cmd/mb-apiserver/main.go --config configs/mb-apiserver.yaml"
echo ""
echo "2. 或者使用 air 进行开发："
echo "   air"
echo ""
echo "3. 验证服务："
echo "   curl http://localhost:5555/healthz"
