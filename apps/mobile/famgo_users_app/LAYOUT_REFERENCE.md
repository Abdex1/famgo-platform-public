# 📍 Visual Layout Reference - Fixed Overflow

## Layout Structure Comparison

### ❌ BEFORE (Broken - Bottom Overflow)

```
┌─────────────────────────────────┐
│  AppBar: "Set Dropoff Location" │ ← Fixed height
├─────────────────────────────────┤
│  SingleChildScrollView {        │
│    Column {                     │
│      ┌───────────────────────┐  │
│      │ Location Input Card   │  │
│      │ [Pickup] [Delete]     │  │
│      │ [Destination] [Delete]│  │
│      └───────────────────────┘  │
│                                 │
│      Loading Indicator          │
│      Predictions List           │
│      Error Messages             │
│      Manual Entry Prompt  ❌ CUT│  ← OVERFLOW! (26px)
│    }                            │
│  }                              │
└─────────────────────────────────┘
         Screen Bottom
```

**Problem:** Column expands beyond screen, last items get clipped

---

### ✅ AFTER (Fixed - No Overflow)

```
┌─────────────────────────────────┐
│  AppBar: "Set Pickup & Dropoff" │ ← Fixed height
├─────────────────────────────────┤
│  ┌───────────────────────────┐   │
│  │ Location Input Card       │   │ ← FIXED (Not scrollable)
│  │ [Pickup] [Delete]         │   │
│  │ [Destination] [Delete]    │   │
│  └───────────────────────────┘   │
├─────────────────────────────────┤ ← Divider
│                                 │
│  Expanded {                     │
│    SingleChildScrollView {      │
│      Column {                   │
│        Loading Indicator        │
│        Predictions List         │
│        Error Messages           │
│        Manual Entry Prompt ✅ OK│ ← Fully visible!
│        SizedBox(height: 20) ✅  │ ← Bottom spacing
│      }                          │
│    }                            │
│  }                              │
│                                 │
└─────────────────────────────────┘
         Screen Bottom ✅ Perfect!
```

**Solution:** Separated fixed card from scrollable content

---

## Location Input Card Layout

```
    ┌─ Card (Elevation: 2) ─────────────────────┐
    │                                            │
    │  Column(mainAxisSize: MainAxisSize.min)   │
    │  ├─ Row                                   │
    │  │  ├─ [📍] (Icon)                        │
    │  │  └─ TextField (Editable)               │
    │  │     ├─ Pickup Address                  │
    │  │     ├─ Hint: "Pickup Address"          │
    │  │     └─ Clear button (✕)                │
    │  │                                         │
    │  ├─ SizedBox(height: 14)                  │
    │  ├─ Divider                               │
    │  ├─ SizedBox(height: 14)                  │
    │  │                                         │
    │  └─ Row                                   │
    │     ├─ [📍] (Icon)                        │
    │     └─ TextField (Editable)               │
    │        ├─ Destination Address             │
    │        ├─ Hint: "Destination Address"     │
    │        └─ Clear button (✕)                │
    │                                            │
    └────────────────────────────────────────────┘
         Padding: All 16px
         Margin: Horizontal 12px, Vertical 12px
         Border Radius: 12px
```

---

## Search Results Display

### When API Works (Predictions Available)

```
Location Input Card
├─── Divider
├─ Expanded
│  └─ SingleChildScrollView
│     └─ Column
│        ├─ [📍] Prediction 1
│        │   "Gulele, Addis Ababa"
│        ├─ Divider
│        ├─ [📍] Prediction 2
│        │   "Guleletin Street, Addis Ababa"
│        ├─ Divider
│        ├─ [📍] Prediction 3
│        │   "Gullele Botanical Garden"
│        └─ SizedBox(height: 20) ✅
```

### When API Fails (Manual Entry Fallback)

```
Location Input Card
├─── Divider
├─ Expanded
│  └─ SingleChildScrollView
│     └─ Column
│        ├─ ⚠️ Container (Amber bg, amber border)
│        │  ├─ [ℹ️] Icon
│        │  ├─ "Location search unavailable"
│        │  ├─ "Enter your location manually."
│        │  └─ [Confirm Address] Button
│        │
│        └─ SizedBox(height: 20) ✅
```

---

## Scrollable Area Behavior

### Small Screen (4.5" - 480px height)

```
AppBar: 56px (fixed)
Card: 140px (fixed)
Divider: 1px
────────────────────
Available for scroll: 480 - 56 - 140 - 1 = 283px

Scrollable Content:
├─ Loading: 40px
├─ Predictions: 
│  ├─ Prediction 1: 50px
│  ├─ Divider: 8px
│  ├─ Prediction 2: 50px
│  └─ ... (repeats)
└─ SizedBox: 20px ✅

Result: Content scrolls smoothly, last item visible ✅
```

### Large Screen (6.5" - 720px height)

```
AppBar: 56px (fixed)
Card: 140px (fixed)
Divider: 1px
────────────────────
Available for scroll: 720 - 56 - 140 - 1 = 523px

Result: All content fits without scroll (even better!) ✅
```

---

## State Management Diagram

### Pickup Location State

