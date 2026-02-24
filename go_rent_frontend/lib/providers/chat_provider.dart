import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:intl/intl.dart';
import '../services/api_service.dart';

class ChatMessage {
  final String text;
  final bool isUser;
  final String timestamp;
  final String? intent;
  final double? confidence;
  final List<String>? suggestedFollowups;

  ChatMessage({
    required this.text,
    required this.isUser,
    required this.timestamp,
    this.intent,
    this.confidence,
    this.suggestedFollowups,
  });

  Map<String, dynamic> toJson() => {
        'text': text,
        'isUser': isUser,
        'timestamp': timestamp,
        'intent': intent,
        'confidence': confidence,
        'suggestedFollowups': suggestedFollowups,
      };

  factory ChatMessage.fromJson(Map<String, dynamic> json) => ChatMessage(
        text: json['text'],
        isUser: json['isUser'],
        timestamp: json['timestamp'],
        intent: json['intent'],
        confidence: json['confidence']?.toDouble(),
        suggestedFollowups: json['suggestedFollowups'] != null
            ? List<String>.from(json['suggestedFollowups'])
            : null,
      );
}

class ChatProvider with ChangeNotifier {
  List<ChatMessage> _messages = [];
  bool _isLoading = false;
  String? _selectedTenantId;
  final ApiService _apiService = ApiService();

  List<ChatMessage> get messages => _messages;
  bool get isLoading => _isLoading;
  String? get selectedTenantId => _selectedTenantId;

  ChatProvider() {
    _loadHistory();
  }

  Future<void> _loadHistory() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final saved = prefs.getStringList('chat_history') ?? [];
      
      _messages = saved
          .map((json) => ChatMessage.fromJson(jsonDecode(json)))
          .toList();
      
      notifyListeners();
    } catch (e) {
      print('Error loading chat history: $e');
    }
  }

  Future<void> _saveHistory() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final jsonList = _messages.map((msg) => jsonEncode(msg.toJson())).toList();
      await prefs.setStringList('chat_history', jsonList);
    } catch (e) {
      print('Error saving chat history: $e');
    }
  }

  Future<void> sendMessage(String text, {String? tenantId}) async {
    // Add user message
    final userMessage = ChatMessage(
      text: text,
      isUser: true,
      timestamp: DateFormat('HH:mm').format(DateTime.now()),
    );
    
    _messages.add(userMessage);
    _isLoading = true;
    notifyListeners();

    try {
      // Get base URL - chatbot uses same backend (same port)
      // The backend handles both regular API and chatbot on the same server
      final baseUrl = ApiService.baseUrl; // Use same base URL as main API
      
      // Prepare request
      final url = Uri.parse('$baseUrl/chat');
      final requestBody = {
        'message': text,
        'tenant_id': tenantId ?? _selectedTenantId ?? '',
        'user_id': 'flutter_user',
      };
      final body = jsonEncode(requestBody);
      
      print('Chatbot: Sending request to $url');
      print('Chatbot: Request body: $body');

      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: body,
      ).timeout(const Duration(seconds: 30));
      
      print('Chatbot: Response status: ${response.statusCode}');
      print('Chatbot: Response body: ${response.body.substring(0, response.body.length > 200 ? 200 : response.body.length)}');

      if (response.statusCode == 200) {
        Map<String, dynamic> data;
        try {
          data = jsonDecode(response.body) as Map<String, dynamic>;
        } catch (e) {
          print('Error parsing JSON response: $e');
          _addErrorMessage('Invalid response format from server');
          return;
        }
        
        // Parse response text - handle both string and null
        String responseText = 'No response';
        if (data['response_text'] != null) {
          responseText = data['response_text'].toString();
        }
        
        print('Chatbot: Parsed response_text: ${responseText.substring(0, responseText.length > 50 ? 50 : responseText.length)}...');
        
        // Parse confidence - handle both number and null
        double? confidence;
        if (data['confidence'] != null) {
          if (data['confidence'] is double) {
            confidence = data['confidence'];
          } else if (data['confidence'] is int) {
            confidence = (data['confidence'] as int).toDouble();
          } else {
            confidence = double.tryParse(data['confidence'].toString());
          }
        }
        
        // Parse suggested followups
        List<String>? suggestedFollowups;
        if (data['suggested_followups'] != null) {
          try {
            suggestedFollowups = List<String>.from(data['suggested_followups']);
          } catch (e) {
            print('Error parsing suggested_followups: $e');
            suggestedFollowups = null;
          }
        }
        
        // Add bot response
        final botMessage = ChatMessage(
          text: responseText,
          isUser: false,
          timestamp: DateFormat('HH:mm').format(DateTime.now()),
          intent: data['intent']?.toString(),
          confidence: confidence,
          suggestedFollowups: suggestedFollowups,
        );
        
        _messages.add(botMessage);
        
        // Save to history (keep last 100 messages)
        if (_messages.length > 100) {
          _messages = _messages.sublist(_messages.length - 100);
        }
        
        await _saveHistory();
      } else {
        final errorBody = response.body.isNotEmpty 
            ? response.body.substring(0, response.body.length > 100 ? 100 : response.body.length)
            : 'Unknown error';
        _addErrorMessage('Error ${response.statusCode}: $errorBody');
      }
    } catch (e) {
      print('Chat error details: $e');
      _addErrorMessage('Network error: ${e.toString()}');
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void _addErrorMessage(String error) {
    _messages.add(ChatMessage(
      text: '❌ $error',
      isUser: false,
      timestamp: DateFormat('HH:mm').format(DateTime.now()),
    ));
    notifyListeners();
  }

  void clearHistory() {
    _messages.clear();
    _saveHistory();
    notifyListeners();
  }

  void setSelectedTenant(String tenantId) {
    _selectedTenantId = tenantId;
    notifyListeners();
  }

  Map<String, dynamic> getStatistics() {
    final total = _messages.length;
    final userMessages = _messages.where((m) => m.isUser).length;
    final botMessages = total - userMessages;
    
    final intents = <String, int>{};
    for (final msg in _messages) {
      if (msg.intent != null) {
        intents[msg.intent!] = (intents[msg.intent!] ?? 0) + 1;
      }
    }

    return {
      'total_messages': total,
      'user_messages': userMessages,
      'bot_messages': botMessages,
      'intents': intents,
      'first_message': _messages.isNotEmpty ? _messages.first.timestamp : null,
      'last_message': _messages.isNotEmpty ? _messages.last.timestamp : null,
    };
  }
}

