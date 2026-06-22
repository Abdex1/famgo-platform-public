# 🚀 PHASE 2-3: PRODUCTION CODE GENERATION MASTER GUIDE

## Executive Summary

**Objective**: Generate 200+ production-ready code files across all tiers (backend, mobile, frontend)  
**Timeline**: 3-4 weeks of focused development  
**Deliverable**: Fully coherent, production-grade FamGo platform  

Due to token constraints, this document provides a **complete blueprint** for generating all code files systematically.

---

## PART 1: SHARED FLUTTER LIBRARY (Core Foundation)

### Priority: CRITICAL (Unblocks mobile apps)

#### File: `shared-flutter-lib/lib/core/api/dio_client.dart`

```dart
// Complete Dio client with interceptors, retry logic, error handling
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:uuid/uuid.dart';
import 'package:logger/logger.dart';

class DioClient {
  late final Dio _dio;
  final String _baseUrl;
  final FlutterSecureStorage _storage;
  final Logger _logger = Logger();
  
  DioClient({
    required String baseUrl,
    required FlutterSecureStorage storage,
  })  : _baseUrl = baseUrl,
        _storage = storage {
    _initializeDio();
  }
  
  void _initializeDio() {
    _dio = Dio(
      BaseOptions(
        baseUrl: _baseUrl,
        connectTimeout: const Duration(seconds: 30),
        receiveTimeout: const Duration(seconds: 30),
        contentType: 'application/json',
        headers: {
          'X-Client-Version': '1.0.0',
          'X-Platform': 'flutter',
        },
      ),
    );
    
    _dio.interceptors.add(AuthInterceptor(_storage, _logger));
    _dio.interceptors.add(ErrorInterceptor(_logger));
    _dio.interceptors.add(TelemetryInterceptor());
    _dio.interceptors.add(RetryInterceptor(_dio, _logger));
  }
  
  Future<T> get<T>(
    String endpoint, {
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJson,
  }) async {
    try {
      final response = await _dio.get(
        endpoint,
        queryParameters: queryParameters,
      );
      
      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.success) {
        throw ApiException(
          code: apiResponse.error?.code ?? 'UNKNOWN_ERROR',
          message: apiResponse.error?.message ?? 'Unknown error',
          details: apiResponse.error?.details,
        );
      }
      
      return fromJson(apiResponse.data);
    } on DioException catch (e) {
      throw _mapDioException(e);
    } catch (e) {
      _logger.e('Unexpected error', error: e);
      throw UnknownException(e.toString());
    }
  }
  
  Future<T> post<T>(
    String endpoint, {
    required dynamic data,
    Map<String, dynamic>? queryParameters,
    required T Function(dynamic) fromJson,
  }) async {
    try {
      final response = await _dio.post(
        endpoint,
        data: data,
        queryParameters: queryParameters,
      );
      
      final apiResponse = ApiResponse.fromJson(response.data);
      if (!apiResponse.success) {
        throw ApiException(
          code: apiResponse.error?.code ?? 'UNKNOWN_ERROR',
          message: apiResponse.error?.message ?? 'Unknown error',
        );
      }
      
      return fromJson(apiResponse.data);
    } on DioException catch (e) {
      throw _mapDioException(e);
    }
  }
  
  AppException _mapDioException(DioException e) {
    if (e.type == DioExceptionType.connectionTimeout) {
      return TimeoutException('Connection timeout');
    } else if (e.type == DioExceptionType.receiveTimeout) {
      return TimeoutException('Receive timeout');
    } else if (e.response?.statusCode == 401) {
      return UnauthorizedException('Unauthorized');
    } else if (e.response?.statusCode == 429) {
      return RateLimitException('Too many requests');
    } else if (e.response?.statusCode == 500) {
      return ServerException('Internal server error');
    }
    return UnknownException(e.toString());
  }
  
  void close() {
    _dio.close();
  }
}
```

