# MiniBlog

一个小型博客系统，展示 Go 项目开发的最佳实践。

## 项目简介

MiniBlog 是一个基于 Go 语言开发的小型博客系统，旨在展示如何构建一个功能完整的 Go 项目。项目采用简洁架构设计，遵循标准化的目录结构和众多开发规范，集成了当下 Go 项目开发中常用的包和工具。

## 功能特性

- **简洁架构**：采用清晰、易维护的架构设计
- **标准目录结构**：遵循 project-layout 规范
- **认证与授权**：基于 JWT 的认证和基于 Casbin 的授权
- **日志与错误处理**：独立的日志包和错误码管理
- **丰富的 Web 功能**：请求 ID、优雅关停、中间件、跨域处理、异常恢复等
- **多服务器支持**：HTTP/HTTPS/gRPC 服务器和 HTTP 反向代理
- **API 设计**：遵循 RESTful API 规范，提供 OpenAPI/Swagger 文档
- **代码质量保证**：通过 golangci-lint 进行静态检查
- **完善的测试**：单元测试、性能测试、模糊测试和示例测试

## 技术栈

项目使用了众多 Go 生态中流行的包：

- **Web 框架**：gin
- **命令行**：cobra, pflag
- **配置管理**：viper
- **ORM**：gorm
- **认证授权**：jwt-go, casbin
- **日志**：zap
- **API 文档**：swagger
- **RPC**：grpc, protobuf, grpc-gateway
- **其他**：govalidator, uuid, pprof 等

## 快速开始

### 环境要求

- Go 1.20+
- Make

### 构建与运行

```bash
# 克隆仓库
git clone https://github.com/clin211/miniblog-v2.git
cd miniblog-v2

# 编译项目
make build

# 运行项目
./_output/mb-apiserver
```

## 开发指南

### 常用命令

```bash
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

## 开源许可

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

> 本仓库代码是「长林啊」学习「孔令飞」大佬的《云原生AI实战营》专栏的实践，如果想的深入理解本仓库项目功能设计背后的设计理念与实现逻辑，可以访问 [https://t.zsxq.com/k1cHj](https://t.zsxq.com/k1cHj) 链接！
