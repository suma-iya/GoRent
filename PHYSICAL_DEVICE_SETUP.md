# Running App on Physical Device

## ✅ Configuration

The app is now configured to run on your physical Android device (RMX3085).

### Backend URL Configuration

- **Android Emulator**: `http://10.0.2.2:8081`
- **Physical Android Device**: `http://192.168.0.232:8081` (your Mac's IP)
- **iOS Simulator**: `http://localhost:8081`

### Your Mac's IP Address
- **Current IP**: `192.168.0.232`
- **Backend Port**: `8081`

## 📱 Running on Physical Device

```bash
cd go_rent_frontend
flutter run -d BQ4X6595TWZ5DUYH
```

Or use the device name:
```bash
flutter run -d RMX3085
```

## 🔧 Important Notes

### 1. Ensure Same Network
Your phone and Mac must be on the **same Wi-Fi network** for the app to connect to the backend.

### 2. Firewall Settings
Make sure your Mac's firewall allows connections on port 8081:
- System Settings → Network → Firewall
- Allow incoming connections for Docker or port 8081

### 3. If IP Address Changes
If your Mac's IP address changes, update it in:
- `go_rent_frontend/lib/services/api_service.dart`
- Line ~25: Change `192.168.0.232` to your new IP

To find your Mac's IP:
```bash
ipconfig getifaddr en0
# or
ipconfig getifaddr en1
```

### 4. Test Backend Connection
From your Mac, test that backend is accessible:
```bash
curl http://localhost:8081/test/fcm-connection-public
```

From your phone's browser, test:
```
http://192.168.0.232:8081/test/fcm-connection-public
```

## 🚀 Quick Start

1. **Ensure Docker is running**:
   ```bash
   docker-compose ps
   ```

2. **Connect your phone via USB** (or ensure it's on same Wi-Fi)

3. **Check device is connected**:
   ```bash
   flutter devices
   ```

4. **Run the app**:
   ```bash
   cd go_rent_frontend
   flutter run -d BQ4X6595TWZ5DUYH
   ```

## ✅ Verification

After running, you should see:
- App installing on your phone
- Login screen appearing
- App can connect to backend at `192.168.0.232:8081`

## 🔍 Troubleshooting

### App can't connect to backend:
1. Check phone and Mac are on same Wi-Fi
2. Verify Mac's IP hasn't changed: `ipconfig getifaddr en0`
3. Test from phone browser: `http://192.168.0.232:8081/test/fcm-connection-public`
4. Check Mac firewall settings

### Device not detected:
```bash
# Enable USB debugging on phone
# Settings → Developer Options → USB Debugging

# Check connection
adb devices
```

### Backend not accessible:
```bash
# Check Docker is running
docker-compose ps

# Check backend logs
docker-compose logs backend
```