```
┌─ Text: "" ────────────────────────────────┐
│                                           │
└─ User types "My Home" ────────────────────┘
                ↓
        ┌─ _onPickUpChanged() ──────────┐
        │ - Clear debounce timer        │
        │ - Check length > 2            │
        │ - Start new debounce (600ms)  │
        └──────────────────────────────┘
                ↓
    ┌─ _searchPickUpLocation() ─────────┐
    │ - Check if billing error occurred │
    │ - Make API call                   │
    └──────────────────────────────────┘
                ↓
        ┌─ API Response ───────────────────────┐
        │                                      │
        ├─ Status: "OK"                       │
        │  └─ _handlePickUpSuccessfulSearch() │
        │     └─ Show predictions             │
        │                                      │
        ├─ Status: "REQUEST_DENIED"           │
        │  └─ _handlePickUpBillingError()     │
        │     ├─ Set _pickUpBillingErrorOccurred = true
        │     └─ Show manual entry prompt     │
        │                                      │
        └─ Other Status ──────────────────────┘
           └─ Retry or show error

User Confirms Manual Entry
        ↓
    ┌─ _confirmManualPickUp() ──────┐
    │ - Get text from controller    │
    │ - Create AddressModel         │
    │ - Save to AppInfo             │
    │ - Show snackbar               │
    │ - requestFocus(destination)   │
    └───────────────────────────────┘
        ↓
    Focus moves to Destination
```

### Destination Location State

```
Similar to Pickup, but:
- _activeField = 'destination'
- Uses destinationTextEditingController
- Uses dropOffPredictionsPlacesList
- Uses _billingErrorOccurred (different flag)
- On confirm: Navigator.pop(context, "placeSelected")
```

---

## Manual Entry Prompt UI

```
┌─────────────────────────────────────┐
│  Container(amber[50] bg)            │
│  ┌─────────────────────────────────┐│
│  │ ℹ️ Icon (amber[700])              ││
│  │                                   ││
│  │ "Location search unavailable"    ││
│  │ (Font: Bold, amber[900])         ││
│  │                                   ││
│  │ "Enter your location manually."  ││
│  │ (Font: Regular, amber[700])      ││
│  │ (Align: Center)                  ││
│  │                                   ││
│  │ ┌─ TextField ──────────────────┐ ││
│  │ │ "Enter location address"     │ ││
│  │ │ (gray[100] bg, gray border)  │ ││
│  │ └──────────────────────────────┘ ││
│  │                                   ││
│  │ ┌──────────────────────────────┐ ││
│  │ │ [Confirm Location] Button    │ ││
│  │ │ (amber[700] bg)              │ ││
│  │ │ (Full width, 44px height)    │ ││
│  │ └──────────────────────────────┘ ││
│  │                                   ││
│  └─────────────────────────────────┘│
│  Border: amber[300]                 │
│  Radius: 8px                        │
│  Padding: All 16px                  │
└─────────────────────────────────────┘
  Margin: Horizontal 16px, Vertical 16px
```

---

## Focus Flow Diagram

```
Screen Loads
    ↓
┌─ Focus on Destination (default)
│  destinationFocusNode.requestFocus()
│
├─ User taps Pickup field
│  pickUpFocusNode receives focus
│  _activeField = 'pickup'
│  Pickup predictions shown
│
├─ User selects pickup prediction (or confirms manual)
│  pickUpTextEditingController.text = selected address
│  AppInfo.pickUpLocation = address
│  destinationFocusNode.requestFocus() ✅
│
└─ Focus on Destination
   _activeField = 'destination'
   Destination predictions shown
   User types destination
   User confirms
   Navigator.pop() → Screen closes
```

---

## Error Handling Flow

```
┌─ API Call Made ──────────────────────┐
│                                      │
└─ Response Received ──────────────────┘
         ↓
    ┌─ Parse Response ────────────────┐
    │ - Extract 'status' field        │
    │ - Check for null                │
    └─────────────────────────────────┘
         ↓
    ┌─ Status Check ──────────────────────────────┐
    │                                             │
    ├─ "OK" ────→ Show predictions              │
    │                                             │
    ├─ "ZERO_RESULTS" ────→ Clear list          │
    │                                             │
    ├─ "REQUEST_DENIED" ────→ Show manual entry │
    │                          (billing issue)    │
    │                                             │
    ├─ "OVER_QUERY_LIMIT" ────→ Retry once      │
    │                           Then show error   │
    │                                             │
    ├─ "INVALID_REQUEST" ────→ Show error       │
    │                                             │
    └─ Other ────→ Retry once, then show error  │
    
Each error wrapped in: if (!mounted) return;
```

---

## Widget Tree Summary

```
Scaffold
├── appBar: AppBar
│   ├── title: "Set Pickup & Dropoff"
│   └── leading: Back button
│
└── body: Column
    ├── LocationInputCard (Card)
    │   └── Column (min size)
    │       ├── Pickup TextField
    │       ├── Divider
    │       └── Destination TextField
    │
    ├── Divider
    │
    └── Expanded
        └── SingleChildScrollView
            └── Column
                ├── LoadingIndicator (if _isLoading)
                ├── ErrorMessage (if _errorMessage)
                ├── PickupManualPrompt (if _pickUpBillingErrorOccurred)
                ├── PickupPredictions (if predictions available)
                ├── DestinationManualPrompt (if _billingErrorOccurred)
                ├── DestinationPredictions (if predictions available)
                └── SizedBox(height: 20) ✅
```

---

**All layouts properly constrained. No overflow possible! ✅**
