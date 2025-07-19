package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go-rent/config"
	"golang.org/x/oauth2/google"
	"go-rent/utils"
)

// FCM message structure for HTTP v1 API
type FCMV1Message struct {
	Message FCMV1MessageData `json:"message"`
}

type FCMV1MessageData struct {
	Token        string                 `json:"token,omitempty"`
	Topic        string                 `json:"topic,omitempty"`
	Notification FCMV1Notification      `json:"notification,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Android      FCMV1AndroidConfig     `json:"android,omitempty"`
	APNS         FCMV1APNSConfig        `json:"apns,omitempty"`
}

type FCMV1Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type FCMV1AndroidConfig struct {
	Priority string `json:"priority"`
}

type FCMV1APNSConfig struct {
	Headers FCMV1APNSHeaders `json:"headers"`
}

type FCMV1APNSHeaders struct {
	APNSPriority string `json:"apns-priority"`
}

// Legacy FCM message structure (for reference)
type FCMMessage struct {
	To           string                 `json:"to,omitempty"`
	Topic        string                 `json:"topic,omitempty"`
	Notification FCMNotification        `json:"notification,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Priority     string                 `json:"priority,omitempty"`
}

type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound,omitempty"`
}

// FCM response structure
type FCMResponse struct {
	Success int `json:"success"`
	Failure int `json:"failure"`
	Results []struct {
		MessageID string `json:"message_id,omitempty"`
		Error     string `json:"error,omitempty"`
	} `json:"results"`
}

// Service account credentials (you can move this to a config file)
const serviceAccountKey = `{
  "type": "service_account",
  "project_id": "rental-app-9e26c",
  "private_key_id": "f3802557709dc7d020e565a8acf1b503786c3aef",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCqeljcf8OtALeZ\nEj8KBuGMzfy8vPL9q/6/0UOGL3ewiAL+J7M7tXjaiCmNrl+n7Ek5TzORIAh7Fo3Q\nlJ+QUjVIVCuVSLQUBKis91VtrdcZZN+XD36K60Z04kdZJb+ttx8l+otfeK1BtSMd\nGvHiAX2LGd1gvdg6m4huGr/39ZI5wbQOuKcM1eI1v8S8u9yOXBfzq6zH1xH/ibPS\n5yQOp8qqlQreXNqZkQQuB+BW73j03Boz7tJZzv3mdQFU9JwJev6Javfj5a/wC0DE\nlfdP1b2fF6bz1ATZi6yB1zjMI1Xafie2pXm9GQwHw9qtckELJv00kHSs0KG5Ln4L\nV3TzBGhPAgMBAAECggEAFBuOM6PWPqehACsIyHP7UpJRRym6N3V7/MAACOm8YlQI\nllW0bEyBTrVUwWNZ4QKbuGjBGlL+7RXB8GI3V9x6cGeNJHSI2PubsZsStncUsegV\n/0lSkawiiVLPO5KaQzEgOWiN29ruBJwK4jn1YcTvO+L3G5wIzxDdTK9XCTYr4yfm\nSWce+393nGc1WajvKAriuJzIHoCphLcq3FGgoni5jdeHiDzeaxtIe+oTyJr/rV88\npwv1mO4K2e3UHNGQSM1KQtpSUILPaHiEGYAuNFG843lb7QHKYiZU5FmQwBf3ls5D\nOeTywkh9xHsNPgB3NEgjRJGf/mXaQrwv+ZUTD1XlQQKBgQDlD0+Wo85xH6M9fS0P\nvg2APMpeqVpF5knZUs0/j/X4HAwY3TLI9h4RmLG/z297/cAQgxbW5iWDswDhwe6D\n85wF7lfwyJf+mkIYqM6S8KXE+fUKyau5GP7a9TAskBoZaCXp4toP3b+uZjzRqkKt\n36qhSlQlycQKFbdamwKA+cgeMQKBgQC+hzZyjPri8wu2cgURzHbFsso7IhrEondw\naIFSHH0ZJR5TtaSEOgMW6C2qzmmGZDw6O2NR/fk7aSXl3J9CmLM4c7YPjqX8vB8+\nbebMJW6XCkS9u+KBAeaZbOkdRLBsxd9cW054zQbVrwMra4QvEDfz5uTsyp0HsH+5\nvBIIsRLOfwKBgDojSrYXWXyI5SvkK53FRTP5OfwQy+LV8oSAaavqZCnXJJLjAdLT\n9QnOUb83bTpxS2BlhVSCEZ99vYmPaXSATmeK+TMaFsn+aSxNHDFbdxepwbI9QaK3\nX2g/tzx4TseIEadtdp90TwR62pD0v/vVuz842GbG8UUGAgWzNk16GHrhAoGACvjC\nW+pecD9Kx2Ddhd7eYBghqTIXlIc+lYyPFelqEs6eZnepV6v3jZQlPRbR4NlY1omg\n+JHFjnRJqGkCCtW8TF3teAvg5yL2MaQmjE8DhVMkDkEJlCBF5UPuUK8p8bmbWTgw\n1qgH4rpHVnLEk+k9L6B2QmSQkmbJlCqOZ027JYUCgYEA21ALKodyQShlz7C1DXOb\nvqm5JYm4XVBS4418VTa6NawxL2HEyYsx2/MdGKREekJhN+jXs25fo/uHXxA16Rhu\nD4i27cVJ+QU0loRwwVrOy0OFP4Ywc43uG56WVqG+36Ze50ZJ4X7T91A6GC5v+QIv\nPYf9iCANZEGLsgti0dOpiUk=\n-----END PRIVATE KEY-----\n",
  "client_email": "fcm-server@rental-app-9e26c.iam.gserviceaccount.com",
  "client_id": "113657692866788191739",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/fcm-server%40rental-app-9e26c.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`

