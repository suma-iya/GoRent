package handlers

import (
	
	"database/sql"
	"encoding/json"
	"fmt"
	"go-rent/config"
	"go-rent/utils"
	
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"log"
	"regexp"
)

// getUserIDFromContext retrieves the user ID from the request context
func getUserIDFromContext(r *http.Request) int64 {
	if userID, ok := r.Context().Value("userID").(int64); ok {
		return userID
	}
	return 0
}

type PropertyRequest struct {
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Photo   *string `json:"photo,omitempty"`
}

type PropertyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	PropertyID int64 `json:"property_id,omitempty"`
}

type Property struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Photo     *string `json:"photo,omitempty"`
	CreatedAt string  `json:"created_at"`
}

type UserPropertiesResponse struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Properties []Property `json:"properties"`
}

type SinglePropertyResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Property Property `json:"property,omitempty"`
	Floors   []Floor  `json:"floors,omitempty"`
	IsManager bool    `json:"is_manager,omitempty"`
}

type Floor struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Rent      int    `json:"rent"`
	CreatedAt string `json:"created_at"`
	Tenant    *int64 `json:"tenant,omitempty"`
	TenantName *string `json:"tenant_name,omitempty"`
	Status    string `json:"status,omitempty"`
	NotificationID *int64 `json:"notification_id,omitempty"`
	HasPendingAdvancePayment bool `json:"has_pending_advance_payment"`
}

type FloorRequest struct {
	Name              string `json:"name"`
	Rent             int    `json:"rent"`
	Tenant           *int64 `json:"tenant,omitempty"`
	ReceivedMoney    int    `json:"received_money,omitempty"`
}

type FloorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	FloorID int64  `json:"floor_id,omitempty"`
}

type UserPhone struct {
	ID    int64  `json:"id"`
	Phone string `json:"phone_number"`
}

type UserPhonesResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Users   []UserPhone `json:"users"`
}

type UserIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int64  `json:"user_id,omitempty"`
}

type PaymentRequest struct {
	Rent          int     `json:"rent"`
	ReceivedMoney int     `json:"received_money"`
	FullPayment   bool    `json:"full_payment"`
	ElectricityBill *int  `json:"electricity_bill,omitempty"`
}

type PaymentResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	PaymentID int64  `json:"payment_id,omitempty"`
}

type TenantRequestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Notification struct {
	ID        int64  `json:"id"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Property  struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"property"`
	Floor struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"floor"`
	ShowActions bool `json:"show_actions"`
	IsRead      bool `json:"is_read"`
	Comment     *string `json:"comment,omitempty"`
	SenderID    *int64  `json:"sender_id,omitempty"`
	ReceiverID  *int64  `json:"receiver_id,omitempty"`
	SenderName  *string `json:"sender_name,omitempty"`
	ReceiverName *string `json:"receiver_name,omitempty"`
}

type NotificationsResponse struct {
	Success       bool          `json:"success"`
	Message       string        `json:"message"`
	Notifications []Notification `json:"notifications"`
}

type PaymentNotificationRequest struct {
	Amount int     `json:"amount"`
	Month  *int    `json:"month"`
	PaidElectricityBill *int `json:"paid_electricity_bill,omitempty"`
}

type AdvancePaymentRequest struct {
	AdvanceUID int64 `json:"advance_uid"`
	Money      int   `json:"money"`
}

type AdvancePaymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	AdvanceID int64 `json:"advance_id,omitempty"`
}

type AdvancePaymentCheckResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	HasPending  bool   `json:"has_pending"`
}

type AdvanceDetail struct {
	ID        int64  `json:"id"`
	AdvanceUID int64 `json:"advance_uid"`
	UserName  string `json:"user_name"`
	Money     int    `json:"money"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type AdvanceDetailsResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Advances []AdvanceDetail `json:"advances"`
}

type PaginationInfo struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
	Limit       int `json:"limit"`
	HasNextPage bool `json:"has_next_page"`
	HasPrevPage bool `json:"has_prev_page"`
}

type PaymentHistoryResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Payments   []PaymentHistory `json:"payments"`
	Pagination PaginationInfo  `json:"pagination"`
}

type PaymentHistory struct {
	ID           int64   `json:"id"`
	NewAddedRent float64 `json:"new_added_rent"`
	Rent         float64 `json:"rent"`
	ReceivedMoney float64 `json:"received_money"`
	DueRent      float64 `json:"due_rent"`
	FullPayment  bool    `json:"full_payment"`
	CreatedAt    string  `json:"created_at"`
	NewAddedElectricityBill *float64 `json:"new_added_electricity_bill,omitempty"`
	PaidElectricityBill     *float64 `json:"paid_electricity_bill,omitempty"`
	DueElectricityBill      *float64 `json:"due_electricity_bill,omitempty"`
	ElectricityBill         *float64 `json:"electricity_bill,omitempty"`
}

func AddPropertyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Add Property Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Method not allowed", 0})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(PropertyResponse{false, "User not authenticated", 0})
		return
	}

	var req PropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Invalid request body", 0})
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Property name is required", 0})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Database connection error", 0})
		return
	}

	// Generate random ID for property
	randomID, err := utils.GenerateRandomID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Error generating property ID", 0})
		return
	}

	// Insert property into database
	_, err = db.Exec(
		`INSERT INTO property (id, name, address, photo, created_at, created_by, updated_at, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		randomID,
		req.Name,
		req.Address,
		req.Photo,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Error adding property", 0})
		return
	}

	// Insert into takes_care_of table
	takesCareID, err := utils.GenerateRandomID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Error generating care ID", 0})
		return
	}

	_, err = db.Exec(
		`INSERT INTO takes_care_of (id, uid, pid, created_at, created_by, updated_at, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		takesCareID,
		userID,
		randomID,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PropertyResponse{false, "Error saving property care details", 0})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PropertyResponse{
		Success: true,
		Message: "Property added successfully",
		PropertyID: randomID,
	})
}

func GetUserPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get User Properties Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Method not allowed", nil})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "User not authenticated", nil})
		return
	}

	fmt.Printf("Fetching properties for user ID: %d\n", userID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Database connection error", nil})
		return
	}

	// Query to get all properties for the user
	query := `
		SELECT p.id, p.name, p.address, p.photo, p.created_at 
		FROM property p
		INNER JOIN takes_care_of t ON p.id = t.pid
		WHERE t.uid = ?
		ORDER BY p.created_at DESC`
	
	fmt.Printf("Executing query: %s with userID: %d\n", query, userID)
	
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Printf("Error querying properties: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Error fetching properties", nil})
		return
	}
	defer rows.Close()

	var properties []Property
	for rows.Next() {
		var prop Property
		var photo sql.NullString
		if err := rows.Scan(&prop.ID, &prop.Name, &prop.Address, &photo, &prop.CreatedAt); err != nil {
			fmt.Printf("Error scanning property row: %v\n", err)
			continue
		}
		if photo.Valid {
			prop.Photo = &photo.String
		}
		properties = append(properties, prop)
		fmt.Printf("Found property: ID=%d, Name=%s\n", prop.ID, prop.Name)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating property rows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Error processing properties", nil})
		return
	}

	fmt.Printf("Found %d properties\n", len(properties))

	response := UserPropertiesResponse{
		Success: true,
		Message: "Properties retrieved successfully",
		Properties: properties,
	}

	// Log the response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Sending response: %s\n", string(responseJSON))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPropertyByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Property By ID Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "Method not allowed", Property{}, nil, false})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "User not authenticated", Property{}, nil, false})
		return
	}

	// Extract property ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "Invalid property ID format", Property{}, nil, false})
		return
	}

	propertyID, err := strconv.ParseInt(pathParts[2], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "Invalid property ID", Property{}, nil, false})
		return
	}

	fmt.Printf("Fetching property ID: %d for user ID: %d\n", propertyID, userID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "Database connection error", Property{}, nil, false})
		return
	}

	// Query to get the specific property and verify user has access
	query := `
		SELECT p.id, p.name, p.address, p.photo, p.created_at 
		FROM property p
		WHERE p.id = ? AND (
			EXISTS (
				SELECT 1 FROM takes_care_of t 
				WHERE t.pid = p.id AND t.uid = ?
			) OR EXISTS (
				SELECT 1 FROM floor f 
				WHERE f.pid = p.id AND f.tenant = ?
			)
		)`
	
	fmt.Printf("Executing query: %s with propertyID: %d, userID: %d\n", query, propertyID, userID, userID)
	
	var prop Property
	var photo sql.NullString
	err = db.QueryRow(query, propertyID, userID, userID).Scan(&prop.ID, &prop.Name, &prop.Address, &photo, &prop.CreatedAt)
	if err != nil {
		fmt.Printf("Error querying property: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "Property not found or access denied", Property{}, nil, false})
		return
	}
	if photo.Valid {
		prop.Photo = &photo.String
	}

	// Get all floors for this property with tenant names
	floorsQuery := `
		SELECT f.id, f.name, f.rent, f.created_at, f.tenant,
		       CASE 
		           WHEN u.name IS NULL OR u.name = '' THEN CONCAT('User ', u.id)
		           ELSE u.name 
		       END as tenant_name,
		       EXISTS(
		           SELECT 1 FROM notification n 
		           WHERE n.fid = f.id AND n.status = 'pending'
		           AND n.message NOT LIKE 'Advance payment request:%'
		       ) as has_pending_request,
		       (
		           SELECT n.id 
		           FROM notification n 
		           WHERE n.fid = f.id AND n.status = 'pending'
		           AND n.message NOT LIKE 'Advance payment request:%'
		           LIMIT 1
		       ) as notification_id,
		       EXISTS(
		           SELECT 1 FROM advance a 
		           WHERE a.fid = f.id AND a.status = 'pending'
		       ) as has_pending_advance_payment
		FROM floor f
		LEFT JOIN user u ON f.tenant = u.id
		WHERE f.pid = ?
		ORDER BY f.created_at DESC`
	
	rows, err := db.Query(floorsQuery, propertyID)
	if err != nil {
		fmt.Printf("Error querying floors: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SinglePropertyResponse{false, "Error fetching floors", prop, nil, false})
		return
	}
	defer rows.Close()

	var floors []Floor
	for rows.Next() {
		var floor Floor
		var tenant sql.NullInt64
		var tenantName sql.NullString
		var hasPendingRequest bool
		var notificationID sql.NullInt64
		var hasPendingAdvancePayment bool
		if err := rows.Scan(&floor.ID, &floor.Name, &floor.Rent, &floor.CreatedAt, &tenant, &tenantName, &hasPendingRequest, &notificationID, &hasPendingAdvancePayment); err != nil {
			fmt.Printf("Error scanning floor row: %v\n", err)
			continue
		}
		if tenant.Valid {
			floor.Tenant = &tenant.Int64
		}
		if tenantName.Valid {
			floor.TenantName = &tenantName.String
		}
		if hasPendingRequest {
			floor.Status = "pending"
			if notificationID.Valid {
				floor.NotificationID = &notificationID.Int64
			}
		}
		floor.HasPendingAdvancePayment = hasPendingAdvancePayment
		floors = append(floors, floor)
	}

	fmt.Printf("Found property: ID=%d, Name=%s with %d floors\n", prop.ID, prop.Name, len(floors))

	// Check if user is a manager of this property
	var isManager bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE uid = ? AND pid = ?
		)`, userID, propertyID).Scan(&isManager)
	
	if err != nil {
		fmt.Printf("Error checking manager status: %v\n", err)
		isManager = false
	}

	fmt.Printf("User %d is manager of property %d: %v\n", userID, propertyID, isManager)

	response := SinglePropertyResponse{
		Success: true,
		Message: "Property retrieved successfully",
		Property: prop,
		Floors: floors,
		IsManager: isManager,
	}

	// Log the response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Sending response: %s\n", string(responseJSON))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandlePropertyRoutes handles all property-related routes
func HandlePropertyRoutes(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	
	// Remove empty strings from path parts
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	// Handle different routes
	switch {
	case len(cleanParts) == 1 && cleanParts[0] == "property":
		// /property
		GetUserPropertiesHandler(w, r)
	case len(cleanParts) == 2 && cleanParts[0] == "property":
		if cleanParts[1] == "tenant" {
			// /property/tenant
			GetUserTenantPropertiesHandler(w, r)
		} else {
			// /property/{id}
			GetPropertyByIDHandler(w, r)
		}
	case len(cleanParts) == 3 && cleanParts[0] == "property" && cleanParts[2] == "floor":
		// /property/{id}/floor
		switch r.Method {
		case http.MethodGet:
			GetFloorsHandler(w, r)
		case http.MethodPost:
			AddFloorHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(FloorResponse{false, "Method not allowed", 0})
		}
	case len(cleanParts) == 4 && cleanParts[0] == "property" && cleanParts[2] == "floor":
		// /property/{id}/floor/{floor_id}
		switch r.Method {
		case http.MethodGet:
			GetFloorByIDHandler(w, r)
		case http.MethodPut:
			UpdateFloorHandler(w, r)
		case http.MethodPost:
			SendTenantRequestHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(FloorResponse{false, "Method not allowed", 0})
		}
	case len(cleanParts) == 5 && cleanParts[0] == "property" && cleanParts[2] == "floor" && cleanParts[4] == "payment":
		// /property/{id}/floor/{floor_id}/payment
		switch r.Method {
		case http.MethodPost:
			CreatePaymentHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(PaymentResponse{false, "Method not allowed", 0})
		}
	default:
		http.NotFound(w, r)
	}
}

func AddFloorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Add Floor Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(FloorResponse{false, "Method not allowed", 0})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(FloorResponse{false, "User not authenticated", 0})
		return
	}

	// Extract property ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid URL format", 0})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid property ID", 0})
		return
	}

	var req FloorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid request body", 0})
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Floor name is required", 0})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Database connection error", 0})
		return
	}

	// Verify user has access to the property
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE pid = ? AND uid = ?
		)`, propertyID, userID).Scan(&exists)
	
	if err != nil || !exists {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(FloorResponse{false, "Access denied to property", 0})
		return
	}

	// Generate random ID for floor
	floorID, err := utils.GenerateRandomID()
	if err != nil {
		fmt.Printf("Error generating random ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Error generating floor ID", 0})
		return
	}

	// Insert floor into database
	_, err = db.Exec(`
		INSERT INTO floor (id, name, rent, created_at, created_by, updated_at, updated_by, pid)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		floorID,
		req.Name,
		req.Rent,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
		propertyID,
	)
	if err != nil {
		fmt.Printf("Error inserting floor: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Error adding floor", 0})
		return
	}

	fmt.Printf("Successfully added floor ID: %d to property ID: %d\n", floorID, propertyID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(FloorResponse{
		Success: true,
		Message: "Floor added successfully",
		FloorID: floorID,
	})
}

// GetFloorsHandler handles GET requests for floors of a property
func GetFloorsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Floors Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(FloorResponse{false, "User not authenticated", 0})
		return
	}

	// Extract property ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid URL format", 0})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid property ID", 0})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Database connection error", 0})
		return
	}

	// Verify user has access to the property
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE pid = ? AND uid = ?
		)`, propertyID, userID).Scan(&exists)
	
	if err != nil || !exists {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(FloorResponse{false, "Access denied to property", 0})
		return
	}

	// Get all floors for this property
	rows, err := db.Query(`
		SELECT f.id, f.name, f.rent, f.created_at, f.tenant,
		       CASE 
		           WHEN u.name IS NULL OR u.name = '' THEN CONCAT('User ', u.id)
		           ELSE u.name 
		       END as tenant_name,
		       EXISTS(
		           SELECT 1 FROM notification n 
		           WHERE n.fid = f.id AND n.status = 'pending' 
		           AND n.message NOT LIKE 'Advance payment request:%'
		       ) as has_pending_request,
		       EXISTS(
		           SELECT 1 FROM advance a 
		           WHERE a.fid = f.id AND a.status = 'pending'
		       ) as has_pending_advance_payment
		FROM floor f
		LEFT JOIN user u ON f.tenant = u.id
		WHERE f.pid = ?
		ORDER BY f.created_at DESC`, propertyID)
	
	if err != nil {
		fmt.Printf("Error querying floors: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Error fetching floors", 0})
		return
	}
	defer rows.Close()

	var floors []Floor
	for rows.Next() {
		var floor Floor
		var tenant sql.NullInt64
		var tenantName sql.NullString
		var hasPendingRequest bool
		var hasPendingAdvancePayment bool
		if err := rows.Scan(&floor.ID, &floor.Name, &floor.Rent, &floor.CreatedAt, &tenant, &tenantName, &hasPendingRequest, &hasPendingAdvancePayment); err != nil {
			fmt.Printf("Error scanning floor row: %v\n", err)
			continue
		}
		if tenant.Valid {
			floor.Tenant = &tenant.Int64
		}
		if tenantName.Valid {
			floor.TenantName = &tenantName.String
		}
		if hasPendingRequest {
			floor.Status = "pending"
		}
		floor.HasPendingAdvancePayment = hasPendingAdvancePayment
		floors = append(floors, floor)
	}

	fmt.Printf("Found %d floors for property ID: %d\n", len(floors), propertyID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Floors retrieved successfully",
		"floors": floors,
	})
}

