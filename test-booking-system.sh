#!/bin/bash

# LuxSUV Booking System Test Commands
# Run these commands to validate all the fixes

echo "🚗 LuxSUV Booking System - Test Commands"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080"

echo -e "\n${BLUE}1. HEALTH CHECK${NC}"
echo "==============="
echo "curl -X GET $BASE_URL/health"
curl -X GET $BASE_URL/health
echo -e "\n"

echo -e "\n${BLUE}2. API INFO${NC}"
echo "==========="
echo "curl -X GET $BASE_URL/api/info"
curl -X GET $BASE_URL/api/info
echo -e "\n"

echo -e "\n${BLUE}3. USER REGISTRATION & AUTHENTICATION${NC}"
echo "====================================="

echo -e "\n${YELLOW}3.1 Register a Rider${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"testrider\",\"email\":\"rider@test.com\",\"password\":\"password123\",\"role\":\"rider\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testrider","email":"rider@test.com","password":"password123","role":"rider"}'
echo -e "\n"

echo -e "\n${YELLOW}3.2 Register a Driver${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"testdriver\",\"email\":\"driver@test.com\",\"password\":\"password123\",\"role\":\"driver\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testdriver","email":"driver@test.com","password":"password123","role":"driver"}'
echo -e "\n"

echo -e "\n${YELLOW}3.3 Register an Admin${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"testadmin\",\"email\":\"admin@test.com\",\"password\":\"password123\",\"role\":\"admin\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testadmin","email":"admin@test.com","password":"password123","role":"admin"}'
echo -e "\n"

echo -e "\n${YELLOW}3.4 Login as Rider (Save token for later use)${NC}"
echo "curl -X POST $BASE_URL/login \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"email\":\"rider@test.com\",\"password\":\"password123\"}'"
RIDER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"rider@test.com","password":"password123"}')
echo $RIDER_RESPONSE
RIDER_TOKEN=$(echo $RIDER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Rider Token: $RIDER_TOKEN${NC}"

echo -e "\n${YELLOW}3.5 Login as Driver (Save token for later use)${NC}"
echo "curl -X POST $BASE_URL/login \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"email\":\"driver@test.com\",\"password\":\"password123\"}'"
DRIVER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"driver@test.com","password":"password123"}')
echo $DRIVER_RESPONSE
DRIVER_TOKEN=$(echo $DRIVER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Driver Token: $DRIVER_TOKEN${NC}"

echo -e "\n${YELLOW}3.6 Login as Admin (Save token for later use)${NC}"
echo "curl -X POST $BASE_URL/login \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"email\":\"admin@test.com\",\"password\":\"password123\"}'"
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"password123"}')
echo $ADMIN_RESPONSE
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Admin Token: $ADMIN_TOKEN${NC}"

echo -e "\n${BLUE}4. BOOKING SYSTEM TESTS${NC}"
echo "======================="

