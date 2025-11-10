#!/bin/bash

# 测试 OpenWeatherMap API Key 是否有效
# 使用方法: ./test_api_key.sh [your_api_key]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 尝试从 .env 文件读取配置
if [ -f "${SCRIPT_DIR}/.env" ]; then
    set -a
    source "${SCRIPT_DIR}/.env"
    set +a
fi

API_KEY="${1:-${WEATHER_API_KEY}}"

if [ -z "$API_KEY" ]; then
    echo "错误: 请提供 API Key"
    echo "使用方法:"
    echo "  ./test_api_key.sh your_api_key"
    echo "  或者设置环境变量: export WEATHER_API_KEY=your_api_key && ./test_api_key.sh"
    exit 1
fi

echo "正在测试 API Key: ${API_KEY:0:10}..."
echo ""

# 测试查询北京的天气
CITY="Beijing"
URL="https://api.openweathermap.org/data/2.5/weather?q=${CITY}&units=metric&appid=${API_KEY}"

echo "请求 URL: ${URL}"
echo ""

# 发送请求
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "${URL}")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

echo "HTTP 状态码: $HTTP_STATUS"
echo "响应内容:"
echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
echo ""

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✅ API Key 有效！可以正常使用。"
    exit 0
elif [ "$HTTP_STATUS" = "401" ]; then
    echo "❌ API Key 无效或已过期。"
    echo ""
    echo "可能的原因:"
    echo "1. API Key 不正确"
    echo "2. API Key 尚未激活（新注册的 Key 可能需要几分钟）"
    echo "3. API Key 已过期或被撤销"
    echo ""
    echo "请检查:"
    echo "- 访问 https://home.openweathermap.org/api_keys 确认 Key 是否正确"
    echo "- 如果是新注册的 Key，请等待几分钟后重试"
    exit 1
elif [ "$HTTP_STATUS" = "429" ]; then
    echo "⚠️  请求过于频繁，请稍后再试。"
    exit 1
else
    echo "❌ 请求失败，状态码: $HTTP_STATUS"
    exit 1
fi

