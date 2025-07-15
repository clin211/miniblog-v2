# MiniBlog

一个小型博客系统，展示 Go 项目开发的最佳实践。

## 项目简介

MiniBlog 是一个基于 Go 语言开发的小型博客系统，旨在展示如何构建一个功能完整的 Go 项目。项目采用简洁架构设计，遵循标准化的目录结构和众多开发规范，集成了当下 Go 项目开发中常用的包和工具。

## 功能特性

- **简洁架构**：采用清晰、易维护的架构设计
- **标准目录结构**：遵循 project-layout 规范
- **认证与授权**：基于 JWT 的认证和基于 Casbin 的授权
- **多数据库支持**：集成 MySQL（主数据库）、MongoDB（文档存储）双数据库架构
- **高性能缓存**：基于 Redis 的分布式缓存系统，支持数据缓存和会话管理
- **日志与错误处理**：独立的日志包和错误码管理
- **丰富的 Web 功能**：请求 ID、优雅关停、中间件、跨域处理、异常恢复等
- **多服务器支持**：HTTP/HTTPS/gRPC 服务器和 HTTP 反向代理
- **API 设计**：遵循 RESTful API 规范，提供 OpenAPI/Swagger 文档
- **代码质量保证**：通过 golangci-lint 进行静态检查

## 技术栈

项目使用了众多 Go 生态中流行的包：

- **Web 框架**：gin
- **命令行**：cobra, pflag
- **配置管理**：viper
- **数据库**：MySQL, MongoDB
- **ORM**：gorm
- **缓存**：Redis
- **认证授权**：jwt-go, casbin
- **日志**：zap
- **API 文档**：swagger
- **RPC**：grpc, protobuf, grpc-gateway
- **容器化**：Docker, Docker Compose
- **其他**：govalidator, uuid, pprof 等

## 快速开始

### 环境要求

- Go 1.20+
- Docker & Docker Compose
- Make

### 启动步骤

#### 1. 克隆仓库

```sh
git clone https://github.com/clin211/miniblog-v2.git
cd miniblog-v2
```

#### 2. 启动数据库服务

使用 Docker Compose 启动 MySQL、MongoDB 和 Redis 服务：

```sh
# 进入 Docker 配置目录
cd deployment/docker/dev

# 启动所有数据库服务
docker compose -f docker-compose.env.yml up -d

# 检查服务状态
docker compose -f docker-compose.env.yml ps
```

#### 3. 初始化数据库

启动数据库服务后，需要初始化 MySQL 数据库表结构。提供两种方式：

- **方式一：使用自动化脚本（推荐）**

  ```sh
  # 回到项目根目录
  cd ../../../

  # 运行自动化初始化脚本
  ./scripts/init-database.sh
  ```

- **方式二：手动执行**

  ```sh
  # 回到项目根目录
  cd ../../../

  # 等待 MySQL 完全启动（重要！）
  echo "等待 MySQL 启动完成..."
  sleep 10

  # 导入数据库表结构
  docker exec -i miniblog-mysql mysql -u root -p123456 miniblog_v2 < configs/miniblog.sql

  # 验证数据库表是否创建成功
  docker exec -it miniblog-mysql mysql -u root -p123456 -e "USE miniblog_v2; SHOW TABLES;"
  ```

**故障排除**：如果遇到连接错误，请：

1. 确认容器正在运行：`docker compose -f deployment/docker/dev/docker-compose.env.yml ps`
2. 检查 MySQL 日志：`docker logs miniblog-mysql --tail 20`
3. 等待更长时间再重试（MySQL 首次启动可能需要 1-2 分钟）

#### 4. 启动应用程序

```sh
# 方式一：直接运行（推荐用于开发）
air

# 方式二：编译后运行
make build
./_output/mb-apiserver --config configs/mb-apiserver.yaml
```

#### 5. 验证服务状态

- 检查应用程序健康状态

  可以在 `mb-apiserver.yaml` 文件中修改 server-mode 的模式，一共支持三种模式：`grpc`、`grpc-gateway`、`gin`！

  ```sh
  # http 或者 grpc-gateway 模式
  curl http://localhost:5555/healthz

  # grpc 模式
  go run ./examples/client/health/main.go
  ```

- 检查 MySQL 连接（可选）

  ```sh
  docker exec -it miniblog-mysql mysql -u root -p123456 miniblog_v2 -e "SELECT COUNT(*) FROM user;"
  ```

- 检查 MongoDB 连接（可选）

  ```sh
  docker exec -it miniblog-mongo mongosh -u root -p r8SggC783Xh1 --authenticationDatabase admin miniblog_v2 --eval "db.getName()"
  ```

- 检查 Redis 缓存（可选）

  ```sh
  ./scripts/redis-inspect.sh
  ```

### 服务端口说明

- **HTTP 服务器**: `http://localhost:5555`
- **gRPC 服务器**: `localhost:6666`
- **MySQL 数据库**: `localhost:11006`
- **MongoDB 数据库**: `localhost:27001`
- **Redis 缓存**: `localhost:63790`

### 停止服务

```sh
# 停止应用程序
# 使用 Ctrl+C 停止正在运行的应用程序

# 停止数据库、缓存
cd deployment/docker/dev
docker compose -f docker-compose.env.yml down
```

## 开发指南

### 开发流程

本项目的开发流程遵循分层架构设计，以下是新增功能（以 Comment 资源为例）的完整开发流程：

1. **定义 API 接口**
   - 设计 RESTful API 接口规范
   - 定义 Proto 文件中的接口结构

