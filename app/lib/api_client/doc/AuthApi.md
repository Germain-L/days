# days_api_client.api.AuthApi

## Load the API package
```dart
import 'package:days_api_client/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**apiAuthLoginPost**](AuthApi.md#apiauthloginpost) | **POST** /api/auth/login | User login


# **apiAuthLoginPost**
> ServicesLoginResponse apiAuthLoginPost(credentials)

User login

Authenticate user and return JWT token

### Example
```dart
import 'package:days_api_client/api.dart';

final api_instance = AuthApi();
final credentials = ServicesLoginRequest(); // ServicesLoginRequest | Login credentials

try {
    final result = api_instance.apiAuthLoginPost(credentials);
    print(result);
} catch (e) {
    print('Exception when calling AuthApi->apiAuthLoginPost: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **credentials** | [**ServicesLoginRequest**](ServicesLoginRequest.md)| Login credentials | 

### Return type

[**ServicesLoginResponse**](ServicesLoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

