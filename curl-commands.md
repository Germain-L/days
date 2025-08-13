# Days API Quick Test Commands
# Copy and paste these curl commands to test your API

# 1. Health Check
curl -X GET "https://days.germainleignel.com/health"

# 2. Create User
curl -X POST "https://days.germainleignel.com/api/users" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123"
  }'

# 3. Login (save the token from response)
curl -X POST "https://days.germainleignel.com/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "testpassword123"
  }'

# 4. Get User Profile (replace USER_ID and JWT_TOKEN)
curl -X GET "https://days.germainleignel.com/api/users/USER_ID" \
  -H "Authorization: Bearer JWT_TOKEN"

# 5. Create Calendar (replace JWT_TOKEN)
curl -X POST "https://days.germainleignel.com/api/calendars" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{
    "name": "My Test Calendar",
    "description": "A calendar for testing"
  }'

# 6. Get All Calendars (replace JWT_TOKEN)
curl -X GET "https://days.germainleignel.com/api/calendars" \
  -H "Authorization: Bearer JWT_TOKEN"

# 7. Get Specific Calendar (replace CALENDAR_ID and JWT_TOKEN)
curl -X GET "https://days.germainleignel.com/api/calendars/CALENDAR_ID" \
  -H "Authorization: Bearer JWT_TOKEN"

# 8. Update Calendar (replace CALENDAR_ID and JWT_TOKEN)
curl -X PUT "https://days.germainleignel.com/api/calendars/CALENDAR_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer JWT_TOKEN" \
  -d '{
    "name": "Updated Calendar Name",
    "description": "Updated description"
  }'

# 9. Delete Calendar (replace CALENDAR_ID and JWT_TOKEN)
curl -X DELETE "https://days.germainleignel.com/api/calendars/CALENDAR_ID" \
  -H "Authorization: Bearer JWT_TOKEN"

# 10. Access Swagger Documentation
curl -X GET "https://days.germainleignel.com/swagger/"

# ---
# Example workflow:
# 1. Run health check
# 2. Create user (or skip if exists)
# 3. Login and copy the JWT token
# 4. Replace JWT_TOKEN in subsequent commands
# 5. Create calendar and copy the calendar ID
# 6. Replace CALENDAR_ID in subsequent commands
# 7. Test CRUD operations on calendars
