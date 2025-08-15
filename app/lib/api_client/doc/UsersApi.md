# days_api_client.api.UsersApi

## Load the API package
```dart
import 'package:days_api_client/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**apiUsersIdGet**](UsersApi.md#apiusersidget) | **GET** /api/users/{id} | Get user by ID
[**apiUsersPost**](UsersApi.md#apiuserspost) | **POST** /api/users | Create a new user


# **apiUsersIdGet**
> ServicesUserResponse apiUsersIdGet(id)

Get user by ID

Retrieve user information by user ID (self only)

### Example
```dart
import 'package:days_api_client/api.dart';
// TODO Configure API key authorization: BearerAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKeyPrefix = 'Bearer';

final api_instance = UsersApi();
final id = id_example; // String | User ID

try {
    final result = api_instance.apiUsersIdGet(id);
    print(result);
} catch (e) {
    print('Exception when calling UsersApi->apiUsersIdGet: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| User ID | 

### Return type

[**ServicesUserResponse**](ServicesUserResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiUsersPost**
> ServicesUserResponse apiUsersPost(user)

Create a new user

Create a new user account with email and password

### Example
```dart
import 'package:days_api_client/api.dart';

final api_instance = UsersApi();
final user = ServicesCreateUserRequest(); // ServicesCreateUserRequest | User creation request

try {
    final result = api_instance.apiUsersPost(user);
    print(result);
} catch (e) {
    print('Exception when calling UsersApi->apiUsersPost: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | [**ServicesCreateUserRequest**](ServicesCreateUserRequest.md)| User creation request | 

### Return type

[**ServicesUserResponse**](ServicesUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

