# FamGo Enterprise Platform — Practical Extraction & Implementation Guide
## Phase-by-Phase Code Extraction & Service Building

---

## QUICK START: What to Do This Week

### Week 1 Task: Build NestJS Service Template

**Source:** Ceng-Carpool `backend/src/modules/`

This is your foundation. Every service will copy this structure.

#### Step 1: Create Service Template Directory

```bash
mkdir -p services/_template/{src,tests}
cd services/_template
npm init -y
```

#### Step 2: Install Dependencies

```json
// services/_template/package.json
{
  "name": "@famgo/template-service",
  "version": "0.0.1",
  "description": "FamGo Service Template",
  "main": "dist/main.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/main.js",
    "start:dev": "ts-node src/main.ts",
    "test": "jest",
    "test:watch": "jest --watch",
    "migration:generate": "typeorm migration:generate",
    "migration:run": "typeorm migration:run",
    "migration:revert": "typeorm migration:revert"
  },
  "dependencies": {
    "@nestjs/common": "^10.0.0",
    "@nestjs/core": "^10.0.0",
    "@nestjs/jwt": "^10.0.0",
    "@nestjs/passport": "^10.0.0",
    "@nestjs/swagger": "^7.0.0",
    "@nestjs/typeorm": "^9.0.0",
    "class-transformer": "^0.5.1",
    "class-validator": "^0.14.0",
    "pg": "^8.11.0",
    "redis": "^4.6.0",
    "typeorm": "^0.3.16",
    "pino": "^8.16.0",
    "pino-http": "^8.5.0"
  },
  "devDependencies": {
    "@types/jest": "^29.5.0",
    "@types/node": "^20.0.0",
    "jest": "^29.5.0",
    "ts-jest": "^29.1.0",
    "ts-node": "^10.9.0",
    "typescript": "^5.0.0"
  }
}
```

#### Step 3: Create Module Structure (FROM CENG-CARPOOL)

```typescript
// services/_template/src/main.ts
import { NestFactory } from '@nestjs/core';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { Logger } from '@nestjs/common';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Swagger API Documentation
  const config = new DocumentBuilder()
    .setTitle('FamGo Service Template')
    .setDescription('Template service for FamGo Platform')
    .setVersion('1.0')
    .addBearerAuth()
    .build();
  
  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api/docs', app, document);
  
  const port = process.env.PORT || 3000;
  await app.listen(port);
  Logger.log(`Service running on port ${port}`);
}

bootstrap();
```

```typescript
// services/_template/src/app.module.ts
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: process.env.DB_HOST || 'localhost',
      port: parseInt(process.env.DB_PORT) || 5432,
      username: process.env.DB_USER || 'postgres',
      password: process.env.DB_PASSWORD || 'postgres',
      database: process.env.DB_NAME || 'famgo',
      entities: [__dirname + '/**/*.entity{.ts,.js}'],
      synchronize: process.env.NODE_ENV === 'development',
      logging: true,
    }),
  ],
})
export class AppModule {}
```

```typescript
// services/_template/src/modules/example/example.module.ts
import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ExampleController } from './example.controller';
import { ExampleService } from './example.service';
import { Example } from './entities/example.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Example])],
  controllers: [ExampleController],
  providers: [ExampleService],
  exports: [ExampleService],
})
export class ExampleModule {}
```

```typescript
// services/_template/src/modules/example/example.controller.ts
import { Controller, Get, Post, Body, Param } from '@nestjs/common';
import { ApiTags, ApiOperation } from '@nestjs/swagger';
import { ExampleService } from './example.service';
import { CreateExampleDto } from './dtos/create-example.dto';

@ApiTags('example')
@Controller('example')
export class ExampleController {
  constructor(private readonly exampleService: ExampleService) {}

  @Get()
  @ApiOperation({ summary: 'Get all examples' })
  findAll() {
    return this.exampleService.findAll();
  }

  @Post()
  @ApiOperation({ summary: 'Create example' })
  create(@Body() createDto: CreateExampleDto) {
    return this.exampleService.create(createDto);
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get example by ID' })
  findOne(@Param('id') id: string) {
    return this.exampleService.findOne(id);
  }
}
```

