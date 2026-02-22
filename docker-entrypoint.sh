#!/bin/sh
# Map Railway's PORT to PICOCLAW_GATEWAY_PORT if not already set
if [ -n "$PORT" ] && [ -z "$PICOCLAW_GATEWAY_PORT" ]; then
    export PICOCLAW_GATEWAY_PORT="$PORT"
fi

# Generate config.json from environment variables if an API key is provided
CONFIG_DIR="$HOME/.picoclaw"
CONFIG_FILE="$CONFIG_DIR/config.json"

# Priority: OpenRouter (primary) > DeepSeek > Kimi
if [ -n "$OPENROUTER_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-openrouter/deepseek/deepseek-chat-v3-0324}"
    mkdir -p "$CONFIG_DIR"
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
      "api_base": "https://openrouter.ai/api/v1",
      "api_key": "$OPENROUTER_API_KEY"
    }
  ]
}
EOCFG
elif [ -n "$DEEPSEEK_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-deepseek/deepseek-chat}"
    mkdir -p "$CONFIG_DIR"
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
      "api_base": "https://api.deepseek.com/v1",
      "api_key": "$DEEPSEEK_API_KEY"
    }
  ]
}
EOCFG
elif [ -n "$KIMI_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-moonshot/kimi-k2.5}"
    mkdir -p "$CONFIG_DIR"
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
      "api_base": "https://api.moonshot.cn/v1",
      "api_key": "$KIMI_API_KEY"
    }
  ]
}
EOCFG
fi

exec picoclaw "$@"