// Update FCM token handler
func UpdateFCMTokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== Update FCM Token Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Parse request body
	var request struct {
		FCMToken string `json:"fcm_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get database connection
	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	// Update FCM token in database
	_, err = db.Exec(`
		UPDATE user 
		SET fcm_token = ?, updated_at = NOW() 
		WHERE id = ?
	`, request.FCMToken, userID)

	if err != nil {
		fmt.Printf("Error updating FCM token: %v\n", err)
		http.Error(w, "Failed to update FCM token", http.StatusInternalServerError)
		return
	}

	fmt.Printf("FCM token updated for user %d: %s", userID, request.FCMToken)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "FCM token updated successfully",
	})
}

// Test push notification handler
func TestPushNotificationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== Test Push Notification Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Parse request body
	var request struct {
		Title string                 `json:"title"`
		Body  string                 `json:"body"`
		Data  map[string]interface{} `json:"data,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use default values if not provided
	if request.Title == "" {
		request.Title = "Test Notification"
	}
	if request.Body == "" {
		request.Body = "This is a test push notification!"
	}
	if request.Data == nil {
		request.Data = map[string]interface{}{
			"type":      "test",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		}
	}

	// Ensure the data contains navigation information
	if request.Data["type"] == nil {
		request.Data["type"] = "test"
	}

	// Send push notification
	err := SendPushNotification(userID, request.Title, request.Body, request.Data)
	if err != nil {
		fmt.Printf("Error sending push notification: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Failed to send push notification",
			"error":   err.Error(),
		})
		return
	}

	fmt.Printf("Test push notification sent to user %d", userID)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Test push notification sent successfully",
	})
}

// Send push notification to user
func SendPushNotification(userID int64, title, body string, data map[string]interface{}) error {
	db, err := config.GetDBConnection()
	if err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}

	// Get user's FCM token
	var fcmToken string
	err = db.QueryRow("SELECT fcm_token FROM user WHERE id = ?", userID).Scan(&fcmToken)
	if err != nil {
		return fmt.Errorf("failed to get FCM token: %v", err)
	}

	if fcmToken == "" {
		return fmt.Errorf("no FCM token found for user %d", userID)
	}

	// Send notification
	return sendFCMNotification(fcmToken, title, body, data)
}

