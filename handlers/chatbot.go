package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-rent/config"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Tenant represents a tenant in the chatbot system
type ChatbotTenant struct {
	ID                  string  `json:"id"`
	RiskProbability     float64 `json:"risk_probability"`
	RiskLevel           string  `json:"risk_level"`
	PreviousLateCount   int     `json:"previous_late_count"`
	AvgDelayDays        float64 `json:"avg_delay_days"`
	LastPaymentDate     string  `json:"last_payment_date"`
	CurrentRentAmount   float64 `json:"current_rent_amount"`
	RentToIncomeRatio   float64 `json:"rent_to_income_ratio"`
	TenancyMonths       int     `json:"tenancy_months"`
	EmploymentStatus    string  `json:"employment_status"`
	CreditScore         int     `json:"credit_score"`
	ComplaintCount      int     `json:"complaint_count"`
	PropertyDamageCount int     `json:"property_damage_count"`
	Notes               string  `json:"notes"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	// Enhanced risk factors
	PartialPaymentRatio float64 `json:"partial_payment_ratio,omitempty"`
	PaymentTrend        float64 `json:"payment_trend,omitempty"`
}

// ChatRequest represents a chat message from the user
type ChatRequest struct {
	Message  string `json:"message"`
	TenantID string `json:"tenant_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// ChatResponse represents the bot's response
type ChatResponse struct {
	Intent             string   `json:"intent"`
	Confidence         float64  `json:"confidence"`
	ResponseText       string   `json:"response_text"`
	SuggestedFollowups []string `json:"suggested_followups,omitempty"`
	ProcessingTimeMS   int64    `json:"processing_time_ms"`
	Timestamp          string   `json:"timestamp"`
	Data               any      `json:"data,omitempty"`
}

// IntentDetector handles intent detection
type IntentDetector struct {
	patterns map[string][]*regexp.Regexp
}

// NewIntentDetector creates a new intent detector
func NewIntentDetector() *IntentDetector {
	id := &IntentDetector{
		patterns: make(map[string][]*regexp.Regexp),
	}

	// Define intent patterns
	id.patterns["EXPLAIN_RISK"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)why.*risk`),
		regexp.MustCompile(`(?i)explain.*risk`),
		regexp.MustCompile(`(?i)what.*makes.*risk`),
		regexp.MustCompile(`(?i)risk.*factors`),
		regexp.MustCompile(`(?i)why.*high.*risk`),
		regexp.MustCompile(`(?i)why.*critical.*risk`),
	}

	id.patterns["RECOMMEND_ACTION"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)what.*should.*do`),
		regexp.MustCompile(`(?i)recommend.*action`),
		regexp.MustCompile(`(?i)how.*handle`),
		regexp.MustCompile(`(?i)next.*step`),
		regexp.MustCompile(`(?i)suggest.*action`),
	}

	id.patterns["LIST_HIGH_RISK"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)list.*high.*risk`),
		regexp.MustCompile(`(?i)who.*high.*risk`),
		regexp.MustCompile(`(?i)show.*risky`),
		regexp.MustCompile(`(?i)high.*risk.*tenant`),
		regexp.MustCompile(`(?i)list.*critical.*risk`),
		regexp.MustCompile(`(?i)critical.*tenant`),
	}

	id.patterns["MONTHLY_SUMMARY"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)monthly.*summary`),
		regexp.MustCompile(`(?i)this.*month.*report`),
		regexp.MustCompile(`(?i)risk.*summary`),
		regexp.MustCompile(`(?i)overview.*this.*month`),
	}

	id.patterns["COMPARE_TENANTS"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)compare.*tenant`),
		regexp.MustCompile(`(?i)vs.*tenant`),
		regexp.MustCompile(`(?i)compare.*T\d+.*and.*T\d+`),
		regexp.MustCompile(`(?i)difference.*between`),
	}

	id.patterns["PAYMENT_HISTORY"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)payment.*history`),
		regexp.MustCompile(`(?i)late.*payment`),
		regexp.MustCompile(`(?i)when.*last.*pay`),
		regexp.MustCompile(`(?i)past.*payment`),
	}

	id.patterns["LEASE_RENEWAL"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)lease.*renew`),
		regexp.MustCompile(`(?i)renew.*lease`),
		regexp.MustCompile(`(?i)should.*i.*renew`),
		regexp.MustCompile(`(?i)extend.*lease`),
	}

	return id
}

// Detect intent from message
func (id *IntentDetector) Detect(message string) (string, float64) {
	message = strings.ToLower(strings.TrimSpace(message))

	for intent, patterns := range id.patterns {
		for _, pattern := range patterns {
			if pattern.MatchString(message) {
				confidence := 0.9
				if strings.Contains(message, "?") {
					confidence = 0.95
				}
				return intent, confidence
			}
		}
	}

	// Keyword fallback
	keywords := map[string]string{
		"risk":      "EXPLAIN_RISK",
		"late":      "EXPLAIN_RISK",
		"recommend": "RECOMMEND_ACTION",
		"list":      "LIST_HIGH_RISK",
		"summary":   "MONTHLY_SUMMARY",
		"compare":   "COMPARE_TENANTS",
		"payment":   "PAYMENT_HISTORY",
		"lease":     "LEASE_RENEWAL",
	}

	for keyword, intent := range keywords {
		if strings.Contains(message, keyword) {
			return intent, 0.7
		}
	}

	// Check for tenant ID pattern
	if matched, _ := regexp.MatchString(`(?i)t\d+`, message); matched {
		if strings.Contains(message, "risk") {
			return "EXPLAIN_RISK", 0.8
		}
	}

	return "UNKNOWN", 0.0
}

// Extract tenant IDs from message
func (id *IntentDetector) ExtractTenantIDs(message string) []string {
	re := regexp.MustCompile(`(?i)t\d+`)
	return re.FindAllString(strings.ToUpper(message), -1)
}

// Extract phone numbers from message
func (id *IntentDetector) ExtractPhoneNumbers(message string) []string {
	// Pattern 1: Match 11-digit numbers starting with 0 (01794002263 format)
	re1 := regexp.MustCompile(`0?1[3-9]\d{9}`)
	matches1 := re1.FindAllString(message, -1)
	
	// Pattern 2: Match 10-digit numbers starting with 1 (1794002263 format, no leading 0)
	re2 := regexp.MustCompile(`1[3-9]\d{8}`)
	matches2 := re2.FindAllString(message, -1)
	
	// Pattern 3: Match with country code (+8801712345678 or 8801712345678)
	re3 := regexp.MustCompile(`(?:\+?880|880)1[3-9]\d{8}`)
	matches3 := re3.FindAllString(message, -1)
	
	// Combine all matches
	allMatches := append(matches1, matches2...)
	allMatches = append(allMatches, matches3...)
	
	log.Printf("ExtractPhoneNumbers: raw matches1=%v, matches2=%v, matches3=%v", matches1, matches2, matches3)
	
	// Deduplicate and normalize
	unique := make(map[string]bool)
	result := []string{}
	for _, match := range allMatches {
		// Normalize phone number (remove spaces, dashes)
		normalized := strings.ReplaceAll(match, " ", "")
		normalized = strings.ReplaceAll(normalized, "-", "")
		
		// Remove country code prefixes
		normalized = strings.TrimPrefix(normalized, "+880")
		normalized = strings.TrimPrefix(normalized, "880")
		
		// Remove leading 0 if present (Bangladeshi format: 01794002263 -> 1794002263)
		if strings.HasPrefix(normalized, "0") {
			normalized = normalized[1:]
		}
		
		// Must be exactly 10 digits and start with 1 (Bangladeshi mobile numbers)
		if len(normalized) == 10 && strings.HasPrefix(normalized, "1") {
			// Check if second digit is 3-9 (valid Bangladeshi mobile prefix)
			if len(normalized) >= 2 {
				secondDigit := normalized[1]
				if secondDigit >= '3' && secondDigit <= '9' {
					if !unique[normalized] {
						unique[normalized] = true
						result = append(result, normalized)
					}
				}
			}
		}
	}
	
	log.Printf("ExtractPhoneNumbers: message=%s, found=%v", message, result)
	return result
}

// ResponseGenerator generates responses based on intent
type ResponseGenerator struct {
	tenants map[string]*ChatbotTenant
	mu      sync.RWMutex
}

// NewResponseGenerator creates a new response generator
func NewResponseGenerator() *ResponseGenerator {
	rg := &ResponseGenerator{
		tenants: make(map[string]*ChatbotTenant),
	}
	rg.initializeSampleData()
	return rg
}

// Initialize sample tenant data
func (rg *ResponseGenerator) initializeSampleData() {
	rg.mu.Lock()
	defer rg.mu.Unlock()

	// Create sample tenants with 4 risk bands
	sampleTenants := []*ChatbotTenant{
		{
			ID:                  "01712345679",
			RiskProbability:     0.78,
			RiskLevel:           "High",
			PreviousLateCount:   3,
			AvgDelayDays:        6.2,
			LastPaymentDate:     time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
			CurrentRentAmount:   1200,
			RentToIncomeRatio:   0.42,
			TenancyMonths:       12,
			EmploymentStatus:    "Full-time",
			CreditScore:         580,
			ComplaintCount:      1,
			PropertyDamageCount: 0,
			PartialPaymentRatio: 0.10,
			PaymentTrend:        2.0,
			Notes:               "3 late payments in last 6 months",
			CreatedAt:           time.Now().Format(time.RFC3339),
			UpdatedAt:           time.Now().Format(time.RFC3339),
		},
		{
			ID:                  "01712345678",
			RiskProbability:     0.50,
			RiskLevel:           "Medium",
			PreviousLateCount:   1,
			AvgDelayDays:        2.0,
			LastPaymentDate:     time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			CurrentRentAmount:   850,
			RentToIncomeRatio:   0.28,
			TenancyMonths:       24,
			EmploymentStatus:    "Full-time",
			CreditScore:         720,
			ComplaintCount:      0,
			PropertyDamageCount: 0,
			PartialPaymentRatio: 0.05,
			PaymentTrend:        0.0,
			Notes:               "1 late payment in last 12 months",
			CreatedAt:           time.Now().Format(time.RFC3339),
			UpdatedAt:           time.Now().Format(time.RFC3339),
		},
		{
			ID:                  "01987654321",
			RiskProbability:     0.92,
			RiskLevel:           "Critical",
			PreviousLateCount:   5,
			AvgDelayDays:        12.5,
			LastPaymentDate:     time.Now().AddDate(0, -2, 0).Format("2006-01-02"),
			CurrentRentAmount:   1500,
			RentToIncomeRatio:   0.55,
			TenancyMonths:       6,
			EmploymentStatus:    "Part-time",
			CreditScore:         520,
			ComplaintCount:      3,
			PropertyDamageCount: 2,
			PartialPaymentRatio: 0.30,
			PaymentTrend:        3.0,
			Notes:               "Multiple late payments and complaints",
			CreatedAt:           time.Now().Format(time.RFC3339),
			UpdatedAt:           time.Now().Format(time.RFC3339),
		},
		{
			ID:                  "01712345675",
			RiskProbability:     0.20,
			RiskLevel:           "Low",
			PreviousLateCount:   0,
			AvgDelayDays:        0.0,
			LastPaymentDate:     time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			CurrentRentAmount:   950,
			RentToIncomeRatio:   0.25,
			TenancyMonths:       36,
			EmploymentStatus:    "Full-time",
			CreditScore:         780,
			ComplaintCount:      0,
			PropertyDamageCount: 0,
			PartialPaymentRatio: 0.0,
			PaymentTrend:        0.0,
			Notes:               "Excellent payment history",
			CreatedAt:           time.Now().Format(time.RFC3339),
			UpdatedAt:           time.Now().Format(time.RFC3339),
		},
	}

	for _, tenant := range sampleTenants {
		rg.tenants[tenant.ID] = tenant
	}
}

// Get tenant by ID
func (rg *ResponseGenerator) GetTenant(id string) (*ChatbotTenant, bool) {
	rg.mu.RLock()
	defer rg.mu.RUnlock()

	tenant, exists := rg.tenants[strings.ToUpper(id)]
	if !exists {
		// Generate synthetic tenant
		tenant = rg.generateSyntheticTenant(id)
		rg.tenants[tenant.ID] = tenant
	}

	return tenant, exists
}

// Get tenant by phone number from database - ENHANCED VERSION
func (rg *ResponseGenerator) GetTenantByPhone(phoneNumber string) (*ChatbotTenant, error) {
	db, err := config.GetDBConnection()
	if err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}

	// Normalize phone number
	normalizedPhone := strings.ReplaceAll(phoneNumber, " ", "")
	normalizedPhone = strings.ReplaceAll(normalizedPhone, "-", "")
	normalizedPhone = strings.TrimPrefix(normalizedPhone, "+880")
	normalizedPhone = strings.TrimPrefix(normalizedPhone, "880")
	if strings.HasPrefix(normalizedPhone, "0") && len(normalizedPhone) == 11 {
		normalizedPhone = normalizedPhone[1:]
	}

	// Get user ID by phone number
	var userID int64
	var userName sql.NullString
	err = db.QueryRow(`
		SELECT id, name 
		FROM user 
		WHERE phone_number = ?`, normalizedPhone).Scan(&userID, &userName)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found with phone number: %s", phoneNumber)
		}
		return nil, fmt.Errorf("error querying user: %v", err)
	}

	// Get tenant's floor and payment data
	var floorID sql.NullInt64
	var rent sql.NullInt64
	var floorCreatedAt sql.NullString
	err = db.QueryRow(`
		SELECT f.id, f.rent, f.created_at
		FROM floor f
		WHERE f.tenant = ?
		ORDER BY f.created_at DESC
		LIMIT 1`, userID).Scan(&floorID, &rent, &floorCreatedAt)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error querying floor: %v", err)
	}

	// Initialize variables
	var previousLateCount int
	var avgDelayDays float64
	var lastPaymentDate sql.NullString
	var currentRentAmount float64
	var tenancyMonths int
	var partialPaymentRatio float64
	var paymentTrend float64

	if floorID.Valid {
		// Get current rent
		if rent.Valid {
			currentRentAmount = float64(rent.Int64)
		}

		// Calculate tenancy months
		if floorCreatedAt.Valid {
			createdTime, err := time.Parse("2006-01-02 15:04:05", floorCreatedAt.String)
			if err == nil {
				tenancyMonths = int(time.Since(createdTime).Hours() / 24 / 30)
			}
		}

		// ✅ ENHANCED: Try comprehensive payment analysis first
		var totalPayments sql.NullInt64
		var lateCount sql.NullInt64
		var avgDelay sql.NullFloat64
		var lastPayment sql.NullString
		var partialPayments sql.NullInt64
		
		err = db.QueryRow(`
			SELECT 
				COUNT(*) as total_payments,
				SUM(CASE WHEN recieved_money < rent THEN 1 ELSE 0 END) as late_count,
				COALESCE(AVG(CASE WHEN recieved_money < rent THEN DATEDIFF(created_at, DATE_SUB(created_at, INTERVAL DAY(created_at)-1 DAY)) ELSE 0 END), 0) as avg_delay,
				MAX(created_at) as last_payment,
				SUM(CASE WHEN recieved_money > 0 AND recieved_money < rent THEN 1 ELSE 0 END) as partial_payments
			FROM payment
			WHERE fid = ? AND uid = ?
		`, floorID.Int64, userID).Scan(&totalPayments, &lateCount, &avgDelay, &lastPayment, &partialPayments)
		
		if err == nil {
			// Successfully got enhanced metrics
			if lateCount.Valid {
				previousLateCount = int(lateCount.Int64)
			}
			if avgDelay.Valid {
				avgDelayDays = avgDelay.Float64
			}
			if lastPayment.Valid {
				lastPaymentDate = lastPayment
			}
			
			// Calculate partial payment ratio
			if totalPayments.Valid && totalPayments.Int64 > 0 && partialPayments.Valid {
				partialPaymentRatio = float64(partialPayments.Int64) / float64(totalPayments.Int64)
			}

			// ✅ ENHANCED: Calculate payment trend
			var recentLate sql.NullInt64
			var olderLate sql.NullInt64
			
			err = db.QueryRow(`
				SELECT 
					SUM(CASE WHEN recieved_money < rent AND created_at >= DATE_SUB(NOW(), INTERVAL 3 MONTH) THEN 1 ELSE 0 END) as recent_late,
					SUM(CASE WHEN recieved_money < rent AND created_at < DATE_SUB(NOW(), INTERVAL 3 MONTH) AND created_at >= DATE_SUB(NOW(), INTERVAL 6 MONTH) THEN 1 ELSE 0 END) as older_late
				FROM payment
				WHERE fid = ? AND uid = ?
			`, floorID.Int64, userID).Scan(&recentLate, &olderLate)
			
			if err == nil {
				// Payment trend: positive = getting worse, negative = improving
				if recentLate.Valid && olderLate.Valid {
					if olderLate.Int64 > 0 {
						paymentTrend = float64(recentLate.Int64) - float64(olderLate.Int64)
					} else if recentLate.Int64 > 0 {
						paymentTrend = float64(recentLate.Int64)
					}
				}
			} else {
				log.Printf("Warning: Could not calculate payment trend: %v (using fallback)", err)
			}

		} else if err != sql.ErrNoRows {
			// ⚠️ FALLBACK: Use simple metrics if enhanced query fails
			log.Printf("Warning: Enhanced metrics query failed, using simple fallback: %v", err)
			
			err = db.QueryRow(`
				SELECT 
					COUNT(*) as late_count,
					COALESCE(AVG(CASE WHEN recieved_money < rent THEN DATEDIFF(created_at, DATE_SUB(created_at, INTERVAL DAY(created_at)-1 DAY)) ELSE 0 END), 0) as avg_delay,
					MAX(created_at) as last_payment
				FROM payment
				WHERE fid = ? AND uid = ? AND recieved_money < rent
			`, floorID.Int64, userID).Scan(&lateCount, &avgDelay, &lastPayment)
			
			if err == nil {
				if lateCount.Valid {
					previousLateCount = int(lateCount.Int64)
				}
				if avgDelay.Valid {
					avgDelayDays = avgDelay.Float64
				}
				if lastPayment.Valid {
					lastPaymentDate = lastPayment
				}
			} else if err != sql.ErrNoRows {
				log.Printf("Error calculating even simple risk metrics: %v", err)
			}
		}
	}

	// ✅ Calculate risk using enhanced formula (with fallback support)
	riskProbability := rg.calculateEnhancedRiskProbability(
		previousLateCount,
		avgDelayDays,
		tenancyMonths,
		partialPaymentRatio,
		paymentTrend,
		currentRentAmount,
	)
	
	// Determine risk level - 4 bands
	riskLevel := "Low"
	if riskProbability >= 0.85 {
		riskLevel = "Critical"
	} else if riskProbability >= 0.65 {
		riskLevel = "High"
	} else if riskProbability >= 0.35 {
		riskLevel = "Medium"
	}

	tenantName := normalizedPhone
	if userName.Valid && userName.String != "" {
		tenantName = userName.String
	}

	tenant := &ChatbotTenant{
		ID:                normalizedPhone,
		RiskProbability:   riskProbability,
		RiskLevel:         riskLevel,
		PreviousLateCount: previousLateCount,
		AvgDelayDays:      avgDelayDays,
		LastPaymentDate: func() string {
			if lastPaymentDate.Valid {
				return lastPaymentDate.String
			}
			return "No payments yet"
		}(),
		CurrentRentAmount:   currentRentAmount,
		RentToIncomeRatio:   0.3,
		TenancyMonths:       tenancyMonths,
		EmploymentStatus:    "Unknown",
		CreditScore:         650,
		ComplaintCount:      0,
		PropertyDamageCount: 0,
		PartialPaymentRatio: partialPaymentRatio,
		PaymentTrend:        paymentTrend,
		Notes:               fmt.Sprintf("Tenant: %s", tenantName),
		CreatedAt:           time.Now().Format(time.RFC3339),
		UpdatedAt:           time.Now().Format(time.RFC3339),
	}

	// Cache the tenant
	rg.mu.Lock()
	rg.tenants[normalizedPhone] = tenant
	rg.mu.Unlock()

	return tenant, nil
}

// ✅ ENHANCED RISK CALCULATION - Matches Figure 2
func (rg *ResponseGenerator) calculateEnhancedRiskProbability(
	lateCount int, 
	avgDelay float64, 
	tenancyMonths int,
	partialPaymentRatio float64,
	paymentTrend float64,
	rentAmount float64,
) float64 {
	risk := 0.0

	// 1. Late payment count factor (0-0.10) - 10% weight
	if lateCount > 0 {
		lateRisk := float64(lateCount) * 0.02
		if lateRisk > 0.10 {
			lateRisk = 0.10
		}
		risk += lateRisk
	}

	// 2. Average delay factor (0-0.25) - 25% weight
	if avgDelay > 0 {
		delayRisk := (avgDelay / 30.0) * 0.25
		if delayRisk > 0.25 {
			delayRisk = 0.25
		}
		risk += delayRisk
	}

	// 3. Payment severity/trend factor (0-0.15) - 15% weight
	if paymentTrend > 0 {
		// Positive trend = getting worse
		trendRisk := paymentTrend * 0.05
		if trendRisk > 0.15 {
			trendRisk = 0.15
		}
		risk += trendRisk
	} else if paymentTrend < 0 {
		// Negative trend = improving
		trendBonus := paymentTrend * 0.03
		if trendBonus < -0.10 {
			trendBonus = -0.10
		}
		risk += trendBonus
	}

	// 4. Partial payment ratio factor (0-0.15) - 15% weight
	if partialPaymentRatio > 0 {
		partialRisk := partialPaymentRatio * 0.15
		if partialRisk > 0.15 {
			partialRisk = 0.15
		}
		risk += partialRisk
	}

	// 5. Rent-to-income ratio proxy (0-0.10) - 10% weight
	if rentAmount > 1000 {
		incomeRisk := 0.05
		if rentAmount > 1500 {
			incomeRisk = 0.10
		}
		risk += incomeRisk
	}

	// 6. Tenancy duration bonus (0 to -0.15) - 15% reduction
	if tenancyMonths > 12 {
		tenancyBonus := float64(tenancyMonths-12) / 60.0 * 0.15
		if tenancyBonus > 0.15 {
			tenancyBonus = 0.15
		}
		risk -= tenancyBonus
	}

	// Clamp to [0, 1]
	if risk < 0 {
		risk = 0
	}
	if risk > 1 {
		risk = 1
	}

	return risk
}

// ⚠️ FALLBACK: Simple risk calculation for backward compatibility
func (rg *ResponseGenerator) calculateRiskProbability(lateCount int, avgDelay float64, tenancyMonths int) float64 {
	// Use enhanced calculation with default values for missing metrics
	return rg.calculateEnhancedRiskProbability(lateCount, avgDelay, tenancyMonths, 0, 0, 0)
}

// Generate synthetic tenant for unknown IDs
func (rg *ResponseGenerator) generateSyntheticTenant(id string) *ChatbotTenant {
	rand.Seed(time.Now().UnixNano())

	riskProb := rand.Float64()
	
	// 4 risk bands
	riskLevel := "Low"
	if riskProb >= 0.85 {
		riskLevel = "Critical"
	} else if riskProb >= 0.65 {
		riskLevel = "High"
	} else if riskProb >= 0.35 {
		riskLevel = "Medium"
	}

	return &ChatbotTenant{
		ID:                  strings.ToUpper(id),
		RiskProbability:     riskProb,
		RiskLevel:           riskLevel,
		PreviousLateCount:   rand.Intn(5),
		AvgDelayDays:        rand.Float64() * 10,
		LastPaymentDate:     time.Now().AddDate(0, 0, -rand.Intn(30)).Format("2006-01-02"),
		CurrentRentAmount:   500 + rand.Float64()*2000,
		RentToIncomeRatio:   0.2 + rand.Float64()*0.5,
		TenancyMonths:       rand.Intn(60),
		EmploymentStatus:    []string{"Full-time", "Part-time", "Self-employed", "Unemployed"}[rand.Intn(4)],
		CreditScore:         300 + rand.Intn(550),
		ComplaintCount:      rand.Intn(3),
		PropertyDamageCount: rand.Intn(2),
		PartialPaymentRatio: rand.Float64() * 0.3,
		PaymentTrend:        (rand.Float64() - 0.5) * 4,
		Notes:               "Synthetic tenant data",
		CreatedAt:           time.Now().Format(time.RFC3339),
		UpdatedAt:           time.Now().Format(time.RFC3339),
	}
}

// Get all tenants
func (rg *ResponseGenerator) GetAllTenants() []*ChatbotTenant {
	rg.mu.RLock()
	defer rg.mu.RUnlock()

	tenants := make([]*ChatbotTenant, 0, len(rg.tenants))
	for _, tenant := range rg.tenants {
		tenants = append(tenants, tenant)
	}

	sort.Slice(tenants, func(i, j int) bool {
		return tenants[i].RiskProbability > tenants[j].RiskProbability
	})

	return tenants
}

// FormatPhoneNumberForDisplay adds leading zero to 10-digit phone numbers
func FormatPhoneNumberForDisplay(phoneNumber string) string {
	if len(phoneNumber) == 10 && strings.HasPrefix(phoneNumber, "1") {
		return "0" + phoneNumber
	}
	if len(phoneNumber) == 11 && strings.HasPrefix(phoneNumber, "0") {
		return phoneNumber
	}
	return phoneNumber
}

// Explain risk response (using tenant object)
func (rg *ResponseGenerator) ExplainRiskForTenant(tenant *ChatbotTenant) string {

	riskFactors := []string{}

	if tenant.PreviousLateCount >= 3 {
		riskFactors = append(riskFactors, fmt.Sprintf("%d previous late payments", tenant.PreviousLateCount))
	} else if tenant.PreviousLateCount > 0 {
		riskFactors = append(riskFactors, fmt.Sprintf("%d late payment(s)", tenant.PreviousLateCount))
	}

	if tenant.AvgDelayDays > 5 {
		riskFactors = append(riskFactors, fmt.Sprintf("average delay of %.1f days", tenant.AvgDelayDays))
	} else if tenant.AvgDelayDays > 2 {
		riskFactors = append(riskFactors, fmt.Sprintf("moderate payment delays (%.1f days)", tenant.AvgDelayDays))
	}

	// ✅ ENHANCED: Include partial payment ratio
	if tenant.PartialPaymentRatio > 0.2 {
		riskFactors = append(riskFactors, fmt.Sprintf("high partial payment rate (%.0f%%)", tenant.PartialPaymentRatio*100))
	} else if tenant.PartialPaymentRatio > 0.1 {
		riskFactors = append(riskFactors, fmt.Sprintf("some partial payments (%.0f%%)", tenant.PartialPaymentRatio*100))
	}

	// ✅ ENHANCED: Include payment trend
	if tenant.PaymentTrend > 1 {
		riskFactors = append(riskFactors, "worsening payment pattern")
	} else if tenant.PaymentTrend < -1 {
		riskFactors = append(riskFactors, "improving payment pattern")
	}

	if tenant.RentToIncomeRatio > 0.4 {
		riskFactors = append(riskFactors, fmt.Sprintf("high rent burden (%.2f ratio)", tenant.RentToIncomeRatio))
	} else if tenant.RentToIncomeRatio > 0.3 {
		riskFactors = append(riskFactors, fmt.Sprintf("moderate rent burden (%.2f ratio)", tenant.RentToIncomeRatio))
	}

	if tenant.CreditScore < 600 {
		riskFactors = append(riskFactors, fmt.Sprintf("low credit score (%d)", tenant.CreditScore))
	}

	if tenant.EmploymentStatus == "Unemployed" || tenant.EmploymentStatus == "Part-time" {
		riskFactors = append(riskFactors, "unstable employment")
	}

	if tenant.ComplaintCount > 0 {
		riskFactors = append(riskFactors, fmt.Sprintf("%d complaint(s)", tenant.ComplaintCount))
	}

	displayID := FormatPhoneNumberForDisplay(tenant.ID)
	riskText := fmt.Sprintf("Tenant %s is classified as **%s risk** (score: %.2f).",
		displayID, tenant.RiskLevel, tenant.RiskProbability)

	if len(riskFactors) > 0 {
		riskText += " Key factors: " + strings.Join(riskFactors, ", ") + "."
	} else {
		riskText += " No significant risk factors identified."
	}
	
	return riskText
}

// Explain risk response (backward compatibility)
func (rg *ResponseGenerator) ExplainRisk(tenantID string) string {
	tenant, _ := rg.GetTenant(tenantID)
	return rg.ExplainRiskForTenant(tenant)
}

// Recommend action response (using tenant object)
func (rg *ResponseGenerator) RecommendActionForTenant(tenant *ChatbotTenant) string {

	actions := []string{}

	// 4 risk bands with specific interventions
	if tenant.RiskProbability >= 0.85 {
		// Critical Risk (≥ 0.85)
		actions = append(actions,
			"⚠️ URGENT: Immediate contact required",
			"Consider legal consultation",
			"Document all communications",
			"Review lease termination options",
			"Offer last-chance payment plan with strict conditions",
			"Flag for daily monitoring",
		)
	} else if tenant.RiskProbability >= 0.65 {
		// High Risk (0.65 - 0.85)
		actions = append(actions,
			"Send early reminder 10 days before due date",
			"Offer flexible payment plan options",
			"Schedule in-person meeting",
			"Flag for weekly monitoring",
			"Consider additional security deposit on renewal",
		)
	} else if tenant.RiskProbability >= 0.35 {
		// Medium Risk (0.35 - 0.65)
		actions = append(actions,
			"Monitor payment closely this month",
			"Send reminder if late by 1 day",
			"Consider gentle follow-up call",
			"Review payment history trends",
		)
	} else {
		// Low Risk (< 0.35)
		actions = append(actions,
			"Standard reminder 3 days before due date",
			"Regular monthly monitoring",
			"Maintain positive relationship",
		)
	}

	displayID := FormatPhoneNumberForDisplay(tenant.ID)
	actionText := fmt.Sprintf("**Recommended Actions for %s (%s Risk)**\n\n", displayID, tenant.RiskLevel)
	for i, action := range actions {
		actionText += fmt.Sprintf("%d. %s\n", i+1, action)
	}

	return actionText
}

// Recommend action response (backward compatibility)
func (rg *ResponseGenerator) RecommendAction(tenantID string) string {
	tenant, _ := rg.GetTenant(tenantID)
	return rg.RecommendActionForTenant(tenant)
}

// Monthly summary
func (rg *ResponseGenerator) MonthlySummary() string {
	allTenants := rg.GetAllTenants()

	// Count all 4 risk bands
	criticalRisk := 0
	highRisk := 0
	mediumRisk := 0
	lowRisk := 0
	totalRent := 0.0
	totalAtRisk := 0.0

	for _, tenant := range allTenants {
		totalRent += tenant.CurrentRentAmount

		if tenant.RiskProbability >= 0.85 {
			criticalRisk++
			totalAtRisk += tenant.CurrentRentAmount
		} else if tenant.RiskProbability >= 0.65 {
			highRisk++
			totalAtRisk += tenant.CurrentRentAmount
		} else if tenant.RiskProbability >= 0.35 {
			mediumRisk++
		} else {
			lowRisk++
		}
	}

	summary := fmt.Sprintf("**Monthly Risk Summary - %s**\n\n", time.Now().Format("January 2006"))
	summary += fmt.Sprintf("**Overview**:\n")
	summary += fmt.Sprintf("- Total Tenants: %d\n", len(allTenants))
	summary += fmt.Sprintf("- Total Monthly Rent: $%.2f\n", totalRent)
	summary += fmt.Sprintf("- Rent at Risk (Critical + High): $%.2f\n\n", totalAtRisk)

	summary += fmt.Sprintf("**Risk Distribution**:\n")
	summary += fmt.Sprintf("- 🔴 Critical Risk: %d tenants\n", criticalRisk)
	summary += fmt.Sprintf("- 🟠 High Risk: %d tenants\n", highRisk)
	summary += fmt.Sprintf("- 🟡 Medium Risk: %d tenants\n", mediumRisk)
	summary += fmt.Sprintf("- 🟢 Low Risk: %d tenants\n\n", lowRisk)

	// Top at-risk tenants
	if criticalRisk + highRisk > 0 {
		summary += "**Top At-Risk Tenants**:\n"
		count := 0
		for _, tenant := range allTenants {
			if tenant.RiskProbability >= 0.65 && count < 5 {
				displayID := FormatPhoneNumberForDisplay(tenant.ID)
				summary += fmt.Sprintf("- %s - %s (%.2f) - $%.0f/month\n",
					displayID, tenant.RiskLevel, tenant.RiskProbability, tenant.CurrentRentAmount)
				count++
			}
		}
	}

	return summary
}

// List high-risk tenants
func (rg *ResponseGenerator) ListHighRisk() string {
	allTenants := rg.GetAllTenants()

	criticalTenants := []string{}
	highRiskTenants := []string{}
	
	for _, tenant := range allTenants {
		displayID := FormatPhoneNumberForDisplay(tenant.ID)
		if tenant.RiskProbability >= 0.85 {
			criticalTenants = append(criticalTenants,
				fmt.Sprintf("%s (%.2f - Critical)", displayID, tenant.RiskProbability))
		} else if tenant.RiskProbability >= 0.65 {
			highRiskTenants = append(highRiskTenants,
				fmt.Sprintf("%s (%.2f - High)", displayID, tenant.RiskProbability))
		}
	}

	if len(criticalTenants) == 0 && len(highRiskTenants) == 0 {
		return "No tenants currently exceed the high-risk threshold (0.65)."
	}

	result := ""
	if len(criticalTenants) > 0 {
		result += "**🔴 Critical Risk Tenants:** " + strings.Join(criticalTenants, ", ") + "\n\n"
	}
	if len(highRiskTenants) > 0 {
		result += "**🟠 High Risk Tenants:** " + strings.Join(highRiskTenants, ", ")
	}

	return result
}

// Compare tenants
func (rg *ResponseGenerator) CompareTenants(tenantIDs []string) string {
	log.Printf("CompareTenants called with: %v", tenantIDs)
	
	if len(tenantIDs) < 2 {
		log.Printf("CompareTenants: Not enough tenant IDs provided: %d", len(tenantIDs))
		return "Please specify at least two tenants to compare."
	}

	tenants := []*ChatbotTenant{}
	errors := []string{}
	
	for _, id := range tenantIDs {
		log.Printf("CompareTenants: Looking up tenant with ID/phone: %s", id)
		if tenant, err := rg.GetTenantByPhone(id); err == nil {
			log.Printf("CompareTenants: Found tenant by phone: %s", id)
			tenants = append(tenants, tenant)
		} else if tenant, exists := rg.GetTenant(id); exists {
			log.Printf("CompareTenants: Found tenant by ID: %s", id)
			tenants = append(tenants, tenant)
		} else {
			log.Printf("CompareTenants: Could not find tenant: %s, error: %v", id, err)
			errors = append(errors, id)
		}
	}

	if len(tenants) < 2 {
		if len(errors) > 0 {
			formattedErrors := make([]string, len(errors))
			for i, err := range errors {
				formattedErrors[i] = FormatPhoneNumberForDisplay(err)
			}
			return fmt.Sprintf("Could not find enough tenants to compare. Could not find tenants with phone numbers: %s. Please check the phone numbers and try again.", strings.Join(formattedErrors, ", "))
		}
		return "Could not find enough tenants to compare. Please provide valid phone numbers."
	}

	comparison := "**Tenant Comparison**\n\n"
	
	for i, tenant := range tenants {
		if i > 0 {
			comparison += "\n"
		}
		displayID := FormatPhoneNumberForDisplay(tenant.ID)
		comparison += fmt.Sprintf("**Tenant %s**\n", displayID)
		comparison += fmt.Sprintf("- Risk Score: %.2f\n", tenant.RiskProbability)
		comparison += fmt.Sprintf("- Risk Level: %s\n", tenant.RiskLevel)
		comparison += fmt.Sprintf("- Late Payments: %d\n", tenant.PreviousLateCount)
		comparison += fmt.Sprintf("- Avg Delay: %.1f days\n", tenant.AvgDelayDays)
		comparison += fmt.Sprintf("- Rent: $%.0f\n", tenant.CurrentRentAmount)
		
		// ✅ ENHANCED: Show additional factors if available
		if tenant.PartialPaymentRatio > 0 {
			comparison += fmt.Sprintf("- Partial Payments: %.0f%%\n", tenant.PartialPaymentRatio*100)
		}
		if tenant.PaymentTrend != 0 {
			trend := "stable"
			if tenant.PaymentTrend > 0 {
				trend = "worsening"
			} else if tenant.PaymentTrend < 0 {
				trend = "improving"
			}
			comparison += fmt.Sprintf("- Payment Trend: %s\n", trend)
		}
	}

	log.Printf("CompareTenants: Returning comparison with %d tenants", len(tenants))
	return comparison
}

// Unknown intent response
func (rg *ResponseGenerator) UnknownIntent() string {
	return `I can help you with:

