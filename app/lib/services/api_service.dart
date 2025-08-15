import '../api_client/lib/api.dart';

class ApiService {
  static ApiService? _instance;
  late final ApiClient _apiClient;
  late final AuthApi _authApi;
  late final CalendarsApi _calendarsApi;
  late final UsersApi _usersApi;

  String? _authToken;

  // Singleton pattern
  static ApiService get instance {
    _instance ??= ApiService._internal();
    return _instance!;
  }

  ApiService._internal() {
    _initializeApiClient();
  }

  void _initializeApiClient() {
    // Initialize the API client with your backend URL
    _apiClient = ApiClient(basePath: 'http://localhost:8080');

    // Initialize API endpoints
    _authApi = AuthApi(_apiClient);
    _calendarsApi = CalendarsApi(_apiClient);
    _usersApi = UsersApi(_apiClient);
  }

  // Update the base URL if needed
  void updateBaseUrl(String baseUrl) {
    _apiClient = ApiClient(basePath: baseUrl);
    _authApi = AuthApi(_apiClient);
    _calendarsApi = CalendarsApi(_apiClient);
    _usersApi = UsersApi(_apiClient);
  }

  // Set authentication token
  void setAuthToken(String token) {
    _authToken = token;
    _apiClient.addDefaultHeader('Authorization', 'Bearer $token');
  }

  // Clear authentication token
  void clearAuthToken() {
    _authToken = null;
    // Reinitialize the client to clear the authorization header
    _initializeApiClient();
  }

  // Check if user is authenticated
  bool get isAuthenticated => _authToken != null;

  // Get the auth token
  String? get authToken => _authToken;

  // Auth API methods
  Future<ServicesLoginResponse?> login(String email, String password) async {
    try {
      final request = ServicesLoginRequest(email: email, password: password);

      final response = await _authApi.apiAuthLoginPost(request);

      if (response?.token != null) {
        setAuthToken(response!.token!);
      }

      return response;
    } catch (e) {
      rethrow;
    }
  }

  // User API methods
  Future<ServicesUserResponse?> createUser(
    String email,
    String password,
  ) async {
    try {
      final request = ServicesCreateUserRequest(
        email: email,
        password: password,
      );

      return await _usersApi.apiUsersPost(request);
    } catch (e) {
      rethrow;
    }
  }

  Future<ServicesUserResponse?> getUserById(String userId) async {
    try {
      return await _usersApi.apiUsersIdGet(userId);
    } catch (e) {
      rethrow;
    }
  }

  // Calendar API methods
  Future<List<ServicesCalendarResponse>?> getUserCalendars() async {
    try {
      return await _calendarsApi.apiCalendarsGet();
    } catch (e) {
      rethrow;
    }
  }

  Future<ServicesCalendarResponse?> createCalendar(
    String name,
    String? description,
  ) async {
    try {
      final request = ServicesCreateCalendarRequest(
        name: name,
        description: description,
      );

      return await _calendarsApi.apiCalendarsPost(request);
    } catch (e) {
      rethrow;
    }
  }

  Future<ServicesCalendarResponse?> getCalendarById(String calendarId) async {
    try {
      return await _calendarsApi.apiCalendarsIdGet(calendarId);
    } catch (e) {
      rethrow;
    }
  }

  Future<ServicesCalendarResponse?> updateCalendar(
    String calendarId,
    String name,
    String? description,
  ) async {
    try {
      final request = ServicesUpdateCalendarRequest(
        name: name,
        description: description,
      );

      return await _calendarsApi.apiCalendarsIdPut(calendarId, request);
    } catch (e) {
      rethrow;
    }
  }

  Future<void> deleteCalendar(String calendarId) async {
    try {
      await _calendarsApi.apiCalendarsIdDelete(calendarId);
    } catch (e) {
      rethrow;
    }
  }
}
