# Days API Test Script (PowerShell) - Fixed Version
# This script tests the Days Calendar API endpoints

param(
    [string]$BaseUrl = "https://days.germainleignel.com"
)

# Configuration
$HealthUrl = "$BaseUrl/health"
$ApiUrl = "$BaseUrl/api"

# Function to print colored output
function Write-Status {
    param(
        [string]$Status,
        [string]$Message
    )
    
    switch ($Status) {
        "SUCCESS" { 
            Write-Host "✓ " -ForegroundColor Green -NoNewline
            Write-Host $Message
        }
        "ERROR" { 
            Write-Host "✗ " -ForegroundColor Red -NoNewline
            Write-Host $Message
        }
        "INFO" { 
            Write-Host "ℹ " -ForegroundColor Blue -NoNewline
            Write-Host $Message
        }
        "WARNING" { 
            Write-Host "⚠ " -ForegroundColor Yellow -NoNewline
            Write-Host $Message
        }
    }
}

# Function to make HTTP requests
function Invoke-ApiRequest {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Body = $null,
        [hashtable]$Headers = @{}
    )
    
    Write-Host ""
    Write-Host "→ $Method $Url" -ForegroundColor Blue
    
    try {
        $requestParams = @{
            Uri = $Url
            Method = $Method
            Headers = $Headers
            ContentType = "application/json"
        }
        
        if ($Body) {
            $requestParams.Body = $Body
        }
        
        $response = Invoke-RestMethod @requestParams
        Write-Host "Status: 200 (Success)" -ForegroundColor Green
        
        if ($response) {
            Write-Host "Response:" -ForegroundColor Cyan
            $response | ConvertTo-Json -Depth 3 | Write-Host
        }
        
        return @{
            Success = $true
            Data = $response
            StatusCode = 200
        }
    }
    catch {
        $statusCode = if ($_.Exception.Response) { $_.Exception.Response.StatusCode.Value__ } else { 0 }
        Write-Host "Status: $statusCode" -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        
        return @{
            Success = $false
            Data = $null
            StatusCode = $statusCode
        }
    }
}

Clear-Host
Write-Host "=====================================" -ForegroundColor Blue
Write-Host "   Days Calendar API Test Script" -ForegroundColor Blue
Write-Host "=====================================" -ForegroundColor Blue
Write-Host ""

Write-Status "INFO" "Testing API at: $BaseUrl"

# Test 1: Health Check
Write-Status "INFO" "Testing health endpoint..."
$healthResult = Invoke-ApiRequest -Method "GET" -Url $HealthUrl

if ($healthResult.Success) {
    Write-Status "SUCCESS" "Health check passed"
} else {
    Write-Status "ERROR" "Health check failed"
    exit 1
}

# Test 2: Create User
Write-Status "INFO" "Creating test user..."
$userData = @{
    email = "test@example.com"
    password = "testpassword123"
} | ConvertTo-Json

$createUserResult = Invoke-ApiRequest -Method "POST" -Url "$ApiUrl/users" -Body $userData

if ($createUserResult.Success) {
    Write-Status "SUCCESS" "User created successfully"
    $userCreated = $true
} elseif ($createUserResult.StatusCode -eq 400) {
    Write-Status "WARNING" "User might already exist (400 error)"
    $userCreated = $false
} else {
    Write-Status "ERROR" "Failed to create user"
    $userCreated = $false
}

# Test 3: Login
Write-Status "INFO" "Logging in..."
$loginData = @{
    email = "test@example.com"
    password = "testpassword123"
} | ConvertTo-Json

$loginResult = Invoke-ApiRequest -Method "POST" -Url "$ApiUrl/auth/login" -Body $loginData

