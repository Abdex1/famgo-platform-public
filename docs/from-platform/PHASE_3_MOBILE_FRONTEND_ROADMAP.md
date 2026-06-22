# 📱 PHASE 3: MOBILE APPS & FRONTEND - IMPLEMENTATION ROADMAP

## Executive Summary
**Goal**: Build coherent mobile apps and web dashboards  
**Timeline**: 3 weeks (120 hours)  
**Deliverables**: Fully functional Rider App, Driver App, and Admin Dashboard  

---

## ARCHITECTURE: Coherent Mobile & Frontend Stack

```
┌─────────────────────────────────────────────────────────┐
│                  Presentation Layer                     │
├──────────────────┬──────────────────┬──────────────────┤
│  Flutter Rider   │  Flutter Driver  │  React Dashboard │
│      App         │       App        │    (Web)         │
└────────┬─────────┴────────┬─────────┴────────┬─────────┘
         │                  │                  │
         └──────────────────┼──────────────────┘
                            │
         ┌──────────────────▼─────────────────┐
         │  Shared Flutter Library            │
         │  (DioClient, WebSocket, Storage)   │
         └──────────────────┬─────────────────┘
                            │
         ┌──────────────────▼─────────────────┐
         │     API Gateway (Kong)             │
         │  (Rate limiting, JWT, CORS)        │
         └──────────────────┬─────────────────┘
                            │
    ┌───────────────────────┼───────────────────────┐
    │                       │                       │
┌───▼──────┐         ┌──────▼──────┐       ┌──────▼──────┐
│ gRPC     │         │   REST      │       │  WebSocket  │
│ Services │         │  (Legacy)   │       │  (Real-time)│
└──────────┘         └─────────────┘       └─────────────┘
```

---

## WEEK 1: Shared Flutter Infrastructure (40 hours)

### Shared Flutter Library Setup

#### Directory Structure
```
shared-flutter-lib/
├── lib/
│   ├── core/
│   │   ├── api/
│   │   │   ├── dio_client.dart        # HTTP client wrapper
│   │   │   ├── interceptors.dart      # Auth, Error, Telemetry
│   │   │   ├── api_response.dart      # Standard response model
│   │   │   └── exceptions.dart        # Exception hierarchy
│   │   │
│   │   ├── services/
│   │   │   ├── websocket_service.dart # Real-time events
│   │   │   ├── storage_service.dart   # Local persistence
│   │   │   ├── location_service.dart  # GPS
│   │   │   ├── auth_service.dart      # JWT handling
│   │   │   ├── notification_service.dart # Push notifications
│   │   │   ├── logger_service.dart    # Structured logging
│   │   │   └── telemetry_service.dart # OpenTelemetry
│   │   │
│   │   ├── di/
│   │   │   └── service_locator.dart   # GetIt setup
│   │   │
│   │   ├── models/
│   │   │   ├── ride.dart
│   │   │   ├── driver.dart
│   │   │   ├── payment.dart
│   │   │   ├── location.dart
│   │   │   ├── user.dart
│   │   │   └── exception.dart
│   │   │
│   │   ├── constants/
│   │   │   ├── api_endpoints.dart
│   │   │   ├── app_constants.dart
│   │   │   └── error_messages.dart
│   │   │
│   │   └── utils/
│   │       ├── validators.dart
│   │       ├── formatters.dart
│   │       ├── extensions.dart
│   │       └── logger.dart
│   │
│   └── pubspec.yaml
│
├── test/
│   ├── unit/
│   │   ├── dio_client_test.dart
│   │   ├── websocket_service_test.dart
│   │   └── validators_test.dart
│   │
│   └── mock/
│       ├── mock_dio_adapter.dart
│       └── mock_websocket.dart
│
└── README.md
```

#### 1. Core: DioClient with Interceptors (8 hours)

```dart
// shared-flutter-lib/lib/core/api/dio_client.dart
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:uuid/uuid.dart';

class DioClient {
  late final Dio _dio;
  final String _baseUrl;
  final FlutterSecureStorage _storage;
  
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
          'X-Request-ID': _generateRequestId(),
        },
      ),
    );
    
    // Add interceptors in order
    _dio.interceptors.add(AuthInterceptor(_storage));
    _dio.interceptors.add(ErrorInterceptor());
    _dio.interceptors.add(TelemetryInterceptor());
    _dio.interceptors.add(LoggingInterceptor());
  }
  
  String _generateRequestId() => const Uuid().v4();
  
  // Generic GET with automatic deserialization
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
        );
      }
      
      return fromJson(apiResponse.data);
    } on DioException catch (e) {
      throw _mapDioException(e);
    }
  }
  
  // Generic POST
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
  
  // Map DioException to AppException
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
  
  // Close Dio instance
  void close() {
    _dio.close();
  }
}
```

