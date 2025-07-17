# LuxSUV Backend API

A Go-based REST API for a luxury ride-sharing service with user authentication, ride booking management, password reset, and email notifications powered by MailerSend.

## 🚀 Features

- User registration and authentication with JWT tokens
- Password reset functionality with beautiful email notifications
- Role-based access control (rider, driver, admin)
- Admin user management with full CRUD operations
- **Ride booking system** with create, update, and cancel functionality
- **Guest booking support** with secure email-based update links
- **24-hour advance booking** and cancellation policy
- **Driver booking acceptance** system
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

### 🚗 Ride Booking Endpoints

#### 1. Create Booking (Public - Authenticated or Guest)
```bash
# Create booking as authenticated user
curl -X POST http://localhost:8080/book-ride \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "your_name": "John Doe",
    "email": "john@example.com",
    "phone_number": "+1234567890",
    "ride_type": "Airport Transfer",
    "pickup_location": "123 Main St, City",
    "dropoff_location": "Airport Terminal 1",
    "date": "2025-01-15",
    "time": "14:30",
    "number_of_passengers": 2,
    "number_of_luggage": 3,
    "additional_notes": "Flight AA123 at 6 PM"
  }'

# Create booking as guest user (no authentication required)
curl -X POST http://localhost:8080/book-ride \
  -H "Content-Type: application/json" \
  -d '{
    "your_name": "Jane Smith",
    "email": "jane@example.com",
    "phone_number": "+1987654321",
    "ride_type": "City Tour",
    "pickup_location": "Hotel Downtown",
    "dropoff_location": "City Center",
    "date": "2025-01-20",
    "time": "10:00",
    "number_of_passengers": 4,
    "number_of_luggage": 0
  }'
```

#### 2. Get Bookings by Email (Public)
```bash
# Retrieve all bookings for an email address
curl -X GET http://localhost:8080/bookings/email/john%40example.com
```

#### 3. Generate Secure Update Link (Public)
```bash
# Request secure update link for guest users
curl -X POST http://localhost:8080/bookings/123/update-link \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```
**Response:** Beautiful email sent with secure update link

#### 4. Update Booking (Authenticated Users)
```bash
# Update booking as authenticated user
curl -X PUT http://localhost:8080/bookings/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pickup_location": "Updated Pickup Location",
    "date": "2025-01-16",
    "time": "15:00",
    "number_of_passengers": 3
  }'
```

#### 5. Update Booking with Secure Token (Guest Users)
```bash
# Update booking using secure token from email
curl -X PUT "http://localhost:8080/bookings/123/update?token=SECURE_TOKEN_FROM_EMAIL" \
  -H "Content-Type: application/json" \
  -d '{
    "your_name": "John Updated",
    "phone_number": "+1234567899",
    "additional_notes": "Updated flight information"
  }'
```

#### 6. Cancel Booking (Authenticated Users)
```bash
# Cancel booking as authenticated user
curl -X DELETE http://localhost:8080/bookings/123/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Plans changed"
  }'
```

#### 7. Cancel Booking with Secure Token (Guest Users)
```bash
# Cancel booking using secure token from email
curl -X DELETE "http://localhost:8080/bookings/123/cancel?token=SECURE_TOKEN_FROM_EMAIL" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Emergency came up"
  }'
```

#### 8. Get My Bookings (Authenticated Users)
```bash
# Get all bookings for authenticated user
curl -X GET http://localhost:8080/bookings/my \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 9. Accept Booking (Driver Only)
```bash
# Driver accepts a pending booking
curl -X PUT http://localhost:8080/driver/bookings/123/accept \
  -H "Authorization: Bearer DRIVER_JWT_TOKEN"
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

### Booking Management Emails
- **Booking Update Links**: Secure, time-limited links for guest users to update bookings
- **Professional Design**: Consistent branding with booking details and clear call-to-action buttons
- **Security Features**: 24-hour token expiry, booking-specific access control

### Beautiful Email Templates
- **Password Reset**: Professional design with gradient headers, security notices, and clear call-to-action buttons
- **Booking Updates**: Secure update links with booking details and professional styling
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
- **Booking Security**: Secure tokens for guest booking updates with email verification
- **24-Hour Policies**: Advance booking and cancellation restrictions
- **Rate Limiting**: 
  - General endpoints: 5 req/sec, burst of 10
  - Auth endpoints: 2 req/sec, burst of 5
- **CORS Protection**: Configurable for development and production
- **Input Validation**: Comprehensive validation and sanitization
- **SQL Injection Protection**: Parameterized queries throughout
- **Password Requirements**: Minimum 8 chars, letters + numbers required

## 📊 API Response Formats

### Booking Response Example
```json
{
  "id": 123,
  "user_id": 456,
  "your_name": "John Doe",
  "email": "john@example.com",
  "phone_number": "+1234567890",
  "ride_type": "Airport Transfer",
  "pickup_location": "123 Main St",
  "dropoff_location": "Airport Terminal 1",
  "date": "2025-01-15",
  "time": "14:30",
  "number_of_passengers": 2,
  "number_of_luggage": 3,
  "additional_notes": "Flight AA123",
  "book_status": "Pending",
  "ride_status": "Pending",
  "created_at": "2025-01-10T12:00:00Z",
  "updated_at": "2025-01-10T12:00:00Z"
}
```

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