2. **编译 Protobuf 文件**
   - 使用 protoc 编译 .proto 文件
   - 生成 Go 代码和相关接口定义

3. **数据库设计与 Model 生成**
   - 在数据库中创建对应的数据表（如 comment 表）
   - 修改 `cmd/gen-gorm-model/gen_gorm_model.go` 文件，添加新表的 GORM Model 生成代码
   - 运行 `go run cmd/gen-gorm-model/gen_gorm_model.go` 命令生成 GORM Model

4. **请求参数处理**
   - 完善 API 接口请求参数的默认值设置方法
   - 修改 `pkg/api/apiserver/v1/*.pb.defaults.go` 文件

5. **参数校验**
   - 实现 API 接口的请求参数校验方法
   - 在 `internal/apiserver/pkg/validation/` 目录中实现对应的校验逻辑

6. **Store 层实现**
   - 实现资源的 Store 层代码（数据访问层）
   - 在 `internal/apiserver/store/` 目录中实现数据库操作逻辑

7. **数据转换**
   - 实现资源的 Model 和 Proto 的转换函数
   - 在 `internal/apiserver/pkg/conversion/` 目录中实现转换逻辑

8. **Biz 层实现**
   - 实现资源的 Biz 层代码（业务逻辑层）
   - 在 `internal/apiserver/biz/v1/` 目录中实现业务逻辑

9. **Handler 层实现**
   - 实现资源的 Handler 层代码（控制器层）
   - 在 `internal/apiserver/handler/` 目录中实现 HTTP 请求处理

#### 架构层次

本项目采用清晰的分层架构：

```txt
Handler 层 -> Biz 层 -> Store 层 -> 数据库
     ↑           ↑         ↑
   HTTP请求    业务逻辑   数据访问
```

每一层都有明确的职责分工，确保代码的可维护性和可扩展性。

### 常用命令

```sh
# 执行下面所有伪目标（因为设置了 .DEFAULT_GOAL）
make

# 仅构建
make build

# 格式化代码
make format

# 添加版权头信息
make add-copyright

# 更新依赖
make tidy

# 清理构建产物
make clean
```

### 生成 GORM Model

因为编写了 gen-gorm-model 工具用来快速生成 Model 文件，所以，在每次数据库有字段增删改的时候，都可以运行 gen-gorm-model 来生成 Model 文件。运行以下命令，可以查看 gen-gorm-model 的使用方式：

```sh
$ go run cmd/gen-gorm-model/gen_gorm_model.go -h
Usage: main [flags] arg [arg...]

This is a pflag example.

Flags:
  -a, --addr string             MySQL host address. (default "127.0.0.1:3306")
      --component strings       Generated model code's for specified component. (default [mb])
  -d, --db string               Database name to connect to. (default "miniblog")
  -h, --help                    Show this help message.
      --model-pkg-path string   Generated model code's package name.
  -p, --password string         Password to use when connecting to the database. (default "miniblog1234")
  -u, --username string         Username to connect to the database. (default "miniblog")
```

使用方式：

```sh
cd cmd/gen-gorm-model
go run ./gen_gorm_model.go -a <host>:<port> -u <username> -p <password> -d <database_name>
```

> 以 docker-compose 中 MySQL 配置为例！生成 GORM Model 的命令如下：
>
> ```sh
> $ cd cmd/gen-gorm-model
> $ go run ./gen_gorm_model.go -a 127.0.0.1:3306 -u 'root' -p 'root' -d 'miniblog_v2'
> 2025/06/13 09:01:50 got 24 columns from table <user>
> 2025/06/13 09:01:50 got 20 columns from table <post>
> 2025/06/13 09:01:50 got 9 columns from table <category>
> 2025/06/13 09:01:50 got 6 columns from table <tag>
> 2025/06/13 09:01:50 got 5 columns from table <post_tag>
> 2025/06/13 09:01:50 got 8 columns from table <casbin_rule>
> 2025/06/13 09:01:50 Start generating code.
> 2025/06/13 09:01:50 generate model file(table <category> -> {model.CategoryM}): /Users/forest/code/backend/Go/miniblog-v2/internal/apiserver/model/category.gen.go
> 2025/06/13 09:01:50 generate model file(table <tag> -> {model.TagM}): /Users/forest/code/backend/Go/miniblog-v2/internal/apiserver/model/tag.gen.go
> 2025/06/13 09:01:50 generate model file(table <post_tag> -> {model.PostTagM}): /Users/forest/code/backend/Go/miniblog-v2/internal/apiserver/model/post_tag.gen.go
> 2025/06/13 09:01:50 generate model file(table <casbin_rule> -> {model.CasbinRuleM}): /Users/forest/code/backend/Go/miniblog-v2/internal/apiserver/model/casbin_rule.gen.go
> 2025/06/13 09:01:50 generate model file(table <post> -> {model.PostM}): /Users/forest/code/backend/Go/miniblog-v2/internal/apiserver/model/post.gen.go
> 2025/06/13 09:01:50 generate model file(table <user> -> {model.UserM}): /Users/forest/code/backend/Go/miniblog-v2/internal/apiserver/model/user.gen.go
> 2025/06/13 09:01:50 Generate code done.
> ```
>

## 开源许可

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

> 本仓库代码是「长林啊」学习「孔令飞」大佬的《云原生AI实战营》专栏的实践，如果想的深入理解本仓库项目功能设计背后的设计理念与实现逻辑，可以访问 [https://t.zsxq.com/k1cHj](https://t.zsxq.com/k1cHj) 链接！
