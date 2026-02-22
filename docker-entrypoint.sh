#!/bin/sh
# Map Railway's PORT to PICOCLAW_GATEWAY_PORT if not already set
if [ -n "$PORT" ] && [ -z "$PICOCLAW_GATEWAY_PORT" ]; then
    export PICOCLAW_GATEWAY_PORT="$PORT"
fi

# Generate config.json from environment variables
CONFIG_DIR="$HOME/.picoclaw"
CONFIG_FILE="$CONFIG_DIR/config.json"
mkdir -p "$CONFIG_DIR"

# Detect which API key is available
# Priority: OpenRouter (primary) > DeepSeek > Kimi
API_KEY=""
API_BASE=""
PROVIDER_NAME=""

if [ -n "$OPENROUTER_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-openrouter/deepseek/deepseek-chat-v3-0324}"
    API_KEY="$OPENROUTER_API_KEY"
    API_BASE="https://openrouter.ai/api/v1"
    PROVIDER_NAME="OpenRouter"
elif [ -n "$DEEPSEEK_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-deepseek/deepseek-chat}"
    API_KEY="$DEEPSEEK_API_KEY"
    API_BASE="https://api.deepseek.com/v1"
    PROVIDER_NAME="DeepSeek"
elif [ -n "$KIMI_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-moonshot/kimi-k2.5}"
    API_KEY="$KIMI_API_KEY"
    API_BASE="https://api.moonshot.cn/v1"
    PROVIDER_NAME="Kimi/Moonshot"
fi

if [ -z "$API_KEY" ]; then
    echo "============================================"
    echo "  WARNING: No LLM API key configured!"
    echo "============================================"
    echo ""
    echo "  Set one of these environment variables:"
    echo "    OPENROUTER_API_KEY  (recommended, https://openrouter.ai/keys)"
    echo "    DEEPSEEK_API_KEY    (https://platform.deepseek.com/)"
    echo "    KIMI_API_KEY        (https://platform.moonshot.cn/)"
    echo ""
    echo "  Without an API key, the bot cannot respond to messages."
    echo "============================================"
    # Still generate a config so the app starts (web UI works, just no LLM)
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-openrouter/deepseek/deepseek-chat-v3-0324}"
    cat > "$CONFIG_FILE" <<EOCFG
{
  "agents": {
    "defaults": {
      "model": "$PICOCLAW_MODEL",
      "max_tokens": 8192
    }
  },
  "model_list": []
}
EOCFG
else
    echo "✓ LLM Provider: $PROVIDER_NAME"
    echo "✓ Model: $PICOCLAW_MODEL"
    cat > "$CONFIG_FILE" <<EOCFG
{
  "agents": {
    "defaults": {
      "model": "$PICOCLAW_MODEL",
      "max_tokens": 8192
    }
  },
  "model_list": [
    {
      "model_name": "$PICOCLAW_MODEL",
      "model": "$PICOCLAW_MODEL",
      "api_base": "$API_BASE",
      "api_key": "$API_KEY"
    }
  ]
}
EOCFG
fi

exec picoclaw "$@"
