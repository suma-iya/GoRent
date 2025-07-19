import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter/material.dart';
import 'dart:convert';
import 'api_service.dart';
import '../firebase_options.dart';
import '../main.dart';

class NotificationService {
  static final FlutterLocalNotificationsPlugin _localNotifications =
      FlutterLocalNotificationsPlugin();
  static final FirebaseMessaging _firebaseMessaging = FirebaseMessaging.instance;
  
  static String? _fcmToken;
  static bool _isInitialized = false;
  static RemoteMessage? _pendingInitialMessage;

  // Getter for FCM token
  static String? get fcmToken => _fcmToken;

  // Initialize notification service
  static Future<void> initialize() async {
    if (_isInitialized) return;

    try {
      // Initialize Firebase
      await Firebase.initializeApp(
        options: DefaultFirebaseOptions.currentPlatform,
      );

      // Request permission for iOS
      NotificationSettings settings = await _firebaseMessaging.requestPermission(
        alert: true,
        announcement: false,
        badge: true,
        carPlay: false,
        criticalAlert: false,
        provisional: false,
        sound: true,
      );

      print('User granted permission: ${settings.authorizationStatus}');

      // Get FCM token
      _fcmToken = await _firebaseMessaging.getToken();
      print('FCM Token: $_fcmToken');

      // Listen for token refresh
      _firebaseMessaging.onTokenRefresh.listen((token) {
        _fcmToken = token;
        print('FCM Token refreshed: $_fcmToken');
        // Here you can send the new token to your backend
        _sendTokenToServer(token);
      });

      // Initialize local notifications
      await _initializeLocalNotifications();

      // Handle background messages
      FirebaseMessaging.onBackgroundMessage(_firebaseMessagingBackgroundHandler);

      // Handle foreground messages
      FirebaseMessaging.onMessage.listen(_handleForegroundMessage);

      // Handle notification tap when app is in background
      FirebaseMessaging.onMessageOpenedApp.listen(_handleNotificationTap);

      // Check if app was opened from notification
      RemoteMessage? initialMessage = await _firebaseMessaging.getInitialMessage();
      if (initialMessage != null) {
        print('App opened from notification: ${initialMessage.data}');
        // Store the initial message to handle after app is fully loaded
        _pendingInitialMessage = initialMessage;
        // Delay navigation to ensure app is fully initialized
        Future.delayed(const Duration(seconds: 3), () {
          if (_pendingInitialMessage != null) {
            _handleNotificationTap(_pendingInitialMessage!);
            _pendingInitialMessage = null;
          }
        });
      }

      _isInitialized = true;
      
      // Test local notification after initialization
      await Future.delayed(Duration(seconds: 2));
      await showCustomNotification(
        title: 'Notification Test',
        body: 'Local notifications are working!',
        payload: '{"type": "test"}',
      );
      
    } catch (e) {
      print('Error initializing notification service: $e');
      
      // Even if Firebase fails, try to initialize local notifications
      try {
        await _initializeLocalNotifications();
        await showCustomNotification(
          title: 'Local Notification Test',
          body: 'Local notifications work even without Firebase!',
          payload: '{"type": "local_test"}',
        );
      } catch (localError) {
        print('Error initializing local notifications: $localError');
      }
    }
  }

  // Initialize local notifications
  static Future<void> _initializeLocalNotifications() async {
    const AndroidInitializationSettings initializationSettingsAndroid =
        AndroidInitializationSettings('@mipmap/ic_launcher');

    const DarwinInitializationSettings initializationSettingsIOS =
        DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );

    const InitializationSettings initializationSettings =
        InitializationSettings(
      android: initializationSettingsAndroid,
      iOS: initializationSettingsIOS,
    );

    await _localNotifications.initialize(
      initializationSettings,
      onDidReceiveNotificationResponse: _onNotificationTapped,
    );

    // Create notification channel for Android
    const AndroidNotificationChannel channel = AndroidNotificationChannel(
      'high_importance_channel',
      'High Importance Notifications',
      description: 'This channel is used for important notifications.',
      importance: Importance.high,
    );

