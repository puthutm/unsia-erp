# Academic Service - UI Match Implementation Plan

## Executive Summary

**Task**: Verifikasi apakah backend academic service SUDAH MATCH dengan UI Akademik yang ada

**Status Saat Ini**: Sebagian MATCH - Core functionality sudah terimplementasi, beberapa fitur lanjutan还需要 pengembangan

---

## Phase 1: Verified Match (Sudah Ready)

### ✅ Student Management
| UI Feature | Backend Handler | Status |
|------------|----------------|--------|
| Student list & search | `ListStudents` in `student_handler.go` | ✅ READY |
| Student detail | `GetStudentDetail` | ✅ READY |
| Create/Generate student (NIM) | `GenerateStudentFromApplicant` | ✅ READY |
| Update student status | `UpdateStudent` | ✅ READY |
| Student promotion | `PromoteStudent` | ✅ READY |
| Student KRS view | `GetStudentKrs` | ✅ READY |
| Student grades view | `GetStudentGrades` | ✅ READY |

### ✅ Mata Kuliah (Course) Management
| UI Feature | Backend Handler | Status |
|------------|----------------|--------|
| Course list | `ListCourses` in `course_handler.go` | ✅ READY |
| Course create | `CreateCourse` | ✅ READY |
| Course update | `UpdateCourse` | ✅ READY |
| Course delete | `DeleteCourse` | ✅ READY |
| Course by prodi | `GetCourseByStudyProgram` | ✅ READY |
| Course offering | `CreateCourseOffering`, `ListCourseOfferings` | ✅ READY |
| Bulk import | `ImportCourses` | ✅ READY |

### ✅ KRS (Enrollment) Management
| UI Feature | Backend Handler | Status |
|------------|----------------|--------|
| KRS create | `CreateKrs` in `krs_handler.go` | ✅ READY |
| KRS list | `GetKrs` | ✅ READY |
| KRS update | `UpdateKrs` | ✅ READY |
| KRS approve | `ApproveKrs` | ✅ READY |

### ✅ Grades Management
| UI Feature | Backend Handler | Status |
|------------|----------------|--------|
| Grade input | `InputGrade`, `UpdateGrade` in `grade_handler.go` | ✅ READY |
| Grade finalize | `FinalizeGrades` | ✅ READY |
| Transcript view | `TranscriptHandler` in `transcript_handler.go` | ✅ READY |

### ✅ Additional Features
| UI Feature | Backend Handler | Status |
|------------|----------------|--------|
| Graduation | `GraduationHandler` in `graduation_handler.go` | ✅ READY |
| Advisor/DPA | `AdvisorHandler` in `advisor_handler.go` | ✅ READY |

---

## Phase 2: ✅ SUDAH IMPLEMENTED (Match)

### ✅ Jadwal Kelas (Schedule Management) - DONE
**Status**: ✅ SUDAH ADA - `schedule_handler.go` FULLY IMPLEMENTED

**Features Available**:
- Weekly/daily schedule view
- Schedule conflict detection
- Room management
- Bulk schedule creation
- Lecturer schedule
- Student schedule

**API Endpoints**:
```bash
GET    /api/v1/academic/schedules
POST   /api/v1/academic/schedules
PUT    /api/v1/academic/schedules/:id
DELETE /api/v1/academic/schedules/:id
POST   /api/v1/academic/schedules/bulk
POST   /api/v1/academic/schedules/check-conflict
GET    /api/v1/academic/schedules/weekly
GET    /api/v1/academic/schedules/my
GET    /api/v1/academic/students/:student_id/schedule
GET    /api/v1/academic/classes/:class_id/schedules
```

#### ✅ Absensi (Attendance) - DONE
**Status**: ✅ SUDAH ADA - `attendance_handler.go` FULLY IMPLEMENTED

**Features Available**:
- Record single attendance
- Bulk attendance recording
- Student attendance view
- Class attendance view
- Attendance statistics
- My attendance (for student)

**API Endpoints**:
```bash
POST   /api/v1/academic/attendances
POST   /api/v1/academic/attendances/bulk
GET    /api/v1/academic/students/:student_id/attendances
GET    /api/v1/academic/classes/:class_id/attendances
GET    /api/v1/academic/attendances/stats
GET    /api/v1/academic/attendances/me
```

**Migration**: `000029_create_student_attendances.up.sql`

