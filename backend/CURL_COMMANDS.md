# Complete cURL Commands for License MNM API

Base URL: `http://localhost:8080`

## Authentication Endpoints (No Auth Required)

### 1. Admin Login
```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**Response:** Copy the `token` from response for admin endpoints.

### 2. Customer Login
```bash
curl -X POST http://localhost:8080/api/customer/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

**Response:** Copy the `token` from response for customer endpoints.

### 3. Customer Signup
```bash
curl -X POST http://localhost:8080/api/customer/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "phone": "+1234567890"
  }'
```

---

## Admin Endpoints (JWT Bearer Token Required)

**Note:** Replace `YOUR_ADMIN_JWT_TOKEN` with the token from admin login.

### 4. Get Admin Dashboard
```bash
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 5. List Customers
```bash
# Basic list
curl -X GET "http://localhost:8080/api/v1/admin/customers?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"

# With search
curl -X GET "http://localhost:8080/api/v1/admin/customers?page=1&limit=10&search=john" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 6. Create Customer
```bash
curl -X POST http://localhost:8080/api/v1/admin/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "+1987654321"
  }'
```

**Note:** Default password is `password123`

### 7. Get Customer Details
```bash
curl -X GET http://localhost:8080/api/v1/admin/customers/1 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 8. Update Customer
```bash
curl -X PUT http://localhost:8080/api/v1/admin/customers/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "name": "Jane Smith Updated",
    "phone": "+1987654321"
  }'
```

### 9. Delete Customer (Soft Delete)
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/customers/1 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 10. List Subscription Packs
```bash
curl -X GET "http://localhost:8080/api/v1/admin/subscription-packs?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 11. Create Subscription Pack
```bash
curl -X POST http://localhost:8080/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "name": "Premium Plan",
    "description": "Full access to all features",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }'
```

### 12. Update Subscription Pack
```bash
curl -X PUT http://localhost:8080/api/v1/admin/subscription-packs/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "name": "Premium Plan Updated",
    "description": "Updated description",
    "price": 39.99,
    "validity_months": 12
  }'
```

### 13. Delete Subscription Pack (Soft Delete)
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/subscription-packs/1 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 14. List All Subscriptions
```bash
# All subscriptions
curl -X GET "http://localhost:8080/api/v1/admin/subscriptions?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"

# Filter by status
curl -X GET "http://localhost:8080/api/v1/admin/subscriptions?page=1&limit=10&status=active" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 15. Approve Subscription Request
```bash
curl -X POST http://localhost:8080/api/v1/admin/subscriptions/1/approve \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

### 16. Assign Subscription to Customer
```bash
curl -X POST http://localhost:8080/api/v1/admin/customers/1/assign-subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN" \
  -d '{
    "pack_id": 1
  }'
```

### 17. Unassign Subscription
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/customers/1/subscription/1 \
  -H "Authorization: Bearer YOUR_ADMIN_JWT_TOKEN"
```

---

## Customer Endpoints (JWT Bearer Token Required)

**Note:** Replace `YOUR_CUSTOMER_JWT_TOKEN` with the token from customer login.

### 18. Get Current Subscription
```bash
curl -X GET http://localhost:8080/api/v1/customer/subscription \
  -H "Authorization: Bearer YOUR_CUSTOMER_JWT_TOKEN"
```

### 19. Request New Subscription
```bash
curl -X POST http://localhost:8080/api/v1/customer/subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_CUSTOMER_JWT_TOKEN" \
  -d '{
    "sku": "premium-plan"
  }'
```

### 20. Deactivate Current Subscription
```bash
curl -X DELETE http://localhost:8080/api/v1/customer/subscription \
  -H "Authorization: Bearer YOUR_CUSTOMER_JWT_TOKEN"
```

### 21. Get Subscription History
```bash
# Default (descending order)
curl -X GET "http://localhost:8080/api/v1/customer/subscription-history?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_CUSTOMER_JWT_TOKEN"

# Ascending order
curl -X GET "http://localhost:8080/api/v1/customer/subscription-history?page=1&limit=10&sort=asc" \
  -H "Authorization: Bearer YOUR_CUSTOMER_JWT_TOKEN"
```

---

## SDK Endpoints

### 22. SDK Login (Get API Key)
```bash
curl -X POST http://localhost:8080/sdk/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

