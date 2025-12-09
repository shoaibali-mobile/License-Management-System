# How to Run the API Locally and Test in Postman

## Step 1: Start the API Server

### Option A: Run directly with Go
```bash
cd /Users/shoaibali/Documents/License-MNM/backend
go run main.go
```

### Option B: Build and run
```bash
cd /Users/shoaibali/Documents/License-MNM/backend
go build -o license-mnm
./license-mnm
```

The server will start on **http://localhost:8080**

## Step 2: Create Admin User (First Time Only)

Before testing, create an admin user:

```bash
cd /Users/shoaibali/Documents/License-MNM/backend
go run cmd/seed/main.go
```

This creates:
- **Email**: admin@example.com
- **Password**: admin123

## Step 3: Test in Postman

### Import OpenAPI Specification

1. Open Postman
2. Click **Import** (top left)
3. Select **File** tab
4. Choose `openapi.yaml` from the project root
5. Click **Import**

This will create a collection with all endpoints.

### Set Up Environment Variables (Optional but Recommended)

1. Click **Environments** (left sidebar)
2. Click **+** to create new environment
3. Name it "License MNM Local"
4. Add variables:
   - `base_url`: `http://localhost:8080`
   - `jwt_token`: (leave empty, will be set after login)
   - `api_key`: (leave empty, will be set after SDK login)
5. Click **Save**

### Testing Flow

#### 1. Admin Login (Get JWT Token)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/admin/login`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "email": "admin@example.com",
    "password": "admin123"
  }
  ```

**Response:** You'll get a JWT token. Copy it!

**Expected Response:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "admin@example.com",
  "expires_in": 3600
}
```

#### 2. Create a Subscription Pack (Admin)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/admin/subscription-packs`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
  ```
- Body (raw JSON):
  ```json
  {
    "name": "Premium Plan",
    "description": "Full access to all features",
    "sku": "premium-plan",
    "price": 29.99,
    "validity_months": 12
  }
  ```

#### 3. Create a Customer (Admin)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/admin/customers`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer YOUR_JWT_TOKEN_HERE
  ```
- Body (raw JSON):
  ```json
  {
    "name": "John Doe",
    "email": "customer@example.com",
    "phone": "+1234567890"
  }
  ```

**Note:** This creates a customer with default password `password123`

#### 4. Customer Signup (Alternative)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/customer/signup`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "name": "Jane Smith",
    "email": "jane@example.com",
    "password": "password123",
    "phone": "+1987654321"
  }
  ```

#### 5. Customer Login (Get JWT Token)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/customer/login`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "email": "customer@example.com",
    "password": "password123"
  }
  ```

#### 6. Request Subscription (Customer)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/customer/subscription`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer CUSTOMER_JWT_TOKEN_HERE
  ```
- Body (raw JSON):
  ```json
  {
    "sku": "premium-plan"
  }
  ```

#### 7. Approve Subscription (Admin)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/admin/subscriptions/1/approve`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer ADMIN_JWT_TOKEN_HERE
  ```

#### 8. Assign Subscription Directly (Admin)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/v1/admin/customers/1/assign-subscription`
- Headers:
  ```
  Content-Type: application/json
  Authorization: Bearer ADMIN_JWT_TOKEN_HERE
  ```
- Body (raw JSON):
  ```json
  {
    "pack_id": 1
  }
  ```

#### 9. Get Customer Subscription (Customer)

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/customer/subscription`
- Headers:
  ```
  Authorization: Bearer CUSTOMER_JWT_TOKEN_HERE
  ```

#### 10. Get Admin Dashboard

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/v1/admin/dashboard`
- Headers:
  ```
  Authorization: Bearer ADMIN_JWT_TOKEN_HERE
  ```

### SDK Endpoints Testing

#### 1. SDK Login (Get API Key)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/sdk/auth/login`
- Headers:
  ```
  Content-Type: application/json
  ```
- Body (raw JSON):
  ```json
  {
    "email": "customer@example.com",
    "password": "password123"
  }
  ```

**Response:** You'll get an `api_key`. Copy it!

#### 2. Get Current Subscription (SDK)

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/sdk/v1/subscription`
- Headers:
  ```
  X-API-Key: YOUR_API_KEY_HERE
  ```

#### 3. Request Subscription (SDK)

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/sdk/v1/subscription`
- Headers:
  ```
  Content-Type: application/json
  X-API-Key: YOUR_API_KEY_HERE
  ```
- Body (raw JSON):
  ```json
  {
    "pack_sku": "premium-plan"
  }
  ```

## Quick Test Sequence

1. **Start server**: `go run main.go`
2. **Create admin**: `go run cmd/seed/main.go`
3. **Admin login** → Get JWT token
4. **Create subscription pack** → Get pack ID/SKU
5. **Create customer** → Get customer ID
6. **Customer login** → Get customer JWT token
7. **Request subscription** (customer) → Subscription in "requested" status
8. **Approve subscription** (admin) → Status changes to "approved"
9. **Assign subscription** (admin) → Status changes to "active"
10. **Get subscription** (customer) → View active subscription

## Troubleshooting

### Port Already in Use
```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process or use a different port
# Edit main.go: r.Run(":8081")
```

### Database Errors
- Make sure `license_mnm.db` file exists in the backend directory
- Delete `license_mnm.db` to reset the database (all data will be lost)

### Authentication Errors
- Make sure you're using `Bearer ` prefix before the token
- Check that the token hasn't expired (24 hours)
- Verify the email/password are correct

### CORS Errors
- The API has CORS enabled for all origins
- If you still see errors, check the browser console

## Common Endpoints Summary

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/admin/login` | POST | None | Admin login |
| `/api/customer/login` | POST | None | Customer login |
| `/api/customer/signup` | POST | None | Customer registration |
| `/api/v1/admin/dashboard` | GET | JWT (Admin) | Admin dashboard |
| `/api/v1/admin/customers` | GET/POST | JWT (Admin) | List/Create customers |
| `/api/v1/admin/subscription-packs` | GET/POST | JWT (Admin) | List/Create packs |
| `/api/v1/customer/subscription` | GET/POST/DELETE | JWT (Customer) | Get/Request/Deactivate subscription |
| `/sdk/auth/login` | POST | None | SDK login (get API key) |
| `/sdk/v1/subscription` | GET/POST/DELETE | API Key | SDK subscription operations |