**Risk Analysis**:
- Explain why a tenant is at specific risk level
- Compare risk between tenants
- List Critical, High, Medium, or Low risk tenants

**Action Planning**:
- Get recommended actions for specific tenants
- Understand intervention strategies for each risk level

**Reporting**:
- Monthly risk summaries with 4-band breakdown
- Payment history details

**Examples**:
- "Why is tenant with phone 01712345678 high risk?"
- "What should I do for tenant 01987654321?"
- "List high risk tenants"
- "List critical risk tenants"
- "Compare tenants 01712345678 and 01987654321"
- "Show monthly summary"

Please rephrase your question using one of these topics.`
}

// Process a chat message
func (rg *ResponseGenerator) ProcessMessage(message string, tenantID string) *ChatResponse {
	startTime := time.Now()

	detector := NewIntentDetector()
	intent, confidence := detector.Detect(message)

	extractedPhones := detector.ExtractPhoneNumbers(message)
	extractedIDs := detector.ExtractTenantIDs(message)
	
	if tenantID == "" {
		if len(extractedPhones) > 0 {
			tenantID = extractedPhones[0]
		} else if len(extractedIDs) > 0 {
			tenantID = extractedIDs[0]
		}
	}

	responseText := ""
	suggestedFollowups := []string{}
	var data any

	switch intent {
	case "EXPLAIN_RISK":
		if tenantID != "" {
			var tenant *ChatbotTenant
			var err error
			if tenant, err = rg.GetTenantByPhone(tenantID); err != nil {
				var exists bool
				tenant, exists = rg.GetTenant(tenantID)
				if !exists {
					responseText = fmt.Sprintf("Could not find tenant with phone number %s. Please check the phone number and try again.", tenantID)
					break
				}
			}
			responseText = rg.ExplainRiskForTenant(tenant)
			suggestedFollowups = []string{
				fmt.Sprintf("What should I do for tenant %s?", tenantID),
				fmt.Sprintf("Compare tenant %s with another tenant", tenantID),
				"List high risk tenants",
			}
			data = map[string]string{"tenant_id": tenantID}
		} else {
			responseText = "Please specify a tenant phone number to explain risk."
		}

	case "RECOMMEND_ACTION":
		if tenantID != "" {
			var tenant *ChatbotTenant
			var err error
			if tenant, err = rg.GetTenantByPhone(tenantID); err != nil {
				var exists bool
				tenant, exists = rg.GetTenant(tenantID)
				if !exists {
					responseText = fmt.Sprintf("Could not find tenant with phone number %s. Please check the phone number and try again.", tenantID)
					break
				}
			}
			responseText = rg.RecommendActionForTenant(tenant)
			suggestedFollowups = []string{
				fmt.Sprintf("Show risk factors for tenant %s", tenantID),
				"What documents do I need?",
				"When should I escalate?",
			}
			data = map[string]string{"tenant_id": tenantID}
		} else {
			responseText = "Please specify a tenant phone number to get recommended actions."
		}

	case "LIST_HIGH_RISK":
		responseText = rg.ListHighRisk()
		suggestedFollowups = []string{
			"Show monthly summary",
			"Compare the top 2 high-risk tenants",
			"What actions for all high-risk tenants?",
		}

	case "MONTHLY_SUMMARY":
		responseText = rg.MonthlySummary()
		suggestedFollowups = []string{
			"List high risk tenants",
			"List critical risk tenants",
			"Compare T100 and T087",
			"Show payment trends",
		}

	case "COMPARE_TENANTS":
		ids := extractedPhones
		if len(ids) == 0 {
			ids = extractedIDs
		}
		if len(ids) == 0 && tenantID != "" {
			ids = append(ids, tenantID)
		}
		log.Printf("COMPARE_TENANTS: extractedPhones=%v, extractedIDs=%v, tenantID=%s, final ids=%v", 
			extractedPhones, extractedIDs, tenantID, ids)
		responseText = rg.CompareTenants(ids)

	case "UNKNOWN":
		responseText = rg.UnknownIntent()

	default:
		responseText = "I understand you're asking about " + strings.ToLower(intent) + ". This feature is under development."
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &ChatResponse{
		Intent:             intent,
		Confidence:         confidence,
		ResponseText:       responseText,
		SuggestedFollowups: suggestedFollowups,
		ProcessingTimeMS:   processingTime,
		Timestamp:          time.Now().Format(time.RFC3339),
		Data:               data,
	}
}

// Global response generator instance
var chatbotResponseGen *ResponseGenerator
var chatbotOnce sync.Once

func getChatbotResponseGenerator() *ResponseGenerator {
	chatbotOnce.Do(func() {
		chatbotResponseGen = NewResponseGenerator()
	})
	return chatbotResponseGen
}

// ChatHandler handles POST requests to process chat messages
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding chat request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	responseGen := getChatbotResponseGenerator()
	response := responseGen.ProcessMessage(req.Message, req.TenantID)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding chat response: %v", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}

// ChatHealthHandler handles health check for chatbot
func ChatHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"status":      "healthy",
		"timestamp":   time.Now().Format(time.RFC3339),
		"version":     "1.0.0-enhanced",
		"tenants":     len(getChatbotResponseGenerator().GetAllTenants()),
		"risk_bands":  "4 (Low < 0.35, Medium 0.35-0.65, High 0.65-0.85, Critical ≥ 0.85)",
		"features":    "Enhanced risk calculation with 6 factors (late count, delay, trend, partial payments, rent burden, tenancy)",
	}
	json.NewEncoder(w).Encode(response)
}