**Files to create**:
- `shared-flutter-lib/lib/core/api/dio_client.dart`
- `shared-flutter-lib/lib/core/api/interceptors.dart`
- `shared-flutter-lib/lib/core/api/api_response.dart`
- `shared-flutter-lib/lib/core/api/exceptions.dart`

---

#### 2. WebSocket Service (8 hours)

```dart
// shared-flutter-lib/lib/core/services/websocket_service.dart
import 'package:socket_io_client/socket_io_client.dart' as IO;

class WebSocketService {
  late IO.Socket _socket;
  final String _baseUrl;
  final String _token;
  
  final Map<String, List<Function(dynamic)>> _listeners = {};
  final bool _isConnected = false;
  
  WebSocketService({
    required String baseUrl,
    required String token,
  })  : _baseUrl = baseUrl,
        _token = token;
  
  Future<void> connect() async {
    _socket = IO.io(
      _baseUrl,
      IO.OptionBuilder()
        .setTransports(['websocket'])
        .disableAutoConnect()
        .setReconnectionDelay(1000)
        .setReconnectionDelayMax(5000)
        .setReconnectionAttempts(10)
        .extraHeaders({'Authorization': 'Bearer $_token'})
        .build(),
    );
    
    _socket.on('connect', (_) {
      print('WebSocket connected');
      _onConnected();
    });
    
    _socket.on('disconnect', (_) {
      print('WebSocket disconnected');
    });
    
    _socket.on('error', (error) {
      print('WebSocket error: $error');
    });
    
    _socket.connect();
  }
  
  void _onConnected() {
    // Join user-specific room
    _socket.emit('join_room', {'room': 'user_$_userId'});
  }
  
  // Subscribe to event
  void on(String event, Function(dynamic) callback) {
    if (!_listeners.containsKey(event)) {
      _listeners[event] = [];
      _socket.on(event, (data) {
        for (final cb in _listeners[event]!) {
          cb(data);
        }
      });
    }
    _listeners[event]!.add(callback);
  }
  
  // Emit event
  void emit(String event, dynamic data) {
    _socket.emit(event, data);
  }
  
  // Unsubscribe
  void off(String event) {
    _socket.off(event);
    _listeners.remove(event);
  }
  
  // Disconnect
  void disconnect() {
    _socket.disconnect();
  }
}
```

**Files to create**:
- `shared-flutter-lib/lib/core/services/websocket_service.dart`
- `shared-flutter-lib/lib/core/services/location_service.dart`
- `shared-flutter-lib/lib/core/services/storage_service.dart`
- `shared-flutter-lib/lib/core/services/auth_service.dart`

---

#### 3. Dependency Injection Setup (4 hours)

```dart
// shared-flutter-lib/lib/core/di/service_locator.dart
import 'package:get_it/get_it.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

final getIt = GetIt.instance;

void setupServiceLocator() {
  // Storage
  getIt.registerSingleton<FlutterSecureStorage>(
    const FlutterSecureStorage(),
  );
  
  // API Client
  getIt.registerSingleton<DioClient>(
    DioClient(
      baseUrl: 'http://api.famgo.et',
      storage: getIt<FlutterSecureStorage>(),
    ),
  );
  
  // WebSocket Service (lazy - initialized on demand)
  getIt.registerLazySingleton<WebSocketService>(() {
    return WebSocketService(
      baseUrl: 'http://api.famgo.et',
      token: getIt<AuthService>().token,
    );
  });
  
  // Location Service
  getIt.registerSingleton<LocationService>(
    LocationService(),
  );
  
  // Auth Service
  getIt.registerSingleton<AuthService>(
    AuthService(
      storage: getIt<FlutterSecureStorage>(),
      dioClient: getIt<DioClient>(),
    ),
  );
}
```

**Files to create**:
- `shared-flutter-lib/lib/core/di/service_locator.dart`
- `shared-flutter-lib/lib/core/models/` (all data models)
- `shared-flutter-lib/lib/core/constants/` (all constants)
- `shared-flutter-lib/lib/core/utils/` (all utilities)

#### 4. Shared Models (8 hours)

