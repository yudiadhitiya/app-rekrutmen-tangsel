# API Contract — Rekrutmen Tangsel

Base URL: `http://localhost:8080/api`
All responses are JSON. Auth via session cookie (`rekr_session`, HttpOnly).
Roles: `admin` (everything), `petugas` (CRUD rekrutmen only).

## Auth

| Method | Path | Body | Response |
|---|---|---|---|
| POST | `/login` | `{username, password}` | `{user:{id,username,nama,role}}` |
| POST | `/logout` | — | `{ok:true}` |
| GET | `/me` | — | `{user:{id,username,nama,role}}` |
| PUT | `/me/password` | `{current_password, new_password}` | `{ok:true}` |

401 on bad credentials or missing/invalid session.

`PUT /me/password` mengubah password akun yang sedang masuk. Semua role boleh
mengubah passwordnya sendiri; tidak ada pembatasan level. `current_password`
wajib benar (400 bila salah), `new_password` minimal 12 karakter, maksimal 200
karakter, tidak boleh berupa password yang umum (misalnya yang memuat
"password", "admin123", "petugas123", atau "123456"), tidak boleh sama dengan
username akun, dan harus berbeda dari password lama (400 bila tidak).
Password baru ikut ditulis ke tab `USERS` bila akunnya ada di aplikasi
sebelumnya, jadi satu password berlaku di kedua aplikasi.

Aturan password yang sama berlaku pada pendaftaran mandiri: `POST /register`
menolak password yang terlalu pendek, yang umum, atau yang sama dengan
username.

## Master wilayah (read-only, all roles)

`GET /wilayah` →
```json
{ "kecamatan": [
  { "id": "Serpong", "nama": "Serpong",
    "kelurahan": [ {"id":"Buaran","nama":"Buaran"} ] }
]}
```
7 kecamatan (Serpong, Serpong Utara, Pondok Aren, Ciputat, Ciputat Timur,
Pamulang, Setu), 54 kelurahan total.

## Dashboard

`GET /dashboard` →

Menampilkan rekapitulasi data pendaftaran. Endpoint ini hanya mengembalikan
**agregat**, tidak pernah satu record pun, jadi tidak bisa dipakai untuk
mengakses data akun lain.

Query `filter` berlaku untuk anggota (level 1):

| Filter | Arti |
|---|---|
| `ALL` (bawaan) | seluruh kecamatan |
| `KECAMATAN` | hanya kecamatan akun |
| `ANGGOTA` | hanya record yang dibuat akun itu sendiri |

Cakupan mengikuti level akun: level 3 = seluruh data, level 2 = kecamatannya
saja, level 1 = mengikuti filter di atas.

```json
{ "scope": "Seluruh data pendaftaran",
  "filter": "ALL",
  "total": 120,
  "methods": [ {"name":"Halaqoh","total":50,"percent":41.7} ],
  "methodSummary": [ {"name":"Halaqoh","total":50,"berhasil":30,"proses":15,
                       "gagalPending":5,"percent":41.7} ],
  "acceptance": [ {"name":"Terbuka","total":60,"percent":50} ],
  "status": [ {"name":"Berhasil","total":70,"percent":58.3} ],
  "statusDetail": [ {"name":"Berhasil","total":70,"percent":58.3} ],
  "kecamatanSummary": [ {"kecamatan":"Serpong","target":30,"berhasil":20,
                          "proses":6,"persenBerhasil":66.7,"persenProses":20} ],
  "top10": [ {"nama":"Ahmad Fauzi","total":3} ] }
```

`kecamatanSummary` selalu memuat 7 kecamatan Tangsel, walau belum ada data.

## Master aktivitas (job/activity categories)

| Method | Path | Role | Notes |
|---|---|---|---|
| GET | `/dashboard` | any | rekapitulasi data (agregat, tanpa paginasi) — query `filter` |
| GET | `/aktivitas` | any | `[{id,nama,aktif}]` sorted by nama |
| POST | `/aktivitas` | admin penuh (level 3) | `{nama}` → created |
| PUT | `/aktivitas/:id` | admin penuh (level 3) | `{nama, aktif}` |
| DELETE | `/aktivitas/:id` | admin penuh (level 3) | 409 if referenced by a rekrutmen record |

## Rekrutmen

| Method | Path | Role | Notes |
|---|---|---|---|
| GET | `/rekrutmen` | any | query: `q`, `kecamatan`, `kelurahan`, `metode`, `status`, `page`, `limit` |
| GET | `/rekrutmen/:id` | any | |
| POST | `/rekrutmen` | any | body = record, no `id` |
| PUT | `/rekrutmen/:id` | any | |
| DELETE | `/rekrutmen/:id` | anggota (data sendiri) atau admin penuh (level 3) | admin kecamatan (level 2) → 403 |
| GET | `/rekrutmen/check-wa?no=628123` | any | `{exists:bool, id?:string}` — warn-before-save check |
| GET | `/rekrutmen/export` | any | unduh data — query filter sama seperti `/rekrutmen` (tanpa paginasi); dibatasi scope level |
| POST | `/rekrutmen/batch` | anggota (level 1) | upload massal — array baris; lihat § Upload massal |

