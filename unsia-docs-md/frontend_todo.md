# Frontend Implementation TODO

## Priority Sequence:
1. SSO & Portal
2. Reference Data
3. PMB (Admissions)
4. Finance
5. Academic
6. LMS

## Task Breakdown

### Phase 1: SSO & Portal - COMPLETED
- [x] Set up Next.js 14 project structure
- [x] Configure Tailwind CSS with UNSIA design tokens
- [x] Implement shared component library (Button, Input, Select, Card, Badge, DataTable)
- [x] Create utility functions (cn, formatCurrency, etc.)
- [x] Build login page with UNSIA branding
- [x] Create auth context and hooks
- [x] Create main layout with dynamic sidebar navigation (per-module)
- [x] Implement reference context and hooks
- [x] Create module hooks (use-pmb, use-finance, use-academic, use-lms)
- [x] Set up protected routes (middleware)

### Phase 2: Reference Data - IN PROGRESS
- [x] API client configuration
- [ ] Create reference data management pages
- [ ] Dropdown/selector components for prodi, fakultas

### Phase 3: PMB Module (Week 3-5)
- [ ] Dashboard with analytics widgets
- [ ] Applicant management table
- [ ] Document verification interface
- [ ] Payment tracking panel
- [ ] Wave/gelombang management
- [ ] Communication tools

### Phase 4: Finance Module (Week 5-8)
- [ ] Treasury dashboard
- [ ] Student billing interface
- [ ] Invoice management
- [ ] Payment verification (auto-reconciliation)
- [ ] Payment gateway integration (Midtrans)
- [ ] Payroll disbursement workflow
- [ ] Tax calculation interface
- [ ] Reporting (Neraca, L/R, Arus Kas)
- [ ] Budget management
- [ ] Journal/Buku Besar viewer

### Phase 5: Academic Module (Week 8-11)
- [ ] Student lifecycle dashboard
- [ ] Student management CRUD
- [ ] KRS interface
- [ ] Schedule management
- [ ] Grade entry
- [ ] Academic transcript viewer
- [ ] Graduation workflow

### Phase 6: LMS Module (Week 11-14)
- [ ] Course catalog
- [ ] Class session management
- [ ] Learning material upload
- [ ] Assignment creation
- [ ] Quiz/exam interface
- [ ] Student progress tracking
- [ ] Assessment results

## Notes
- UI designs available in UI/ folder as HTML reference
- Start implementation in frontend/unsia-portal-web
- Use existing Next.js project as base
- Follow design tokens from HTML mockups
