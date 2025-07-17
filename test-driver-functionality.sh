#!/bin/bash

# LuxSUV Driver Functionality Test Commands
# Run these commands to test the new driver system

echo "🚗 LuxSUV Driver Functionality - Test Commands"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080"

echo -e "\n${BLUE}1. REGISTER USERS WITH NEW ROLES${NC}"
echo "================================="

echo -e "\n${YELLOW}1.1 Register Regular Driver${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"regulardriver\",\"email\":\"regulardriver@test.com\",\"password\":\"password123\",\"role\":\"driver\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"regulardriver","email":"regulardriver@test.com","password":"password123","role":"driver"}'
echo -e "\n"

echo -e "\n${YELLOW}1.2 Register Super Driver${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"superdriver\",\"email\":\"superdriver@test.com\",\"password\":\"password123\",\"role\":\"super_driver\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"superdriver","email":"superdriver@test.com","password":"password123","role":"super_driver"}'
echo -e "\n"

echo -e "\n${YELLOW}1.3 Register Dispatcher${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"dispatcher\",\"email\":\"dispatcher@test.com\",\"password\":\"password123\",\"role\":\"dispatcher\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"dispatcher","email":"dispatcher@test.com","password":"password123","role":"dispatcher"}'
echo -e "\n"

echo -e "\n${YELLOW}1.4 Register Second Driver${NC}"
echo "curl -X POST $BASE_URL/register \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\":\"driver2\",\"email\":\"driver2@test.com\",\"password\":\"password123\",\"role\":\"driver\"}'"
curl -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"driver2","email":"driver2@test.com","password":"password123","role":"driver"}'
echo -e "\n"

echo -e "\n${BLUE}2. LOGIN AND GET TOKENS${NC}"
echo "======================="

