#!/bin/bash

# 环境变量配置脚本
# 用于设置 WEATHER_API_KEY 环境变量

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"
ENV_EXAMPLE="${SCRIPT_DIR}/.env.example"

echo "=== MCP Weather Server 环境配置 ==="
echo ""

# 检查是否已有 .env 文件
if [ -f "$ENV_FILE" ]; then
    echo "发现已存在的 .env 文件"
    CURRENT_KEY=$(grep "^WEATHER_API_KEY=" "$ENV_FILE" | cut -d'=' -f2- | tr -d '"' | tr -d "'")
    if [ -n "$CURRENT_KEY" ] && [ "$CURRENT_KEY" != "your_api_key_here" ]; then
        echo "当前配置的 API Key: ${CURRENT_KEY:0:10}... (长度: ${#CURRENT_KEY})"
        echo ""
        read -p "是否要更新 API Key? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "保持现有配置"
            exit 0
        fi
    fi
fi

# 提示用户输入 API Key
echo "请输入你的 OpenWeatherMap API Key"
echo "（如果还没有，请访问 https://openweathermap.org/api 注册获取）"
echo ""
read -p "API Key: " API_KEY

# 验证输入
if [ -z "$API_KEY" ]; then
    echo "错误: API Key 不能为空"
    exit 1
fi

if [ ${#API_KEY} -lt 20 ]; then
    echo "警告: API Key 长度似乎不正确（通常为 32 个字符）"
    read -p "是否继续？(y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 创建或更新 .env 文件
if [ -f "$ENV_FILE" ]; then
    # 更新现有文件
    if grep -q "^WEATHER_API_KEY=" "$ENV_FILE"; then
        # 使用 sed 更新（macOS 兼容）
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s|^WEATHER_API_KEY=.*|WEATHER_API_KEY=${API_KEY}|" "$ENV_FILE"
        else
            sed -i "s|^WEATHER_API_KEY=.*|WEATHER_API_KEY=${API_KEY}|" "$ENV_FILE"
        fi
    else
        echo "WEATHER_API_KEY=${API_KEY}" >> "$ENV_FILE"
    fi
else
    # 创建新文件
    cat > "$ENV_FILE" << EOF
# OpenWeatherMap API Key
# 获取方式：
# 1. 访问 https://openweathermap.org/api
# 2. 注册/登录账号
# 3. 在 https://home.openweathermap.org/api_keys 创建 API Key
# 4. 复制 API Key 到下面的配置中
WEATHER_API_KEY=${API_KEY}
EOF
fi

echo ""
echo "✅ 配置已保存到 .env 文件"
echo ""
echo "现在你可以："
echo "  1. 运行 ./start_server.sh 启动服务器"
echo "  2. 或者运行 ./test_api_key.sh 测试 API Key"
echo ""


