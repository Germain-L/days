# AuthApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**apiAuthLoginPost**](AuthApi.md#apiAuthLoginPost) | **POST** api/auth/login | User login |



User login

Authenticate user and return JWT token

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(AuthApi::class.java)
val credentials : ServicesLoginRequest =  // ServicesLoginRequest | Login credentials

launch(Dispatchers.IO) {
    val result : ServicesLoginResponse = webService.apiAuthLoginPost(credentials)
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **credentials** | [**ServicesLoginRequest**](ServicesLoginRequest.md)| Login credentials | |

### Return type

[**ServicesLoginResponse**](ServicesLoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

