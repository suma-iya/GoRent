import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:shared_preferences/shared_preferences.dart';

class LocalizationService extends ChangeNotifier {
  static const String _languageKey = 'selected_language';
  static const String _defaultLanguage = 'en';
  
  Locale _currentLocale = const Locale('en');
  
  Locale get currentLocale => _currentLocale;
  
  static final List<Locale> supportedLocales = [
    const Locale('en'), // English
    const Locale('bn'), // Bangla
  ];
  
  static const List<LocalizationsDelegate<dynamic>> localizationsDelegates = [
    GlobalMaterialLocalizations.delegate,
    GlobalWidgetsLocalizations.delegate,
    GlobalCupertinoLocalizations.delegate,
  ];
  
  static const List<String> supportedLanguageCodes = ['en', 'bn'];
  
  LocalizationService() {
    _loadSavedLanguage();
  }
  
  Future<void> _loadSavedLanguage() async {
    final prefs = await SharedPreferences.getInstance();
    final savedLanguage = prefs.getString(_languageKey) ?? _defaultLanguage;
    await changeLanguage(savedLanguage);
  }
  
  Future<void> changeLanguage(String languageCode) async {
    if (!supportedLanguageCodes.contains(languageCode)) {
      languageCode = _defaultLanguage;
    }
    
    _currentLocale = Locale(languageCode);
    
    // Save to preferences
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_languageKey, languageCode);
    
    notifyListeners();
  }
  
  String getLanguageName(String languageCode) {
    switch (languageCode) {
      case 'en':
        return 'English';
      case 'bn':
        return 'বাংলা';
      default:
        return 'English';
    }
  }
  
  String getCurrentLanguageName() {
    return getLanguageName(_currentLocale.languageCode);
  }
} 