```typescript
// services/_template/src/modules/example/example.service.ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Example } from './entities/example.entity';
import { CreateExampleDto } from './dtos/create-example.dto';

@Injectable()
export class ExampleService {
  constructor(
    @InjectRepository(Example)
    private readonly repository: Repository<Example>,
  ) {}

  async findAll(): Promise<Example[]> {
    return this.repository.find();
  }

  async create(dto: CreateExampleDto): Promise<Example> {
    const entity = this.repository.create(dto);
    return this.repository.save(entity);
  }

  async findOne(id: string): Promise<Example> {
    return this.repository.findOneBy({ id });
  }
}
```

```typescript
// services/_template/src/modules/example/entities/example.entity.ts
import { Entity, PrimaryGeneratedColumn, Column, CreateDateColumn } from 'typeorm';

@Entity()
export class Example {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  name: string;

  @CreateDateColumn()
  createdAt: Date;
}
```

```typescript
// services/_template/src/modules/example/dtos/create-example.dto.ts
import { IsString, IsNotEmpty } from 'class-validator';

export class CreateExampleDto {
  @IsString()
  @IsNotEmpty()
  name: string;
}
```

#### Step 4: Create Dockerfile

```dockerfile
# services/_template/Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY --from=builder /app/dist ./dist

EXPOSE 3000

CMD ["node", "dist/main.js"]
```

#### Step 5: Create Makefile

```makefile
# services/_template/Makefile
.PHONY: install build start dev test clean docker

install:
	npm install

build:
	npm run build

start:
	npm start

dev:
	npm run start:dev

test:
	npm run test

test-watch:
	npm run test:watch

migration-generate:
	npm run migration:generate

migration-run:
	npm run migration:run

docker-build:
	docker build -t famgo/template:latest .

docker-run:
	docker run -p 3000:3000 famgo/template:latest

clean:
	rm -rf dist node_modules
```

#### Step 6: Environment Configuration

```bash
# services/_template/.env.example
NODE_ENV=development
PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=famgo

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Kafka
KAFKA_BROKERS=localhost:9092

# Observability
LOG_LEVEL=debug
JAEGER_ENABLED=true
```

---

### Week 2 Task: Extract First Real Service (AUTH SERVICE)

**Source:** Ceng-Carpool `backend/src/modules/auth/`

#### Extract Auth Service Structure

```typescript
// services/auth-service/src/modules/auth/auth.service.ts
// FROM: Ceng-Carpool auth service + ORider KYC concepts

import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { UserService } from '../user/user.service';

@Injectable()
export class AuthService {
  constructor(
    private readonly jwtService: JwtService,
    private readonly userService: UserService,
  ) {}

  // FROM CENG-CARPOOL: JWT token generation
  async generateTokens(userId: string) {
    const payload = { sub: userId };
    
    const accessToken = this.jwtService.sign(payload, {
      expiresIn: '15m',
    });
    
    const refreshToken = this.jwtService.sign(payload, {
      expiresIn: '7d',
      secret: process.env.JWT_REFRESH_SECRET,
    });
    
    return { accessToken, refreshToken };
  }

  // FROM CENG-CARPOOL: WeChat OAuth2
  async wechatLogin(code: string) {
    // 1. Exchange code for access token
    const wechatTokenResponse = await fetch(
      `https://api.weixin.qq.com/sns/oauth2/access_token`,
      {
        params: {
          appid: process.env.WECHAT_APPID,
          secret: process.env.WECHAT_SECRET,
          code,
          grant_type: 'authorization_code',
        },
      },
    );

    const { access_token, openid } = await wechatTokenResponse.json();

    // 2. Get user info
    const userInfoResponse = await fetch(
      `https://api.weixin.qq.com/sns/userinfo`,
      {
        params: { access_token, openid },
      },
    );

    const userInfo = await userInfoResponse.json();

    // 3. Find or create user
    let user = await this.userService.findByWechatId(openid);
    if (!user) {
      user = await this.userService.create({
        wechatId: openid,
        name: userInfo.nickname,
        avatar: userInfo.headimgurl,
      });
    }

    // 4. Generate tokens
    return this.generateTokens(user.id);
  }

  // FROM ORIDER: KYC/Real name verification
  async verifyKYC(userId: string, idNumber: string, fullName: string) {
    // Call KYC provider (e.g., Shuvi, IDify)
    const kycResult = await this.kycProvider.verify({
      idNumber,
      fullName,
    });

    if (kycResult.verified) {
      await this.userService.updateKYCStatus(userId, 'verified');
      return { verified: true };
    }

    return { verified: false };
  }

  // NEW: Device fingerprinting (safety)
  async registerDevice(userId: string, fingerprint: string) {
    // Store device fingerprint for fraud detection
    await this.userService.addDevice(userId, {
      fingerprint,
      registeredAt: new Date(),
      trusted: false,
    });
  }
}
```

---

## EXTRACTING DISPATCH SERVICE (WEEKS 9-10)

**Source:** Carpooling Platform `SVCOrderTempGrab` + DriveMind scoring

### Extract Concurrent Matching Algorithm

```go
// services/dispatch-service/internal/domain/matching.go
// FROM: Carpooling Platform + DriveMind scoring

