# 🎨 FamGo Color Migration - Quick Reference

## How to Use Green Colors

### In Widgets

**Instead of hardcoded colors:**
```dart
// ❌ OLD
FloatingActionButton(
  backgroundColor: Colors.red,
  child: Icon(Icons.add, color: Colors.white),
)

// ✅ NEW
FloatingActionButton(
  backgroundColor: FamGoColors.primary,
  child: Icon(Icons.add, color: FamGoColors.white),
)
```

**For themed buttons:**
```dart
ElevatedButton(
  style: ElevatedButton.styleFrom(
    backgroundColor: FamGoColors.primary, // Green
    foregroundColor: FamGoColors.white,   // White text
  ),
  onPressed: () {},
  child: Text('Confirm'),
)
```

**For text colors:**
```dart
Text(
  'Hello',
  style: TextStyle(
    color: FamGoColors.textDark,      // Dark grey text
    fontWeight: FontWeight.bold,
  ),
)
```

**For containers & cards:**
```dart
Container(
  color: FamGoColors.backgroundColor,  // Light grey background
  child: Card(
    color: FamGoColors.cardBackground, // White
    child: ...,
  ),
)
```

### Material Theme Integration

**In main.dart:**
```dart
MaterialApp(
  title: 'FamGo',
  theme: FamGoColors.getLightTheme(), // ✅ Applies green theme globally
  home: SplashScreen(),
)
```

This automatically makes:
- AppBars green
- Buttons green
- Form fields green on focus
- All Material widgets themed

## Color Reference

```
PRIMARY COLORS (Use these most)
├── FamGoColors.primary           → #2ECC71 (Main green)
├── FamGoColors.primaryDark       → #27AE60 (Darker green)
└── FamGoColors.primaryLight      → #A9DFBF (Lighter green)

STATUS COLORS (For specific states)
├── FamGoColors.success           → #2ECC71 (Green - Success)
├── FamGoColors.warning           → #F39C12 (Orange - Warning)
├── FamGoColors.error             → #E74C3C (Red - Error only)
└── FamGoColors.info              → #3498DB (Blue - Info)

TEXT COLORS (For text)
├── FamGoColors.textDark          → #2C3E50 (Dark grey)
├── FamGoColors.textLight         → #7F8C8D (Light grey)
├── FamGoColors.textGrey          → #BCDBC7 (Very light)
├── FamGoColors.white             → #FFFFFF (White)
└── FamGoColors.black             → #000000 (Black)

BACKGROUND COLORS (For backgrounds)
├── FamGoColors.backgroundColor   → #FAFAFA (Light grey)
├── FamGoColors.cardBackground    → #FFFFFF (White)
└── FamGoColors.dividerColor      → #ECF0F1 (Divider)
```

## Search & Replace Patterns

Use your IDE's Find & Replace to convert:

```
Find: Colors.red
Replace: FamGoColors.primary

Find: Color(0xFFDC143C)
Replace: FamGoColors.primary

Find: Colors.deepOrange
Replace: FamGoColors.primary

Find: Colors.white
Replace: FamGoColors.white

Find: Colors.grey\[300\]
Replace: FamGoColors.dividerColor
```

## Common Patterns

### App Bar
```dart
AppBar(
  backgroundColor: FamGoColors.primary,  // Green
  foregroundColor: FamGoColors.white,    // White text/icons
  title: Text('Title'),
)
```

### Buttons
```dart
// Primary action
ElevatedButton(
  style: ElevatedButton.styleFrom(
    backgroundColor: FamGoColors.primary,
  ),
  onPressed: () {},
  child: Text('Primary'),
)

// Secondary action
OutlinedButton(
  style: OutlinedButton.styleFrom(
    side: BorderSide(color: FamGoColors.primary),
  ),
  onPressed: () {},
  child: Text('Secondary'),
)

// Tertiary action
TextButton(
  onPressed: () {},
  child: Text('Tertiary'),
)
```

