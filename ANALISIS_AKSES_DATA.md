# Analisis: Admin tidak melihat semua data (Rekrutmen Tangsel)

**Ruang lingkup:** analisis + implementasi filter role. Sumber: `old-app/` (logika lama), `backend/`, `frontend/`.

**Status implementasi (selesai):**
- Role `admin` → melihat & mengubah semua data. Role `petugas` → hanya data dengan `created_by` = username-nya (daftar, detail, ubah). Hapus hanya untuk admin.
- Perubahan: `model.Filter.CreatedBy`, filter di 3 store (local/sheets/supabase), `canAccess()` di `api.go`, rute `DELETE` → `requireAdmin`, tombol "Hapus" disembunyikan untuk petugas.
- **Belum dikerjakan:** perluasan `LegacyRole` (ADMIN KECAMATAN), importer header lama, pemecahan UPLS/UPLM, dan pengisian `created_by` untuk data impor.


---

## 1. Ringkasan Eksekutif

Ada **dua hal terpisah** yang tercampur dalam keluhan ini:

1. **Aplikasi baru tidak punya-filter hak akses sama sekali** untuk daftar data. `handleListRekrutmen` meneruskan filter dari query apa adanya ke store, tanpa membatasi berdasarkan role. Artinya **petugas juga melihat semua data** — bukan hanya admin. Yang hilang justru pembedaan yang dulu ada.
2. **Data "hilang" bukan karena role, tapi karena sumber data & pemetaan kolom.** Aplikasi baru membaca tab `Rekrutmen` (header baru, lowercase) sedangkan data historis ada di tab `PENDAFTARAN` milik aplikasi lama dengan header berbeda (`NAMA TARGET`, `NO WHATSAPP`, `NAMA KECAMATAN`, `AKTIFITAS`, `TINGKAT PENERIMAAN`, `UPLS`/`UPLM`). Karena `readSheet` meng-index baris **memakai nama header apa adanya** dan `valuesToRekrutmen` meminta kunci lowercase persis, **tidak ada satu pun kolom lama yang cocok**.

Jadi: `system_admin` memang sudah.role `admin` (lihat §3), tetapi data yang ia minta tidak ada di tab yang dibaca aplikasi.

---

## 2. Logika yang ada di `old-app` (yang seharusnya dipertahankan)

### 2.1 Model role lama — `old-app/code.gs:6-11`
```js
ROLES: {
  ANGGOTA: 'ANGGOTA',
  ADMIN_KECAMATAN: 'ADMIN KECAMATAN',
  ADMIN_DPD: 'ADMIN PTD',
  SYSTEM_ADMIN: 'SYSTEM ADMIN'
}
```
Session menyimpan `username`, `role`, dan `kecamatan` (`code.gs:210-213`).

### 2.2 Aturan akses data — `old-app/code.gs:373-392` (dari `getRecruitmentsPaged`)
| Role | Aturan | Sumber |
|------|--------|--------|
| ANGGOTA | `createdBy === myUsername` — hanya data yang diinput sendiri | `code.gs:382-383` |
| ADMIN KECAMATAN | `rowKecamatan === myKecamatan` — kolom `NAMA KECAMATAN` (index 4) | `code.gs:384-385` |
| ADMIN PTD | `allowed = true` — semua data | `code.gs:387-388` |
| SYSTEM ADMIN | `allowed = true` — semua data | `code.gs:388-389` |

Dokumentasi resmi di `code.gs:317-323`:
```
* - ANGGOTA   : hanya data yang dibuat/diinput oleh username sendiri.
* - ADMIN KECAMATAN : hanya data dengan NAMA KECAMATAN yang sama dengan KECAMATAN user.
* - ADMIN DPD : semua data.
* - SYSTEM ADMIN: semua data (read-only untuk data pendaftaran).
```
Penjelasan yang sama diulang di `old-app/index.html:531-548` untuk kolom "scope" pada UI.

### 2.3 Sumber data lama
`old-app/code.gs:2` → `DATA_SHEET: 'PENDAFTARAN'`, dengan 20 kolom (`code.gs:11-19`):
```
ID, NAMA ANGGOTA, NAMA TARGET, NO WHATSAPP, NAMA KECAMATAN, ALAMAT, RT, RW,
KECAMATAN, KELURAHAN, UPLS, UPLM, AKTIFITAS, TINGKAT PENERIMAAN, METODE,
STATUS, KETERANGAN, CREATED_AT, UPDATED_AT, CREATED_BY
```
Baris kosong dan baris tanpa `NAMA TARGET` dilewati (`code.gs:364-369`).

### 2.4 Perilaku hak akses lain di UI lama (`old-app/index.html`)
- `index.html:523-524` — filter dashboard hanya untuk ANGGOTA.
- `index.html:829, 893-896` — ANGGOTA dan ADMIN KECAMATAN tidak boleh menghapus data; hanya ADMIN PTD/SYSTEM ADMIN yang punya tombol Hapus.
- `index.html:1159-1160` — field "Nama Anggota" & "Kecamatan" hanya relevan untuk role tertentu.

