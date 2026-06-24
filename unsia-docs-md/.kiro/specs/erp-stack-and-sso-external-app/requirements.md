# Requirements Document

## Introduction

Dokumen ini mendefinisikan kebutuhan teknis untuk dua perubahan besar pada ERP Pendidikan / SIAKAD Terintegrasi UNSIA:

1. **Migrasi Stack Teknologi** — Seluruh service backend bermigrasi dari Node.js/TypeScript ke Go (Golang) dengan Gin framework, mengganti Prisma dengan GORM + sqlc, dan mengganti arsitektur repo yang ada dengan pola clean architecture Go. Frontend tetap menggunakan Next.js namun beralih dari Pages Router ke App Router.

2. **Fitur SSO External App Registration** — Penambahan mekanisme registrasi OAuth 2.0 / OIDC untuk aplikasi eksternal (vendor, mitra, integrasi pihak ketiga) agar dapat terhubung ke ekosistem UNSIA secara aman dan terkontrol, dengan model persetujuan admin.

ERP UNSIA terdiri dari 10 modul (Core, Referensi, CRM, PMB, Finance, Akademik, HRIS, LMS, Assessment, Portal) yang masing-masing memiliki database fisik terpisah dan berkomunikasi melalui API contract, event contract, outbox/inbox, dan snapshot/read model.

---

## Glossary

- **Go_Service**: Service backend yang dibangun dengan Go + Gin framework.
- **GORM**: ORM untuk Go yang digunakan sebagai database layer utama.
- **sqlc**: Tool codegen SQL-to-Go untuk query kompleks dan type-safe.
- **golang-migrate**: Tool migrasi database berbasis Go yang menggantikan Prisma Migrate.
- **Clean_Architecture**: Pola arsitektur dengan layer `cmd/`, `internal/domain`, `internal/application`, `internal/infrastructure`, `internal/handler`, `internal/middleware`.
- **Shared_Go_Module**: Go module terpisah yang berisi kode reusable lintas service.
- **Next_App_Router**: Paradigma routing Next.js berbasis App Router (bukan Pages Router).
- **TanStack_Query**: Library data fetching dan state management untuk React/Next.js.
- **OAuth_Client**: Aplikasi eksternal yang terdaftar sebagai OAuth 2.0 client di sistem UNSIA.
- **Authorization_Code_Flow**: Alur OAuth 2.0 untuk web app eksternal berbasis redirect.
- **Client_Credentials_Flow**: Alur OAuth 2.0 untuk machine-to-machine tanpa intervensi user.
- **Dynamic_Client_Registration**: Mekanisme self-service registrasi OAuth client melalui API publik.
- **PKCE**: Proof Key for Code Exchange — ekstensi keamanan wajib untuk Authorization Code Flow.
- **OIDC**: OpenID Connect — lapisan identitas di atas OAuth 2.0.
- **OAuth_Server**: Komponen Core Service yang mengelola seluruh alur OAuth 2.0 dan OIDC.
- **Super_Admin**: Role dengan akses global penuh termasuk manajemen OAuth client.
- **Developer**: Pemilik atau operator aplikasi eksternal yang mendaftarkan OAuth client.
- **core_db**: Database fisik milik Core Service, tempat tabel OAuth disimpan.
- **RabbitMQ**: Message broker untuk komunikasi event asynchronous antar service.
- **Outbox_Worker**: Go service terpisah yang memproses outbox events dan mempublish ke RabbitMQ.
- **RS256**: Algoritma tanda tangan JWT menggunakan RSA 256-bit.
- **JWKS_Endpoint**: Endpoint publik yang menyediakan JSON Web Key Set untuk verifikasi token.
- **Audit_Log**: Catatan lengkap aksi sensitif berisi actor, role, timestamp, old value, new value.
- **Integration_Worker**: Go service terpisah yang menjalankan outbox publisher, inbox consumer, retry, DLQ, dan reconciliation.

---

## Requirements

## Bagian A — Migrasi Stack Teknologi


### Kebutuhan A-1: Backend Service Berbasis Go + Gin

**User Story:** Sebagai Technical Lead, saya ingin semua service backend dibangun menggunakan Go dan Gin framework, agar tim mendapatkan performa tinggi, type safety, dan ekosistem Go yang konsisten di seluruh ERP UNSIA.

#### Kriteria Penerimaan

