# PMB Reference UI Implementation - TODO List

## Task Summary
- Task: Modul PMB Referensi UIintegration with Reference Service APIs
- User Confirmation: All reference data types, use existing APIs only
- Priority: Fokus di Backend (BE) dulu

## Implementation Steps

### Step 1: Analyze Available APIs from Reference Service
- [x] Study Programs - GET /api/v1/ref/study-programs
- [x] Academic Years - GET /api/v1/ref/academic-years
- [x] Academic Periods - GET /api/v1/ref/academic-periods
- [x] PMB Waves - GET /api/v1/ref/pmb-waves
- [x] Admission Paths - Need to check if available
- [x] Religions - Need to check if available
- [x] Regions (Provinces, Cities, Districts, Villages) - Need to check
- [x] Document Types - GET /api/v1/ref/document-types
- [x] Payment Components - GET /api/v1/ref/payment-components
- [x] Payment Methods - GET /api/v1/ref/payment-methods

### Step 2: Create PMB Reference UI HTML File
- Location: unsia-docs-md/UI/PMB/
- Features needed:
  - Tab/section navigation for each reference type
  - Data tables displaying reference data
  - Search/filter functionality
  - Integration with Reference Service APIs (GET endpoints)

### Step 3: Add Missing API Endpoints to Reference Service (if needed)
- Admission Paths - Add list & create endpoints
- Religions - Add list & create endpoints  
- Regions (Provinces, Cities, Districts, Villages) - Add hierarchical endpoints

### Step 4: Test Integration
- Verify APIs can be called from PMB UI
- Verify data displays correctly

## Reference Data Types to Include in PMB Reference UI
1. Program Studi (Study Programs)
2. Gelombang PMB (PMB Waves)
3. Jalur Masuk (Admission Paths)
4. Agama (Religions)
5. Wilayah (Provinces, Cities, Districts, Villages)
6. Periode Akademik (Academic Periods)
7. Tahun Ajaran (Academic Years)
8. Dokumen Persyaratan (Document Types)
9. Komponen Pembayaran (Payment Components)
10. Metode Pembayaran (Payment Methods)

## API Endpoints to Use
- Reference Service Base URL: http://localhost:8002
- All GET endpoints are protected (require valid JWT token)
