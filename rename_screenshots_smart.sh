#!/bin/bash

# Smart screenshot renaming script
# This script renames screenshots based on timestamp patterns and folder organization

SCREENSHOTS_DIR="screenshots"

echo "Starting smart screenshot renaming..."
echo "======================================"

# Function to rename file
rename_file() {
    local source="$1"
    local target="$2"
    
    if [ -f "$SCREENSHOTS_DIR/$source" ]; then
        mv "$SCREENSHOTS_DIR/$source" "$SCREENSHOTS_DIR/$target"
        echo "✓ Renamed: $source -> $target"
        return 0
    fi
    return 1
}

# Check folders first (some screenshots might be in subfolders)
if [ -d "$SCREENSHOTS_DIR/login Screen" ]; then
    find "$SCREENSHOTS_DIR/login Screen" -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" \) | head -1 | while read file; do
        mv "$file" "$SCREENSHOTS_DIR/login.png" 2>/dev/null && echo "✓ Moved login screenshot"
    done
fi

if [ -d "$SCREENSHOTS_DIR/Register Screen" ]; then
    find "$SCREENSHOTS_DIR/Register Screen" -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" \) | head -1 | while read file; do
        mv "$file" "$SCREENSHOTS_DIR/register.png" 2>/dev/null && echo "✓ Moved register screenshot"
    done
fi

if [ -d "$SCREENSHOTS_DIR/Managed Property Screen" ]; then
    find "$SCREENSHOTS_DIR/Managed Property Screen" -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" \) | head -1 | while read file; do
        mv "$file" "$SCREENSHOTS_DIR/properties.png" 2>/dev/null && echo "✓ Moved properties screenshot"
    done
fi

if [ -d "$SCREENSHOTS_DIR/Tenant Property screen" ]; then
    find "$SCREENSHOTS_DIR/Tenant Property screen" -type f \( -name "*.jpg" -o -name "*.jpeg" -o -name "*.png" \) | head -1 | while read file; do
        mv "$file" "$SCREENSHOTS_DIR/tenant_properties.png" 2>/dev/null && echo "✓ Moved tenant properties screenshot"
    done
fi

echo ""
echo "Now processing timestamped files..."
echo "Please review the files and manually identify which screenshots match which features."
echo ""
echo "To help identify, here are the expected screenshots:"
echo "  Authentication: login.png, register.png"
echo "  Properties: properties.png, add_property.png, property_details.png"
echo "  Floors: add_floor.png, update_floor.png, tenant_request.png"
echo "  Payments: payment_details.png, payment_history.png, send_payment.png"
echo "  Advance: advance_request.png, advance_details.png, adjust_rent.png"
echo "  Chatbot: chat.png, risk_analysis.png, monthly_summary.png, recommendations.png, high_risk.png, compare.png"
echo "  Notifications: notifications.png, notification_actions.png, notification_accepted.png"
echo "  Settings: settings.png, language.png, bengali_ui.png"
echo ""
echo "You can manually rename files using:"
echo "  mv 'screenshots/2025-12-30 19.XX.XX.jpg' 'screenshots/desired_name.png'"