if ($loginResult.Success) {
    Write-Status "SUCCESS" "Login successful"
    
    # Extract JWT token
    $jwtToken = $loginResult.Data.token
    
    if ($jwtToken) {
        Write-Status "SUCCESS" "JWT token extracted"
        $authHeaders = @{
            "Authorization" = "Bearer $jwtToken"
        }
        
        # Extract user ID
        $userId = $loginResult.Data.user.id
        
        # Test 4: Get User Profile (requires auth)
        if ($userId) {
            Write-Status "INFO" "Getting user profile..."
            $userProfileResult = Invoke-ApiRequest -Method "GET" -Url "$ApiUrl/users/$userId" -Headers $authHeaders
            
            if ($userProfileResult.Success) {
                Write-Status "SUCCESS" "User profile retrieved"
            } else {
                Write-Status "ERROR" "Failed to get user profile"
            }
        }
        
        # Test 5: Create Calendar (requires auth)
        Write-Status "INFO" "Creating test calendar..."
        $calendarData = @{
            name = "Test Calendar"
            description = "A test calendar for API testing"
        } | ConvertTo-Json
        
        $createCalendarResult = Invoke-ApiRequest -Method "POST" -Url "$ApiUrl/calendars" -Body $calendarData -Headers $authHeaders
        
        if ($createCalendarResult.Success) {
            Write-Status "SUCCESS" "Calendar created successfully"
            $calendarId = $createCalendarResult.Data.id
        } else {
            Write-Status "ERROR" "Failed to create calendar"
            $calendarId = $null
        }
        
        # Test 6: Get Calendars (requires auth)
        Write-Status "INFO" "Getting user calendars..."
        $getCalendarsResult = Invoke-ApiRequest -Method "GET" -Url "$ApiUrl/calendars" -Headers $authHeaders
        
        if ($getCalendarsResult.Success) {
            Write-Status "SUCCESS" "Calendars retrieved successfully"
        } else {
            Write-Status "ERROR" "Failed to get calendars"
        }
        
        # Test 7: Get Specific Calendar (if we have an ID)
        if ($calendarId) {
            Write-Status "INFO" "Getting specific calendar..."
            $getCalendarResult = Invoke-ApiRequest -Method "GET" -Url "$ApiUrl/calendars/$calendarId" -Headers $authHeaders
            
            if ($getCalendarResult.Success) {
                Write-Status "SUCCESS" "Calendar retrieved successfully"
            } else {
                Write-Status "ERROR" "Failed to get specific calendar"
            }
            
            # Test 8: Update Calendar
            Write-Status "INFO" "Updating calendar..."
            $updateData = @{
                name = "Updated Test Calendar"
                description = "Updated description for testing"
            } | ConvertTo-Json
            
            $updateCalendarResult = Invoke-ApiRequest -Method "PUT" -Url "$ApiUrl/calendars/$calendarId" -Body $updateData -Headers $authHeaders
            
            if ($updateCalendarResult.Success) {
                Write-Status "SUCCESS" "Calendar updated successfully"
            } else {
                Write-Status "ERROR" "Failed to update calendar"
            }
        }
        
    } else {
        Write-Status "ERROR" "Could not extract JWT token from login response"
    }
} else {
    Write-Status "ERROR" "Login failed"
}

# Test 9: Test Swagger Documentation
Write-Status "INFO" "Testing Swagger documentation..."
try {
    $swaggerResult = Invoke-WebRequest -Uri "$BaseUrl/swagger/" -Method GET
    if ($swaggerResult.StatusCode -eq 200) {
        Write-Status "SUCCESS" "Swagger documentation is accessible"
    }
} catch {
    Write-Status "WARNING" "Swagger documentation might not be available"
}

Write-Host ""
Write-Host "=====================================" -ForegroundColor Blue
Write-Status "INFO" "API testing completed!"
Write-Host "=====================================" -ForegroundColor Blue
Write-Host ""

# Cleanup instructions
if ($userCreated -and $calendarId) {
    Write-Host "Cleanup Commands:" -ForegroundColor Yellow
    Write-Host "To delete the test calendar, run:"
    Write-Host "Invoke-RestMethod -Uri '$ApiUrl/calendars/$calendarId' -Method DELETE -Headers @{'Authorization'='Bearer $jwtToken'}" -ForegroundColor Gray
    Write-Host ""
}

Write-Status "INFO" "You can now test your API manually at: $BaseUrl/swagger/"
