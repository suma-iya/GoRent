# Chatbot Troubleshooting Guide

## Common Issues and Fixes

### Issue 1: Text Not Displaying / Encoding Errors
**Symptoms:** Text shows as "â¢" or other garbled characters

**Fix Applied:**
- Replaced bullet points (•) with dashes (-) in backend
- Added text cleaning in Flutter to fix encoding issues
- Fallback to plain text if markdown fails

**If still seeing issues:**
1. Hot restart the app (not just hot reload)
2. Clear chat history (tap delete icon in chatbot)
3. Check console logs for errors

### Issue 2: Markdown Not Rendering
**Symptoms:** Text appears as raw markdown instead of formatted

**Fix Applied:**
- Added proper MarkdownStyleSheet configuration
- Added error handling with fallback to SelectableText
- Added shrinkWrap for proper layout

**If still seeing issues:**
1. Verify `flutter_markdown` is installed: `flutter pub get`
2. Check if message.text contains valid markdown
3. Look for errors in console starting with "Markdown rendering error"

### Issue 3: Messages Not Appearing
**Symptoms:** Messages sent but not showing in chat

**Fix Applied:**
- Added proper error handling in sendMessage
- Added debug logging
- Added bounds checking in ListView

**If still seeing issues:**
1. Check console for "Chatbot:" log messages
2. Verify backend is running: `curl http://localhost:8081/chat/health`
3. Check network connectivity

### Issue 4: App Crashes When Opening Chatbot
**Symptoms:** App crashes when tapping "AI Chatbot" in menu

**Possible Causes:**
- Missing import for ChatScreen
- Provider not registered
- Missing dependencies

**Fixes Applied:**
- Added ChatProvider to main.dart providers
- Added import for chat_screen.dart in properties_screen.dart
- Added all required dependencies to pubspec.yaml

## Debug Steps

1. **Check Backend:**
   ```bash
   curl -X POST http://localhost:8081/chat \
     -H "Content-Type: application/json" \
     -d '{"message":"test"}'
   ```

2. **Check Flutter Logs:**
   - Look for "Chatbot:" prefixed messages
   - Check for any red error messages
   - Verify base URL is correct

3. **Clear and Restart:**
   ```bash
   # Clear Flutter build cache
   cd go_rent_frontend
   flutter clean
   flutter pub get
   flutter run
   ```

4. **Verify Dependencies:**
   ```bash
   cd go_rent_frontend
   flutter pub get
   # Should install: flutter_markdown, url_launcher
   ```

## Expected Behavior

✅ **Working:**
- Messages appear in chat bubbles
- Markdown formatting (bold, lists) renders correctly
- Suggested follow-ups appear as chips
- Intent and confidence show below bot messages
- Text is selectable and readable

❌ **Not Working:**
- Garbled characters (encoding issue)
- Raw markdown showing (rendering issue)
- Messages not appearing (network/parsing issue)
- App crashes (missing imports/dependencies)

## Quick Fixes

### If text shows encoding errors:
The backend now uses dashes (-) instead of bullets. If you still see "â¢", the text cleaning should fix it. If not, restart the app.

### If markdown doesn't render:
The code falls back to plain text automatically. Check console for "Markdown rendering error" to see why.

### If messages don't send:
Check the console for "Chatbot: Sending request" and "Chatbot: Response status" messages to debug.

## Still Having Issues?

Please provide:
1. Exact error message from console
2. What happens when you try to use the chatbot
3. Screenshot if possible (though I can't see it, it helps describe the issue)

