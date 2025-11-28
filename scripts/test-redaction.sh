#!/bin/bash

echo "====================================="
echo "Cilium-Shield Redaction Test Script"
echo "====================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Backend URL
BACKEND_URL="${BACKEND_URL:-http://localhost:3001}"

echo "Testing backend at: $BACKEND_URL"
echo ""

# Test 1: Health Check
echo "Test 1: Health Check"
echo "---------------------"
RESPONSE=$(curl -s "${BACKEND_URL}/")
if echo "$RESPONSE" | grep -q "ok"; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend is not responding${NC}"
    exit 1
fi
echo ""

# Test 2: Submit Credit Card Redaction Event
echo "Test 2: Submit Credit Card Redaction Event"
echo "-------------------------------------------"
curl -X POST "${BACKEND_URL}/api/v1/events" \
  -H "Content-Type: application/json" \
  -d '{
    "source_pod_ip": "10.0.1.42",
    "destination_url": "https://api.openai.com/v1/completions",
    "redacted_type": "REDACTED_CREDIT_CARD"
  }'
echo ""
echo -e "${GREEN}✓ Event submitted${NC}"
echo ""

# Test 3: Submit API Key Redaction Event
echo "Test 3: Submit API Key Redaction Event"
echo "---------------------------------------"
curl -X POST "${BACKEND_URL}/api/v1/events" \
  -H "Content-Type: application/json" \
  -d '{
    "source_pod_ip": "10.0.1.43",
    "destination_url": "https://api.anthropic.com/v1/messages",
    "redacted_type": "REDACTED_API_KEY"
  }'
echo ""
echo -e "${GREEN}✓ Event submitted${NC}"
echo ""

# Test 4: Submit Email Redaction Event
echo "Test 4: Submit Email Redaction Event"
echo "-------------------------------------"
curl -X POST "${BACKEND_URL}/api/v1/events" \
  -H "Content-Type: application/json" \
  -d '{
    "source_pod_ip": "10.0.1.44",
    "destination_url": "https://api.openai.com/v1/chat/completions",
    "redacted_type": "REDACTED_EMAIL"
  }'
echo ""
echo -e "${GREEN}✓ Event submitted${NC}"
echo ""

# Test 5: List All Events
echo "Test 5: List All Events"
echo "-----------------------"
EVENTS=$(curl -s "${BACKEND_URL}/api/v1/events/list")
EVENT_COUNT=$(echo "$EVENTS" | jq '. | length' 2>/dev/null || echo "0")
echo "Total events: $EVENT_COUNT"
echo "$EVENTS" | jq '.' 2>/dev/null || echo "$EVENTS"
echo ""

# Test 6: Get Statistics
echo "Test 6: Get Event Statistics"
echo "-----------------------------"
STATS=$(curl -s "${BACKEND_URL}/api/v1/events/stats")
echo "$STATS" | jq '.' 2>/dev/null || echo "$STATS"
echo ""

echo "====================================="
echo -e "${GREEN}All tests completed!${NC}"
echo "====================================="
echo ""
echo "Open the dashboard at http://localhost:3000 to view the events"