### Complete Booking Workflow Examples

#### Guest User Booking Flow
```bash
# Step 1: Create booking as guest
curl -X POST http://localhost:8080/book-ride \
  -H "Content-Type: application/json" \
  -d '{
    "your_name": "Guest User",
    "email": "guest@example.com",
    "phone_number": "+1234567890",
    "ride_type": "City Tour",
    "pickup_location": "Hotel",
    "dropoff_location": "Mall",
    "date": "2025-01-20",
    "time": "10:00",
    "number_of_passengers": 2,
    "number_of_luggage": 1
  }'

# Step 2: Request update link
curl -X POST http://localhost:8080/bookings/123/update-link \
  -H "Content-Type: application/json" \
  -d '{"email": "guest@example.com"}'

# Step 3: Use token from email to update
curl -X PUT "http://localhost:8080/bookings/123/update?token=TOKEN_FROM_EMAIL" \
  -H "Content-Type: application/json" \
  -d '{"pickup_location": "Updated Hotel Location"}'

# Step 4: Cancel if needed (using same token)
curl -X DELETE "http://localhost:8080/bookings/123/cancel?token=TOKEN_FROM_EMAIL" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Plans changed"}'
```

#### Authenticated User Booking Flow
```bash
# Step 1: Login and get token
TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }' | jq -r '.token')

# Step 2: Create booking
curl -X POST http://localhost:8080/book-ride \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "your_name": "Authenticated User",
    "email": "user@example.com",
    "phone_number": "+1234567890",
    "ride_type": "Business Meeting",
    "pickup_location": "Office",
    "dropoff_location": "Conference Center",
    "date": "2025-01-18",
    "time": "09:00",
    "number_of_passengers": 1,
    "number_of_luggage": 1
  }'

# Step 3: View my bookings
curl -X GET http://localhost:8080/bookings/my \
  -H "Authorization: Bearer $TOKEN"

# Step 4: Update booking directly
curl -X PUT http://localhost:8080/bookings/123 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"time": "09:30"}'

# Step 5: Cancel if needed
curl -X DELETE http://localhost:8080/bookings/123/cancel \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Meeting rescheduled"}'
```

#### Driver Workflow
```bash
# Step 1: Login as driver
DRIVER_TOKEN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "driver@luxsuv.com",
    "password": "driverpass123"
  }' | jq -r '.token')

# Step 2: Accept pending booking
curl -X PUT http://localhost:8080/driver/bookings/123/accept \
  -H "Authorization: Bearer $DRIVER_TOKEN"
```

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

### Booking Permissions

#### Guest Users (No Authentication)
- Create bookings
- View bookings by email
- Request secure update links
- Update/cancel bookings with secure tokens

#### Authenticated Users
- All guest permissions
- Direct booking updates without tokens
- View personal booking history
- Automatic user association with bookings

### Rider (Default)
- Register and login
- View own profile
- Change own password
- **Book rides and manage bookings**
- **Update and cancel own bookings**

### Driver
- All rider permissions
- **Accept pending bookings**
- **View assigned rides**
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
- **Booking Audit Trail**: Complete tracking of booking creation, updates, and cancellations
- **Log Levels**: INFO, WARN, ERROR with timestamps
- **File + Console**: Logs written to both `app.log` and console
- **Request Tracking**: Full request/response logging with timing

## 🔄 Business Rules

### Booking Rules
- **24-Hour Advance Booking**: All bookings must be scheduled at least 24 hours in the future
- **24-Hour Cancellation Policy**: Bookings can only be cancelled 24+ hours before scheduled time
- **Status Protection**: Cannot update/cancel completed or already cancelled bookings
- **Driver Assignment**: Only drivers can accept bookings, and only pending bookings can be accepted
- **Guest Security**: Guest users receive secure, time-limited tokens (24-hour expiry) for updates

### Booking Statuses
- **Book Status**: `Pending` → `Accepted` → `Completed` or `Cancelled`
- **Ride Status**: `Pending` → `Assigned` → `In Progress` → `Completed` or `Cancelled`

## 🔮 Future Enhancements

- [x] ~~Ride booking and management system~~ ✅ **COMPLETED**
- [x] ~~Guest booking support with secure updates~~ ✅ **COMPLETED**
- [x] ~~24-hour booking and cancellation policies~~ ✅ **COMPLETED**
- [ ] Real-time ride tracking with WebSockets
- [ ] Payment integration (Stripe)
- [ ] Driver location tracking
- [ ] Push notifications
- [ ] Booking confirmation SMS
- [ ] Driver rating system
- [ ] Ride history and receipts
- [ ] Mobile app API support
- [ ] Advanced analytics and reporting

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**LuxSUV Backend API** - Built with ❤️ using Go, Echo, PostgreSQL, and MailerSend