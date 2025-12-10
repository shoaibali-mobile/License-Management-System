# License Management System

A subscription and license management API built with Go/Gin framework. Supports both web frontend (JWT authentication) and mobile SDK (API key authentication) integrations.

## Table of Contents

1. [Core Components](#core-components)
2. [Installation & Setup](#installation--setup)
3. [API Endpoints](#api-endpoints)
4. [Authentication](#authentication)
5. [Subscription Status](#subscription-status)
6. [Testing](#testing)
7. [Configuration](#configuration)
8. [Troubleshooting](#troubleshooting)

---

## Core Components

### User Management
- **Admin**: Full system access, manages customers, subscriptions, and packs
- **Customer**: Self-service registration, login, subscription requests, and deactivation

### Subscription Pack Management
- Create, list, update, and delete subscription packs
- Attributes: Name, Description, SKU, Price, Validity (1-12 months)

### Customer Management
- CRUD operations for customer profiles
- Attributes: Name, Email, Phone, Subscription History

### Subscription Lifecycle
- **Status Flow**: `requested` → `approved` → `active` → `inactive`/`expired`
- **Business Rules**: Only one active subscription per customer
- **Operations**: Request, Approve, Assign, Deactivate, Unassign

### SDK Integration
- **Platforms**: Android, iOS, JavaScript
- **Authentication**: API Key (persistent, no expiration)
- **Operations**: Get subscription, Request subscription, Deactivate, View history

## Installation & Setup

### Prerequisites
- Go 1.21+
- Git
- SQLite (included)

### Setup Steps
1. Clone repository and navigate to backend directory
2. Install dependencies: `go mod download`
3. Create admin user: `go run cmd/seed/main.go`
   - Email: `admin@example.com`
   - Password: `admin123`
4. Start server: `go run main.go`
   - Server runs on `http://0.0.0.0:8080` (accessible from network)
   - Local: `http://localhost:8080`
   - Mobile: `http://YOUR_IP_ADDRESS:8080`

## API Endpoints

### Public Endpoints (No Authentication)
- `POST /api/admin/login` - Admin login (returns JWT)
- `POST /api/customer/login` - Customer login (returns JWT)
- `POST /api/customer/signup` - Customer registration
- `POST /sdk/auth/login` - SDK login (returns API key)

### Admin Endpoints (JWT Required)
- `GET /api/v1/admin/dashboard` - Dashboard statistics
- `GET /api/v1/admin/customers` - List all customers
- `POST /api/v1/admin/customers` - Create customer
- `GET /api/v1/admin/customers/:id` - Get customer details
- `PUT /api/v1/admin/customers/:id` - Update customer
- `DELETE /api/v1/admin/customers/:id` - Delete customer
- `GET /api/v1/admin/subscription-packs` - List packs
- `POST /api/v1/admin/subscription-packs` - Create pack
- `PUT /api/v1/admin/subscription-packs/:id` - Update pack
- `DELETE /api/v1/admin/subscription-packs/:id` - Delete pack
- `GET /api/v1/admin/subscriptions` - List all subscriptions
- `POST /api/v1/admin/subscriptions/:id/approve` - Approve subscription
- `POST /api/v1/admin/customers/:id/assign-subscription` - Assign subscription
- `DELETE /api/v1/admin/customers/:id/subscription/:id` - Unassign subscription

### Customer Endpoints (JWT Required)
- `GET /api/v1/customer/subscription` - Get current subscription
- `POST /api/v1/customer/subscription` - Request subscription
- `DELETE /api/v1/customer/subscription` - Deactivate subscription
- `GET /api/v1/customer/subscription-history` - Get history

### SDK Endpoints (API Key Required)
- `GET /sdk/v1/subscription` - Get current subscription
- `POST /sdk/v1/subscription` - Request subscription
- `DELETE /sdk/v1/subscription` - Deactivate subscription
- `GET /sdk/v1/subscription-history` - Get history

## Authentication

### JWT Authentication (Web Frontend)
- Used for admin and customer web endpoints
- Token expires in 24 hours
- Header: `Authorization: Bearer <token>`
- Get token via login endpoints

### API Key Authentication (Mobile SDK)
- Used for SDK endpoints only
- Never expires (persistent)
- Header: `X-API-Key: <api_key>`
- Get API key via `/sdk/auth/login` endpoint

## Subscription Status

### Status Definitions
| Status | Meaning | Who Can Change | When It Happens |
|--------|---------|----------------|-----------------|
| `requested` | Customer requested subscription | Customer → Admin | Customer submits request |
| `approved` | Admin approved the request | Admin → Admin | Admin approves subscription |
| `active` | Subscription is active and valid | Admin → Customer/Admin | Admin assigns subscription |
| `inactive` | Subscription deactivated | Customer/Admin | Customer deactivates or admin unassigns |
| `expired` | Validity period ended | System | `expires_at` date passed |

### Status Transitions
**Normal Flow:** `requested` → `approved` → `active` → `inactive`/`expired`

**Who Can Deactivate:**
- **Customers**: Can only deactivate `active` subscriptions
- **Admin**: Can unassign any subscription (any status)

**Business Rules:**
- Only one `active` subscription per customer at a time
- Cannot request new subscription while `active` exists
- `approved` subscriptions cannot be deactivated by customers (not active yet)
- `requested` subscriptions cannot be deactivated (waiting for approval)

## Testing

### Manual Testing
See `backend/CURL_COMMANDS.md` for complete curl command reference.

### Automated Testing
Run the test script: `./quick_test.sh` in the backend directory.

### Testing in Postman
1. Import `openapi.yaml` into Postman
2. Set up environment variables: `base_url`, `admin_token`, `customer_token`, `api_key`

## Configuration

### Environment Variables
Create a `.env` file in the backend directory with:
- `PORT=8080`
- `HOST=0.0.0.0`
- `DB_TYPE=sqlite`
- `DB_PATH=license_mnm.db`
- `JWT_SECRET=your-secret-key-change-in-production`
- `CORS_ALLOW_ORIGINS=*`

### Database Configuration
- **SQLite (Default)**: Database file `license_mnm.db`, no additional configuration needed
- **PostgreSQL (Production)**: Update database connection in `database/database.go`

## Deployment

### Production Considerations
1. Change JWT Secret with a strong, randomly generated secret
2. Use PostgreSQL instead of SQLite
3. Enable HTTPS using reverse proxy (nginx) with SSL certificates
4. Configure CORS to allow only trusted domains
5. Use environment variables for sensitive data
6. Implement proper logging

## Troubleshooting

### Server Won't Start
- Check if port 8080 is already in use
- Verify Go version (requires 1.21+)
- Check database file permissions

### Authentication Fails
- Verify JWT token hasn't expired (24 hours)
- Check API key is correct for SDK endpoints
- Ensure Authorization header format: `Bearer TOKEN`

### Database Errors
- Check database file exists: `license_mnm.db`
- Verify database permissions
- Run seed script to create admin user

### Network Issues
- Ensure server is listening on `0.0.0.0:8080` (not just localhost)
- Check firewall settings
- Verify client and server are on same network

## Base URL Configuration

### Development
- **Local**: `http://localhost:8080`
- **Network**: `http://YOUR_IP_ADDRESS:8080` (e.g., `http://192.168.2.202:8080`)
- **Android Emulator**: `http://10.0.2.2:8080`

### Production
- **HTTPS**: `https://your-domain.com`
- **Port**: Usually 443 (HTTPS) or 8080 (if behind reverse proxy)

## Common Workflows

### Customer Workflow
1. Signup → `POST /api/customer/signup`
2. SDK Login → `POST /sdk/auth/login` → Get API key
3. Check Subscription → `GET /sdk/v1/subscription`
4. Request Subscription → `POST /sdk/v1/subscription` (if no active)
5. Wait for Admin → Admin approves and assigns
6. View Subscription → `GET /sdk/v1/subscription` (now active)
7. Deactivate → `DELETE /sdk/v1/subscription` (if needed)

### Admin Workflow
1. Login → `POST /api/admin/login` → Get JWT token
2. View Dashboard → `GET /api/v1/admin/dashboard`
3. Create Pack → `POST /api/v1/admin/subscription-packs`
4. List Subscriptions → `GET /api/v1/admin/subscriptions?status=requested`
5. Approve → `POST /api/v1/admin/subscriptions/:id/approve`
6. Assign → `POST /api/v1/admin/customers/:id/assign-subscription`

## Additional Resources

- **API Documentation**: See `API_DOCUMENTATION.md` for Android/mobile integration
- **cURL Commands**: See `backend/CURL_COMMANDS.md` for all API examples
- **OpenAPI Specification**: See `openapi.yaml` for complete API specification

---
