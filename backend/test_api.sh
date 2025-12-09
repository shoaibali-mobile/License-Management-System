#!/bin/bash

# Complete API Test Script using cURL
# Usage: bash test_api.sh

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "License MNM API Test Script"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== 1. Admin Login ===${NC}"
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -z "$ADMIN_TOKEN" ]; then
  echo "❌ Admin login failed"
  exit 1
fi
echo -e "${GREEN}✓ Admin Token: ${ADMIN_TOKEN:0:50}...${NC}"
echo ""

echo -e "${BLUE}=== 2. Create Subscription Pack ===${NC}"
PACK_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Premium Plan",
    "description": "Full access to all features",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }')
echo "$PACK_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Subscription pack created${NC}" || echo -e "${YELLOW}⚠ Pack may already exist${NC}"
echo ""

echo -e "${BLUE}=== 3. Create Customer ===${NC}"
CUSTOMER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "customer@example.com",
    "phone": "+1234567890"
  }')
echo "$CUSTOMER_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Customer created${NC}" || echo -e "${YELLOW}⚠ Customer may already exist${NC}"
echo ""

echo -e "${BLUE}=== 4. Customer Login ===${NC}"
CUSTOMER_LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/customer/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }')
CUSTOMER_TOKEN=$(echo $CUSTOMER_LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
if [ -z "$CUSTOMER_TOKEN" ]; then
  echo "❌ Customer login failed"
  exit 1
fi
echo -e "${GREEN}✓ Customer Token: ${CUSTOMER_TOKEN:0:50}...${NC}"
echo ""

echo -e "${BLUE}=== 5. Request Subscription (Customer) ===${NC}"
REQUEST_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/customer/subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -d '{
    "sku": "premium-plan"
  }')
SUBSCRIPTION_ID=$(echo $REQUEST_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "$REQUEST_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Subscription requested (ID: $SUBSCRIPTION_ID)${NC}" || echo -e "${YELLOW}⚠ Request may have failed${NC}"
echo ""

if [ ! -z "$SUBSCRIPTION_ID" ]; then
  echo -e "${BLUE}=== 6. Approve Subscription (Admin) ===${NC}"
  APPROVE_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/subscriptions/$SUBSCRIPTION_ID/approve \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  echo "$APPROVE_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Subscription approved${NC}" || echo -e "${YELLOW}⚠ Approval may have failed${NC}"
  echo ""
fi

echo -e "${BLUE}=== 7. Assign Subscription (Admin) ===${NC}"
ASSIGN_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/admin/customers/1/assign-subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "pack_id": 1
  }')
echo "$ASSIGN_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Subscription assigned${NC}" || echo -e "${YELLOW}⚠ Assignment may have failed (customer may already have active subscription)${NC}"
echo ""

echo -e "${BLUE}=== 8. Get Customer Subscription ===${NC}"
GET_SUB_RESPONSE=$(curl -s -X GET $BASE_URL/api/v1/customer/subscription \
  -H "Authorization: Bearer $CUSTOMER_TOKEN")
echo "$GET_SUB_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Subscription retrieved${NC}" || echo -e "${YELLOW}⚠ No active subscription found${NC}"
echo "$GET_SUB_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GET_SUB_RESPONSE"
echo ""

echo -e "${BLUE}=== 9. Get Admin Dashboard ===${NC}"
DASHBOARD_RESPONSE=$(curl -s -X GET $BASE_URL/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "$DASHBOARD_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ Dashboard data retrieved${NC}" || echo -e "${YELLOW}⚠ Dashboard request failed${NC}"
echo "$DASHBOARD_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$DASHBOARD_RESPONSE"
echo ""

echo -e "${BLUE}=== 10. SDK Login ===${NC}"
SDK_RESPONSE=$(curl -s -X POST $BASE_URL/sdk/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }')
API_KEY=$(echo $SDK_RESPONSE | grep -o '"api_key":"[^"]*' | cut -d'"' -f4)
if [ -z "$API_KEY" ]; then
  echo "❌ SDK login failed"
else
  echo -e "${GREEN}✓ API Key: $API_KEY${NC}"
  echo ""
  
  echo -e "${BLUE}=== 11. Get Subscription (SDK) ===${NC}"
  SDK_SUB_RESPONSE=$(curl -s -X GET $BASE_URL/sdk/v1/subscription \
    -H "X-API-Key: $API_KEY")
  echo "$SDK_SUB_RESPONSE" | grep -q "success" && echo -e "${GREEN}✓ SDK subscription retrieved${NC}" || echo -e "${YELLOW}⚠ No active subscription found${NC}"
  echo "$SDK_SUB_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$SDK_SUB_RESPONSE"
fi

echo ""
echo -e "${GREEN}=========================================="
echo "Test Script Completed!"
echo "==========================================${NC}"