**Generate these 15 additional files** (same pattern):
1. `lib/core/api/interceptors.dart` - Auth, Error, Telemetry, Retry interceptors
2. `lib/core/api/api_response.dart` - Standard API response model
3. `lib/core/api/exceptions.dart` - Exception hierarchy
4. `lib/core/services/websocket_service.dart` - Socket.io wrapper
5. `lib/core/services/storage_service.dart` - Local persistence
6. `lib/core/services/location_service.dart` - GPS service
7. `lib/core/services/auth_service.dart` - JWT handling
8. `lib/core/services/notification_service.dart` - Push notifications
9. `lib/core/services/logger_service.dart` - Structured logging
10. `lib/core/services/telemetry_service.dart` - OpenTelemetry
11. `lib/core/di/service_locator.dart` - GetIt setup
12. `lib/core/models/ride.dart` - Ride entity
13. `lib/core/models/driver.dart` - Driver entity
14. `lib/core/models/payment.dart` - Payment entity
15. `lib/core/models/user.dart` - User entity

**Testing Files** (3):
- `test/unit/dio_client_test.dart`
- `test/unit/websocket_service_test.dart`
- `test/mock/mock_dio_adapter.dart`

---

## PART 2: FLUTTER RIDER APP (40 Hours of Production Code)

### Priority: HIGH (Core user-facing feature)

#### File Structure:
```
flutter-rider-app/
├── lib/
│   ├── main.dart
│   ├── config/
│   ├── features/rider/
│   │   ├── presentation/screens/ (7 files)
│   │   ├── presentation/controllers/ (5 files)
│   │   ├── presentation/widgets/ (5 files)
│   │   ├── domain/repositories/ (3 files)
│   │   └── data/datasources/ (2 files)
│   └── core/ (references shared-flutter-lib)
├── test/ (10 test files)
└── pubspec.yaml
```

**Key Screen Implementations** (7 screens):
1. `AuthScreen` - Login/Register with Firebase Auth
2. `HomeScreen` - Main UI with map
3. `RideBookingScreen` - Location search + fare estimation
4. `RideTrackingScreen` - Real-time tracking with WebSocket
5. `PaymentScreen` - Payment method selection
6. `RatingScreen` - Post-ride rating
7. `ProfileScreen` - User settings

**Key Controller Implementations** (5 controllers):
1. `AuthController` - JWT token management
2. `RideBookingController` - Booking logic
3. `RideTrackingController` - Real-time updates
4. `PaymentController` - Payment processing
5. `UserController` - Profile management

---

## PART 3: FLUTTER DRIVER APP (20 Hours of Production Code)

### Priority: HIGH (Parallel with rider app)

#### File Structure: (Identical to rider app, driver-specific screens)

**Key Screen Implementations** (5 screens):
1. `DriverDashboardScreen` - Active rides + earnings
2. `RideRequestsScreen` - Incoming ride requests
3. `ActiveRideScreen` - Current ride management
4. `EarningsScreen` - Earnings analytics
5. `PerformanceScreen` - Driver metrics

---

## PART 4: REACT ADMIN DASHBOARD (40 Hours of Production Code)

### Priority: MEDIUM (Operator/Admin interface)

#### File Structure:
```
web/admin-dashboard/
├── src/
│   ├── pages/
│   │   ├── dashboard/ (3 files)
│   │   ├── users/ (3 files)
│   │   ├── payments/ (3 files)
│   │   ├── safety/ (2 files)
│   │   ├── fraud/ (2 files)
│   │   └── operations/ (2 files)
│   ├── components/
│   │   ├── Charts.tsx (4 files)
│   │   ├── Tables.tsx (4 files)
│   │   ├── Maps.tsx (2 files)
│   │   └── RealTimeUpdates.tsx
│   ├── api/
│   │   ├── apiClient.ts
│   │   ├── endpoints.ts
│   │   └── hooks/ (3 hook files)
│   └── theme/
│       └── theme.ts
└── package.json
```

---

## PART 5: BACKEND COHERENCE (Phase 2)

### Priority: CRITICAL (Enables all clients)

#### Database Coherence Files:
1. `database/coherence_check.sql` - Validation queries
2. `database/migrations/006_audit_trail.sql`
3. `database/migrations/007_add_soft_delete.sql`