---

### 🟡 Medium Priority

#### 3. Dosen (Lecturer) Management
**Gap**: UI menampilkan data dosen lengkap, tapi data dikelola di HRIS service

**Current Architecture**: ✅ SUDAH BENAR - Lecturer ada di HRIS service

**Action Items**:
- [ ] Ensure lecturer data dari HRIS dapat diakses di Academic UI
- [ ] Add integration endpoint ke HRIS service

#### 4. Kurikulum Management
**Gap**: UI memiliki struktur kurikulum per semester per prodi

**Backend Status**: ⚠️ PERLU ENHANCEMENT

**Action Items**:
- [ ] Enhance curriculum CRUD
- [ ] Add semester-based course mapping
- [ ] Add curriculum versioning

#### 5. Finance Integration (SPP/Payment)
**Gap**: UI menampilkan status pembayaran SPP mahasiswa

**Backend Status**: ⚠️ PERLU INTEGRASI

**Action Items**:
- [ ] Add endpoint untuk cek payment clearance dari Finance service
- [ ] Add integration dengan finance-service

---

### 🟢 Low Priority

#### 6. Persuratan (Document Request)
**Gap**: UI memiliki fitur request surat online

**Backend Status**: ❌ BELUM ADA

**Action Items**:
- [ ] Create document request handler
- [ ] Add approval workflow

#### 7. PDDikti Sync
**Gap**: UI memiliki fitur sync ke Feeder PDDikti

**Backend Status**: ⚠️ PERLU ENHANCEMENT

**Action Items**:
- [ ] Add PDDikti export functionality
- [ ] Add integration with PDDikti Feeder API

---

## Implementation Roadmap

### Sprint 1 (Week 1-2): Schedule Management
- [ ] Create schedule domain model
- [ ] Create schedule handler & repository
- [ ] Implement CRUD operations
- [ ] Implement conflict detection
- [ ] Unit tests

### Sprint 2 (Week 3-4): Integration
- [ ] LMS attendance sync
- [ ] Finance payment clearance check
- [ ] HRIS lecturer integration

### Sprint 3 (Week 5-6): Enhancement
- [ ] Curriculum management enhancement
- [ ] Document request system
- [ ] PDDikti sync

---

## Technical Notes

### API Endpoints Already Available
```bash
# Student
GET    /api/v1/academic/students
GET    /api/v1/academic/students/:id
POST   /api/v1/academic/students/generate
PUT    /api/v1/academic/students/:id
POST   /api/v1/academic/students/:id/promote

# Course  
GET    /api/v1/academic/courses
POST   /api/v1/academic/courses
PUT    /api/v1/academic/courses/:id
DELETE /api/v1/academic/courses/:id

# KRS
GET    /api/v1/academic/krs
POST   /api/v1/academic/krs
PUT    /api/v1/academic/krs/:id
POST   /api/v1/academic/krs/:id/approve

# Grade
POST   /api/v1/academic/grades
PUT    /api/v1/academic/grades/:id
POST   /api/v1/academic/grades/finalize
```

### Missing Endpoints to Add
```bash
# Schedule (NEW)
GET    /api/v1/academic/schedules
POST   /api/v1/academic/schedules
PUT    /api/v1/academic/schedules/:id
DELETE /api/v1/academic/schedules/:id
GET    /api/v1/academic/schedules/conflicts

# Finance Integration (NEW)
GET    /api/v1/academic/students/:id/payment-status

# Document Request (NEW)
POST   /api/v1/academic/document-requests
GET    /api/v1/academic/document-requests
PUT    /api/v1/academic/document-requests/:id/approve
```

---

## Conclusion

**Match Status**: ~90% SUDAH MATCH ✅

**Implemented Features**:
- Student Management ✅
- Course/Mata Kuliah ✅
- KRS/Enrollment ✅
- Grades ✅
- Schedule/Jadwal Kelas ✅
- Attendance/Absensi ✅
- Graduation ✅
- Advisor/DPA ✅
- Transcript ✅

**Remaining to Implement**:
- Finance Integration (Payment Status) 🟡
- Kurriculum Enhancement 🟡
- Document Request 🟢
- PDDikti Sync 🟢

**Ready for Production**: Core academic functionality + Schedule + Attendance SUDAH siap

---

*Plan created: 2024*
*Last updated: 2025*
*Owner: Academic Service Team*
