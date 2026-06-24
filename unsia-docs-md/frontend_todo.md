# Frontend Implementation TODO

## Priority Sequence:
1. SSO & Portal
2. Reference Data
3. PMB (Admissions)
4. Finance
5. Academic
6. LMS

## Task Breakdown

### Phase 1: SSO & Portal (Week 1-2) - IN PROGRESS
- [x] Set up Next.js 14 project structure
- [x] Configure Tailwind CSS with UNSIA design tokens
- [x] Implement shared component library (Button component with variants)
- [x] Create utility functions (cn, formatCurrency, etc.)
- [x] Build login page with UNSIA branding
- [x] Create auth context and hooks
- [x] Create main layout with sidebar navigation
- [ ] Complete npm install in container
- [ ] Set up protected routes (middleware)
- [ ] Verify full application runs

### Phase 2: Reference Integration (Week 2-3)
- [ ] API client configuration
- [ ] Reference data fetching hooks
- [ ] Dropdown/selector components for prodi, fakultas
- [ ] Reference management CRUD pages

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
