#!/bin/bash

# Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/clin211/miniblog-v2.git.

# Docker 文件传输工具脚本

set -e

CONTAINER_NAME="miniblog-app"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 显示帮助信息
show_help() {
    cat <<EOF
Docker 文件传输工具

用法:
    $0 [选项] [操作]

操作:
    copy-config [file]     复制配置文件到容器
    copy-logs             从容器复制日志文件
    copy-to-container [src] [dst]  复制指定文件到容器
    copy-from-container [src] [dst]  从容器复制指定文件
    list-container-files [path]     列出容器中的文件
    exec-container [cmd]           在容器中执行命令

选项:
    -c, --container NAME   指定容器名称 (默认: $CONTAINER_NAME)
    -h, --help            显示帮助信息

示例:
    $0 copy-config mb-apiserver.yaml
    $0 copy-logs
    $0 copy-to-container ./configs/app.yaml /app/configs/app.yaml
    $0 copy-from-container /app/logs/app.log ./logs/app.log
    $0 list-container-files /app/configs
    $0 exec-container "ls -la /app"
EOF
}

# 检查容器是否存在且运行中
check_container() {
    if ! docker ps | grep -q "$CONTAINER_NAME"; then
        echo "错误: 容器 '$CONTAINER_NAME' 不存在或未运行"
        echo "请先启动容器: docker-compose up -d"
        exit 1
    fi
}

# 复制配置文件到容器
copy_config() {
    local config_file="$1"
    if [ -z "$config_file" ]; then
        config_file="mb-apiserver.yaml"
    fi

    local src_path="$PROJECT_ROOT/configs/$config_file"
    local dst_path="/app/configs/$config_file"

    if [ ! -f "$src_path" ]; then
        echo "错误: 配置文件 '$src_path' 不存在"
        exit 1
    fi

    echo "复制配置文件: $src_path -> $CONTAINER_NAME:$dst_path"
    docker cp "$src_path" "$CONTAINER_NAME:$dst_path"
    echo "配置文件复制完成"

    # 询问是否重启容器
    read -p "是否重启容器以应用新配置? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "重启容器中..."
        docker restart "$CONTAINER_NAME"
        echo "容器重启完成"
    fi
}

# 从容器复制日志文件
copy_logs() {
    local log_dir="$PROJECT_ROOT/logs"
    mkdir -p "$log_dir"

    echo "从容器复制日志文件到: $log_dir"
    docker cp "$CONTAINER_NAME:/app/logs/." "$log_dir/"
    echo "日志文件复制完成"

    echo "可用的日志文件:"
    ls -la "$log_dir"
}

# 复制文件到容器
copy_to_container() {
    local src="$1"
    local dst="$2"

    if [ -z "$src" ] || [ -z "$dst" ]; then
        echo "错误: 请指定源文件和目标路径"
        echo "用法: $0 copy-to-container <源文件> <容器内目标路径>"
        exit 1
    fi

    echo "复制文件: $src -> $CONTAINER_NAME:$dst"
    docker cp "$src" "$CONTAINER_NAME:$dst"
    echo "文件复制完成"
}

# 从容器复制文件
copy_from_container() {
    local src="$1"
    local dst="$2"

    if [ -z "$src" ] || [ -z "$dst" ]; then
        echo "错误: 请指定容器内源文件和目标路径"
        echo "用法: $0 copy-from-container <容器内源文件> <目标路径>"
        exit 1
    fi

    echo "复制文件: $CONTAINER_NAME:$src -> $dst"
    docker cp "$CONTAINER_NAME:$src" "$dst"
    echo "文件复制完成"
}

# 列出容器中的文件
list_container_files() {
    local path="$1"
    if [ -z "$path" ]; then
        path="/app"
    fi

    echo "列出容器 '$CONTAINER_NAME' 中 '$path' 的文件:"
    docker exec "$CONTAINER_NAME" ls -la "$path"
}

# 在容器中执行命令
exec_container() {
    local cmd="$1"
    if [ -z "$cmd" ]; then
        echo "错误: 请指定要执行的命令"
        exit 1
    fi

    echo "在容器 '$CONTAINER_NAME' 中执行: $cmd"
    docker exec "$CONTAINER_NAME" sh -c "$cmd"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
    -c | --container)
        CONTAINER_NAME="$2"
        shift 2
        ;;
    -h | --help)
        show_help
        exit 0
        ;;
    copy-config)
        check_container
        copy_config "$2"
        exit 0
        ;;
    copy-logs)
        check_container
        copy_logs
        exit 0
        ;;
    copy-to-container)
        check_container
        copy_to_container "$2" "$3"
        exit 0
        ;;
    copy-from-container)
        check_container
        copy_from_container "$2" "$3"
        exit 0
        ;;
    list-container-files)
        check_container
        list_container_files "$2"
        exit 0
        ;;
    exec-container)
        check_container
        exec_container "$2"
        exit 0
        ;;
    *)
        echo "未知选项: $1"
        show_help
        exit 1
        ;;
    esac
done

# 如果没有指定操作，显示帮助
show_help
