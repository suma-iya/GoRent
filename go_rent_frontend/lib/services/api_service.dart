import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/property.dart';
import '../models/floor.dart';
import '../models/notification.dart' as models;
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:io';

class ApiService {
    static final ApiService _instance = ApiService._internal();
  factory ApiService() => _instance;
  ApiService._internal();

  // Detect base URL dynamically
  static String get baseUrl {
    String url;
    if (Platform.isAndroid) {
      // For Android devices, we need to detect emulator vs physical device
      // Emulators use 10.0.2.2 to access host machine
      // Physical devices need the Mac's actual IP address on the network
      
      // Try to detect emulator - emulators typically have "emulator" in device name
      // or we can check environment variables
      final androidSerial = Platform.environment['ANDROID_SERIAL'] ?? '';
      final isEmulator = androidSerial.contains('emulator') || 
                        androidSerial.startsWith('emulator-');
      
      if (isEmulator) {
        url = 'http://10.0.2.2:8081'; // Emulator: use special IP
      } else {
        // Physical device: use your Mac's IP address
        // IMPORTANT: Update this IP (192.168.0.232) to match your Mac's IP on your network
        // Find your Mac's IP with: ipconfig getifaddr en0
        url = 'http://192.168.0.232:8081'; // Physical device: use Mac's IP
      }
    } else if (Platform.isIOS) {
      // iOS Simulator can use localhost directly
      url = 'http://localhost:8081'; // Docker backend on port 8081
    } else {
      // Web or desktop (debug)
      url = 'http://localhost:8081'; // Docker backend on port 8081
    }
    print('Using base URL: $url (Android Serial: ${Platform.environment['ANDROID_SERIAL'] ?? 'not set'})');
    return url;
  }
  String? _sessionToken;
  final _client = http.Client();
  int? _currentUserId;

  int? get currentUserId => _currentUserId;

  Future<void> setSessionToken(String token) async {
    _sessionToken = token;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('session_token', token);
    print('Session token set: $_sessionToken');
  }

  Future<void> loadSessionToken() async {
    final prefs = await SharedPreferences.getInstance();
    _sessionToken = prefs.getString('session_token');
    print('Loaded session token: $_sessionToken');
  }

  Future<void> clearSessionToken() async {
    _sessionToken = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('session_token');
    print('Session token cleared');
  }

  Map<String, String> get _headers {
    final headers = {
      'Content-Type': 'application/json',
    };
    
    if (_sessionToken != null) {
      headers['Cookie'] = 'sessiontoken=$_sessionToken';
      print('Adding session token to headers: $_sessionToken');
    }
    
    return headers;
  }