**Response:** Copy the `api_key` from response for SDK endpoints.

### 23. Get Current Subscription (SDK)
```bash
curl -X GET http://localhost:8080/sdk/v1/subscription \
  -H "X-API-Key: YOUR_API_KEY_HERE"
```

### 24. Request Subscription (SDK)
```bash
curl -X POST http://localhost:8080/sdk/v1/subscription \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY_HERE" \
  -d '{
    "pack_sku": "premium-plan"
  }'
```

### 25. Deactivate Subscription (SDK)
```bash
curl -X DELETE http://localhost:8080/sdk/v1/subscription \
  -H "X-API-Key: YOUR_API_KEY_HERE"
```

### 26. Get Subscription History (SDK)
```bash
# Default (descending order)
curl -X GET "http://localhost:8080/sdk/v1/subscription-history?page=1&limit=10" \
  -H "X-API-Key: YOUR_API_KEY_HERE"

# Ascending order
curl -X GET "http://localhost:8080/sdk/v1/subscription-history?page=1&limit=10&sort=asc" \
  -H "X-API-Key: YOUR_API_KEY_HERE"
```

---

## Complete Test Workflow (cURL Script)

Save this as `test_api.sh` and run: `bash test_api.sh`

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== 1. Admin Login ==="
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Admin Token: $ADMIN_TOKEN"
echo ""

echo "=== 2. Create Subscription Pack ==="
curl -s -X POST $BASE_URL/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Premium Plan",
    "description": "Full access to all features",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }' | jq .
echo ""

echo "=== 3. Create Customer ==="
curl -s -X POST $BASE_URL/api/v1/admin/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "customer@example.com",
    "phone": "+1234567890"
  }' | jq .
echo ""

echo "=== 4. Customer Login ==="
CUSTOMER_RESPONSE=$(curl -s -X POST $BASE_URL/api/customer/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }')
CUSTOMER_TOKEN=$(echo $CUSTOMER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Customer Token: $CUSTOMER_TOKEN"
echo ""

echo "=== 5. Request Subscription (Customer) ==="
curl -s -X POST $BASE_URL/api/v1/customer/subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" \
  -d '{
    "sku": "premium-plan"
  }' | jq .
echo ""

echo "=== 6. Approve Subscription (Admin) ==="
curl -s -X POST $BASE_URL/api/v1/admin/subscriptions/1/approve \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
echo ""

echo "=== 7. Assign Subscription (Admin) ==="
curl -s -X POST $BASE_URL/api/v1/admin/customers/1/assign-subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "pack_id": 1
  }' | jq .
echo ""

echo "=== 8. Get Customer Subscription ==="
curl -s -X GET $BASE_URL/api/v1/customer/subscription \
  -H "Authorization: Bearer $CUSTOMER_TOKEN" | jq .
echo ""

echo "=== 9. Get Admin Dashboard ==="
curl -s -X GET $BASE_URL/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
echo ""

echo "=== 10. SDK Login ==="
SDK_RESPONSE=$(curl -s -X POST $BASE_URL/sdk/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }')
API_KEY=$(echo $SDK_RESPONSE | grep -o '"api_key":"[^"]*' | cut -d'"' -f4)
echo "API Key: $API_KEY"
echo ""

echo "=== 11. Get Subscription (SDK) ==="
curl -s -X GET $BASE_URL/sdk/v1/subscription \
  -H "X-API-Key: $API_KEY" | jq .
```

---

## Using Variables in cURL

### Save tokens to variables (bash)
```bash
# Admin login and save token
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Use token in subsequent requests
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### Pretty print JSON responses
Add `| jq .` at the end of any curl command to format JSON:
```bash
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

### Save response to file
```bash
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -o response.json
```

### Verbose output (see headers)
```bash
curl -v -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## Quick Reference

### Authentication Headers

**Frontend APIs (JWT):**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**SDK APIs (API Key):**
```
X-API-Key: YOUR_API_KEY
```

### Common Response Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation error)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

### Tips

1. **Replace IDs**: Replace `1` in URLs with actual IDs from previous responses
2. **Check token expiration**: JWT tokens expire after 24 hours
3. **Use jq for JSON**: Install `jq` for pretty JSON output: `brew install jq` (macOS)
4. **Save responses**: Use `-o filename.json` to save responses
5. **Debug**: Use `-v` flag to see full request/response details





