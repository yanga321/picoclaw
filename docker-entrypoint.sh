#!/bin/sh
# Map Railway's PORT to PICOCLAW_GATEWAY_PORT if not already set
if [ -n "$PORT" ] && [ -z "$PICOCLAW_GATEWAY_PORT" ]; then
    export PICOCLAW_GATEWAY_PORT="$PORT"
fi

exec picoclaw "$@"
