#!/bin/bash

# Days API Test Script
# This script tests the Days Calendar API endpoints

# Configuration
BASE_URL="https://days.germainleignel.com"
HEALTH_URL="$BASE_URL/health"
API_URL="$BASE_URL/api"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS") echo -e "${GREEN}✓${NC} $message" ;;
        "ERROR") echo -e "${RED}✗${NC} $message" ;;
        "INFO") echo -e "${BLUE}ℹ${NC} $message" ;;
        "WARNING") echo -e "${YELLOW}⚠${NC} $message" ;;
    esac
}

# Function to make HTTP requests and handle responses
make_request() {
    local method=$1
    local url=$2
    local data=$3
    local auth_header=$4
    
    echo -e "\n${BLUE}→${NC} $method $url"
    
    if [ -n "$auth_header" ]; then
        if [ -n "$data" ]; then
            response=$(curl -s -w "%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -H "$auth_header" \
                -d "$data" \
                "$url")
        else
            response=$(curl -s -w "%{http_code}" -X "$method" \
                -H "$auth_header" \
                "$url")
        fi
    else
        if [ -n "$data" ]; then
            response=$(curl -s -w "%{http_code}" -X "$method" \
                -H "Content-Type: application/json" \
                -d "$data" \
                "$url")
        else
            response=$(curl -s -w "%{http_code}" "$url")
        fi
    fi
    
    # Extract status code (last 3 characters)
    status_code="${response: -3}"
    # Extract body (everything except last 3 characters)
    body="${response%???}"
    
    echo "Status: $status_code"
    if [ -n "$body" ]; then
        echo "Response: $body" | jq . 2>/dev/null || echo "Response: $body"
    fi
    
    return $status_code
}

echo -e "${BLUE}=====================================\n"
echo -e "   Days Calendar API Test Script\n"
echo -e "=====================================${NC}\n"

print_status "INFO" "Testing API at: $BASE_URL"

# Test 1: Health Check
print_status "INFO" "Testing health endpoint..."
make_request "GET" "$HEALTH_URL"
if [ $? -eq 200 ]; then
    print_status "SUCCESS" "Health check passed"
else
    print_status "ERROR" "Health check failed"
    exit 1
fi

# Test 2: Create User
print_status "INFO" "Creating test user..."
user_data='{
    "email": "test@example.com",
    "password": "testpassword123"
}'

make_request "POST" "$API_URL/users" "$user_data"
if [ $? -eq 201 ]; then
    print_status "SUCCESS" "User created successfully"
    USER_CREATED=true
elif [ $? -eq 400 ]; then
    print_status "WARNING" "User might already exist (400 error)"
    USER_CREATED=false
else
    print_status "ERROR" "Failed to create user"
    USER_CREATED=false
fi

# Test 3: Login
print_status "INFO" "Logging in..."
login_data='{
    "email": "test@example.com",
    "password": "testpassword123"
}'

make_request "POST" "$API_URL/auth/login" "$login_data"
login_status=$?

if [ $login_status -eq 200 ]; then
    print_status "SUCCESS" "Login successful"
    
    # Extract JWT token from response (assuming it's in the response)
    # You might need to adjust this based on your actual response format
    JWT_TOKEN=$(echo "$body" | jq -r '.token' 2>/dev/null)
    
    if [ "$JWT_TOKEN" != "null" ] && [ -n "$JWT_TOKEN" ]; then
        print_status "SUCCESS" "JWT token extracted"
        AUTH_HEADER="Authorization: Bearer $JWT_TOKEN"
        
        # Test 4: Get User Profile (requires auth)
        print_status "INFO" "Getting user profile..."
        # You'll need to extract user ID from login response or create user response
        USER_ID=$(echo "$body" | jq -r '.user.id' 2>/dev/null)
        
        if [ "$USER_ID" != "null" ] && [ -n "$USER_ID" ]; then
            make_request "GET" "$API_URL/users/$USER_ID" "" "$AUTH_HEADER"
            if [ $? -eq 200 ]; then
                print_status "SUCCESS" "User profile retrieved"
            else
                print_status "ERROR" "Failed to get user profile"
            fi
        else
            print_status "WARNING" "Could not extract user ID from login response"
        fi
        
        # Test 5: Create Calendar (requires auth)
        print_status "INFO" "Creating test calendar..."
        calendar_data='{
            "name": "Test Calendar",
            "description": "A test calendar for API testing"
        }'
        
        make_request "POST" "$API_URL/calendars" "$calendar_data" "$AUTH_HEADER"
        if [ $? -eq 201 ]; then
            print_status "SUCCESS" "Calendar created successfully"
            CALENDAR_ID=$(echo "$body" | jq -r '.id' 2>/dev/null)
        else
            print_status "ERROR" "Failed to create calendar"
        fi
        
        # Test 6: Get Calendars (requires auth)
        print_status "INFO" "Getting user calendars..."
        make_request "GET" "$API_URL/calendars" "" "$AUTH_HEADER"
        if [ $? -eq 200 ]; then
            print_status "SUCCESS" "Calendars retrieved successfully"
        else
            print_status "ERROR" "Failed to get calendars"
        fi
        
        # Test 7: Get Specific Calendar (if we have an ID)
        if [ "$CALENDAR_ID" != "null" ] && [ -n "$CALENDAR_ID" ]; then
            print_status "INFO" "Getting specific calendar..."
            make_request "GET" "$API_URL/calendars/$CALENDAR_ID" "" "$AUTH_HEADER"
            if [ $? -eq 200 ]; then
                print_status "SUCCESS" "Calendar retrieved successfully"
            else
                print_status "ERROR" "Failed to get specific calendar"
            fi
            
            # Test 8: Update Calendar
            print_status "INFO" "Updating calendar..."
            update_data='{
                "name": "Updated Test Calendar",
                "description": "Updated description for testing"
            }'
            
            make_request "PUT" "$API_URL/calendars/$CALENDAR_ID" "$update_data" "$AUTH_HEADER"
            if [ $? -eq 200 ]; then
                print_status "SUCCESS" "Calendar updated successfully"
            else
                print_status "ERROR" "Failed to update calendar"
            fi
        fi
        
    else
        print_status "ERROR" "Could not extract JWT token from login response"
    fi
else
    print_status "ERROR" "Login failed"
fi

# Test 9: Test Swagger Documentation
print_status "INFO" "Testing Swagger documentation..."
make_request "GET" "$BASE_URL/swagger/"
if [ $? -eq 200 ]; then
    print_status "SUCCESS" "Swagger documentation is accessible"
else
    print_status "WARNING" "Swagger documentation might not be available"
fi

echo -e "\n${BLUE}=====================================\n"
print_status "INFO" "API testing completed!"
echo -e "${BLUE}=====================================${NC}\n"

# Cleanup instructions
if [ "$USER_CREATED" = true ] && [ "$CALENDAR_ID" != "null" ] && [ -n "$CALENDAR_ID" ]; then
    echo -e "${YELLOW}Cleanup Commands:${NC}"
    echo "To delete the test calendar:"
    echo "curl -X DELETE -H \"$AUTH_HEADER\" \"$API_URL/calendars/$CALENDAR_ID\""
    echo ""
fi

print_status "INFO" "You can now test your API manually at: $BASE_URL/swagger/"
