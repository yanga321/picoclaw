#!/bin/sh
# Map Railway's PORT to PICOCLAW_GATEWAY_PORT if not already set
if [ -n "$PORT" ] && [ -z "$PICOCLAW_GATEWAY_PORT" ]; then
    export PICOCLAW_GATEWAY_PORT="$PORT"
fi

# Generate config.json from environment variables if an API key is provided
CONFIG_DIR="$HOME/.picoclaw"
CONFIG_FILE="$CONFIG_DIR/config.json"

# Priority: DeepSeek > Kimi > OpenRouter
if [ -n "$DEEPSEEK_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-deepseek/deepseek-chat}"
    mkdir -p "$CONFIG_DIR"
    cat > "$CONFIG_FILE" <<EOCFG
{
  "agents": {
    "defaults": {
      "model": "$PICOCLAW_MODEL"
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
      "model": "$PICOCLAW_MODEL"
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
elif [ -n "$OPENROUTER_API_KEY" ]; then
    PICOCLAW_MODEL="${PICOCLAW_MODEL:-google/gemini-2.0-flash-exp:free}"
    mkdir -p "$CONFIG_DIR"
    cat > "$CONFIG_FILE" <<EOCFG
{
  "agents": {
    "defaults": {
      "model": "$PICOCLAW_MODEL"
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
fi

exec picoclaw "$@"
