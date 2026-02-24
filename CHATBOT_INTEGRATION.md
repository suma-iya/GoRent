# Chatbot Integration Guide

## Overview
The AI Chatbot feature has been successfully integrated into the Go Rent app. It provides intelligent risk analysis and recommendations for tenant management.

## What Was Added

### 1. Flutter Frontend
- **Chat Screen** (`lib/screens/chat_screen.dart`)
  - Beautiful chat interface with markdown support
  - Quick action buttons for common queries
  - Message history persistence
  - Suggested follow-up questions

- **Chat Provider** (`lib/providers/chat_provider.dart`)
  - State management for chat messages
  - API communication with backend
  - Message history storage using SharedPreferences

- **Menu Integration**
  - Added "AI Chatbot" option in the drawer menu (between Notifications and Settings)
  - Purple icon (Icons.smart_toy_rounded) for easy identification

### 2. Go Backend
- **Chatbot Handler** (`handlers/chatbot.go`)
  - Intent detection using regex patterns
  - Response generation based on tenant risk data
  - Sample tenant data (T100, T087, T200, T050)
  - Support for multiple intents:
    - EXPLAIN_RISK
    - RECOMMEND_ACTION
    - LIST_HIGH_RISK
    - MONTHLY_SUMMARY
    - COMPARE_TENANTS
    - PAYMENT_HISTORY
    - LEASE_RENEWAL

- **Routes Added** (`main.go`)
  - `POST /chat` - Process chat messages
  - `GET /chat/health` - Health check for chatbot service

### 3. Dependencies
- Added `flutter_markdown: ^0.6.18` for markdown rendering
- Added `url_launcher: ^6.2.2` for link handling

## How to Use

### For Users
1. Open the app and navigate to the drawer menu (hamburger icon)
2. Tap on "AI Chatbot" (purple robot icon)
3. Start asking questions like:
   - "Why is T100 high risk?"
   - "What should I do for T087?"
   - "List high risk tenants"
   - "Show monthly summary"
   - "Compare T100 and T087"

### For Developers

#### Running the Backend
The chatbot routes are automatically available when you start the Go backend:
```bash
# If using Docker
docker-compose up -d backend

# If running directly
go run main.go
```

The chatbot uses the same backend server and port (8081 from Docker, or 8080 directly).

#### Testing the Chatbot
```bash
# Test health endpoint
curl http://localhost:8081/chat/health

# Test chat endpoint
curl -X POST http://localhost:8081/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Why is T100 high risk?", "user_id": "test_user"}'
```

#### Flutter Development
```bash
cd go_rent_frontend
flutter pub get  # Install new dependencies
flutter run      # Run the app
```

## Features

### Intent Detection
The chatbot uses pattern matching to detect user intents:
- Risk explanations
- Action recommendations
- Tenant comparisons
- Summary reports
- Payment history queries

### Sample Tenants
The chatbot comes with 4 sample tenants:
- **T100**: High risk (0.78) - 3 late payments
- **T087**: Medium risk (0.45) - 1 late payment
- **T200**: Critical risk (0.92) - 5 late payments
- **T050**: Low risk (0.15) - Excellent history

### Response Features
- Markdown formatting support
- Suggested follow-up questions
- Intent confidence scores
- Processing time metrics

## Customization

### Adding More Tenants
Edit `handlers/chatbot.go` in the `initializeSampleData()` function to add more sample tenants.

### Adding New Intents
1. Add regex patterns to `IntentDetector` in `handlers/chatbot.go`
2. Add a case in the `ProcessMessage` switch statement
3. Create a response generation method

### Changing UI
Modify `lib/screens/chat_screen.dart` to customize:
- Colors and themes
- Message bubble styles
- Quick action buttons
- Layout and spacing

## Notes

- The chatbot currently uses sample data. To integrate with real tenant data, modify `handlers/chatbot.go` to fetch from your database.
- Chat history is stored locally using SharedPreferences (last 100 messages).
- The chatbot is accessible without authentication (public route). To make it protected, move the routes to `protectedRouter` in `main.go`.

## Troubleshooting

### Chatbot not responding
1. Check if backend is running: `curl http://localhost:8081/chat/health`
2. Check Flutter logs for API errors
3. Verify base URL in `ApiService.baseUrl` matches your backend

### Messages not saving
- Check SharedPreferences permissions
- Verify `flutter pub get` installed all dependencies

### Markdown not rendering
- Ensure `flutter_markdown` is in `pubspec.yaml`
- Run `flutter pub get`

## Future Enhancements

- [ ] Integrate with real tenant database
- [ ] Add authentication requirement
- [ ] Support for multiple languages
- [ ] Voice input support
- [ ] Export chat history
- [ ] Analytics and insights dashboard

