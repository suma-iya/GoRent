import 'package:flutter/material.dart';

class AppLocalizations {
  final Locale locale;
  
  AppLocalizations(this.locale);
  
  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }
  
  static const Map<String, Map<String, String>> _localizedValues = {
    'en': {
      'appTitle': 'Go Rent',
      'managedProperties': 'Managed Properties',
      'tenantProperties': 'Tenant Properties',
      'addProperty': 'Add Property',
      'propertyName': 'Property Name',
      'propertyNameHint': 'Enter property name',
      'address': 'Address',
      'addressHint': 'Enter property address',
      'propertyPhoto': 'Property Photo (Optional)',
      'addPhotoMessage': 'Add a photo of your property',
      'camera': 'Camera',
      'gallery': 'Gallery',
      'remove': 'Remove',
      'change': 'Change',
      'cancel': 'Cancel',
      'confirm': 'Confirm',
      'addPropertySuccess': 'Property added successfully',
      'fillRequiredFields': 'Please fill all required fields',
      'failedToAddProperty': 'Failed to add property',
      'noManagedProperties': 'No managed property yet',
      'noManagedPropertiesSubtitle': 'Add your first property to get started',
      'noTenantProperties': 'No tenant property yet',
      'noTenantPropertiesSubtitle': 'You\'ll see properties where you\'re a tenant here',
      'loadingManagedProperties': 'Loading managed properties...',
      'loadingTenantProperties': 'Loading tenant properties...',
      'somethingWentWrong': 'Something went wrong',
      'tryAgain': 'Try Again',
      'managedProperty': 'Managed Property',
      'tenantProperty': 'Tenant Property',
      'notifications': 'Notifications',
      'language': 'Language',
      'logout': 'Logout',
      'logoutConfirmation': 'Are you sure you want to logout?',
      'selectLanguage': 'Select Language',
      'english': 'English',
      'bangla': 'বাংলা',
      'manager': 'Manager',
      'tenant': 'Tenant',
      'failedToLoadProperties': 'Failed to load properties. Please try again.',
      'add': 'Add',
      'addFloor': 'Add Floor',
      'addTenant': 'Add Tenant',
      'update': 'Update',
      'send': 'Send',
      'sendRequest': 'Send Request',
      'sendPaymentNotification': 'Send Payment',
      'sendAdvancePaymentRequest': 'Send Advance Payment Request',
      'yes': 'Yes',
      'no': 'No',
      'retry': 'Retry',
      'addFirstFloor': 'Add First Floor',
      'occupied': 'Occupied',
      'pending': 'Pending',
      'available': 'Available',
      'tapToView': 'Tap to view',
      'requestPending': 'Request Pending',
      'rent': 'Rent',
      'unknownUser': 'Unknown User',
      'floorName': 'Floor Name',
      'floorNameHint': 'e.g., Ground Floor, First Floor',
      'monthlyRent': 'Monthly Rent',
      'rentHint': 'Enter amount in rupees',
      'tenantName': 'Tenant Name',
      'tenantNameHint': 'Enter tenant name',
      'phoneNumber': 'Phone Number',
      'phoneNumberHint': 'Enter tenant phone number',
      'removeTenant': 'Remove Tenant',
      'removeTenantConfirmation': 'Are you sure you want to remove the tenant from this floor?',
      'updateFloor': 'Update Floor',
      'sendTenantRequest': 'Send Tenant Request',
      'cancelRequest': 'Cancel Request',
      'cancelRequestConfirmation': 'Are you sure you want to cancel this tenant request?',
      'paymentAmount': 'Payment Amount',
      'paymentAmountHint': 'Enter payment amount',
      'electricityBill': 'Electricity Bill',
      'electricityBillHint': 'Enter electricity bill amount to pay',
      'advanceAmount': 'Amount',
      'advanceAmountHint': 'Enter advance payment amount',
      'selectMonth': 'Select a month',
      'cancelAdvancePayment': 'Cancel Advance Payment Request',
      'cancelAdvancePaymentConfirmation': 'Are you sure you want to cancel the pending advance payment request for this floor?',
      'yesCancel': 'Yes, Cancel',
      'noFloorsAdded': 'No floors added yet',
      'fillAllFields': 'Please fill all fields',
      'invalidRentAmount': 'Invalid rent amount',
      'failedToAddFloor': 'Failed to add floor',
      'failedToAddTenant': 'Failed to add tenant',
      'failedToRemoveTenant': 'Failed to remove tenant',
      'tenantRemovedSuccessfully': 'Tenant removed successfully',
      'failedToUpdateFloor': 'Failed to update floor',
      'pleaseEnterPhoneNumber': 'Please enter a phone number',
      'tenantRequestSentSuccessfully': 'Tenant request sent successfully',
      'requestCancelledSuccessfully': 'Request cancelled successfully',
      'failedToCancelRequest': 'Failed to cancel request',
      'pleaseEnterValidAmount': 'Please enter a valid amount',
      'pleaseEnterValidElectricityBill': 'Please enter a valid electricity bill amount',
      'paymentNotificationSent': 'Payment notification sent!',
      'pleaseSelectMonth': 'Please select a month',
      'advancePaymentRequestSentSuccessfully': 'Advance payment request sent successfully!',
      'failedToSendAdvancePaymentRequest': 'Failed to send advance payment request',
      'advancePaymentRequestCancelledSuccessfully': 'Advance payment request cancelled successfully!',
      'failedToCancelAdvancePaymentRequest': 'Failed to cancel advance payment request',
      'amount': 'Amount',
      'month': 'Month',
      'january': 'January',
      'february': 'February',
      'march': 'March',
      'april': 'April',
      'may': 'May',
      'june': 'June',
      'july': 'July',
      'august': 'August',
      'september': 'September',
      'october': 'October',
      'november': 'November',
      'december': 'December',
      'enterTenantName': 'Enter tenant name',
      'enterPhoneNumber': 'Enter tenant phone number',
      'errorPrefix': 'Error',
      'enterPaymentAmount': 'Enter payment amount',
      'enterElectricityBill': 'Enter electricity bill amount to pay',
      'cancelAdvancePaymentRequest': 'Cancel Advance Payment Request',
      'cancelAdvancePaymentRequestConfirmation': 'Are you sure you want to cancel the pending advance payment request for this floor?',
      'yesCancelAdvancePaymentRequest': 'Yes, Cancel',
      'noFloorsAddedYet': 'No floors added yet',
      'enterAdvancePaymentAmount': 'Enter advance payment amount',
  'showAdvanceDetails': 'Show Advance Details',
  'advanceDetails': 'Advance Details',
  'noAdvanceDetails': 'No advance details found for this floor',
  'advanceDetailsForFloor': 'Advance Details for {floorName}',
      'currency': 'tk/month',

      "paymentDetails": "Payment details",
      "tk": "tk",
      "adjustDueRent": "Adjust due rent",
      "subtract": "Subtract",
      "rentAmount": "Rent amount",
      "electricityBillOptional": "Electricity bill",
      "preview": "Preview",
      "current": "Current",
      "newTotal": "New total",
      "addAmount": "Add amount",
      "paymentHistory": "Payment history",
      "date": "Date",
      "newAddedRent": "New added rent",
      "receivedMoney": "Received money",
      "dueRent": "Due rent",
      "newAddedElectricityBill": "New added electricity bill",
      "paidElectricity": "Paid electricity",
      "dueElectricityBill": "Due electricity bill",
      "records": "Records",
      "noMonth": "No month selected",
      "paymentRecord": "Payment record",
      "pleaseEnterValidPositiveAmount": "Please enter a valid positive amount",
      "dueRentCannotBeNegative": "Due rent cannot be negative",
      "increased": "Increased",
      "decreased": "Decreased",
      "successfully": "Successfully",
      "subtractAmount": "Subtract amount",
      "loadingNotifications": "Loading notifications...",
      "noNotificationsYet": "No notifications yet",
      "youllSeeYourNotificationsHere": "You'll see your notifications here",
      "accept": "Accept",
      "reject": "Reject",
      "addComment": "Comment",
      "addYourComment": "Add your comment",
      "pleaseEnterAComment": "Please enter a comment",
      "commentSentSuccessfully": "Comment sent successfully",
      "failedToSendComment": "Failed to send comment",
      "failedToProcessTheAction": "Failed to process the action"
        },
    'bn': {
      'appTitle': 'গো রেন্ট',
      'managedProperties': 'পরিচালিত সম্পত্তি',
      'tenantProperties': 'ভাড়াটে সম্পত্তি',
      'addProperty': 'সম্পত্তি যোগ করুন',
      'propertyName': 'সম্পত্তির নাম',
      'propertyNameHint': 'সম্পত্তির নাম লিখুন',
      'address': 'ঠিকানা',
      'addressHint': 'সম্পত্তির ঠিকানা লিখুন',
      'propertyPhoto': 'সম্পত্তির ছবি (ঐচ্ছিক)',
      'addPhotoMessage': 'আপনার সম্পত্তির একটি ছবি যোগ করুন',
      'camera': 'ক্যামেরা',
      'gallery': 'গ্যালারি',
      'remove': 'সরান',
      'change': 'পরিবর্তন',
      'cancel': 'বাতিল',
      'confirm': 'নিশ্চিত করুন',
      'addPropertySuccess': 'সম্পত্তি সফলভাবে যোগ করা হয়েছে',
      'fillRequiredFields': 'অনুগ্রহ করে সমস্ত প্রয়োজনীয় ক্ষেত্র পূরণ করুন',
      'failedToAddProperty': 'সম্পত্তি যোগ করতে ব্যর্থ',
      'noManagedProperties': 'এখনও কোন পরিচালিত সম্পত্তি নেই',
      'noManagedPropertiesSubtitle': 'শুরু করতে আপনার প্রথম সম্পত্তি যোগ করুন',
      'noTenantProperties': 'কোন ভাড়াটে সম্পত্তি পাওয়া যায়নি',
      'noTenantPropertiesSubtitle': 'আপনি যেখানে ভাড়াটে সেখানে সম্পত্তি দেখতে পাবেন',
      'loadingManagedProperties': 'পরিচালিত সম্পত্তি লোড হচ্ছে...',
      'loadingTenantProperties': 'ভাড়াটে সম্পত্তি লোড হচ্ছে...',
      'somethingWentWrong': 'কিছু ভুল হয়েছে',
      'tryAgain': 'আবার চেষ্টা করুন',
      'managedProperty': 'পরিচালিত সম্পত্তি',
      'tenantProperty': 'ভাড়াটে সম্পত্তি',
      'notifications': 'বিজ্ঞপ্তি',
      'language': 'ভাষা',
      'logout': 'লগআউট',
      'logoutConfirmation': 'আপনি কি নিশ্চিত যে আপনি লগআউট করতে চান?',
      'selectLanguage': 'ভাষা নির্বাচন করুন',
      'english': 'English',
      'bangla': 'বাংলা',
      'manager': 'পরিচালক',
      'tenant': 'ভাড়াটে',
      'failedToLoadProperties': 'সম্পত্তি লোড করতে ব্যর্থ। অনুগ্রহ করে আবার চেষ্টা করুন।',
      'add': 'যোগ করুন',
      'addFloor': 'ফ্লোর যোগ করুন',
      'addTenant': 'ভাড়াটে যোগ করুন',
      'update': 'আপডেট করুন',
      'send': 'পাঠান',
      'sendRequest': 'অনুরোধ পাঠান',
      'sendPaymentNotification': 'পেমেন্ট পাঠান',
      'sendAdvancePaymentRequest': 'অগ্রিম পেমেন্ট অনুরোধ পাঠান',
      'yes': 'হ্যাঁ',
      'no': 'না',
      'retry': 'পুনরায় চেষ্টা করুন',
      'addFirstFloor': 'প্রথম ফ্লোর যোগ করুন',
      'occupied': 'দখলকৃত',
      'pending': 'বিচারাধীন',
      'available': 'খালি',
      'tapToView': 'দেখতে ট্যাপ করুন',
      'requestPending': 'অনুরোধ বিচারাধীন',
      'rent': 'ভাড়া',
      'unknownUser': 'অজানা ব্যবহারকারী',
      'floorName': 'ফ্লোরের নাম',
      'floorNameHint': 'যেমন, নিচতলা, প্রথম তলা',
      'monthlyRent': 'মাসিক ভাড়া',
      'rentHint': 'টাকায় পরিমাণ লিখুন',
      'tenantName': 'ভাড়াটের নাম',
      'tenantNameHint': 'ভাড়াটের নাম লিখুন',
      'phoneNumber': 'ফোন নম্বর',
      'phoneNumberHint': 'ভাড়াটের ফোন নম্বর লিখুন',
      'removeTenant': 'ভাড়াটে সরান',
      'removeTenantConfirmation': 'আপনি কি নিশ্চিত যে আপনি এই ফ্লোর থেকে ভাড়াটেকে সরাতে চান?',
      'updateFloor': 'ফ্লোর আপডেট করুন',
      'sendTenantRequest': 'ভাড়াটে অনুরোধ পাঠান',
      'cancelRequest': 'অনুরোধ বাতিল করুন',
      'cancelRequestConfirmation': 'আপনি কি নিশ্চিত যে আপনি এই ভাড়াটে অনুরোধ বাতিল করতে চান?',
      'paymentAmount': 'পেমেন্ট পরিমাণ',
      'paymentAmountHint': 'পেমেন্ট পরিমাণ লিখুন',
      'electricityBill': 'বিদ্যুৎ বিল',
      'electricityBillHint': 'পরিশোধ করার বিদ্যুৎ বিলের পরিমাণ লিখুন',
      'advanceAmount': 'পরিমাণ',
      'advanceAmountHint': 'অগ্রিম পেমেন্টের পরিমাণ লিখুন',
      'selectMonth': 'একটি মাস নির্বাচন করুন',
      'cancelAdvancePayment': 'অগ্রিম পেমেন্ট অনুরোধ বাতিল করুন',
      'cancelAdvancePaymentConfirmation': 'আপনি কি নিশ্চিত যে আপনি এই ফ্লোরের জন্য অপেক্ষমান অগ্রিম পেমেন্ট অনুরোধ বাতিল করতে চান?',
      'yesCancel': 'হ্যাঁ, বাতিল করুন',
      'noFloorsAdded': 'এখনও কোন ফ্লোর যোগ করা হয়নি',
      'fillAllFields': 'অনুগ্রহ করে সমস্ত ক্ষেত্র পূরণ করুন',
      'invalidRentAmount': 'অবৈধ ভাড়ার পরিমাণ',
      'failedToAddFloor': 'ফ্লোর যোগ করতে ব্যর্থ',
      'failedToAddTenant': 'ভাড়াটে যোগ করতে ব্যর্থ',
      'failedToRemoveTenant': 'ভাড়াটে সরাতে ব্যর্থ',
      'tenantRemovedSuccessfully': 'ভাড়াটে সফলভাবে সরানো হয়েছে',
      'failedToUpdateFloor': 'ফ্লোর আপডেট করতে ব্যর্থ',
      'pleaseEnterPhoneNumber': 'অনুগ্রহ করে একটি ফোন নম্বর লিখুন',
      'tenantRequestSentSuccessfully': 'ভাড়াটে অনুরোধ সফলভাবে পাঠানো হয়েছে',
      'requestCancelledSuccessfully': 'অনুরোধ সফলভাবে বাতিল করা হয়েছে',
      'failedToCancelRequest': 'অনুরোধ বাতিল করতে ব্যর্থ',
      'pleaseEnterValidAmount': 'অনুগ্রহ করে একটি বৈধ পরিমাণ লিখুন',
      'pleaseEnterValidElectricityBill': 'অনুগ্রহ করে একটি বৈধ বিদ্যুৎ বিলের পরিমাণ লিখুন',
      'paymentNotificationSent': 'পেমেন্ট নোটিফিকেশন পাঠানো হয়েছে!',
      'pleaseSelectMonth': 'অনুগ্রহ করে একটি মাস নির্বাচন করুন',
      'advancePaymentRequestSentSuccessfully': 'অগ্রিম পেমেন্ট অনুরোধ সফলভাবে পাঠানো হয়েছে!',
      'failedToSendAdvancePaymentRequest': 'অগ্রিম পেমেন্ট অনুরোধ পাঠাতে ব্যর্থ',
      'advancePaymentRequestCancelledSuccessfully': 'অগ্রিম পেমেন্ট অনুরোধ সফলভাবে বাতিল করা হয়েছে!',
      'failedToCancelAdvancePaymentRequest': 'অগ্রিম পেমেন্ট অনুরোধ বাতিল করতে ব্যর্থ',
      'amount': 'টাকা',
      'month': 'মাস',
      'january': 'জানুয়ারি',
      'february': 'ফেব্রুয়ারি',
      'march': 'মার্চ',
      'april': 'এপ্রিল',
      'may': 'মে',
      'june': 'জুন',
      'july': 'জুলাই',
      'august': 'আগস্ট',
      'september': 'সেপ্টেম্বর',
      'october': 'অক্টোবর',
      'november': 'নভেম্বর',
      'december': 'ডিসেম্বর',
      'enterTenantName': 'ভাড়াটের নাম লিখুন',
      'enterPhoneNumber': 'ভাড়াটের ফোন নম্বর লিখুন',
      'errorPrefix': 'ত্রুটি',
      'enterPaymentAmount': 'পেমেন্ট পরিমাণ লিখুন',
      'enterElectricityBill': 'পরিশোধ করার বিদ্যুৎ বিলের পরিমাণ লিখুন',
      'cancelAdvancePaymentRequest': 'অগ্রিম পেমেন্ট অনুরোধ বাতিল করুন',
      'cancelAdvancePaymentRequestConfirmation': 'আপনি কি নিশ্চিত যে আপনি এই ফ্লোরের জন্য অপেক্ষমান অগ্রিম পেমেন্ট অনুরোধ বাতিল করতে চান?',
      'yesCancelAdvancePaymentRequest': 'হ্যাঁ, বাতিল করুন',
      'noFloorsAddedYet': 'এখনও কোন ফ্লোর যোগ করা হয়নি',
      'enterAdvancePaymentAmount': 'অগ্রিম পেমেন্টের পরিমাণ লিখুন',
  'showAdvanceDetails': 'অগ্রিম বিবরণ দেখুন',
  'advanceDetails': 'অগ্রিম বিবরণ',
  'noAdvanceDetails': 'এই ফ্লোরের জন্য কোন অগ্রিম বিবরণ পাওয়া যায়নি',
  'advanceDetailsForFloor': '{floorName} এর অগ্রিম বিবরণ',
      'currency' : 'টাকা/মাসে',
      "paymentDetails": "পেমেন্টের বিবরণ",
      "tk": "টাকা",
      "adjustDueRent": "বকেয়া ভাড়া সমন্বয় করুন",
      "subtract": "বিয়োগ করুন",
      "rentAmount": "ভাড়ার পরিমাণ",
      "electricityBillOptional": "বিদ্যুত বিল",
      "preview": "পূর্বরূপ দেখুন",
      "current": "বর্তমান",
      "newTotal": "নতুন মোট",
      "addAmount": "পরিমাণ যোগ করুন",
      "paymentHistory": "পেমেন্ট ইতিহাস",
      "date": "তারিখ",
      "newAddedRent": "নতুন যোগকৃত ভাড়া",
      "receivedMoney": "প্রাপ্ত অর্থ",
      "dueRent": "বকেয়া ভাড়া",
      "newAddedElectricityBill": "নতুন যোগকৃত বিদ্যুত বিল",
      "paidElectricity": "পরিশোধিত বিদ্যুত বিল",
      "dueElectricityBill": "বকেয়া বিদ্যুত বিল",
      "records": "রেকর্ডসমূহ",
      "noMonth": "কোন মাস নির্বাচিত হয়নি",
      "paymentRecord": "পেমেন্ট রেকর্ড",
      "pleaseEnterValidPositiveAmount": "অনুগ্রহ করে একটি সঠিক পজিটিভ পরিমাণ দিন",
      "dueRentCannotBeNegative": "বাকি ভাড়া ঋণাত্মক হতে পারে না",
      "increased": "বৃদ্ধি পেয়েছে",
      "decreased": "হ্রাস পেয়েছে",
      "successfully": "সফলভাবে",
      "subtractAmount": "পরিমাণ বিয়োগ করুন",
      "loadingNotifications": "বিজ্ঞপ্তি লোড হচ্ছে...",
      "noNotificationsYet": "এখনও কোন বিজ্ঞপ্তি নেই",
      "youllSeeYourNotificationsHere": "আপনি এখানে আপনার বিজ্ঞপ্তি দেখতে পাবেন",
      "accept": "গ্রহণ",
      "reject": "প্রত্যাখ্যান",
      "addComment": "মন্তব্য",
      "addYourComment": "আপনার মন্তব্য যোগ করুন",
      "pleaseEnterAComment": "অনুগ্রহ করে একটি মন্তব্য লিখুন",
      "commentSentSuccessfully": "মন্তব্য সফলভাবে পাঠানো হয়েছে",
      "failedToSendComment": "মন্তব্য পাঠাতে ব্যর্থ",
      "failedToProcessTheAction": "কর্ম প্রক্রিয়া করতে ব্যর্থ"
    },
  };
  
  String get(String key) {
    return _localizedValues[locale.languageCode]?[key] ?? 
           _localizedValues['en']![key] ?? 
           key;
  }
  
  String get appTitle => get('appTitle');
  String get managedProperties => get('managedProperties');
  String get tenantProperties => get('tenantProperties');
  String get addProperty => get('addProperty');
  String get propertyName => get('propertyName');
  String get propertyNameHint => get('propertyNameHint');
  String get address => get('address');
  String get addressHint => get('addressHint');
  String get propertyPhoto => get('propertyPhoto');
  String get addPhotoMessage => get('addPhotoMessage');
  String get camera => get('camera');
  String get gallery => get('gallery');
  String get remove => get('remove');
  String get change => get('change');
  String get cancel => get('cancel');
  String get confirm => get('confirm');
  String get addPropertySuccess => get('addPropertySuccess');
  String get fillRequiredFields => get('fillRequiredFields');
  String get failedToAddProperty => get('failedToAddProperty');
  String get noManagedProperties => get('noManagedProperties');
  String get noManagedPropertiesSubtitle => get('noManagedPropertiesSubtitle');
  String get noTenantProperties => get('noTenantProperties');
  String get noTenantPropertiesSubtitle => get('noTenantPropertiesSubtitle');
  String get loadingManagedProperties => get('loadingManagedProperties');
  String get loadingTenantProperties => get('loadingTenantProperties');
  String get somethingWentWrong => get('somethingWentWrong');
  String get tryAgain => get('tryAgain');
  String get managedProperty => get('managedProperty');
  String get tenantProperty => get('tenantProperty');
  String get notifications => get('notifications');
  String get language => get('language');
  String get logout => get('logout');
  String get logoutConfirmation => get('logoutConfirmation');
  String get selectLanguage => get('selectLanguage');
  String get english => get('english');
  String get bangla => get('bangla');
  String get manager => get('manager');
  String get tenant => get('tenant');
  String get failedToLoadProperties => get('failedToLoadProperties');
  String get add => get('add');
  String get addFloor => get('addFloor');
  String get addTenant => get('addTenant');
  String get update => get('update');
  String get send => get('send');
  String get sendRequest => get('sendRequest');
  String get sendPaymentNotification => get('sendPaymentNotification');
  String get sendAdvancePaymentRequest => get('sendAdvancePaymentRequest');
  String get yes => get('yes');
  String get no => get('no');
  String get retry => get('retry');
  String get addFirstFloor => get('addFirstFloor');
  String get occupied => get('occupied');
  String get pending => get('pending');
  String get available => get('available');
  String get tapToView => get('tapToView');
  String get requestPending => get('requestPending');
  String get rent => get('rent');
  String get unknownUser => get('unknownUser');
  String get floorName => get('floorName');
  String get floorNameHint => get('floorNameHint');
  String get monthlyRent => get('monthlyRent');
  String get rentHint => get('rentHint');
  String get tenantName => get('tenantName');
  String get tenantNameHint => get('tenantNameHint');
  String get phoneNumber => get('phoneNumber');
  String get phoneNumberHint => get('phoneNumberHint');
  String get removeTenant => get('removeTenant');
  String get removeTenantConfirmation => get('removeTenantConfirmation');
  String get updateFloor => get('updateFloor');
  String get sendTenantRequest => get('sendTenantRequest');
  String get cancelRequest => get('cancelRequest');
  String get cancelRequestConfirmation => get('cancelRequestConfirmation');
  String get paymentAmount => get('paymentAmount');
  String get paymentAmountHint => get('paymentAmountHint');
  String get electricityBill => get('electricityBill');
  String get electricityBillHint => get('electricityBillHint');
  String get advanceAmount => get('advanceAmount');
  String get advanceAmountHint => get('advanceAmountHint');
  String get selectMonth => get('selectMonth');
  String get cancelAdvancePayment => get('cancelAdvancePayment');
  String get cancelAdvancePaymentConfirmation => get('cancelAdvancePaymentConfirmation');
  String get yesCancel => get('yesCancel');
  String get noFloorsAdded => get('noFloorsAdded');
  String get fillAllFields => get('fillAllFields');
  String get invalidRentAmount => get('invalidRentAmount');
  String get failedToAddFloor => get('failedToAddFloor');
  String get failedToAddTenant => get('failedToAddTenant');
  String get failedToRemoveTenant => get('failedToRemoveTenant');
  String get tenantRemovedSuccessfully => get('tenantRemovedSuccessfully');
  String get failedToUpdateFloor => get('failedToUpdateFloor');
  String get pleaseEnterPhoneNumber => get('pleaseEnterPhoneNumber');
  String get tenantRequestSentSuccessfully => get('tenantRequestSentSuccessfully');
  String get requestCancelledSuccessfully => get('requestCancelledSuccessfully');
  String get failedToCancelRequest => get('failedToCancelRequest');
  String get pleaseEnterValidAmount => get('pleaseEnterValidAmount');
  String get pleaseEnterValidElectricityBill => get('pleaseEnterValidElectricityBill');
  String get paymentNotificationSent => get('paymentNotificationSent');
  String get pleaseSelectMonth => get('pleaseSelectMonth');
  String get advancePaymentRequestSentSuccessfully => get('advancePaymentRequestSentSuccessfully');
  String get failedToSendAdvancePaymentRequest => get('failedToSendAdvancePaymentRequest');
  String get advancePaymentRequestCancelledSuccessfully => get('advancePaymentRequestCancelledSuccessfully');
  String get failedToCancelAdvancePaymentRequest => get('failedToCancelAdvancePaymentRequest');
  String get amount => get('amount');
  String get month => get('month');
  String get january => get('january');
  String get february => get('february');
  String get march => get('march');
  String get april => get('april');
  String get may => get('may');
  String get june => get('june');
  String get july => get('july');
  String get august => get('august');
  String get september => get('september');
  String get october => get('october');
  String get november => get('november');
  String get december => get('december');
  String get enterTenantName => get('enterTenantName');
  String get enterPhoneNumber => get('enterPhoneNumber');
  String get errorPrefix => get('errorPrefix');
  String get enterPaymentAmount => get('enterPaymentAmount');
  String get enterElectricityBill => get('enterElectricityBill');
  String get cancelAdvancePaymentRequest => get('cancelAdvancePaymentRequest');
  String get cancelAdvancePaymentRequestConfirmation => get('cancelAdvancePaymentRequestConfirmation');
  String get yesCancelAdvancePaymentRequest => get('yesCancelAdvancePaymentRequest');
  String get noFloorsAddedYet => get('noFloorsAddedYet');
  String get monthJanuary => get('january');
  String get monthFebruary => get('february');
  String get monthMarch => get('march');
  String get monthApril => get('april');
  String get monthMay => get('may');
  String get monthJune => get('june');
  String get monthJuly => get('july');
  String get monthAugust => get('august');
  String get monthSeptember => get('september');
  String get monthOctober => get('october');
  String get monthNovember => get('november');
  String get monthDecember => get('december');
  String get enterAdvancePaymentAmount => get('enterAdvancePaymentAmount');
  String get showAdvanceDetails => get('showAdvanceDetails');
  String get advanceDetails => get('advanceDetails');
  String get noAdvanceDetails => get('noAdvanceDetails');
  String get advanceDetailsForFloor => get('advanceDetailsForFloor');
  String get currency => get('currency');
  String get paymentDetails => get('paymentDetails');
  String get tk => get('tk');
  String get adjustDueRent => get('adjustDueRent');
  String get subtract => get('subtract');
  String get rentAmount => get('rentAmount');
  String get electricityBillOptional => get('electricityBillOptional');
  String get preview => get('preview');
  String get current => get('current');
  String get newTotal => get('newTotal');
  String get addAmount => get('addAmount');
  String get paymentHistory => get('paymentHistory');
  String get date => get('date');
  String get newAddedRent => get('newAddedRent');
  String get receivedMoney => get('receivedMoney');
  String get records => get('records');
  String get noMonth => get('noMonth');
  String get paymentRecord => get('paymentRecord');
  String get newAddedElectricityBill => get('newAddedElectricityBill');
  String get paidElectricity => get('paidElectricity');
  String get dueElectricityBill => get('dueElectricityBill');
  String get dueRent => get('dueRent'); 
  String get pleaseEnterValidPositiveAmount => get('pleaseEnterValidPositiveAmount');
  String get dueRentCannotBeNegative => get('dueRentCannotBeNegative');
  String get increased => get('increased');
  String get decreased => get('decreased');
  String get successfully => get('successfully');
  String get subtractAmount => get('subtractAmount');
  String get loadingNotifications => get('loadingNotifications');
  String get noNotificationsYet => get('noNotificationsYet');
  String get youllSeeYourNotificationsHere => get('youllSeeYourNotificationsHere');
  String get accept => get('accept');
  String get reject => get('reject');
  String get addComment => get('addComment');
  String get addYourComment => get('addYourComment');
  String get pleaseEnterAComment => get('pleaseEnterAComment');
  String get commentSentSuccessfully => get('commentSentSuccessfully');
  String get failedToSendComment => get('failedToSendComment');
  String get failedToProcessTheAction => get('failedToProcessTheAction');
  
  String notificationAcceptedSuccessfully(String action) {
    if (locale.languageCode == 'bn') {
      return action == 'accepted' ? 'বিজ্ঞপ্তি সফলভাবে গ্রহণ করা হয়েছে' : 'বিজ্ঞপ্তি সফলভাবে প্রত্যাখ্যান করা হয়েছে';
    }
    return action == 'accepted' ? 'Notification accepted successfully' : 'Notification rejected successfully';
  }
  
  String error(String message) {
    if (locale.languageCode == 'bn') {
      return 'ত্রুটি: $message';
    }
    return 'Error: $message';
  }

}

class AppLocalizationsDelegate extends LocalizationsDelegate<AppLocalizations> {
  const AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    return ['en', 'bn'].contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) async {
    return AppLocalizations(locale);
  }

  @override
  bool shouldReload(AppLocalizationsDelegate old) => false;
} 