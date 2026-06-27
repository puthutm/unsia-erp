# Frontend Build Fixes TODO

## Status: Analysis Complete

### ✅ SUCCESS (Fixed):
- [x] unsia-crm - Module alias paths fixed
- [x] unsia-academic - Already working
- [x] unsia-finance - Already working

### ❌ FAILED - Needs Fixing:

1. **[unsia-core]** - Tailwind v4 error
   - Error: Cannot apply unknown utility class 'border-border'
   - File: `frontend/unsia-core/app/globals.css`
   - Fix: Update to Tailwind v4 syntax in globals.css

2. **[unsia-pmb]** - TypeScript null check error  
   - Error: 'stats' is possibly 'null'
   - File: `frontend/unsia-pmb/app/commandcenter/page.tsx:33`
   - Fix: Add null check for stats

3. **[unsia-lms]** - Missing root layout
   - Error: page.tsx doesn't have a root layout
   - File: `frontend/unsia-lms/app/page.tsx`
   - Fix: Create layout.tsx or ensure proper app directory structure

4. **[unsia-hris]** - Missing contexts
   - Error: Cannot find module '@/contexts/reference-context'
   - File: `frontend/unsia-hris/hooks/index.ts`
   - Fix: Create missing context files OR update paths in tsconfig.json

5. **[unsia-assessment]** - Missing property
   - Error: Property 'currentAttempt' does not exist
   - File: `frontend/unsia-assessment/app/exam/page.tsx`
   - Fix: Add currentAttempt to useAssessment hook

6. **[unsia-reference]** - Missing contexts
   - Error: Cannot find module '@/contexts/auth-context'
   - File: `frontend/unsia-reference/hooks/index.ts`
   - Fix: Create missing context files OR update paths
