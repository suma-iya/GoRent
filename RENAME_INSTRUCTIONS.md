# Screenshot Renaming Instructions

## Quick Start

To run the renaming script, use the **full filename**:

```bash
./rename_screenshots.sh
```

**Note**: Make sure you're in the project root directory (`rentApp/`)

## Manual Renaming Guide

Since your screenshots have timestamp-based names, here's a guide to help you identify and rename them:

### Step 1: Identify Your Screenshots

Based on the image descriptions, here's what each screenshot should be named:

#### Authentication & Registration
- **login.png** - Login screen (shows phone number and password fields)
- **register.png** - Registration screen (shows Full Name, Phone Number, Password fields)

#### Property Management  
- **properties.png** - Managed Properties list (shows "Concord Tower" card)
- **add_property.png** - Add Property dialog (shows property photo, name, address fields)
- **property_details.png** - Property details screen (shows "ConcordTower" header with floors)
- **add_floor.png** - Add Floor dialog (shows Floor Name and Monthly Rent fields)
- **update_floor.png** - Update Floor dialog (similar to add floor but for editing)
- **tenant_request.png** - Send Tenant Request dialog (shows phone number input)

#### Payment Management
- **payment_details.png** - Payment details screen (shows rent amount and payment history)
- **payment_history.png** - Payment history table (shows date, rent, received money columns)
- **send_payment.png** - Send Payment dialog (shows Payment Amount and Electricity Bill fields)
- **adjust_rent.png** - Adjust due rent dialog (shows Add/Subtract toggle, rent amount, electricity bill)
- **advance_request.png** - Advance Payment Request dialog (shows phone number and amount)
- **advance_details.png** - Advance Details screen (shows tenant name, amount, date, Deduct Advance button)

#### AI Chatbot
- **chat.png** - Chat interface (shows "Rent Decision C..." title, empty chat or conversation)
- **risk_analysis.png** - Risk analysis response (shows high-risk tenants list)
- **monthly_summary.png** - Monthly Risk Summary (shows total tenants, rent at risk, risk distribution)
- **recommendations.png** - Recommended Actions (shows action list for a tenant)
- **high_risk.png** - High Risk Tenants list (shows tenant phone numbers with risk scores)
- **compare.png** - Tenant Comparison (shows comparison between tenants)

#### Notifications
- **notifications.png** - Notifications list (shows multiple notification cards)
- **notification_actions.png** - Notification with action buttons (shows Accept/Reject/Comment buttons)
- **notification_accepted.png** - Accepted notification (shows green "accepted" status)

#### Settings & Localization
- **settings.png** - Settings screen (shows Push Notifications and App Information sections)
- **language.png** - Language Selection dialog (shows English and Bengali options)
- **bengali_ui.png** - Bengali UI example (shows interface in Bengali language)

#### Additional
- **tenant_properties.png** - Tenant Properties screen (shows properties where user is a tenant)

### Step 2: Rename Files

You can rename files using one of these methods:

#### Method 1: Using Terminal (Recommended)

```bash
cd screenshots

# Example renaming (replace with your actual filenames):
mv "2025-12-30 19.31.24.jpg" "login.png"
mv "2025-12-30 19.34.20.jpg" "register.png"
# ... continue for all files
```

#### Method 2: Using Finder (macOS)

1. Open the `screenshots` folder in Finder
2. Click on each file to preview it
3. Right-click → Rename
4. Change to the appropriate name (e.g., `login.png`)

#### Method 3: Batch Rename Script

Create a file `batch_rename.sh` with your mappings:

```bash
#!/bin/bash
cd screenshots

# Add your mappings here based on what you see in each file
mv "2025-12-30 19.31.24.jpg" "login.png"
mv "2025-12-30 19.34.20.jpg" "register.png"
# ... add all your mappings
```

Then run: `chmod +x batch_rename.sh && ./batch_rename.sh`

### Step 3: Verify

After renaming, verify you have all required files:

```bash
cd screenshots
ls -1 *.png | sort
```

You should see files like:
- login.png
- register.png
- properties.png
- add_property.png
- ... (all the files listed above)

## Tips

1. **Preview First**: Before renaming, preview each image to confirm what it shows
2. **Keep Originals**: Consider keeping a backup of original filenames
3. **Use PNG Format**: The README expects `.png` files, so convert `.jpg` files if needed:
   ```bash
   # Convert JPG to PNG (requires ImageMagick or sips on macOS)
   sips -s format png "file.jpg" --out "file.png"
   ```
4. **Chronological Order**: Files are timestamped, so earlier timestamps likely correspond to earlier screens (login → register → properties → etc.)

## Need Help?

If you're unsure which file is which:
1. Open each image in Preview/Photos
2. Match it to the descriptions above
3. Rename accordingly

Or, you can tell me which files you're unsure about and I can help identify them based on the descriptions!

