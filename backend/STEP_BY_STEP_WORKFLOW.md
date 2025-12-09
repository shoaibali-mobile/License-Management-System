# Step-by-Step API Workflow with cURL Commands

Follow these steps in order to test the complete API workflow.

---

## 🎯 **STEP 1: Admin Login** (Get Admin JWT Token)

**Purpose:** Get authentication token for admin operations

```bash
curl -X POST http://localhost:8080/api/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "admin@example.com",
  "expires_in": 3600
}
```

**⚠️ IMPORTANT:** Copy the `token` value! You'll need it for all admin endpoints.

**Save to variable (optional):**
```bash
ADMIN_TOKEN="YOUR_TOKEN_HERE"
```

---

## 🎯 **STEP 2: Create Subscription Pack** (Admin)

**Purpose:** Create a subscription plan that customers can subscribe to

```bash
curl -X POST http://localhost:8080/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "Premium Plan",
    "description": "Full access to all features",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "pack": {
    "id": 1,
    "name": "Premium Plan",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }
}
```

**⚠️ IMPORTANT:** Note the `id` and `sku` - you'll need them later!

**With variable:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Premium Plan",
    "description": "Full access to all features",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }'
```

---

## 🎯 **STEP 3: Create Another Subscription Pack** (Optional - Admin)

**Purpose:** Create a second plan for variety

```bash
curl -X POST http://localhost:8080/api/v1/admin/subscription-packs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "Basic Plan",
    "description": "Basic features access",
    "sku": "basic-plan",
    "price": 9.99,
    "validity_months": 6
  }'
```

---

## 🎯 **STEP 4: List All Subscription Packs** (Admin)

**Purpose:** Verify packs were created and see all available plans

```bash
curl -X GET "http://localhost:8080/api/v1/admin/subscription-packs?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "packs": [
    {
      "id": 1,
      "name": "Premium Plan",
      "sku": "premium-plan",
      "price": 29.99
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

---

## 🎯 **STEP 5: Create Customer** (Admin)

**Purpose:** Create a customer account (admin can create customers directly)

```bash
curl -X POST http://localhost:8080/api/v1/admin/customers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "customer@example.com",
    "phone": "+1234567890"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "customer": {
    "id": 1,
    "name": "John Doe",
    "email": "customer@example.com",
    "phone": "+1234567890"
  }
}
```

**⚠️ IMPORTANT:** 
- Note the customer `id`
- Default password is `password123`
- Customer can now login with: `customer@example.com` / `password123`

---

## 🎯 **STEP 6: Customer Signup** (Alternative to Step 5)

**Purpose:** Customer can register themselves (no auth required)

```bash
curl -X POST http://localhost:8080/api/customer/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "password": "password123",
    "phone": "+1987654321"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Account created successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "name": "Jane Smith",
  "phone": "+1987654321",
  "expires_in": 3600
}
```

**Note:** This also returns a JWT token, so customer can immediately use it!

---

## 🎯 **STEP 7: Customer Login** (Get Customer JWT Token)

**Purpose:** Get authentication token for customer operations

```bash
curl -X POST http://localhost:8080/api/customer/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "name": "John Doe",
  "phone": "+1234567890",
  "expires_in": 3600
}
```

**⚠️ IMPORTANT:** Copy the `token` value for customer endpoints!

**Save to variable:**
```bash
CUSTOMER_TOKEN="YOUR_CUSTOMER_TOKEN_HERE"
```

---

## 🎯 **STEP 8: Customer Requests Subscription**

**Purpose:** Customer requests to subscribe to a plan (requires admin approval)

```bash
curl -X POST http://localhost:8080/api/v1/customer/subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN" \
  -d '{
    "sku": "premium-plan"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Subscription request submitted successfully",
  "subscription": {
    "id": 1,
    "status": "requested",
    "requested_at": "2024-12-08T10:00:00Z"
  }
}
```

**⚠️ IMPORTANT:** Note the subscription `id` - you'll need it in Step 9!

**Status:** Subscription is now in `"requested"` status, waiting for admin approval.

---

## 🎯 **STEP 9: Admin Approves Subscription Request**

**Purpose:** Admin approves the customer's subscription request

```bash
curl -X POST http://localhost:8080/api/v1/admin/subscriptions/1/approve \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Replace `1` with the subscription ID from Step 8!**

**Expected Response:**
```json
{
  "success": true,
  "message": "Subscription approved successfully"
}
```

**Status:** Subscription is now in `"approved"` status, but not yet active.

---

## 🎯 **STEP 10: Admin Assigns Subscription** (Makes it Active)

**Purpose:** Admin assigns the subscription to customer (makes it active)

```bash
curl -X POST http://localhost:8080/api/v1/admin/customers/1/assign-subscription \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "pack_id": 1
  }'