// GetFloorByIDHandler handles GET requests for a specific floor
func GetFloorByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Floor By ID Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(FloorResponse{false, "User not authenticated", 0})
		return
	}

	// Extract property ID and floor ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid URL format", 0})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid property ID", 0})
		return
	}

	floorID, err := strconv.ParseInt(cleanParts[3], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid floor ID", 0})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Database connection error", 0})
		return
	}

	// Verify user has access to the property
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE pid = ? AND uid = ?
		)`, propertyID, userID).Scan(&exists)
	
	if err != nil || !exists {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(FloorResponse{false, "Access denied to property", 0})
		return
	}

	// Get floor details
	var floor Floor
	var tenant sql.NullInt64
	var tenantName sql.NullString
	err = db.QueryRow(`
		SELECT f.id, f.name, f.rent, f.created_at, f.tenant, 
		       CASE 
		           WHEN u.name IS NULL OR u.name = '' THEN CONCAT('User ', u.id)
		           ELSE u.name 
		       END as tenant_name
		FROM floor f
		LEFT JOIN user u ON f.tenant = u.id
		WHERE f.id = ? AND f.pid = ?`, floorID, propertyID).Scan(
		&floor.ID, &floor.Name, &floor.Rent, &floor.CreatedAt, &tenant, &tenantName)
	
	if err != nil {
		fmt.Printf("Error querying floor: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(FloorResponse{false, "Floor not found", 0})
		return
	}

	if tenant.Valid {
		floor.Tenant = &tenant.Int64
	}
	if tenantName.Valid {
		floor.TenantName = &tenantName.String
	}

	fmt.Printf("Found floor: ID=%d, Name=%s\n", floor.ID, floor.Name)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Floor retrieved successfully",
		"floor": floor,
	})
}

// UpdateFloorHandler handles PUT requests to update a floor
func UpdateFloorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Update Floor Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(FloorResponse{false, "User not authenticated", 0})
		return
	}

	// Extract property ID and floor ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid URL format", 0})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid property ID", 0})
		return
	}

	floorID, err := strconv.ParseInt(cleanParts[3], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid floor ID", 0})
		return
	}

	var req FloorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Invalid request body", 0})
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FloorResponse{false, "Floor name is required", 0})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Database connection error", 0})
		return
	}

	// Verify user has access to the property
	var exists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE pid = ? AND uid = ?
		)`, propertyID, userID).Scan(&exists)
	
	if err != nil || !exists {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(FloorResponse{false, "Access denied to property", 0})
		return
	}

	// Update floor
	_, err = db.Exec(`
		UPDATE floor 
		SET name = ?, rent = ?, tenant = ?, updated_at = ?, updated_by = ?
		WHERE id = ? AND pid = ?`,
		req.Name, req.Rent, req.Tenant, time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"), userID, floorID, propertyID)
	
	if err != nil {
		fmt.Printf("Error updating floor: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FloorResponse{false, "Error updating floor", 0})
		return
	}

	// If tenant is being added, create a payment record
	if req.Tenant != nil {
		// Generate random ID for payment
		paymentID, err := utils.GenerateRandomID()
		if err != nil {
			fmt.Printf("Error generating payment ID: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(FloorResponse{false, "Error generating payment ID", 0})
			return
		}

		// Insert payment record
		_, err = db.Exec(`
			INSERT INTO payment (
                                id, rent, recieved_money, 
				full_payment, created_at, created_by, updated_at, updated_by,
				fid, uid
                        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			paymentID,
			0, // due_rent
			0, // received_money
			true, // full_payment
			time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
			userID,
			time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
			userID,
			floorID,
			*req.Tenant,
		)

		if err != nil {
			fmt.Printf("Error creating payment record: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(FloorResponse{false, fmt.Sprintf("Error creating payment record: %v", err), 0})
			return
		}

		fmt.Printf("Successfully created payment record for floor ID: %d and tenant ID: %d\n", floorID, *req.Tenant)
	}

	fmt.Printf("Successfully updated floor ID: %d\n", floorID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FloorResponse{
		Success: true,
		Message: "Floor updated successfully",
		FloorID: floorID,
	})
}

// GetUserPhonesHandler handles GET requests for all users' phone numbers
func GetUserPhonesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get User Phones Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(UserPhonesResponse{false, "Method not allowed", nil})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UserPhonesResponse{false, "User not authenticated", nil})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPhonesResponse{false, "Database connection error", nil})
		return
	}

	// Query to get all users' phone numbers
	query := `
		SELECT id, phone_number 
		FROM user 
		WHERE phone_number IS NOT NULL AND phone_number != ''
		ORDER BY id DESC`
	
	fmt.Printf("Executing query: %s\n", query)
	
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error querying users: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPhonesResponse{false, "Error fetching users", nil})
		return
	}
	defer rows.Close()

	var users []UserPhone
	for rows.Next() {
		var user UserPhone
		if err := rows.Scan(&user.ID, &user.Phone); err != nil {
			fmt.Printf("Error scanning user row: %v\n", err)
			continue
		}
		users = append(users, user)
		fmt.Printf("Found user: ID=%d, Phone=%s\n", user.ID, user.Phone)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating user rows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPhonesResponse{false, "Error processing users", nil})
		return
	}

	fmt.Printf("Found %d users with phone numbers\n", len(users))

	response := UserPhonesResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Users:   users,
	}

	// Log the response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Sending response: %s\n", string(responseJSON))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserIDByPhoneHandler handles GET requests to get user ID by phone number
func GetUserIDByPhoneHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get User ID By Phone Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(UserIDResponse{false, "Method not allowed", 0})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UserIDResponse{false, "User not authenticated", 0})
		return
	}

	// Extract phone number from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserIDResponse{false, "Invalid URL format", 0})
		return
	}

	phoneNumber := cleanParts[2]
	if phoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UserIDResponse{false, "Phone number is required", 0})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserIDResponse{false, "Database connection error", 0})
		return
	}

	// Query to get user ID by phone number
	var foundUserID int64
	err = db.QueryRow(`
		SELECT id 
		FROM user 
		WHERE phone_number = ?`, phoneNumber).Scan(&foundUserID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No user found with phone number: %s\n", phoneNumber)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(UserIDResponse{false, "User not found", 0})
			return
		}
		fmt.Printf("Error querying user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserIDResponse{false, "Error fetching user", 0})
		return
	}

	fmt.Printf("Found user ID: %d for phone number: %s\n", foundUserID, phoneNumber)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(UserIDResponse{
		Success: true,
		Message: "User ID retrieved successfully",
		UserID:  foundUserID,
	})
}

// CreatePaymentHandler handles POST requests to create a payment record
func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Create Payment Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Method not allowed", 0})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	fmt.Printf("User ID from session: %d\n", userID)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(PaymentResponse{false, "User not authenticated", 0})
		return
	}

	// Extract property ID and floor ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	fmt.Printf("URL parts: %v\n", cleanParts)

	if len(cleanParts) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Invalid URL format", 0})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Invalid property ID", 0})
		return
	}

	floorID, err := strconv.ParseInt(cleanParts[3], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Invalid floor ID", 0})
		return
	}

	fmt.Printf("Property ID: %d, Floor ID: %d\n", propertyID, floorID)

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Invalid request body", 0})
		return
	}

	fmt.Printf("Request body: %+v\n", req)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Database connection error", 0})
		return
	}

	// Check if user is a manager of the property
	var isManager bool
	managerQuery := "SELECT EXISTS(SELECT 1 FROM takes_care_of WHERE uid = ? AND pid = ?)"
	fmt.Printf("Executing manager check query: %s with userID: %d, propertyID: %d\n", managerQuery, userID, propertyID)
	
	err = db.QueryRow(managerQuery, userID, propertyID).Scan(&isManager)
	if err != nil {
		fmt.Printf("Error checking manager status: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Database error", 0})
		return
	}

	fmt.Printf("Is user %d manager of property %d: %v\n", userID, propertyID, isManager)

	if !isManager {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Only managers can create payments", 0})
		return
	}

	// Get tenant ID from floor
	var tenantID sql.NullInt64
	err = db.QueryRow(`
		SELECT tenant 
		FROM floor 
		WHERE id = ? AND pid = ?`, floorID, propertyID).Scan(&tenantID)
	
	if err != nil {
		fmt.Printf("Error getting tenant ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Error getting tenant information", 0})
		return
	}

	if !tenantID.Valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{false, "No tenant assigned to this floor", 0})
		return
	}
	
	// Calculate after_receiving_money
        afterReceivingMoney := req.ReceivedMoney - req.Rent
	
	// Set full_payment based on after_receiving_money
	fullPayment := afterReceivingMoney == 0

        fmt.Printf("Rent: %d, Received: %d, After receiving: %d, Full payment: %v\n", 
                req.Rent, req.ReceivedMoney, afterReceivingMoney, fullPayment)

	// Always create a new payment record
	paymentID, err := utils.GenerateRandomID()
	if err != nil {
		fmt.Printf("Error generating payment ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PaymentResponse{false, "Error generating payment ID", 0})
		return
	}

	        	// Insert new payment record
	_, err = db.Exec(`
		INSERT INTO payment (
                        id, rent, recieved_money, 
			full_payment, created_at, created_by, updated_at, updated_by,
			fid, uid, electricity_bill, paid_bill
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		paymentID,
                req.Rent,
		req.ReceivedMoney,
		fullPayment,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
		time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
		userID,
		floorID,
		tenantID.Int64,
		req.ElectricityBill,
		0, // paid_bill is 0 (NULL) for new payments
	)

	if err != nil {
		fmt.Printf("Error creating payment record: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PaymentResponse{false, fmt.Sprintf("Error creating payment record: %v", err), 0})
		return
	}

	fmt.Printf("Successfully created new payment record ID: %d for floor ID: %d and tenant ID: %d\n", paymentID, floorID, tenantID.Int64)

	fmt.Printf("Successfully created payment record ID: %d for floor ID: %d and tenant ID: %d\n", paymentID, floorID, tenantID.Int64)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PaymentResponse{
		Success:   true,
		Message:   "Payment record created successfully",
		PaymentID: paymentID,
	})
}

// SendTenantRequestHandler handles POST requests to send a tenant request
func SendTenantRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Send Tenant Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Method not allowed"})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not authenticated"})
		return
	}

	// Extract property ID and floor ID from URL using Gorilla Mux's Vars
	vars := mux.Vars(r)
	propertyID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid floor ID", http.StatusBadRequest)
		return
	}

	// Get phone number from request body
	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Invalid request body"})
		return
	}

	if req.PhoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Phone number is required"})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Database connection error"})
		return
	}

	// Check if user is a manager of the property
	var isManager bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE uid = ? AND pid = ?
		)`, userID, propertyID).Scan(&isManager)
	
	if err != nil || !isManager {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Only managers can send tenant requests"})
		return
	}

	// Get property and floor details
	var propertyName, floorName string
	err = db.QueryRow(`
		SELECT p.name, f.name
		FROM property p
		JOIN floor f ON p.id = f.pid
		WHERE p.id = ? AND f.id = ?`, propertyID, floorID).Scan(&propertyName, &floorName)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error getting property details"})
		return
	}

	// Get tenant ID from phone number
	var tenantID int64
	err = db.QueryRow(`
		SELECT id FROM user 
		WHERE phone_number = ?`, req.PhoneNumber).Scan(&tenantID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not found with this phone number"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error finding user"})
		return
	}

	// Check if there's already a pending tenant request for this floor (excluding advance payment requests)
	var pendingExists bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM notification 
			WHERE fid = ? AND status = 'pending'
			AND message NOT LIKE 'Advance payment request:%'
		)`, floorID).Scan(&pendingExists)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error checking pending notifications"})
		return
	}

	if pendingExists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "A pending request already exists for this floor"})
		return
	}

	// Create notification message
	message := fmt.Sprintf("Tenant request for %s - %s", propertyName, floorName)

	// Create notification with push notification
	err = SendNotificationWithPush(userID, tenantID, propertyID, floorID, message, "pending", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error creating notification"})
		return
	}

	// Update floor status to 'pending' to show "Request Pending"
	// TODO: Uncomment this after adding status column to floor table
	// _, err = db.Exec(`
	// 	UPDATE floor 
	// 	SET status = 'pending', updated_at = NOW(), updated_by = ?
	// 	WHERE id = ?
	// `, userID, floorID)
	
	// if err != nil {
	// 	log.Printf("Error updating floor status: %v", err)
	// 	// Don't fail the entire request if floor status update fails
	// } else {
	// 	log.Printf("Updated floor ID %d status to 'pending'", floorID)
	// }

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TenantRequestResponse{
		Success: true,
			Message: "Tenant request sent successfully",
	})
}

// GetUserNotificationsHandler handles GET requests to get all notifications for a user
func GetUserNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get User Notifications Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Method not allowed", nil})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "User not authenticated", nil})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Database connection error", nil})
		return
	}

	// Get all notifications for the user (only received notifications)
	query := `
		SELECT 
			n.id, n.message, n.status, n.created_at,
			p.id as property_id, p.name as property_name,
			f.id as floor_id, f.name as floor_name,
			CASE 
				WHEN n.status = 'pending' THEN true
				ELSE false
			END as show_actions,
			COALESCE(n.is_read, false) as is_read,
			n.comment,
			n.sender, n.receiver,
			u1.name as sender_name,
			u2.name as receiver_name
		FROM notification n
		JOIN property p ON n.pid = p.id
		JOIN floor f ON n.fid = f.id
		JOIN user u1 ON n.sender = u1.id
		JOIN user u2 ON n.receiver = u2.id
		WHERE n.receiver = ?
		ORDER BY n.created_at DESC
	`
	
	fmt.Printf("Executing notification query for user %d\n", userID)
	
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Printf("Error querying notifications: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Error fetching notifications", nil})
		return
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		var senderID, receiverID int64
		var senderName, receiverName string
		if err := rows.Scan(
			&n.ID, &n.Message, &n.Status, &n.CreatedAt,
			&n.Property.ID, &n.Property.Name,
			&n.Floor.ID, &n.Floor.Name,
			&n.ShowActions,
			&n.IsRead,
			&n.Comment,
			&senderID, &receiverID,
			&senderName, &receiverName,
		); err != nil {
			fmt.Printf("Error scanning notification row: %v\n", err)
			continue
		}
		
		n.SenderID = &senderID
		n.ReceiverID = &receiverID
		n.SenderName = &senderName
		n.ReceiverName = &receiverName
		
		// Debug logging for each notification
		fmt.Printf("Notification: ID=%d, Message='%s', Status='%s', ShowActions=%v, IsRead=%v, Comment=%v\n", 
			n.ID, n.Message, n.Status, n.ShowActions, n.IsRead, n.Comment)
		
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating notification rows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Error processing notifications", nil})
		return
	}

	fmt.Printf("Found %d notifications for user %d\n", len(notifications), userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NotificationsResponse{
		Success:       true,
		Message:       "Notifications retrieved successfully",
		Notifications: notifications,
	})
}