---

## 3. Kondisi aplikasi baru (verified)

| Aspek | Implementasi baru | Lokasi |
|-------|-------------------|--------|
| Role yang tersedia | hanya `admin` dan `petugas` | `backend/internal/model/model.go:99-110` |
| Pemetaan role lama → baru | `SYSTEM ADMIN`, `ADMIN PTD`, `ADMIN DPD`, `ADMIN KECAMATAN` → **admin**; sisanya → petugas | `backend/internal/store/sheets/users.go:47-54` |
| Data user yang disimpan | hanya `ID, Username, Nama, Role` — **tidak ada `kecamatan`** | `model.go:99-104` |
| Filtering daftar | tidak ada; `model.Filter` hanya berisi q/kecamatan/kelurahan/metode/status/page/limit | `api.go:453-487`, `model.go:88-96` |
| Otorisasi | `requireAdmin` hanya untuk hapus & pengaturan | `api.go:206-214`, `api.go:252-253` |
| Store | `local` & `sheets` menerapkan hanya filter query | `store/local/store.go:98-150`, `store/sheets/store.go:787-838` |

**Konsekuensi (dua arah):**
- `system_admin` → `admin` ✅ benar, dan **tetap bisa melihat semua baris** karena memang tidak ada pembatasan. Jadi peran admin **bukan** penyebab data tidak lengkap.
- Sebaliknya, **ADMIN KECAMATAN lama kini diperlakukan sebagai admin biasa** → kini ia melihat data seluruh kecamatan (regresi hak akses), dan `ANGGOTA` kini melihat semua data, bukan hanya miliknya sendiri.

---

## 4. Mengapa data tidak "keluar semua" — 4 penyebab terverifikasi

### (a) Tab yang dibaca ≠ tab yang berisi data historis — **penyebab utama**
- Default tab data aplikasi baru: `Rekrutmen` (`store/sheets/store.go:34`, `DEPLOYMENT_PLAN.md:118`).
- Data historis ada di tab `PENDAFTARAN` (aplikasi lama). Aplikasi baru **tidak pernah membaca tab itu**; tab `Rekrutmen` dibuat otomatis saat kosong (`store.go:219-240`).
- Akibatnya hanya baris yang dibuat lewat aplikasi baru (atau hasil impor) yang muncul.

### (b) Header lama tidak cocok dengan kunci yang dibaca — **data kosong buta**
- `readSheet` menyimpan baris dengan **nama header apa adanya**: `row[h] = vals[i]` (`store/sheets/store.go:626`).
- `valuesToRekrutmen` membaca **kunci lowercase persis**: `row["nama"]`, `row["no_wa"]`, `row["kecamatan"]`, `row["aktivitas"]`, `row["tingkat_penerimaan"]`, `row["upl"]`, … (`store/sheets/store.go:744-768`).
- Tab `PENDAFTARAN` memakai `NAMA TARGET`, `NO WHATSAPP`, `NAMA KECAMATAN`, `AKTIFITAS`, `TINGKAT PENERIMAAN`, `UPLS`, `UPLM` → **tidak satu pun cocok** (`'ID' != 'id'`). Bila tab diarahkan ke `PENDAFTARAN`, baris muncul tetapi semua field kosong.
- `SHEETS_INTEGRATION_PLAN.md:156` (D-128) sudah menandai kebutuhan ini: konversi harus mapping-aware, bukan urutan tetap.

### (c) Kolom `UPL` terpecah di data lama
- Lama: dua kolom terpisah `UPLS` dan `UPLM`.
- Baru: satu field `upl` dengan nilai `''`/`UPLS`/`UPLM` (`frontend/src/constants.js:27-31`, `mapping.go:78`).
- Perlu aturan: ambil kolom yang terisi.

### (d) Importer tidak mengenal header lama
- `FormColumnSpec` (`store/sheets/mapping.go:63-80`) tidak punya alias `nama target` maupun `aktifitas`; `nama` ditandai **Required**.
- Impor tab `PENDAFTARAN` akan gagal dengan `kolom wajib tidak ditemukan: nama` — atau kehilangan nama & aktivitas bila dipaksa jalan.

### (e) Mode `form` tidak menyuap data sumber ke daftar
- `ListRekrutmen` selalu membaca `s.tabRekrutmen` (`store.go:789`), **`tabSumber` tidak pernah dipakai di jalur daftar**; `tabSumber` hanya dipakai oleh tool `impor-responses`.
- Konsekuensi: di mode `form`, daftar hanya menampilkan hasil impor di tab `Rekrutmen`.

