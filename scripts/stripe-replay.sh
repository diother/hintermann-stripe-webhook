#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Usage: ./stripe-replay.sh <event.json>"
  exit 1
fi

if [ -z "$WEBHOOK_SECRET" ]; then
  echo "Error: WEBHOOK_SECRET is not set"
  echo "Run: export WEBHOOK_SECRET=whsec_..."
  exit 1
fi

EVENT_FILE="$1"

if [ ! -f "$EVENT_FILE" ]; then
  echo "File not found: $EVENT_FILE"
  exit 1
fi

payload=$(cat "$EVENT_FILE")

timestamp=$(date +%s)

signed_payload="${timestamp}.${payload}"

signature=$(printf "%s" "$signed_payload" \
  | openssl dgst -sha256 -hmac "$WEBHOOK_SECRET" \
  | sed 's/^.* //')

curl http://localhost:8080/webhook \
  -s \
  -H "Content-Type: application/json" \
  -H "Stripe-Signature: t=${timestamp},v1=${signature}" \
  -d "$payload"

echo "Sent: $EVENT_FILE"
