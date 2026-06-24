# UNSIA Frontend Development Plan

## Overview
This document outlines the frontend development plan for UNSIA ERP system, building modules in sequence: SSO → Reference → PMB → Finance → Academic → LMS.

## Module Dependencies

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend Portal                    │
│                  (Next.js 14 App)                      │
└─────────────────────┬───────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
    ┌────────┐ ┌─────────┐ ┌──────────┐
    │   SSO  │ │   SSO   │ │  Other  │
    │ Login │ │ Session │ │ Modules │
    └───┬────┘ └────┬───┘ └───┬────┘
        │           │        │
        │   ┌──────▼──────┐  │
        │  │ Reference  │  │
        │  │   Data    │──┘
        │  └──────┬─────┘
        │         │
        │   ┌────▼─────────────────────┐
        │  │        PMB Module      │──► Finance
        │  └──────┬─────────────────┘
        │         │
        │   ┌────▼──────────────────────────┐
        │  │      Academic Module           │──► LMS
        │  └──────────────────────────────┘
        │
        │ Note: HRIS, Assessment, CRM, Portal can be added later
```

## Current Status

### ✅ Completed
- [x] Login page UI (unsia-portal-web/app/(auth)/login/page.tsx)
- [x] Auth context with SSO integration (contexts/auth-context.tsx)
- [x] API utilities (lib/api.ts)
- [x] Constants for API endpoints (lib/constants.ts)
- [x] Portal layout (app/(portal)/layout.tsx)
- [x] Dashboard page (app/(portal)/dashboard/page.tsx)

### 🔄 In Progress
- [ ] SSO Login API integration with real backend

### 📋 To Do

#### Phase 1: SSO & Auth (Priority: HIGH)
- [ ] Login API connection to unsia-core-service
- [ ] Token refresh logic
- [ ] Role switching UI
- [ ] Logout functionality
- [ ] Protected routes

#### Phase 2: Reference Data (Priority: HIGH)
- [ ] Reference API hooks
- [ ] Study Programs list
- [ ] Academic Years list
- [ ] Academic Periods list
- [ ] Payment Components list
- [ ] Document Types list
- [ ] Payment Methods list
- [ ] PMB Waves list (for PMB)
- [ ] Provinces/Cities/Districts/Villages (region data)

#### Phase 3: PMB Module (Priority: HIGH)
- [ ] PMB Dashboard
- [ ] Wave Management
- [ ] Applicant Registration Form
- [ ] Document Upload
- [ ] Payment Verification
- [ ] Selection Results
- [ ] Applicant List (Admin)

#### Phase 4: Finance Module (Priority: MEDIUM)
- [ ] Finance Dashboard
- [ ] Invoice Management
- [ ] Payment Records
- [ ] Student Clearance

#### Phase 5: Academic Module (Priority: MEDIUM)
- [ ] Academic Dashboard
- [ ] Student Management
- [ ] Course Management
- [ ] KRS (Enrollment)
- [ ] Grades
- [ ] Schedules

#### Phase 6: LMS Module (Priority: MEDIUM)
- [ ] LMS Dashboard
- [ ] Course List
- [ ] Session/Material Browser
- [ ] Assignments
- [ ] Attendance

## API Mapping

### unsia-core-service (Port 8001) - Auth/SSO
| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/auth/login | POST | Login with username/password |
| /api/v1/auth/refresh | POST | Refresh access token |
| /api/v1/auth/me | GET | Get current user info |
| /api/v1/auth/switch-role | POST | Switch active role |

### unsia-reference-service (Port 8007) - Reference Data
| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/reference/study-programs | GET | List study programs |
| /api/v1/reference/academic-years | GET | List academic years |
| /api/v1/reference/academic-periods | GET | List academic periods |
| /api/v1/reference/payment-components | GET | List payment components |
| /api/v1/reference/payment-methods | GET | List payment methods |
| /api/v1/reference/document-types | GET | List document types |
| /api/v1/reference/pmb-waves | GET | List PMB waves |
| /api/v1/reference/provinces | GET | List provinces |
| /api/v1/reference/cities | GET | List cities |
| /api/v1/reference/districts | GET | List districts |
| /api/v1/reference/villages | GET | List villages |

### unsia-pmb-service (Port 8003) - PMB
| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/pmb/applicants | GET/POST | List/Create applicants |
| /api/v1/pmb/applicants/:id | GET/PUT | Get/Update applicant |
| /api/v1/pmb/waves | GET/POST | List/Create waves |
| /api/v1/pmb/selection/results | GET | Selection results |

### unsia-finance-service (Port 8002) - Finance
| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/finance/invoices | GET/POST | List/Create invoices |
| /api/v1/finance/payments | GET | List payments |
| /api/v1/finance/clearance | GET | Student clearance |

### unsia-academic-service (Port 8004) - Academic
| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/academic/students | GET/POST | List/Create students |
| /api/v1/academic/krs | GET | Student KRS |
| /api/v1/academic/schedules | GET | Class schedules |
| /api/v1/academic/grades | GET | Student grades |

### unsia-lms-service (Port 8005) - LMS
| Endpoint | Method | Description |
|----------|--------|-------------|
| /api/v1/lms/courses | GET | List courses |
| /api/v1/lms/sessions | GET | Course sessions |
| /api/v1/lms/assignments | GET | Course assignments |

## Frontend File Structure

```
unsia-portal-web/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   │   └── page.tsx         ✅ Login page
│   │   ├── register/
│   │   │   └── page.tsx         ⬜ Registration
│   │   └── layout.tsx
│   ├── (portal)/
│   │   ├── layout.tsx            ✅ Portal layout
│   │   ├── dashboard/
│   │   │   └── page.tsx         ✅ Dashboard
│   │   ├── pmb/
│   │   │   ├── applicants/
│   │   │   │   └── page.tsx
│   │   │   ├── waves/
│   │   │   │   └── page.tsx
│   │   │   └── page.tsx
│   │   ├── finance/
│   │   │   ├── invoices/
│   │   │   │   └── page.tsx
│   │   │   └── page.tsx
│   │   ├── academic/
│   │   │   ├── students/
│   │   │   │   └── page.tsx
│   │   │   ├── courses/
│   │   │   │   └── page.tsx
│   │   │   └── page.tsx
│   │   └── lms/
│   │       ├── courses/
│   │       │   └── page.tsx
│   │       └── page.tsx
│   ├── globals.css                  ✅ Global styles
│   ├── layout.tsx
│   └── page.tsx                   → Redirect to login
├── components/
│   ├── ui/                       ✅ UI components
│   ├── layout/
│   │   ├── sidebar.tsx
│   │   ├── header.tsx
│   │   └── footer.tsx
│   └── features/
│       ├── pmb/
│       ├── finance/
│       ├── academic/
│       └── lms/
├── contexts/
│   ├── auth-context.tsx          ✅ Auth context
│   └── reference-context.tsx       ⬜ Reference data
├── hooks/
│   ├── use-auth.ts
│   ├── use-reference.ts
│   ├── use-pmb.ts
│   ├── use-finance.ts
│   ├── use-academic.ts
│   └── use-lms.ts
├── lib/
│   ├── api.ts                   ✅ API utilities
│   ├── constants.ts             ✅ API endpoints
│   └── utils.ts
└── types/
    ├── user.ts
    ├── reference.ts
    ├── pmb.ts
    ├── finance.ts
    ├── academic.ts
    └── lms.ts
```

## Key Design Decisions

1. **Route Structure**: Use Next.js 14 App Router with route groups
2. **Authentication**: JWT tokens with refresh token pattern
3. **State Management**: React Context + TanStack Query for server state
4. **UI Framework**: Custom Tailwind + Radix UI primitives
5. **Forms**: React Hook Form + Zod validation
6. **API Pattern**: Centralized API utilities with typed responses

## Environment Variables

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_AUTH_API=http://localhost:8001
NEXT_PUBLIC_REFERENCE_API=http://localhost:8007
NEXT_PUBLIC_PMB_API=http://localhost:8003
NEXT_PUBLIC_FINANCE_API=http://localhost:8002
NEXT_PUBLIC_ACADEMIC_API=http://localhost:8004
NEXT_PUBLIC_LMS_API=http://localhost:8005
```

## Next Steps

1. Fix login to use `username` field instead of `email`
2. Create reference data hooks
3. Build PMB module pages
4. Build Finance module pages
5. Build Academic module pages
6. Build LMS module pages

---

*Last Updated: 2024*
*Status: In Development*