```dart
// shared-flutter-lib/lib/core/models/ride.dart
class Ride {
  final String id;
  final String riderId;
  final String? driverId;
  final Location pickupLocation;
  final Location dropoffLocation;
  final RideStatus status;
  final double estimatedFare;
  final double? actualFare;
  final int estimatedDurationMinutes;
  final DateTime createdAt;
  final DateTime? completedAt;
  
  Ride({
    required this.id,
    required this.riderId,
    this.driverId,
    required this.pickupLocation,
    required this.dropoffLocation,
    required this.status,
    required this.estimatedFare,
    this.actualFare,
    required this.estimatedDurationMinutes,
    required this.createdAt,
    this.completedAt,
  });
  
  factory Ride.fromJson(Map<String, dynamic> json) {
    return Ride(
      id: json['id'] as String,
      riderId: json['rider_id'] as String,
      driverId: json['driver_id'] as String?,
      pickupLocation: Location.fromJson(json['pickup_location']),
      dropoffLocation: Location.fromJson(json['dropoff_location']),
      status: RideStatus.fromString(json['status']),
      estimatedFare: (json['estimated_fare'] as num).toDouble(),
      actualFare: (json['actual_fare'] as num?)?.toDouble(),
      estimatedDurationMinutes: json['estimated_duration_minutes'] as int,
      createdAt: DateTime.parse(json['created_at']),
      completedAt: json['completed_at'] != null 
        ? DateTime.parse(json['completed_at']) 
        : null,
    );
  }
  
  Map<String, dynamic> toJson() => {
    'id': id,
    'rider_id': riderId,
    'driver_id': driverId,
    'pickup_location': pickupLocation.toJson(),
    'dropoff_location': dropoffLocation.toJson(),
    'status': status.value,
    'estimated_fare': estimatedFare,
    'actual_fare': actualFare,
    'estimated_duration_minutes': estimatedDurationMinutes,
    'created_at': createdAt.toIso8601String(),
    'completed_at': completedAt?.toIso8601String(),
  };
}

enum RideStatus {
  pending('pending'),
  accepted('accepted'),
  started('started'),
  completed('completed'),
  cancelled('cancelled'),
  noShow('no_show');
  
  final String value;
  
  const RideStatus(this.value);
  
  factory RideStatus.fromString(String value) {
    return RideStatus.values.firstWhere(
      (status) => status.value == value,
      orElse: () => RideStatus.pending,
    );
  }
}

class Location {
  final double latitude;
  final double longitude;
  final String address;
  
  Location({
    required this.latitude,
    required this.longitude,
    required this.address,
  });
  
  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      address: json['address'] as String,
    );
  }
  
  Map<String, dynamic> toJson() => {
    'latitude': latitude,
    'longitude': longitude,
    'address': address,
  };
}
```

**Files to create** (all core models):
- `shared-flutter-lib/lib/core/models/ride.dart`
- `shared-flutter-lib/lib/core/models/driver.dart`
- `shared-flutter-lib/lib/core/models/payment.dart`
- `shared-flutter-lib/lib/core/models/user.dart`
- `shared-flutter-lib/lib/core/models/location.dart`
- `shared-flutter-lib/lib/core/models/wallet.dart`

#### 5. Shared Library Tests (4 hours)

```dart
// shared-flutter-lib/test/unit/dio_client_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:shared_flutter_lib/core/api/dio_client.dart';

void main() {
  group('DioClient', () {
    late DioClient dioClient;
    late MockHttpClient mockHttpClient;
    
    setUp(() {
      // Setup mock
      mockHttpClient = MockHttpClient();
      dioClient = DioClient(
        baseUrl: 'http://test.com',
        httpClient: mockHttpClient,
      );
    });
    
    test('GET request returns data on success', () async {
      // Arrange
      when(mockHttpClient.get(...)).thenAnswer(
        (_) async => http.Response('{"success": true, "data": {}}', 200),
      );
      
      // Act
      final result = await dioClient.get('/test');
      
      // Assert
      expect(result, isNotNull);
    });
    
    test('GET request throws exception on 401', () async {
      // Arrange
      when(mockHttpClient.get(...)).thenAnswer(
        (_) async => http.Response('{}', 401),
      );
      
      // Act & Assert
      expect(
        () => dioClient.get('/test'),
        throwsA(isA<UnauthorizedException>()),
      );
    });
  });
}
```

**Files to create**:
- `shared-flutter-lib/test/unit/dio_client_test.dart`
- `shared-flutter-lib/test/unit/websocket_service_test.dart`
- `shared-flutter-lib/test/mock/mock_dio_adapter.dart`

---

## WEEK 2: Flutter Rider App (40 hours)