// DeleteNotificationHandler handles DELETE requests to remove a notification
func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Delete Notification Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Method not allowed"})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not authenticated"})
		return
	}

	// Extract notification ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Invalid URL format"})
		return
	}

	notificationID, err := strconv.ParseInt(cleanParts[2], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Invalid notification ID"})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Database connection error"})
		return
	}

	// Check if user is the sender or receiver of the notification and if it's pending
	var canDelete bool
	var floorID int64
	var message string
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM notification 
			WHERE id = ? AND (sender = ? OR receiver = ?) AND status = 'pending'
		)`, notificationID, userID, userID).Scan(&canDelete)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error checking notification"})
		return
	}

	if !canDelete {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "You can only delete your own pending notifications"})
		return
	}

	// Get the notification details before deleting
	err = db.QueryRow(`
		SELECT fid, message
		FROM notification 
		WHERE id = ? AND (sender = ? OR receiver = ?) AND status = 'pending'`, notificationID, userID, userID).Scan(&floorID, &message)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error getting notification details"})
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error starting transaction"})
		return
	}
	defer tx.Rollback()

	// Delete the notification
	_, err = tx.Exec(`
		DELETE FROM notification 
		WHERE id = ? AND (sender = ? OR receiver = ?) AND status = 'pending'`,
		notificationID, userID, userID)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error deleting notification"})
		return
	}

	// Check if this was a payment notification and clear floor status if no more pending notifications
	if strings.HasPrefix(message, "Payment amount:") {
		// Check if there are any remaining pending notifications for this floor (excluding advance payment requests)
		var remainingNotifications int
		err = tx.QueryRow(`
			SELECT COUNT(*) 
			FROM notification 
			WHERE fid = ? AND status = 'pending' 
			AND message NOT LIKE 'Advance payment request:%'`, floorID).Scan(&remainingNotifications)
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error checking remaining notifications"})
			return
		}

		// If no more pending notifications, clear the floor status
		if remainingNotifications == 0 {
			// TODO: Uncomment this after adding status column to floor table
			// _, err = tx.Exec(`
			// 	UPDATE floor 
			// 	SET status = NULL, updated_at = NOW(), updated_by = ?
			// 		WHERE id = ?`, userID, floorID)
			
			// if err != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error clearing floor status"})
			// 	return
			// }
			
			// fmt.Printf("Cleared pending status from floor ID: %d after deleting payment notification", floorID)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error committing transaction"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TenantRequestResponse{
		Success: true,
		Message: "Notification deleted successfully",
	})
}

// SendMonthlyNotifications sends notifications to all tenants on the 5th of each month
func SendMonthlyNotifications() {
	// Get current date in Bangladesh timezone
	loc, _ := time.LoadLocation("Asia/Dhaka")
	now := time.Now().In(loc)

	// Only send notifications if the current hour and minute match (for immediate test)
	
	// Only send notifications on the 5th of each month at 9:00 AM
	if !(now.Day() == 5 && now.Hour() == 9 && now.Minute() == 0) {
        return
    }

	fmt.Println("=== Sending Monthly Notifications ===")

	// Get database connection
	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		return
	}

	// Get all floors with tenants
	query := `
		SELECT f.id, f.pid, f.tenant, f.name, p.name as property_name
		FROM floor f
		JOIN property p ON f.pid = p.id
		WHERE f.tenant IS NOT NULL
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error querying floors: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var floorID, propertyID, tenantID int64
		var floorName, propertyName string
		if err := rows.Scan(&floorID, &propertyID, &tenantID, &floorName, &propertyName); err != nil {
			fmt.Printf("Error scanning floor: %v\n", err)
			continue
		}

		// Get latest payment for this floor
		var rent float64
		paymentQuery := `
			SELECT rent
			FROM payment
			WHERE fid = ?
			ORDER BY created_at DESC
			LIMIT 1
		`
		err := db.QueryRow(paymentQuery, floorID).Scan(&rent)
		if err != nil {
			if err == sql.ErrNoRows {
				// No payment record found, use default values
				rent = 0
			} else {
			fmt.Printf("Error querying payment: %v\n", err)
			continue
			}
		}

		// Get manager's user id (sender)
		var managerID int64
		managerQuery := `SELECT uid FROM takes_care_of WHERE pid = ? LIMIT 1`
		err = db.QueryRow(managerQuery, propertyID).Scan(&managerID)
		if err != nil {
			fmt.Printf("Error getting manager for property %d: %v\n", propertyID, err)
			continue
		}

		// Compose notification message
		message := fmt.Sprintf(
			"Monthly rent reminder for %s - %s:\nDue Rent: %.2f tk",
			propertyName, floorName, rent,
		)

		// Create notification with push notification
		err = SendNotificationWithPush(managerID, tenantID, propertyID, floorID, message, "", nil)
		if err != nil {
			fmt.Printf("Error creating notification: %v\n", err)
			continue
		}

		fmt.Printf("Created notification for tenant %d in property %s, floor %s\n", tenantID, propertyName, floorName)
	}

	fmt.Println("Monthly notifications sent successfully.")
}

// TestSendNotifications is a test function to manually trigger notifications
func TestSendNotifications() {
	fmt.Println("=== Testing Send Notifications ===")

	// Get database connection
	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		return
	}

	// Get all floors with tenants
	query := `
		SELECT f.id, f.pid, f.tenant, f.name, p.name as property_name
		FROM floor f
		JOIN property p ON f.pid = p.id
		WHERE f.tenant IS NOT NULL
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error querying floors: %v\n", err)
		return
	}
	defer rows.Close()

	notificationCount := 0
	for rows.Next() {
		var floorID, propertyID, tenantID int64
		var floorName, propertyName string
		if err := rows.Scan(&floorID, &propertyID, &tenantID, &floorName, &propertyName); err != nil {
			fmt.Printf("Error scanning floor: %v\n", err)
			continue
		}

		// Get latest payment for this floor
		var rent float64
		paymentQuery := `
			SELECT rent
			FROM payment
			WHERE fid = ?
			ORDER BY created_at DESC
			LIMIT 1
		`
		err := db.QueryRow(paymentQuery, floorID).Scan(&rent)
		if err != nil {
			if err == sql.ErrNoRows {
				// No payment record found, use default values
				rent = 0
			} else {
				fmt.Printf("Error querying payment: %v\n", err)
				continue
			}
		}

		// Get manager's user id (sender)
		var managerID int64
		managerQuery := `SELECT uid FROM takes_care_of WHERE pid = ? LIMIT 1`
		err = db.QueryRow(managerQuery, propertyID).Scan(&managerID)
		if err != nil {
			fmt.Printf("Error getting manager for property %d: %v\n", propertyID, err)
			continue
		}

		// Compose notification message
		message := fmt.Sprintf(
			"TEST: Monthly rent reminder for %s - %s:\nDue Rent: %.2f tk",
			propertyName, floorName, rent,
		)

		// Create notification with push notification
		err = SendNotificationWithPush(managerID, tenantID, propertyID, floorID, message, "", nil)
		if err != nil {
			fmt.Printf("Error creating notification: %v\n", err)
			continue
		}

		fmt.Printf("Created test notification for tenant %d in property %s, floor %s\n", tenantID, propertyName, floorName)
		notificationCount++
	}

	fmt.Printf("Test completed. Created %d notifications successfully.\n", notificationCount)
}

// TestSendNotificationsHandler handles the test endpoint to manually trigger notifications
func TestSendNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== Test Send Notifications Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
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

	// Call the test function
	TestSendNotifications()

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Test notifications sent successfully",
	})
}