echo -e "\n${YELLOW}2.1 Login as Regular Driver${NC}"
DRIVER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"regulardriver@test.com","password":"password123"}')
echo $DRIVER_RESPONSE
DRIVER_TOKEN=$(echo $DRIVER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Driver Token: $DRIVER_TOKEN${NC}"

echo -e "\n${YELLOW}2.2 Login as Super Driver${NC}"
SUPER_DRIVER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"superdriver@test.com","password":"password123"}')
echo $SUPER_DRIVER_RESPONSE
SUPER_DRIVER_TOKEN=$(echo $SUPER_DRIVER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Super Driver Token: $SUPER_DRIVER_TOKEN${NC}"

echo -e "\n${YELLOW}2.3 Login as Dispatcher${NC}"
DISPATCHER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dispatcher@test.com","password":"password123"}')
echo $DISPATCHER_RESPONSE
DISPATCHER_TOKEN=$(echo $DISPATCHER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Dispatcher Token: $DISPATCHER_TOKEN${NC}"

echo -e "\n${YELLOW}2.4 Login as Second Driver${NC}"
DRIVER2_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"email":"driver2@test.com","password":"password123"}')
echo $DRIVER2_RESPONSE
DRIVER2_TOKEN=$(echo $DRIVER2_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo -e "\n${GREEN}Driver2 Token: $DRIVER2_TOKEN${NC}"

echo -e "\n${BLUE}3. CREATE TEST BOOKINGS${NC}"
echo "======================"

echo -e "\n${YELLOW}3.1 Create First Booking${NC}"
BOOKING1_RESPONSE=$(curl -s -X POST $BASE_URL/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"Test Customer 1","email":"customer1@test.com","phone_number":"+1234567890","ride_type":"Airport Transfer","pickup_location":"Hotel A","dropoff_location":"Airport","date":"2025-01-25","time":"10:00","number_of_passengers":2,"number_of_luggage":3}')
echo $BOOKING1_RESPONSE
BOOKING1_ID=$(echo $BOOKING1_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo -e "\n${GREEN}Booking 1 ID: $BOOKING1_ID${NC}"

echo -e "\n${YELLOW}3.2 Create Second Booking${NC}"
BOOKING2_RESPONSE=$(curl -s -X POST $BASE_URL/book-ride \
  -H "Content-Type: application/json" \
  -d '{"your_name":"Test Customer 2","email":"customer2@test.com","phone_number":"+1987654321","ride_type":"City Tour","pickup_location":"Hotel B","dropoff_location":"Downtown","date":"2025-01-26","time":"14:00","number_of_passengers":1,"number_of_luggage":1}')
echo $BOOKING2_RESPONSE
BOOKING2_ID=$(echo $BOOKING2_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo -e "\n${GREEN}Booking 2 ID: $BOOKING2_ID${NC}"

echo -e "\n${BLUE}4. DISPATCHER FUNCTIONALITY${NC}"
echo "==========================="

echo -e "\n${YELLOW}4.1 Dispatcher - View All Bookings${NC}"
echo "curl -X GET $BASE_URL/dispatcher/bookings/all \\"
echo "  -H \"Authorization: Bearer \$DISPATCHER_TOKEN\""
curl -X GET $BASE_URL/dispatcher/bookings/all \
  -H "Authorization: Bearer $DISPATCHER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}4.2 Dispatcher - View Available Bookings${NC}"
echo "curl -X GET $BASE_URL/dispatcher/bookings/available \\"
echo "  -H \"Authorization: Bearer \$DISPATCHER_TOKEN\""
curl -X GET $BASE_URL/dispatcher/bookings/available \
  -H "Authorization: Bearer $DISPATCHER_TOKEN"
echo -e "\n"

# Get driver IDs from the login responses for assignment
DRIVER1_ID=$(echo $DRIVER_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
DRIVER2_ID=$(echo $DRIVER2_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)

echo -e "\n${YELLOW}4.3 Dispatcher - Assign Booking 1 to Driver 1${NC}"
echo "curl -X POST $BASE_URL/dispatcher/bookings/$BOOKING1_ID/assign \\"
echo "  -H \"Authorization: Bearer \$DISPATCHER_TOKEN\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"driver_id\":$DRIVER1_ID,\"notes\":\"Assigned by dispatcher\"}'"
curl -X POST $BASE_URL/dispatcher/bookings/$BOOKING1_ID/assign \
  -H "Authorization: Bearer $DISPATCHER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"driver_id\":$DRIVER1_ID,\"notes\":\"Assigned by dispatcher\"}"
echo -e "\n"

echo -e "\n${YELLOW}4.4 Dispatcher - View Driver 1 Bookings${NC}"
echo "curl -X GET $BASE_URL/dispatcher/bookings/driver/$DRIVER1_ID \\"
echo "  -H \"Authorization: Bearer \$DISPATCHER_TOKEN\""
curl -X GET $BASE_URL/dispatcher/bookings/driver/$DRIVER1_ID \
  -H "Authorization: Bearer $DISPATCHER_TOKEN"
echo -e "\n"

echo -e "\n${BLUE}5. SUPER DRIVER FUNCTIONALITY${NC}"
echo "============================="

echo -e "\n${YELLOW}5.1 Super Driver - View Available Bookings${NC}"
echo "curl -X GET $BASE_URL/super-driver/bookings/available \\"
echo "  -H \"Authorization: Bearer \$SUPER_DRIVER_TOKEN\""
curl -X GET $BASE_URL/super-driver/bookings/available \
  -H "Authorization: Bearer $SUPER_DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}5.2 Super Driver - Assign Booking 2 to Driver 2${NC}"
echo "curl -X POST $BASE_URL/super-driver/bookings/$BOOKING2_ID/assign \\"
echo "  -H \"Authorization: Bearer \$SUPER_DRIVER_TOKEN\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"driver_id\":$DRIVER2_ID,\"notes\":\"Assigned by super driver\"}'"
curl -X POST $BASE_URL/super-driver/bookings/$BOOKING2_ID/assign \
  -H "Authorization: Bearer $SUPER_DRIVER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"driver_id\":$DRIVER2_ID,\"notes\":\"Assigned by super driver\"}"
echo -e "\n"

echo -e "\n${YELLOW}5.3 Super Driver - View Driver 2 Bookings${NC}"
echo "curl -X GET $BASE_URL/super-driver/bookings/driver/$DRIVER2_ID \\"
echo "  -H \"Authorization: Bearer \$SUPER_DRIVER_TOKEN\""
curl -X GET $BASE_URL/super-driver/bookings/driver/$DRIVER2_ID \
  -H "Authorization: Bearer $SUPER_DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${BLUE}6. REGULAR DRIVER FUNCTIONALITY${NC}"
echo "==============================="

echo -e "\n${YELLOW}6.1 Driver 1 - View Assigned Bookings${NC}"
echo "curl -X GET $BASE_URL/driver/bookings/assigned \\"
echo "  -H \"Authorization: Bearer \$DRIVER_TOKEN\""
curl -X GET $BASE_URL/driver/bookings/assigned \
  -H "Authorization: Bearer $DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}6.2 Driver 2 - View Assigned Bookings${NC}"
echo "curl -X GET $BASE_URL/driver/bookings/assigned \\"
echo "  -H \"Authorization: Bearer \$DRIVER2_TOKEN\""
curl -X GET $BASE_URL/driver/bookings/assigned \
  -H "Authorization: Bearer $DRIVER2_TOKEN"
echo -e "\n"

echo -e "\n${BLUE}7. MANAGEMENT ENDPOINTS (DISPATCHER OR SUPER-DRIVER)${NC}"
echo "===================================================="

echo -e "\n${YELLOW}7.1 Management - View Available Bookings (using Dispatcher token)${NC}"
echo "curl -X GET $BASE_URL/management/bookings/available \\"
echo "  -H \"Authorization: Bearer \$DISPATCHER_TOKEN\""
curl -X GET $BASE_URL/management/bookings/available \
  -H "Authorization: Bearer $DISPATCHER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}7.2 Management - View Available Bookings (using Super Driver token)${NC}"
echo "curl -X GET $BASE_URL/management/bookings/available \\"
echo "  -H \"Authorization: Bearer \$SUPER_DRIVER_TOKEN\""
curl -X GET $BASE_URL/management/bookings/available \
  -H "Authorization: Bearer $SUPER_DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${BLUE}8. ACCESS CONTROL TESTS${NC}"
echo "======================"

echo -e "\n${YELLOW}8.1 Regular Driver trying to access dispatcher endpoint (should fail)${NC}"
echo "curl -X GET $BASE_URL/dispatcher/bookings/all \\"
echo "  -H \"Authorization: Bearer \$DRIVER_TOKEN\""
curl -X GET $BASE_URL/dispatcher/bookings/all \
  -H "Authorization: Bearer $DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}8.2 Regular Driver trying to access super-driver endpoint (should fail)${NC}"
echo "curl -X GET $BASE_URL/super-driver/bookings/available \\"
echo "  -H \"Authorization: Bearer \$DRIVER_TOKEN\""
curl -X GET $BASE_URL/super-driver/bookings/available \
  -H "Authorization: Bearer $DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${YELLOW}8.3 Regular Driver trying to access management endpoint (should fail)${NC}"
echo "curl -X GET $BASE_URL/management/bookings/available \\"
echo "  -H \"Authorization: Bearer \$DRIVER_TOKEN\""
curl -X GET $BASE_URL/management/bookings/available \
  -H "Authorization: Bearer $DRIVER_TOKEN"
echo -e "\n"

echo -e "\n${GREEN}✅ All driver functionality tests completed!${NC}"
echo -e "\n${BLUE}Summary of what was tested:${NC}"
echo "- User registration with new roles (driver, super_driver, dispatcher)"
echo "- Dispatcher functionality (view all bookings, assign to drivers)"
echo "- Super-driver functionality (view available bookings, assign to drivers)"
echo "- Regular driver functionality (view assigned bookings only)"
echo "- Management endpoints (accessible by both dispatcher and super-driver)"
echo "- Access control (regular drivers cannot access elevated endpoints)"
echo "- Booking assignment workflow"
echo ""
echo -e "${YELLOW}Note: Make sure your server is running on localhost:8080${NC}"
echo -e "${YELLOW}The system now supports hierarchical driver management!${NC}"