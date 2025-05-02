#!/bin/sh

set -e

echo "🔄 Running ClickHouse migrations..."

for file in /migrations/*.sql; do
  echo "📄 Executing migration: $file"
  echo "$CLICKHOUSE_USER | $CLICKHOUSE_HOST"
  curl -sS -u "$CLICKHOUSE_USER:$CLICKHOUSE_PASSWORD" \
    --data-binary "@$file" \
    "http://$CLICKHOUSE_HOST:$CLICKHOUSE_PORT/?database=$CLICKHOUSE_DB" \
    && echo "✅ Applied: $file" \
    || echo "❌ Error: $file"
done
