# Firebase Security Rules - Quick Reference & Troubleshooting

## 🚀 Quick Start

### 1. Deploy Rules (3 steps)
```bash
# Option A: Automatic (Recommended)
chmod +x deploy_firebase_rules.sh
./deploy_firebase_rules.sh

# Option B: Manual
# 1. Go to Firebase Console → Database → Rules
# 2. Copy content from firebase_rules_production.json
# 3. Paste and Publish
```

### 2. Create Admin User
In Firebase Console, add to database manually:
```
/admins/{your_uid} = true
```

### 3. Initialize Metadata
```
/admins/{your_uid} = true
/_metadata/blockList = {}
/_metadata/stats = {totalUsers: 0, totalDrivers: 0, totalTrips: 0}
```

---

## 📋 Rules Summary Table

| Collection | Read | Write | Key Fields | Notes |
|-----------|------|-------|-----------|-------|
| `users/$uid` | Self only | Self only | id, name, email, phone | Email immutable |
| `drivers/$uid` | Self or if available | Self only | license, vehicle | Location tracking |
| `onlineDrivers` | Public | Drivers only | g (geohash), l (coordinates) | GeoFire indexes |
| `tripRequest/$id` | Involved parties | Creator or involved | userID, status, fare | Status flow validated |
| `ratings/$id` | Public | Rating user only | userId, driverId, rating | 1-5 scale |
| `payments/$id` | Self/admin | Creator/admin | userId, amount, status | Immutable audit |
| `admins/$uid` | Admins only | None | boolean | Hard-coded |
| `_metadata` | Admins only | Admins only | blockList, stats | System config |

---

## 🔒 Field Validation Rules

### Email Validation
```
Pattern: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
Max Length: 254 characters
Immutable: After creation
```

### Phone Validation
```
Pattern: ^\+?[0-9]{10,20}$ (E.164 format)
Min Length: 10 digits
Max Length: 20 digits
Examples: +12025551234, +917898765432
```

### Coordinates Validation
```
Latitude:  -90 to 90
Longitude: -180 to 180
Type: Number (not string)
Example: {latitude: 28.7041, longitude: 77.1025}
```

### Enum Fields Validation
| Field | Allowed Values | Example |
|-------|---|---------|
| blockStatus | 'yes', 'no' | 'no' |
| userType | 'passenger', 'driver' | 'passenger' |
| vehicleType | 'car', 'auto', 'bike' | 'car' |
| status | 'new', 'accepted', 'arrived', 'ontrip', 'ended', 'cancelled' | 'ontrip' |
| paymentMethod | 'cash', 'card', 'wallet' | 'cash' |
| verificationStatus | 'verified', 'pending', 'rejected' | 'pending' |

---

## ⚠️ Common Errors & Solutions

### Error 1: "PERMISSION_DENIED"
**Cause**: User not authenticated or lacks permission
**Solution**:
```dart
// Ensure user is logged in
if (FirebaseAuth.instance.currentUser == null) {
  // Show login screen
}

// Ensure correct user ID
String uid = FirebaseAuth.instance.currentUser!.uid;
db.ref('users/$uid').set(data);  // ← Must match auth.uid
```

### Error 2: "Validation Error: Field 'email' invalid"
**Cause**: Invalid email format
**Solution**:
```dart
// Validate before writing
String email = "user@example.com";  // ✓ Valid
// String email = "invalid";  // ✗ Invalid

final emailRegex = RegExp(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');
if (!emailRegex.hasMatch(email)) {
  showError("Invalid email format");
  return;
}

db.ref('users/$uid').update({'email': email});
```

### Error 3: "Validation Error: Field 'phone' invalid"
**Cause**: Phone number not in E.164 format
**Solution**:
```dart
// Format phone numbers correctly
String phone = "+12025551234";  // ✓ Valid (E.164)
// String phone = "2025551234";  // ✗ Missing country code
// String phone = "+1 202-555-1234";  // ✗ Invalid characters

db.ref('users/$uid').update({'phone': phone});
```

### Error 4: "Validation Error: Field 'latitude' invalid"
**Cause**: Latitude out of range or wrong type
**Solution**:
```dart
// Latitude must be -90 to 90
double lat = 28.7041;  // ✓ Valid
// double lat = 100;  // ✗ Out of range
// String lat = "28.7041";  // ✗ Wrong type

db.ref('drivers/$uid').update({
  'latitude': lat,
  'longitude': 77.1025
});
```

### Error 5: "Validation Error: Field 'status' invalid"
**Cause**: Status not in allowed enum values
**Solution**:
```dart
// Only these status values allowed
String status = "ontrip";  // ✓ Valid
// String status = "in_progress";  // ✗ Not allowed

// Valid statuses:
const validStatuses = ['new', 'accepted', 'arrived', 'ontrip', 'ended', 'cancelled'];

db.ref('tripRequest/$tripId').update({'status': status});
```

### Error 6: "User blocked - Cannot write"
**Cause**: User is in block list
**Solution**:
```dart
// Only admin can remove from block list
// In Firebase Console:
// 1. Go to _metadata → blockList
// 2. Find and delete the user's entry

// Or programmatically (admin only):
// admin.database().ref('_metadata/blockList/{userId}').remove();
```