### Form Fields
```dart
TextField(
  decoration: InputDecoration(
    filled: true,
    fillColor: FamGoColors.backgroundColor,
    enabledBorder: OutlineInputBorder(
      borderSide: BorderSide(color: FamGoColors.dividerColor),
    ),
    focusedBorder: OutlineInputBorder(
      borderSide: BorderSide(
        color: FamGoColors.primary, // Green when focused
        width: 2,
      ),
    ),
  ),
)
```

### Badges & Status
```dart
Container(
  padding: EdgeInsets.symmetric(horizontal: 12, vertical: 6),
  decoration: BoxDecoration(
    color: FamGoColors.primary,     // Green background
    borderRadius: BorderRadius.circular(20),
  ),
  child: Text(
    'Active',
    style: TextStyle(color: FamGoColors.white),
  ),
)
```

### Loading Indicators
```dart
CircularProgressIndicator(
  valueColor: AlwaysStoppedAnimation<Color>(FamGoColors.primary),
)
```

### Cards with Borders
```dart
Card(
  shape: RoundedRectangleBorder(
    borderRadius: BorderRadius.circular(12),
    side: BorderSide(
      color: FamGoColors.primary,
      width: 2,
    ),
  ),
  child: ...,
)
```

## Transparency & Opacity

```dart
// 50% transparent green
FamGoColors.primary.withOpacity(0.5)

// Or use helper
FamGoColors.withOpacity(FamGoColors.primary, 0.5)

// Examples
Container(
  color: FamGoColors.primary.withOpacity(0.1), // Very light green
  child: ...,
)

Container(
  color: FamGoColors.primary.withOpacity(0.7), // Dark green
  child: ...,
)
```

## Gradients

```dart
// Use predefined gradient
Container(
  decoration: BoxDecoration(
    gradient: FamGoColors.primaryGradient,
  ),
  child: ...,
)

// Or create custom
Container(
  decoration: BoxDecoration(
    gradient: LinearGradient(
      colors: [
        FamGoColors.primary,
        FamGoColors.primaryDark,
      ],
    ),
  ),
)
```

## Ride Type Colors

```dart
// Different ride types have different colors
Container(
  color: FamGoColors.rideEconomy,   // Grey
)

Container(
  color: FamGoColors.rideStandard,  // Green (primary)
)

Container(
  color: FamGoColors.rideShare,     // Blue
)

Container(
  color: FamGoColors.ridePremium,   // Red
)
```

## Map Colors

```dart
// Current location (blue)
Marker(
  icon: BitmapDescriptor.defaultMarkerWithHue(
    BitmapDescriptor.hueBlue,
  ),
)

// Pickup location (green)
Marker(
  infoWindow: InfoWindow(title: 'Pickup'),
  icon: BitmapDescriptor.defaultMarkerWithHue(
    BitmapDescriptor.hueGreen, // Use mapMarkerPickup color
  ),
)

// Polyline route (green)
Polyline(
  polylineId: PolylineId('route'),
  color: FamGoColors.mapPolyline, // Green
  width: 5,
  points: ...,
)
```

## Testing Color Implementation

**Checklist:**
- [ ] All buttons are green (not red)
- [ ] All app bars are green
- [ ] Form focus borders are green
- [ ] Success states use green
- [ ] Error states still use red (only for errors)
- [ ] Text contrast is good
- [ ] Errors display properly
- [ ] On light backgrounds
- [ ] On dark backgrounds (if applicable)
- [ ] On map overlays

## Import Statement

Always add at top of file:
```dart
import 'package:famgo_passenger_app/core/app_colors.dart';
```

## Folder Structure

```
lib/
├── core/
│   └── app_colors.dart ✅ (Centralized colors)
├── authentication/
│   ├── otp_screen.dart (Update colors)
│   ├── register_screen.dart (Update colors)
│   └── user_information_screen.dart (Update colors)
├── pages/
│   └── home_page.dart (Update colors + layout fix)
├── screens/
│   └── splash_screen.dart (Update colors)
└── ...
```

---

**Quick Start:** Replace `Colors.red` with `FamGoColors.primary` everywhere!