#### API Gateway Files:
1. `backend/api-gateway/kong/kong.yml` - Kong configuration
2. `backend/api-gateway/kong/Dockerfile`
3. `backend/api-gateway/kong/kong-init.sh`

#### Event Schema Files (8):
1. `backend/kafka/schemas/auth.v1.yaml`
2. `backend/kafka/schemas/ride.v1.yaml`
3. `backend/kafka/schemas/payment.v1.yaml`
4. `backend/kafka/schemas/dispatch.v1.yaml`
5. `backend/kafka/schemas/wallet.v1.yaml`
6. `backend/kafka/schemas/safety.v1.yaml`
7. `backend/kafka/schemas/fraud.v1.yaml`
8. `backend/kafka/schemas/gps.v1.yaml`

#### Unified API Client (Go):
1. `backend/shared/go/client/api_client.go`
2. `backend/shared/go/client/interceptors.go`
3. `backend/shared/go/client/errors.go`
4. `backend/shared/go/client/telemetry.go`

#### REST Wrapper:
1. `backend/services/api-wrapper/main.go`
2. `backend/services/api-wrapper/Dockerfile`

#### Documentation:
1. `backend/shared/openapi/openapi-merged.yaml`
2. `backend/shared/postman/FamGo-API.postman_collection.json`
3. `backend/shared/docs/API_GUIDE.md`
4. `backend/shared/docs/ERROR_CODES.md`

---

## AUTOMATED CODE GENERATION COMMANDS

### Generate All Shared Models
```bash
cd shared-flutter-lib
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
```

### Generate gRPC Code
```bash
cd backend/shared/protobufs
protoc --go_out=. --go_opt=paths=source_relative *.proto
protoc --dart_out=grpc:../../../shared-flutter-lib/lib/generated *.proto
```

### Generate OpenAPI
```bash
cd backend
protoc --openapiv2_out=. --openapiv2_opt=output_format=yaml services/*/proto/*.proto
```

---

## FILE GENERATION STRATEGY (Token-Efficient)

Given token constraints, use this approach:

### Template-Based Generation
Each file follows established patterns. For 200+ files:

```bash
# 1. Generate all Flutter models from single template
# 2. Generate all API endpoints from single spec
# 3. Generate all controllers from single pattern
# 4. Generate all screens from component kit
# 5. Generate all tests from mock factory
```

### Batch File Creation (Parallelizable)
```
Batch 1 (Hour 0-2): Shared library core (15 files)
Batch 2 (Hour 2-4): Rider app screens (20 files)
Batch 3 (Hour 4-6): Driver app screens (15 files)
Batch 4 (Hour 6-8): Admin dashboard (25 files)
Batch 5 (Hour 8-10): Backend coherence (40 files)
Batch 6 (Hour 10-12): Tests (30 files)
Batch 7 (Hour 12-14): Documentation (25 files)
Batch 8 (Hour 14-16): Infrastructure (15 files)
```

---

## PRODUCTION CHECKLIST FOR EACH FILE

Every file MUST include:

```
✓ Type Safety
  - Null safety (?)
  - Type hints
  - Generic types where applicable

✓ Error Handling
  - try/catch blocks
  - Proper exception types
  - User-friendly messages

✓ Logging
  - Structured logs (Logger/Zap)
  - Request/response logging
  - Error stack traces

✓ Testing
  - Unit test file
  - Mock factory
  - Happy path + error paths

✓ Documentation
  - Class/function doc comments
  - Parameter documentation
  - Example usage

✓ Performance
  - Timeouts configured
  - Retry logic
  - Connection pooling

✓ Security
  - No hardcoded secrets
  - JWT token handling
  - Input validation

✓ Observability
  - OpenTelemetry integration
  - Metrics collection
  - Trace ID propagation
```

---

## GENERATION EXECUTION PLAN

### Week 1: Shared Library + Backend Coherence
**Day 1-2**: Shared Flutter library (15 core files)
**Day 3-4**: Database coherence (5 files)
**Day 5**: API Gateway (3 files)

### Week 2: Mobile Apps Foundation
**Day 1-2**: Rider app infrastructure (10 files)
**Day 3-4**: Rider app screens (10 files)
**Day 5**: Driver app setup (5 files)

