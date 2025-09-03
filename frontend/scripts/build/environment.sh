#!/bin/bash
set -euo pipefail

DOMAIN=${DOMAIN:-"localhost:8080"}

echo "Setting DOMAIN to: $DOMAIN in config.ts"

CONFIG_FILE="src/ts/config/config.ts"

if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "❌ Error: Config file not found at $CONFIG_FILE"
    exit 1
fi

TEMP_FILE=$(mktemp)

sed "s|const HOST = 'localhost:8080';|const HOST = '$DOMAIN';|" "$CONFIG_FILE" > "$TEMP_FILE"

mv "$TEMP_FILE" "$CONFIG_FILE"

echo "✅ Successfully updated config.ts with HOST: $DOMAIN"