### Cara verifikasi cepat (2 menit, tanpa kode)
1. Buka **Pengaturan → Informasi Provider** (dump JSON): cek `active_provider` dan `google_sheets.rekrutmen` → tab mana yang dibaca.
2. `curl -b cookie -H 'Accept: application/json' 'http://localhost:8080/api/rekrutmen?limit=100'` → bandingkan `total` vs `data.length`.
3. Buka spreadsheet: hitung baris di tab `PENDAFTARAN` vs tab `Rekrutmen`. Kalau `PENDAFTARAN` jauh lebih besar → penyebab (a).

---

## 5. Rekomendasi pemulihan logika lama (belum diimplementasikan)

### 5.1 Kembalikan granularity role
- `model.User` perlu tambahan: `RoleDetail` (atau `Level`) dan `Kecamatan`.
- Pemetaan dari role lama (`users.go:47-54`) diperluas:
  | Role lama | Level | Scope data |
  |-----------|-------|------------|
  | SYSTEM ADMIN | 3 | semua |
  | ADMIN PTD / ADMIN DPD | 3 | semua |
  | ADMIN KECAMATAN / TARBIYAH KECAMATAN | 2 | `kecamatan` miliknya |
  | ANGGOTA | 1 | `created_by` = username miliknya |
- Level 3 (admin) tetap mendapat akses penuh → keluhan "admin harus lihat semua" satisfied.

### 5.2 Terapkan scope di API, bukan di store
- Tambahkan parameter scope pada `model.Filter` (mis. `ScopeByUsername`, `ScopeKecamatan`) atau lakukan penyaringan di `api.go:453-487` sebelum memanggil store.
- Filter role-old diterapkan **sebelum** pagination agar `total` juga benar.
- Jangan diubah di layer store: `store` tetap agnostik terhadap role.

### 5.3 Jaga hak hapus seperti lama
- `requireAdmin` saat ini mengizinkan hapus oleh semua admin; versi lama hanya ADMIN PTD/SYSTEM ADMIN (`index.html:893-896`). Pertimbangkan `requireAdmin(level>=3)`.

### 5.4 Migrasi data historis sekali jalan
Buat importer yang memetakan header lama → field baru:

| Header lama | Field baru | Catatan |
|-------------|------------|---------|
| `ID` | `id` | bila kosong, generate |
| `NAMA TARGET` | `nama` | **wajib**, baris kosong dilewati (perilaku lama) |
| `NO WHATSAPP` | `no_wa` | normalisasi digit (`mapping.go:202`) |
| `NAMA KECAMATAN` | `kecamatan` | kolom 5 (lihat catatan) |
| `KECAMATAN` | `kecamatan` | kolom 9 — pilih satu, jangan dua-duanya |
| `KELURAHAN` | `kelurahan` | |
| `AKTIFITAS` | `aktivitas` | typo lama harus diterima |
| `TINGKAT PENERIMAAN` | `tingkat_penerimaan` | |
| `UPLS` / `UPLM` | `upl` | ambil yang terisi |
| `CREATED_BY` | `created_by` | sumber untuk scoping ANGGOTA |
| `CREATED_AT` / `UPDATED_AT` | `created_at` / `updated_at` | |
| `ALAMAT`, `RT`, `RW`, `METODE`, `STATUS`, `KETERANGAN` | sama | |

Arah yang disarankan (sesuai `SHEETS_INTEGRATION_PLAN.md`): **impor sekali** `PENDAFTARAN` → tab `Rekrutmen` dengan schema aplikasi, lalu aplikasi tetap memakai tab `Rekrutmen`. Ini menghindari header warisan selamanya dan membuat `created_by` terisi (syarat scoping ANGGOTA).

### 5.5 Alias importer diperluas
Tambahkan ke `FormColumnSpec`: `nama target`, `nama anggota`, `aktifitas`, `upls`, `uplm`, `nama kecamatan` — agar impor tab lama tidak gagal.

---

## 6. Ringkasan akar masalah

| Item | Dampak | Prioritas |
|--------|--------|-----------|
| Tidak ada filter role di daftar data | petugas melihat semua data; scoping lama hilang | P1 |
| Tab data historis (`PENDAFTARAN`) tidak pernah dibaca | data lama tidak muncul sama sekali | P1 |
| Header warisan tidak cocok dengan kunci reader | data tampil kosong bila tab diarahkan ke `PENDAFTARAN` | P1 |
| `UPL` terpecah (UPLS/UPLM) | nilai UPL hilang saat impor | P2 |
| Alias importer tidak mengenal `NAMA TARGET`/`AKTIFITAS` | impor gagal / data kehilangan field | P2 |
| `created_by` tidak terisi pada data lama | scoping ANGGOTA tidak mungkin ditegakkan | P2 |
| Role `ADMIN KECAMATAN` dipromosikan menjadi `admin` | privilege escalation vs aplikasi lama | P1 |

> Catatan: dokumen ini murni analisis. Perbaikan kode (filter role, importer, alias kolom) dilakukan pada tahap terpisah setelah keputusan scope disepakati.
