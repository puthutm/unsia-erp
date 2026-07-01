# ERD Final - Portal Module

## 1. Skema Database
Database `portal_db` mengelola notifikasi pengguna, penandaan notifikasi dibaca, preferensi personalisasi layout, pintasan halaman (*shortcuts*), serta read model teragregasi dashboard.

## 2. Tabel Utama
* **`portal.notifications`**: Rekaman log pesan notifikasi masuk.
* **`portal.notification_reads`**: Status penanda notifikasi yang telah dibaca.
* **`portal.user_preferences`**: Pengaturan visual profil (tema/bahasa).
* **`portal.menu_shortcuts`**: Pintasan akses modul cepat dashboard.

## 3. Script DBML (Copy-paste ke dbdiagram.io)
```dbml
Table portal.notifications {
  id uuid [pk]
  user_id uuid [not null]
  title varchar
  message text
  module_source varchar
  target_url text
  sent_at timestamp
}

Table portal.notification_reads {
  id uuid [pk]
  notification_id uuid [not null]
  user_id uuid [not null]
  read_at timestamp
}

Table portal.user_preferences {
  id uuid [pk]
  user_id uuid [not null]
  preference_key varchar
  preference_value jsonb
  updated_at timestamp
}

Table portal.menu_shortcuts {
  id uuid [pk]
  user_id uuid [not null]
  menu_code varchar
  menu_label varchar
  target_url text
  sort_order int
}

Ref: portal.notification_reads.notification_id > portal.notifications.id
```

## 4. Hubungan Logis Lintas Modul (Tanpa Database FK)
* `portal.notifications.user_id` merujuk logis ke `core_db.users.id`.
* `portal.menu_shortcuts.user_id` merujuk logis ke `core_db.users.id`.
