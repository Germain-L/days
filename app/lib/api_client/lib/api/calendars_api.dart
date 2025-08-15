//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//
// @dart=2.18

// ignore_for_file: unused_element, unused_import
// ignore_for_file: always_put_required_named_parameters_first
// ignore_for_file: constant_identifier_names
// ignore_for_file: lines_longer_than_80_chars

part of openapi.api;


class CalendarsApi {
  CalendarsApi([ApiClient? apiClient]) : apiClient = apiClient ?? defaultApiClient;

  final ApiClient apiClient;

  /// Get user calendars
  ///
  /// Retrieve all calendars for the authenticated user
  ///
  /// Note: This method returns the HTTP [Response].
  Future<Response> apiCalendarsGetWithHttpInfo() async {
    // ignore: prefer_const_declarations
    final path = r'/api/calendars';

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get user calendars
  ///
  /// Retrieve all calendars for the authenticated user
  Future<List<ServicesCalendarResponse>?> apiCalendarsGet() async {
    final response = await apiCalendarsGetWithHttpInfo();
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      final responseBody = await _decodeBodyBytes(response);
      return (await apiClient.deserializeAsync(responseBody, 'List<ServicesCalendarResponse>') as List)
        .cast<ServicesCalendarResponse>()
        .toList(growable: false);

    }
    return null;
  }

  /// Delete calendar
  ///
  /// Delete a calendar and all associated data (user must own the calendar)
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Calendar ID
  Future<Response> apiCalendarsIdDeleteWithHttpInfo(String id,) async {
    // ignore: prefer_const_declarations
    final path = r'/api/calendars/{id}'
      .replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'DELETE',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Delete calendar
  ///
  /// Delete a calendar and all associated data (user must own the calendar)
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Calendar ID
  Future<void> apiCalendarsIdDelete(String id,) async {
    final response = await apiCalendarsIdDeleteWithHttpInfo(id,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
  }

  /// Get calendar by ID
  ///
  /// Retrieve a specific calendar by ID (user must own the calendar)
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Calendar ID
  Future<Response> apiCalendarsIdGetWithHttpInfo(String id,) async {
    // ignore: prefer_const_declarations
    final path = r'/api/calendars/{id}'
      .replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>[];


    return apiClient.invokeAPI(
      path,
      'GET',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Get calendar by ID
  ///
  /// Retrieve a specific calendar by ID (user must own the calendar)
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Calendar ID
  Future<ServicesCalendarResponse?> apiCalendarsIdGet(String id,) async {
    final response = await apiCalendarsIdGetWithHttpInfo(id,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ServicesCalendarResponse',) as ServicesCalendarResponse;
    
    }
    return null;
  }

  /// Update calendar
  ///
  /// Update calendar name and description (user must own the calendar)
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Calendar ID
  ///
  /// * [ServicesUpdateCalendarRequest] calendar (required):
  ///   Calendar update request
  Future<Response> apiCalendarsIdPutWithHttpInfo(String id, ServicesUpdateCalendarRequest calendar,) async {
    // ignore: prefer_const_declarations
    final path = r'/api/calendars/{id}'
      .replaceAll('{id}', id);

    // ignore: prefer_final_locals
    Object? postBody = calendar;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'PUT',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Update calendar
  ///
  /// Update calendar name and description (user must own the calendar)
  ///
  /// Parameters:
  ///
  /// * [String] id (required):
  ///   Calendar ID
  ///
  /// * [ServicesUpdateCalendarRequest] calendar (required):
  ///   Calendar update request
  Future<ServicesCalendarResponse?> apiCalendarsIdPut(String id, ServicesUpdateCalendarRequest calendar,) async {
    final response = await apiCalendarsIdPutWithHttpInfo(id, calendar,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ServicesCalendarResponse',) as ServicesCalendarResponse;
    
    }
    return null;
  }

  /// Create a new calendar
  ///
  /// Create a new calendar for the authenticated user
  ///
  /// Note: This method returns the HTTP [Response].
  ///
  /// Parameters:
  ///
  /// * [ServicesCreateCalendarRequest] calendar (required):
  ///   Calendar creation request
  Future<Response> apiCalendarsPostWithHttpInfo(ServicesCreateCalendarRequest calendar,) async {
    // ignore: prefer_const_declarations
    final path = r'/api/calendars';

    // ignore: prefer_final_locals
    Object? postBody = calendar;

    final queryParams = <QueryParam>[];
    final headerParams = <String, String>{};
    final formParams = <String, String>{};

    const contentTypes = <String>['application/json'];


    return apiClient.invokeAPI(
      path,
      'POST',
      queryParams,
      postBody,
      headerParams,
      formParams,
      contentTypes.isEmpty ? null : contentTypes.first,
    );
  }

  /// Create a new calendar
  ///
  /// Create a new calendar for the authenticated user
  ///
  /// Parameters:
  ///
  /// * [ServicesCreateCalendarRequest] calendar (required):
  ///   Calendar creation request
  Future<ServicesCalendarResponse?> apiCalendarsPost(ServicesCreateCalendarRequest calendar,) async {
    final response = await apiCalendarsPostWithHttpInfo(calendar,);
    if (response.statusCode >= HttpStatus.badRequest) {
      throw ApiException(response.statusCode, await _decodeBodyBytes(response));
    }
    // When a remote server returns no body with a status of 204, we shall not decode it.
    // At the time of writing this, `dart:convert` will throw an "Unexpected end of input"
    // FormatException when trying to decode an empty string.
    if (response.body.isNotEmpty && response.statusCode != HttpStatus.noContent) {
      return await apiClient.deserializeAsync(await _decodeBodyBytes(response), 'ServicesCalendarResponse',) as ServicesCalendarResponse;
    
    }
    return null;
  }
}
