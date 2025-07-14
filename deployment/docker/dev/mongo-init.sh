#!/bin/bash

# MongoDB 初始化脚本
# 在MongoDB容器启动时自动执行，创建miniblog_v2数据库的应用用户
# 此脚本会被放置在/docker-entrypoint-initdb.d/目录中

echo "🚀 Starting MongoDB user initialization for miniblog_v2..."

# 使用mongosh创建应用用户
# 注意：此时MongoDB已经启动，root用户已存在，可以直接连接
mongosh --username "$MONGO_INITDB_ROOT_USERNAME" \
  --password "$MONGO_INITDB_ROOT_PASSWORD" \
  --authenticationDatabase admin \
  --eval "
        // 切换到 miniblog_v2 数据库
        print('📝 Switching to miniblog_v2 database...');
        use miniblog_v2;

        // 创建应用用户
        print('👤 Creating application user...');
        db.createUser({
            user: 'root',
            pwd: 'r8SggC783Xh1',
            roles: [
                {
                    role: 'readWrite',
                    db: 'miniblog_v2'
                },
                {
                    role: 'dbAdmin',
                    db: 'miniblog_v2'
                }
            ]
        });

        print('✅ User \"root\" created successfully for database \"miniblog_v2\"');

        // 验证用户创建
        const users = db.getUsers();
        print('📋 Total users in miniblog_v2 database: ' + users.length);

        // 创建一个测试集合以确保数据库被实际创建
        db.test_collection.insertOne({message: 'Database initialized', timestamp: new Date()});
        print('🗄️  Test collection created in miniblog_v2 database');
    "

if [ $? -eq 0 ]; then
  echo "🎉 MongoDB initialization completed successfully!"
  echo "🔧 Database 'miniblog_v2' is ready with user 'root'"
else
  echo "❌ MongoDB initialization failed!"
  exit 1
fi
