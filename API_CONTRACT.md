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

401 on bad credentials or missing/invalid session.

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

## Master aktivitas (job/activity categories)

| Method | Path | Role | Notes |
|---|---|---|---|
| GET | `/aktivitas` | any | `[{id,nama,aktif}]` sorted by nama |
| POST | `/aktivitas` | admin | `{nama}` → created |
| PUT | `/aktivitas/:id` | admin | `{nama, aktif}` |
| DELETE | `/aktivitas/:id` | admin | 409 if referenced by a rekrutmen record |

## Rekrutmen

| Method | Path | Role | Notes |
|---|---|---|---|
| GET | `/rekrutmen` | any | query: `q`, `kecamatan`, `kelurahan`, `metode`, `status`, `page`, `limit` |
| GET | `/rekrutmen/:id` | any | |
| POST | `/rekrutmen` | any | body = record, no `id` |
| PUT | `/rekrutmen/:id` | any | |
| DELETE | `/rekrutmen/:id` | any | |
| GET | `/rekrutmen/check-wa?no=628123` | any | `{exists:bool, id?:string}` — warn-before-save check |

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
| `/` | ListView (redirect `/`) | any |
| `/tambah` | FormView (create) | any |
| `/ubah/:id` | FormView (edit) | any |
| `/aktivitas` | AktivitasView (master CRUD) | any read; admin write |
| `/pengaturan` | SettingsView | admin |