package domain

import (
  "context"
  "errors"
  "time"
  "github.com/go-redis/redis/v8"
)

type Ride struct {
  ID              string
  PickupLocation  Location
  DropoffLocation Location
  PassengerID     string
  Status          string // pending, matched, accepted, started
  CreatedAt       time.Time
}

type Driver struct {
  ID              string
  CurrentLocation Location
  Rating          float32
  VehicleType     string
}

type MatchingEngine struct {
  redis  *redis.Client
  logger Logger
}

// FROM CARPOOLING PLATFORM: Concurrent grab + temporary lock
func (m *MatchingEngine) MatchRide(ctx context.Context, ride *Ride) (driverId string, err error) {
  // Step 1: Find eligible drivers
  drivers, err := m.findEligibleDrivers(ctx, ride)
  if err != nil {
    return "", err
  }

  // Step 2: Score each driver (FROM DRIVEMIND pattern)
  scores := m.scoreDrivers(drivers, ride)

  // Step 3: Select top driver
  topDriver := scores[0].driver

  // Step 4: Create temporary lock (FROM CARPOOLING)
  // This replaces Java's synchronized block with Redis Redlock
  lockKey := fmt.Sprintf("ride_grab:%s", ride.ID)
  lockValue := topDriver.ID
  
  // Attempt to acquire lock (30 second timeout)
  acquired := m.redis.SetNX(ctx, lockKey, lockValue, 30*time.Second).Val()
  if !acquired {
    return "", errors.New("ride already grabbed")
  }

  // Step 5: Send push notification (Firebase)
  err = m.notificationService.SendRideOffer(ctx, topDriver.ID, ride)
  if err != nil {
    m.redis.Del(ctx, lockKey) // Release lock on error
    return "", err
  }

  // Step 6: Wait for driver response (30 second window)
  // If no response within 30s, lock auto-expires in Redis
  // Other drivers can grab

  return topDriver.ID, nil
}

// FROM DRIVEMIND: Scoring formula
type MatchScore struct {
  driver Driver
  score  float32
}

func (m *MatchingEngine) scoreDrivers(drivers []Driver, ride *Ride) []MatchScore {
  scores := make([]MatchScore, len(drivers))
  
  for i, driver := range drivers {
    // Weight formula:
    // 40% - Distance to pickup (lower is better)
    // 30% - ETA (lower is better)
    // 20% - Driver rating (higher is better)
    // 10% - Preference match
    
    distScore := 1 - (getDistance(driver.CurrentLocation, ride.PickupLocation) / 5000)
    etaScore := 1 - (getETA(driver, ride) / 600) // 10 min max
    ratingScore := driver.Rating / 5.0
    
    totalScore := (distScore * 0.4) + (etaScore * 0.3) + (ratingScore * 0.2)
    
    scores[i] = MatchScore{driver, totalScore}
  }
  
  // Sort by score (descending)
  sort.Slice(scores, func(i, j int) bool {
    return scores[i].score > scores[j].score
  })
  
  return scores
}

func (m *MatchingEngine) findEligibleDrivers(ctx context.Context, ride *Ride) ([]Driver, error) {
  // FROM DRIVEMIND: Redis GEO index
  // Find drivers near pickup within 2km radius
  location := ride.PickupLocation
  
  drivers, err := m.geoService.FindNearby(ctx, location, 2000) // 2km
  if err != nil {
    return nil, err
  }
  
  // Filter: only online drivers
  eligible := make([]Driver, 0)
  for _, driver := range drivers {
    if driver.Status == "online" {
      eligible = append(eligible, driver)
    }
  }
  
  return eligible, nil
}
```

---

## EXTRACTING POOLING SERVICE (WEEK 10)

**Source:** Ceng-Carpool smart allocation + DriveMind routing

### Smart Pool Matching Algorithm

```go
// services/pooling-service/internal/domain/pool_matcher.go
// FROM: Ceng-Carpool smart allocation + DriveMind route optimization

