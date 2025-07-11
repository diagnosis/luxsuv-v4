# LuxSUV Backend API

A Go-based REST API for a luxury ride-sharing service with user authentication, password reset, and email notifications.

## Features

- User registration and authentication with JWT tokens
- Password reset functionality with email notifications
- Role-based access control (rider, driver, admin)
- Admin user management
- Email service integration
- Comprehensive logging
- Database migrations with Goose
- Rate limiting and CORS support

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- SMTP email service (Gmail, Outlook, etc.)

## Setup

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd luxsuv-v4
go mod download
```

### 2. Environment Configuration

Copy the `.env` file and configure your settings:

```bash
# Database
DATABASE_URL=postgresql://username:password@host:port/database?sslmode=require

# JWT Secret (generate a secure 32+ character string)
JWT_SECRET=your-secure-jwt-secret-key-here

# Server
PORT=8080
ENVIRONMENT=development

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@gmail.com
```

### 3. Email Service Setup

#### For Gmail:
1. Enable 2-factor authentication on your Google account
2. Generate an App Password:
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a password for "Mail"
   - Use this password in `SMTP_PASSWORD`

#### For Outlook/Hotmail:
- Use `smtp-mail.outlook.com` as SMTP_HOST
- Use your regular email password

#### For Yahoo:
- Use `smtp.mail.yahoo.com` as SMTP_HOST
- Generate an app password in Yahoo account settings

### 4. Database Setup

The application will automatically run migrations on startup. Make sure your PostgreSQL database is accessible.

### 5. Run the Application

```bash
go run cmd/server/main.go
```

## API Endpoints

### Public Endpoints

- `POST /register` - User registration
- `POST /login` - User login
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password with token
- `GET /health` - Health check

### Protected Endpoints (Require Authentication)

- `GET /users/me` - Get current user info
- `PUT /users/me/password` - Change password

### Admin Endpoints (Require Admin Role)

- `GET /admin/users` - List all users (with pagination)
- `GET /admin/users/by-email?email=user@example.com` - Get user by email
- `GET /admin/users/:id` - Get user by ID
- `PUT /admin/users/:id/role` - Update user role
- `DELETE /admin/users/:id` - Delete user

## Usage Examples

### Register a New User

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123",
    "role": "rider"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Request Password Reset

```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

### Reset Password

```bash
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "reset_token": "your-reset-token-here",
    "new_password": "newpassword123"
  }'
```

### Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer your-jwt-token-here"
```

## Email Templates

The application includes HTML email templates for:

- **Password Reset**: Professional email with reset link and expiration notice
- **Welcome Email**: Sent to new users upon registration

## Security Features

- Password hashing with bcrypt
- JWT token authentication
- Rate limiting (5 req/sec general, 2 req/sec for auth endpoints)
- CORS protection
- Input validation and sanitization
- SQL injection protection with parameterized queries

## Development

### Running with Air (Hot Reload)

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

### Database Migrations

Migrations are automatically applied on startup. To create new migrations:

```bash
# Create a new migration
goose -dir migrations create migration_name sql
```

## Production Deployment

1. Set `ENVIRONMENT=production` in your environment variables
2. Use a strong JWT secret (32+ characters)
3. Configure proper CORS origins in `internal/middleware/cors_middleware.go`
4. Set up proper logging and monitoring
5. Use environment variables instead of `.env` file
6. Configure SSL/TLS for HTTPS

## Troubleshooting

### Email Not Sending

1. Check SMTP credentials and host settings
2. Verify app password for Gmail (not regular password)
3. Check firewall settings for SMTP port (587/465)
4. Review application logs for detailed error messages

### Database Connection Issues

1. Verify DATABASE_URL format and credentials
2. Ensure PostgreSQL is running and accessible
3. Check network connectivity and firewall settings

### JWT Token Issues

1. Ensure JWT_SECRET is at least 32 characters
2. Check token expiration (24 hours by default)
3. Verify Authorization header format: `Bearer <token>`

## License

This project is licensed under the MIT License.