// HandleTenantRequestAction handles POST requests to accept/reject tenant requests and payment notifications
func HandleTenantRequestAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Notification Action Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Method not allowed"})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not authenticated"})
		return
	}

	// Parse request body
	var request struct {
		NotificationID int64 `json:"notification_id"`
		Accept        bool  `json:"accept"`
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
	if db == nil {
		http.Error(w, "Database connection is nil", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Transaction start error: %v\n", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get notification details
	var notification struct {
		ID      int64
		Message string
		Status  string
		FloorID int64
		PID     int64
		Sender  int64
		Receiver int64
	}

	err = tx.QueryRow(`
		SELECT n.id, n.message, n.status, n.fid, n.pid, n.sender, n.receiver
		FROM notification n
		WHERE n.id = ? AND n.receiver = ?
	`, request.NotificationID, userID).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Status,
		&notification.FloorID,
		&notification.PID,
		&notification.Sender,
		&notification.Receiver,
	)

	if err != nil {
		fmt.Printf("Error getting notification: %v\n", err)
		if err == sql.ErrNoRows {
			http.Error(w, "Notification not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get notification", http.StatusInternalServerError)
		}
		return
	}

	if notification.Status != "pending" {
		http.Error(w, "Notification is not pending", http.StatusBadRequest)
		return
	}

	// Update notification status
	newStatus := "rejected"
	if request.Accept {
		newStatus = "accepted"
	}

	_, err = tx.Exec(`
		UPDATE notification 
		SET status = ?, updated_at = NOW(), updated_by = ?
		WHERE id = ?
	`, newStatus, userID, request.NotificationID)

	if err != nil {
		fmt.Printf("Error updating notification: %v\n", err)
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	// Check if this is a payment notification, advance payment notification, or tenant request
	isPaymentNotification := strings.HasPrefix(notification.Message, "Payment amount:")
	isAdvancePaymentNotification := strings.HasPrefix(notification.Message, "Advance payment request:")
	
	if isPaymentNotification {
		// Handle payment notification
		if request.Accept {
			// Payment accepted - create payment record
			fmt.Printf("Payment notification accepted: ID=%d, creating payment record", notification.ID)
			
			// Extract amount from the notification message
			// Message format: "Payment amount: X tk" or "Payment amount: X tk for MonthName"
			amountRegex := regexp.MustCompile(`Payment amount:\s*([0-9]+)\s*tk`)
			
			amountMatches := amountRegex.FindStringSubmatch(notification.Message)
			
			var amount int
			
			if len(amountMatches) >= 2 {
				if parsedAmount, err := strconv.Atoi(amountMatches[1]); err == nil {
					amount = parsedAmount
				} else {
					fmt.Printf("Error parsing amount from message: %v\n", err)
					http.Error(w, "Failed to parse payment amount", http.StatusInternalServerError)
					return
				}
			} else {
				fmt.Printf("Could not extract amount from message: %s\n", notification.Message)
				http.Error(w, "Failed to extract payment amount from message", http.StatusInternalServerError)
				return
			}
			
			// Extract electricity bill from the notification message if present
			// Message format: "Payment amount: X tk, Paid electricity bill: Y tk"
			electricityBillRegex := regexp.MustCompile(`Paid electricity bill:\s*([0-9]+)\s*tk`)
			electricityBillMatches := electricityBillRegex.FindStringSubmatch(notification.Message)
			
			var electricityBill *int
			if len(electricityBillMatches) >= 2 {
				if parsedElectricityBill, err := strconv.Atoi(electricityBillMatches[1]); err == nil {
					electricityBill = &parsedElectricityBill
					fmt.Printf("Extracted electricity bill: %d tk", parsedElectricityBill)
				}
			}
			
			// Get tenant ID from floor table
			var tenantID int64
			err = tx.QueryRow("SELECT tenant FROM floor WHERE id = ?", notification.FloorID).Scan(&tenantID)
			if err != nil {
				fmt.Printf("Error getting tenant ID: %v\n", err)
				http.Error(w, "Failed to get tenant information", http.StatusInternalServerError)
				return
			}
			
			// Generate payment ID
			paymentID, err := utils.GenerateRandomID()
			if err != nil {
				fmt.Printf("Error generating payment ID: %v\n", err)
				http.Error(w, "Failed to generate payment ID", http.StatusInternalServerError)
				return
			}
			
			// Create payment record according to requirements:
			// rent = 0, received_money = amount from message, electricity_bill = 0, paid_bill = electricity bill from message
			_, err = tx.Exec(`
				INSERT INTO payment (
					id, rent, recieved_money, full_payment, created_at, created_by, updated_at, updated_by, 
					fid, uid, electricity_bill, paid_bill
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
				paymentID,
				0, // rent = 0 (as specified)
				amount, // received_money = amount from message
				true, // full_payment = true since it's accepted
				time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
				userID,
				time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02 15:04:05"),
				userID,
				notification.FloorID, // fid from notification
				tenantID, // uid from floor.tenant
				0, // electricity_bill = 0
				electricityBill, // paid_bill = electricity bill from message
			)
			
			if err != nil {
				fmt.Printf("Error creating payment record: %v\n", err)
				http.Error(w, "Failed to create payment record", http.StatusInternalServerError)
				return
			}
			
			fmt.Printf("Created payment record: ID=%d, Floor=%d, Tenant=%d, Amount=%d, ElectricityBill=%v", 
				paymentID, notification.FloorID, tenantID, amount, electricityBill)
		}
		// If rejected, just update the notification status (already done above)
		
	} else if isAdvancePaymentNotification {
		// Handle advance payment notification
		if request.Accept {
			// Advance payment accepted - update the advance record status
			_, err = tx.Exec(`
				UPDATE advance 
				SET status = 'accepted', updated_at = NOW(), updated_by = ?
				WHERE fid = ? AND status = 'pending'
			`, userID, notification.FloorID)
			
			if err != nil {
				fmt.Printf("Error updating advance payment status: %v\n", err)
				http.Error(w, "Failed to update advance payment status", http.StatusInternalServerError)
				return
			}
			
			fmt.Printf("Advance payment accepted: notification ID=%d, floor=%d", notification.ID, notification.FloorID)
		} else {
			// Advance payment rejected - update the advance record status
			_, err = tx.Exec(`
				UPDATE advance 
				SET status = 'rejected', updated_at = NOW(), updated_by = ?
				WHERE fid = ? AND status = 'pending'
			`, userID, notification.FloorID)
			
			if err != nil {
				fmt.Printf("Error updating advance payment status: %v\n", err)
				http.Error(w, "Failed to update advance payment status", http.StatusInternalServerError)
				return
			}
			
			fmt.Printf("Advance payment rejected: notification ID=%d, floor=%d", notification.ID, notification.FloorID)
		}
		
	} else {
		// Handle tenant request
		if request.Accept {
			// Check if floor is already occupied
			var isOccupied bool
			err = tx.QueryRow(`
				SELECT COUNT(*) > 0
				FROM floor 
				WHERE id = ? AND tenant IS NOT NULL
			`, notification.FloorID).Scan(&isOccupied)

			if err != nil {
				fmt.Printf("Error checking floor status: %v\n", err)
				http.Error(w, "Failed to check floor status", http.StatusInternalServerError)
				return
			}

			if isOccupied {
				http.Error(w, "Floor is already occupied", http.StatusConflict)
				return
			}

			// Update floor with tenant (receiver of the notification, not sender)
			_, err = tx.Exec(`
				UPDATE floor 
				SET tenant = ?, updated_at = NOW(), updated_by = ?
				WHERE id = ?
			`, notification.Receiver, userID, notification.FloorID)

			if err != nil {
				fmt.Printf("Error updating floor: %v\n", err)
				http.Error(w, "Failed to update floor", http.StatusInternalServerError)
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Create auto-generated response notification
	{
		// Create auto-generated response notification
		var responseMessage string
		if isPaymentNotification {
			// Extract amount from the original message for payment notifications
			// Original message format: "Payment amount: X tk" or "Payment amount: $X"
			amountRegex := regexp.MustCompile(`Payment amount:\s*([0-9]+)\s*(tk|\$|৳)`)
			matches := amountRegex.FindStringSubmatch(notification.Message)
			if len(matches) >= 2 {
				amount := matches[1]
				if request.Accept {
					responseMessage = fmt.Sprintf("Payment of %s tk is accepted", amount)
				} else {
					responseMessage = fmt.Sprintf("Payment of %s tk is rejected", amount)
				}
			} else {
				if request.Accept {
					responseMessage = "Payment request is accepted"
				} else {
					responseMessage = "Payment request is rejected"
				}
			}
		} else if isAdvancePaymentNotification {
			// Extract amount from the original message for advance payment notifications
			// Original message format: "Advance payment request: X tk"
			advanceRegex := regexp.MustCompile(`Advance payment request:\s*([0-9]+)\s*tk`)
			matches := advanceRegex.FindStringSubmatch(notification.Message)
			if len(matches) >= 2 {
				amount := matches[1]
				if request.Accept {
					responseMessage = fmt.Sprintf("Advance payment of %s tk is accepted", amount)
				} else {
					responseMessage = fmt.Sprintf("Advance payment of %s tk is rejected", amount)
				}
			} else {
				if request.Accept {
					responseMessage = "Advance payment request is accepted"
				} else {
					responseMessage = "Advance payment request is rejected"
				}
			}
		} else {
			// Handle tenant request
			if request.Accept {
				responseMessage = "Tenant request is accepted"
			} else {
				responseMessage = "Tenant request is rejected"
			}
		}
		
		// Determine sender and receiver for the response notification
		// The response should go from the person who took the action to the original sender
		responseSender := notification.Receiver // Person who accepted/rejected
		responseReceiver := notification.Sender // Original sender of the request
		
		// Create auto-response notification with push notification
		err = SendNotificationWithPush(responseSender, responseReceiver, notification.PID, notification.FloorID, responseMessage, newStatus, nil)
		if err != nil {
			fmt.Printf("Error creating auto-response notification: %v\n", err)
			// Don't fail the whole request, just log the error
		} else {
			fmt.Printf("Created auto-response notification with push: Message='%s', from user=%d to user=%d", 
				responseMessage, responseSender, responseReceiver)
		}
	}

	// Send response
	actionType := "payment notification"
	if isAdvancePaymentNotification {
		actionType = "advance payment notification"
	} else if !isPaymentNotification {
		actionType = "tenant request"
	}
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": actionType + " " + newStatus + " (auto-response notification sent, use /notifications/send-comment to add comments)",
	})
}

// GetAdvanceDetailsHandler handles GET requests to get advance details for a floor
func GetAdvanceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Advance Details Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AdvanceDetailsResponse{false, "Method not allowed", nil})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AdvanceDetailsResponse{false, "User not authenticated", nil})
		return
	}

	// Extract floor ID from URL
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvanceDetailsResponse{false, "Invalid floor ID", nil})
		return
	}

	fmt.Printf("Fetching advance details for floor ID: %d\n", floorID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvanceDetailsResponse{false, "Database connection error", nil})
		return
	}

	// Query to get advance details for the floor where money is not zero (including pending status)
	query := `
		SELECT a.id, a.advance_uid, u.name, a.money, a.created_at, a.status
		FROM advance a
		INNER JOIN user u ON a.advance_uid = u.id
		WHERE a.fid = ? AND a.money > 0
		ORDER BY a.created_at DESC`
	
	fmt.Printf("Executing query: %s with floorID: %d\n", query, floorID)
	
	rows, err := db.Query(query, floorID)
	if err != nil {
		fmt.Printf("Error querying advance details: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvanceDetailsResponse{false, "Error fetching advance details", nil})
		return
	}
	defer rows.Close()

	var advances []AdvanceDetail
	for rows.Next() {
		var advance AdvanceDetail
		if err := rows.Scan(&advance.ID, &advance.AdvanceUID, &advance.UserName, &advance.Money, &advance.CreatedAt, &advance.Status); err != nil {
			fmt.Printf("Error scanning advance detail row: %v\n", err)
			continue
		}
		advances = append(advances, advance)
		fmt.Printf("Found advance: ID=%d, User=%s, Money=%d\n", advance.ID, advance.UserName, advance.Money)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating advance detail rows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvanceDetailsResponse{false, "Error processing advance details", nil})
		return
	}

	fmt.Printf("Found %d advance details for floor\n", len(advances))

	response := AdvanceDetailsResponse{
		Success: true,
		Message: "Advance details retrieved successfully",
		Advances: advances,
	}

	// Log the response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Sending response: %s\n", string(responseJSON))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserTenantPropertiesHandler handles GET requests to get all properties where the user is a tenant
func GetUserTenantPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get User Tenant Properties Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Method not allowed", nil})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "User not authenticated", nil})
		return
	}

	fmt.Printf("Fetching properties where user ID: %d is a tenant\n", userID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Database connection error", nil})
		return
	}

	// Query to get all properties where the user is a tenant
	query := `
		SELECT DISTINCT p.id, p.name, p.address, p.photo, p.created_at 
		FROM property p
		INNER JOIN floor f ON p.id = f.pid
		WHERE f.tenant = ?
		ORDER BY p.created_at DESC`
	
	fmt.Printf("Executing query: %s with userID: %d\n", query, userID)
	
	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Printf("Error querying properties: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Error fetching properties", nil})
		return
	}
	defer rows.Close()

	var properties []Property
	for rows.Next() {
		var prop Property
		var photo sql.NullString
		if err := rows.Scan(&prop.ID, &prop.Name, &prop.Address, &photo, &prop.CreatedAt); err != nil {
			fmt.Printf("Error scanning property row: %v\n", err)
			continue
		}
		if photo.Valid {
			prop.Photo = &photo.String
		}
		properties = append(properties, prop)
		fmt.Printf("Found property: ID=%d, Name=%s\n", prop.ID, prop.Name)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating property rows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UserPropertiesResponse{false, "Error processing properties", nil})
		return
	}

	fmt.Printf("Found %d properties where user is a tenant\n", len(properties))

	response := UserPropertiesResponse{
		Success: true,
		Message: "Properties retrieved successfully",
		Properties: properties,
	}

	// Log the response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Sending response: %s\n", string(responseJSON))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RemoveTenantHandler handles POST requests to remove a tenant from a floor
func RemoveTenantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Method not allowed"})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	vars := mux.Vars(r)
	propertyID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}
	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid floor ID", http.StatusBadRequest)
		return
	}

	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not authenticated"})
		return
	}

	// Get database connection
	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	if db == nil {
		http.Error(w, "Database connection is nil", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Transaction start error: %v\n", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Check if user is the manager of the property
	var isManager bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE uid = ? AND pid = ?
		)`, userID, propertyID).Scan(&isManager)

	if err != nil {
		fmt.Printf("Error checking property manager: %v\n", err)
		http.Error(w, "Failed to check authorization", http.StatusInternalServerError)
		return
	}

	if !isManager {
		http.Error(w, "You are not authorized to manage this property", http.StatusForbidden)
		return
	}

	// Check if there's a tenant in the floor
	var hasTenant bool
	err = tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM floor 
			WHERE id = ? AND pid = ? AND tenant IS NOT NULL
		)`, floorID, propertyID).Scan(&hasTenant)

	if err != nil {
		fmt.Printf("Error checking tenant: %v\n", err)
		http.Error(w, "Failed to check tenant", http.StatusInternalServerError)
		return
	}

	if !hasTenant {
		http.Error(w, "No tenant found in this floor", http.StatusBadRequest)
		return
	}

	// Update floor to remove tenant
	_, err = tx.Exec(`
		UPDATE floor 
		SET tenant = NULL, updated_at = NOW(), updated_by = ?
		WHERE id = ? AND pid = ?
	`, userID, floorID, propertyID)

	if err != nil {
		fmt.Printf("Error removing tenant: %v\n", err)
		http.Error(w, "Failed to remove tenant", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tenant removed successfully",
	})
}

type ManagerCheckResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	IsManager bool   `json:"is_manager"`
}

// CheckUserManagerHandler handles GET requests to check if user is a manager of a property
func CheckUserManagerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Check User Manager Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ManagerCheckResponse{false, "Method not allowed", false})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ManagerCheckResponse{false, "User not authenticated", false})
		return
	}

	// Extract property ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ManagerCheckResponse{false, "Invalid URL format", false})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ManagerCheckResponse{false, "Invalid property ID", false})
		return
	}

	fmt.Printf("Checking if user ID: %d is manager of property ID: %d\n", userID, propertyID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ManagerCheckResponse{false, "Database connection error", false})
		return
	}

	// Check if user is a manager of the property
	var isManager bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM takes_care_of 
			WHERE uid = ? AND pid = ?
		)`, userID, propertyID).Scan(&isManager)
	
	if err != nil {
		fmt.Printf("Error checking manager status: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ManagerCheckResponse{false, "Error checking manager status", false})
		return
	}

	fmt.Printf("User %d is manager of property %d: %v\n", userID, propertyID, isManager)

	response := ManagerCheckResponse{
		Success:   true,
		Message:   "Manager check completed",
		IsManager: isManager,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SendPaymentNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	propertyID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}
	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid floor ID", http.StatusBadRequest)
		return
	}

	var req PaymentNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromContext(r) // sender (tenant)
	db, err := config.GetDBConnection()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if user is the tenant of this floor
	var isTenant bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM floor WHERE id = ? AND pid = ? AND tenant = ?)`, floorID, propertyID, userID).Scan(&isTenant)
	log.Printf("[Tenant Check] floorID=%d, propertyID=%d, userID=%d, isTenant=%v, err=%v", floorID, propertyID, userID, isTenant, err)
	
	if err != nil || !isTenant {
		if err != nil {
			log.Printf("[Tenant Check] DB error: %v", err)
		}
		if !isTenant {
			log.Printf("[Tenant Check] User is not the tenant of this floor.")
		}
		http.Error(w, "You are not the tenant of this floor", http.StatusForbidden)
		return
	}

	// Get property manager (receiver)
	var managerID int64
	err = db.QueryRow(`SELECT uid FROM takes_care_of WHERE pid = ? LIMIT 1`, propertyID).Scan(&managerID)
	if err != nil {
		http.Error(w, "Manager not found", http.StatusInternalServerError)
		return
	}

	// Get property and floor names
	var propertyName, floorName string
	err = db.QueryRow(`SELECT name FROM property WHERE id = ?`, propertyID).Scan(&propertyName)
	if err != nil {
		http.Error(w, "Property not found", http.StatusInternalServerError)
		return
	}
	err = db.QueryRow(`SELECT name FROM floor WHERE id = ?`, floorID).Scan(&floorName)
	if err != nil {
		http.Error(w, "Floor not found", http.StatusInternalServerError)
		return
	}

	// Compose message with structured format
	var message string
	if req.Month != nil {
		monthNames := []string{"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December"}
		monthName := monthNames[*req.Month-1]
		message = fmt.Sprintf("Payment amount: %d tk for %s", req.Amount, monthName)
	} else {
		message = fmt.Sprintf("Payment amount: %d tk", req.Amount)
	}
	
	// Add electricity bill info to message if provided
	if req.PaidElectricityBill != nil {
		message += fmt.Sprintf(", Paid electricity bill: %d tk", *req.PaidElectricityBill)
	}

	// Generate notification ID
	notificationID, err := utils.GenerateRandomID()
	if err != nil {
		http.Error(w, "Error generating notification ID", http.StatusInternalServerError)
		return
	}

	// Create notification with push notification
	err = SendNotificationWithPush(userID, managerID, propertyID, floorID, message, "pending", nil)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	// Update floor status to 'pending' to show "Request Pending"
	// TODO: Uncomment this after adding status column to floor table
	// _, err = db.Exec(`
	// 	UPDATE floor 
	// 	SET status = 'pending', updated_at = NOW(), updated_by = ?
	// 	WHERE id = ?
	// `, userID, floorID)
	
	// if err != nil {
	// 	log.Printf("Error updating floor status: %v", err)
	// 	// Don't fail the entire request if floor status update fails
	// } else {
	// 	log.Printf("Updated floor ID %d status to 'pending'", floorID)
	// }

	// Note: Payment records are now only created when the manager accepts the notification
	// This prevents duplicate payment records from being created

	log.Printf("Successfully sent payment notification: ID=%d, from user=%d to manager=%d", notificationID, userID, managerID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Payment notification sent"})
}

// HandlePaymentNotificationAction handles POST requests to accept/reject payment notifications
func HandlePaymentNotificationAction(w http.ResponseWriter, r *http.Request) {
	// This function is no longer needed as payment notifications are handled in HandleTenantRequestAction
	http.Error(w, "Use /notifications/action endpoint instead", http.StatusMovedPermanently)
}

// GetPaymentDetailsHandler handles GET requests to get payment details for a floor
func GetPaymentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Payment Details Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
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

	// Extract floor ID from URL using Gorilla Mux variables
	vars := mux.Vars(r)
	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid floor ID",
		})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database connection error",
		})
		return
	}

	// Get the tenant ID for this floor
	var tenantID sql.NullInt64
	err = db.QueryRow(`SELECT tenant FROM floor WHERE id = ?`, floorID).Scan(&tenantID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Floor not found",
		})
		return
	}

	// If no tenant is assigned, return 0 outstanding rent
	if !tenantID.Valid {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Payment details retrieved successfully",
			"payment": map[string]interface{}{
				"rent":           0,
				"received_money": 0,
				"full_payment":   false,
			},
		})
		return
	}

	// Calculate total outstanding rent: sum(rent - received_money) for all payments of this floor and tenant
	var totalOutstandingRent sql.NullFloat64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(rent - recieved_money), 0) as total_outstanding
		FROM payment
		WHERE fid = ? AND uid = ?
	`, floorID, tenantID.Int64).Scan(&totalOutstandingRent)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error calculating outstanding rent",
		})
		return
	}

	// Get the latest payment record for this floor (for other details)
	var rent, receivedMoney sql.NullInt64
	var fullPayment bool
	var paymentID int64
	var createdAt string
	err = db.QueryRow(`
		SELECT id, rent, recieved_money, full_payment, created_at
		FROM payment
		WHERE fid = ?
		ORDER BY created_at DESC
		LIMIT 1
	`, floorID).Scan(&paymentID, &rent, &receivedMoney, &fullPayment, &createdAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// No payment record found, return default values
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"message": "Payment details retrieved successfully",
				"payment": map[string]interface{}{
					"rent":           0,
					"received_money": 0,
					"full_payment":   false,
				},
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error retrieving payment details",
		})
		return
	}

			// Convert NULL values to 0 if they are NULL
		var rentValue, receivedMoneyValue int64
		if rent.Valid {
			rentValue = rent.Int64
		} else {
			rentValue = 0
		}
		if receivedMoney.Valid {
			receivedMoneyValue = receivedMoney.Int64
		} else {
			receivedMoneyValue = 0
		}

		// Use the calculated total outstanding rent
		totalOutstanding := int64(0)
		if totalOutstandingRent.Valid {
			totalOutstanding = int64(totalOutstandingRent.Float64)
		}

		fmt.Printf("Payment details for floor %d: Total Outstanding=%d, Latest Rent=%d, Received=%d, Created=%s\n",
				floorID, totalOutstanding, rentValue, receivedMoneyValue, createdAt)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Payment details retrieved successfully",
			"payment": map[string]interface{}{
				"rent":          totalOutstanding, // Use total outstanding rent instead of latest payment rent
				"received_money": receivedMoneyValue,
				"full_payment":   fullPayment,
			},
		})
}

// GetPendingPaymentNotificationsHandler handles GET requests to get pending payment notifications for a floor
func GetPendingPaymentNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Pending Payment Notifications Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Method not allowed", nil})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "User not authenticated", nil})
		return
	}

	// Extract property ID and floor ID from URL using Gorilla Mux's Vars
	vars := mux.Vars(r)
	propertyID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid floor ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("Fetching pending payment notifications for floor ID: %d, property ID: %d, user ID: %d\n", floorID, propertyID, userID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Database connection error", nil})
		return
	}

	// Check if user is the tenant of this floor
	var isTenant bool
	err = db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM floor 
			WHERE id = ? AND pid = ? AND tenant = ?
		)`, floorID, propertyID, userID).Scan(&isTenant)
	
	if err != nil {
		fmt.Printf("Error checking tenant status: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Error checking tenant status", nil})
		return
	}

	if !isTenant {
		fmt.Printf("User %d is not a tenant of floor %d in property %d\n", userID, floorID, propertyID)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "You are not a tenant of this floor", nil})
		return
	}

	// Get pending payment notifications for this floor where the user is the sender
	query := `
		SELECT 
			n.id, n.message, n.status, n.created_at,
			p.id as property_id, p.name as property_name,
			f.id as floor_id, f.name as floor_name,
			false as show_actions
		FROM notification n
		JOIN property p ON n.pid = p.id
		JOIN floor f ON n.fid = f.id
		WHERE n.fid = ? 
		AND n.sender = ? 
		AND n.message LIKE 'Payment amount:%'
		AND n.status = 'pending'
		ORDER BY n.created_at DESC
	`
	
	fmt.Printf("Executing query: %s with floorID: %d, userID: %d\n", query, floorID, userID)
	
	rows, err := db.Query(query, floorID, userID)
	if err != nil {
		fmt.Printf("Error querying notifications: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Error fetching notifications", nil})
		return
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(
			&n.ID, &n.Message, &n.Status, &n.CreatedAt,
			&n.Property.ID, &n.Property.Name,
			&n.Floor.ID, &n.Floor.Name,
			&n.ShowActions,
		); err != nil {
			fmt.Printf("Error scanning notification row: %v\n", err)
			continue
		}
		
		// Debug logging for each notification
		fmt.Printf("Payment notification: ID=%d, Message='%s', Status='%s'\n", 
			n.ID, n.Message, n.Status)
		
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating notification rows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(NotificationsResponse{false, "Error processing notifications", nil})
		return
	}

	fmt.Printf("Found %d pending payment notifications for floor %d\n", len(notifications), floorID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NotificationsResponse{
		Success:       true,
		Message:       "Pending payment notifications retrieved successfully",
		Notifications: notifications,
	})
}

