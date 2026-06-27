# TODO - Akademik Match UI Task

## Task
Match akademik frontend features & menus with UI/AKADEMIK design

The task is to ensure: **frontend unsia-academic** has the same features and menus as shown in **UI/AKADEMIK** design

## UI/AKADEMIK/ADMIN Menus (from design):
1. Beranda
   - Dashboard
2. Master Akademik
   - Tahun Ajaran
   - Periode Akademik
   - Kurikulum
   - Mata Kuliah ⚠️ MISSING
3. Operasional
   - Kelas Kuliah ⚠️ MISSING
   - Jadwal & Sesi
   - Nilai & KHS
4. SDM & Mahasiswa
   - Data Mahasiswa
   - Dosen & Pengampu
   - Persuratan
5. Lain-Lain
   - Laporan
   - Pengaturan

## Status: 12 of 14 menus implemented (86%)

## Current Frontend Pages (unsia-academic/app):
| UI Category | UI Menu Item | Frontend Page | Status |
|---|---|---|---
| Beranda | Dashboard | admin/ | ✅ DONE |
| Master Akademik | Tahun Ajaran | tahun-ajaran/ | ✅ DONE |
| | Periode Akademik | periode/ | ✅ DONE |
| | Kurikulum | kurikulum/ | ✅ DONE |
| | Mata Kuliah | _MISSING_ | ❌ NEED CREATE |
| Operasional | Kelas Kuliah | _MISSING_ | ❌ NEED CREATE |
| | Jadwal & Sesi | schedule/ | ✅ DONE |
| | Nilai & KHS | grade/ | ✅ DONE |
| SDM & Mahasiswa | Data Mahasiswa | student/ | ✅ DONE |
| | Dosen & Pengampu | lecturer/ | ✅ DONE |
| | Persuratan | persuratan/ | ✅ DONE |
| Lain-Lain | Laporan | laporan/ | ✅ DONE |
| | Pengaturan | pengaturan/ | ✅ DONE |

## Next Steps - CREATE MISSING PAGES:
- [x] Create matakuliah/ page (Mata Kuliah - Master Akademik)
- [x] Create kelas/ page (Kelas Kuliah - Operasional)

## ✅ COMPLETED - ALL 14 MENUS MATCHED (100%)

All menus from UI/AKADEMIK/ADMIN now implemented in unsia-academic frontend.