### Error 7: "Cannot write to immutable field 'email'"
**Cause**: Trying to modify email after creation
**Solution**:
```dart
// Email is immutable - set only during user creation
db.ref('users/$uid').set({
  'id': uid,
  'email': 'user@example.com',  // ✓ Set on creation
  'phone': '+12025551234',
  'name': 'John Doe',
  'createdAt': DateTime.now().millisecondsSinceEpoch
});

// Later: Cannot update email
db.ref('users/$uid').update({'email': 'new@example.com'});  // ✗ DENIED
```

### Error 8: "No index defined - Query requires index"
**Cause**: Querying on field without index
**Solution**:
```dart
// For Realtime Database, indexes are predefined in rules
// These fields already have indexes:
// - users: email, phone, blockStatus, createdAt
// - drivers: isAvailable, rating, latitude, longitude
// - tripRequest: userID, driverId, status, createdAt
// - ratings: driverId, userId, createdAt
// - payments: userId, driverId, status, createdAt

// Valid query:
db.ref('users').orderByChild('email').equalTo('user@example.com').once();
```

---

## 🔐 Security Checklist

### Before Going to Production

- [ ] All validation rules in place
- [ ] Admin users created and tested
- [ ] Block list structure initialized
- [ ] Backup created
- [ ] Client app validates data before writing
- [ ] Error handling implemented
- [ ] Rate limiting on backend (not in rules)
- [ ] Admin access controlled
- [ ] Sensitive fields immutable
- [ ] All enum values tested
- [ ] Coordinate validation tested

### After Deployment

- [ ] Monitor Firebase logs for errors
- [ ] Check for permission denied errors
- [ ] Verify validation failures
- [ ] Test with different user types
- [ ] Test block list functionality
- [ ] Test admin operations
- [ ] Test edge cases (invalid coords, etc)

---

## 📊 Testing Database Writes

### Test Script for Flutter
```dart
Future<void> testDatabaseRules() async {
  final db = FirebaseDatabase.instance.ref();
  final uid = FirebaseAuth.instance.currentUser!.uid;
  
  // Test 1: Create user (should work)
  try {
    await db.child('users/$uid').set({
      'id': uid,
      'name': 'Test User',
      'email': 'test@example.com',
      'phone': '+12025551234',
      'blockStatus': 'no',
      'createdAt': DateTime.now().millisecondsSinceEpoch,
    });
    print('✓ User creation: PASSED');
  } catch (e) {
    print('✗ User creation failed: $e');
  }
  
  // Test 2: Read own user (should work)
  try {
    final snap = await db.child('users/$uid').get();
    if (snap.exists) {
      print('✓ User read: PASSED');
    }
  } catch (e) {
    print('✗ User read failed: $e');
  }
  
  // Test 3: Invalid email (should fail)
  try {
    await db.child('users/$uid').update({
      'email': 'invalid'  // Invalid format
    });
    print('✗ Invalid email validation: FAILED (should have been rejected)');
  } catch (e) {
    print('✓ Invalid email validation: PASSED (correctly rejected)');
  }
  
  // Test 4: Invalid coordinates (should fail)
  try {
    await db.child('drivers/$uid').update({
      'latitude': 150  // Out of range
    });
    print('✗ Coordinate validation: FAILED (should have been rejected)');
  } catch (e) {
    print('✓ Coordinate validation: PASSED (correctly rejected)');
  }
}
```

---

## 🛠️ Advanced Administration

### Block a User (Admin Only)
```dart
Future<void> blockUser(String userId) async {
  final db = FirebaseDatabase.instance.ref();
  
  // Mark user as blocked in users collection
  await db.child('users/$userId').update({
    'blockStatus': 'yes'
  });
  
  // Add to block list for real-time prevention
  await db.child('_metadata/blockList/$userId').set(true);
  
  print('User $userId blocked successfully');
}
```

### Unblock a User (Admin Only)
```dart
Future<void> unblockUser(String userId) async {
  final db = FirebaseDatabase.instance.ref();
  
  // Update user status
  await db.child('users/$userId').update({
    'blockStatus': 'no'
  });
  
  // Remove from block list
  await db.child('_metadata/blockList/$userId').remove();
  
  print('User $userId unblocked successfully');
}
```

### Update User Statistics
```dart
Future<void> updateStats() async {
  final db = FirebaseDatabase.instance.ref();
  
  // This should be done by backend periodically
  // Counting is expensive on Realtime Database
  
  final users = await db.child('users').get();
  final drivers = await db.child('drivers').get();
  final trips = await db.child('tripRequest').get();
  
  await db.child('_metadata/stats').update({
    'totalUsers': users.children.length,
    'totalDrivers': drivers.children.length,
    'totalTrips': trips.children.length,
  });
}
```

---

## 📞 Support & Escalation

### If Issues Persist:

1. **Check Firebase Console**
   - Real-time Database → Rules (view current rules)
   - Logs (check for validation errors)

2. **Review Error Messages**
   - Copy exact error message
   - Check against solutions above

3. **Contact Firebase Support**
   - Firebase Console → Support → Create issue
   - Include project ID, error code, reproduction steps

4. **Community Resources**
   - [Stack Overflow](https://stackoverflow.com/questions/tagged/firebase-realtime-database)
   - [Firebase Slack Community](https://firebase.community)
   - [Firebase Documentation](https://firebase.google.com/docs/rules)

---

**Last Updated**: 2024
**Current Version**: 1.0.0
**Status**: Production Ready