  // LOGIN
  Future<Map<String, dynamic>?> login(String phoneNumber, String password) async {
    try {
      print('Attempting login with phone: $phoneNumber');
      final response = await _client.post(
        Uri.parse('$baseUrl/login'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'phone_number': phoneNumber,
          'password': password,
        }),
      );
      
      print('Login response status: ${response.statusCode}');
      print('Login response headers: ${response.headers}');
      print('Login response body: ${response.body}');

      if (response.statusCode == 200) {
        // Parse response body to get user ID
        final responseData = json.decode(response.body);
        final userId = responseData['user_id'];
        final userName = responseData['name'];
        
        print('Extracted user ID: $userId, name: $userName');
        
        // Extract session cookie
        final setCookie = response.headers['set-cookie'];
        print('Set-Cookie header: $setCookie');
        
        if (setCookie != null) {
          final sessionMatch = RegExp(r'sessiontoken=([^;]+)').firstMatch(setCookie);
          if (sessionMatch != null) {
            await setSessionToken(sessionMatch.group(1)!);
            return {
              'success': true,
              'user_id': userId,
              'name': userName,
            };
          }
        }
        print('No session token found in response');
      }
      return null;
    } catch (e) {
      print('Login error: $e');
      throw Exception('Login error: $e');
    }
  }

  // REGISTRATION
  Future<bool> register(String phoneNumber, String password, {required String name}) async {
    try {
      print('Attempting registration with phone: $phoneNumber');
      final response = await _client.post(
        Uri.parse('$baseUrl/register'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'phone_number': phoneNumber,
          'name': name,
          'password': password,
        }),
      );
      
      print('Registration response status: ${response.statusCode}');
      print('Registration response body: ${response.body}');
      
      return response.statusCode == 201;
    } catch (e) {
      print('Registration error: $e');
      throw Exception('Registration error: $e');
    }
  }

  Future<List<Property>> getProperties() async {
    try {
      print('Fetching properties with headers: $_headers');
      final response = await _client.get(
        Uri.parse('$baseUrl/properties'),
        headers: _headers,
      );

      print('Properties response status: ${response.statusCode}');
      print('Properties response headers: ${response.headers}');
      print('Properties response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic>? propertiesJson = data['properties'];
          return propertiesJson?.map((json) => Property.fromJson(json)).toList() ?? [];
        } else {
          throw Exception(data['message'] ?? 'Failed to load properties');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching properties: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<Property>> getTenantProperties() async {
    try {
      print('Fetching tenant properties with headers: $_headers');
      final response = await _client.get(
        Uri.parse('$baseUrl/property/tenant'),
        headers: _headers,
      );

      print('Tenant properties response status: ${response.statusCode}');
      print('Tenant properties response headers: ${response.headers}');
      print('Tenant properties response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic>? propertiesJson = data['properties'];
          return propertiesJson?.map((json) => Property.fromJson(json)).toList() ?? [];
        } else {
          throw Exception(data['message'] ?? 'Failed to load tenant properties');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching tenant properties: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<Property> getPropertyById(int id) async {
    try {
      print('Fetching property from: $baseUrl/property/$id');
      final response = await _client.get(
        Uri.parse('$baseUrl/property/$id'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          return Property.fromJson(data['property']);
        } else {
          throw Exception(data['message'] ?? 'Failed to load property');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching property: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> addProperty(String name, String address, {String? photoBase64}) async {
    try {
      print('Adding property to: $baseUrl/property');
      print('Adding session token to headers: $_sessionToken');

      final Map<String, dynamic> requestBody = {
        'name': name,
        'address': address,
      };
      
      if (photoBase64 != null) {
        requestBody['photo'] = photoBase64;
      }

      final response = await _client.post(
        Uri.parse('$baseUrl/property'),
        headers: _headers,
        body: json.encode(requestBody),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error adding property: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<Map<String, dynamic>> getPropertyDetails(int id) async {
    try {
      print('Fetching property details from: $baseUrl/property/$id');
      final response = await _client.get(
        Uri.parse('$baseUrl/property/$id'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final property = Property.fromJson(data['property']);
          final List<dynamic> floorsJson = data['floors'] ?? [];
          final floors = floorsJson.map((json) {
            // Check if there's a pending notification for this floor
            if (json['status'] == 'pending' && json['notification'] != null) {
              json['notification_id'] = json['notification']['id'];
              json['status'] = 'pending';
            }
            return Floor.fromJson(json);
          }).toList();
          return {
            'property': property,
            'floors': floors,
            'is_manager': data['is_manager'] ?? false,
          };
        } else {
          throw Exception(data['message'] ?? 'Failed to load property details');
        }
      } else if (response.statusCode == 401) {
        await clearSessionToken(); // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching property details: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> addFloor(int propertyId, String name, int rent) async {
    try {
      print('Adding floor to property: $propertyId');
      final response = await _client.post(
        Uri.parse('$baseUrl/property/$propertyId/floor'),
        headers: _headers,
        body: json.encode({
          'name': name,
          'rent': rent,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'];
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error adding floor: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> addTenantToFloor(int propertyId, int floorId, String tenantName, String phoneNumber) async {
    try {
      print('Adding tenant to floor: $floorId in property: $propertyId');
      final response = await _client.post(
        Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/tenant'),
        headers: _headers,
        body: json.encode({
          'name': tenantName,
          'phone_number': phoneNumber,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'];
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error adding tenant: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> removeTenantFromFloor(int propertyId, int floorId) async {
    try {
      print('Removing tenant from floor: $floorId in property: $propertyId');
      final response = await _client.delete(
        Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/tenant'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error removing tenant: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> updateFloor(int propertyId, int floorId, String name, int rent) async {
    try {
      print('Updating floor: $floorId in property: $propertyId');
      final response = await _client.put(
        Uri.parse('$baseUrl/property/$propertyId/floor/$floorId'),
        headers: _headers,
        body: json.encode({
          'name': name,
          'rent': rent,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'];
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error updating floor: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<int> getUserIdByPhone(String phoneNumber) async {
    try {
      print('Getting user ID for phone: $phoneNumber');
      final response = await _client.get(
        Uri.parse('$baseUrl/users/phones/$phoneNumber'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          return data['user_id'];
        } else {
          throw Exception(data['message'] ?? 'Failed to get user ID');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error getting user ID: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> sendTenantRequest(int propertyId, int floorId, String phoneNumber) async {
    try {
      print('Sending tenant request for property: $propertyId, floor: $floorId');
      final response = await _client.post(
       Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/request'),
        headers: _headers,
        body: json.encode({
          'phone_number': phoneNumber,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else if (response.statusCode == 409) {
        final Map<String, dynamic> data = json.decode(response.body);
        throw Exception(data['message'] ?? 'A pending request already exists for this floor');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error sending tenant request: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<String>> getUserPhones() async {
    try {
      print('Fetching user phones');
      final response = await _client.get(
        Uri.parse('$baseUrl/users/phones'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic> users = data['users'];
          return users.map((user) => user['phone_number'] as String).toList();
        } else {
          throw Exception(data['message'] ?? 'Failed to load user phones');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching user phones: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> sendNotification(int userId, String message) async {
    try {
      print('Sending notification to user: $userId');
      final response = await _client.post(
        Uri.parse('$baseUrl/notifications'),
        headers: _headers,
        body: json.encode({
          'user_id': userId,
          'message': message,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'];
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error sending notification: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> deleteNotification(int notificationId) async {
    try {
      print('Deleting notification: $notificationId');
      final response = await _client.delete(
        Uri.parse('$baseUrl/notifications/delete/$notificationId'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error deleting notification: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<models.AppNotification>> getNotifications() async {
    try {
      final response = await _client.get(
        Uri.parse('$baseUrl/notifications'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          // Handle null notifications by returning an empty list
          final notifications = data['notifications'] as List<dynamic>?;
          if (notifications == null) {
            return [];
          }
          return notifications.map((n) => models.AppNotification.fromJson(n)).toList();
        } else {
          throw Exception(data['message'] ?? 'Failed to load notifications');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching notifications: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> handleTenantRequestAction(int notificationId, bool accept) async {
    try {
      print('Handling tenant request action for notification: $notificationId');
      final response = await _client.post(
        Uri.parse('$baseUrl/notifications/action'),
        headers: _headers,
        body: json.encode({
          'notification_id': notificationId,
          'accept': accept,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error handling tenant request action: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<Property>> getUserProperties() async {
    try {
      print('Fetching user properties with headers: $_headers');
      final response = await _client.get(
        Uri.parse('$baseUrl/properties'),
        headers: _headers,
      );

      print('User properties response status: ${response.statusCode}');
      print('User properties response headers: ${response.headers}');
      print('User properties response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic>? propertiesJson = data['properties'];
          return propertiesJson?.map((json) => Property.fromJson(json)).toList() ?? [];
        } else {
          throw Exception(data['message'] ?? 'Failed to load properties');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching user properties: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<Property>> getUserTenantProperties() async {
    try {
      print('Fetching tenant properties with headers: $_headers');
      final response = await _client.get(
        Uri.parse('$baseUrl/properties/tenant'),
        headers: _headers,
      );

      print('Tenant properties response status: ${response.statusCode}');
      print('Tenant properties response headers: ${response.headers}');
      print('Tenant properties response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic> propertiesJson = data['properties'] ?? [];
          final properties = propertiesJson.map((json) => Property.fromJson(json)).toList();
          print('Parsed ${properties.length} tenant properties');
          return properties;
        } else {
          throw Exception(data['message'] ?? 'Failed to load tenant properties');
        }
      } else if (response.statusCode == 401) {
        await clearSessionToken(); // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching tenant properties: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<void> setCurrentUserId(int userId) async {
    _currentUserId = userId;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setInt('current_user_id', userId);
  }

  Future<void> loadCurrentUserId() async {
    final prefs = await SharedPreferences.getInstance();
    _currentUserId = prefs.getInt('current_user_id');
  }

  Future<void> sendPaymentNotification({
    required int propertyId,
    required int floorId,
    required int amount,
    int? electricityBill,
  }) async {
    final response = await _client.post(
      Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/payment-notification'),
      headers: _headers,
      body: json.encode({
        'amount': amount,
        if (electricityBill != null) 'paid_electricity_bill': electricityBill,
      }),
    );
    if (response.statusCode != 200) {
      throw Exception('Failed to send payment notification');
    }
  }

  Future<List<models.AppNotification>> getPendingPaymentNotifications(int propertyId, int floorId) async {
    try {
      print('Fetching pending payment notifications for property: $propertyId, floor: $floorId');
      final response = await _client.get(
        Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/pending-payments'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final notifications = data['notifications'] as List<dynamic>?;
          if (notifications == null) {
            return [];
          }
          return notifications.map((n) => models.AppNotification.fromJson(n)).toList();
        } else {
          throw Exception(data['message'] ?? 'Failed to load pending payment notifications');
        }
      } else if (response.statusCode == 401) {
        await clearSessionToken();
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching pending payment notifications: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> markNotificationsAsRead() async {
    try {
      print('Marking notifications as read');
      final response = await _client.post(
        Uri.parse('$baseUrl/notifications/mark-read'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        await clearSessionToken();
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error marking notifications as read: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> sendComment(int notificationId, String comment) async {
    try {
      print('Sending comment for notification: $notificationId');
      final response = await _client.post(
        Uri.parse('$baseUrl/notifications/send-comment'),
        headers: _headers,
        body: json.encode({
          'notification_id': notificationId,
          'comment': comment,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error sending comment: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  // Test push notification
  Future<bool> testPushNotification({
    required String title,
    required String body,
    Map<String, dynamic>? data,
  }) async {
    try {
      print('Testing push notification: $title - $body');
      final response = await _client.post(
        Uri.parse('$baseUrl/test/push-notification'),
        headers: _headers,
        body: json.encode({
          'title': title,
          'body': body,
          if (data != null) 'data': data,
        }),
      );

      print('Test notification response status: ${response.statusCode}');
      print('Test notification response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> responseData = json.decode(response.body);
        return responseData['success'] ?? false;
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error testing push notification: $e');
      throw Exception('Error: $e');
    }
  }

  // Update FCM token on server
  Future<bool> updateFcmToken(String fcmToken) async {
    try {
      print('Updating FCM token: $fcmToken');
      final response = await _client.post(
        Uri.parse('$baseUrl/user/fcm-token'),
        headers: _headers,
        body: json.encode({
          'fcm_token': fcmToken,
        }),
      );

      print('FCM token update response status: ${response.statusCode}');
      print('FCM token update response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error updating FCM token: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<Map<String, dynamic>>> getConversationHistory(int floorId) async {
    try {
      print('Fetching conversation history for floor: $floorId');
      final response = await _client.get(
        Uri.parse('$baseUrl/notifications/conversation?floor_id=$floorId'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final conversations = data['conversations'] as List<dynamic>?;
          if (conversations == null) {
            return [];
          }
          return conversations.cast<Map<String, dynamic>>();
        } else {
          throw Exception(data['message'] ?? 'Failed to load conversation history');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching conversation history: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> createAdvancePaymentRequest({
    required int propertyId,
    required int floorId,
    required int advanceUid,
    required int money,
  }) async {
    try {
      print('Creating advance payment request for property: $propertyId, floor: $floorId');
      final response = await _client.post(
        Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/advance-payment'),
        headers: _headers,
        body: json.encode({
          'advance_uid': advanceUid,
          'money': money,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error creating advance payment request: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> checkPendingAdvancePayment(int floorId) async {
    try {
      print('Checking pending advance payment for floor: $floorId');
      final response = await _client.get(
        Uri.parse('$baseUrl/floor/$floorId/advance-payment/check'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['has_pending'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error checking pending advance payment: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> cancelAdvancePaymentRequest(int floorId) async {
    try {
      print('Canceling advance payment request for floor: $floorId');
      final response = await _client.delete(
        Uri.parse('$baseUrl/floor/$floorId/advance-payment'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error canceling advance payment request: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<List<Map<String, dynamic>>> getAdvanceDetails(int floorId) async {
    try {
      print('Fetching advance details for floor: $floorId');
      final response = await _client.get(
        Uri.parse('$baseUrl/floor/$floorId/advance-details'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic> advancesJson = data['advances'] ?? [];
          return advancesJson.map((json) => Map<String, dynamic>.from(json)).toList();
        } else {
          throw Exception(data['message'] ?? 'Failed to load advance details');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching advance details: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<Map<String, dynamic>> getPaymentDetails(int floorId) async {
    try {
      print('Fetching payment details for floor: $floorId');
      final response = await _client.get(
        Uri.parse('$baseUrl/floor/$floorId/payment'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          return data['payment'] ?? {};
        } else {
          throw Exception(data['message'] ?? 'Failed to load payment details');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching payment details: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<Map<String, dynamic>> getPaymentHistory(int floorId, {int page = 1, int limit = 25}) async {
    try {
      print('Fetching payment history for floor: $floorId, page: $page, limit: $limit');
      final uri = Uri.parse('$baseUrl/floor/$floorId/payment-history').replace(
        queryParameters: {
          'page': page.toString(),
          'limit': limit.toString(),
        },
      );
      
      final response = await _client.get(uri, headers: _headers);

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        if (data['success']) {
          final List<dynamic> paymentsJson = data['payments'] ?? [];
          final List<Map<String, dynamic>> payments = paymentsJson.map((json) => Map<String, dynamic>.from(json)).toList();
          
          return {
            'payments': payments,
            'pagination': data['pagination'] ?? {},
          };
        } else {
          throw Exception(data['message'] ?? 'Failed to load payment history');
        }
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching payment history: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> createPayment({
    required int propertyId,
    required int floorId,
    required int dueRent,
    required int receivedMoney,
    required bool fullPayment,
    int? electricityBill,
  }) async {
    try {
      print('Creating payment for floor: $floorId');
      final response = await _client.post(
        Uri.parse('$baseUrl/property/$propertyId/floor/$floorId/payment'),
        headers: _headers,
        body: json.encode({
          'rent': dueRent,
          'received_money': receivedMoney,
          'full_payment': fullPayment,
          if (electricityBill != null) 'electricity_bill': electricityBill,
        }),
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 201) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['success'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        // Try to parse error message from response
        try {
          final Map<String, dynamic> errorData = json.decode(response.body);
          final errorMessage = errorData['message'] ?? 'Unknown error';
          throw Exception(errorMessage);
        } catch (parseError) {
          // If parsing fails, use generic error
          throw Exception('Server returned status code ${response.statusCode}');
        }
      }
    } catch (e) {
      print('Error creating payment: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  Future<bool> checkUserManager(int propertyId) async {
    try {
      print('Checking if user is manager for property: $propertyId');
      final response = await _client.get(
        Uri.parse('$baseUrl/property/$propertyId/manager'),
        headers: _headers,
      );

      print('Response status code: ${response.statusCode}');
      print('Response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = json.decode(response.body);
        return data['is_manager'] ?? false;
      } else if (response.statusCode == 401) {
        _sessionToken = null; // Clear invalid session token
        throw Exception('Session expired. Please login again.');
      } else {
        throw Exception('Server returned status code ${response.statusCode}');
      }
    } catch (e) {
      print('Error checking user manager status: $e');
      if (e.toString().contains('SocketException')) {
        throw Exception('Could not connect to server. Please check if the server is running and you have internet connection.');
      }
      throw Exception('Error: $e');
    }
  }

  void dispose() {
    _client.close();
  }
} 