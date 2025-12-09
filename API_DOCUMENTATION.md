# API Documentation for Android App Development

Complete API reference for License Management System Android SDK integration.

**Base URL:** `http://localhost:8080` (Development)  


---

## 📋 Table of Contents

1. [Authentication](#authentication)
2. [Subscription Management](#subscription-management)
3. [Data Models](#data-models)
4. [Error Handling](#error-handling)
5. [Example Requests](#example-requests)

---

## 🔐 Authentication

### SDK Login

Authenticate customer and get API key for SDK endpoints.

**Endpoint:** `POST /sdk/auth/login`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "customer@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "api_key": "sk-sdk-1f7ae96e807f3bfef29afc113756c496",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "name": "John Doe",
  "phone": "+1234567890",
  "expires_in": 3600
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Invalid credentials"
}
```

**Important Notes:**
- Store the `api_key` securely (EncryptedSharedPreferences)
- API key does NOT expire - use it for all SDK requests
- JWT token is optional (expires in 24 hours)
- Use `api_key` in `X-API-Key` header for all SDK endpoints

---

## 📱 Subscription Management

### Get Current Subscription

Get customer's active subscription details.

**Endpoint:** `GET /sdk/v1/subscription`

**Headers:**
```
X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496
```

**Response (200 OK):**
```json
{
  "success": true,
  "subscription": {
    "id": 4,
    "pack_name": "Premium Plan",
    "pack_sku": "premium-plan",
    "price": 29.99,
    "status": "active",
    "assigned_at": "2024-12-08T11:18:23.791469+05:30",
    "expires_at": "2025-12-08T11:18:23.791469+05:30",
    "is_valid": true
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "No active subscription found"
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Invalid API key"
}
```

---

### Request New Subscription

Request a subscription for a specific pack.

**Endpoint:** `POST /sdk/v1/subscription`

**Headers:**
```
Content-Type: application/json
X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496
```

**Request Body:**
```json
{
  "pack_sku": "premium-plan"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Subscription request submitted successfully",
  "subscription": {
    "id": 5,
    "status": "requested",
    "requested_at": "2024-12-08T12:00:00.000000+05:30"
  }
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Customer already has an active subscription"
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Subscription pack not found"
}
```

**Status Flow:**
- `requested` → Customer requested, waiting for admin approval
- `approved` → Admin approved, waiting for assignment
- `active` → Subscription is active and valid
- `inactive` → Subscription deactivated
- `expired` → Subscription validity ended

---

### Deactivate Current Subscription

Deactivate customer's active subscription.

**Endpoint:** `DELETE /sdk/v1/subscription`

**Headers:**
```
X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Subscription deactivated successfully",
  "deactivated_at": "2024-12-08T12:00:00.000000+05:30"
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "No active subscription found"
}
```

---

### Get Subscription History

Get paginated list of customer's subscription history.

**Endpoint:** `GET /sdk/v1/subscription-history`

**Headers:**
```
X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496
```

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `limit` (optional, default: 10) - Items per page
- `sort` (optional, default: "desc") - Sort order: "asc" or "desc"

**Example Request:**
```
GET /sdk/v1/subscription-history?page=1&limit=10&sort=desc
```

**Response (200 OK):**
```json
{
  "success": true,
  "history": [
    {
      "id": 4,
      "pack_name": "Premium Plan",
      "status": "active",
      "assigned_at": "2024-12-08T11:18:23.791469+05:30",
      "expires_at": "2025-12-08T11:18:23.791469+05:30"
    },
    {
      "id": 3,
      "pack_name": "Premium Plan",
      "status": "expired",
      "assigned_at": "2023-12-08T11:18:23.791469+05:30",
      "expires_at": "2024-12-08T11:18:23.791469+05:30"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2
  }
}
```

---

## 📊 Data Models

### LoginRequest
```json
{
  "email": "string (required, email format)",
  "password": "string (required, min 6 characters)"
}
```

### LoginResponse
```json
{
  "success": "boolean",
  "api_key": "string (required for SDK)",
  "token": "string (optional JWT token)",
  "name": "string (customer name)",
  "phone": "string (customer phone)",
  "expires_in": "integer (token expiration in seconds)"
}
```

### Subscription
```json
{
  "id": "integer",
  "pack_name": "string",
  "pack_sku": "string",
  "price": "number (float)",
  "status": "string (requested|approved|active|inactive|expired)",
  "assigned_at": "string (ISO 8601 datetime)",
  "expires_at": "string (ISO 8601 datetime)",
  "is_valid": "boolean"
}
```

### SubscriptionRequest
```json
{
  "pack_sku": "string (required)"
}
```

### SubscriptionHistoryItem
```json
{
  "id": "integer",
  "pack_name": "string",
  "status": "string",
  "assigned_at": "string (ISO 8601 datetime)",
  "expires_at": "string (ISO 8601 datetime)"
}
```

### Pagination
```json
{
  "page": "integer",
  "limit": "integer",
  "total": "integer"
}
```

---

## ⚠️ Error Handling

### HTTP Status Codes

| Status Code | Meaning | Description |
|------------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request data or business rule violation |
| 401 | Unauthorized | Invalid or missing API key |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error |

### Error Response Format

All error responses follow this format:

```json
{
  "success": false,
  "message": "Error description"
}
```

### Common Error Messages

- `"Invalid credentials"` - Wrong email/password
- `"Invalid API key"` - API key not found or invalid
- `"No active subscription found"` - Customer has no active subscription
- `"Customer already has an active subscription"` - Cannot request new subscription while one is active
- `"Subscription pack not found"` - Invalid pack_sku
- `"X-API-Key header required"` - Missing API key header

---

## 🔧 Example Requests

### cURL Examples

#### 1. SDK Login
```bash
curl -X POST http://localhost:8080/sdk/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "password123"
  }'
```

#### 2. Get Current Subscription
```bash
curl -X GET http://localhost:8080/sdk/v1/subscription \
  -H "X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496"
```

#### 3. Request Subscription
```bash
curl -X POST http://localhost:8080/sdk/v1/subscription \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496" \
  -d '{
    "pack_sku": "premium-plan"
  }'
```

#### 4. Deactivate Subscription
```bash
curl -X DELETE http://localhost:8080/sdk/v1/subscription \
  -H "X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496"
```

#### 5. Get Subscription History
```bash
curl -X GET "http://localhost:8080/sdk/v1/subscription-history?page=1&limit=10&sort=desc" \
  -H "X-API-Key: sk-sdk-1f7ae96e807f3bfef29afc113756c496"
```

---

## 📱 Android Implementation Notes

### 1. Base URL Configuration

**Development (Emulator):**
```
http://10.0.2.2:8080
```

**Development (Physical Device):**
```
http://YOUR_COMPUTER_IP:8080
```

**Production:**
```
https://your-production-url.com
```

### 2. API Key Storage

- Store API key in **EncryptedSharedPreferences**
- Never store in plain SharedPreferences
- Use Android Keystore for additional security
- Clear API key on logout

### 3. Request Headers

All SDK endpoints require:
```
X-API-Key: <your_api_key>
```

Content-Type for POST requests:
```
Content-Type: application/json
```

### 4. Network Configuration

Add to `AndroidManifest.xml`:
```xml
<application
    android:networkSecurityConfig="@xml/network_security_config"
    ...>
</application>
```

Create `res/xml/network_security_config.xml`:
```xml
<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
    <!-- Development only -->
    <domain-config cleartextTrafficPermitted="true">
        <domain includeSubdomains="true">10.0.2.2</domain>
        <domain includeSubdomains="true">localhost</domain>
    </domain-config>
    
    <!-- Production: HTTPS only -->
    <base-config cleartextTrafficPermitted="false">
        <trust-anchors>
            <certificates src="system" />
        </trust-anchors>
    </base-config>
</network-security-config>
```

### 5. Error Handling

Always check:
1. HTTP status code
2. `success` field in response
3. `message` field for error details

Example:
```kotlin
if (response.isSuccessful && response.body()?.success == true) {
    // Handle success
} else {
    // Handle error
    val errorMessage = response.body()?.message ?: "Unknown error"
}
```

### 6. API Key Lifecycle

- **Get API Key**: Call `/sdk/auth/login` on first login
- **Store Securely**: Save to EncryptedSharedPreferences
- **Use for Requests**: Include in `X-API-Key` header
- **Check Validity**: API key doesn't expire, but check if user is logged in
- **Clear on Logout**: Remove from storage when user logs out

### 7. Subscription Status Handling

| Status | Meaning | User Action |
|--------|---------|-------------|
| `requested` | Waiting for admin approval | Show "Pending Approval" |
| `approved` | Approved but not active | Show "Approved - Waiting Activation" |
| `active` | Subscription is active | Show subscription details |
| `inactive` | Deactivated by user/admin | Show "Deactivated" |
| `expired` | Validity period ended | Show "Expired - Request New" |

---

## 🔄 Complete Flow Example

### 1. User Login
```
POST /sdk/auth/login
→ Get api_key
→ Store securely
```

### 2. Check Subscription
```
GET /sdk/v1/subscription
→ Show subscription if active
→ Show "No subscription" if none
```

### 3. Request Subscription (if needed)
```
POST /sdk/v1/subscription
→ Status: "requested"
→ Wait for admin approval
```

### 4. View History
```
GET /sdk/v1/subscription-history
→ Show all past subscriptions
```

### 5. Deactivate (if needed)
```
DELETE /sdk/v1/subscription
→ Status changes to "inactive"
```

---

## 📝 Quick Reference

### Endpoints Summary

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/sdk/auth/login` | None | Get API key |
| GET | `/sdk/v1/subscription` | API Key | Get current subscription |
| POST | `/sdk/v1/subscription` | API Key | Request new subscription |
| DELETE | `/sdk/v1/subscription` | API Key | Deactivate subscription |
| GET | `/sdk/v1/subscription-history` | API Key | Get subscription history |

### Required Headers

**All SDK endpoints (except login):**
```
X-API-Key: <api_key>
```

**POST requests:**
```
Content-Type: application/json
X-API-Key: <api_key>
```

---

## 🚨 Important Notes

1. **API Key Never Expires**: Once obtained, it works forever (until user changes password or admin revokes)
2. **No Token Refresh Needed**: Unlike JWT, API key doesn't need refresh
3. **Store Securely**: Always use EncryptedSharedPreferences
4. **Handle Offline**: Show cached data if API call fails
5. **Error Messages**: Always show user-friendly error messages
6. **Loading States**: Show loading indicators during API calls
7. **Network Errors**: Handle network timeouts and connection errors

---

## 📞 Support

For API issues or questions:
- Check error messages in response
- Verify API key is correct
- Ensure network connectivity
- Check server status

---

**Last Updated:** December 2024  
**API Version:** 1.0.0



