# UNSIA Core Service

Core authentication and authorization service for UNSIA ERP ecosystem.

## Overview

UNSIA Core Service is a comprehensive OAuth2/OpenID Connect server that provides:
- User authentication and session management
- Role-based access control (RBAC)
- Service-to-service authentication
- OAuth2/OIDC implementation
- External application management
- Audit logging
- Webhook management

## Features

### Authentication
- JWT-based authentication
- RS256 token signing
- Token refresh mechanism
- Session management
- Multi-role support
- Impersonation

### Authorization
- Role-based access control (RBAC)
- Permission system
- Scope-based authorization
- Service tokens for internal calls

### User Management
- User CRUD operations
- Role assignment
- User activation/deactivation
- Password management

### External Integrations
- External application management
- OAuth2 client registration
- Webhook system
- Service tokens

### Observability
- Health checks (/health, /ready, /live)
- Request logging
- Correlation ID tracking
- Audit logging

## Architecture

```
cmd/core-service/main.go      - Entry point
internal/
├── domain/models.go         - Domain models
├── handler/                  - HTTP handlers
│   ├── auth_handler.go
│   ├── user_handler.go
│   ├── role_handler.go
│   ├── session_handler.go
│   ├── application_handler.go
│   ├── service_token_handler.go
│   ├── audit_handler.go
│   ├── webhook_handler.go
│   ├── health_handler.go
│   └── external_app_handler.go
├── service/                  - Business logic
│   ├── audit_service.go
│   ├── webhook_service.go
│   ├── pagination.go
│   ├── validation_service.go
│   ├── notification_service.go
│   ├── report_service.go
│   ├── export_service.go
│   ├── cache_service.go
│   ├── queue_service.go
│   ├── config_service.go
│   ├── preferences_service.go
│   └── external_app_service.go
├── middleware/               - HTTP middleware
│   ├── auth_middleware.go
│   ├── service_token_middleware.go
│   ├── ratelimit_middleware.go
│   ├── cors_middleware.go
│   └── logging_middleware.go
├── event/                    - Event system
│   └── event_bus.go
├── router/                   - Router setup
│   └── router.go
├── infrastructure/           - Infrastructure
│   ├── database/
│   └── keys/
└── migrations/               - Database migrations
```

## API Endpoints

### Public Routes
```
POST   /api/v1/auth/login              - User login
POST   /api/v1/auth/refresh            - Refresh token
GET    /.well-known/jwks.json           - JWKS endpoint
GET    /.well-known/openid-configuration - OIDC discovery
GET    /health                         - Health check
GET    /ready                         - Readiness check
```

### Protected Routes (JWT Required)
```
# Auth
POST   /api/v1/auth/logout              - User logout
POST   /api/v1/auth/change-password    - Change password
GET    /api/v1/auth/me                 - Get current user

# Users
GET    /api/v1/users                   - List users
POST   /api/v1/users                   - Create user
GET    /api/v1/users/:id                - Get user
PUT    /api/v1/users/:id                - Update user
DELETE /api/v1/users/:id               - Delete user
POST   /api/v1/users/:id/activate     - Activate user
POST   /api/v1/users/:id/deactivate    - Deactivate user

# Roles
GET    /api/v1/roles                   - List roles
POST   /api/v1/roles                   - Create role
GET    /api/v1/roles/:id                - Get role
PUT    /api/v1/roles/:id                - Update role
DELETE /api/v1/roles/:id               - Delete role
POST   /api/v1/roles/:id/assign        - Assign role
POST   /api/v1/roles/:id/revoke        - Revoke role

# Sessions
GET    /api/v1/sessions                 - List sessions
GET    /api/v1/sessions/:id             - Get session
DELETE /api/v1/sessions/:id            - Delete session
DELETE /api/v1/sessions                - Revoke all sessions

# Applications
GET    /api/v1/applications             - List applications
POST   /api/v1/applications            - Create application
GET    /api/v1/applications/:id        - Get application
PUT    /api/v1/applications/:id        - Update application
DELETE /api/v1/applications/:id       - Delete application

# Service Tokens
GET    /api/v1/service-tokens          - List tokens
POST   /api/v1/service-tokens          - Create token
GET    /api/v1/service-tokens/:id     - Get token
DELETE /api/v1/service-tokens/:id    - Revoke token

# Audit Logs
GET    /api/v1/audit-logs              - List logs
GET    /api/v1/audit-logs/:id         - Get log

# Webhooks
GET    /api/v1/webhooks               - List webhooks
POST   /api/v1/webhooks              - Create webhook
GET    /api/v1/webhooks/:id           - Get webhook
PUT    /api/v1/webhooks/:id          - Update webhook
DELETE /api/v1/webhooks/:id          - Delete webhook
POST   /api/v1/webhooks/:id/test     - Test webhook
POST   /api/v1/webhooks/trigger     - Trigger event

# External Apps
GET    /api/v1/external-apps         - List external apps
POST   /api/v1/external-apps         - Create external app
GET    /api/v1/external-apps/:id    - Get external app
PUT    /api/v1/external-apps/:id     - Update external app
DELETE /api/v1/external-apps/:id     - Deactivate external app
POST   /api/v1/external-apps/:id/secret - Regenerate secret

# Internal (Service-to-Service)
POST   /internal/validate-token       - Validate token
```

## Database Tables

### Core Tables
- `users` - User accounts
- `roles` - Role definitions
- `user_roles` - User-role mapping
- `permissions` - Permission definitions
- `role_permissions` - Role-permission mapping
- `sessions` - User sessions
- `applications` - OAuth2 applications
- `service_tokens` - Service tokens
- `audit_logs` - Audit trail
- `webhooks` - Webhook configurations
- `external_apps` - External applications

## Environment Variables

```env
# Server
PORT=8001
SERVER_HOST=0.0.0.0

# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/core_db

# JWT
JWT_SECRET=your-secret-key

# RSA Keys
RSA_PRIVATE_KEY_PATH=keys/private.pem
RSA_PUBLIC_KEY_PATH=keys/public.pem

# Cache
REDIS_ADDR=localhost:6379

# Observability
JWKS_URL=http://localhost:8001/.well-known/jwks.json
```

## Running

```bash
# With existing makefile
make run-core

# Or directly
go run cmd/core-service/main.go
```

## Testing

```bash
# Run tests
make test-core

# Run with coverage
make test-core-coverage
```

## Dependencies

- `gin` - HTTP framework
- `gorm` - ORM
- `postgres` - Database driver
- `shared-auth` - JWT utilities
- `shared-audit` - Audit logging
- `shared-idempotency` - Idempotency
- `shared-observability` - Logging & metrics
