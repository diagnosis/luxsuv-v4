# LuxSUV Backend API

A Go-based REST API for a luxury ride-sharing service with user authentication, password reset, and email notifications powered by MailerSend.

## 🚀 Features

- User registration and authentication with JWT tokens
- Password reset functionality with beautiful email notifications
- Role-based access control (rider, driver, admin)
- Admin user management with full CRUD operations
- Professional email service integration with MailerSend
- Comprehensive logging and monitoring
- Database migrations with Goose
- Rate limiting and CORS support
- Beautiful HTML email templates

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- MailerSend account (free tier includes 12,000 emails/month)

## 🛠️ Setup

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd luxsuv-v4
go mod download
```

### 2. Environment Configuration

Create/update your `.env` file:

```bash
# Database
DATABASE_URL=postgresql://username:password@host:port/database?sslmode=require

# JWT Secret (generate a secure 32+ character string)
JWT_SECRET=your-secure-jwt-secret-key-here

# Server
PORT=8080
ENVIRONMENT=development

# Email Configuration (MailerSend)
MAILERSEND_API_KEY=mlsn.your-api-key-here
MAILERSEND_FROM_EMAIL=noreply@yourdomain.com
MAILERSEND_FROM_NAME=LuxSUV Support
```

### 3. MailerSend Setup

1. **Sign up for MailerSend**: Go to [MailerSend](https://www.mailersend.com/)
2. **Get API Key**: Dashboard → API Tokens → Create new token
3. **Verify Domain**: Add and verify your sending domain
4. **Update .env**: Add your API key and verified email address

### 4. Database Setup

The application automatically runs migrations on startup. Ensure your PostgreSQL database is accessible.

### 5. Run the Application

```bash
go run cmd/server/main.go
```

## 📚 Complete API Documentation

### 🌐 Public Endpoints

#### 1. Health Check
```bash
# Check if server is running
curl -X GET http://localhost:8080/health
```
**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-10T12:00:00Z"
}
```

#### 2. User Registration
```bash
# Register a new rider
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123",
    "role": "rider"
  }'

# Register a driver
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "janedoe",
    "email": "jane@example.com",
    "password": "driverpass123",
    "role": "driver"
  }'

# Register an admin
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@luxsuv.com",
    "password": "adminpass123",
    "role": "admin"
  }'
```

**Success Response:**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "rider",
    "is_admin": false,
    "created_at": "2025-01-10T12:00:00Z"
  }
}
```

#### 3. User Login
```bash
# Login with email and password
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```
**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "rider",
    "is_admin": false,
    "created_at": "2025-01-10T12:00:00Z"
  }
}
```

#### 4. Password Reset Request
```bash
# Request password reset (sends beautiful email)
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

**Response:**
```json
{
  "message": "if the email exists, a password reset link has been sent"
}
```

#### 5. Reset Password with Token
```bash
# Reset password using token from email
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "reset_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "new_password": "mynewpassword123"
  }'
```

### 🔐 Protected Endpoints (Require Authentication)

#### 6. Get Current User Profile
```bash
# Get your own profile information
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 7. Change Password (Authenticated)
```bash
# Change your password while logged in
curl -X PUT http://localhost:8080/users/me/password \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "oldpassword123",
    "new_password": "newpassword123"
  }'
```

### 👑 Admin Endpoints (Require Admin Role)

#### 8. List All Users (with Pagination)
```bash
# Get all users (default: page 1, limit 10)
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"

# Get users with pagination
curl -X GET "http://localhost:8080/admin/users?page=2&limit=5" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"
```
**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "role": "rider",
      "is_admin": false,
      "created_at": "2025-01-10T12:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_count": 50,
    "limit": 10
  }
}
```

#### 9. Get User by Email
```bash
# Find user by email address
curl -X GET "http://localhost:8080/admin/users/by-email?email=john@example.com" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"
```

#### 10. Get User by ID
```bash
# Get specific user by ID
curl -X GET http://localhost:8080/admin/users/123 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"
```

#### 11. Update User Role
```bash
# Promote user to driver
curl -X PUT http://localhost:8080/admin/users/123/role \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "driver"
  }'

# Promote user to admin
curl -X PUT http://localhost:8080/admin/users/123/role \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin"
  }'

# Demote user to rider
curl -X PUT http://localhost:8080/admin/users/123/role \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "rider"
  }'
```

#### 12. Delete User
```bash
# Delete a user account (admin only)
curl -X DELETE http://localhost:8080/admin/users/123 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN_HERE"
```

## 🔄 Complete Workflow Examples

### Registration → Login → Access Protected Resource
```bash
# Step 1: Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "password123",
    "role": "rider"
  }'

# Step 2: Login (save the token)
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123"
  }' | jq -r '.token')

# Step 3: Use token to access protected endpoint
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### Complete Password Reset Flow
```bash
# Step 1: Request reset (user receives beautiful email)
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'

# Step 2: User clicks link in email or uses token directly
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "reset_token": "TOKEN_FROM_EMAIL",
    "new_password": "mynewpassword123"
  }'

# Step 3: Login with new password
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "mynewpassword123"
  }'
```

