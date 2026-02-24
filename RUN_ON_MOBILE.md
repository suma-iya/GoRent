# Commands to Run App on Mobile Device

## 🛑 Stop Current App

```bash
# Stop Flutter process
pkill -f "flutter run"

# Or press 'q' in the terminal where flutter run is active

# Force stop the app on device
adb -s BQ4X6595TWZ5DUYH shell am force-stop com.example.go_rent_frontend
```

## 📱 Run App on Physical Device

### Step 1: Check Device is Connected
```bash
flutter devices
```

You should see:
```
RMX3085 (mobile) • BQ4X6595TWZ5DUYH • android-arm64 • Android 13 (API 33)
```

### Step 2: Ensure Backend is Running
```bash
docker-compose ps
```

Both `rent-backend` and `rent-mysql` should show as "Up"

### Step 3: Run the App
```bash
cd go_rent_frontend
flutter run -d BQ4X6595TWZ5DUYH
```

**OR** use device name:
```bash
cd go_rent_frontend
flutter run -d RMX3085
```

## 🚀 Quick One-Liner

```bash
cd go_rent_frontend && flutter run -d BQ4X6595TWZ5DUYH
```

## 🔧 Alternative Commands

### Run in Release Mode (faster, no hot reload)
```bash
cd go_rent_frontend
flutter run -d BQ4X6595TWZ5DUYH --release
```

### Run with Verbose Output (for debugging)
```bash
cd go_rent_frontend
flutter run -d BQ4X6595TWZ5DUYH --verbose
```

### Clean Build and Run
```bash
cd go_rent_frontend
flutter clean
flutter pub get
flutter run -d BQ4X6595TWZ5DUYH
```

## 🛑 Stop the App

### While Running in Terminal
- Press `q` to quit
- Or press `Ctrl+C`

### Force Stop from Command Line
```bash
# Stop Flutter process
pkill -f "flutter run"

# Force stop app on device
adb -s BQ4X6595TWZ5DUYH shell am force-stop com.example.go_rent_frontend
```

### Uninstall App
```bash
adb -s BQ4X6595TWZ5DUYH uninstall com.example.go_rent_frontend
```

## ✅ Verify Setup

### Check Backend is Running
```bash
curl http://localhost:8081/test/fcm-connection-public
```

### Check Device Connection
```bash
adb devices
```

### Check Flutter Devices
```bash
flutter devices
```

## 📝 Important Notes

1. **Same Wi-Fi Network**: Your phone and Mac must be on the same Wi-Fi network
2. **Backend URL**: App is configured to use `http://192.168.0.232:8081` (your Mac's IP)
3. **If IP Changes**: Update IP in `go_rent_frontend/lib/services/api_service.dart` line ~25

## 🔍 Troubleshooting

### Device Not Found
```bash
# Enable USB debugging on phone
# Settings → Developer Options → USB Debugging

# Check connection
adb devices
```

### Can't Connect to Backend
1. Check Mac's IP: `ipconfig getifaddr en0`
2. Test from phone browser: `http://192.168.0.232:8081/test/fcm-connection-public`
3. Check firewall settings on Mac

### App Won't Install
```bash
# Uninstall old version
adb -s BQ4X6595TWZ5DUYH uninstall com.example.go_rent_frontend

# Clean and rebuild
cd go_rent_frontend
flutter clean
flutter pub get
flutter run -d BQ4X6595TWZ5DUYH
```


