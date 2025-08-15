# days_api_client.api.CalendarsApi

## Load the API package
```dart
import 'package:days_api_client/api.dart';
```

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**apiCalendarsGet**](CalendarsApi.md#apicalendarsget) | **GET** /api/calendars | Get user calendars
[**apiCalendarsIdDelete**](CalendarsApi.md#apicalendarsiddelete) | **DELETE** /api/calendars/{id} | Delete calendar
[**apiCalendarsIdGet**](CalendarsApi.md#apicalendarsidget) | **GET** /api/calendars/{id} | Get calendar by ID
[**apiCalendarsIdPut**](CalendarsApi.md#apicalendarsidput) | **PUT** /api/calendars/{id} | Update calendar
[**apiCalendarsPost**](CalendarsApi.md#apicalendarspost) | **POST** /api/calendars | Create a new calendar


# **apiCalendarsGet**
> List<ServicesCalendarResponse> apiCalendarsGet()

Get user calendars

Retrieve all calendars for the authenticated user

### Example
```dart
import 'package:days_api_client/api.dart';
// TODO Configure API key authorization: BearerAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKeyPrefix = 'Bearer';

final api_instance = CalendarsApi();

try {
    final result = api_instance.apiCalendarsGet();
    print(result);
} catch (e) {
    print('Exception when calling CalendarsApi->apiCalendarsGet: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**List<ServicesCalendarResponse>**](ServicesCalendarResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiCalendarsIdDelete**
> apiCalendarsIdDelete(id)

Delete calendar

Delete a calendar and all associated data (user must own the calendar)

### Example
```dart
import 'package:days_api_client/api.dart';
// TODO Configure API key authorization: BearerAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKeyPrefix = 'Bearer';

final api_instance = CalendarsApi();
final id = id_example; // String | Calendar ID

try {
    api_instance.apiCalendarsIdDelete(id);
} catch (e) {
    print('Exception when calling CalendarsApi->apiCalendarsIdDelete: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| Calendar ID | 

### Return type

void (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiCalendarsIdGet**
> ServicesCalendarResponse apiCalendarsIdGet(id)

Get calendar by ID

Retrieve a specific calendar by ID (user must own the calendar)

### Example
```dart
import 'package:days_api_client/api.dart';
// TODO Configure API key authorization: BearerAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKeyPrefix = 'Bearer';

final api_instance = CalendarsApi();
final id = id_example; // String | Calendar ID

try {
    final result = api_instance.apiCalendarsIdGet(id);
    print(result);
} catch (e) {
    print('Exception when calling CalendarsApi->apiCalendarsIdGet: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| Calendar ID | 

### Return type

[**ServicesCalendarResponse**](ServicesCalendarResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiCalendarsIdPut**
> ServicesCalendarResponse apiCalendarsIdPut(id, calendar)

Update calendar

Update calendar name and description (user must own the calendar)

### Example
```dart
import 'package:days_api_client/api.dart';
// TODO Configure API key authorization: BearerAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKeyPrefix = 'Bearer';

final api_instance = CalendarsApi();
final id = id_example; // String | Calendar ID
final calendar = ServicesUpdateCalendarRequest(); // ServicesUpdateCalendarRequest | Calendar update request

try {
    final result = api_instance.apiCalendarsIdPut(id, calendar);
    print(result);
} catch (e) {
    print('Exception when calling CalendarsApi->apiCalendarsIdPut: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| Calendar ID | 
 **calendar** | [**ServicesUpdateCalendarRequest**](ServicesUpdateCalendarRequest.md)| Calendar update request | 

### Return type

[**ServicesCalendarResponse**](ServicesCalendarResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiCalendarsPost**
> ServicesCalendarResponse apiCalendarsPost(calendar)

Create a new calendar

Create a new calendar for the authenticated user

### Example
```dart
import 'package:days_api_client/api.dart';
// TODO Configure API key authorization: BearerAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('BearerAuth').apiKeyPrefix = 'Bearer';

final api_instance = CalendarsApi();
final calendar = ServicesCreateCalendarRequest(); // ServicesCreateCalendarRequest | Calendar creation request

try {
    final result = api_instance.apiCalendarsPost(calendar);
    print(result);
} catch (e) {
    print('Exception when calling CalendarsApi->apiCalendarsPost: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **calendar** | [**ServicesCreateCalendarRequest**](ServicesCreateCalendarRequest.md)| Calendar creation request | 

### Return type

[**ServicesCalendarResponse**](ServicesCalendarResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

