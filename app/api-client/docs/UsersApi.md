# UsersApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**apiUsersIdGet**](UsersApi.md#apiUsersIdGet) | **GET** api/users/{id} | Get user by ID |
| [**apiUsersPost**](UsersApi.md#apiUsersPost) | **POST** api/users | Create a new user |



Get user by ID

Retrieve user information by user ID (self only)

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(UsersApi::class.java)
val id : kotlin.String = id_example // kotlin.String | User ID

launch(Dispatchers.IO) {
    val result : ServicesUserResponse = webService.apiUsersIdGet(id)
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **id** | **kotlin.String**| User ID | |

### Return type

[**ServicesUserResponse**](ServicesUserResponse.md)

### Authorization



### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


Create a new user

Create a new user account with email and password

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(UsersApi::class.java)
val user : ServicesCreateUserRequest =  // ServicesCreateUserRequest | User creation request

launch(Dispatchers.IO) {
    val result : ServicesUserResponse = webService.apiUsersPost(user)
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **user** | [**ServicesCreateUserRequest**](ServicesCreateUserRequest.md)| User creation request | |

### Return type

[**ServicesUserResponse**](ServicesUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

