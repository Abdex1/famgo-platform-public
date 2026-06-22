# 🔌 PROTOBUF REGISTRY: gRPC Service Definitions

**Status:** Task 2 Phase 2.3 Complete  
**Location:** shared/contracts/protobufs/  
**Status:** Ready for Future gRPC Services (Tasks 5+)

---

## PROTOBUF STRUCTURE

```
shared/contracts/protobufs/
├── README.md (this will be created in Task X)
├── common/
│   └── v1/
│       ├── envelope.proto
│       ├── error.proto
│       └── status.proto
├── auth/
│   └── v1/
│       ├── auth_service.proto
│       └── auth_messages.proto
├── ride/
│   └── v1/
│       ├── ride_service.proto
│       └── ride_messages.proto
├── driver/
│   └── v1/
│       ├── driver_service.proto
│       └── driver_messages.proto
├── gps/
│   └── v1/
│       ├── gps_service.proto
│       └── gps_messages.proto
└── payment/
    └── v1/
        ├── payment_service.proto
        └── payment_messages.proto
```

---

## CURRENT STATUS

### ✅ Structure Ready
- [x] Directory layout prepared
- [x] Location validated
- [x] Ready for service definitions

### ⏳ To Be Implemented (Future Tasks)

**Task 5+ (GPS Service):**
```protobuf
service GPSService {
  rpc UpdateLocation(UpdateLocationRequest) returns (UpdateLocationResponse);
  rpc GetNearbyDrivers(GetNearbyRequest) returns (GetNearbyResponse);
  rpc GetTripRoute(GetTripRouteRequest) returns (GetTripRouteResponse);
}
```

**Task 8+ (Dispatch Service):**
```protobuf
service DispatchService {
  rpc AssignRide(AssignRideRequest) returns (AssignRideResponse);
  rpc GetAvailableDrivers(GetAvailableRequest) returns (GetAvailableResponse);
  rpc CancelAssignment(CancelAssignmentRequest) returns (CancelAssignmentResponse);
}
```

---

## PROTOBUF BEST PRACTICES

### 1. File Organization

**Pattern:**
```
shared/contracts/protobufs/
└── {domain}/
    └── v1/
        ├── {domain}_service.proto (service definitions)
        └── {domain}_messages.proto (message definitions)
```

**Rationale:**
- Separates service interfaces from data structures
- Versioning by directory (v1, v2, etc.)
- Clear domain ownership

### 2. Proto File Naming

```
✅ auth_service.proto (service file)
✅ auth_messages.proto (messages file)
✅ auth_errors.proto (error types)
```

### 3. Message Naming

```protobuf
// ✅ Good
message CreateRideRequest {
  string user_id = 1;
  double pickup_lat = 2;
}

// ❌ Avoid
message Ride {
  // Ambiguous: is this a request or response?
}

message RideMsg {
  // Unclear abbreviation
}
```

### 4. Field Numbering

```protobuf
message RideRequest {
  string trip_id = 1;        // Never reuse numbers!
  string user_id = 2;
  double pickup_lat = 3;
  double pickup_lng = 4;
  // If you remove a field later: don't reuse its number
  // reserved 5; (marks number as reserved)
}
```

**Important:** Never reuse field numbers. If field removed, mark as reserved.

### 5. Versioning

```
shared/contracts/protobufs/ride/
├── v1/
│   ├── ride_service.proto (initial)
│   └── ride_messages.proto
└── v2/
    ├── ride_service.proto (breaking changes)
    └── ride_messages.proto
```

---

## USING PROTOBUFS IN SERVICES

### When to Use

**Use gRPC (protobuf) when:**
- Internal service-to-service communication
- Performance critical (low latency required)
- High-volume communication
- Bidirectional streaming needed

**Examples:**
- GPS service → Dispatch service (frequent location updates)
- Auth service → All services (token validation)
- Payment service → Wallet service (transaction processing)

### When NOT to Use

**Use HTTP/JSON when:**
- External APIs (mobile apps, third-party integrations)
- Webhook/callback patterns
- Simple request-response patterns
- Human-readable logging needed

**Examples:**
- Mobile app → API Gateway
- Third-party payment providers → Payment service
- Analytics dashboards → Analytics service

---

## PROTOBUF COMPILATION

### Compilation Target

**Go Services:**
```bash
# Generate Go code
protoc --go_out=. --go-grpc_out=. shared/contracts/protobufs/**/*.proto
```

**Output Files:**
```
shared/contracts/protobufs/ride/v1/
├── ride_service.proto
├── ride_messages.proto
├── ride_service_grpc.pb.go (generated)
├── ride_messages.pb.go (generated)
└── ...
```

### CI/CD Integration

**In Dockerfile or CI pipeline:**
```bash
# Compile all protobufs during build
RUN protoc --go_out=. --go-grpc_out=. shared/contracts/protobufs/**/*.proto
```

**Failure:** If proto compilation fails, build fails. Ensures type safety.

---

## EXAMPLE PROTOBUFS (For Reference)

### Auth Service Example

```protobuf
// shared/contracts/protobufs/auth/v1/auth_service.proto
syntax = "proto3";

package auth.v1;

option go_package = "github.com/famgo/shared/contracts/protobufs/auth/v1";

service AuthService {
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (RefreshTokenResponse);
  rpc RevokeSession(RevokeSessionRequest) returns (RevokeSessionResponse);
}

message ValidateTokenRequest {
  string token = 1;
  string device_id = 2;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  repeated string roles = 3;
  int64 expires_at = 4;
}

message RefreshTokenRequest {
  string refresh_token = 1;
}

message RefreshTokenResponse {
  string access_token = 1;
  string refresh_token = 2;
  int64 expires_in = 3;
}

message RevokeSessionRequest {
  string session_id = 1;
}

message RevokeSessionResponse {
  bool success = 1;
}
```