1. THE Go_Service SHALL dibangun menggunakan Go versi minimum 1.22 dan Gin framework versi minimum 1.9, dan versi ini SHALL dikunci di file `go.mod` setiap service.
2. WHEN sebuah service baru dibuat, THE Go_Service SHALL menggunakan Gin sebagai HTTP router dan middleware chain utama; penggunaan router atau HTTP mux lain secara langsung di luar Gin TIDAK DIIZINKAN.
3. THE Go_Service SHALL menggunakan GORM untuk operasi database yang tidak memerlukan JOIN lintas tabel atau subquery (operasi CRUD tunggal); THE Go_Service SHALL menggunakan sqlc untuk operasi yang memerlukan JOIN lintas tabel, subquery, atau agregasi kompleks agar menghasilkan kode Go bertipe-kuat.
4. THE Go_Service SHALL menggunakan golang-migrate sebagai satu-satunya tool untuk manajemen migrasi schema database; penggunaan `AutoMigrate` GORM di environment non-development TIDAK DIIZINKAN.
5. WHEN sebuah migration dijalankan, THE Go_Service SHALL mengeksekusi file migrasi sesuai urutan nomor versi yang menaik secara deterministik; IF migration yang sama dijalankan ulang THEN golang-migrate SHALL melewatinya tanpa error (idempoten).
6. THE Go_Service SHALL menggunakan golang-jwt/jwt dengan algoritma RS256 (RSA 2048-bit minimum) untuk pembuatan dan validasi JSON Web Token; penggunaan algoritma HS256 atau `none` TIDAK DIIZINKAN.
7. THE Go_Service (Core Service) SHALL menyediakan JWKS endpoint publik di path `/.well-known/jwks.json` yang mengembalikan JSON Web Key Set berformat RFC 7517 berisi semua public key aktif.
8. WHEN public key dirotasi, THE Go_Service SHALL mempertahankan public key lama di JWKS endpoint selama minimal 5 menit (TTL cache default) sebelum dihapus, agar token yang diterbitkan dengan key lama tetap dapat divalidasi selama masa berlakunya.
9. THE Go_Service SHALL mengembalikan setiap response HTTP menggunakan format envelope standar `{"success": bool, "data": {}, "error": {code: "", message: ""}, "meta": {"trace_id": "", "timestamp": ""}}` untuk semua status code sukses maupun gagal.
10. IF sebuah request gagal karena validasi input, autentikasi, otorisasi, atau business rule, THEN THE Go_Service SHALL mengembalikan error envelope dengan `error.code` berupa string kapital yang stabil (contoh: `INVALID_INPUT`, `FORBIDDEN_SCOPE`) dan `error.message` dalam bahasa Indonesia yang dapat dipahami pengguna.
11. THE `meta.trace_id` SHALL bernilai UUID v4 yang unik per-request dan SHALL sama dengan nilai header `X-Correlation-Id` yang masuk atau baru di-generate jika header tidak ada.

---

### Kebutuhan A-2: Struktur Direktori Clean Architecture Go

**User Story:** Sebagai Backend Developer, saya ingin setiap repo service mengikuti struktur clean architecture Go yang standar, agar kode mudah dipahami, diuji, dan di-maintain oleh seluruh anggota tim.

#### Kriteria Penerimaan

1. THE Go_Service SHALL mengorganisasikan kode dalam struktur direktori: `cmd/`, `internal/domain/`, `internal/application/`, `internal/infrastructure/`, `internal/handler/`, `internal/middleware/`.
2. THE Go_Service SHALL menempatkan entitas domain, value object, enum, dan state machine di dalam `internal/domain/`.
3. THE Go_Service SHALL menempatkan use case, command handler, dan query handler di dalam `internal/application/`.
4. THE Go_Service SHALL menempatkan implementasi database repository, external HTTP client, event publisher, dan storage adapter di dalam `internal/infrastructure/`.
5. THE Go_Service SHALL menempatkan controller HTTP, request validator, dan presenter response di dalam `internal/handler/`.
6. THE Go_Service SHALL menempatkan seluruh middleware Gin (auth, RBAC, correlation ID, idempotency, rate limit) di dalam `internal/middleware/`.
7. THE Go_Service SHALL menempatkan binary entry point di dalam `cmd/{service-name}/main.go`.
8. WHEN sebuah Go_Service memanggil database, THE Go_Service SHALL hanya mengakses database miliknya sendiri sesuai prinsip database per modul.
9. THE Go_Service SHALL menyertakan direktori `migrations/` berisi file SQL terurut untuk golang-migrate.
10. THE Go_Service SHALL menyertakan direktori `tests/` dengan sub-direktori `unit/`, `integration/`, dan `contract/`.

---

### Kebutuhan A-3: Shared Go Module

**User Story:** Sebagai Backend Developer, saya ingin semua service menggunakan shared library Go yang telah distandardisasi, agar implementasi auth, RBAC, audit, idempotency, event, HTTP client, error, dan observability konsisten di seluruh modul.

#### Kriteria Penerimaan

1. THE Shared_Go_Module SHALL tersedia sebagai Go module terpisah dengan nama: `shared-auth`, `shared-rbac`, `shared-audit`, `shared-idempotency`, `shared-event`, `shared-httpclient`, `shared-errorenvelope`, dan `shared-observability`.
2. THE `shared-auth` SHALL menyediakan fungsi validasi JWT, caching JWKS, dan validasi service token antar modul.
3. THE `shared-rbac` SHALL menyediakan fungsi pengecekan permission dan resolver data scope berdasarkan active role.
4. THE `shared-audit` SHALL menyediakan fungsi pencatatan audit log dengan field: actor, role, action, old_value, new_value, reason, correlation_id, timestamp.
5. THE `shared-idempotency` SHALL menyediakan fungsi simpan request hash, cache response, locking, dan expiry untuk mencegah duplikasi command.
6. THE `shared-event` SHALL menyediakan outbox writer, inbox consumer, fungsi retry, DLQ handler, dan event envelope builder.
7. THE `shared-httpclient` SHALL menyediakan HTTP client service-to-service dengan konfigurasi timeout, retry, dan circuit breaker.
8. THE `shared-errorenvelope` SHALL menyediakan format error code dan error response konsisten yang digunakan seluruh service.
9. THE `shared-observability` SHALL menyediakan structured logging, trace ID propagation, correlation ID middleware, dan metrics exporter.
10. WHEN sebuah Go_Service diinisialisasi, THE Go_Service SHALL mengimpor shared modules yang relevan alih-alih mengimplementasikan ulang fungsi yang sama.


---

### Kebutuhan A-4: Integration Worker sebagai Go Service Terpisah

**User Story:** Sebagai Backend Developer, saya ingin outbox/inbox event processing dijalankan oleh Go service terpisah menggunakan RabbitMQ, agar proses event tidak mengganggu service utama dan dapat di-scale secara independen.

#### Kriteria Penerimaan

