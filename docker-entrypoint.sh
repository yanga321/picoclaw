#!/bin/sh
# Map Railway's PORT to PICOCLAW_GATEWAY_PORT if not already set
if [ -n "$PORT" ] && [ -z "$PICOCLAW_GATEWAY_PORT" ]; then
    export PICOCLAW_GATEWAY_PORT="$PORT"
fi

# Generate config.json from environment variables if an API key is provided
CONFIG_DIR="$HOME/.picoclaw"
CONFIG_FILE="$CONFIG_DIR/config.json"

if [ -n "$PICOCLAW_PROVIDERS_ANTHROPIC_API_KEY" ]; then
    mkdir -p "$CONFIG_DIR"
    cat > "$CONFIG_FILE" <<EOCFG
{
  "agents": {
    "defaults": {
      "model": "claude-sonnet-4-6"
    }
  },
  "model_list": [
    {
      "model_name": "claude-sonnet-4-6",
      "model": "anthropic/claude-sonnet-4-6",
      "api_base": "https://api.anthropic.com/v1",
      "api_key": "$PICOCLAW_PROVIDERS_ANTHROPIC_API_KEY"
    }
  ]
}
EOCFG
fi

exec picoclaw "$@"
