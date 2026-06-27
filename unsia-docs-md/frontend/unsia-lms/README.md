# UNSIA LMS Frontend

Learning Management System frontend untuk Universitas Sintang.

## Struktur Folder

```
unsia-lms/
├── app/
│   ├── page.tsx              # Halaman utama LMS
│   ├── layout.tsx            # Layout LMS
│   ├── globals.css          # Global styles
│   └── session/
│       └── [id]/
│           └── page.tsx    # Detail sesi/material
├── hooks/
│   └── use-lms.ts          # LMS API hooks
├── components/            # Reusable components
├── contexts/              # React contexts
└── lib/                   # Utilities
```

## Fitur

1. **Dashboard LMS**
   - Tampilan statistik kursus
   - Daftar sesi kuliyah
   - Monitoring progress

2. **Manajemen Kursus**
   - Daftar kursus yang diambil
   - Detail informasi kursus
   - Progres pembelajaran

3. **Sesi Kuliyah**
   - Jadwal sesi
   - Materi pembelajaran
   - Tugas dan quiz
   - Diskusi

4. **Materi**
   - Upload/download materi
   - Video conference
   - Dokumen PDF

5. **Tugas**
   - Pengumpulan tugas
   - Penilaian
   - Feedback

6. **Diskusi**
   - Forum diskusi
   - Q&A dengan dosen

## API Integration

Menggunakan LMS Service di `http://localhost:8081/api/v1/lms`

### Endpoints

```typescript
// Courses
GET    /api/v1/lms/courses
POST   /api/v1/lms/courses

// Sessions  
GET    /api/v1/lms/sessions
POST   /api/v1/lms/sessions
GET    /api/v1/lms/sessions/:id

// Materials
GET    /api/v1/lms/sessions/:id/materials
POST   /api/v1/lms/sessions/:id/materials

// Assignments
GET    /api/v1/lms/sessions/:id/assignments
POST   /api/v1/lms/sessions/:id/assignments

// Discussions
GET    /api/v1/lms/sessions/:id/discussions
POST   /api/v1/lms/sessions/:id/discussions
POST   /api/v1/lms/discussions/:id/replies

// Attendance
GET    /api/v1/lms/sessions/:id/attendance
POST   /api/v1/lms/sessions/:id/attendance
```

## Tech Stack

- **Framework**: Next.js 13+ (App Router)
- **Styling**: Tailwind CSS
- **State**: React Hooks
- **API**: REST

## Getting Started

```bash
# Install dependencies
npm install

# Run development
npm run dev
```

## Konfigurasi

Buat file `.env.local`:

```env
NEXT_PUBLIC_LMS_API_URL=http://localhost:8081/api/v1/lms
