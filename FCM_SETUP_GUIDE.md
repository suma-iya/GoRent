# Firebase Cloud Messaging (FCM) Setup Guide

## Overview
This guide will help you set up Firebase Cloud Messaging (FCM) for push notifications in your Go Rent application.

## Prerequisites
1. Google Cloud Console access
2. Firebase project: `rental-app-9e26c`
3. Service account with Firebase Admin permissions

## Step 1: Enable Firebase Cloud Messaging API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Select your project: `rental-app-9e26c`
3. Go to **APIs & Services > Library**
4. Search for "Firebase Cloud Messaging API"
5. Click on it and press **Enable**

## Step 2: Create Service Account

1. Go to **IAM & Admin > Service Accounts**
2. Click **Create Service Account**
3. Name: `fcm-server`
4. Description: `Service account for FCM push notifications`
5. Click **Create and Continue**

## Step 3: Grant Permissions

1. Add the following roles:
   - **Firebase Admin**
   - **Cloud Messaging Admin**
2. Click **Continue**
3. Click **Done**

## Step 4: Create Service Account Key

1. Find your service account in the list
2. Click on the email address
3. Go to **Keys** tab
4. Click **Add Key > Create new key**
5. Choose **JSON** format
6. Click **Create**
7. Download the JSON file

## Step 5: Update Configuration

1. Replace the content of `config/firebase-service-account.json` with your downloaded JSON
2. Make sure the file contains the actual private key (not placeholder)

## Step 6: Test FCM Connection

1. Start your Go server
2. Test the connection:
   ```bash
   curl -X GET http://localhost:8080/test/fcm-connection \
     -H "Cookie: sessiontoken=YOUR_SESSION_TOKEN"
   ```

## Step 7: Test Push Notification

1. Make sure your Flutter app is running and you're logged in
2. Press the "🧪 TEST PUSH NOTIFICATION" button in the app
3. Check the server logs for FCM response

## Troubleshooting

### Common Issues:

1. **"FCM authentication failed"**
   - Check if the service account JSON file is correct
   - Verify the Firebase Admin role is assigned
   - Ensure the FCM API is enabled

2. **"FCM API returned status 403"**
   - Check if the service account has proper permissions
   - Verify the project ID matches

3. **"No FCM token found for user"**
   - Make sure the user has logged in and FCM token is saved
   - Check if the FCM token is being sent to the server

### Debug Steps:

1. Check server logs for detailed error messages
2. Verify the FCM token is being generated in the Flutter app
3. Test the FCM connection endpoint
4. Check if the user's FCM token is saved in the database

## API Endpoints

- `GET /test/fcm-connection` - Test FCM authentication
- `POST /test/push-notification` - Send test push notification
- `POST /user/fcm-token` - Update user's FCM token

## Flutter App Integration

The Flutter app is already configured with:
- Firebase Core
- Firebase Messaging
- Local Notifications
- FCM token generation and sending to server

## Security Notes

1. Keep the service account JSON file secure
2. Don't commit the actual private key to version control
3. Use environment variables for production
4. Regularly rotate service account keys 