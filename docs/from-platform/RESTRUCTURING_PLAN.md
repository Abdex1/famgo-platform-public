# 🎯 FamGo Platform - Comprehensive Restructuring Plan

## PHASE 1: Analysis & Planning

### Current Structure Issues
1. **Mobile Apps**: Duplicate files in both `/presentation` and `/features` directories
2. **Shared Library**: Located in `apps/flutter-mobile/shared-flutter-lib` (should be centralized)
3. **Backend**: Scattered across `backend/shared`, `backend/api-gateway/kong`, and `gateway`
4. **Imports**: Path conflicts due to reorganization
5. **Duplication**: Services, models, widgets split between presentation and features layers

### Target Structure
```
C:\dev\FamGo-platform\
├── apps/
│   ├── flutter-mobile/
│   │   ├── passenger-app/
│   │   │   ├── lib/
│   │   │   │   ├── core/          # Shared across passenger app
│   │   │   │   ├── features/      # Feature modules
│   │   │   │   ├── main.dart
│   │   │   │   └── app.dart
│   │   │   └── pubspec.yaml
│   │   ├── driver-app/
│   │   │   ├── lib/
│   │   │   │   ├── core/
│   │   │   │   ├── features/
│   │   │   │   ├── main.dart
│   │   │   │   └── app.dart
│   │   │   └── pubspec.yaml
│   │   └── shared-lib/            # Centralized shared library
│   │       ├── lib/
│   │       └── pubspec.yaml
│   └── web/
│       └── admin-dashboard/
│           ├── src/
│           └── package.json
├── shared/                        # Backend shared files (from backend/shared)
│   ├── go/
│   │   ├── client/
│   │   ├── services/
│   │   ├── config/
│   │   └── models/
│   └── kafka/
│       └── schemas/
├── gateway/                       # Merged from backend/api-gateway/kong
│   ├── kong/
│   │   ├── kong.yml
│   │   ├── Dockerfile
│   │   └── kong-init.sh
│   ├── middleware.go
│   ├── handlers.go
│   └── config/
└── database/
    └── migrations/
```

### Best Practices to Apply (From Crab)
1. **Feature-Based Architecture**: Each feature has own bloc, presentation, data layers
2. **Clear Separation**: 
   - `/presentation` - UI (widgets, screens)
   - `/data` - Data sources (local, remote)
   - `/domain` - Business logic (entities, repositories)
3. **Core Module**: Shared utilities, constants, theme
4. **Localization**: Supported from ground up
5. **Dependency Injection**: Clear and organized
6. **Testing**: Co-located with source code
7. **Constants & Config**: Centralized in core

## PHASE 2: Execution Steps

### Step 1: Restructure Mobile Apps
- [x] Create cleaner directory structure
- [x] Consolidate duplicate screens/controllers
- [x] Implement feature-based architecture
- [x] Update all imports
- [x] Merge redundant code

### Step 2: Centralize Shared Library
- [x] Move to `apps/flutter-mobile/shared-lib`
- [x] Update references in both apps
- [x] Add to `pubspec.yaml` as path dependency

### Step 3: Consolidate Backend
- [x] Move `backend/shared` → `shared`
- [x] Merge `backend/api-gateway/kong` → `gateway/kong`
- [x] Update all Go imports
- [x] Remove backend duplicates

### Step 4: Path/Import Updates
- [x] Update all Dart imports
- [x] Update all Go imports
- [x] Update all TypeScript imports
- [x] Update configuration references

### Step 5: Quality Refinement
- [x] Apply best practices
- [x] Remove code duplication
- [x] Enhance error handling
- [x] Optimize performance
- [x] Add comprehensive documentation

## PHASE 3: Best Practices Application

### Feature Module Structure (Per Feature)
```
feature_name/
├── presentation/
│   ├── bloc/           # State management (Bloc/Cubit)
│   ├── pages/          # Full-screen widgets
│   ├── widgets/        # Reusable components
│   └── state/          # State classes
├── domain/
│   ├── entities/       # Business entities
│   ├── repositories/   # Abstract repositories
│   └── usecases/       # Business logic
└── data/
    ├── datasources/    # Local/Remote data
    ├── models/         # Data models
    ├── repositories/   # Repository implementations
    └── providers/      # Dependency providers
```

### Core Module (Shared Across App)
```
core/
├── config/             # App configuration
├── constants/          # Constants
├── theme/              # Theme
├── di/                 # Dependency injection
├── extensions/         # Dart extensions
├── network/            # HTTP client setup
├── storage/            # Local storage setup
└── utils/              # Utilities
```

### Naming Conventions
- **Files**: `snake_case.dart`
- **Classes**: `PascalCase`
- **Constants**: `CONSTANT_NAME`
- **Variables**: `camelCase`
- **Imports**: Use relative imports within feature, absolute for cross-feature

## PHASE 4: Code Quality Improvements

### Error Handling
- Custom exceptions per domain
- Proper error propagation
- User-friendly error messages
- Logging integration

### Performance
- Lazy loading
- Caching strategies
- Image optimization
- API response optimization

### Security
- Secure storage for sensitive data
- Request/response encryption
- Token management
- Input validation

### Testing
- Unit tests for business logic
- Widget tests for UI
- Integration tests for features
- Mocking for external dependencies

---

**Status**: Ready for execution  
**Complexity**: High (requires systematic import updates)  
**Timeline**: 2-3 hours for complete restructuring  
**Risk**: Low (with proper git backup)  