1. THE Integration_Worker SHALL dibangun sebagai Go service terpisah (repo `unsia-integration-worker`) menggunakan RabbitMQ sebagai event broker.
2. THE Integration_Worker SHALL mengimplementasikan outbox pattern: membaca `outbox_events` berstatus PENDING dari setiap database modul dan mempublish payload ke RabbitMQ.
3. THE Integration_Worker SHALL mengimplementasikan inbox pattern: menerima event dari RabbitMQ, memeriksa duplikat berdasarkan `event_key` di `inbox_events`, dan memproses event jika belum ada.
4. WHEN sebuah event gagal diproses, THE Integration_Worker SHALL melakukan retry dengan strategi exponential backoff maksimal 5 kali sebelum memindahkan event ke Dead Letter Queue (DLQ).
5. WHEN sebuah event masuk DLQ, THE Integration_Worker SHALL mencatat `event_key`, `last_error`, `retry_count`, dan `failed_at` untuk keperluan monitoring dan replay manual.
6. THE Integration_Worker SHALL mendukung replay DLQ dengan aksi manual yang membutuhkan `reason` dan mencatat actor di audit log.
7. WHEN sebuah event duplikat diterima (event_key sudah ada di inbox_events), THE Integration_Worker SHALL menandai event sebagai IGNORED_DUPLICATE dan tidak memproses ulang.
8. THE Integration_Worker SHALL mengekspos health check endpoint untuk monitoring status koneksi RabbitMQ dan lag queue.

---

### Kebutuhan A-5: Frontend Next.js App Router

**User Story:** Sebagai Frontend Developer, saya ingin portal UNSIA menggunakan Next.js App Router dan TanStack Query, agar arsitektur frontend modern, mendukung server components, dan data fetching lebih efisien.

#### Kriteria Penerimaan

1. THE Next_App_Router SHALL digunakan sebagai paradigma routing utama di seluruh aplikasi frontend UNSIA, menggantikan Pages Router.
2. THE Next_App_Router SHALL mengorganisasikan route dalam direktori `app/` menggunakan file `page.tsx`, `layout.tsx`, dan `loading.tsx` sesuai konvensi Next.js App Router.
3. THE Next_App_Router SHALL menggunakan TanStack Query sebagai library utama untuk data fetching, caching, dan state sinkronisasi ke server.
4. WHEN data berasal dari read model atau snapshot, THE Next_App_Router SHALL menampilkan label `synced_at` atau `refreshed_at` yang menunjukkan waktu sinkronisasi terakhir kepada user.
5. THE Next_App_Router SHALL menerapkan pola route group `(auth)/` untuk halaman login dan `(portal)/` untuk halaman setelah autentikasi.
6. WHEN service backend tidak tersedia, THE Next_App_Router SHALL menampilkan degraded state yang informatif dengan pesan jelas dan trace_id untuk debugging.
7. THE Next_App_Router SHALL menegakkan akses halaman berdasarkan permission dan active role di sisi client, sebagai pelengkap validasi backend.
8. THE Next_App_Router SHALL menampilkan mandatory field marker, pesan validasi inline, dan dialog konfirmasi untuk aksi sensitif pada semua form.

---

## Bagian B — SSO External App Registration


### Kebutuhan B-1: Registrasi OAuth Client Self-Service

**User Story:** Sebagai Developer aplikasi eksternal, saya ingin mendaftarkan aplikasi saya sebagai OAuth client UNSIA melalui API publik, agar aplikasi saya dapat mengintegrasikan login via SSO UNSIA tanpa perlu intervensi manual dari tim IT.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL menyediakan endpoint publik `POST /api/v1/oauth/register` yang dapat diakses tanpa autentikasi (tidak memerlukan Bearer token) untuk registrasi client baru.
2. WHEN Developer mengirim permintaan registrasi yang lolos validasi, THE OAuth_Server SHALL menyimpan data client di tabel `oauth_clients` dengan status awal `PENDING` dan membuat entri baru di `client_registration_requests` sebagai audit trail dalam satu transaksi database yang sama.
3. IF `redirect_uris` yang didaftarkan mengandung karakter wildcard (`*`, `?`), skema selain `https` (kecuali `http://localhost` untuk development), atau format URI tidak valid, THEN THE OAuth_Server SHALL mengembalikan error `INVALID_REDIRECT_URI` dan menolak permintaan.
4. IF `grant_types` yang diminta mengandung nilai selain `authorization_code` atau `client_credentials`, THEN THE OAuth_Server SHALL mengembalikan error `UNSUPPORTED_GRANT_TYPE` dan menolak permintaan.
5. IF `allowed_scopes` yang diminta mengandung scope yang tidak terdaftar dalam daftar scope resmi sistem UNSIA (misalnya `academic:read`, `finance:invoices:read`, `pmb:applicants:read`, `profile`, `email`, `openid`), THEN THE OAuth_Server SHALL mengembalikan error `INVALID_SCOPE` dan menolak permintaan.
6. WHEN validasi seluruh field berhasil, THE OAuth_Server SHALL mengembalikan HTTP 202 dengan response body berisi `registration_id` (UUID), `status: "PENDING"`, dan pesan bahwa permintaan sedang menunggu persetujuan admin.
7. IF Developer mengirim permintaan registrasi dengan `owner_email` yang sama dan sudah ada record dengan status `PENDING` di `oauth_clients`, THEN THE OAuth_Server SHALL mengembalikan error `DUPLICATE_REGISTRATION_REQUEST` dengan HTTP 409 dan menolak pembuatan data duplikat.
8. THE OAuth_Server SHALL menyimpan data kontak pemilik — `owner_name` (wajib, string 1–255 karakter), `owner_email` (wajib, format email valid), `owner_organization` (wajib, string 1–255 karakter) — sebagai bagian dari data `oauth_clients`.
9. THE request body untuk `POST /api/v1/oauth/register` SHALL mewajibkan field: `client_name` (string, 1–100 karakter), `owner_name`, `owner_email`, `owner_organization`, `redirect_uris` (array, minimal 1 item jika `grant_types` mengandung `authorization_code`), `grant_types` (array, minimal 1 item), dan `allowed_scopes` (array, minimal 1 item); IF field wajib tidak ada THEN SHALL mengembalikan error `MISSING_REQUIRED_FIELD`.