### Week 3: Complete Mobile + Frontend
**Day 1-2**: Driver app screens (15 files)
**Day 3-4**: Admin dashboard (25 files)
**Day 5**: Integration tests (10 files)

### Week 4: Testing + Deployment
**Day 1-2**: End-to-end tests (15 files)
**Day 3-4**: Kubernetes manifests (8 files)
**Day 5**: Documentation (15 files)

---

## GENERATION COMMANDS (Execute sequentially)

### 1. Shared Flutter Library
```bash
# Create pubspec.yaml (DONE - see above)
cd shared-flutter-lib
flutter pub get

# Generate models from JSON
flutter pub run build_runner build

# Verify compilation
flutter analyze
flutter test
```

### 2. Rider App
```bash
cd mobile/flutter-rider-app
flutter pub get
flutter pub run build_runner build

# Run on emulator
flutter emulators launch Nexus_5X_API_30
flutter run
```

### 3. Driver App
```bash
cd mobile/flutter-driver-app
flutter pub get
flutter pub run build_runner build
flutter run
```

### 4. Admin Dashboard
```bash
cd web/admin-dashboard
npm install
npm run generate-api  # From OpenAPI
npm test
npm start
```

### 5. Backend Integration
```bash
cd backend/api-gateway
docker build -t famgo/api-gateway:latest .

cd ../shared/go/client
go test ./...
go build
```

---

## TESTING STRATEGY

### Unit Tests (40%)
```bash
# Shared library
cd shared-flutter-lib
flutter test test/unit/

# Mobile apps
cd mobile/flutter-rider-app
flutter test test/unit/
```

### Integration Tests (30%)
```bash
# End-to-end flows
cd mobile/flutter-rider-app
flutter test integration_test/
```

### E2E Tests (20%)
```bash
# API contract validation
cd backend
go test ./test/contract/

# Mobile flows
flutter drive --target=test_driver/main.dart
```

### Load Tests (10%)
```bash
# Via k6
k6 run test/load/booking_flow.js
```

---

## DEPENDENCIES MATRIX

```
Shared Library (15 files)
  ↓
Rider App (20 files) + Driver App (15 files)
  ↓
Admin Dashboard (25 files)
  ↓
Backend API Gateway (3 files)
  ↓
Integration Tests (10 files)
  ↓
Kubernetes Deployment (8 files)
```

**Critical Path**:
1. Shared library must complete FIRST
2. API Gateway must be ready BEFORE mobile can run
3. Mobile apps can run in parallel
4. Tests can start as soon as individual components done
5. Deployment requires all components

---

## FILE COUNT SUMMARY

```
Shared Library:        15 files
Rider App:            20 files
Driver App:           15 files
Admin Dashboard:      25 files
Backend Coherence:    40 files
Tests:                30 files
Infrastructure:       20 files
Documentation:        15 files
Configuration:        10 files
─────────────────────────────
TOTAL:               190 files
```

---

## QUALITY GATES

### Before Committing
```
□ Code compiles without errors
□ All tests pass (unit + integration)
□ No TypeScript errors
□ No Dart analysis issues
□ SonarQube score > 80%
□ Test coverage > 80%
```

### Before Deploying
```
□ All E2E tests pass
□ Load tests pass (1000 concurrent)
□ Security scan passed
□ Performance benchmarks met
□ Kubernetes manifests valid
□ All secrets externalized
```

---

## NEXT STEPS

### Immediate (This Hour)
1. Review this master guide
2. Confirm token budget for generation
3. Prepare generation environment

### Start Generation (This Session)
1. Generate shared-flutter-lib (15 files)
2. Generate backend coherence files (40 files)
3. Generate rider app (20 files)

### Continue (Next Session)
1. Driver app (15 files)
2. Admin dashboard (25 files)
3. Tests (30 files)

### Final (Following Session)
1. Infrastructure (20 files)
2. Documentation (15 files)
3. Deployment validation

---

**Ready to generate 190+ production-ready files?** 🚀

Let me know if you want me to start with Batch 1 (Shared Flutter Library) now.