// MarkNotificationsAsReadHandler handles POST requests to mark notifications as read
func MarkNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Mark Notifications As Read Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Method not allowed"})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not authenticated"})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Database connection error"})
		return
	}

	// Mark all notifications for this user as read
	_, err = db.Exec(`
		UPDATE notification 
		SET is_read = true, updated_at = NOW(), updated_by = ?
		WHERE receiver = ? AND is_read = false
	`, userID, userID)

	if err != nil {
		fmt.Printf("Error marking notifications as read: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error marking notifications as read"})
		return
	}

	fmt.Printf("Marked all notifications as read for user %d\n", userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TenantRequestResponse{
		Success: true,
		Message: "Notifications marked as read successfully",
	})
}

// AddTenantToFloorHandler handles POST requests to add a tenant to a floor
func AddTenantToFloorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Method not allowed"})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not authenticated"})
		return
	}

	// Extract property ID and floor ID from URL
	vars := mux.Vars(r)
	propertyID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Invalid property ID"})
		return
	}
	floorID, err := strconv.ParseInt(vars["floor_id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Invalid floor ID"})
		return
	}

	// Parse request body
	var req struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Invalid request body"})
		return
	}
	if req.Name == "" || req.PhoneNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Name and phone number are required"})
		return
	}

	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Database connection error"})
		return
	}

	// Find tenant user by phone number
	var tenantID int64
	err = db.QueryRow(`SELECT id FROM user WHERE phone_number = ?`, req.PhoneNumber).Scan(&tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(TenantRequestResponse{false, "User not found with this phone number"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error finding user"})
		return
	}

	// Update floor with tenant
	_, err = db.Exec(`UPDATE floor SET tenant = ?, updated_at = NOW(), updated_by = ? WHERE id = ? AND pid = ?`, tenantID, userID, floorID, propertyID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(TenantRequestResponse{false, "Error updating floor with tenant"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TenantRequestResponse{true, "Tenant added to floor successfully"})
}

// SendCommentHandler handles POST requests to send comments to notifications
func SendCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Send Comment Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
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
		NotificationID int64  `json:"notification_id"`
		Comment        string `json:"comment"`
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

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("Transaction start error: %v\n", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Get original notification details
	var originalNotification struct {
		ID       int64
		Message  string
		Status   string
		Sender   int64
		Receiver int64
		PID      int64
		FloorID  int64
	}

	err = tx.QueryRow(`
		SELECT id, message, status, sender, receiver, pid, fid
		FROM notification
		WHERE id = ? AND (sender = ? OR receiver = ?)
	`, request.NotificationID, userID, userID).Scan(
		&originalNotification.ID,
		&originalNotification.Message,
		&originalNotification.Status,
		&originalNotification.Sender,
		&originalNotification.Receiver,
		&originalNotification.PID,
		&originalNotification.FloorID,
	)

	if err != nil {
		fmt.Printf("Error getting notification: %v\n", err)
		if err == sql.ErrNoRows {
			http.Error(w, "Notification not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get notification", http.StatusInternalServerError)
		}
		return
	}

	// Determine the new notification's sender and receiver
	// The new notification should be from the user who made the comment to the other user
	var newSender, newReceiver int64
	newSender = userID // The person who made the comment
	if userID == originalNotification.Sender {
		newReceiver = originalNotification.Receiver
	} else {
		newReceiver = originalNotification.Sender
	}

	// Determine status for the new notification
	var newStatus string
	if originalNotification.Status == "accepted" {
		newStatus = "accepted"
	} else if originalNotification.Status == "rejected" {
		newStatus = "rejected"
	} else if originalNotification.Status == "pending" {
		// For pending notifications, create a comment notification that doesn't show Accept/Reject buttons
		// Use "comment" status to indicate this is just a comment, not an actionable notification
		newStatus = "comment"
	} else {
		newStatus = "pending" // If original status is NULL or any other value
	}

	// Create the new message with comment
	var newMessage string
	if request.Comment != "" {
		newMessage = request.Comment
	} else {
		newMessage = "Response sent"
	}

	// Generate new notification ID
	newNotificationID, err := utils.GenerateRandomID()
	if err != nil {
		fmt.Printf("Error generating notification ID: %v\n", err)
		http.Error(w, "Error generating notification ID", http.StatusInternalServerError)
		return
	}

	// Create notification with push notification (outside transaction since SendNotificationWithPush manages its own DB connection)
	err = SendNotificationWithPush(newSender, newReceiver, originalNotification.PID, originalNotification.FloorID, newMessage, newStatus, nil)
	if err != nil {
		fmt.Printf("Error creating notification: %v\n", err)
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	// Update the original notification with the comment
	_, err = tx.Exec(`
		UPDATE notification 
		SET comment = ?, updated_at = NOW(), updated_by = ?
		WHERE id = ?
	`, request.Comment, userID, request.NotificationID)

	if err != nil {
		fmt.Printf("Error updating original notification comment: %v\n", err)
		http.Error(w, "Failed to update original notification comment", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Created comment notification: ID=%d, Message='%s', from user=%d to user=%d", 
		newNotificationID, newMessage, newSender, newReceiver)

	// Send response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Comment sent successfully",
		"notification_id": newNotificationID,
	})
}

// GetConversationHistoryHandler handles GET requests to retrieve conversation history between users
func GetConversationHistoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Conversation History Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
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

	// Get the floor ID from query parameters
	floorIDStr := r.URL.Query().Get("floor_id")
	if floorIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "floor_id parameter is required",
		})
		return
	}

	floorID, err := strconv.ParseInt(floorIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid floor_id parameter",
		})
		return
	}

	// Get database connection
	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}

	// Query to get all notifications for the same floor between the same users
	rows, err := db.Query(`
		SELECT 
			n.id,
			n.message,
			n.status,
			n.created_at,
			n.comment,
			n.sender,
			n.receiver,
			p.id as property_id,
			p.name as property_name,
			f.id as floor_id,
			f.name as floor_name,
			u1.name as sender_name,
			u2.name as receiver_name
		FROM notification n
		JOIN property p ON n.pid = p.id
		JOIN floor f ON n.fid = f.id
		JOIN user u1 ON n.sender = u1.id
		JOIN user u2 ON n.receiver = u2.id
		WHERE n.fid = ? AND (
			(n.sender = ? AND n.receiver IN (
				SELECT DISTINCT 
					CASE 
						WHEN sender = ? THEN receiver 
						ELSE sender 
					END
				FROM notification 
				WHERE fid = ? AND (sender = ? OR receiver = ?)
			)) OR
			(n.receiver = ? AND n.sender IN (
				SELECT DISTINCT 
					CASE 
						WHEN sender = ? THEN receiver 
						ELSE sender 
					END
				FROM notification 
				WHERE fid = ? AND (sender = ? OR receiver = ?)
			))
		)
		ORDER BY n.created_at DESC
	`, floorID, userID, userID, floorID, userID, userID, userID, userID, floorID, userID, userID)

	if err != nil {
		fmt.Printf("Error querying conversation history: %v\n", err)
		http.Error(w, "Failed to get conversation history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var conversations []map[string]interface{}
	for rows.Next() {
		var conv struct {
			ID           int64
			Message      string
			Status       sql.NullString
			CreatedAt    string
			Comment      sql.NullString
			Sender       int64
			Receiver     int64
			PropertyID   int64
			PropertyName string
			FloorID      int64
			FloorName    string
			SenderName   string
			ReceiverName string
		}

		err := rows.Scan(
			&conv.ID,
			&conv.Message,
			&conv.Status,
			&conv.CreatedAt,
			&conv.Comment,
			&conv.Sender,
			&conv.Receiver,
			&conv.PropertyID,
			&conv.PropertyName,
			&conv.FloorID,
			&conv.FloorName,
			&conv.SenderName,
			&conv.ReceiverName,
		)

		if err != nil {
			fmt.Printf("Error scanning conversation row: %v\n", err)
			continue
		}

		conversation := map[string]interface{}{
			"id":           conv.ID,
			"message":      conv.Message,
			"status":       conv.Status.String,
			"created_at":   conv.CreatedAt,
			"comment":      conv.Comment.String,
			"sender":       conv.Sender,
			"receiver":     conv.Receiver,
			"sender_name":  conv.SenderName,
			"receiver_name": conv.ReceiverName,
			"property": map[string]interface{}{
				"id":   conv.PropertyID,
				"name": conv.PropertyName,
			},
			"floor": map[string]interface{}{
				"id":   conv.FloorID,
				"name": conv.FloorName,
			},
			"is_from_me": conv.Sender == userID,
		}

		conversations = append(conversations, conversation)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating conversation rows: %v\n", err)
		http.Error(w, "Failed to process conversation history", http.StatusInternalServerError)
		return
	}

	// Send response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Conversation history retrieved successfully",
		"conversations": conversations,
	})
}

