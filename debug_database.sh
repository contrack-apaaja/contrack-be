#!/bin/bash

# Debug database to check if data exists
BASE_URL="http://localhost:8080/api"

echo "=== Debug Database ==="

# 1. Login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}')

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# 2. Check if we can get analyses
echo "2. Checking existing analyses..."
ANALYSES_RESPONSE=$(curl -s -X GET "${BASE_URL}/ai/analyses" \
  -H "Authorization: Bearer $TOKEN")

echo "Analyses response: $ANALYSES_RESPONSE"

# 3. Check stats
echo "3. Checking AI stats..."
STATS_RESPONSE=$(curl -s -X GET "${BASE_URL}/ai/stats" \
  -H "Authorization: Bearer $TOKEN")

echo "Stats response: $STATS_RESPONSE"

echo "=== Debug Complete ==="
