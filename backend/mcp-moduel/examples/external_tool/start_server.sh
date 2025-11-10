#!/bin/bash

# MCP Weather Server 启动脚本
# 支持从 .env 文件或环境变量读取配置

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_DIR="${SCRIPT_DIR}/server"

# 检查 .env 文件是否存在
if [ -f "${SCRIPT_DIR}/.env" ]; then
    echo "从 .env 文件加载配置..."
    # 读取 .env 文件并导出环境变量
    set -a
    source "${SCRIPT_DIR}/.env"
    set +a
fi

# 检查 API Key 是否设置
if [ -z "$WEATHER_API_KEY" ]; then
    echo "错误: WEATHER_API_KEY 环境变量未设置！"
    echo ""
    echo "请选择以下方式之一："
    echo ""
    echo "方式 1: 创建 .env 文件（推荐）"
    echo "  1. 复制 .env.example 为 .env:"
    echo "     cp .env.example .env"
    echo "  2. 编辑 .env 文件，填入你的 API Key"
    echo "  3. 重新运行此脚本"
    echo ""
    echo "方式 2: 直接设置环境变量"
    echo "  export WEATHER_API_KEY=your_api_key_here"
    echo "  ./start_server.sh"
    echo ""
    echo "方式 3: 启动时临时设置"
    echo "  WEATHER_API_KEY=your_api_key_here ./start_server.sh"
    echo ""
    echo "获取 API Key:"
    echo "  访问 https://openweathermap.org/api 注册并获取免费 API Key"
    exit 1
fi

# 检查 API Key 格式（至少 20 个字符）
if [ ${#WEATHER_API_KEY} -lt 20 ]; then
    echo "警告: API Key 长度似乎不正确（通常为 32 个字符）"
    echo "当前长度: ${#WEATHER_API_KEY}"
    echo ""
    read -p "是否继续？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "✅ API Key 已配置 (长度: ${#WEATHER_API_KEY})"
echo ""

# 切换到 server 目录
cd "${SERVER_DIR}" || exit 1

# 检查是否已编译
if [ ! -f "./server" ]; then
    echo "正在编译服务器..."
    go build -o server main.go
    if [ $? -ne 0 ]; then
        echo "编译失败！"
        exit 1
    fi
    echo "编译完成"
    echo ""
fi

# 启动服务器
echo "启动 MCP Weather Server..."
echo "按 Ctrl+C 停止服务器"
echo ""
./server


