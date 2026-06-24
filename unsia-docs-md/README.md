# UNSIA ERP Documentation Markdown Pack

Paket ini berisi dokumen Markdown terpisah dari baseline ERP UNSIA, mulai dari PRD sampai SRS, termasuk API, Event Contract, DBML, UAT, developer specification, rencana kerja, dan struktur repo.

## Struktur Folder

- [PRD Global UNSIA](01-prd/PRD_Global_UNSIA_v6_5_Event_Contract_Updated.md)
- [BRD UNSIA](02-brd/BRD_UNSIA_v1_1_Event_Contract_Updated.md)
- [FSD Per Modul UNSIA](03-fsd/FSD_Per_Modul_UNSIA_v1_0_Event_Contract_Updated.md)
- [UAT Scenario dan QA Test Plan Detail UNSIA](07-uat/UAT_Scenario_QA_Test_Plan_Detail_ERP_UNSIA_v1_0_1_Event_Contract_Updated.md)
- [Developer Implementation Specification ERP UNSIA](08-developer/Developer_Implementation_Specification_ERP_UNSIA.md)
- [Rencana Kerja Developer ERP UNSIA](09-workplan/Rencana_Kerja_Developer_ERP_UNSIA.md)
- [Struktur Repo ERP UNSIA](10-repo-structure/Struktur_Repo_ERP_UNSIA.md)
- [SRS ERP UNSIA](11-srs/SRS_ERP_UNSIA.md)
- [API Contract / OpenAPI](04-api-contract/OpenAPI_Swagger_Final_ERP_UNSIA_v1_0_1_Event_Contract_Updated.md)
- [Event Contract ERP UNSIA](05-event-contract/Event_Contract_ERP_UNSIA.md)
- [ERD/DBML Global UNSIA](06-erd-dbml/DBML_Global_UNSIA_v1_0_1_Event_Contract_Updated.md)
- [API Endpoints List](10-repo-structure/API_LIST.md)

## API Endpoints Summary

| Service | Port | Handlers | APIs | Status |
|---------|------|---------|-------|--------|
| Core Service | 8001 | 8 | ~13 | ✅ DONE |
| Academic Service | 8002 | 5 | ~15 | ✅ DONE |
| Finance Service | 8004 | 15 | ~40+ | ✅ DONE |
| PMB Service | 8005 | 7 | ~20+ | ✅ DONE |
| LMS Service | 8003 | 6 | ~12+ | ✅ DONE |
| Assessment Service | 8006 | 1 | ~12+ | ✅ DONE |
| HRIS Service | 8007 | 7 | ~15+ | ✅ DONE |
| CRM Service | 8008 | 2 | ~15+ | ✅ DONE |
| Reference Service | 8009 | 1 | ~20+ | ✅ DONE |
| Portal Service | 8010 | 1 | ~10+ | ✅ DONE |

### Services Status:

1. **Core Service (Port 8001)**: ✅ DONE - Auth, Session, Application, Service Token, Audit, Webhook, External App
2. **Academic Service (Port 8002)**: ✅ DONE - Student, Course, Grade, KRS, Academic (SUDAH MATCH dengan LMS untuk Grade sync)
3. **Finance Service (Port 8004)**: ✅ DONE - Invoice, Payment, Budget, Cash, Journal, Clearance, Disbursement, Payroll, dll
4. **PMB Service (Port 8005)**: ✅ DONE - Applicant, Document, Wave, Study Program, Selection, Public, Dashboard
5. **LMS Service (Port 8003)**: ✅ DONE - Class Sync, Enrollment Sync, Grade Sync ke Academic
6. **Assessment Service (Port 8006)**: ✅ DONE - Question Bank, Question Sets, Assessment Sessions
7. **HRIS Service (Port 8007)**: ✅ DONE - Employee, Attendance, Leave Request, BKD, Performance Review
8. **CRM Service (Port 8008)**: ✅ DONE - Campaign, Agent, Referral, Lead, Commission
9. **Reference Service (Port 8009)**: ✅ DONE - provinces, cities, religions, study programs
10. **Portal Service (Port 8010)**: ✅ DONE - Notification, User Preferences, Menu Shortcuts

📌 **Note**: Task "akademik udah match" = Academic dengan LMS sudah terintegrasi untuk sync nilai (Grade)

## Catatan Implementasi

- Dokumen Markdown ini disiapkan agar bisa langsung dimasukkan ke repo `unsia-docs`.
- OpenAPI JSON dan DBML tetap dipertahankan dalam bentuk code block agar tidak kehilangan detail teknis.
- Jika ada revisi pada PRD/BRD/FSD, regenerasi file Markdown agar traceability tetap sinkron.
- Semua API endpoints详见 [API_LIST.md](10-repo-structure/API_LIST.md)
