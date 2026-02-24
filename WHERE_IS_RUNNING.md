# Where Is The Project Running? 🚀

## 📍 Current Running Services

### ✅ Backend (Go API Server)
- **Status**: ✅ Running in Docker
- **Container Name**: `rent-backend`
- **Host URL**: `http://localhost:8081`
- **Container Port**: 8080 (internal)
- **Access from**: 
  - Your Mac: `http://localhost:8081`
  - Android Emulator: `http://10.0.2.2:8081`
  - iOS Simulator: `http://localhost:8081`

**Test it:**
```bash
curl http://localhost:8081/test/fcm-connection-public
```

### ✅ Database (MySQL)
- **Status**: ✅ Running in Docker
- **Container Name**: `rent-mysql`
- **Host Port**: 3307
- **Container Port**: 3306 (internal)
- **Access from**: 
  - Your Mac: `localhost:3307`
  - Docker containers: `mysql:3306` (service name)

### ✅ Flutter App (Frontend)
- **Status**: ✅ Running on Android Emulator
- **Device**: `emulator-5554` (sdk gphone64 arm64)
- **App Name**: Your rent app
- **Backend URL**: `http://10.0.2.2:8081` (configured in `api_service.dart`)

## 🔍 How to See What's Running

### 1. Check Docker Containers
```bash
docker-compose ps
```

**Output shows:**
- Container names
- Status (Up/Down)
- Port mappings
- Health status

### 2. Check Docker Logs
```bash
# Backend logs
docker-compose logs -f backend

# MySQL logs
docker-compose logs -f mysql

# All services
docker-compose logs -f
```

### 3. Check Running Ports
```bash
# See what's using ports 8081 and 3307
lsof -i :8081
lsof -i :3307

# Or use netstat
netstat -an | grep -E "8081|3307"
```

### 4. Check Flutter App Status
```bash
# List connected devices
flutter devices

# Check if Flutter is running
ps aux | grep flutter
```

### 5. Access Backend API Directly
Open in browser or use curl:
```bash
# Test endpoint
curl http://localhost:8081/test/fcm-connection-public

# Or open in browser
open http://localhost:8081/test/fcm-connection-public
```

## 📱 How to See the Flutter App

### On Android Emulator:
1. **Look at your Android Emulator window** - the app should be visible there
2. If not visible, check:
   ```bash
   flutter devices
   # Should show: emulator-5554
   ```

### If App is Not Visible:
```bash
cd go_rent_frontend
flutter run -d emulator-5554
```

## 🌐 Access Points Summary

| Service | Local Access | Emulator Access | Container Access |
|---------|-------------|-----------------|------------------|
| **Backend API** | `localhost:8081` | `10.0.2.2:8081` | `backend:8080` |
| **MySQL** | `localhost:3307` | N/A | `mysql:3306` |
| **Flutter App** | Emulator Screen | Emulator Screen | N/A |

## 🔧 Quick Commands

### View All Running Services
```bash
# Docker containers
docker-compose ps

# Flutter processes
ps aux | grep flutter

# Network connections
lsof -i :8081 -i :3307
```

### View Logs
```bash
# Backend logs (real-time)
docker-compose logs -f backend

# Flutter logs (if running in terminal)
# Check the terminal where you ran `flutter run`
```

### Restart Services
```bash
# Restart Docker services
docker-compose restart

# Restart Flutter app
# Press 'r' in Flutter terminal, or stop and restart:
flutter run -d emulator-5554
```

## ✅ Verification Checklist

- [ ] Docker containers running: `docker-compose ps`
- [ ] Backend responding: `curl http://localhost:8081/test/fcm-connection-public`
- [ ] Flutter app visible on emulator
- [ ] Can login to app
- [ ] App connects to backend (check Flutter console logs)

## 🎯 Current Status

Based on current checks:
- ✅ **Backend**: Running on `localhost:8081`
- ✅ **Database**: Running on `localhost:3307`
- ✅ **Flutter**: Running on `emulator-5554`
- ✅ **All services**: Healthy and connected

**Your app should be visible on the Android emulator!** 🎉


