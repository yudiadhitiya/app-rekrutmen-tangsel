# Testing Plan — Fitur Ubah Password (Change Password)

**Status dokumen:** selesai dieksekusi — seluruh test otomatis lulus
**Terakhir diperbarui:** 2026-11-16
**Hasil eksekusi:** 18/18 test backend lulus (11 API + 7 unit auth),
7/7 test E2E lulus, 3/4 uji manual lulus (MN-04 terbatas pada jalur otomatis),
5 catatan ditemukan
**Fitur:** `PUT /api/me/password` → tombol "Password" di header
(`frontend/src/App.vue`, `frontend/src/components/PasswordDialog.vue`)
**Kode terkait:**
`backend/internal/api/api.go` (`handleChangePassword`),
`backend/internal/auth/auth.go` (`Manager.ChangePassword`),
`backend/internal/store/sheets/users.go` (`SetLegacyPassword`),
`backend/internal/validate/validate.go` (`Password`),
`frontend/src/api.js` (`api.password.change`)
**Dokumen induk:** [`TESTING_PLAN.md`](TESTING_PLAN.md) (plan repo-wide),
[`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md) §13 (D-31, D-84),
[`API_CONTRACT.md`](API_CONTRACT.md) § "Auth"
**Plan fitur lain:** [`testing_plan_export.md`](testing_plan_export.md) (Unduh Data)

---

## 1. Ruang lingkup

Yang diuji:

- Endpoint `PUT /api/me/password` (backend): otorisasi, validasi, dan efek samping.
- `auth.Manager.ChangePassword`: verifikasi password lama, penulisan hash baru
  (PBKDF2 + salt acak) ke `data/users.json`, dan penulisan hash ke tab `USERS`
  aplikasi lama.
- Dialog "Ubah Password" di frontend: validasi konfirmasi, pesan galat, toast,
  dan alur masuk ulang dengan password baru.

Di luar lingkup:

- Pemulihan password lewat pertanyaan keamanan (hanya ada di aplikasi lama).
- Rate limiting / lockout brute-force pada endpoint ini (terpisah, lihat D-83).
- Penghapusan password default (D-85).
- Perilaku aplikasi lama (`old-app`) selain kompatibilitas hash di tab `USERS`.

## 2. Kriteria lulus (exit criteria)

| # | Kriteria | Hasil |
|---|---|---|
| 1 | Semua test otomatis lulus (`go test ./...`, Playwright) | ✅ `go test ./...` hijau, Playwright 7/7 lulus |
| 2 | Semua level akun (1, 2, 3) dapat mengganti passwordnya sendiri | ✅ BE-02/03/04, E2E-01 |
| 3 | Password lama salah → ditolak (400) dan password lama tetap berlaku | ✅ BE-05 |
| 4 | Password baru minimal 6 karakter, harus berbeda dari yang lama | ✅ BE-06/07/08 |
| 5 | Password baru berlaku untuk masuk berikutnya; password lama ditolak | ✅ BE-11, E2E-07/08 |
| 6 | Sesi yang sedang berjalan tidak diputus setelah password diganti | ✅ BE-10, MN-02 |
| 7 | Hash baru tidak pernah dikirim ke browser (respons hanya `{ok:true}`) | ✅ respons `{"ok":true}` |
| 8 | Setiap temuan di §6 selesai diperbaiki atau diputuskan diterima | ✅ semua temuan diputuskan |

## 3. Prasyarat

- Backend + frontend berjalan: `cd backend && ./run-dev.sh`, `cd frontend && npm run build`.
- Akun uji: `admin` (level 3), `petugas` (level 1), `kec` (level 2, kecamatan Serpong).
- Akun uji khusus E2E (`pwpetugas`, `pwkec`, `pwadmin`) supaya password akun
  bawaan tidak berubah dan test lain tidak terganggu. Akun ini dibuat dengan
  helper `backend/cmd/e2eseed` (lihat §7).
- Browser Chromium (Playwright terpasang di root repo).

## 4. Matriks test

### 4.1 Backend API (Go, `backend/internal/api/password_test.go`)

| ID | Skenario | Ekspektasi | Status | Test |
|---|---|---|---|---|
| BE-01 | Ganti password tanpa sesi | 401 `sesi tidak valid` | ✅ | `TestChangePasswordNeedsSession` |
| BE-02 | Anggota level 1 mengganti passwordnya | 200 `{ok:true}` | ✅ | `TestChangePasswordOfMemberAccount` |
| BE-03 | Admin kecamatan level 2 mengganti passwordnya | 200 `{ok:true}` | ✅ | `TestChangePasswordOfEveryLevel` |
| BE-04 | Admin penuh level 3 mengganti passwordnya | 200 `{ok:true}` | ✅ | `TestChangePasswordOfEveryLevel` |
| BE-05 | Password lama salah | 400 `password lama salah` | ✅ | `TestChangePasswordRejectsWrongCurrentPassword` |
| BE-06 | Password baru < 6 karakter | 400 + `fields.password` | ✅ | `TestChangePasswordRejectsShortPassword` |
| BE-07 | Password baru sama dengan password lama | 400 | ✅ | `TestChangePasswordRejectsSamePassword` |
| BE-08 | `current_password` kosong | 400 | ✅ | `TestChangePasswordNeedsCurrentPassword` |
| BE-09 | Body JSON tidak valid | 400 `body tidak valid` | ✅ | `TestChangePasswordRejectsBrokenBody` |
| BE-10 | Sesi tetap berlaku setelah password diganti | `GET /api/me` tetap 200 | ✅ | `TestChangePasswordKeepsTheSession` |
| BE-11 | Password baru dipakai untuk masuk lagi, password lama ditolak | 200 / 401 | ✅ | `TestChangePasswordIsUsedByTheNextLogin` |
| BE-12 | Akun yang hanya ada di tab `USERS` (warisan) | hash baru ikut ditulis ke store warisan | ✅ | `TestChangePasswordUpdatesTheAccountOfThePreviousApplication` |

### 4.2 Backend auth (unit, `backend/internal/auth/auth_test.go`)

| ID | Skenario | Ekspektasi | Status | Test |
|---|---|---|---|---|
| UT-01 | Hash lokal diperbarui dan disimpan ke `users.json` | file memuat hash baru | ✅ | `TestChangePasswordUpdatesTheLocalAccount` |
| UT-02 | Password lama ditolak setelah diganti | `ErrInvalidCredentials` | ✅ | `TestChangePasswordRejectsWrongCurrentPassword` |
| UT-03 | Mode `sheet`: akun hanya di sheet, password baru ditulis ke sheet | setter warisan dipanggil | ✅ | `TestChangePasswordInSheetOnlyMode` |
| UT-04 | Akun tidak dikenal di store mana pun | `ErrAccountNotFound` | ✅ | `TestChangePasswordOfUnknownAccount` |
| UT-05 | Akun ada di lokal dan di sheet | keduanya diperbarui | ✅ | `TestChangePasswordUpdatesBothStores` |

### 4.3 E2E (Playwright, `tests/e2e/ubah-password.spec.ts`)

| ID | Skenario | Ekspektasi | Status | Test |
|---|---|---|---|---|
| E2E-01 | Tombol "Password" terlihat untuk anggota, admin kecamatan, dan admin penuh | tombol ada | ✅ | `Tombol Password tersedia untuk semua level akun` |
| E2E-02 | Dialog dibuka dari header | form 3 field (lama, baru, ulang) | ✅ | `Dialog mengubah password menampilkan tiga field` |
| E2E-03 | Konfirmasi password tidak sama | pesan galat, tidak ada permintaan | ✅ | `Konfirmasi password yang tidak sama ditolak di sisi klien` |
| E2E-04 | Password lama salah | pesan `password lama salah` dari server | ✅ | `Password lama yang salah ditolak oleh server` |
| E2E-05 | Password baru < 6 karakter | pesan validasi klien | ✅ | `Password baru terlalu pendek ditolak di sisi klien` |
| E2E-06 | Ganti password berhasil | toast sukses, dialog tertutup | ✅ | `Password berhasil diganti dan dipakai untuk masuk lagi` |
| E2E-07 | Keluar lalu masuk lagi dengan password baru | masuk berhasil | ✅ | `Password berhasil diganti dan dipakai untuk masuk lagi` |
| E2E-08 | Masuk lagi dengan password lama | ditolak | ✅ | `Password berhasil diganti dan dipakai untuk masuk lagi` |
| E2E-09 | Dialog ditutup dengan tombol Batal, isian dikosongkan | form tertutup, isian kosong | ✅ | `Dialog bisa ditutup dan isiannya dikosongkan` |

### 4.4 Manual / eksploratif

| ID | Skenario | Ekspektasi | Status | Cara uji |
|---|---|---|---|---|
| MN-01 | Ketik password di dialog | karakter tersamar (type=password) | ✅ | Playwright: `type` ketiga field = `password` |
| MN-02 | Sesi di browser lain | tetap aktif (tidak di-invalidate) | ✅ | 2 cookie jar (curl): ganti password di sesi A → `GET /api/me` sesi B tetap 200 |
| MN-03 | Tutup dialog dengan tombol Batal / klik luar | form tertutup, isian dikosongkan | ✅ | Playwright: klik backdrop → dialog hilang; buka lagi → ketiga field kosong |
| MN-04 | Ganti password pada akun yang juga ada di aplikasi lama | bisa masuk di kedua aplikasi dengan password baru | ⚠️ | Terverifikasi otomatis lewat BE-12/UT-03/UT-05 (fake store warisan). Uji terhadap Google Sheets sungguhan tidak dijalankan — menulis ke sheet produksi butuh persetujuan |

## 5. Hasil uji manual yang sudah dijalankan

Tanggal: 2026-11-16. Server uji: `ADDR=:8100`, `DATA_DIR=/tmp/e2e-pw-data`,
`LEGACY_AUTH_MODE=off`, frontend dari `frontend/dist`.

| ID | Langkah yang dijalankan | Hasil yang dilihat |
|---|---|---|
| MN-01 | Buka dialog, periksa `type` ketiga field | `password`, `password`, `password` — tidak pernah ditampilkan polos |
| MN-02 | Masuk dari dua "browser" (dua cookie jar), ubah password di yang satu | Sesi lain tetap `200` pada `GET /api/me`; login dengan password lama `401`, dengan password baru `200` |
| MN-03 | Isi sebagian form → klik "Batal" / klik area di luar dialog | Dialog tertutup; dibuka lagi ketiga field kosong; pesan galat lama hilang |
| MN-04 | Jalur warisan (`SetLegacyPassword` → SHA-256 di tab `USERS`) | Lulus di test otomatis (BE-12, UT-03, UT-05); belum diuji ke Sheets sungguhan |

Catatan tambahan:

- Respons `PUT /api/me/password` hanya `{"ok":true}` — hash baru tidak pernah
  dikirim ke browser (kriteria §2 no. 7 terpenuhi).
- Sesi yang sedang berjalan tetap berlaku setelah penggantian password
  (tidak ada invalidasi token).
- Pesan galat yang terverifikasi: `password lama salah` (server),
  `Password baru minimal 6 karakter.`, `Konfirmasi password tidak sama.`,
  `Password baru harus berbeda dari password lama.` (klien).

## 6. Temuan

| ID | Temuan | Dampak | Saran |
|---|---|---|---|
| T-1 | Hash warisan ditulis ke tab `USERS` sebagai SHA-256 tanpa salt (`sha256Hex`), mengikuti format aplikasi lama | Rendah — format lama tidak bisa diubah tanpa memindahkan data sheet | Diterima: pertahankan kompatibilitas; beri tahu pengguna agar mengganti password di kedua aplikasi |
| T-2 | Endpoint `PUT /api/me/password` belum punya rate limiting / lockout (D-83) | Sedang — endpoint bisa dipakai untuk menebak password lama | Terapkan rate limiting sesuai D-83 pada tahap terpisah |
| T-3 | Ganti password tidak mencabut sesi lain di perangkat lain | Rendah — sesuai perilaku standar; password lama yang bocor masih bisa dipakai sampai session berakhir | Tambahkan opsi "keluar dari semua perangkat" bila dibutuhkan |
| T-4 | Dialog tidak punya tombol "lihat password"; tidak ada indikator kekuatan password | Rendah — kualitas UX | Tambahkan show/hide password (opsional, bukan blokir rilis) |
| T-5 | E2E bergantung pada akun uji yang harus disiapkan manual | Sedang — test gagal di mesin baru | Buka `backend/cmd/e2eseed` untuk membuat akun uji, lalu tambahkan skrip setup E2E (langkah lanjutan, di luar cakupan dokumen ini) |

Catatan: T-3 dan T-5 pada tabel di atas hanya catatan; tidak ada perubahan kode
yang dilakukan pada tahap pengujian ini.

## 7. Perintah verifikasi

Perintah yang benar-benar dijalankan pada eksekusi 2026-11-16:

```bash
# Backend — build, vet, dan seluruh test
cd backend && go build ./... && go vet ./...
cd backend && go test ./...                 # ✅ semua paket hijau
cd backend && go test ./... -run Password -v   # ✅ 18 test password lulus

# Frontend
cd frontend && npm run build               # ✅ sukses

# Akun uji E2E (level 1, 2, 3)
cd backend && go run ./cmd/e2eseed /tmp/e2e-pw-data

# Server uji untuk E2E
cd backend && ADDR=:8100 DATA_DIR=/tmp/e2e-pw-data \
  FRONTEND_DIR=../frontend/dist LEGACY_AUTH_MODE=off \
  SESSION_SECRET=e2e-test-secret go run ./cmd/server

# E2E
cd .. && npx playwright test tests/e2e/ubah-password.spec.ts   # ✅ 7/7 lulus
```

## 8. Urutan pengerjaan test yang disarankan

1. Backend unit + API (BE-01…BE-12, UT-01…UT-05) — cepat, tanpa infra baru. ✅
2. E2E (E2E-01…E2E-09) — butuh server + build frontend + akun uji khusus. ✅
3. Manual (MN-01…MN-04). ✅ (MN-04 terbatas pada jalur otomatis)
4. Perbarui status di dokumen ini dan di [`TESTING_PLAN.md`](TESTING_PLAN.md). ✅
