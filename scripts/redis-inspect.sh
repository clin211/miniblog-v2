#!/bin/bash

# Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/clin211/miniblog-v2.git.

# Redis 数据查看脚本
# 用于快速查看 Redis 中的缓存数据

echo "🔍 Redis 数据查看工具"
echo "====================="

# Redis 连接信息
REDIS_CONTAINER="miniblog-redis"
REDIS_PASSWORD="v6yJM2kAZpOc"

echo "📊 Redis 基本信息："
docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning INFO server | grep -E "(redis_version|uptime_in_seconds|connected_clients)"

echo ""
echo "🔢 数据库统计："
docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning DBSIZE

echo ""
echo "🔑 所有键列表："
docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning KEYS '*'

echo ""
echo "💾 内存使用情况："
docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning INFO memory | grep -E "(used_memory_human|used_memory_peak_human)"

echo ""
echo "⏰ 键的过期时间（如果有）："
KEYS=$(docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning KEYS '*' | tr -d '\r')
for key in $KEYS; do
  if [ ! -z "$key" ]; then
    ttl=$(docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning TTL "$key" | tr -d '\r')
    if [ "$ttl" != "-1" ] && [ "$ttl" != "-2" ]; then
      echo "  $key: ${ttl}s"
    fi
  fi
done

echo ""
echo "🎯 详细数据查看："
echo "使用以下命令查看特定键的值："
echo "docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning GET <key_name>"
echo "docker exec -it ${REDIS_CONTAINER} redis-cli -a ${REDIS_PASSWORD} --no-auth-warning HGETALL <key_name>"
