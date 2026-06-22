# FamGo Service Template

## Overview

This is the standardized NestJS service template for all FamGo platform microservices. Every service MUST follow this structure to ensure:

- Consistency across the platform
- Easy onboarding for new developers
- Standardized deployment and testing
- Enterprise-grade observability

## Directory Structure

```
services/_template/
├── src/
│   ├── main.ts                 # Service entry point
│   ├── app.module.ts           # Root module
│   ├── config/                 # Configuration management
│   ├── modules/                # Business logic modules
│   │   └── example/
│   │       ├── example.module.ts
│   │       ├── example.controller.ts
│   │       ├── example.service.ts
│   │       ├── entities/
│   │       │   └── example.entity.ts
│   │       ├── dtos/
│   │       │   ├── create-example.dto.ts
│   │       │   └── update-example.dto.ts
│   │       ├── repositories/
│   │       │   └── example.repository.ts
│   │       └── tests/
│   │           └── example.service.spec.ts
│   ├── common/
│   │   ├── filters/            # Exception filters
│   │   ├── guards/             # Auth guards
│   │   ├── interceptors/       # HTTP interceptors
│   │   ├── middleware/         # HTTP middleware
│   │   └── decorators/         # Custom decorators
│   ├── database/
│   │   ├── migrations/         # Database migrations
│   │   └── seeds/              # Seed data
│   └── shared/
│       ├── constants/
│       ├── interfaces/
│       ├── utils/
│       └── types/
├── test/
│   └── app.e2e-spec.ts
├── package.json
├── tsconfig.json
├── Dockerfile
├── .dockerignore
├── .gitignore
├── .env.example
├── Makefile
└── README.md
```

## Quick Start

### Installation

```bash
# Create new service from template
cp -r services/_template services/my-new-service
cd services/my-new-service
npm install
```

### Development

```bash
# Start in development mode with hot reload
npm run start:dev

# Run tests
npm test

# Run tests with coverage
npm test:cov

# Watch mode
npm run test:watch
```

### Database

```bash
# Generate new migration
npm run migration:generate -- -n CreateMyTable

# Run migrations
npm run migration:run

# Revert last migration
npm run migration:revert
```

### Building

```bash
# Build for production
npm run build

# Start production build
npm start
```

## Architecture

### Modules

Each business domain is a NestJS Module containing:
- **Controller**: HTTP endpoints
- **Service**: Business logic
- **Entity**: Database model
- **Repository**: Data access (optional, use TypeORM repository)
- **DTO**: Request/Response schemas
- **Tests**: Unit tests

### Layering

```
Controller (HTTP Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
Entity (Database)
```

## Key Principles

1. **Single Responsibility**: Each class has one reason to change
2. **Dependency Injection**: Use NestJS DI everywhere
3. **Type Safety**: Strict TypeScript, no `any` types
4. **Testing**: Minimum 80% coverage
5. **Documentation**: API docs via Swagger
6. **Error Handling**: Standardized exception handling
7. **Logging**: Structured logging (Pino)
8. **Validation**: DTO validation with class-validator

## Configuration

Configuration is environment-based, using `.env` files.

### .env Example

```
NODE_ENV=development
PORT=3000
SERVICE_NAME=template-service

# Database
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=famgo

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Observability
LOG_LEVEL=debug
JAEGER_ENABLED=false

# JWT
JWT_SECRET=your-secret-here
JWT_EXPIRATION=15m
JWT_REFRESH_SECRET=refresh-secret-here
JWT_REFRESH_EXPIRATION=7d
```

## Observability

Services include built-in observability:

- **Logging**: Pino structured logging
- **Metrics**: Prometheus (via decorators)
- **Tracing**: Jaeger (optional)
- **Health Checks**: Liveness/readiness probes

## Testing

```bash
# Unit tests
npm test

# Watch mode
npm run test:watch

# Coverage report
npm run test:cov
```

## Docker

```bash
# Build image
docker build -t famgo/template-service:latest .

# Run container
docker run -p 3000:3000 \
  -e DB_HOST=postgres \
  -e REDIS_HOST=redis \
  famgo/template-service:latest
```

## API Documentation

After starting the service, view Swagger docs at:
```
http://localhost:3000/api/docs
```

## Deployment

Services are deployed via Kubernetes. See `/infra/kubernetes/` for manifests.

## Common Tasks

### Adding a New Module

```bash
# Generate module scaffolding
nest generate module modules/my-module
nest generate controller modules/my-module
nest generate service modules/my-module
```

### Adding Database Entity

1. Create entity in `src/modules/my-module/entities/`
2. Add to `TypeOrmModule.forFeature()`
3. Generate migration: `npm run migration:generate -- -n AddMyTable`
4. Run migration: `npm run migration:run`

### Calling Another Service

Use service-to-service gRPC or REST via dedicated SDK in `packages/`.

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Port already in use | Change PORT in .env or kill existing process |
| Database connection error | Verify DB_HOST, DB_PORT, DB_USER, DB_PASSWORD |
| Module not found | Check imports in app.module.ts |
| Migration fails | Check database permissions and SQL syntax |

## Contributing

1. Follow NestJS best practices
2. Maintain >80% test coverage
3. Update API docs when changing endpoints
4. Add database migrations for schema changes
5. Use conventional commits: `feat:`, `fix:`, `docs:`, etc.

## References

- [NestJS Documentation](https://docs.nestjs.com/)
- [TypeORM Documentation](https://typeorm.io/)
- [Class Validator](https://github.com/typestack/class-validator)
- [Pino Logger](https://getpino.io/)

---

**Version:** 1.0  
**Last Updated:** 2025  
**Maintained by:** FamGo Platform Team
