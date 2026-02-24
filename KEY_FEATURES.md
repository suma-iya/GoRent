# Key Features of RentApp - Property Rental Management System

## 1. User Management & Authentication
- **User Registration**: Register with email, phone number, and password
- **Secure Login**: JWT-based authentication system
- **Role-Based Access Control**: 
  - Property Manager: Full access to all management features
  - Tenant: Limited access to view assigned properties and payment history
- **Password Security**: Bcrypt hashing for secure password storage
- **Phone Number Lookup**: Search for users by phone number

## 2. Property Management
- **Add Properties**: Create new properties with name, address, and optional photo upload
- **View Properties**: List all properties managed by the user
- **Property Details**: Detailed view showing property information and associated floors/units
- **Property Photos**: Upload and display property images

## 3. Floor/Unit Management
- **Add Floors/Units**: Create multiple floors or units within a property
- **Set Rent Amounts**: Assign monthly rent for each floor/unit
- **Update Floor Information**: Modify floor details and rent amounts
- **Floor Status Tracking**: Track occupancy status (vacant, occupied, pending)

## 4. Tenant Management
- **Send Tenant Requests**: Managers can send rental requests to tenants via phone number
- **Accept/Reject Requests**: Tenants can accept or reject rental requests through notifications
- **Direct Tenant Assignment**: Assign tenants directly to floors/units
- **Remove Tenants**: Remove tenants from floors when needed
- **Tenant Information View**: View tenant details and payment history
- **Tenant Properties View**: Tenants can view all properties where they are assigned

## 5. Payment Management
- **Record Payments**: Track rental payments with amount, date, and timestamps
- **Payment History**: Complete history of all payments for each floor/unit
- **Payment Details**: Detailed view of individual payment records
- **Late Payment Tracking**: Automatic detection and tracking of late payments
- **Payment History Export**: Download payment history as CSV file for external analysis
- **Payment Analytics**: View payment trends and statistics

## 6. Advance Payment System
- **Request Advance Payments**: Managers can request advance payments from tenants
- **Process Advance Payments**: Tenants can accept or reject advance payment requests
- **Advance Payment History**: Track all advance payment requests and their status
- **Advance Payment Details**: View detailed information about advance payments
- **Cancel Advance Payments**: Cancel pending advance payment requests

## 7. AI-Powered Chatbot
- **Intelligent Risk Analysis**: AI chatbot provides tenant risk assessment based on payment history
- **Intent Detection**: Natural language understanding with pattern matching for:
  - Risk explanations (EXPLAIN_RISK)
  - Action recommendations (RECOMMEND_ACTION)
  - High-risk tenant listings (LIST_HIGH_RISK)
  - Monthly summaries (MONTHLY_SUMMARY)
  - Tenant comparisons (COMPARE_TENANTS)
  - Payment history queries (PAYMENT_HISTORY)
  - Lease renewal guidance (LEASE_RENEWAL)
- **Risk Calculation Algorithm**: 
  - Calculates risk probability based on late payments, average delay, and tenancy duration
  - Risk levels: High (≥0.7), Medium (0.4-0.7), Low (<0.4)
- **Phone Number Recognition**: Extracts and normalizes Bangladeshi phone numbers from queries
- **Markdown Support**: Rich text formatting in chatbot responses
- **Suggested Follow-ups**: Provides contextual follow-up questions
- **Chat History**: Persistent message history stored locally (last 100 messages)
- **Quick Action Buttons**: Pre-defined queries for common questions

## 8. Real-Time Notification System
- **Firebase Cloud Messaging**: Push notifications for mobile devices
- **In-App Notifications**: Notification center within the application
- **Notification Types**:
  - Tenant request notifications
  - Payment reminders
  - Payment confirmations
  - Advance payment requests
  - System notifications
- **Read/Unread Status**: Track notification read status
- **Notification Actions**: Accept/reject actions directly from notifications
- **Notification History**: View all past notifications
- **Delete Notifications**: Remove unwanted notifications
- **Conversation Threads**: View conversation history for tenant requests
- **Automated Monthly Reminders**: Cron-based automated payment reminders

## 9. Pending Payment Notifications
- **Track Pending Payments**: View all pending payment notifications for a floor
- **Payment Notification Management**: Send and manage payment notifications

## 10. Multi-Language Support (Localization)
- **English & Bengali**: Support for both English and Bengali (বাংলা) languages
- **Dynamic Language Switching**: Change language without app restart
- **Localized Content**: All UI elements and messages are localized

## 11. Data Export & Sharing
- **CSV Export**: Download payment history as CSV file
- **Share Functionality**: Share exported data via other apps

## 12. Settings & Preferences
- **Language Settings**: Change app language preference
- **User Preferences**: Manage app settings and preferences

## 13. Security Features
- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: Bcrypt for password security
- **CORS Protection**: Cross-origin request protection
- **Rate Limiting**: API rate limiting to prevent abuse
- **Input Validation**: Server-side validation for all inputs
- **Protected Routes**: Middleware-based route protection

## 14. Automated Scheduling
- **Monthly Payment Reminders**: Automated cron jobs for payment reminders
- **Scheduled Notifications**: Background task scheduling for automated notifications

## 15. Cross-Platform Mobile Application
- **Android & iOS Support**: Single codebase for both platforms using Flutter
- **Material Design 3**: Modern, beautiful UI components
- **Responsive Design**: Works on various screen sizes

## 16. Backend Architecture
- **RESTful API**: Comprehensive REST API with Go backend
- **MySQL Database**: Relational database for data persistence
- **Docker Support**: Containerized deployment for easy setup
- **Connection Pooling**: Efficient database connection management
- **Error Handling**: Comprehensive error handling and logging

## 17. Developer Features
- **Health Check Endpoints**: Monitor service health
- **Test Endpoints**: Testing utilities for FCM and notifications
- **API Documentation**: Well-structured API endpoints

## Technical Stack Highlights
- **Backend**: Go (Golang) with Gorilla Mux router
- **Frontend**: Flutter (cross-platform mobile framework)
- **Database**: MySQL 8.0
- **Authentication**: JWT tokens
- **Notifications**: Firebase Cloud Messaging
- **Deployment**: Docker & Docker Compose
- **State Management**: Provider pattern (Flutter)
- **Scheduling**: Cron jobs for automated tasks

---

**Note**: This is an intelligent property rental management system that combines traditional property management with AI-powered risk analysis, automated notifications, and comprehensive payment tracking to reduce manual administrative overhead by approximately 60%.


