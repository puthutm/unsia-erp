# TODO - Frontend Implementation Progress

## Current Task: Fix Finance Page & Continue Development

### Step 1: Fix API Port Configuration (COMPLETED)
- Fixed: Finance service port from 8002 → 8005 in constants.ts ✅
- Fixed: LMS service port from 8005 → 8006 in constants.ts ✅

### Step 2: Fix Finance Page Issues (IN PROGRESS)
- Issue: Finance page tries to use stats but doesn't fetch stats
- Issue: getStatusBadge has Record type error
- Need to add proper stats fetching

### Step 3: Create Protected Route Wrapper (PENDING)
- Create components/auth-guard.tsx for protected routes

### Step 4: Install Dependencies (PENDING)
- Run npm install in frontend project

---

# TODO - Implementasi Dosen PA (Pembimbing Akademik)

## Plan
1. Add AdvisorID field to Student model ✅ DONE
2. Add repository methods for student advisor operations ✅ DONE
3. Add handler methods for assigning/removing PA ✅ DONE
4. Add router endpoints (pending router file update)

## Steps
- [x] 1. Update Student model - add advisor_id field
- [x] 2. Add repository methods - GetStudentsByAdvisor, UpdateStudentAdvisor
- [x] 3. Add handler methods - AssignAdvisor, RemoveAdvisor, GetStudentsByAdvisor
- [x] 4. Add router endpoints in main.go

---

# TODO - API List to README

## Task: Add API list summary to README.md
- [x] 1. Read API_LIST.md file
- [x] 2. Add API Endpoints Summary section to README.md
- [x] 3. Add Service Status table to README.md
- [x] 4. Add note about "akademik udah match" to README.md

## Status: DONE ✅

---

# TODO - API Unplanned List

## Task: Create plan for APIs that are still "Planned" in API_LIST.md
- [x] 1. Analyze API_LIST.md to identify "Planned" endpoints
- [x] 2. Identify services still needing implementation
- [x] 3. Create API_UNPLANNED.md with implementation plan

## Status: DONE ✅
