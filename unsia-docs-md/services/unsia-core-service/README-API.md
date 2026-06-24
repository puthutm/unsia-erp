# UNSIA Core Service API Documentation

## Service Tokens API

### Authentication

Service tokens use `X-Service-Token` header or `Authorization: ServiceToken <token>` for service-to-service authentication.

### Endpoints

#### Create Service Token
```
POST /api/v1/service-tokens
Authorization: Bearer <admin_token> or ServiceToken <service_token>

Request:
{
  "application_id": "uuid",
  "scopes": ["finance:read", "finance:write"],
  "expires_at": "2026-01-01T00:00:00Z" // optional, default 1 year
}

Response (201):
{
  "success": true,
  "data": {
    "id": "uuid",
    "application_id": "uuid",
    "token": "base64_encoded_token",
    "scopes": ["finance:read"],
    "expires_at": "2026-01-01T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### List Service Tokens
```
GET /api/v1/service-tokens?application_id=<uuid>
Authorization: Bearer <admin_token>

Response (200):
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "application_id": "uuid",
      "scopes": "[]",
      "expires_at": "2026-01-01T00:00:00Z",
      "revoked_at": null,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### Revoke Service Token
```
POST /api/v1/service-tokens/:id/revoke
Authorization: Bearer <admin_token>

Response (200):
{
  "success": true,
  "data": {
    "message": "Service token berhasil dicabut"
  }
}
```

#### Rotate Service Token
```
POST /api/v1/service-tokens/:id/rotate
Authorization: Bearer <admin_token>

Response (200):
{
  "success": true,
  "data": {
    "id": "uuid",
    "application_id": "uuid",
    "token": "new_base64_encoded_token",
    "scopes": ["finance:read"],
    "expires_at": "2027-01-01T00:00:00Z",
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

#### Validate Service Token
```
POST /api/v1/service-tokens/validate
X-Service-Token: <token>

Request:
{
  "token": "service_token"
}

Response (200):
{
  "success": true,
  "data": {
    "valid": true,
    "application_id": "uuid",
    "scopes": "[]",
    "expires_at": "2026-01-01T00:00:00Z"
  }
}
```

---

## Audit Logs API

### Endpoints

#### List Audit Logs
```
GET /api/v1/audit-logs
Authorization: Bearer <admin_token>
Query Parameters:
- user_id: filter by user
- actor_user_id: filter by actor
- module: filter by module (e.g., "finance", "academic")
- action: filter by action
- entity_name: filter by entity type
- entity_id: filter by entity ID
- application_id: filter by application
- start_date: filter by start date (RFC3339)
- end_date: filter by end date (RFC3339)
- limit: pagination limit (default 50)
- offset: pagination offset

Response (200):
{
  "success": true,
  "data": [...],
  "total": 100,
  "limit": 50,
  "offset": 0
}
```

#### Get Audit Log by ID
```
GET /api/v1/audit-logs/:id
Authorization: Bearer <admin_token>

Response (200):
{
  "success": true,
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "actor_user_id": "uuid",
    "module": "finance",
    "action": "invoice.create",
    "entity_name": "invoice",
    "entity_id": "uuid",
    "old_value": null,
    "new_value": {...},
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

#### Get Entity Audit Logs
```
GET /api/v1/audit-logs/entity/:entity_name/:entity_id
Authorization: Bearer <admin_token>
Query: ?limit=20

Response (200):
{
  "success": true,
  "data": [...]
}
```

#### Get User Audit Logs
```
GET /api/v1/audit-logs/user/:user_id
Authorization: Bearer <admin_token>
Query: ?limit=20

Response (200):
{
  "success": true,
  "data": [...]
}
```

#### Create Audit Log (Manual)
```
POST /api/v1/audit-logs
Authorization: ServiceToken <service_token>

Request:
{
  "user_id": "uuid",
  "actor_user_id": "uuid",
  "module": "finance",
  "action": "invoice.create",
  "entity_name": "invoice",
  "entity_id": "uuid",
  "new_value": {"amount": 1000000},
  "reason": " Pembuatan invoice baru"
}

Response (201):
{
  "success": true,
  "data": {...}
}
```

---

## Module and Action Standards

### Finance Module
- `finance.invoice.create` - Create invoice
- `finance.invoice.update` - Update invoice
- `finance.invoice.delete` - Delete invoice
- `finance.invoice.paid` - Mark invoice as paid
- `finance.payment.create` - Create payment
- `finance.payment.refund` - Refund payment

### Academic Module
- `academic.student.create` - Create student
- `academic.student.update` - Update student
- `academic.student.enroll` - Enroll student
- `academic.grade.input` - Input grade
- `academic.transcript.request` - Request transcript

### HR Module
- `hr.employee.create` - Create employee
- `hr.employee.update` - Update employee
- `hr.salary.process` - Process salary
- `hr.leave.approve` - Approve leave

### Core Module
- `core.user.login` - User login
- `core.user.logout` - User logout
- `core.user.impersonate` - Impersonate user
- `core.role.assign` - Assign role
- `core.permission.grant` - Grant permission