---

### Kebutuhan B-2: Persetujuan Admin untuk OAuth Client

**User Story:** Sebagai Super Admin, saya ingin menyetujui atau menolak permintaan registrasi OAuth client dari developer/vendor, agar hanya aplikasi yang terverifikasi yang dapat menggunakan SSO UNSIA.

#### Kriteria Penerimaan

1. WHEN Super_Admin mengakses endpoint `PATCH /api/v1/admin/oauth-clients/{id}/approve`, THE OAuth_Server SHALL mengubah status `oauth_clients` dari PENDING menjadi ACTIVE.
2. WHEN status berubah menjadi ACTIVE, THE OAuth_Server SHALL menghasilkan `client_id` (UUID) dan `client_secret` (random string 32 byte) untuk client tersebut.
3. THE OAuth_Server SHALL menampilkan `client_secret` satu kali dalam response approve, kemudian menyimpannya dalam bentuk hash bcrypt sehingga tidak dapat dibaca ulang.
4. WHEN Super_Admin mengakses endpoint `PATCH /api/v1/admin/oauth-clients/{id}/suspend`, THE OAuth_Server SHALL mengubah status menjadi SUSPENDED dan memblokir seluruh alur OAuth untuk client tersebut.
5. WHEN Super_Admin mengakses endpoint `DELETE /api/v1/admin/oauth-clients/{id}/revoke`, THE OAuth_Server SHALL mengubah status menjadi REVOKED, membatalkan semua token aktif client tersebut, dan menutup akses permanen.
6. THE OAuth_Server SHALL membutuhkan field `reason` pada setiap aksi approve, suspend, dan revoke.
7. THE OAuth_Server SHALL mencatat setiap aksi admin (approve, suspend, revoke) ke audit_log dengan field: actor, role, action, client_id, reason, timestamp.
8. IF Super_Admin mencoba approve client dengan status bukan PENDING, THEN THE OAuth_Server SHALL mengembalikan error dengan kode `INVALID_CLIENT_STATUS_TRANSITION`.
9. THE OAuth_Server SHALL membutuhkan permission `oauth:clients:approve` untuk aksi approve dan `oauth:clients:manage` untuk aksi suspend dan revoke.

---

### Kebutuhan B-3: Authorization Code Flow dengan PKCE

**User Story:** Sebagai Developer web app eksternal, saya ingin menggunakan Authorization Code Flow agar pengguna UNSIA dapat login ke aplikasi saya tanpa aplikasi saya perlu menyimpan credential pengguna.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL menyediakan authorization endpoint `GET /api/v1/oauth/authorize` yang menerima parameter query: `response_type=code` (wajib), `client_id` (wajib), `redirect_uri` (wajib), `scope` (wajib), `state` (wajib), `code_challenge` (wajib), dan `code_challenge_method=S256` (wajib).
2. THE OAuth_Server SHALL mewajibkan parameter `code_challenge` dan `code_challenge_method`; `code_challenge` SHALL berformat Base64URL-encoded tanpa padding dengan panjang 43–128 karakter yang merupakan hash SHA-256 dari `code_verifier`.
3. IF parameter `code_challenge` atau `code_challenge_method` tidak ada dalam request, THEN THE OAuth_Server SHALL mengembalikan error `PKCE_REQUIRED` dengan HTTP 400 dan tidak melanjutkan proses otorisasi.
4. IF user belum terautentikasi ketika mengakses authorization endpoint, THEN THE OAuth_Server SHALL menyimpan seluruh parameter request otorisasi ke session sementara, kemudian mengarahkan user ke halaman login UNSIA; setelah login berhasil, THE OAuth_Server SHALL melanjutkan proses otorisasi dengan parameter yang disimpan.
5. IF user telah terautentikasi ketika mengakses authorization endpoint, THEN THE OAuth_Server SHALL menampilkan halaman consent yang mencantumkan nama aplikasi client dan deskripsi setiap scope yang diminta sebelum melanjutkan.
6. WHEN user memberikan consent, THE OAuth_Server SHALL menghasilkan `authorization_code` berupa string acak kriptografis (minimum 32 byte), menyimpannya di `oauth_authorization_codes` dengan `expires_at` = waktu sekarang + 10 menit, kemudian mengarahkan ke `redirect_uri` yang terdaftar dengan parameter `code` dan `state` (nilai `state` harus sama dengan `state` pada request awal).
7. THE OAuth_Server SHALL memvalidasi `redirect_uri` dalam request sebagai case-sensitive exact string match dengan salah satu `redirect_uri` yang terdaftar di `oauth_clients`; IF tidak match THEN SHALL mengembalikan error `REDIRECT_URI_MISMATCH` dengan HTTP 400 dan tidak melakukan redirect.
8. IF `client_id` tidak ditemukan di `oauth_clients` atau status client bukan `ACTIVE`, THEN THE OAuth_Server SHALL mengembalikan error `INVALID_CLIENT` dengan HTTP 400 dan tidak melanjutkan proses otorisasi.
9. IF parameter `state` tidak disertakan dalam request, THEN THE OAuth_Server SHALL mengembalikan error `STATE_REQUIRED` dengan HTTP 400 untuk mencegah CSRF.
10. WHEN Developer menukar `authorization_code` ke token via `POST /api/v1/oauth/token`, THE OAuth_Server SHALL menghitung SHA-256 dari `code_verifier` yang dikirim dan membandingkannya dengan `code_challenge` yang tersimpan; IF tidak cocok THEN SHALL mengembalikan error `INVALID_PKCE_VERIFIER` dengan HTTP 400.
11. IF `code_verifier` gagal validasi PKCE, THEN THE OAuth_Server SHALL mengembalikan error `INVALID_PKCE_VERIFIER` dengan HTTP 400 dan tidak menerbitkan token.
12. WHEN `authorization_code` sudah digunakan satu kali, THE OAuth_Server SHALL menandai `used_at` pada record dan menolak penggunaan ulang dengan error `AUTHORIZATION_CODE_ALREADY_USED`; deteksi penggunaan ulang SHALL memicu revocation semua token yang sebelumnya diterbitkan dari code tersebut dalam waktu maksimal 5 detik.
13. IF `authorization_code` sudah expired (waktu sekarang > `expires_at` yang diterbitkan), THEN THE OAuth_Server SHALL menolak penukaran dengan error `AUTHORIZATION_CODE_EXPIRED` dengan HTTP 400.
14. WHEN token berhasil diterbitkan, THE OAuth_Server SHALL menyimpan record access token di `oauth_access_tokens`, refresh token di `oauth_refresh_tokens`, dan menandai `used_at` pada `oauth_authorization_codes` untuk audit trail.


