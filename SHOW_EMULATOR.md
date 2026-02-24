# How to Show Emulator in Android Studio

## 🔍 The Emulator is Running!

Your emulator (`emulator-5554`) is running, but the window might be hidden. Here's how to show it:

## Method 1: Show Emulator Window in Android Studio

1. **Open Android Studio**
2. **Look for the emulator icon** in the toolbar (usually at the top)
3. **Click on the emulator dropdown** - you should see `emulator-5554` listed
4. **Click on the emulator name** - this will bring the emulator window to the front

## Method 2: Use Android Studio Device Manager

1. **Open Android Studio**
2. **Go to**: `Tools` → `Device Manager` (or click the device manager icon)
3. **Find your emulator** in the list (should show as "Medium_Phone_API_36.0" or similar)
4. **Click the play/start icon** next to it, or **right-click** → `Show Emulator Window`

## Method 3: Show All Windows

1. **Press `Cmd + Tab`** (Mac) to see all open windows
2. **Look for "Android Emulator"** or "qemu" in the list
3. **Click on it** to bring it to the front

## Method 4: Restart Emulator Window

If the window is completely gone:

1. **In Android Studio Device Manager**:
   - Find your emulator
   - Click the **stop icon** (square) to stop it
   - Click the **play icon** (triangle) to start it again

2. **Or use command line**:
   ```bash
   # Stop emulator
   adb -s emulator-5554 emu kill
   
   # Start emulator (from Android Studio or AVD Manager)
   ```

## Method 5: Launch App Directly

The app is already installed! You can launch it directly:

```bash
# Launch the app
adb -s emulator-5554 shell am start -n com.example.go_rent_frontend/.MainActivity

# Or use Flutter
cd go_rent_frontend
flutter run -d emulator-5554
```

## ✅ Quick Check Commands

```bash
# Check if emulator is running
adb devices

# Check if app is installed
adb -s emulator-5554 shell pm list packages | grep rent

# Launch app
adb -s emulator-5554 shell am start -n com.example.go_rent_frontend/.MainActivity

# View emulator logs
adb -s emulator-5554 logcat | grep -i flutter
```

## 🎯 Current Status

- ✅ **Emulator**: Running (`emulator-5554`)
- ✅ **App**: Installed (`com.example.go_rent_frontend`)
- ✅ **Backend**: Running on `localhost:8081`
- ✅ **Database**: Running on `localhost:3307`

**The app should be visible on your emulator!** If you still don't see it:

1. Try launching it with the command above
2. Check Android Studio's emulator window
3. Restart the emulator from Android Studio Device Manager

## 🔧 Troubleshooting

### If emulator window is completely missing:
1. Open Android Studio
2. Go to `Tools` → `Device Manager`
3. Stop and restart the emulator

### If app doesn't appear:
```bash
# Uninstall and reinstall
adb -s emulator-5554 uninstall com.example.go_rent_frontend
cd go_rent_frontend
flutter run -d emulator-5554
```

### If emulator is frozen:
```bash
# Restart emulator
adb -s emulator-5554 emu kill
# Then start from Android Studio Device Manager
```


