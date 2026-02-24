#!/bin/bash

# Script to rename screenshot files to match README naming convention
# Place all your screenshot images in the screenshots/ folder, then run this script

SCREENSHOTS_DIR="screenshots"

# Check if screenshots directory exists
if [ ! -d "$SCREENSHOTS_DIR" ]; then
    echo "Creating screenshots directory..."
    mkdir -p "$SCREENSHOTS_DIR"
fi

# Function to rename file if it exists
rename_file() {
    local old_name="$1"
    local new_name="$2"
    
    if [ -f "$SCREENSHOTS_DIR/$old_name" ]; then
        mv "$SCREENSHOTS_DIR/$old_name" "$SCREENSHOTS_DIR/$new_name"
        echo "✓ Renamed: $old_name -> $new_name"
    elif [ -f "$old_name" ]; then
        mv "$old_name" "$SCREENSHOTS_DIR/$new_name"
        echo "✓ Moved and renamed: $old_name -> $new_name"
    fi
}

echo "Starting screenshot renaming process..."
echo "======================================"

# Authentication & Registration
rename_file "login_screen.png" "login.png"
rename_file "register_screen.png" "register.png"
rename_file "registration.png" "register.png"

# Property Management
rename_file "managed_properties.png" "properties.png"
rename_file "properties_screen.png" "properties.png"
rename_file "add_property.png" "add_property.png"
rename_file "property_details.png" "property_details.png"
rename_file "add_floor.png" "add_floor.png"
rename_file "update_floor.png" "update_floor.png"
rename_file "send_tenant_request.png" "tenant_request.png"
rename_file "tenant_request.png" "tenant_request.png"
rename_file "concord_tower.png" "property_details.png"

# Payment Management
rename_file "payment_details.png" "payment_details.png"
rename_file "payment_history.png" "payment_history.png"
rename_file "send_payment.png" "send_payment.png"
rename_file "adjust_due_rent.png" "adjust_rent.png"
rename_file "adjust_rent.png" "adjust_rent.png"
rename_file "advance_payment_request.png" "advance_request.png"
rename_file "advance_details.png" "advance_details.png"
rename_file "deduct_advance.png" "advance_details.png"

# AI Chatbot
rename_file "chat_screen.png" "chat.png"
rename_file "chatbot.png" "chat.png"
rename_file "rent_decision_chat.png" "chat.png"
rename_file "risk_analysis.png" "risk_analysis.png"
rename_file "monthly_summary.png" "monthly_summary.png"
rename_file "recommendations.png" "recommendations.png"
rename_file "high_risk_tenants.png" "high_risk.png"
rename_file "tenant_comparison.png" "compare.png"
rename_file "compare_tenants.png" "compare.png"

# Notifications
rename_file "notifications.png" "notifications.png"
rename_file "notifications_screen.png" "notifications.png"
rename_file "notification_actions.png" "notification_actions.png"
rename_file "notification_accepted.png" "notification_accepted.png"

# Settings & Localization
rename_file "settings.png" "settings.png"
rename_file "settings_screen.png" "settings.png"
rename_file "language_selection.png" "language.png"
rename_file "select_language.png" "language.png"
rename_file "bengali_ui.png" "bengali_ui.png"
rename_file "tenant_properties.png" "tenant_properties.png"

echo ""
echo "======================================"
echo "Renaming complete!"
echo ""
echo "If some files weren't renamed, you may need to manually rename them."
echo "Expected filenames in README:"
echo "  - login.png, register.png"
echo "  - properties.png, add_property.png, property_details.png"
echo "  - add_floor.png, update_floor.png, tenant_request.png"
echo "  - payment_details.png, payment_history.png, send_payment.png"
echo "  - adjust_rent.png, advance_request.png, advance_details.png"
echo "  - chat.png, risk_analysis.png, monthly_summary.png"
echo "  - recommendations.png, high_risk.png, compare.png"
echo "  - notifications.png, notification_actions.png, notification_accepted.png"
echo "  - settings.png, language.png, bengali_ui.png"