---

### Kebutuhan B-4: Client Credentials Flow

**User Story:** Sebagai Developer sistem backend eksternal, saya ingin menggunakan Client Credentials Flow agar service saya dapat memanggil API UNSIA secara machine-to-machine tanpa keterlibatan pengguna akhir.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL mendukung Client Credentials Flow pada endpoint `POST /api/v1/oauth/token` dengan parameter `grant_type=client_credentials`, `client_id`, `client_secret`, dan `scope`.
2. THE OAuth_Server SHALL memvalidasi bahwa `grant_type=client_credentials` hanya dapat digunakan oleh OAuth client yang mendaftarkan `client_credentials` dalam `grant_types`.
3. THE OAuth_Server SHALL memverifikasi `client_secret` dengan membandingkan nilai hash bcrypt dari secret yang dikirim dengan hash yang tersimpan di `oauth_clients`.
4. WHEN autentikasi client berhasil, THE OAuth_Server SHALL menerbitkan access token dengan scope yang diminta (yang merupakan subset dari `allowed_scopes` client).
5. WHEN scope yang diminta melebihi `allowed_scopes` yang terdaftar, THE OAuth_Server SHALL menolak permintaan dengan error `INVALID_SCOPE`.
6. THE OAuth_Server SHALL tidak menerbitkan refresh token untuk Client Credentials Flow; client harus melakukan autentikasi ulang ketika access token expired.
7. WHILE client berstatus SUSPENDED atau REVOKED, THE OAuth_Server SHALL menolak semua permintaan Client Credentials Flow dengan error `CLIENT_NOT_ACTIVE`.

---

### Kebutuhan B-5: Token Endpoint dan Manajemen Token

**User Story:** Sebagai Developer, saya ingin token yang diterbitkan memiliki masa berlaku yang jelas dan dapat di-refresh, agar sesi aplikasi saya tetap aktif tanpa meminta pengguna login ulang terlalu sering.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL menerbitkan access token dengan masa berlaku default 1 jam (3600 detik) yang dapat dikonfigurasi melalui environment variable `OAUTH_ACCESS_TOKEN_TTL`.
2. THE OAuth_Server SHALL menerbitkan refresh token dengan masa berlaku default 30 hari (2592000 detik) yang dapat dikonfigurasi melalui environment variable `OAUTH_REFRESH_TOKEN_TTL`.
3. WHEN Developer mengirim request ke `POST /api/v1/oauth/token` dengan `grant_type=refresh_token`, THE OAuth_Server SHALL memvalidasi refresh token, menerbitkan access token baru, dan merotasi refresh token (menerbitkan refresh token baru dan membatalkan refresh token lama).
4. IF refresh token sudah expired atau sudah digunakan ulang (rotasi), THEN THE OAuth_Server SHALL menolak dengan error `INVALID_REFRESH_TOKEN` dan mengharuskan autentikasi ulang.
5. THE OAuth_Server SHALL menyimpan semua access token yang diterbitkan di tabel `oauth_access_tokens` dengan field: `jti`, `client_id`, `user_id` (nullable untuk client credentials), `scope`, `expires_at`, `revoked_at`.
6. THE OAuth_Server SHALL menyimpan semua refresh token di tabel `oauth_refresh_tokens` dengan field: `token_hash`, `client_id`, `user_id`, `access_token_jti`, `expires_at`, `used_at`, `revoked_at`.
7. THE OAuth_Server SHALL menerapkan rate limiting pada endpoint `POST /api/v1/oauth/token` dengan batas maksimal 60 permintaan per menit per `client_id`.
8. IF rate limit terlampaui, THEN THE OAuth_Server SHALL mengembalikan HTTP 429 dengan header `Retry-After` yang menunjukkan waktu tunggu dalam detik.

---

### Kebutuhan B-6: Token Introspection dan Revocation

**User Story:** Sebagai modul internal UNSIA, saya ingin dapat memvalidasi token yang diterbitkan untuk aplikasi eksternal, agar saya bisa mengontrol akses resource yang dilindungi secara real-time.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL menyediakan endpoint `POST /api/v1/oauth/introspect` yang menerima parameter `token` dan mengembalikan metadata token jika valid.
2. WHEN token valid dan belum expired, THE OAuth_Server SHALL mengembalikan response `{"active": true, "client_id": "...", "scope": "...", "exp": ..., "sub": "..."}`.
3. WHEN token sudah expired, sudah direvoke, atau tidak dikenali, THE OAuth_Server SHALL mengembalikan response `{"active": false}`.
4. THE OAuth_Server SHALL mengharuskan caller endpoint introspect menggunakan service token internal UNSIA yang valid (bukan token publik) untuk mencegah penyalahgunaan.
5. THE OAuth_Server SHALL menyediakan endpoint `POST /api/v1/oauth/revoke` yang memungkinkan client merevoke access token atau refresh token miliknya sendiri.
6. WHEN token direvoke, THE OAuth_Server SHALL menandai `revoked_at` pada record token dan menolak semua penggunaan selanjutnya via introspection.
7. WHEN Super_Admin merevoke sebuah OAuth client, THE OAuth_Server SHALL secara otomatis merevoke semua access token dan refresh token yang aktif milik client tersebut dalam satu transaksi.