### Core Structure
```
flutter-rider-app/
├── lib/
│   ├── main.dart
│   ├── config/
│   │   ├── app_config.dart
│   │   └── routes.dart
│   ├── features/
│   │   └── rider/
│   │       ├── presentation/
│   │       │   ├── screens/
│   │       │   │   ├── auth_screen.dart
│   │       │   │   ├── home_screen.dart
│   │       │   │   ├── booking_screen.dart
│   │       │   │   ├── tracking_screen.dart
│   │       │   │   ├── payment_screen.dart
│   │       │   │   ├── rating_screen.dart
│   │       │   │   ├── history_screen.dart
│   │       │   │   ├── wallet_screen.dart
│   │       │   │   └── profile_screen.dart
│   │       │   │
│   │       │   ├── controllers/
│   │       │   │   ├── auth_controller.dart
│   │       │   │   ├── booking_controller.dart
│   │       │   │   ├── tracking_controller.dart
│   │       │   │   ├── payment_controller.dart
│   │       │   │   └── wallet_controller.dart
│   │       │   │
│   │       │   └── widgets/
│   │       │       ├── location_search_widget.dart
│   │       │       ├── fare_estimate_widget.dart
│   │       │       ├── driver_card_widget.dart
│   │       │       ├── ride_status_widget.dart
│   │       │       └── rating_widget.dart
│   │       │
│   │       ├── domain/
│   │       │   ├── entities/ (empty - use shared models)
│   │       │   └── repositories/
│   │       │       ├── auth_repository.dart
│   │       │       ├── ride_repository.dart
│   │       │       ├── payment_repository.dart
│   │       │       └── wallet_repository.dart
│   │       │
│   │       └── data/
│   │           ├── datasources/
│   │           │   ├── ride_remote_datasource.dart
│   │           │   ├── ride_local_datasource.dart
│   │           │   └── auth_remote_datasource.dart
│   │           │
│   │           └── models/
│   │               └── (empty - use shared models)
│   │
│   ├── core/
│   │   └── (Shared from shared-flutter-lib)
│   │
│   └── theme/
│       ├── app_theme.dart
│       ├── app_colors.dart
│       └── text_styles.dart
│
└── test/
    └── (integration tests)
```

### Implementation Tasks (40 hours)

**Day 1: Setup & Auth (6 hours)**
- Authentication screens (login, register, verification)
- GetX auth controller
- JWT token storage
- Auto-login on app restart

**Day 2-3: Booking Flow (10 hours)**
- Location search (Google Maps integration)
- Fare estimation
- Ride booking confirmation
- Real-time driver matching (WebSocket)

**Day 3-4: Tracking & Real-time (10 hours)**
- Real-time location tracking
- Driver location on map
- ETA updates
- Real-time notifications

**Day 5: Payment & Rating (8 hours)**
- Payment method selection
- Payment processing
- Receipt
- Rating screen

**Day 6: Wallet & Profile (6 hours)**
- Wallet balance and history
- Profile management
- Settings

---

## WEEK 3: Flutter Driver App + Frontend Dashboard (40 hours)

### Flutter Driver App (20 hours)
- Identical structure to Rider app
- Driver-specific screens:
  - Active rides acceptance
  - Route optimization
  - Earnings tracking
  - SOS response

### React Admin Dashboard (20 hours)
```
web/admin-dashboard/
├── src/
│   ├── pages/
│   │   ├── dashboard/
│   │   │   ├── DashboardPage.tsx
│   │   │   ├── RealTimeMetrics.tsx
│   │   │   ├── MapView.tsx
│   │   │   └── SystemHealth.tsx
│   │   │
│   │   ├── users/
│   │   │   ├── UserManagement.tsx
│   │   │   ├── RidersList.tsx
│   │   │   └── DriversList.tsx
│   │   │
│   │   ├── payments/
│   │   │   ├── PaymentReconciliation.tsx
│   │   │   ├── DisputeResolution.tsx
│   │   │   └── PayoutManagement.tsx
│   │   │
│   │   ├── safety/
│   │   │   ├── SOSIncidents.tsx
│   │   │   ├── IncidentDetails.tsx
│   │   │   └── IncidentTimeline.tsx
│   │   │
│   │   ├── fraud/
│   │   │   ├── FraudAlerts.tsx
│   │   │   ├── RiskAnalysis.tsx
│   │   │   └── BlockedUsers.tsx
│   │   │
│   │   └── operations/
│   │       ├── ServiceAreaManagement.tsx
│   │       ├── DriverOnboarding.tsx
│   │       └── PromoCodes.tsx
│   │
│   ├── components/
│   │   ├── Charts.tsx
│   │   ├── Tables.tsx
│   │   ├── Maps.tsx
│   │   ├── RealTimeUpdates.tsx
│   │   └── Notifications.tsx
│   │
│   ├── api/
│   │   ├── apiClient.ts (uses FamGo API)
│   │   ├── endpoints.ts
│   │   └── hooks/ (useQuery, useMutation)
│   │
│   ├── hooks/
│   │   ├── useWebSocket.ts
│   │   ├── useRealTimeData.ts
│   │   └── usePermissions.ts
│   │
│   └── theme/
│       ├── theme.ts
│       ├── colors.ts
│       └── typography.ts
│
└── package.json
```