// CreateAdvancePaymentRequestHandler handles POST requests to create an advance payment request
func CreateAdvancePaymentRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Create Advance Payment Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Method not allowed", 0})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	fmt.Printf("User ID from session: %d\n", userID)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "User not authenticated", 0})
		return
	}

	// Extract property ID and floor ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	fmt.Printf("URL parts: %v\n", cleanParts)

	if len(cleanParts) != 5 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid URL format", 0})
		return
	}

	propertyID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid property ID", 0})
		return
	}

	floorID, err := strconv.ParseInt(cleanParts[3], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid floor ID", 0})
		return
	}

	fmt.Printf("Property ID: %d, Floor ID: %d\n", propertyID, floorID)

	var req AdvancePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid request body", 0})
		return
	}

	fmt.Printf("Request body: %+v\n", req)

	// Validate request
	if req.AdvanceUID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid advance_uid", 0})
		return
	}

	if req.Money <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Money amount must be greater than 0", 0})
		return
	}



	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database connection error", 0})
		return
	}

	// Check if user is a manager of the property
	var isManager bool
	managerQuery := "SELECT EXISTS(SELECT 1 FROM takes_care_of WHERE uid = ? AND pid = ?)"
	fmt.Printf("Executing manager check query: %s with userID: %d, propertyID: %d\n", managerQuery, userID, propertyID)
	
	err = db.QueryRow(managerQuery, userID, propertyID).Scan(&isManager)
	if err != nil {
		fmt.Printf("Error checking manager status: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database error", 0})
		return
	}

	fmt.Printf("Is user %d manager of property %d: %v\n", userID, propertyID, isManager)

	if !isManager {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Only managers can create advance payment requests", 0})
		return
	}

	// Verify that the floor exists and belongs to the property
	var floorExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM floor WHERE id = ? AND pid = ?)", floorID, propertyID).Scan(&floorExists)
	if err != nil {
		fmt.Printf("Error checking floor existence: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database error", 0})
		return
	}

	if !floorExists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Floor not found", 0})
		return
	}

	// Verify that the advance_uid user exists
	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)", req.AdvanceUID).Scan(&userExists)
	if err != nil {
		fmt.Printf("Error checking user existence: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database error", 0})
		return
	}

	if !userExists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "User not found", 0})
		return
	}

	// Generate random ID for advance payment
	advanceID, err := utils.GenerateRandomID()
	if err != nil {
		fmt.Printf("Error generating advance payment ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Error generating advance payment ID", 0})
		return
	}

	// Insert advance payment record
	currentTime := time.Now().In(time.FixedZone("BDT", 6*60*60)).Format("2006-01-02")
	_, err = db.Exec(`
		INSERT INTO advance (
			id, advance_uid, money, fid, created_at, created_by, updated_at, updated_by, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		advanceID,
		req.AdvanceUID,
		req.Money,
		floorID,
		currentTime,
		userID,
		currentTime,
		userID,
		"pending",
	)

	if err != nil {
		fmt.Printf("Error creating advance payment record: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, fmt.Sprintf("Error creating advance payment record: %v", err), 0})
		return
	}

	// Create notification for the advance payment request
	{
		// Create notification message with payment amount
		message := fmt.Sprintf("Advance payment request: %d tk", req.Money)

		// Create notification with push notification
		err = SendNotificationWithPush(userID, req.AdvanceUID, propertyID, floorID, message, "pending", nil)
		if err != nil {
			fmt.Printf("Error creating notification: %v\n", err)
			// Don't fail the entire request if notification creation fails
			fmt.Printf("Failed to create notification for advance payment")
		} else {
			fmt.Printf("Successfully created notification with push for advance payment request")
		}
	}

	fmt.Printf("Successfully created advance payment record ID: %d for floor ID: %d and user ID: %d\n", advanceID, floorID, req.AdvanceUID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AdvancePaymentResponse{
		Success:   true,
		Message:   "Advance payment request created successfully",
		AdvanceID: advanceID,
	})
}

// CheckPendingAdvancePaymentHandler handles GET requests to check if there's a pending advance payment for a floor
func CheckPendingAdvancePaymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Check Pending Advance Payment Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{false, "Method not allowed", false})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	fmt.Printf("User ID from session: %d\n", userID)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{false, "User not authenticated", false})
		return
	}

	// Extract floor ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	fmt.Printf("URL parts: %v\n", cleanParts)

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{false, "Invalid URL format", false})
		return
	}

	floorID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{false, "Invalid floor ID", false})
		return
	}

	fmt.Printf("Floor ID: %d\n", floorID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{false, "Database connection error", false})
		return
	}

	// Check if there's a pending advance payment for this floor
	var hasPending bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM advance WHERE fid = ? AND status = 'pending')", floorID).Scan(&hasPending)
	if err != nil {
		fmt.Printf("Error checking pending advance payment: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{false, "Database error", false})
		return
	}

	fmt.Printf("Floor %d has pending advance payment: %v\n", floorID, hasPending)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AdvancePaymentCheckResponse{
		Success:    true,
		Message:    "Advance payment status checked successfully",
		HasPending: hasPending,
	})
}

// CancelAdvancePaymentHandler handles DELETE requests to cancel an advance payment request
func CancelAdvancePaymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Cancel Advance Payment Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Method not allowed", 0})
		return
	}

	// Get user ID from session
	userID := getUserIDFromContext(r)
	fmt.Printf("User ID from session: %d\n", userID)
	if userID == 0 {
		fmt.Println("No user ID found in session")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "User not authenticated", 0})
		return
	}

	// Extract floor ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	var cleanParts []string
	for _, part := range pathParts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	fmt.Printf("URL parts: %v\n", cleanParts)

	if len(cleanParts) != 3 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid URL format", 0})
		return
	}

	floorID, err := strconv.ParseInt(cleanParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Invalid floor ID", 0})
		return
	}

	fmt.Printf("Floor ID: %d\n", floorID)

	db, err := config.GetDBConnection()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database connection error", 0})
		return
	}

	// Get the property ID for this floor to check if user is manager
	var propertyID int64
	err = db.QueryRow("SELECT pid FROM floor WHERE id = ?", floorID).Scan(&propertyID)
	if err != nil {
		fmt.Printf("Error getting property ID for floor: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database error", 0})
		return
	}

	// Check if user is a manager of the property
	var isManager bool
	managerQuery := "SELECT EXISTS(SELECT 1 FROM takes_care_of WHERE uid = ? AND pid = ?)"
	fmt.Printf("Executing manager check query: %s with userID: %d, propertyID: %d\n", managerQuery, userID, propertyID)
	
	err = db.QueryRow(managerQuery, userID, propertyID).Scan(&isManager)
	if err != nil {
		fmt.Printf("Error checking manager status: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database error", 0})
		return
	}

	fmt.Printf("Is user %d manager of property %d: %v\n", userID, propertyID, isManager)

	if !isManager {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Only managers can cancel advance payment requests", 0})
		return
	}

	// Delete the pending advance payment record
	result, err := db.Exec("DELETE FROM advance WHERE fid = ? AND status = 'pending'", floorID)
	if err != nil {
		fmt.Printf("Error deleting advance payment record: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, fmt.Sprintf("Error deleting advance payment record: %v", err), 0})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "Database error", 0})
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(AdvancePaymentResponse{false, "No pending advance payment found for this floor", 0})
		return
	}

	fmt.Printf("Successfully deleted advance payment record for floor ID: %d\n", floorID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AdvancePaymentResponse{
		Success:   true,
		Message:   "Advance payment request cancelled successfully",
		AdvanceID: 0,
	})
}

func GetPaymentHistoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== New Get Payment History Request ===")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL)

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Method not allowed",
		})
		return
	}

	// Get user ID from context
	userID := getUserIDFromContext(r)
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	fmt.Printf("User authenticated: ID=%d\n", userID)

	// Parse URL to get floor ID
	urlParts := strings.Split(r.URL.Path, "/")
	fmt.Printf("URL parts: %v\n", urlParts)
	
	if len(urlParts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid URL format",
		})
		return
	}

	floorID, err := strconv.ParseInt(urlParts[2], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid floor ID",
		})
		return
	}

	fmt.Printf("Floor ID: %d\n", floorID)

	// Parse pagination parameters
	page := 1
	limit := 25
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	
	offset := (page - 1) * limit

	// Get database connection
	db, err := config.GetDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database connection error",
		})
		return
	}

	// Get the tenant ID for this floor first
	var tenantID sql.NullInt64
	err = db.QueryRow("SELECT tenant FROM floor WHERE id = ?", floorID).Scan(&tenantID)
	if err != nil {
		fmt.Printf("Error getting tenant for floor: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error getting tenant information",
		})
		return
	}

	if !tenantID.Valid {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PaymentHistoryResponse{
			Success:  true,
			Message:  "No tenant assigned to this floor",
			Payments: []PaymentHistory{},
		})
		return
	}

	// Get payment history for the floor with calculated fields and pagination
	query := `
		SELECT 
			p.id,
			p.rent as new_added_rent,
			p.recieved_money,
			p.full_payment,
			p.created_at,
			p.electricity_bill as new_added_electricity_bill,
			p.paid_bill as paid_electricity_bill,
			COALESCE((
				SELECT SUM(p2.rent - p2.recieved_money)
				FROM payment p2
				WHERE p2.fid = p.fid AND p2.uid = p.uid AND p2.created_at < p.created_at
			), 0) as rent,
			COALESCE((
				SELECT SUM(p3.rent - p3.recieved_money)
				FROM payment p3
				WHERE p3.fid = p.fid AND p3.uid = p.uid AND p3.created_at <= p.created_at
			), 0) as due_rent,
			COALESCE((
				SELECT SUM(COALESCE(p4.electricity_bill, 0) - COALESCE(p4.paid_bill, 0))
				FROM payment p4
				WHERE p4.fid = p.fid AND p4.uid = p.uid AND p4.created_at <= p.created_at
			), 0) as due_electricity_bill,
			COALESCE((
				SELECT SUM(COALESCE(p5.electricity_bill, 0) - COALESCE(p5.paid_bill, 0))
				FROM payment p5
				WHERE p5.fid = p.fid AND p5.uid = p.uid AND p5.created_at < p.created_at
			), 0) as electricity_bill
		FROM payment p
		WHERE p.fid = ? AND p.uid = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`
	fmt.Printf("Executing query: %s with floorID: %d, tenantID: %d, limit: %d, offset: %d\n", query, floorID, tenantID.Int64, limit, offset)
	
	rows, err := db.Query(query, floorID, tenantID.Int64, limit, offset)
	if err != nil {
		fmt.Printf("Error querying payment history: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error querying payment history",
		})
		return
	}
	defer rows.Close()

	var payments []PaymentHistory
	paymentCount := 0
	for rows.Next() {
		var payment PaymentHistory
		var createdAt time.Time
		var newAddedRent, rent, dueRent, dueElectricityBill, electricityBill float64
		var newAddedElectricityBill, paidElectricityBill sql.NullFloat64
		err := rows.Scan(&payment.ID, &newAddedRent, &payment.ReceivedMoney, &payment.FullPayment, &createdAt, &newAddedElectricityBill, &paidElectricityBill, &rent, &dueRent, &dueElectricityBill, &electricityBill)
		if err != nil {
			fmt.Printf("Error scanning payment: %v\n", err)
			continue
		}
		
		// Set new_added_rent (the actual rent amount from payment table)
		payment.NewAddedRent = newAddedRent
		
		// Set rent (sum of previous outstanding amounts)
		payment.Rent = rent
		
		// Set due_rent (sum after this row was inserted)
		payment.DueRent = dueRent
		
		// Set electricity bill fields
		if newAddedElectricityBill.Valid {
			payment.NewAddedElectricityBill = &newAddedElectricityBill.Float64
		}
		if paidElectricityBill.Valid {
			payment.PaidElectricityBill = &paidElectricityBill.Float64
		}
		payment.DueElectricityBill = &dueElectricityBill
		payment.ElectricityBill = &electricityBill
		
		payment.CreatedAt = createdAt.Format("2006-01-02T15:04:05Z")
		payments = append(payments, payment)
		paymentCount++
		
		fmt.Printf("Found payment: ID=%d, NewAddedRent=%.2f, Rent(Previous Outstanding)=%.2f, ReceivedMoney=%.2f, DueRent=%.2f, ElectricityBill=%.2f, FullPayment=%v\n", 
			payment.ID, payment.NewAddedRent, payment.Rent, payment.ReceivedMoney, payment.DueRent, *payment.ElectricityBill, payment.FullPayment)
	}

	fmt.Printf("Total payments found: %d\n", paymentCount)

	// Get total count for pagination
	var totalCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM payment p
		WHERE p.fid = ? AND p.uid = ?
	`, floorID, tenantID.Int64).Scan(&totalCount)
	
	if err != nil {
		fmt.Printf("Error getting total count: %v\n", err)
		totalCount = paymentCount
	}

	// Calculate pagination metadata
	totalPages := (totalCount + limit - 1) / limit
	hasNextPage := page < totalPages
	hasPrevPage := page > 1

	// Return success response
	response := PaymentHistoryResponse{
		Success:  true,
		Message:  "Payment history retrieved successfully",
		Payments: payments,
		Pagination: PaginationInfo{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalCount:   totalCount,
			Limit:        limit,
			HasNextPage:  hasNextPage,
			HasPrevPage:  hasPrevPage,
		},
	}
	
	fmt.Printf("Sending response: %+v\n", response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}