---

### Kebutuhan B-7: OIDC Discovery dan JWKS Endpoint

**User Story:** Sebagai Developer aplikasi eksternal, saya ingin menggunakan OIDC Discovery untuk mengkonfigurasi library OAuth secara otomatis, agar integrasi lebih cepat dan tidak bergantung pada konfigurasi manual.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL menyediakan endpoint `GET /api/v1/.well-known/openid-configuration` yang mengembalikan dokumen OIDC Discovery sesuai spesifikasi OpenID Connect Discovery 1.0.
2. THE OAuth_Server SHALL menyertakan field minimal dalam dokumen OIDC Discovery: `issuer`, `authorization_endpoint`, `token_endpoint`, `introspection_endpoint`, `revocation_endpoint`, `jwks_uri`, `scopes_supported`, `response_types_supported`, `grant_types_supported`, `token_endpoint_auth_methods_supported`.
3. THE OAuth_Server SHALL menyediakan endpoint `GET /api/v1/.well-known/jwks.json` yang mengembalikan JSON Web Key Set berisi public key aktif dalam format JWK.
4. WHEN public key dirotasi, THE OAuth_Server SHALL mempertahankan public key lama di JWKS endpoint selama minimum 10 menit untuk memberikan waktu transisi bagi consumer yang melakukan caching.
5. THE OAuth_Server SHALL menerbitkan access token sebagai signed JWT dengan claim: `iss` (issuer UNSIA), `sub` (user ID atau client ID), `aud` (client ID), `exp`, `iat`, `jti`, `scope`, `client_id`.
6. THE OAuth_Server SHALL mendukung penambahan claim custom OIDC seperti `email`, `name`, dan `roles` pada token via konfigurasi scope `profile` dan `email`.

---

### Kebutuhan B-8: Entitas Data OAuth di core_db

**User Story:** Sebagai DBA, saya ingin entitas OAuth tersimpan di core_db dengan skema yang jelas dan mendukung audit, agar data client dan token dapat ditelusuri dan dikelola dengan aman.

#### Kriteria Penerimaan

1. THE core_db SHALL menyimpan data OAuth client di tabel `oauth_clients` dengan field minimal: `id` (UUID), `client_id` (UUID unik), `client_secret_hash` (bcrypt), `client_name`, `redirect_uris` (array), `allowed_scopes` (array), `grant_types` (array), `status` (ENUM: PENDING/ACTIVE/SUSPENDED/REVOKED), `owner_name`, `owner_email`, `owner_organization`, `created_at`, `updated_at`, `approved_at`, `approved_by`, `suspended_at`, `revoked_at`.
2. THE core_db SHALL menyimpan authorization code di tabel `oauth_authorization_codes` dengan field: `id`, `code_hash`, `client_id`, `user_id`, `redirect_uri`, `scope`, `code_challenge`, `code_challenge_method`, `expires_at`, `used_at`, `created_at`.
3. THE core_db SHALL menyimpan access token di tabel `oauth_access_tokens` dengan field: `id`, `jti` (UUID unik), `client_id`, `user_id` (nullable), `scope`, `expires_at`, `revoked_at`, `created_at`.
4. THE core_db SHALL menyimpan refresh token di tabel `oauth_refresh_tokens` dengan field: `id`, `token_hash`, `client_id`, `user_id` (nullable), `access_token_jti`, `expires_at`, `used_at`, `revoked_at`, `created_at`.
5. THE core_db SHALL menyimpan audit trail registrasi di tabel `client_registration_requests` dengan field: `id`, `client_id` (nullable, diisi setelah approved), `owner_name`, `owner_email`, `owner_organization`, `requested_scopes`, `requested_grant_types`, `requested_redirect_uris`, `status`, `admin_notes`, `created_at`, `reviewed_at`, `reviewed_by`.
6. THE core_db SHALL mendefinisikan index pada `oauth_authorization_codes.code_hash`, `oauth_access_tokens.jti`, dan `oauth_refresh_tokens.token_hash` untuk performa lookup token.
7. THE core_db SHALL menerapkan constraint NOT NULL pada semua field wajib dan constraint CHECK pada field `status` untuk memastikan nilai enum yang valid.

---

### Kebutuhan B-9: Keamanan OAuth

**User Story:** Sebagai Security Officer, saya ingin seluruh implementasi OAuth mengikuti praktik keamanan terbaik, agar ekosistem UNSIA terlindungi dari penyalahgunaan token dan serangan pada alur OAuth.

#### Kriteria Penerimaan

