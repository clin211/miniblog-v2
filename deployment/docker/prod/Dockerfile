# Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/clin211/miniblog-v2.git.

# 第一阶段：构建阶段
FROM golang:1.24.0-alpine AS builder

ARG PROJECT_VERSION

# 安装必要的构建工具
RUN apk add --no-cache git make bash

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
  GOPROXY=https://goproxy.cn,direct \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# 工作目录路径
WORKDIR /app

# 从项目复制go.mod 和 go.sum 文件到工作目录
COPY go.* ./

# 从 go modules 官方镜像获取依赖
RUN go mod download && go mod verify

# 将我们项目中的所有内容复制到工作目录中
COPY . .

# 在工作目录中执行 go build 命令生成二进制文件
RUN make build

# 第二阶段：运行阶段
FROM alpine:latest

# 安装ca证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/_output/mb-apiserver .

# 创建配置文件目录和日志目录
RUN mkdir -p /etc/miniblog /root/_output

# 告诉 docker 我们的代码将暴露端口 5555 和 6666
EXPOSE 5555
EXPOSE 6666

# 当我们运行这个镜像容器时，它将从我们的构建可执行文件开始执行
ENTRYPOINT ["./mb-apiserver", "-c", "/etc/miniblog/mb-apiserver.yaml"]