package domain

import (
  "math"
)

type PoolMatcher struct {
  routeService RouteService // For polyline overlapping
}

// Matching rules from FamGo spec
const (
  MAX_DETOUR_MINUTES    = 10
  MAX_PICKUP_RADIUS_KM  = 2
  MIN_ROUTE_OVERLAP_PCT = 0.70
  MAX_POOL_SIZE         = 3
)

type PoolScore struct {
  Ride1           *Ride
  Ride2           *Ride
  Score           float32
  RouteOverlapPct float32
  TotalDetour     float32
  EstimatedProfit float32
}

// FROM CENG-CARPOOL: Smart allocation algorithm
func (pm *PoolMatcher) FindBestPool(newRide *Ride, existingRides []*Ride) (*Ride, *PoolScore, error) {
  var bestScore *PoolScore
  var bestRide *Ride
  
  for _, ride := range existingRides {
    // Check if poolable
    if !pm.canPool(newRide, ride) {
      continue
    }
    
    // Calculate pool score
    score := pm.calculatePoolScore(newRide, ride)
    
    if bestScore == nil || score.Score > bestScore.Score {
      bestScore = score
      bestRide = ride
    }
  }
  
  return bestRide, bestScore, nil
}

// FROM CENG-CARPOOL: Pooling rules
func (pm *PoolMatcher) canPool(ride1, ride2 *Ride) bool {
  // Rule 1: Same direction (70% overlap minimum)
  overlap := pm.calculateRouteOverlap(ride1, ride2)
  if overlap < MIN_ROUTE_OVERLAP_PCT {
    return false
  }
  
  // Rule 2: Pickup detour within 10 minutes
  pickupDetour := pm.calculatePickupDetour(ride1, ride2)
  if pickupDetour > MAX_DETOUR_MINUTES {
    return false
  }
  
  // Rule 3: Both pickups within 2km
  pickupDist := getDistance(ride1.PickupLocation, ride2.PickupLocation)
  if pickupDist > MAX_PICKUP_RADIUS_KM*1000 {
    return false
  }
  
  // Rule 4: Pool size doesn't exceed 3
  if len(ride1.Passengers)+len(ride2.Passengers) > MAX_POOL_SIZE {
    return false
  }
  
  return true
}

// FROM CENG-CARPOOL: Scoring formula
// score = (overlap × 0.4) + (profitability × 0.3) + (eta_similarity × 0.2) + (distance × 0.1)
func (pm *PoolMatcher) calculatePoolScore(ride1, ride2 *Ride) *PoolScore {
  overlap := pm.calculateRouteOverlap(ride1, ride2)
  profitability := pm.calculateProfitability(ride1, ride2) // Relative to solo
  etaSimilarity := pm.calculateETASimilarity(ride1, ride2)
  pickupDistance := getDistance(ride1.PickupLocation, ride2.PickupLocation)
  
  score := (overlap * 0.4) +
    (profitability * 0.3) +
    (etaSimilarity * 0.2) +
    ((1 - (pickupDistance / 2000)) * 0.1) // Normalize distance
  
  return &PoolScore{
    Ride1:           ride1,
    Ride2:           ride2,
    Score:           score,
    RouteOverlapPct: overlap,
  }
}

func (pm *PoolMatcher) calculateRouteOverlap(ride1, ride2 *Ride) float32 {
  // Use polyline matching
  // Compare ride1 route to ride2 route
  overlap := pm.routeService.PolylineOverlap(
    ride1.PickupLocation,
    ride1.DropoffLocation,
    ride2.PickupLocation,
    ride2.DropoffLocation,
  )
  return float32(overlap) // 0.0 to 1.0
}

