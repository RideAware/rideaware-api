#!/bin/bash

BASE_URL="http://localhost:5000"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║    Testing RideAware API               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"

# Test 1: Health check
echo -e "${YELLOW}1. Health Check${NC}"
curl -s -X GET "$BASE_URL/health"
echo -e "\n\n"

# Test 2: Signup
echo -e "${YELLOW}2. Signup (New User)${NC}"
SIGNUP_RESPONSE=$(curl -s -X POST "$BASE_URL/api/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "blakearidgway",
    "password": "SecurePass123",
    "email": "blakearidgway@gmail.com",
    "first_name": "Blake",
    "last_name": "Ridgway"
  }')
echo "$SIGNUP_RESPONSE" | jq .
ACCESS_TOKEN=$(echo "$SIGNUP_RESPONSE" | jq -r '.access_token // empty')
echo -e "\n"

# Test 3: Login
echo -e "${YELLOW}3. Login${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "blakearidgway",
    "password": "SecurePass123"
  }')
echo "$LOGIN_RESPONSE" | jq .
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token // empty')
echo -e "\n"

# Test 4: Protected route with access token
echo -e "${YELLOW}4. Protected Route (with Access Token)${NC}"
if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" == "null" ]; then
	echo -e "${RED}No access token available${NC}"
else
	echo "Using token: ${ACCESS_TOKEN:0:50}..."
	curl -s -X GET "$BASE_URL/api/protected/profile" \
		-H "Authorization: Bearer $ACCESS_TOKEN" | jq .
fi
echo -e "\n"

# Test 5: Invalid token
echo -e "${YELLOW}5. Protected Route (with Invalid Token - should fail)${NC}"
curl -s -X GET "$BASE_URL/api/protected/profile" \
	-H "Authorization: Bearer invalid_token_here" | jq .
echo -e "\n"

# Test 6: Missing auth header (should fail)
echo -e "${YELLOW}6. Protected Route (without Auth Header - should fail)${NC}"
curl -s -X GET "$BASE_URL/api/protected/profile" | jq .
echo -e "\n"

# Test 7: Password reset request
echo -e "${YELLOW}7. Request Password Reset${NC}"
curl -s -X POST "$BASE_URL/api/password-reset/request" \
	-H "Content-Type: application/json" \
	-d '{
    "email": "blakearidgway@gmail.com"
  }' | jq .
echo -e "\n"

# Test 8: Logout
echo -e "${YELLOW}8. Logout${NC}"
curl -s -X POST "$BASE_URL/api/logout" \
	-H "Content-Type: application/json" | jq .
echo -e "\n"

echo -e "${GREEN}✓ Tests complete!${NC}"