1. THE OAuth_Server SHALL menyimpan `client_secret` hanya dalam bentuk hash bcrypt dengan cost factor minimum 12; plaintext secret TIDAK BOLEH disimpan di database, log, atau response setelah approval.
2. WHEN Super_Admin menyetujui registrasi OAuth client, THE OAuth_Server SHALL menampilkan `client_secret` dalam plaintext satu kali dalam response approve; setelah response terkirim, THE OAuth_Server SHALL NOT menyediakan endpoint atau mekanisme apapun untuk melihat ulang plaintext secret tersebut.
3. THE OAuth_Server SHALL memvalidasi `redirect_uri` sebagai exact string match yang case-sensitive; partial match, prefix match, suffix match, dan wildcard TIDAK DIIZINKAN.
4. THE OAuth_Server SHALL mewajibkan PKCE dengan `code_challenge_method=S256` untuk semua Authorization Code Flow.
5. IF permintaan Authorization Code Flow tidak menyertakan `code_challenge` atau menggunakan `code_challenge_method=plain`, THEN THE OAuth_Server SHALL mengembalikan error `PKCE_REQUIRED` dengan HTTP 400 dan tidak melanjutkan proses otorisasi.
6. THE OAuth_Server SHALL menolak permintaan Authorization Code Flow yang tidak menyertakan parameter `state` dengan error `STATE_REQUIRED` untuk mencegah serangan CSRF.
7. WHEN `authorization_code` dideteksi digunakan lebih dari satu kali, THE OAuth_Server SHALL mengembalikan error `CODE_ALREADY_USED` dan merevoke semua access token serta refresh token yang pernah diterbitkan dari authorization code tersebut dalam waktu maksimal 5 detik.
8. IF revocation semua token terkait gagal dalam 5 detik, THEN THE OAuth_Server SHALL mencatat kegagalan tersebut ke audit log dengan severity HIGH dan memblokir client sementara sambil menunggu proses revocation selesai.
9. IF jumlah permintaan ke endpoint `POST /api/v1/oauth/token` dari satu `client_id` melebihi 60 dalam jendela waktu 60 detik, THEN THE OAuth_Server SHALL mengembalikan HTTP 429 dengan header `Retry-After` dan error `RATE_LIMIT_EXCEEDED`.
10. IF jumlah permintaan ke endpoint `POST /api/v1/oauth/register` dari satu alamat IP melebihi 10 dalam jendela waktu 3600 detik, THEN THE OAuth_Server SHALL mengembalikan HTTP 429 dengan header `Retry-After` dan error `RATE_LIMIT_EXCEEDED`.
11. THE OAuth_Server SHALL mencatat semua aksi admin (approve, suspend, revoke, view credentials) ke `audit_logs` di `core_db` dengan field: actor, active_role, action, resource_id, reason, timestamp; IF pencatatan audit log gagal THEN aksi admin SHALL dibatalkan dan mengembalikan error `AUDIT_LOG_FAILURE`.
12. IF sebuah OAuth client berstatus `SUSPENDED` atau `REVOKED`, THEN THE OAuth_Server SHALL menolak semua permintaan token untuk client tersebut dengan error `CLIENT_NOT_ACTIVE` dan HTTP 401.


---

### Kebutuhan B-10: UI Admin Manajemen OAuth Client

**User Story:** Sebagai Super Admin, saya ingin memiliki halaman admin untuk mengelola OAuth client dari aplikasi eksternal, agar saya dapat melihat, menyetujui, memantau, dan menangguhkan client dengan mudah.

#### Kriteria Penerimaan

1. THE Next_App_Router SHALL menyediakan halaman list OAuth client di path `/admin/oauth-clients` yang menampilkan tabel berisi semua client dengan kolom: nama client, organisasi, status, grant types, tanggal registrasi, dan aksi.
2. THE Next_App_Router SHALL menyediakan filter pada halaman list berdasarkan status (ALL, PENDING, ACTIVE, SUSPENDED, REVOKED) dan fitur pencarian berdasarkan nama client atau email owner.
3. THE Next_App_Router SHALL menampilkan badge berwarna berbeda untuk setiap status: PENDING (kuning), ACTIVE (hijau), SUSPENDED (oranye), REVOKED (merah).
4. THE Next_App_Router SHALL menyediakan halaman detail OAuth client di path `/admin/oauth-clients/{id}` yang menampilkan: informasi lengkap client, allowed scopes, redirect URIs, grant types, histori status, dan log aksi admin.
5. WHEN Super_Admin melakukan aksi approve, suspend, atau revoke pada halaman detail, THE Next_App_Router SHALL menampilkan dialog konfirmasi yang meminta field `reason` sebelum mengirimkan request ke API.
6. THE Next_App_Router SHALL hanya menampilkan tombol aksi yang relevan berdasarkan status client saat ini: tombol "Approve" dan "Reject" hanya tampil untuk client PENDING; tombol "Suspend" hanya tampil untuk client ACTIVE; tombol "Revoke" tampil untuk client ACTIVE dan SUSPENDED.
7. WHEN aksi admin berhasil, THE Next_App_Router SHALL menampilkan notifikasi sukses dan memperbarui tampilan status client tanpa full page reload.
8. THE Next_App_Router SHALL mengharuskan user memiliki permission `oauth:clients:view` untuk mengakses halaman list dan detail.

---

### Kebutuhan B-11: UI Developer — Melihat Credentials

**User Story:** Sebagai Developer, saya ingin melihat `client_id` dan mendapatkan notifikasi ketika aplikasi saya disetujui, agar saya bisa segera mengintegrasikan aplikasi saya dengan SSO UNSIA.

#### Kriteria Penerimaan

1. THE Next_App_Router SHALL menyediakan halaman di path `/developer/oauth-credentials` yang dapat diakses oleh pengguna yang memiliki peran Developer setelah aplikasinya berstatus ACTIVE.
2. THE Next_App_Router SHALL menampilkan `client_id` pada halaman credentials setelah client disetujui.
3. THE Next_App_Router SHALL menampilkan `client_secret` hanya sekali dari response API pada momen pertama approval, disertai pesan peringatan bahwa secret tidak dapat ditampilkan ulang.
4. THE Next_App_Router SHALL menyediakan tombol "Copy to Clipboard" untuk `client_id` dan `client_secret` untuk memudahkan developer.
5. WHEN Developer mengakses halaman credentials setelah sesi pertama (secret sudah tidak tersedia dari API), THE Next_App_Router SHALL menampilkan pesan "Client Secret hanya ditampilkan sekali saat approval. Jika Anda kehilangan secret, hubungi administrator untuk reset."
6. THE Next_App_Router SHALL menampilkan informasi tambahan pada halaman credentials: status client, allowed scopes, redirect URIs, grant types, dan tautan ke dokumentasi OAuth UNSIA.
7. WHILE status client adalah PENDING, THE Next_App_Router SHALL menampilkan halaman status pending yang menginformasikan bahwa permintaan sedang dalam proses review.

