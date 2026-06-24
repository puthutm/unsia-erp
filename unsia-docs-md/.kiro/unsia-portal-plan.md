# Portal Service Implementation Plan

## Service Overview
**Service Name**: unsia-portal-service (Web Portal Backend)
**Port**: 8000 (or shared)
**Database**: portal_db (PostgreSQL) - optional, mostly API aggregation

## Core Features

### 1. Portal Web (Next.js Frontend)
- Student portal (mahasiswa.unsia.ac.id)
- Admin portal (admin.unsia.ac.id)
- Public portal (unsia.ac.id)

### 2. API Gateway
- Route to backend services
- Authentication redirect to Core

### 3. Dashboard
- Aggregated data from all services
- Quick stats

---

## Frontend Structure (Next.js)

```
unsia-portal-web/
├── app/
│   ├── (auth)/        # Login, register
│   ├── (portal)/     # Protected routes
│   │   ├── dashboard/
│   │   ├── pmb/
│   │   ├── akademik/
│   │   ├── keuangan/
│   │   └── ...
│   └── api/         # Route handlers
├── components/
├── lib/
└── services/        # API clients
```

---

## Implementation Timeline: 8-10 weeks
