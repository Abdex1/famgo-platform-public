# Manual Location Entry - Implementation Details

## How It Works (Step-by-Step)

### Scenario: Cloud Billing Not Enabled

**User Interaction Flow:**

```
1. User opens Search Destinations screen
   ↓
2. User types destination (e.g., "Addis Ababa")
   ↓
3. App makes API call to Google Places API
   ↓
4. API Response: {"status": "REQUEST_DENIED", "error_message": "Billing not enabled"}
   ↓
5. _handleBillingError() triggered:
   - Sets _billingErrorOccurred = true
   - Clears predictions list
   - Does NOT retry (max retries set to 0)
   ↓
6. Manual Entry Prompt appears with:
   - Warning message (amber color)
   - Text input field
   - "Confirm Address" button
   ↓
7. User enters address manually (e.g., "123 Main Street")
   ↓
8. User taps "Confirm Address"
   ↓
9. Address saved to app state + Navigator pops
```

---

## Code Walkthrough

### 1. Detecting Billing Error

```dart
Future<void> _searchLocation(String locationName, {bool isRetry = false}) async {
  // If billing error already occurred, skip API and show manual entry UI
  if (_billingErrorOccurred) {
    _showManualEntryFallback();
    return;
  }

  // Make API call...
  final status = responseFromPlacesAPI['status'] as String?;
  
  if (status == "REQUEST_DENIED") {
    _handleBillingError(errorMessage, isRetry);
  }
}
```

### 2. Handling Billing Error

```dart
void _handleBillingError(String? errorMessage, bool isRetry) {
  debugPrint('❌ API Request Denied (Billing Issue): $errorMessage');
  
  // Mark billing error and show manual entry fallback immediately
  if (mounted) {
    setState(() {
      _billingErrorOccurred = true;
      _errorMessage = '';  // Clear any previous errors
      dropOffPredictionsPlacesList = [];  // Clear predictions
      _retryCount = _maxRetries;  // Prevent further retries
    });
  }
}
```

### 3. Manual Address Confirmation

```dart
void _confirmManualAddress() {
  final destination = destinationTextEditingController.text.trim();
  
  if (destination.isEmpty) {
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Please enter a destination')),
    );
    return;
  }

  try {
    final appInfo = Provider.of<AppInfoClass>(context, listen: false);
    
    // Create manual address object
    final manualAddress = AddressModel(
      humanReadableAddress: destination,
      placeID: 'manual_${DateTime.now().millisecondsSinceEpoch}',
      latitudePosition: 0.0,
      longitudePosition: 0.0,
    );
    
    // Save to app state
    appInfo.dropOffLocation = manualAddress;
    
    // Show confirmation
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text('✓ Destination: $destination'),
        duration: const Duration(seconds: 2),
        backgroundColor: Colors.green[700],
      ),
    );
    
    // Return to previous screen after confirmation
    Future.delayed(const Duration(milliseconds: 600), () {
      if (mounted) {
        Navigator.pop(context, "placeSelected");
      }
    });
  } catch (e) {
    debugPrint('Error confirming manual address: $e');
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Error: ${e.toString()}')),
    );
  }
}
```

### 4. Conditional UI Rendering

```dart
@override
Widget build(BuildContext context) {
  return Scaffold(
    appBar: _buildAppBar(),
    body: SafeArea(
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.only(bottom: 16),
          child: Column(
            children: [
              _buildLocationInputCard(),
              if (_isLoading) _buildLoadingIndicator(),
              if (_errorMessage.isNotEmpty) _buildErrorMessage(),
              
              // Show manual entry prompt when:
              // 1. No predictions found AND no error, OR
              // 2. Billing error occurred
              if ((dropOffPredictionsPlacesList.isEmpty && !_isLoading && _errorMessage.isEmpty) || _billingErrorOccurred)
                _buildManualEntryPrompt(),
              
              // Show API predictions only if NO billing error
              if (dropOffPredictionsPlacesList.isNotEmpty && !_billingErrorOccurred)
                _buildPredictionsList(),
            ],
          ),
        ),
      ),
    ),
  );
}
```

### 5. Enhanced Manual Entry Prompt Widget

