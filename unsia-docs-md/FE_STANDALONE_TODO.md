# FE Standalone TODO

## Task: Make each frontend module run in its own container

## Progress Summary

| Module | Dockerfile | next.config.js | docker-compose | Status |
|--------|-----------|----------------|---------------|--------|
| unsia-portal-web | ✅ Exists | ⏳ TBD | ⏳ TBD | In Progress |
| unsia-pmb | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-academic | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-finance | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-lms | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-hris | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-assessment | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-crm | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-reference | ✅ Exists | ⏳ TBD | ⏳ TBD | Pending |
| unsia-core | ✅ Created | ⏳ TBD | ⏳ TBD | Pending |

## Action Items

### 1. Update next.config.js for each module
For each frontend module, update next.config.js to enable standalone mode:

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  // ... other config
}

module.exports = nextConfig
```

### 2. Add each FE to docker-compose.yml
Add each frontend service to docker-compose.yml with:
- Build context pointing to respective directory
- Environment variables for backend URLs
- Port mapping
- Network connectivity

### 3. Test and verify
- Verify each container builds successfully
- Verify each container runs independently
- Verify backend communication works

## Files Modified

- `unsia-docs-md/frontend/unsia-core/Dockerfile` - Created
- `unsia-docs-md/FE_STANDALONE_PLAN.md` - Created
- `unsia-docs-md/FE_STANDALONE_TODO.md` - This file

## Notes
- Run `npm run build` in each frontend directory before container builds
- Ensure all required files (package.json, tsconfig.json, etc.) exist in each module directory
- Each Next.js standalone container needs server.js generated during build