func (pm *PoolMatcher) calculateProfitability(ride1, ride2 *Ride) float32 {
  // More passengers = better utilization = higher profitability
  // Compare: solo fare vs pooled fare
  soloFare := ride1.FareCalculation.TotalAmount
  poolFarePerPassenger := (ride1.FareCalculation.TotalAmount + ride2.FareCalculation.TotalAmount) / 2
  
  profitability := poolFarePerPassenger / soloFare
  return min(profitability, 1.0) // Cap at 1.0
}

func (pm *PoolMatcher) calculateETASimilarity(ride1, ride2 *Ride) float32 {
  // If ETAs are similar, they're good candidates
  etaDiff := math.Abs(float64(ride1.ETA - ride2.ETA))
  
  // If within 5 minutes, perfect match
  if etaDiff < 5*60 {
    return 1.0
  }
  
  // Otherwise scale down
  return float32(1.0 - (etaDiff / 600)) // 10 min max
}
```

---

## EXTRACTING PAYMENT SERVICE (WEEK 11)

**Source:** ORider smart contract escrow pattern + Ceng-Carpool pricing

### Immutable Wallet Ledger Implementation

```go
// services/wallet-service/internal/domain/ledger.go
// FROM: ORider smart contract escrow + Ceng-Carpool payment modes

package domain

import (
  "crypto/sha256"
  "encoding/hex"
  "time"
)

type WalletTransaction struct {
  ID              string    // Unique UUID
  WalletID        string
  Amount          decimal.Decimal
  Type            string    // "deposit", "ride_payment", "driver_earning", "refund", "bonus"
  RideID          string    // Links to specific ride
  Status          string    // "pending", "locked", "committed", "failed"
  
  // FROM ORIDER: GPS verification
  GPSVerified     bool
  PickupGPS       Location
  DropoffGPS      Location
  ActualEndGPS    Location
  DistanceToEnd   float32 // meters
  
  // FROM ORIDER: Immutability
  PreviousHash    string    // Hash of previous transaction (blockchain-like)
  CurrentHash     string    // Hash of this transaction
  
  CreatedAt       time.Time
  VerifiedAt      time.Time
  CommittedAt     time.Time
}

// FROM ORIDER: Immutable insert pattern
// NEVER UPDATE transactions, only INSERT
func (lt *LedgerStore) InsertTransaction(tx *WalletTransaction) error {
  // Calculate transaction hash
  tx.CurrentHash = lt.calculateHash(tx)
  
  // Validate chain (get previous transaction)
  prevTx, err := lt.getLastTransaction(tx.WalletID)
  if err == nil {
    tx.PreviousHash = prevTx.CurrentHash
  }
  
  // Insert (never update)
  err = lt.db.Create(tx).Error
  if err != nil {
    return err
  }
  
  return nil
}

func (lt *LedgerStore) calculateHash(tx *WalletTransaction) string {
  // Create hash of transaction data
  data := fmt.Sprintf("%s:%s:%s:%s:%d:%s",
    tx.WalletID,
    tx.Amount.String(),
    tx.Type,
    tx.RideID,
    tx.CreatedAt.Unix(),
    tx.PreviousHash,
  )
  
  hash := sha256.Sum256([]byte(data))
  return hex.EncodeToString(hash[:])
}

// FROM ORIDER: GPS-based payment verification
type PaymentEscrow struct {
  ID              string
  RideID          string
  DriverID        string
  PassengerID     string
  Amount          decimal.Decimal
  Status          string // "locked", "verified", "released", "refunded"
  
  PickupGPS       Location
  DropoffGPS      Location
  ActualEndGPS    Location
  
  LockedAt        time.Time
  VerificationKey string // One-time key for driver
}

// FROM ORIDER: Verify trip completion by GPS
func (pe *PaymentEscrow) VerifyCompletion(actualEndGPS Location) (bool, error) {
  // Check if within 100m of destination
  distance := getDistance(pe.DropoffGPS, actualEndGPS)
  
  if distance <= 100 { // 100m threshold
    pe.Status = "verified"
    pe.ActualEndGPS = actualEndGPS
    return true, nil
  }
  
  return false, errors.New("destination not reached")
}

