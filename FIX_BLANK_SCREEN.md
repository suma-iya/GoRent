# Fixing Blank Screen Issue

## 🔍 Diagnosis

The app is installed and running, but showing a blank screen. This is usually caused by:

1. **Firebase initialization blocking** - The app might be stuck initializing Firebase
2. **Network connection issue** - App can't connect to backend
3. **Widget rendering error** - Silent error preventing UI from showing
4. **Hot reload issue** - App needs a full restart

## ✅ Solutions to Try

### Solution 1: Full Rebuild and Restart

```bash
# Stop the app
adb -s emulator-5554 shell am force-stop com.example.go_rent_frontend

# Clean and rebuild
cd go_rent_frontend
flutter clean
flutter pub get
flutter run -d emulator-5554
```

### Solution 2: Check if Backend is Accessible

The app needs to connect to `http://10.0.2.2:8081`. Verify:

```bash
# Test backend from your Mac
curl http://localhost:8081/test/fcm-connection-public

# Should return: {"message":"FCM connection successful","success":true}
```

### Solution 3: Check Flutter Logs

```bash
# View real-time Flutter logs
cd go_rent_frontend
flutter logs --device-id=emulator-5554

# Or use adb
adb -s emulator-5554 logcat | grep flutter
```

### Solution 4: Restart Emulator

Sometimes the emulator needs a restart:

1. **In Android Studio**: 
   - Go to `Tools` → `Device Manager`
   - Stop the emulator (click stop icon)
   - Start it again (click play icon)

2. **Or via command line**:
   ```bash
   adb -s emulator-5554 emu kill
   # Then start from Android Studio
   ```

### Solution 5: Check Firebase Configuration

The app initializes Firebase on startup. If Firebase config is missing or incorrect, it might cause a blank screen.

Check if `firebase_options.dart` exists:
```bash
ls go_rent_frontend/lib/firebase_options.dart
```

### Solution 6: Run with Debug Output

```bash
cd go_rent_frontend
flutter run -d emulator-5554 --verbose
```

Look for any error messages in the output.

## 🎯 Quick Fix Commands

Run these in order:

```bash
# 1. Force stop the app
adb -s emulator-5554 shell am force-stop com.example.go_rent_frontend

# 2. Clear app data (fresh start)
adb -s emulator-5554 shell pm clear com.example.go_rent_frontend

# 3. Rebuild and run
cd go_rent_frontend
flutter clean
flutter pub get
flutter run -d emulator-5554
```

## 🔧 If Still Blank

1. **Check if emulator has internet**:
   - Open browser in emulator
   - Try to visit a website
   - If no internet, restart emulator

2. **Check backend is running**:
   ```bash
   docker-compose ps
   # Should show both backend and mysql as "Up"
   ```

3. **Try running on a different device**:
   ```bash
   flutter devices
   flutter run -d chrome  # Try web version
   ```

4. **Check for specific errors**:
   ```bash
   adb -s emulator-5554 logcat | grep -i "error\|exception\|fatal"
   ```

## 📱 Expected Behavior

After running `flutter run`, you should see:
- Login screen with phone number and password fields
- "Go Rent" title in app bar
- Register button at the bottom

If you see a completely white/black screen, the app is likely stuck during initialization.

## 🚨 Common Causes

1. **Firebase initialization timeout** - Check Firebase config
2. **Network unreachable** - Backend not accessible from emulator
3. **Missing dependencies** - Run `flutter pub get`
4. **Build cache issues** - Run `flutter clean`

Try the solutions above and check the logs for specific error messages!