// SendNotificationWithPush creates a notification in the database AND sends a push notification
func SendNotificationWithPush(senderID, receiverID, propertyID, floorID int64, message, status string, comment *string) error {
	db, err := config.GetDBConnection()
	if err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}

	// Generate notification ID
	notificationID, err := utils.GenerateRandomID()
	if err != nil {
		return fmt.Errorf("failed to generate notification ID: %v", err)
	}

	// Insert notification into database
	_, err = db.Exec(`
		INSERT INTO notification (
			id, message, sender, receiver, pid, fid,
			status, comment, created_at, created_by, updated_at, updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		notificationID,
		message,
		senderID,
		receiverID,
		propertyID,
		floorID,
		status,
		comment,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		senderID,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		senderID,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification in database: %v", err)
	}

	// Get property and floor names for better push notification titles
	var propertyName, floorName string
	err = db.QueryRow("SELECT p.name, f.name FROM property p JOIN floor f ON p.id = f.pid WHERE p.id = ? AND f.id = ?", propertyID, floorID).Scan(&propertyName, &floorName)
	if err != nil {
		// If we can't get property/floor names, use generic names
		propertyName = "Property"
		floorName = "Floor"
	}

	// Send push notification

	title := ""
	body := message
	
	// Customize title based on notification type with property and floor info
	if strings.Contains(message, "Tenant request") {
		title = fmt.Sprintf("New Tenant Request - %s %s", propertyName, floorName)
	} else if strings.Contains(message, "Payment amount") {
		title = fmt.Sprintf("Payment Notification - %s %s", propertyName, floorName)
	} else if strings.Contains(message, "Advance payment") {
		title = fmt.Sprintf("Advance Payment Request - %s %s", propertyName, floorName)
	} else if strings.Contains(message, "accepted") || strings.Contains(message, "rejected") {
		title = fmt.Sprintf("Request Update - %s %s", propertyName, floorName)
	} else if strings.Contains(message, "Monthly rent reminder") {
		title = fmt.Sprintf("Monthly Rent Reminder - %s %s", propertyName, floorName)
	} else {
		title = fmt.Sprintf("New Notification! - %s %s", propertyName, floorName)
	}

	// Determine notification type based on message content
	notificationType := "notification"
	if strings.Contains(message, "Monthly rent reminder") {
		notificationType = "monthly_reminder"
	} else if strings.Contains(message, "Payment amount") {
		notificationType = "payment"
	} else if strings.Contains(message, "Advance payment") {
		notificationType = "advance_payment"
	} else if strings.Contains(message, "Tenant request") {
		notificationType = "tenant_request"
	}

	data := map[string]interface{}{
		"notification_id": fmt.Sprintf("%d", notificationID),
		"type":            notificationType,
		"property_id":     fmt.Sprintf("%d", propertyID),
		"floor_id":        fmt.Sprintf("%d", floorID),
		"timestamp":       fmt.Sprintf("%d", time.Now().Unix()),
	}

	// Send push notification to receiver
	err = SendPushNotification(receiverID, title, body, data)
	if err != nil {
		// Log the error but don't fail the entire operation
		fmt.Printf("Warning: Failed to send push notification to user %d: %v\n", receiverID, err)
		// Don't return error here as the notification was successfully created in database
	} else {
		fmt.Printf("Push notification sent successfully to user %d for notification %d\n", receiverID, notificationID)
	}

	return nil
}

// Send FCM notification using HTTP v1 API
func sendFCMNotification(token, title, body string, data map[string]interface{}) error {
	// Convert all data values to strings (FCM requirement)
	stringData := make(map[string]interface{})
	for key, value := range data {
		switch v := value.(type) {
		case string:
			stringData[key] = v
		case int, int64, float64:
			stringData[key] = fmt.Sprintf("%v", v)
		default:
			stringData[key] = fmt.Sprintf("%v", v)
		}
	}

	// Create the FCM message payload for HTTP v1 API
	message := map[string]interface{}{
		"message": map[string]interface{}{
			"token": token,
			"notification": map[string]string{
				"title": title,
				"body":  body,
			},
			"data": stringData,
			"android": map[string]interface{}{
				"priority": "high",
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Create HTTP request
	url := "https://fcm.googleapis.com/v1/projects/rental-app-9e26c/messages:send"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Get access token
	accessToken, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	fmt.Printf("FCM notification sent successfully to token: %s\n", token)
	return nil
}

// Get access token for FCM API
func getAccessToken() (string, error) {
	// Read service account key file
	keyFile := "config/firebase-service-account.json"
	creds, err := os.ReadFile(keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read service account key file: %v", err)
	}

	// Create JWT config
	jwtConfig, err := google.JWTConfigFromJSON(creds, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("failed to create JWT config: %v", err)
	}

	// Get token
	token, err := jwtConfig.TokenSource(context.Background()).Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %v", err)
	}

	return token.AccessToken, nil
}

// Send notification to topic
func SendTopicNotification(topic, title, body string, data map[string]interface{}) error {
	// TODO: Implement topic notification
	fmt.Printf("=== FCM Topic Notification Simulation ===\n")
	fmt.Printf("Topic: %s\n", topic)
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Body: %s\n", body)
	fmt.Printf("Data: %+v\n", data)
	fmt.Printf("=== End Topic Simulation ===\n")
	
	return nil
} 

// Test FCM connection
func TestFCMConnection() error {
	// Try to get an access token
	accessToken, err := getAccessToken()
	if err != nil {
		return fmt.Errorf("FCM authentication failed: %v", err)
	}

	if accessToken == "" {
		return fmt.Errorf("FCM access token is empty")
	}

	fmt.Printf("FCM authentication successful, token length: %d\n", len(accessToken))
	return nil
} 

// Test FCM connection handler
func TestFCMConnectionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== Test FCM Connection Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	// Test FCM connection
	err := TestFCMConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "FCM connection failed",
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "FCM connection successful",
	})
}

// Test FCM connection handler (public, no auth required)
func TestFCMConnectionPublicHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== Test FCM Connection Public Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	// Test FCM connection
	err := TestFCMConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "FCM connection failed",
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "FCM connection successful",
	})
} 