```

**Replace `1` in the URL with customer ID, and `pack_id: 1` with the pack ID!**

**Expected Response:**
```json
{
  "success": true,
  "message": "Subscription assigned successfully"
}
```

**Status:** Subscription is now `"active"` and customer has access!

---

## 🎯 **STEP 11: Customer Views Their Active Subscription**

**Purpose:** Customer checks their current active subscription

```bash
curl -X GET http://localhost:8080/api/v1/customer/subscription \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "subscription": {
    "id": 1,
    "pack": {
      "name": "Premium Plan",
      "sku": "premium-plan",
      "price": 29.99,
      "validity_months": 12
    },
    "status": "active",
    "assigned_at": "2024-12-08T10:00:00Z",
    "expires_at": "2025-12-08T10:00:00Z",
    "is_valid": true
  }
}
```

---

## 🎯 **STEP 12: Customer Views Subscription History**

**Purpose:** Customer sees all their past and current subscriptions

```bash
curl -X GET "http://localhost:8080/api/v1/customer/subscription-history?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "history": [
    {
      "id": 1,
      "pack_name": "Premium Plan",
      "status": "active",
      "assigned_at": "2024-12-08T10:00:00Z",
      "expires_at": "2025-12-08T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

---

## 🎯 **STEP 13: Admin Views Dashboard**

**Purpose:** Admin sees system overview and statistics

```bash
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "total_customers": 1,
    "active_subscriptions": 1,
    "pending_requests": 0,
    "total_revenue": 29.99,
    "recent_activities": [...]
  }
}
```

---

## 🎯 **STEP 14: Admin Lists All Customers**

**Purpose:** Admin sees all registered customers

```bash
curl -X GET "http://localhost:8080/api/v1/admin/customers?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## 🎯 **STEP 15: Admin Lists All Subscriptions**

**Purpose:** Admin sees all subscriptions with their statuses

```bash
curl -X GET "http://localhost:8080/api/v1/admin/subscriptions?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

**Filter by status:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/subscriptions?page=1&limit=10&status=active" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

---

## 🎯 **STEP 16: Customer Deactivates Subscription**

**Purpose:** Customer cancels their active subscription

```bash
curl -X DELETE http://localhost:8080/api/v1/customer/subscription \
  -H "Authorization: Bearer YOUR_CUSTOMER_TOKEN"
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Subscription deactivated successfully",
  "deactivated_at": "2024-12-08T10:00:00Z"
}
```

**Status:** Subscription is now `"inactive"`

---

## 🔐 **SDK Workflow (Alternative Authentication)**

### **STEP 17: SDK Login** (Get API Key)

**Purpose:** Get API key for SDK/mobile app usage

```bash
curl -X POST http://localhost:8080/sdk/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "api_key": "sk-sdk-1234567890abcdef",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "name": "John Doe",
  "phone": "+1234567890",
  "expires_in": 3600
}
```

**⚠️ IMPORTANT:** Copy the `api_key` for SDK endpoints!

**Save to variable:**
```bash
API_KEY="sk-sdk-1234567890abcdef"
```

---

### **STEP 18: Get Subscription via SDK**

**Purpose:** Mobile app gets customer's subscription using API key

```bash
curl -X GET http://localhost:8080/sdk/v1/subscription \
  -H "X-API-Key: YOUR_API_KEY_HERE"
```

---

### **STEP 19: Request Subscription via SDK**

**Purpose:** Mobile app requests subscription using API key

```bash
curl -X POST http://localhost:8080/sdk/v1/subscription \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY_HERE" \
  -d '{
    "pack_sku": "premium-plan"
  }'
```

---

### **STEP 20: Get Subscription History via SDK**

**Purpose:** Mobile app gets subscription history

```bash
curl -X GET "http://localhost:8080/sdk/v1/subscription-history?page=1&limit=10" \
  -H "X-API-Key: YOUR_API_KEY_HERE"
```

---

## 📋 **Complete Workflow Summary**

### **Basic Flow:**
1. ✅ Admin Login → Get Admin Token
2. ✅ Create Subscription Pack → Get Pack SKU
3. ✅ Create Customer → Get Customer Email
4. ✅ Customer Login → Get Customer Token
5. ✅ Customer Requests Subscription → Get Subscription ID
6. ✅ Admin Approves Subscription
7. ✅ Admin Assigns Subscription → Makes it Active
8. ✅ Customer Views Subscription → Verify it's Active

### **Admin Management Flow:**
- View Dashboard
- List Customers
- List Subscriptions
- Create/Update/Delete Packs
- Create/Update/Delete Customers

### **Customer Self-Service Flow:**
- View Current Subscription
- Request New Subscription
- View Subscription History
- Deactivate Subscription

### **SDK Flow:**
- SDK Login → Get API Key
- Use API Key for all SDK endpoints
- Same operations as customer endpoints but with API key auth

---

## 💡 **Pro Tips**

### **1. Use Variables in Bash:**
```bash
# Set tokens
ADMIN_TOKEN="your_admin_token"
CUSTOMER_TOKEN="your_customer_token"
API_KEY="your_api_key"

# Use in requests
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

### **2. Pretty Print JSON:**
```bash
# Install jq first: brew install jq (macOS)
curl ... | jq .
```

### **3. Save Responses:**
```bash
curl ... -o response.json
```

### **4. View Full Request/Response:**
```bash
curl -v ...  # Shows headers and full details
```

---

## 🚨 **Common Issues**

### **401 Unauthorized:**
- Token expired (24 hours)
- Wrong token
- Missing "Bearer " prefix

### **404 Not Found:**
- Wrong ID in URL
- Endpoint doesn't exist
- Resource was deleted

### **400 Bad Request:**
- Missing required fields
- Invalid data format
- Business rule violation (e.g., customer already has active subscription)

---

## ✅ **Quick Test Checklist**

- [ ] Admin can login
- [ ] Admin can create subscription pack
- [ ] Admin can create customer
- [ ] Customer can login
- [ ] Customer can request subscription
- [ ] Admin can approve subscription
- [ ] Admin can assign subscription
- [ ] Customer can view active subscription
- [ ] SDK login works
- [ ] SDK endpoints work with API key

---

**Happy Testing! 🚀**







