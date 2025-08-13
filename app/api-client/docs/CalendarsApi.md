# CalendarsApi

All URIs are relative to *http://localhost:8080*

| Method | HTTP request | Description |
| ------------- | ------------- | ------------- |
| [**apiCalendarsGet**](CalendarsApi.md#apiCalendarsGet) | **GET** api/calendars | Get user calendars |
| [**apiCalendarsIdDelete**](CalendarsApi.md#apiCalendarsIdDelete) | **DELETE** api/calendars/{id} | Delete calendar |
| [**apiCalendarsIdGet**](CalendarsApi.md#apiCalendarsIdGet) | **GET** api/calendars/{id} | Get calendar by ID |
| [**apiCalendarsIdPut**](CalendarsApi.md#apiCalendarsIdPut) | **PUT** api/calendars/{id} | Update calendar |
| [**apiCalendarsPost**](CalendarsApi.md#apiCalendarsPost) | **POST** api/calendars | Create a new calendar |



Get user calendars

Retrieve all calendars for the authenticated user

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(CalendarsApi::class.java)

launch(Dispatchers.IO) {
    val result : kotlin.collections.List<ServicesCalendarResponse> = webService.apiCalendarsGet()
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**kotlin.collections.List&lt;ServicesCalendarResponse&gt;**](ServicesCalendarResponse.md)

### Authorization



### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


Delete calendar

Delete a calendar and all associated data (user must own the calendar)

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(CalendarsApi::class.java)
val id : kotlin.String = id_example // kotlin.String | Calendar ID

launch(Dispatchers.IO) {
    webService.apiCalendarsIdDelete(id)
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **id** | **kotlin.String**| Calendar ID | |

### Return type

null (empty response body)

### Authorization



### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


Get calendar by ID

Retrieve a specific calendar by ID (user must own the calendar)

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(CalendarsApi::class.java)
val id : kotlin.String = id_example // kotlin.String | Calendar ID

launch(Dispatchers.IO) {
    val result : ServicesCalendarResponse = webService.apiCalendarsIdGet(id)
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **id** | **kotlin.String**| Calendar ID | |

### Return type

[**ServicesCalendarResponse**](ServicesCalendarResponse.md)

### Authorization



### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


Update calendar

Update calendar name and description (user must own the calendar)

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(CalendarsApi::class.java)
val id : kotlin.String = id_example // kotlin.String | Calendar ID
val calendar : ServicesUpdateCalendarRequest =  // ServicesUpdateCalendarRequest | Calendar update request

launch(Dispatchers.IO) {
    val result : ServicesCalendarResponse = webService.apiCalendarsIdPut(id, calendar)
}
```

### Parameters
| **id** | **kotlin.String**| Calendar ID | |
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **calendar** | [**ServicesUpdateCalendarRequest**](ServicesUpdateCalendarRequest.md)| Calendar update request | |

### Return type

[**ServicesCalendarResponse**](ServicesCalendarResponse.md)

### Authorization



### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


Create a new calendar

Create a new calendar for the authenticated user

### Example
```kotlin
// Import classes:
//import com.germainleignel.days.api.*
//import com.germainleignel.days.api.infrastructure.*
//import com.germainleignel.days.api.models.*

val apiClient = ApiClient()
val webService = apiClient.createWebservice(CalendarsApi::class.java)
val calendar : ServicesCreateCalendarRequest =  // ServicesCreateCalendarRequest | Calendar creation request

launch(Dispatchers.IO) {
    val result : ServicesCalendarResponse = webService.apiCalendarsPost(calendar)
}
```

### Parameters
| Name | Type | Description  | Notes |
| ------------- | ------------- | ------------- | ------------- |
| **calendar** | [**ServicesCreateCalendarRequest**](ServicesCreateCalendarRequest.md)| Calendar creation request | |

### Return type

[**ServicesCalendarResponse**](ServicesCalendarResponse.md)

### Authorization



### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