```dart
Widget _buildManualEntryPrompt() {
  return Padding(
    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
    child: Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        // Color based on whether billing error occurred
        color: _billingErrorOccurred ? Colors.amber[50] : Colors.blue[50],
        border: Border.all(
          color: _billingErrorOccurred ? Colors.amber[300]! : Colors.blue[300]!
        ),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            _billingErrorOccurred ? Icons.info_outline : Icons.edit_location_alt,
            color: _billingErrorOccurred ? Colors.amber[700] : Colors.blue[700],
            size: 28,
          ),
          const SizedBox(height: 12),
          
          // Dynamic title
          Text(
            _billingErrorOccurred 
                ? 'Location search unavailable' 
                : 'Can\'t find your destination?',
            style: TextStyle(
              color: _billingErrorOccurred ? Colors.amber[900] : Colors.blue[900],
              fontSize: 14,
              fontWeight: FontWeight.w600,
            ),
          ),
          
          const SizedBox(height: 8),
          
          // Dynamic message
          Text(
            _billingErrorOccurred
                ? 'The location service is temporarily offline. No problem! Enter your destination manually below.'
                : 'You can enter your address manually below and continue.',
            textAlign: TextAlign.center,
            style: TextStyle(
              color: _billingErrorOccurred ? Colors.amber[700] : Colors.blue[700],
              fontSize: 12,
            ),
          ),
          
          const SizedBox(height: 16),
          
          // Text input field
          Container(
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.grey[300]!),
            ),
            child: TextField(
              controller: destinationTextEditingController,
              enabled: _billingErrorOccurred,
              decoration: InputDecoration(
                hintText: 'Enter destination address',
                hintStyle: TextStyle(color: Colors.grey[500]),
                border: InputBorder.none,
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 12, 
                  vertical: 12
                ),
                suffixIcon: destinationTextEditingController.text.isNotEmpty
                    ? GestureDetector(
                        onTap: () {
                          destinationTextEditingController.clear();
                          setState(() {});
                        },
                        child: Icon(
                          Icons.clear, 
                          color: Colors.grey[600], 
                          size: 20
                        ),
                      )
                    : null,
              ),
              onChanged: (value) {
                setState(() {}); // Update UI for clear button
              },
            ),
          ),
          
          const SizedBox(height: 12),
          
          // Confirm button
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: _billingErrorOccurred 
                  ? Colors.amber[700] 
                  : Colors.blue[700],
              padding: const EdgeInsets.symmetric(
                horizontal: 24, 
                vertical: 12
              ),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(6),
              ),
              minimumSize: const Size(double.infinity, 44),
            ),
            onPressed: _confirmManualAddress,
            child: const Text(
              'Confirm Address',
              style: TextStyle(
                color: Colors.white,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    ),
  );
}
```

---

## Safety Features Implemented

✅ **Mounted Check**: All context operations wrapped with `if (mounted)` to prevent crashes

✅ **Error Handling**: Try-catch around manual address creation and storage

✅ **No Network Retries**: After billing error detected, API is not called again

✅ **User Feedback**: Clear snackbars for success/error states

✅ **Validation**: Empty address check before confirmation

✅ **State Management**: _billingErrorOccurred flag prevents UI state conflicts

✅ **Graceful Fallback**: No app crash, just manual entry option

---

## Testing Manual Entry

### Test Case 1: Billing Disabled
```
1. Disable billing in Google Cloud Console
2. Open app → Search Destinations
3. Type "Addis Ababa"
4. See: Manual entry prompt with amber warning
5. Enter: "123 Main Street, Addis Ababa"
6. Tap: "Confirm Address"
7. Result: Address saved, screen closes ✅
```

### Test Case 2: Billing Enabled
```
1. Enable billing in Google Cloud Console
2. Open app → Search Destinations
3. Type "Addis Ababa"
4. See: API predictions list (Places API working) ✅
5. Manual entry prompt NOT shown
6. Tap any prediction → address selected ✅
```

### Test Case 3: Mixed Scenario
```
1. Billing enabled, search works normally
2. Disable billing in Google Cloud
3. Search again → Manual entry prompt appears ✅
4. Works without app restart
```

---

## AddressModel Structure (For Reference)

```dart
class AddressModel {
  String? humanReadableAddress;      // "123 Main St, Addis Ababa"
  double? latitudePosition;           // 0.0 for manual entries
  double? longitudePosition;          // 0.0 for manual entries
  String? placeID;                    // 'manual_1706123456789'
  String? placeName;                  // Optional, can be null

  AddressModel({
    this.humanReadableAddress,
    this.latitudePosition,
    this.longitudePosition,
    this.placeID,
    this.placeName
  });
}
```

---

## Integration with Trip Flow

When user confirms manual address:

```
1. Address stored in AppInfoClass.dropOffLocation
2. Navigator.pop(context, "placeSelected")
3. Previous screen receives "placeSelected" return value
4. Trip creation proceeds with:
   - Pickup: From device location
   - Dropoff: From manual entry
   - Coordinates: (0.0, 0.0) for manual address*
   
*Note: When billing is restored, user can update address
      with API predictions which include coordinates
```

---

## Future Enhancement (When Billing Ready)

When Google Cloud Billing is enabled:

1. Remove `_billingErrorOccurred` checks
2. Users will get instant autocomplete predictions
3. Coordinates (latitude/longitude) will be accurate
4. Trip matching with drivers becomes more precise

Manual entry fallback remains as a safety net! ✅
