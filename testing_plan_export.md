# Testing Plan — Fitur Unduh Data (Export Rekrutmen)

**Status dokumen:** draft, siap dieksekusi (belum ada kode test baru — masih tahap rencana)
**Terakhir diperbarui:** 2026-11-15
**Fitur:** `GET /api/rekrutmen/export` → tombol "⬇ Unduh Data" (`frontend/src/views/ListView.vue`)
**Kode terkait:** `backend/internal/api/export.go`, `backend/internal/api/export_test.go`,
`frontend/src/excel.js` (`downloadWorkbook`, `loadSheetJS`), `frontend/src/api.js`
(`api.rekrutmen.exportData`), `frontend/src/views/ListView.vue` (`exportData`)
**Dokumen induk:** [`TESTING_PLAN.md`](TESTING_PLAN.md) (plan repo-wide),
[`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md) §13 (D-197…D-206),
[`API_CONTRACT.md`](API_CONTRACT.md) § "Unduh data"

---

## 1. Ruang lingkup

Yang diuji:

- Endpoint `GET /api/rekrutmen/export` (backend).
- Penyusun XLSX + fallback CSV di frontend.
- Tombol "⬇ Unduh Data" di `ListView.vue` beserta filter aktif, toast, dan status tombol.

Di luar lingkup:heets/Supabase sebagai provider store (uji memakai store `local`),
perilaku Excel desktop (hanya isi file), Keamanan kata sandi.

## 2. Kriteria lulus (exit criteria)

1. Semua test otomatis di bawah lulus (`go test ./...`, Playwright).
2. Cakupan data (level 1/2/3) terbukti sama persis dengan `GET /api/rekrutmen`.
3. Isi file = header (17 kolom) + baris hasil filter, urut & benar.
4. Tidak ada data bocor antar akun pada file hasil unduhan.
5. Setiap temuan di §6 selesai diperbaiki atau memutuskan diterima.

## 3. Prasyarat

- Backend + frontend berjalan: `cd backend && ./run-dev.sh`, `cd frontend && npm run build`.
- Akun uji: `admin` (level 3), `petugas` (level 1), `kec` (level 2, kecamatan Serpong).
- Data uji minimal: 5 record (2 di kecamatan berbeda, 2 status berbeda, 1 record `upl=UPLM`).
- Browser Chromium (Playwright terpasang di root repo).

## 4. Matriks test

### 4.1 Backend (Go, `backend/internal/api/export_test.go`)

| ID | Skenario | Ekspektasi | Status |
|---|---|---|---|
| BE-01 | Export tanpa sesi | 401 `sesi tidak valid` | ✅ ada (lulus) |
| BE-02 | Anggota level 1 export | hanya record miliknya | ✅ ada (lulus) |
| BE-03 | Admin kecamatan level 2 export | hanya record kecamatannya | ✅ ada (lulus) |
| BE-04 | Admin penuh level 3 export | semua record | ✅ ada (lulus) |
| BE-05 | Filter `kecamatan` diterapkan | baris sesuai filter | ✅ ada (lulus) |
| BE-06 | Filter `status` diterapkan | baris sesuai filter | ✅ ada (lulus) |
| BE-07 | Tanpa paginasi | `limit=5&page=2` diabaikan, semua baris ikut | ✅ ada (lulus) |
| BE-08 | Filter `q` (nama/No WA) diterapkan | Like/contains, case-insensitive | ⬜ belum ada |
| BE-09 | Filter `kelurahan` & `metode` (nilai lama `Program`/`Antara`) | ternormalisasi, sama seperti list | ⬜ belum ada |
| BE-10 | Bentuk respons: 17 header, `rows` sel = jumlah header, `filename` pola | konsisten | ⬜ belum ada |
| BE-11 | Pemetaan kolom UPLS/UPLM (hanya satu kolom terisi, `""` → keduanya kosong) | benar | ⬜ belum ada |
| BE-12 | `total` = jumlah baris hasil filter (bukan jumlah sebelum filter) | benar | ⬜ belum ada |
| BE-13 | Batas 10.000 baris | terpotong + `total` tetap jumlah sebenarnya → lihat B-04 | ⬜ belum ada |
| BE-14 | Error store (backend tidak bisa membaca data) | 500 `gagal menyiapkan data unduhan` | ⬜ belum ada |

### 4.2 Frontend unit (belum ada runner test)

| ID | Skenario | Ekspektasi | Status |
|---|---|---|---|
| FE-01 | `toCSV` menutupi nilai ber-koma/ber-kutip/baris baru | hasil CSV valid | ⬜ belum ada |
| FE-02 | `downloadWorkbook` fallback saat SheetJS gagal dimuat | menghasilkan `.csv` + BOM UTF-8 | ⬜ belum ada |
| FE-03 | Nama file fallback mengganti `.xlsx` → `.csv` | benar | ⬜ belum ada |
| FE-04 | `activeFilters` mengirim `q, kecamatan, kelurahan, metode, status` | lengkap, tanpa `page/limit` | ⬜ belum ada |

> Catatan: `frontend/package.json` belum punya test runner (vitest). Perlu
> keputusan: tambah vitest atau cukup E2E Playwright.

### 4.3 E2E (Playwright, `tests/e2e/unduh-data.spec.ts`)

| ID | Skenario | Ekspektasi | Status |
|---|---|---|---|
| E2E-01 | Admin penuh klik "⬇ Unduh Data" | file `.xlsx` terunduh, nama `Data_Rekrutmen_Tangsel_YYYYMMDD_HHMMSS.xlsx` | ⬜ belum ada |
| E2E-02 | Isi XLSX | sheet `PENDAFTARAN`, baris 1 = 17 header, baris = jumlah record | ⬜ belum ada |
| E2E-03 | Anggota level 1 mengunduh | file hanya berisi record miliknya | ⬜ belum ada |
| E2E-04 | Filter aktif (pencarian/kecamatan/status) | isi file = hasil filter, sama dengan layar | ⬜ belum ada |
| E2E-05 | Hasil kosong | tidak ada file, toast info "Tidak ada data yang cocok dengan filter untuk diunduh." | ⬜ belum ada |
| E2E-06 | CDN diblokir (offline) | file `.csv` terunduh + toast peringatan | ⬜ belum ada |
| E2E-07 | Klik ganda saat mengunduh | tombol nonaktif, hanya 1 file | ⬜ belum ada |
| E2E-08 | Sesi berakhir | toast error, tidak ada file | ⬜ belum ada |

### 4.4 Manual / eksploratif

| ID | Skenario | Ekspektasi | Status |
|---|---|---|---|
| MN-01 | Buka file XLSX di Excel/LibreOffice | kolom rapi, teks & angka terbaca | ⬜ belum ada |
| MN-02 | Buka file CSV fallback di Excel | tidak ada karakter aneh (BOM) | ⬜ belum ada |
| MN-03 | Data >10.000 baris | file terpotong; perlu konfirmasi ke pengguna | ⬜ belum ada |
| MN-04 | Nama file pada dua unduhan berdekatan | tidak bentrok | ⬜ belum ada |

## 5. Hasil uji manual yang sudah dijalankan (15-11-2026)

Server uji: store `local`, data uji 7 record (1 dibuat oleh anggota `petugas`).

| ID | Hasil |
|---|---|
| E2E-01 | ✅ file `Data_Rekrutmen_Tangsel_20261115_094536.xlsx`, signature `PK` (XLSX sah), 21.035 byte, toast "7 data berhasil diunduh." |
| E2E-02 | ✅ `dimension A1:Q8` → 17 kolom × 8 baris (1 header + 7 data), header sesuai `exportHeaders` |
| E2E-03 | ✅ anggota unduh 2 baris (1 header + 1 record miliknya); total record 7 |
| E2E-04 | ✅ filter "Ciputat" → 2 baris (1 record); filter "Ekspor Uji 1" → 2 baris (1 record) |
| E2E-05 | ✅ tanpa file + toast "Tidak ada data yang cocok dengan filter untuk diunduh." |
| E2E-06 | ✅ CDN diblokir → file `.csv` (BOM `EFBBBF`) + toast "Pustaka Excel tidak termuat, file diunduh sebagai CSV." |
| E2E-08 | ✅ `GET /api/rekrutmen/export` tanpa sesi → 401 `{"error":"sesi tidak valid, silakan login kembali"}` |
| BE-01…BE-07 | ✅ `go test ./internal/api/ -run Export` → 5 test lulus |

## 6. Temuan (terdaftar, belum diperbaiki)

| ID | Temuan | Dampak | Saran |
|---|---|---|---|
| B-01 | `loadSheetJS()` tidak punya batas waktu; CDN yang lambat/menggantung membuat unduhan menggantung (fallback CSV hanya jalan saat `onerror`) | Tinggi | tambah timeout ±5 s, lalu jatuh ke CSV |
| B-02 | Promise SheetJS yang gagal disimpan selamanya (`sheetJSPromise`), jadi setelah 1 kegagalan semua unduhan berikutnya tetap CSV | Sedang | reset promise saat ditolak |
| B-03 | Pencarian punya debounce 300 ms; klik "Unduh Data" <300 ms setelah mengetik memakai nilai lama | Sedang | kirim `filters.q` langsung, bukan `searchTerm` |
| B-04 | Batas 10.000 baris memotong data tanpa pemberitahuan; `total` di respons tetap jumlah sebenarnya | Sedang | beri penanda "dipotong" atau beri tahu pengguna |
| B-05 | Endpoint export tidak memakai rate limiter (endpoint lain punya `registerLimiter`/`uploadLimiter`) | Sedang | tambah limiter / batasi permintaan export |
| B-06 | Tidak ada test otomatis frontend maupun E2E untuk fitur ini | Sedang | buat `tests/e2e/unduh-data.spec.ts` |
| B-07 | Nama file memakai waktu server, bukan waktu klien | Rendah | tidak perlu diubah |

## 7. Perintah verifikasi

```bash
cd backend && go build ./... && go vet ./... && go test ./... -run Export -v
cd frontend && npm run build
npx playwright test tests/e2e/unduh-data.spec.ts   # setelah spec dibuat
```

## 8. Urutan pengerjaan test yang disarankan

1. Backend: BE-08…BE-12 (lengkap, cepat, tanpa infra baru).
2. E2E: E2E-01…E2E-08 (butuh server + build frontend).
3. Manual: MN-01…MN-04.
4. Setelah B-01…B-05 diperbaiki, ulang seluruh matriks.