---

## TESTING PHASE (Parallel with development)

### Unit Tests (Flutter)
```dart
// flutter-rider-app/test/features/rider/booking_controller_test.dart
void main() {
  group('RideBookingController', () {
    late RideBookingController controller;
    late MockRideRepository mockRepository;
    
    setUp(() {
      mockRepository = MockRideRepository();
      controller = RideBookingController(mockRepository);
    });
    
    test('searchRides updates rideOptions', () async {
      // Arrange
      final mockRides = [
        Ride(...),
        Ride(...),
      ];
      when(mockRepository.searchRides(...))
        .thenAnswer((_) async => mockRides);
      
      // Act
      await controller.searchRides(
        Location(...),
        Location(...),
      );
      
      // Assert
      expect(controller.rideOptions.length, 2);
    });
  });
}
```

### Integration Tests (Flutter)
```dart
// flutter-rider-app/test/integration/booking_flow_test.dart
void main() {
  group('Ride Booking Flow', () {
    testWidgets('Complete booking flow', (WidgetTester tester) async {
      // Launch app
      await tester.pumpWidget(MyApp());
      
      // Login
      await tester.tap(find.byType(ElevatedButton));
      await tester.pumpAndSettle();
      
      // Search ride
      await tester.enterText(find.byType(TextField).first, '2nd Avenue');
      await tester.pumpAndSettle();
      
      // Select driver
      await tester.tap(find.byType(DriverCard).first);
      await tester.pumpAndSettle();
      
      // Confirm booking
      await tester.tap(find.text('Confirm Booking'));
      await tester.pumpAndSettle();
      
      // Verify tracking screen appears
      expect(find.byType(TrackingScreen), findsOneWidget);
    });
  });
}
```

### E2E Tests (Via Postman/API)
```bash
# Run full end-to-end flow
newman run FamGo-E2E-Collection.postman_collection.json \
  -e staging.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export report.json
```

---

## DELIVERY CHECKLIST

### Mobile Apps
```
Flutter Rider App
├─ Authentication (login, register, verification)
├─ Ride booking (location search, fare estimate, confirmation)
├─ Real-time tracking (GPS, driver location, ETA)
├─ Payments (method selection, processing, receipt)
├─ Rating & feedback
├─ Wallet management
├─ Profile settings
├─ Offline capability
├─ Push notifications
├─ Error handling
├─ Logging & telemetry
└─ 80%+ test coverage

Flutter Driver App (identical structure)
├─ Ride acceptance
├─ Active ride management
├─ Route optimization
├─ Real-time earnings
├─ SOS response
├─ Performance metrics
└─ 80%+ test coverage
```

### Web Dashboard
```
Admin Dashboard
├─ Real-time metrics & KPIs
├─ Live map with active rides
├─ User management (riders, drivers)
├─ Payment reconciliation
├─ Safety incident management
├─ Fraud detection alerts
├─ Service area configuration
├─ Promo code management
├─ System health monitoring
└─ Comprehensive reporting
```

---

## TIMELINE SUMMARY

```
Week 1: 40 hours
├─ Shared Flutter library (DioClient, WebSocket, Models)
└─ Deliverable: Mobile apps can connect to backend

Week 2: 40 hours
├─ Flutter Rider App (full implementation)
└─ Deliverable: Fully functional rider app

Week 3a: 20 hours
├─ Flutter Driver App (same structure)
└─ Deliverable: Fully functional driver app

Week 3b: 20 hours
├─ React Admin Dashboard
└─ Deliverable: Complete web interface

TOTAL: 120 hours (3 weeks)
```

---

## SUCCESS CRITERIA

✅ Mobile apps connect to live backend  
✅ All API calls use unified DioClient  
✅ Real-time events work via WebSocket  
✅ Offline mode works identically (both apps)  
✅ Error handling consistent (all platforms)  
✅ 80%+ test coverage (mobile)  
✅ Admin can manage entire platform  
✅ Production-ready deployment  

---

**Phase 3 Completion**: Fully functional mobile and web tier  
**Estimated Hours**: 120 hours (3 weeks parallel with Phase 2)  
**Total Project**: 200 hours (5 weeks) from Phase 1 start to production

---

**Next**: Execute Phase 2 & Phase 3 in parallel for maximum efficiency