### Admin User Management Flow
```bash
# Step 1: Admin login
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@luxsuv.com",
    "password": "adminpass123"
  }' | jq -r '.token')

# Step 2: List all users
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Step 3: Find specific user by email
curl -X GET "http://localhost:8080/admin/users/by-email?email=john@example.com" \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Step 4: Update user role
curl -X PUT http://localhost:8080/admin/users/123/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "driver"}'
```

## 📧 Email Features

### Beautiful Email Templates
- **Password Reset**: Professional design with gradient headers, security notices, and clear call-to-action buttons
- **Welcome Email**: Engaging onboarding email with feature highlights and modern styling
- **Responsive Design**: Looks great on all devices
- **Security Features**: Expiration notices, IP logging, security warnings

### MailerSend Benefits
- **Reliable Delivery**: Better inbox placement than traditional SMTP
- **Analytics**: Track opens, clicks, bounces, and delivery status
- **Professional Templates**: Beautiful, responsive HTML emails
- **Webhooks**: Real-time delivery notifications
- **Suppression Lists**: Automatic bounce and complaint handling

## 🔒 Security Features

- **Password Security**: bcrypt hashing with salt
- **JWT Authentication**: Secure token-based authentication
- **Rate Limiting**: 
  - General endpoints: 5 req/sec, burst of 10
  - Auth endpoints: 2 req/sec, burst of 5
- **CORS Protection**: Configurable for development and production
- **Input Validation**: Comprehensive validation and sanitization
- **SQL Injection Protection**: Parameterized queries throughout
- **Password Requirements**: Minimum 8 chars, letters + numbers required

## 📊 API Response Formats

### Success Responses
```json
{
  "message": "operation successful",
  "data": { ... }
}
```

### Error Responses
```json
{
  "error": "descriptive error message"
}
```

### HTTP Status Codes
- `200`: Success
- `201`: Created (registration)
- `400`: Bad Request (validation errors)
- `401`: Unauthorized (invalid/missing token)
- `403`: Forbidden (insufficient permissions)
- `404`: Not Found
- `429`: Too Many Requests (rate limited)
- `500`: Internal Server Error

## 🚀 Development

### Running with Air (Hot Reload)
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Database Migrations
```bash
# Migrations run automatically on startup
# To create new migrations:
goose -dir migrations create migration_name sql
```

### Testing Endpoints
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test with verbose output
curl -v -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123","role":"rider"}'
```

## 🌐 Production Deployment

### Environment Setup
1. Set `ENVIRONMENT=production`
2. Use strong JWT secret (32+ characters)
3. Configure proper CORS origins
4. Set up SSL/TLS for HTTPS
5. Use environment variables instead of `.env` file

### MailerSend Production Setup
1. Verify your sending domain
2. Set up DNS records (SPF, DKIM, DMARC)
3. Configure webhooks for delivery tracking
4. Set up suppression lists
5. Monitor sender reputation

### Database Production Setup
1. Use connection pooling
2. Set up read replicas if needed
3. Configure backup strategy
4. Monitor performance metrics

## 🔧 Troubleshooting

### Email Issues
- **Not sending**: Check API key and domain verification
- **Spam folder**: Verify DNS records and sender reputation
- **Rate limits**: Free plan allows 12,000 emails/month

### Authentication Issues
- **Invalid token**: Check JWT secret and token expiration (24h)
- **Permissions**: Verify user role and admin status
- **Headers**: Ensure `Authorization: Bearer <token>` format

### Database Issues
- **Connection**: Verify DATABASE_URL and network access
- **Migrations**: Check migration files and database permissions
- **Performance**: Monitor connection pool and query performance

## 📝 User Roles & Permissions

### Rider (Default)
- Register and login
- View own profile
- Change own password
- Book rides (future feature)

### Driver
- All rider permissions
- Accept rides (future feature)
- Driver-specific features (future)

### Admin
- All user permissions
- List all users
- View any user profile
- Update user roles
- Delete users
- Access admin endpoints

## 🎯 API Rate Limits

| Endpoint Type | Rate Limit | Burst Limit |
|---------------|------------|-------------|
| General | 5 req/sec | 10 requests |
| Authentication | 2 req/sec | 5 requests |
| Admin | 5 req/sec | 10 requests |

## 📈 Monitoring & Logging

- **Comprehensive Logging**: All requests, errors, and operations logged
- **Log Levels**: INFO, WARN, ERROR with timestamps
- **File + Console**: Logs written to both `app.log` and console
- **Request Tracking**: Full request/response logging with timing

## 🔮 Future Enhancements

- [ ] Ride booking and management system
- [ ] Real-time ride tracking with WebSockets
- [ ] Payment integration (Stripe)
- [ ] Driver location tracking
- [ ] Push notifications
- [ ] Mobile app API support
- [ ] Advanced analytics and reporting

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**LuxSUV Backend API** - Built with ❤️ using Go, Echo, PostgreSQL, and MailerSend