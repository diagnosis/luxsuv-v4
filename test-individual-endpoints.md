# Individual Test Commands

## Quick Test Commands (Copy & Paste)

### 1. Health Check
```bash
curl -X GET http://localhost:8080/health
```

### 2. Register Users
```bash
# Register Rider
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testrider","email":"rider@test.com","password":"password123","role":"rider"}'

# Register Driver  
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testdriver","email":"driver@test.com","password":"password123","role":"driver"}'

# Register Admin
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testadmin","email":"admin@test.com","password":"password123","role":"admin"}'
```

### 3. Login and Get Tokens
```bash
# Login as Rider (save the token)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"rider@test.com","password":"password123"}'

# Login as Driver (save the token)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"driver@test.com","password":"password123"}'

# Login as Admin (save the token)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"password123"}'
```

### 4. Booking Tests (Replace YOUR_TOKEN with actual token)
```bash
# Create booking as authenticated user
curl -X POST http://localhost:8080/book-ride \
  -H "Authorization: Bearer YOUR_RIDER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"your_name":"John Rider","email":"rider@test.com","phone_number":"+1234567890","ride_type":"Airport Transfer","pickup_location":"123 Main St","dropoff_location":"Airport Terminal 1","date":"2025-01-20","time":"14:30","number_of_passengers":2,"number_of_luggage":3,"additional_notes":"Flight AA123"}'

# Create booking as guest
curl -X POST http://localhost:8080/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"Jane Guest","email":"guest@test.com","phone_number":"+1987654321","ride_type":"City Tour","pickup_location":"Hotel Downtown","dropoff_location":"City Center","date":"2025-01-22","time":"10:00","number_of_passengers":4,"number_of_luggage":0}'

# Get my bookings
curl -X GET http://localhost:8080/bookings/my \
  -H "Authorization: Bearer YOUR_RIDER_TOKEN"

# Get bookings by email
curl -X GET http://localhost:8080/bookings/email/guest%40test.com

# Update booking (replace 1 with actual booking ID)
curl -X PUT http://localhost:8080/bookings/1 \
  -H "Authorization: Bearer YOUR_RIDER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"pickup_location":"Updated Location","time":"15:00"}'

# Generate update link for guest
curl -X POST http://localhost:8080/bookings/1/update-link \
  -H "Content-Type: application/json" \
  -d '{"email":"guest@test.com"}'

# Driver accept booking
curl -X PUT http://localhost:8080/driver/bookings/1/accept \
  -H "Authorization: Bearer YOUR_DRIVER_TOKEN"
```

### 5. Admin Tests
```bash
# List all users
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Get user by email
curl -X GET "http://localhost:8080/admin/users/by-email?email=rider@test.com" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Update user role
curl -X PUT http://localhost:8080/admin/users/1/role \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"driver"}'
```

### 6. Password Reset Tests
```bash
# Request password reset
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"rider@test.com"}'

# Reset password (replace TOKEN with actual reset token)
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"reset_token":"YOUR_RESET_TOKEN","new_password":"newpassword123"}'
```

### 7. Error Handling Tests
```bash
# Test invalid token
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer invalid_token"

# Test missing authorization
curl -X GET http://localhost:8080/users/me

# Test invalid booking data
curl -X POST http://localhost:8080/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"","email":"invalid-email"}'

# Test non-admin access to admin endpoint
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer YOUR_RIDER_TOKEN"
```

## Expected Results

### ✅ Success Cases:
- Health check returns `{"status":"healthy"}`
- Registration returns user object with ID
- Login returns token and user object
- Booking creation returns booking with ID
- Admin endpoints work with admin token
- Driver can accept bookings

### ❌ Error Cases:
- Invalid tokens return 401 Unauthorized
- Missing auth returns 401 Unauthorized  
- Invalid data returns 400 Bad Request
- Non-admin access to admin endpoints returns 403 Forbidden
- Non-driver access to driver endpoints returns 403 Forbidden
- Bookings less than 24 hours in advance are rejected

## Tips:
1. Save tokens from login responses to use in subsequent requests
2. Replace booking IDs with actual IDs from creation responses
3. Make sure server is running on localhost:8080
4. Check server logs for detailed error information
5. Some email features may not work without proper email configuration