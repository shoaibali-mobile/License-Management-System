#!/bin/bash

# Quick Test Script - Step by Step Workflow
# This script demonstrates the complete workflow

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "Complete API Workflow Test"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# STEP 1: Admin Login
echo -e "${BLUE}STEP 1: Admin Login${NC}"
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -z "$ADMIN_TOKEN" ]; then
  echo -e "${RED}❌ Admin login failed${NC}"
  exit 1
fi
echo -e "${GREEN}✓ Admin logged in. Token: ${ADMIN_TOKEN:0:30}...${NC}"
echo ""

# STEP 2: Create Subscription Pack
echo -e "${BLUE}STEP 2: Create Subscription Pack${NC}"
PACK_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Premium Plan","description":"Full access","sku":"premium-plan","price":29.99,"validity_months":12}')
PACK_ID=$(echo $PACK_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ ! -z "$PACK_ID" ]; then
  echo -e "${GREEN}✓ Subscription pack created (ID: $PACK_ID)${NC}"
else
  echo -e "${YELLOW}⚠ Pack may already exist${NC}"
  PACK_ID=1
fi
echo ""

# STEP 3: Create Customer
echo -e "${BLUE}STEP 3: Create Customer${NC}"
CUSTOMER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"John Doe","email":"customer@example.com","phone":"+1234567890"}')
CUSTOMER_ID=$(echo $CUSTOMER_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ ! -z "$CUSTOMER_ID" ]; then
  echo -e "${GREEN}✓ Customer created (ID: $CUSTOMER_ID)${NC}"
else
  echo -e "${YELLOW}⚠ Customer may already exist${NC}"
  CUSTOMER_ID=1
fi
echo ""

# STEP 4: Customer Login
echo -e "${BLUE}STEP 4: Customer Login${NC}"
CUSTOMER_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/customer/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123"}')
CUSTOMER_TOKEN=$(echo $CUSTOMER_LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -z "$CUSTOMER_TOKEN" ]; then
  echo -e "${RED}❌ Customer login failed${NC}"
  exit 1
fi
echo -e "${GREEN}✓ Customer logged in. Token: ${CUSTOMER_TOKEN:0:30}...${NC}"
echo ""

# STEP 5: Customer Requests Subscription
echo -e "${BLUE}STEP 5: Customer Requests Subscription${NC}"
REQUEST_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/customer/subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -d '{"sku":"premium-plan"}')
SUBSCRIPTION_ID=$(echo $REQUEST_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
if [ ! -z "$SUBSCRIPTION_ID" ]; then
  echo -e "${GREEN}✓ Subscription requested (ID: $SUBSCRIPTION_ID)${NC}"
else
  echo -e "${YELLOW}⚠ Request may have failed or customer already has active subscription${NC}"
  echo "Trying to get existing subscription..."
  EXISTING_SUB=$(curl -s -X GET $BASE_URL/api/v1/customer/subscription \
    -H "Authorization: Bearer $CUSTOMER_TOKEN")
  if echo "$EXISTING_SUB" | grep -q "success"; then
    echo -e "${GREEN}✓ Customer already has an active subscription${NC}"
    exit 0
  fi
  exit 1
fi
echo ""

# STEP 6: Admin Approves Subscription
echo -e "${BLUE}STEP 6: Admin Approves Subscription${NC}"
APPROVE_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/subscriptions/$SUBSCRIPTION_ID/approve \
  -H "Authorization: Bearer $ADMIN_TOKEN")
if echo "$APPROVE_RESPONSE" | grep -q "success"; then
  echo -e "${GREEN}✓ Subscription approved${NC}"
else
  echo -e "${YELLOW}⚠ Approval may have failed${NC}"
fi
echo ""

# STEP 7: Admin Assigns Subscription
echo -e "${BLUE}STEP 7: Admin Assigns Subscription (Makes it Active)${NC}"
ASSIGN_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/customers/$CUSTOMER_ID/assign-subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{\"pack_id\":$PACK_ID}")
if echo "$ASSIGN_RESPONSE" | grep -q "success"; then
  echo -e "${GREEN}✓ Subscription assigned and activated${NC}"
else
  echo -e "${YELLOW}⚠ Assignment may have failed (customer may already have active subscription)${NC}"
fi
echo ""

# STEP 8: Customer Views Subscription
echo -e "${BLUE}STEP 8: Customer Views Active Subscription${NC}"
GET_SUB_RESPONSE=$(curl -s -X GET $BASE_URL/api/v1/customer/subscription \
  -H "Authorization: Bearer $CUSTOMER_TOKEN")
if echo "$GET_SUB_RESPONSE" | grep -q "success"; then
  echo -e "${GREEN}✓ Subscription retrieved successfully${NC}"
  echo "$GET_SUB_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_SUB_RESPONSE"
else
  echo -e "${YELLOW}⚠ No active subscription found${NC}"
fi
echo ""

# STEP 9: Admin Views Dashboard
echo -e "${BLUE}STEP 9: Admin Views Dashboard${NC}"
DASHBOARD_RESPONSE=$(curl -s -X GET $BASE_URL/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN")
if echo "$DASHBOARD_RESPONSE" | grep -q "success"; then
  echo -e "${GREEN}✓ Dashboard data retrieved${NC}"
  echo "$DASHBOARD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DASHBOARD_RESPONSE"
fi
echo ""

# STEP 10: SDK Login
echo -e "${BLUE}STEP 10: SDK Login (Get API Key)${NC}"
SDK_RESPONSE=$(curl -s -X POST $BASE_URL/sdk/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"password123"}')
API_KEY=$(echo $SDK_RESPONSE | grep -o '"api_key":"[^"]*' | cut -d'"' -f4)
if [ ! -z "$API_KEY" ]; then
  echo -e "${GREEN}✓ SDK login successful. API Key: $API_KEY${NC}"
  echo ""
  
  echo -e "${BLUE}STEP 11: Get Subscription via SDK${NC}"
  SDK_SUB_RESPONSE=$(curl -s -X GET $BASE_URL/sdk/v1/subscription \
    -H "X-API-Key: $API_KEY")
  if echo "$SDK_SUB_RESPONSE" | grep -q "success"; then
    echo -e "${GREEN}✓ SDK subscription retrieved${NC}"
    echo "$SDK_SUB_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$SDK_SUB_RESPONSE"
  else
    echo -e "${YELLOW}⚠ No active subscription found via SDK${NC}"
  fi
else
  echo -e "${RED}❌ SDK login failed${NC}"
fi

echo ""
echo -e "${GREEN}=========================================="
echo "Workflow Test Completed!"
echo "==========================================${NC}"





