#!/bin/bash

# Script to rename JPG screenshots to PNG names (GitHub will display both)
# This creates symbolic links or renames files to match README expectations

SCREENSHOTS_DIR="screenshots"

echo "Renaming screenshots for GitHub README..."
echo "=========================================="

# Function to rename/copy file
process_file() {
    local source="$1"
    local target="$2"
    
    if [ -f "$SCREENSHOTS_DIR/$source" ]; then
        # Copy JPG to PNG (GitHub displays both formats)
        cp "$SCREENSHOTS_DIR/$source" "$SCREENSHOTS_DIR/$target"
        echo "✓ Created: $target (from $source)"
        return 0
    fi
    return 1
}

# Get list of all JPG files sorted by timestamp
cd "$SCREENSHOTS_DIR" || exit

# Count files
file_count=$(ls -1 *.jpg *.jpeg 2>/dev/null | wc -l | tr -d ' ')
echo "Found $file_count image files"
echo ""
echo "IMPORTANT: Since files have timestamps, you need to manually identify them."
echo "Here's the order you should rename them based on app flow:"
echo ""
echo "1. login.png - Login screen"
echo "2. register.png - Registration screen"  
echo "3. properties.png - Managed properties list"
echo "4. add_property.png - Add property dialog"
echo "5. property_details.png - Property details (ConcordTower)"
echo "6. add_floor.png - Add floor dialog"
echo "7. update_floor.png - Update floor dialog"
echo "8. tenant_request.png - Send tenant request"
echo "9. payment_details.png - Payment details screen"
echo "10. payment_history.png - Payment history table"
echo "11. send_payment.png - Send payment dialog"
echo "12. adjust_rent.png - Adjust due rent dialog"
echo "13. advance_request.png - Advance payment request"
echo "14. advance_details.png - Advance details screen"
echo "15. chat.png - Chat interface"
echo "16. risk_analysis.png - Risk analysis"
echo "17. monthly_summary.png - Monthly summary"
echo "18. recommendations.png - Recommendations"
echo "19. high_risk.png - High risk tenants"
echo "20. compare.png - Tenant comparison"
echo "21. notifications.png - Notifications list"
echo "22. notification_actions.png - Notification with actions"
echo "23. notification_accepted.png - Accepted notification"
echo "24. settings.png - Settings screen"
echo "25. language.png - Language selection"
echo "26. bengali_ui.png - Bengali UI"
echo "27. tenant_properties.png - Tenant properties"
echo ""
echo "To rename manually, use:"
echo "  cd screenshots"
echo "  mv '2025-12-30 19.XX.XX.jpg' 'login.png'"
echo "  (repeat for each file)"
echo ""
echo "Or use this script with a mapping file (see RENAME_INSTRUCTIONS.md)"

cd ..