echo -e "\n${YELLOW}4.1 Create Booking as Authenticated User${NC}"
echo "curl -X POST $BASE_URL/book-ride \\"
echo "  -H \"Authorization: Bearer \$RIDER_TOKEN\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"your_name\":\"John Rider\",\"email\":\"rider@test.com\",\"phone_number\":\"+1234567890\",\"ride_type\":\"Airport Transfer\",\"pickup_location\":\"123 Main St\",\"dropoff_location\":\"Airport Terminal 1\",\"date\":\"2025-01-20\",\"time\":\"14:30\",\"number_of_passengers\":2,\"number_of_luggage\":3,\"additional_notes\":\"Flight AA123\"}'"
BOOKING_RESPONSE=$(curl -s -X POST $BASE_URL/book-ride \
  -H "Authorization: Bearer $RIDER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"your_name":"John Rider","email":"rider@test.com","phone_number":"+1234567890","ride_type":"Airport Transfer","pickup_location":"123 Main St","dropoff_location":"Airport Terminal 1","date":"2025-01-20","time":"14:30","number_of_passengers":2,"number_of_luggage":3,"additional_notes":"Flight AA123"}')
echo $BOOKING_RESPONSE
BOOKING_ID=$(echo $BOOKING_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo -e "\n${GREEN}Booking ID: $BOOKING_ID${NC}"

echo -e "\n${YELLOW}4.2 Create Booking as Guest User${NC}"
echo "curl -X POST $BASE_URL/book-ride \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"your_name\":\"Jane Guest\",\"email\":\"guest@test.com\",\"phone_number\":\"+1987654321\",\"ride_type\":\"City Tour\",\"pickup_location\":\"Hotel Downtown\",\"dropoff_location\":\"City Center\",\"date\":\"2025-01-22\",\"time\":\"10:00\",\"number_of_passengers\":4,\"number_of_luggage\":0}'"
GUEST_BOOKING_RESPONSE=$(curl -s -X POST $BASE_URL/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"Jane Guest","email":"guest@test.com","phone_number":"+1987654321","ride_type":"City Tour","pickup_location":"Hotel Downtown","dropoff_location":"City Center","date":"2025-01-22","time":"10:00","number_of_passengers":4,"number_of_luggage":0}')
echo $GUEST_BOOKING_RESPONSE
GUEST_BOOKING_ID=$(echo $GUEST_BOOKING_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo -e "\n${GREEN}Guest Booking ID: $GUEST_BOOKING_ID${NC}"

echo -e "\n${YELLOW}4.3 Get My Bookings (Authenticated User)${NC}"
echo "curl -X GET $BASE_URL/bookings/my \\"
echo "  -H \"Authorization: Bearer \$RIDER_TOKEN\""
curl -X GET $BASE_URL/bookings/my \
  -H "Authorization: Bearer $RIDER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}4.4 Get Bookings by Email (Public)${NC}"
echo "curl -X GET $BASE_URL/bookings/email/guest%40test.com"
curl -X GET $BASE_URL/bookings/email/guest%40test.com
echo -e "\n"

echo -e "\n${YELLOW}4.5 Update Booking (Authenticated User)${NC}"
echo "curl -X PUT $BASE_URL/bookings/$BOOKING_ID \\"
echo "  -H \"Authorization: Bearer \$RIDER_TOKEN\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"pickup_location\":\"Updated Pickup Location\",\"time\":\"15:00\"}'"
curl -X PUT $BASE_URL/bookings/$BOOKING_ID \
  -H "Authorization: Bearer $RIDER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"pickup_location":"Updated Pickup Location","time":"15:00"}'
echo -e "\n"

echo -e "\n${YELLOW}4.6 Generate Update Link for Guest Booking${NC}"
echo "curl -X POST $BASE_URL/bookings/$GUEST_BOOKING_ID/update-link \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"email\":\"guest@test.com\"}'"
UPDATE_LINK_RESPONSE=$(curl -s -X POST $BASE_URL/bookings/$GUEST_BOOKING_ID/update-link \
  -H "Content-Type: application/json" \
  -d '{"email":"guest@test.com"}')
echo $UPDATE_LINK_RESPONSE
UPDATE_TOKEN=$(echo $UPDATE_LINK_RESPONSE | grep -o '"update_token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Update Token: $UPDATE_TOKEN${NC}"

if [ ! -z "$UPDATE_TOKEN" ]; then
  echo -e "\n${YELLOW}4.7 Update Guest Booking with Token${NC}"
  echo "curl -X PUT \"$BASE_URL/bookings/$GUEST_BOOKING_ID/update?token=\$UPDATE_TOKEN\" \\"
  echo "  -H \"Content-Type: application/json\" \\"
  echo "  -d '{\"your_name\":\"Jane Updated\",\"phone_number\":\"+1234567899\"}'"
  curl -X PUT "$BASE_URL/bookings/$GUEST_BOOKING_ID/update?token=$UPDATE_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"your_name":"Jane Updated","phone_number":"+1234567899"}'
  echo -e "\n"
fi

echo -e "\n${YELLOW}4.8 Driver Accept Booking${NC}"
echo "curl -X PUT $BASE_URL/driver/bookings/$BOOKING_ID/accept \\"
echo "  -H \"Authorization: Bearer \$DRIVER_TOKEN\""
curl -X PUT $BASE_URL/driver/bookings/$BOOKING_ID/accept \
  -H "Authorization: Bearer $DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${BLUE}5. ADMIN FUNCTIONALITY TESTS${NC}"
echo "============================"

echo -e "\n${YELLOW}5.1 List All Users (Admin)${NC}"
echo "curl -X GET $BASE_URL/admin/users \\"
echo "  -H \"Authorization: Bearer \$ADMIN_TOKEN\""
curl -X GET $BASE_URL/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}5.2 Get User by Email (Admin)${NC}"
echo "curl -X GET \"$BASE_URL/admin/users/by-email?email=rider@test.com\" \\"
echo "  -H \"Authorization: Bearer \$ADMIN_TOKEN\""
curl -X GET "$BASE_URL/admin/users/by-email?email=rider@test.com" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
echo -e "\n"

echo -e "\n${BLUE}6. PASSWORD RESET TESTS${NC}"
echo "======================"

echo -e "\n${YELLOW}6.1 Request Password Reset${NC}"
echo "curl -X POST $BASE_URL/auth/forgot-password \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"email\":\"rider@test.com\"}'"
RESET_RESPONSE=$(curl -s -X POST $BASE_URL/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"rider@test.com"}')
echo $RESET_RESPONSE
RESET_TOKEN=$(echo $RESET_RESPONSE | grep -o '"reset_token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Reset Token: $RESET_TOKEN${NC}"

if [ ! -z "$RESET_TOKEN" ]; then
  echo -e "\n${YELLOW}6.2 Reset Password with Token${NC}"
  echo "curl -X POST $BASE_URL/auth/reset-password \\"
  echo "  -H \"Content-Type: application/json\" \\"
  echo "  -d '{\"reset_token\":\"\$RESET_TOKEN\",\"new_password\":\"newpassword123\"}'"
  curl -X POST $BASE_URL/auth/reset-password \
    -H "Content-Type: application/json" \
    -d "{\"reset_token\":\"$RESET_TOKEN\",\"new_password\":\"newpassword123\"}"
  echo -e "\n"
fi

echo -e "\n${BLUE}7. ERROR HANDLING TESTS${NC}"
echo "======================"

echo -e "\n${YELLOW}7.1 Test Invalid Token${NC}"
echo "curl -X GET $BASE_URL/users/me \\"
echo "  -H \"Authorization: Bearer invalid_token\""
curl -X GET $BASE_URL/users/me \
  -H "Authorization: Bearer invalid_token"
echo -e "\n"

echo -e "\n${YELLOW}7.2 Test Missing Authorization${NC}"
echo "curl -X GET $BASE_URL/users/me"
curl -X GET $BASE_URL/users/me
echo -e "\n"

echo -e "\n${YELLOW}7.3 Test Invalid Booking Data${NC}"
echo "curl -X POST $BASE_URL/book-ride \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"your_name\":\"\",\"email\":\"invalid-email\"}'"
curl -X POST $BASE_URL/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"","email":"invalid-email"}'
echo -e "\n"

echo -e "\n${YELLOW}7.4 Test Booking Too Soon (Less than 24 hours)${NC}"
echo "curl -X POST $BASE_URL/book-ride \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"your_name\":\"Test User\",\"email\":\"test@test.com\",\"phone_number\":\"+1234567890\",\"ride_type\":\"Test\",\"pickup_location\":\"A\",\"dropoff_location\":\"B\",\"date\":\"2025-01-11\",\"time\":\"10:00\",\"number_of_passengers\":1,\"number_of_luggage\":0}'"
curl -X POST $BASE_URL/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"Test User","email":"test@test.com","phone_number":"+1234567890","ride_type":"Test","pickup_location":"A","dropoff_location":"B","date":"2025-01-11","time":"10:00","number_of_passengers":1,"number_of_luggage":0}'
echo -e "\n"

echo -e "\n${YELLOW}7.5 Test Non-Admin Access to Admin Endpoint${NC}"
echo "curl -X GET $BASE_URL/admin/users \\"
echo "  -H \"Authorization: Bearer \$RIDER_TOKEN\""
curl -X GET $BASE_URL/admin/users \
  -H "Authorization: Bearer $RIDER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}7.6 Test Non-Driver Access to Driver Endpoint${NC}"
echo "curl -X PUT $BASE_URL/driver/bookings/1/accept \\"
echo "  -H \"Authorization: Bearer \$RIDER_TOKEN\""
curl -X PUT $BASE_URL/driver/bookings/1/accept \
  -H "Authorization: Bearer $RIDER_TOKEN"
echo -e "\n"

echo -e "\n${GREEN}✅ All tests completed!${NC}"
echo -e "\n${BLUE}Summary of what was tested:${NC}"
echo "- Health checks and API info"
echo "- User registration (rider, driver, admin)"
echo "- Authentication and token generation"
echo "- Booking creation (authenticated and guest)"
echo "- Booking retrieval and updates"
echo "- Guest booking update links and tokens"
echo "- Driver booking acceptance"
echo "- Admin user management"
echo "- Password reset functionality"
echo "- Error handling and validation"
echo "- Role-based access control"
echo ""
echo -e "${YELLOW}Note: Make sure your server is running on localhost:8080${NC}"
echo -e "${YELLOW}Some tests may fail if email service is not configured${NC}"