### GPS Service Example

```protobuf
// shared/contracts/protobufs/gps/v1/gps_service.proto
syntax = "proto3";

package gps.v1;

option go_package = "github.com/famgo/shared/contracts/protobufs/gps/v1";

service GPSService {
  rpc UpdateLocation(UpdateLocationRequest) returns (UpdateLocationResponse);
  rpc GetNearbyDrivers(GetNearbyDriversRequest) returns (GetNearbyDriversResponse);
  rpc GetTripRoute(GetTripRouteRequest) returns (GetTripRouteResponse);
}

message UpdateLocationRequest {
  string driver_id = 1;
  double latitude = 2;
  double longitude = 3;
  float accuracy = 4;
  int32 heading = 5;
}

message UpdateLocationResponse {
  bool accepted = 1;
}

message GetNearbyDriversRequest {
  double latitude = 1;
  double longitude = 2;
  float radius_meters = 3;
  int32 limit = 4;
}

message Driver {
  string driver_id = 1;
  double distance_meters = 2;
  int32 eta_seconds = 3;
}

message GetNearbyDriversResponse {
  repeated Driver drivers = 1;
}

message GetTripRouteRequest {
  string trip_id = 1;
}

message GetTripRouteResponse {
  string polyline = 1;
  float distance_meters = 2;
  int32 duration_seconds = 3;
}
```

---

## MIGRATION FROM HTTP → GRPC

### When to Migrate

**Signal:** If service calls are:
- Happening >100 times/second
- Between services (not external)
- Performance critical

**Decision:**
1. Measure current latency (HTTP)
2. Estimate gRPC latency savings
3. If >20% improvement expected → Migrate

### Migration Process

1. **Define .proto files**
   - Create protobufs/[service]/v1/*.proto
   - Define all service methods
   - Define all message types

2. **Generate Go code**
   - Run protoc compiler
   - Import generated packages

3. **Implement gRPC server**
   - Implement service methods
   - Deploy alongside HTTP (dual-running)

4. **Migrate consumers**
   - Update client to use gRPC
   - Remove HTTP client
   - Deploy updated consumer

5. **Deprecate HTTP**
   - Remove HTTP endpoints
   - Archive HTTP code

---

## GOVERNANCE

### Adding New Proto Files

**Checklist:**
- [ ] Create file in shared/contracts/protobufs/
- [ ] Follow naming conventions
- [ ] Use v1 directory
- [ ] Add to this registry
- [ ] Compile successfully
- [ ] Add to MIGRATION.md

### Approving Proto Changes

**Who approves:**
- Tech Lead (architectural impact)
- Service Owner (implementation impact)
- Security Team (if auth/encryption related)

**PR Checklist:**
- [ ] Proto file syntax valid
- [ ] Compiles without errors
- [ ] Backward compatible (or version bumped)
- [ ] Documentation updated
- [ ] Example client code provided

---

## COMMON GOTCHAS

### ❌ Reusing Field Numbers
```protobuf
// WRONG:
message User {
  string user_id = 1;
  // removed: string email = 2; (DON'T reuse 2!)
  string phone = 2; // ❌ WRONG - already used
}
```

**Fix:**
```protobuf
// CORRECT:
message User {
  string user_id = 1;
  reserved 2; // Mark as reserved
  string phone = 3;
}
```

### ❌ Changing Field Types
```protobuf
// WRONG: Changing amount from int64 to float
message PaymentRequest {
  // v1: int64 amount = 1;
  float amount = 1; // ❌ Breaks v1 consumers
}
```

**Fix:**
```protobuf
// CORRECT: Create v2 message
message PaymentRequestV2 {
  float amount = 1; // New version, new field
}
```

### ❌ Missing go_package
```protobuf
// WRONG:
syntax = "proto3";
package gps.v1;
// Missing: option go_package = "...";
```

**Fix:**
```protobuf
// CORRECT:
syntax = "proto3";
package gps.v1;
option go_package = "github.com/famgo/shared/contracts/protobufs/gps/v1";
```

---

## PROTOBUF TESTING

### Unit Tests

```go
// Test proto compilation
func TestProtoCompilation(t *testing.T) {
    // Verify proto files exist
    // Verify generated code exists
    // Verify package imports work
}

// Test gRPC service
func TestGPSService(t *testing.T) {
    client := gps.NewGPSServiceClient(conn)
    resp, err := client.UpdateLocation(ctx, &gps.UpdateLocationRequest{
        DriverID: "driver-123",
        Latitude: 37.7749,
        Longitude: -122.4194,
    })
    assert.NoError(t, err)
    assert.True(t, resp.Accepted)
}
```

### Integration Tests

```go
// Test service communication
func TestServiceCommunication(t *testing.T) {
    // Start gRPC server
    // Connect with client
    // Send request
    // Verify response
}
```

---

## NEXT STEPS (FOR FUTURE TASKS)

### Task 5+: When Implementing GPS Service

1. Create shared/contracts/protobufs/gps/v1/gps_*.proto files
2. Define GPSService with methods
3. Generate Go code
4. Implement server in gps-service/
5. Implement client in dispatch-service/
6. Test communication
7. Deploy

---

**Protobuf Registry:** ✅ READY FOR FUTURE IMPLEMENTATION