---

### Kebutuhan B-12: Integrasi OAuth dengan Modul Internal UNSIA

**User Story:** Sebagai modul internal UNSIA, saya ingin dapat memverifikasi token dari aplikasi eksternal menggunakan introspection, agar saya bisa mengontrol akses ke resource yang saya jaga.

#### Kriteria Penerimaan

1. THE Go_Service SHALL memvalidasi token eksternal melalui endpoint introspect `POST /api/v1/oauth/introspect` menggunakan service token internal UNSIA.
2. WHEN token eksternal terbukti valid melalui introspection, THE Go_Service SHALL menerapkan scope-based access control sesuai scope yang terkandung dalam token.
3. THE OAuth_Server SHALL mendefinisikan scope hierarki yang jelas untuk sumber daya UNSIA, contoh: `academic:read`, `finance:invoices:read`, `pmb:applicants:read`.
4. THE Go_Service SHALL menolak permintaan dari token eksternal yang memiliki scope tidak mencukupi dengan HTTP 403 dan error code `INSUFFICIENT_SCOPE`.
5. THE OAuth_Server SHALL menyediakan cache TTL 60 detik untuk hasil introspection guna mengurangi beban Core Service pada volume tinggi.
6. IF Core Service tidak tersedia untuk introspection, THEN THE Go_Service SHALL jatuh ke mode validasi JWT lokal menggunakan cached JWKS public key, dengan asumsi token belum direvoke.


---

## Ringkasan Endpoint OAuth

| Method | Path | Akses | Fungsi |
|--------|------|-------|--------|
| POST | `/api/v1/oauth/register` | Publik (tanpa auth) | Self-service registrasi OAuth client |
| GET | `/api/v1/oauth/authorize` | User UNSIA terautentikasi | Authorization endpoint (redirect-based) |
| POST | `/api/v1/oauth/token` | OAuth Client | Exchange code/credentials untuk token |
| POST | `/api/v1/oauth/introspect` | Service Token Internal | Validasi token aktif |
| POST | `/api/v1/oauth/revoke` | OAuth Client | Revoke token |
| GET | `/api/v1/.well-known/openid-configuration` | Publik | OIDC Discovery metadata |
| GET | `/api/v1/.well-known/jwks.json` | Publik | Public key untuk verifikasi JWT |
| GET | `/api/v1/admin/oauth-clients` | Super_Admin | List semua registered clients |
| PATCH | `/api/v1/admin/oauth-clients/{id}/approve` | Super_Admin | Approve pending client |
| PATCH | `/api/v1/admin/oauth-clients/{id}/suspend` | Super_Admin | Suspend active client |
| DELETE | `/api/v1/admin/oauth-clients/{id}/revoke` | Super_Admin | Revoke client secara permanen |

---

## Ringkasan Entitas Database (core_db)

| Tabel | Fungsi |
|-------|--------|
| `oauth_clients` | Data client terdaftar, status, secret hash, konfigurasi |
| `oauth_authorization_codes` | Kode sementara untuk Authorization Code Flow |
| `oauth_access_tokens` | Token akses yang diterbitkan |
| `oauth_refresh_tokens` | Refresh token untuk perpanjangan sesi |
| `client_registration_requests` | Audit trail semua permintaan registrasi |

---

## Matriks Testabilitas Kriteria Penerimaan

| Kebutuhan | Tipe Pengujian | Catatan |
|-----------|----------------|---------|
| A-1: Go + Gin Service | Integration test | Verifikasi stack, response envelope, JWT RS256 |
| A-2: Clean Architecture | Unit + code review | Verifikasi struktur direktori dan layer separation |
| A-3: Shared Go Module | Unit test | Verifikasi fungsi auth, RBAC, audit, idempotency |
| A-4: Integration Worker | Integration test + event contract test | Verifikasi outbox/inbox, retry, DLQ idempotency |
| A-5: Next.js App Router | E2E test (Playwright/Cypress) | Verifikasi routing, TanStack Query, degraded state |
| B-1: Registrasi OAuth Client | Integration test + property test | Round-trip: register → status PENDING; error duplikat |
| B-2: Persetujuan Admin | Integration test | State machine: PENDING → ACTIVE/REVOKED |
| B-3: Auth Code Flow + PKCE | Integration test + security test | PKCE mandatory, single-use code, exact URI match |
| B-4: Client Credentials Flow | Integration test | Validasi scope, reject suspended/revoked |
| B-5: Token Management | Property test + integration test | Token TTL, refresh rotation, rate limit |
| B-6: Introspection + Revocation | Integration test | Active/inactive response, cascade revoke |
| B-7: OIDC Discovery | Integration test | Dokumen discovery valid, JWKS tersedia |
| B-8: Entitas core_db | Database migration test | Schema valid, index, constraint |
| B-9: Keamanan OAuth | Security test | bcrypt, exact match URI, PKCE, replay detection |
| B-10: UI Admin | E2E test | Filter, konfirmasi aksi, aksi sesuai status |
| B-11: UI Developer Credentials | E2E test | One-time secret display, status pending |
| B-12: Integrasi Modul Internal | Integration test | Scope-based access, fallback ke JWKS |