// FROM ORIDER: Release payment to driver
func (pe *PaymentEscrow) ReleaseToDriver(ledger *LedgerStore) error {
  // Create immutable ledger entry
  driverTx := &WalletTransaction{
    ID: uuid.New().String(),
    WalletID: pe.DriverID,
    Amount: pe.Amount,
    Type: "driver_earning",
    RideID: pe.RideID,
    Status: "committed",
    GPSVerified: true,
    PickupGPS: pe.PickupGPS,
    DropoffGPS: pe.DropoffGPS,
    ActualEndGPS: pe.ActualEndGPS,
    CreatedAt: time.Now(),
  }
  
  // Insert to ledger (never update)
  return ledger.InsertTransaction(driverTx)
}

// FROM ORIDER: Auto-refund incomplete trip
func (pe *PaymentEscrow) RefundPassenger(ledger *LedgerStore, reason string) error {
  passengerTx := &WalletTransaction{
    ID: uuid.New().String(),
    WalletID: pe.PassengerID,
    Amount: pe.Amount,
    Type: "refund",
    RideID: pe.RideID,
    Status: "committed",
    CreatedAt: time.Now(),
  }
  
  // Insert to ledger
  return ledger.InsertTransaction(passengerTx)
}
```

---

## WEEK 12: SAFETY SERVICE

**Source:** DriveMind ML detection + ORider GPS verification

```go
// services/safety-service/internal/domain/safety_detector.go
// FROM: DriveMind ML patterns + ORider verification

package domain

type TripAnomalyDetector struct {
  ml *MLClient // For DriveMind-style detection
}

// FROM DRIVEMIND: Detect route deviation
func (td *TripAnomalyDetector) DetectRouteDeviation(
  trip *Trip,
  currentGPS Location,
) (bool, float32, error) {
  // Get expected route
  expectedPolyline := trip.ExpectedRoute
  
  // Calculate deviation from route
  deviation := polylineDeviation(currentGPS, expectedPolyline)
  
  // Threshold: if >500m off route, flag
  if deviation > 500 {
    return true, deviation, nil
  }
  
  return false, deviation, nil
}

// FROM DRIVEMIND: Harsh braking detection
func (td *TripAnomalyDetector) DetectHarshBraking(acceleration float32) bool {
  // Acceleration < -6 m/s² is harsh braking
  return acceleration < -6.0
}

// FROM DRIVEMIND: Speed monitoring
func (td *TripAnomalyDetector) DetectExcessiveSpeed(
  speed float32,
  speedLimit float32,
) bool {
  // If speed > limit + 10 km/h, flag
  return speed > (speedLimit + 10)
}

// NEW: SOS panic button
type EmergencyRequest struct {
  ID        string
  TripID    string
  UserID    string
  Location  Location
  CreatedAt time.Time
  Status    string // "pending", "received", "en_route", "resolved"
}

func (td *TripAnomalyDetector) HandleSOS(trip *Trip, userID string) error {
  sos := &EmergencyRequest{
    ID: uuid.New().String(),
    TripID: trip.ID,
    UserID: userID,
    Location: trip.CurrentLocation,
    CreatedAt: time.Now(),
    Status: "pending",
  }
  
  // 1. Store SOS
  // 2. Notify emergency contacts
  // 3. Alert nearest police/hospital
  
  return nil
}
```

---

## PHASE SUMMARY TABLE

| Phase | Week | Service | From Repo | Status |
|-------|------|---------|-----------|--------|
| 0 | 1-2 | Template + SDKs | Ceng-Carpool | Foundation |
| 1 | 3-4 | Auth, Database | Ceng-Carpool + ORider | Security |
| 1 | 5-6 | Observability | DriveMind | Visibility |
| 2 | 7-8 | User, GPS, Ride | Ceng-Carpool + DriveMind | Core |
| 3 | 9-10 | Dispatch, Pooling | Carpooling + Ceng-Carpool | Matching |
| 4 | 11-12 | Payment, Safety | ORider + DriveMind | Security |
| 5 | 13-16 | Subscriptions, ML, Infra | Ceng-Carpool + DriveMind | Advanced |

---

## CRITICAL SUCCESS FACTORS

1. **Week 1-2:** Get NestJS template right (foundation for all services)
2. **Week 3-4:** Implement auth + database (security baseline)
3. **Week 7-8:** Core ride service working (now you can test E2E)
4. **Week 9-10:** Matching algorithm (hardest part, test thoroughly)
5. **Week 11-12:** Payments working + safety (money + safety = trust)
6. **Week 13-16:** ML + production infrastructure (scale ready)

---

**Next:** Start implementing Week 1 NestJS template today!