`GET /rekrutmen` →
```json
{ "total": 120, "page": 1, "limit": 25,
  "data": [ { ...record } ] }
```

Record:
```json
{
  "id": "uuid",
  "nama": "Ahmad Fauzi",
  "no_wa": "628123456789",
  "alamat": "Jl. Melati No. 5",
  "rt": "003", "rw": "005",
  "kecamatan": "Serpong", "kelurahan": "Buaran",
  "aktivitas": "Wiraswasta",
  "tingkat_penerimaan": "Terbuka",
  "metode": "Halaqoh",
  "status": "Berhasil",
  "keterangan": "",
  "upl": "UPLS",
  "created_at": "2026-11-28T22:00:00Z",
  "updated_at": "2026-11-28T22:00:00Z",
  "created_by": "admin", "updated_by": "admin"
}
```

Enums:
- `tingkat_penerimaan`: `Sangat Terbuka` | `Terbuka` | `Netral` | `Kurang Terbuka` | `Tertutup`
- `metode`: `Halaqoh` | `Program Antara` | `Dakwah Fardiyah`
- `status`: `Berhasil` | `Proses` | `Pending` | `Gagal`
- `upl`: `""` | `UPLS` | `UPLM`

### Unduh data (`GET /rekrutmen/export`)

Mengembalikan baris siap unduh (bukan paginasi), mengikuti filter yang
dikirim. Cakupan data ditentukan oleh **level** akun, sama seperti list:

| Level | Cakupan export |
|---|---|
| 1 (anggota) | hanya record yang dibuat sendiri |
| 2 (admin kecamatan) | hanya record kecamatannya |
| 3 (admin penuh) | semua record |

```json
{ "headers": ["NAMA ANGGOTA", "NAMA TARGET", "..."],
  "rows": [["petugas", "Ahmad Fauzi", "628123456789", "..."]],
  "total": 1,
  "filename": "Data_Rekrutmen_Tangsel_20261114_101500.xlsx" }
```

Frontend menyusun file XLSX dari `headers` + `rows` (SheetJS, dimuat dari
CDN). Bila pustaka tidak termuat, file diunduh sebagai CSV.

### Upload massal (`POST /rekrutmen/batch`)

Hanya untuk **anggota** (level 1). Body berisi array baris, masing-masing
satu objek dengan nama kolom (nama field aplikasi ini atau header template
aplikasi lama). Pemilik dan timestamp diambil dari sesi, jadi tidak bisa diisi
dari file.

| Method | Path | Role | Notes |
|---|---|---|---|
| POST | `/rekrutmen/batch` | anggota (level 1) | maks. 1.000 baris, maks. 5 MB |

Berhasil (`201`):
```json
{ "created": 2,
  "skipped": [ { "row": 4, "nama": "Siti Aminah", "no_wa": "628123456780",
                 "reason": "nomor whatsapp sudah terdaftar" } ] }
```

Gagal validasi (`400`) — **tidak ada baris yang disimpan**:
```json
{ "error": "validasi upload gagal. tidak ada data yang disimpan.",
  "error_count": 1,
  "errors": [ { "row": 3, "field": "no_wa", "reason": "hanya boleh angka 0-9" } ] }
```

Batas percobaan (`429`): satu akun boleh melakukan 5 upload per 10 menit.
Jumlahnya dapat dinaikkan lewat environment `UPLOAD_MAX_ATTEMPTS`.

## Settings (admin only)

| Method | Path | Notes |
|---|---|---|
| GET | `/settings` | `{provider:"local"|"supabase"|"google_sheets", ...}` |
| PUT | `/settings` | `{provider}` — switches active provider; does NOT move data |

## Validation rules (enforced on both sides)

- `nama` — required, non-blank after trim
- `no_wa` — required, digits `0-9` only, stored as string, **unique** across records
- `rt`, `rw` — optional, digits `0-9` only
- `kecamatan` — required, must be one of the 7
- `kelurahan` — optional, but if given must belong to `kecamatan`
- `aktivitas` — required, must exist in master aktivitas
- `tingkat_penerimaan`, `metode`, `status` — required, must be in enum above
- `metode` — the retired values `Program` and `Antara` are still accepted on
  write and on filter, and are normalized to `Program Antara`
- `upl` — optional, must be `""`, `UPLS`, or `UPLM`

### Error shape
```json
{ "error": "validasi gagal", "fields": { "no_wa": "hanya boleh angka 0-9" } }
```
- 400 validation → body has `fields`
- 401 unauthenticated → `{error:"..."}`
- 403 wrong role → `{error:"..."}`
- 404 not found
- 409 conflict (duplicate no_wa, delete in use)

## Frontend routes (vue-router)

| Path | View | Access |
|---|---|---|
| `/login` | LoginView | public |
| `/dashboard` | DashboardView | any |
| `/` | ListView (redirect `/`) | any |
| `/tambah` | FormView (create) | any |
| `/ubah/:id` | FormView (edit) | any |
| `/aktivitas` | AktivitasView (master CRUD) | any read; admin write |
| `/pengaturan` | SettingsView | admin |