    await _localNotifications
        .resolvePlatformSpecificImplementation<
            AndroidFlutterLocalNotificationsPlugin>()
        ?.createNotificationChannel(channel);
  }

  // Handle local notification tap
  static void _onNotificationTapped(NotificationResponse response) {
    print('Local notification tapped: ${response.payload}');
    
    if (response.payload != null) {
      Map<String, dynamic> data = json.decode(response.payload!);
      _navigateBasedOnNotification(data);
    }
  }

  // Navigate based on notification data
  static void _navigateBasedOnNotification(Map<String, dynamic> data) {
    print('Navigate based on notification type: ${data['type']}');
    print('Notification data: $data');
    
    // Check if user is logged in (has navigator context)
    if (navigatorKey.currentContext == null) {
      print('No navigator context available, cannot navigate');
      return;
    }
    
    print('Navigator context available, current route: ${ModalRoute.of(navigatorKey.currentContext!)?.settings.name}');
    
    try {
      // Get current route
      final currentRoute = ModalRoute.of(navigatorKey.currentContext!);
      final currentRouteName = currentRoute?.settings.name;
      
      print('Current route name: $currentRouteName');
      
      // If we're on login or register screen, don't navigate (user not logged in)
      if (currentRouteName == '/login' || currentRouteName == '/register') {
        print('User not logged in, skipping notification navigation');
        return;
      }
      
      // If we're already on notifications screen, do nothing
      if (currentRouteName == '/notifications') {
        print('Already on notifications screen');
        return;
      }
      
      // If we're on home screen, navigate directly to notifications
      if (currentRouteName == '/home') {
        print('On home screen, navigating directly to notifications');
        _navigateToNotificationScreen();
        return;
      }
      
      // If we're not on home screen and not on login/register, navigate to home first
      if (currentRouteName != '/login' && currentRouteName != '/register') {
        print('Not on home screen, navigating to home first');
        navigatorKey.currentState?.pushReplacementNamed('/home');
        // Then navigate to notifications after a short delay
        Future.delayed(const Duration(milliseconds: 300), () {
          _navigateToNotificationScreen();
        });
        return;
      }
    } catch (e) {
      print('Error navigating from notification: $e');
    }
  }
  
  // Helper method to navigate to notification screen
  static void _navigateToNotificationScreen() {
    try {
      // Check if we're already on notifications screen
      final currentRoute = ModalRoute.of(navigatorKey.currentContext!);
      if (currentRoute?.settings.name == '/notifications') {
        print('Already on notifications screen, skipping navigation');
        return;
      }
      
      print('Navigating to notifications screen');
      navigatorKey.currentState?.pushNamed('/notifications');
    } catch (e) {
      print('Error navigating to notification screen: $e');
    }
  }
  
  // Public method to handle pending initial message after app is ready
  static void handlePendingInitialMessage() {
    if (_pendingInitialMessage != null) {
      print('Handling pending initial message');
      _handleNotificationTap(_pendingInitialMessage!);
      _pendingInitialMessage = null;
    }
  }

  // Callback to refresh notification count when returning from notifications screen
  static Function? _onNotificationCountChanged;

  // Set callback for notification count changes
  static void setNotificationCountCallback(Function callback) {
    _onNotificationCountChanged = callback;
  }

  // Clear callback
  static void clearNotificationCountCallback() {
    _onNotificationCountChanged = null;
  }

  // Notify that notification count should be refreshed
  static void notifyNotificationCountChanged() {
    if (_onNotificationCountChanged != null) {
      _onNotificationCountChanged!();
    }
  }

  // Request notification permission
  static Future<bool> requestPermission() async {
    try {
      NotificationSettings settings = await _firebaseMessaging.requestPermission(
        alert: true,
        announcement: false,
        badge: true,
        carPlay: false,
        criticalAlert: false,
        provisional: false,
        sound: true,
      );
      
      print('Permission request result: ${settings.authorizationStatus}');
      return settings.authorizationStatus == AuthorizationStatus.authorized;
    } catch (e) {
      print('Error requesting permission: $e');
      return false;
    }
  }

  // Clear FCM token from server
  static Future<void> clearTokenFromServer() async {
    try {
      final apiService = ApiService();
      await apiService.loadSessionToken();
      
      // Send empty token to server to disable notifications
      await apiService.updateFcmToken('');
      print('FCM token cleared from server');
    } catch (e) {
      print('Error clearing FCM token from server: $e');
      throw e;
    }
  }
  
  // Public method to clear pending initial message
  static void clearPendingInitialMessage() {
    if (_pendingInitialMessage != null) {
      print('Clearing pending initial message');
      _pendingInitialMessage = null;
    }
  }

  // Handle background messages
  static Future<void> _firebaseMessagingBackgroundHandler(RemoteMessage message) async {
    await Firebase.initializeApp(
      options: DefaultFirebaseOptions.currentPlatform,
    );
    print('Handling a background message: ${message.messageId}');
    
    // Show local notification for background messages
    await _showLocalNotification(message);
  }

  // Handle foreground messages
  static void _handleForegroundMessage(RemoteMessage message) {
    print('Got a message whilst in the foreground!');
    print('Message data: ${message.data}');

    if (message.notification != null) {
      print('Message also contained a notification: ${message.notification}');
      _showLocalNotification(message);
    }
  }

  // Handle notification tap
  static void _handleNotificationTap(RemoteMessage message) {
    print('Notification tapped: ${message.data}');
    
    // Add a small delay to ensure the app is fully loaded
    Future.delayed(const Duration(milliseconds: 100), () {
      // Navigate based on notification type
      _navigateBasedOnNotification(message.data);
    });
  }

  // Show local notification
  static Future<void> _showLocalNotification(RemoteMessage message) async {
    const AndroidNotificationDetails androidPlatformChannelSpecifics =
        AndroidNotificationDetails(
      'high_importance_channel',
      'High Importance Notifications',
      channelDescription: 'This channel is used for important notifications.',
      importance: Importance.max,
      priority: Priority.high,
      showWhen: true,
    );

    const DarwinNotificationDetails iOSPlatformChannelSpecifics =
        DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );

    const NotificationDetails platformChannelSpecifics = NotificationDetails(
      android: androidPlatformChannelSpecifics,
      iOS: iOSPlatformChannelSpecifics,
    );

    await _localNotifications.show(
      message.hashCode,
      message.notification?.title ?? 'Go Rent',
      message.notification?.body ?? '',
      platformChannelSpecifics,
      payload: json.encode(message.data),
    );
  }

  // Send FCM token to your backend server
  static Future<void> _sendTokenToServer(String token) async {
    try {
      await ApiService().updateFcmToken(token);
      print('FCM token sent to server successfully: $token');
    } catch (e) {
      print('Error sending FCM token to server: $e');
    }
  }

  // Public method to send current FCM token to server
  static Future<void> sendCurrentTokenToServer() async {
    if (_fcmToken != null) {
      await _sendTokenToServer(_fcmToken!);
    } else {
      print('FCM token not available yet');
    }
  }

  // Subscribe to topics
  static Future<void> subscribeToTopic(String topic) async {
    await _firebaseMessaging.subscribeToTopic(topic);
    print('Subscribed to topic: $topic');
  }

  // Unsubscribe from topics
  static Future<void> unsubscribeFromTopic(String topic) async {
    await _firebaseMessaging.unsubscribeFromTopic(topic);
    print('Unsubscribed from topic: $topic');
  }

  // Get current FCM token
  static Future<String?> getToken() async {
    return await _firebaseMessaging.getToken();
  }

  // Clear all notifications
  static Future<void> clearAllNotifications() async {
    await _localNotifications.cancelAll();
  }

  // Show custom local notification
  static Future<void> showCustomNotification({
    required String title,
    required String body,
    String? payload,
  }) async {
    const AndroidNotificationDetails androidPlatformChannelSpecifics =
        AndroidNotificationDetails(
      'high_importance_channel',
      'High Importance Notifications',
      channelDescription: 'This channel is used for important notifications.',
      importance: Importance.max,
      priority: Priority.high,
      showWhen: true,
    );

    const DarwinNotificationDetails iOSPlatformChannelSpecifics =
        DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );

    const NotificationDetails platformChannelSpecifics = NotificationDetails(
      android: androidPlatformChannelSpecifics,
      iOS: iOSPlatformChannelSpecifics,
    );

    await _localNotifications.show(
      DateTime.now().millisecondsSinceEpoch ~/ 1000,
      title,
      body,
      platformChannelSpecifics,
      payload: payload,
    );
  }
} 