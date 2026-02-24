#!/usr/bin/env python3
"""
Smart screenshot renaming script based on image descriptions.
This script renames screenshots to match README naming convention.
"""

import os
import shutil
from pathlib import Path

# Mapping based on image descriptions provided
# Format: (description_keywords, target_filename)
SCREENSHOT_MAPPING = [
    # Authentication
    (["login", "Login"], "login.png"),
    (["register", "Register", "registration"], "register.png"),
    
    # Properties
    (["Managed Properties", "managed properties", "properties list"], "properties.png"),
    (["Add Property", "add property"], "add_property.png"),
    (["Property details", "property details", "ConcordTower", "property information"], "property_details.png"),
    
    # Floors
    (["Add Floor", "add floor"], "add_floor.png"),
    (["Update Floor", "update floor"], "update_floor.png"),
    (["Send Tenant Request", "tenant request", "send tenant"], "tenant_request.png"),
    
    # Payments
    (["Payment details", "payment details"], "payment_details.png"),
    (["Payment history", "payment history"], "payment_history.png"),
    (["Send Payment", "send payment"], "send_payment.png"),
    (["Adjust due rent", "adjust rent", "adjust due"], "adjust_rent.png"),
    (["Advance Payment Request", "advance payment request", "advance request"], "advance_request.png"),
    (["Advance Details", "advance details", "Deduct Advance"], "advance_details.png"),
    
    # Chatbot
    (["Rent Decision", "chat", "chatbot", "AI Chatbot"], "chat.png"),
    (["risk factors", "risk analysis", "high risk tenants"], "risk_analysis.png"),
    (["Monthly Risk Summary", "monthly summary", "monthly risk"], "monthly_summary.png"),
    (["Recommended Actions", "recommendations", "recommended actions"], "recommendations.png"),
    (["High Risk", "high-risk", "list high risk"], "high_risk.png"),
    (["Compare", "comparison", "compare tenants"], "compare.png"),
    
    # Notifications
    (["Notifications", "notifications list"], "notifications.png"),
    (["Accept", "Reject", "notification actions", "action buttons"], "notification_actions.png"),
    (["accepted", "notification accepted"], "notification_accepted.png"),
    
    # Settings
    (["Settings", "settings screen"], "settings.png"),
    (["Select Language", "language selection", "Language"], "language.png"),
    (["Bengali", "বাংলা", "bengali ui"], "bengali_ui.png"),
    
    # Tenant Properties
    (["Tenant Properties", "tenant property"], "tenant_properties.png"),
]

def find_and_rename_screenshots(screenshots_dir):
    """Find and rename screenshots based on folder names and file patterns."""
    screenshots_path = Path(screenshots_dir)
    
    if not screenshots_path.exists():
        print(f"Error: {screenshots_dir} directory not found!")
        return
    
    renamed_count = 0
    
    # First, handle files in subdirectories
    for subdir in screenshots_path.iterdir():
        if subdir.is_dir():
            # Check folder name against mapping
            folder_name = subdir.name.lower()
            target_file = None
            
            if "login" in folder_name:
                target_file = "login.png"
            elif "register" in folder_name:
                target_file = "register.png"
            elif "managed property" in folder_name or "properties" in folder_name:
                target_file = "properties.png"
            elif "tenant property" in folder_name:
                target_file = "tenant_properties.png"
            
            if target_file:
                # Move first image file from subdirectory
                for img_file in subdir.glob("*"):
                    if img_file.is_file() and img_file.suffix.lower() in ['.jpg', '.jpeg', '.png']:
                        target_path = screenshots_path / target_file
                        if not target_path.exists():
                            shutil.move(str(img_file), str(target_path))
                            print(f"✓ Moved {img_file.name} from {subdir.name} -> {target_file}")
                            renamed_count += 1
                        break
    
    print(f"\nRenamed {renamed_count} files from subdirectories.")
    print("\nNote: For timestamped files, you'll need to manually identify and rename them.")
    print("Use the mapping guide in screenshots/README.md to match files to their correct names.")

if __name__ == "__main__":
    find_and_rename_screenshots("screenshots")

