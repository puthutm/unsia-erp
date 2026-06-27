# FE Standalone Plan

## Overview
This document outlines the plan to make each frontend (FE) module stand alone per container in the UNSIA ERP system.

## Background
The current architecture has a monolithic frontend (unsia-portal-web) that serves all modules. To improve scalability, maintainability, and deployment flexibility, each frontend module should run in its own container.

## Goals
1. Each frontend module runs in its own Docker container
2. Independent scaling per module
3. Independent deployment per module
4. Isolated failures (one module going down doesn't affect others)
5. Better resource allocation

## Frontend Modules

| Module | Service Name | Container Name | Port | Status |
|--------|------------|---------------|------|--------|
| unsia-portal-web | portal-web | unsia-portal-web | 3000 | Existing |
| unsia-pmb | pmb-web | unsia-pmb-web | 3001 | ✅ Done |
| unsia-academic | academic-web | unsia-academic-web | 3002 | ✅ Done |
| unsia-finance | finance-web | unsia-finance-web | 3003 | ✅ Done |
| unsia-lms | lms-web | unsia-lms-web | 3004 | ✅ Done |
| unsia-hris | hris-web | unsia-hris-web | 3005 | ✅ Done |
| unsia-assessment | assessment-web | unsia-assessment-web | 3006 | ✅ Done |
| unsia-crm | crm-web | unsia-crm-web | 3007 | ✅ Done |
| unsia-reference | reference-web | unsia-reference-web | 3008 | ✅ Done |
| unsia-core | core-web | unsia-core-web | 3009 | ✅ Done |

## Implementation

### 1. Directory Structure
Each frontend module has its own directory under `frontend/`:
```
frontend/
├── unsia-pmb/
│   ├── package.json
│   ├── Dockerfile
│   ├── next.config.js
│   └── app/
├── unsia-academic/
│   └── ...
└── ...
```

### 2. Required Files per Module
- `package.json` - Dependencies
- `next.config.js` - Next.js configuration
- `tsconfig.json` - TypeScript configuration
- `tailwind.config.js` - Tailwind CSS configuration
- `postcss.config.js` - PostCSS configuration
- `Dockerfile` - Container definition
- `app/` - Next.js pages and components

### 3. Docker Configuration
Each Dockerfile uses multi-stage build:
- Stage 1: Builder (builds the Next.js app)
- Stage 2: Runner (production runtime)

### 4. Environment Variables
Each container needs environment variables for backend service URLs:
```env
NEXT_PUBLIC_API_CORE_URL=http://unsia-core-service:8001
NEXT_PUBLIC_API_REFERENCE_URL=http://unsia-reference-service:8002
# ... other service URLs as needed
```

## docker-compose.yml Updates
The main docker-compose.yml needs to be updated to include all frontend containers. Each frontend service should:
1. Build from its respective directory
2. Expose its assigned port
3. Connect to the backend services it depends on
4. Be part of the unsia-network

## Backend Service Dependencies

| Frontend | Backend Services Needed |
|----------|----------------------|
| unsia-pmb | unsia-pmb-service, unsia-core-service |
| unsia-academic | unsia-academic-service, unsia-core-service, unsia-reference-service |
| unsia-finance | unsia-finance-service, unsia-core-service, unsia-reference-service |
| unsia-lms | unsia-lms-service, unsia-core-service, unsia-academic-service |
| unsia-hris | unsia-hris-service, unsia-core-service, unsia-finance-service |
| unsia-assessment | unsia-assessment-service, unsia-core-service, unsia-lms-service |
| unsia-crm | unsia-crm-service, unsia-core-service, unsia-reference-service |
| unsia-reference | unsia-reference-service, unsia-core-service |
| unsia-core | unsia-core-service |

## Migration Steps

### Phase 1: Setup (Completed)
1. ✅ Create directory structure for each module
2. ✅ Create package.json for each module
3. ✅ Create Next.js configuration files
4. ✅ Create Dockerfiles

### Phase 2: docker-compose Integration (In Progress)
1. ⏳ Update main docker-compose.yml
2. ⏳ Add all frontend containers
3. ⏳ Configure environment variables

### Phase 3: Testing (Pending)
1. ⏳ Test each container independently
2. ⏳ Verify inter-service communication
3. ⏳ Verify routing works

### Phase 4: Deployment (Pending)
1. ⏳ Deploy to production
2. ⏳ Set up load balancing (if needed)
3. ⏳ Configure monitoring

## Notes
- Each frontend module uses Next.js standalone mode for smaller container images
- All containers use node:20-alpine base image
- Environment variables follow Next.js conventions (NEXT_PUBLIC_*)
