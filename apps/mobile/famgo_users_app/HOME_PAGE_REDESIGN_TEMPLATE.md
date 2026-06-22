# Home Page Redesign Implementation Template

This file contains the exact code structure for fixing home_page.dart bottom overflow and adding map zoom controls.

## Structure Overview

```
Scaffold
├── AppBar (Green theme)
├── Body: Stack
│   ├── GoogleMap (with custom zoom controls)
│   │   └── Positioned: Zoom +/- buttons (right side)
│   └── DraggableScrollableSheet (ride options)
│       └── ListView
│           ├── Handle bar
│           ├── Ride options cards
│           ├── Confirm button
│           └── Bottom padding (NO OVERFLOW!)
```

## Code Sections to Update

### 1. Imports (Add to top of home_page.dart)
```dart
import 'package:famgo_passenger_app/core/app_colors.dart';
```

### 2. Replace build() method body with Stack-based layout
**OLD:** SingleChildScrollView or Column with overflow
**NEW:** Stack with GoogleMap + DraggableScrollableSheet

### 3. Add GoogleMapController field
```dart
late GoogleMapController _mapController;
```

### 4. Zoom Control Buttons (Positioned widget)
Place in Stack as child, positioned to right side

### 5. DraggableScrollableSheet
- initialChildSize: 0.35
- minChildSize: 0.2
- maxChildSize: 0.9
- Build ListView inside (never Column!)

### 6. Bottom Padding
- Add SizedBox(height: 20) at end of ListView
- No overflow possible with ListView + padding

## Key Points

✅ **ListView prevents overflow** (built-in scrolling)
✅ **DraggableScrollableSheet handles drag gestures**
✅ **Positioned widgets for zoom controls**
✅ **Proper spacing with SizedBox + padding**
✅ **Green colors from FamGoColors**

## Testing Checklist

- [ ] No overflow on 4.5" screen
- [ ] No overflow on 6.5" screen
- [ ] Zoom + button works
- [ ] Zoom - button works
- [ ] Ride options sheet draggable
- [ ] Ride options scrollable when expanded
- [ ] Confirm button visible
- [ ] All text uses FamGoColors (green theme)

## Migration Path

1. Backup current home_page.dart
2. Replace build() with new Stack layout
3. Update all color references to FamGoColors
4. Test on device
5. Adjust initialChildSize if needed (0.3-0.4 range typically works)
6. Deploy

---

Exact code ready in COMPLETE_REDESIGN_GUIDE.md section "Phase 1"
