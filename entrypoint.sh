#!/bin/sh
set -e

mkdir -p /app/logs
chown -R appuser:appgroup /app/logs || true
chmod 775 /app/logs || true

# 使用 gosu 而非 su
exec gosu appuser ./main
