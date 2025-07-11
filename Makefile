SHELL := /bin/bash

# ==============================================================================
# 定义全局 Makefile 变量方便后面引用

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# Protobuf 文件存放路径
APIROOT=$(PROJ_ROOT_DIR)/pkg/api
# 项目根目录
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output

# ==============================================================================
# 定义版本相关变量

## 指定应用使用的 version 包，会通过 `-ldflags -X` 向该包中指定的变量注入值
VERSION_PACKAGE=github.com/clin211/miniblog-v2/pkg/version
## 定义 VERSION 语义化版本号
ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git describe --tags --always --match='v*')
endif

## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
  	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
	-X $(VERSION_PACKAGE).gitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# ==============================================================================
# 定义默认目标为 all
.DEFAULT_GOAL := all

# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
.PHONY: all
all: tidy format lint build add-copyright

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/mb-apiserver $(PROJ_ROOT_DIR)/cmd/mb-apiserver/main.go

.PHONY: format
format: # 格式化 Go 源码.
	@gofmt -s -w ./

.PHONY: add-copyright
add-copyright: # 添加版权头信息.
	@addlicense -v -f $(PROJ_ROOT_DIR)/scripts/boilerplate.txt $(PROJ_ROOT_DIR) --skip-dirs=third_party,vendor,$(OUTPUT_DIR) --skip-files=\.pb\.go$

.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: protoc
protoc: # 编译 protobuf 文件.
	@echo "===========> Generate protobuf files"
	@echo "APIROOT: $(APIROOT)"
	@echo "Found proto files:"
	@find $(APIROOT) -name "*.proto"
	@protoc \
        --proto_path=$(APIROOT) \
        --proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf \
        --go_out=paths=source_relative:$(APIROOT) \
        --go-grpc_out=paths=source_relative:$(APIROOT) \
        --grpc-gateway_out=allow_delete_body=true,paths=source_relative:$(APIROOT) \
		--openapiv2_out=$(PROJ_ROOT_DIR)/api/openapi \
		--openapiv2_opt=allow_delete_body=true,logtostderr=true \
        $(shell find $(APIROOT) -name "*.proto")
	@find $(PROJ_ROOT_DIR) -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;

.PHONY: ca
ca: # 生成 CA 文件.
	@mkdir -p $(OUTPUT_DIR)/cert
	@# 1. 生成根证书私钥 (CA Key)
	@openssl genrsa -out $(OUTPUT_DIR)/cert/ca.key 4096
	@# 2. 使用根私钥生成证书签名请求 (CA CSR)，有效期为 10 年
	@openssl req -new -nodes -key $(OUTPUT_DIR)/cert/ca.key -sha256 -days 3650 -out $(OUTPUT_DIR)/cert/ca.csr \
        -subj "/C=CN/ST=sichuan/L=chengdu/O=clin211/OU=it/CN=127.0.0.1/emailAddress=767425412lin@gmail.com"
	@# 3. 使用根私钥签发根证书 (CA CRT)，使其自签名
	@openssl x509 -req -days 365 -in $(OUTPUT_DIR)/cert/ca.csr -signkey $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.crt
	@# 4. 生成服务端私钥
	@openssl genrsa -out $(OUTPUT_DIR)/cert/server.key 2048
	@# 5. 使用服务端私钥生成服务端的证书签名请求 (Server CSR)
	@openssl req -new -key $(OUTPUT_DIR)/cert/server.key -out $(OUTPUT_DIR)/cert/server.csr \
        -subj "/C=CN/ST=sichuan/L=chengdu/O=serverdevops/OU=serverit/CN=localhost/emailAddress=767425412lin@gmail.com" \
        -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
	@# 6. 使用根证书 (CA) 签发服务端证书 (Server CRT)
	@openssl x509 -days 365 -sha256 -req -CA $(OUTPUT_DIR)/cert/ca.crt -CAkey $(OUTPUT_DIR)/cert/ca.key \
        -CAcreateserial -in $(OUTPUT_DIR)/cert/server.csr -out $(OUTPUT_DIR)/cert/server.crt -extensions v3_req \
        -extfile <(printf "[v3_req]\nsubjectAltName=DNS:localhost,IP:127.0.0.1")

.PHONY: lint
lint: # 执行静态代码检查.
	@echo "===========> Running golangci to lint source codes"
	@golangci-lint run -c $(PROJ_ROOT_DIR)/.golangci.yaml $(PROJ_ROOT_DIR)/...


.PHONY: dev
dev:
	@docker compose -f $(PROJ_ROOT_DIR)/deployment/docker/dev/docker-compose.env.yml up -d
	@mkdir -p $(PROJ_ROOT_DIR)/tmp
	@air
