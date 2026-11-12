# Audit UI/UX — Rekrutmen Tangsel (Tahap 1: Analisis, tanpa perubahan kode)

**Status:** dokumentasi saja. Tidak ada file sumber yang diubah. Seluruh business logic, API, dan kontrak data tidak tersentuh.

**Metode:**
1. Impeccable detector (deterministik, atas `frontend/src`).
2. Review desain independen oleh design director terpisah (sub-agent), terisolasi dari hasil detector.
3. Pemeriksaan manual: rasio kontras WCAG terhitung, ukuran target sentuh, struktur semantik/ARIA, audit token warna, konsistensi tipografi & spacing.

**Mode produk:** Operate (tool admin/internal, bukan landing page). Pengguna utama: petugas lapangan yang menginput data dari HP.

---

## 1. Ringkasan Eksekutif

Aplikasi sudah punya fondasi yang rapi: satu stylesheet dengan design token, komponen yang konsisten, loading/empty/error state yang lengkap, dan pola responsif tabel→kartu. Yang belum ada adalah **kecukupan aksesibilitas** (kontras, dialog, target sentuh), **sistem tipografi** (belum ada skala; hierarki lemah), dan **karakter produk** (UI masih templat admin generik tanpa identitas Kota Tangerang Selatan). Belum ada `DESIGN.md`/`PRODUCT.md`, jadi keputusan visual tidak pernah didokumentasikan.

**Skor Audit: 11/20 (Acceptable — perlu perbaikan signifikan, fondasi kuat).**
**Skor Heuristik: 24/40 (60% — Acceptable).**

---

## 2. Skor Audit Teknis

| # | Dimensi | Skor | Temuan Kunci |
|---|---------|------|--------------|
| 1 | Accessibility | 1/4 | 4 token warna di bawah 4.5:1; dialog tanpa focus trap/Esc; error form tanpa `aria-invalid`/`aria-describedby`; target sentuh 27–33px |
| 2 | Performance | 3/4 | Bundle kecil, tanpa dependensi UI, animasi ringan; tidak ada `prefers-reduced-motion` |
| 3 | Theming | 2/4 | Token dipakai luas, tetapi hex bocor di stylesheet, `style="..."` inline di 2 view, tanpa dark mode |
| 4 | Responsive | 2/4 | Breakpoint 900/820/480 + tabel→kartu; tabel 10 kolom padat, target sentuh kecil, nav scroll horizontal |
| 5 | Implementation Integrity | 3/4 | Sistem koheren; drift: markup tabel & kartu duplikatif, styling inline, seksi form bernomor dekoratif |
| **Total** | | **11/20** | **Acceptable** |

---

## 3. Skor Heuristik Nielsen

| # | Heuristik | Skor | Catatan |
|---|-----------|------|---------|
| 1 | Visibility of System Status | 3 | Loading, toast, badge status ada; state "filter kosong" tidak dibedakan dari "belum ada data" |
| 2 | Match System / Real World | 3 | Bahasa domain tepat; istilah "UPL", "Tingkat Penerimaan", "Metode" tanpa penjelasan |
| 3 | User Control and Freedom | 3 | Batal, Reset filter, konfirmasi hapus; tidak ada undo, dialog tak bisa ditutup Esc |
| 4 | Consistency and Standards | 3 | Token & kelas konsisten; bocor hex dan style inline |
| 5 | Error Prevention | 3 | Wajib diisi, validasi inline, deteksi duplikat WA sebelum simpan |
| 6 | Recognition Rather Than Recall | 2 | 10 kolom tabel, "UPL" tak dikenal petugas baru |
| 7 | Flexibility & Efficiency | 2 | Debounce pencarian, pagination; tanpa bulk action, ekspor, sortir, pintasan keyboard |
| 8 | Aesthetic & Minimalist | 2 | Tajam secara visual; tabel padat + penomoran seksi/formula "1-2-3-4" |
| 9 | Error Recovery | 2 | Pesan error jelas di dekat field; fokus tidak diarahkan ke field pertama yang gagal, tidak diumumkan ke screen reader |
| 10 | Help & Documentation | 1 | Hanya hint satu baris; tidak ada halaman bantuan, tooltip, atau penjelasan istilah |
| **Total** | | **24/40** | **60% — Acceptable** |

---

## 4. Analisis UI Existing (per Permukaan)

### 4.1 App Shell & Navigasi (`App.vue`)
- App bar sticky, brand "RT" (teks inisial, bukan logo), nav 4 item, blok user + tombol Keluar. Bersih dan rapi.
- **Issue:** "Tambah" muncul sebagai item nav *dan* sebagai CTA utama di `ListView` — dua jalan ke aksi yang sama, dengan hierarki yang saling meniadakan. Item nav "Tambah" juga tak menjelaskan konteks (tambah rekrutmen? aktivitas?).
- Di lebar < 820px, nav menjadi baris horizontal scroll (`overflow-x: auto`, scrollbar disembunyikan) — item di luar layar tidak terdeteksi pengguna, dan target sentuh link nav hanya ±31px.

### 4.2 Daftar Rekrutmen (`ListView.vue`)
- Struktur: PageHead + Card filter + Card tabel + pagination. Estado loading/empty/error lengkap. Card mobile memuat alamat (data tidak ada di tabel desktop) — asimetri isi tabel vs kartu.
- **Issue:** markup tabel (baris 259–309) dan kartu mobile (312–359) merender data yang sama dua kali → setiap perubahan kolom harus disinkronkan di dua tempat (risiko integrity).
- **Issue:** 10 kolom tabel (Nama, WA, Kecamatan, Kelurahan, Aktivitas, Tingkat, Metode, Status, UPL, Aksi) — kepadatan tinggi, tidak ada sortable column, tidak ada bulk action, tidak ada ekspor (padahal ada rencana integrasi Sheets).
- **Issue:** empty state "Tidak ada data" dipakai untuk dua kondisi berbeda: database kosong vs filter tidak cocok. Tidak ada aksi "Reset filter" dari dalam empty state.
- Tidak ada state "memuat" per-aksi (klik Hapus hanya mengganti teks tombol menjadi "...").

### 4.3 Form Tambah/Ubah (`FormView.vue`)
- Empat seksi bernomor (Identitas, Wilayah, Informasi Rekrutmen, Keterangan), grup `.form-row` 2 kolom, radio pill untuk UPL, validasi inline + pesan server.
- **Issue:** penomoran seksi "1/2/3/4" dengan badge lingkaran + heading uppercase adalah pola dekoratif generik, bukan pemandu yang menambah comprehension.
- **Issue:** field Keterangan (textarea panjang) berada di seksi terakhir yang sama dengan UPL — textarea panjang di dalam form satu kolom nichts awkward secara komposisi.
- **Issue:** label error tidak terhubung ke input (`aria-describedby`, `aria-invalid` tidak ada), sehinggaValidate on submit → toast "Periksa kembali isian yang ditandai merah" tapi fokus tidak pindah ke field yang gagal.
- **Issue:** No. WhatsApp `inputmode="numeric"` tapi `type="text"` — beberapa keyboard mobile tidak memunculkan keypad angka; belum ada normalisasi prefiks 62.

### 4.4 Login & Registrasi (`LoginView.vue`)
- Satu kartu 420px; mode login/register memakai satu kartu yang sama.
- **Issue:** formulir registrasi 8 field (username, sandi, nama, kecamatan, 2 pertanyaan + 2 jawaban) ditumpuk vertikal dalam satu kartu sempit → scroll panjang di HP, pertanyaan/jawaban yang berpasangan tidak dikelompokkan secara visual.
- Tidak ada toggle "lihat kata sandi", tidak adarequirement strength, error message global tidak fokus ke field.
- Login dan registrasi berbagi state `username`/`password` yang sama —стоя ошибка/kehilangan input saat berpindah mode.

### 4.5 Master Aktivitas (`AktivitasView.vue`)
- Tambah inline + edit inline (input + checkbox "Aktif") + badge status. Rapi.
- **Issue:** admin melihat checkbox "Aktif" saat edit, non-admin melihat badge "Aktif/Nonaktif" — dua representasi status berbeda untuk data yang sama.
- **Issue:** style inline `flex: 1 1 auto` melanggar sistem token.

### 4.6 Pengaturan (`SettingsView.vue`)
- Pilihan provider sebagai kartu radio kustom (`.provider-option`) — pola yang benar.
- **Issue:** "Informasi Provider" merender seluruh objek `settings` sebagai pasangan key/value mentah (`JSON.stringify` untuk objek) → kebocoran detail internal server ke UI, noisy, dan label teknis ("sheets_id", "provider") tanpa penjelasan.
- **Issue:** style inline `margin-top: 16px` di dua tempat.

### 4.7 Toast & Dialog Konfirmasi
- Toast punya `role="status" aria-live="polite"`, judul per tipe, border kiri 4px berwarna.
- **Issue (detector, 1 temuan):** `style.css:940` — `border-left: 4px solid` pada toast = pola "side-tab accent border", penanda UI generik. Info dikomunikasikan lewat warna saja.
- **Issue:** tombol tutup toast ±17px (di bawah WCAG 2.5.8 minimum 24px).
- **Issue:** `ConfirmDialog` punya `role="dialog" aria-modal="true"` tapi **tanpa focus trap, tanpa pemindahan fokus ke dalam dialog, tanpa dukungan Esc, dan tanpa `aria-describedby`** ke `.modal__message`. Pengguna keyboard dapat mengetik di belakang dialog; pengguna screen reader tidak tahu dialog muncul.
- **Issue:** `aria-live="polite"` dipakai untuk toast error; error sebaiknya `assertive`.

### 4.8 Bahasa Visual (yang sudah bekerja)
- Token warna lengkap (primary/danger/success/warning/info + soft variants), radius bertingkat, elevation 3 tingkat, spacing bertingkat, satu keluarga font.
- `.data-table` punya header uppercase kecil + hover row + strong column — tabular scanning cukup baik.
- `.radio-option` pill untuk pilihan UPL — lebih baik daripada select untuk opsi pendek.

---

## 5. Evaluasi UX & Konsistensi Visual

| Aspek | Temuan | Dampak |
|-------|--------|--------|
| **Kontras teks** | `--muted #7b879b` = **3.63:1** di putih (butuh 4.5:1) dipakai di subtitle halaman, hint field, `result-count`, `pagination__info`, header tabel (11.5px), label kartu mobile (10.5px) | Teks sekunder nyaris tak terbaca — kontras gagal |
| **Kontras badge** | tone-success 4.15:1, tone-warning 3.77:1, tone-danger 4.00:1 di atas background soft-nya (font 12px, bukan large text) | Badge status — elemen paling sering dipindai — di bawah ambang |
| **Ketidakkonsistenan warna status** | Versi gelap dipakai di `.alert--*` (lulus 5.8–7.1:1) tetapi warna terang dipakai di `.badge/.tone-*` (gagal). Satu semantik, dua token | Developer tidak tahu token mana yang benar |
| **Placeholder** | `#a6b0c0` = 2.19:1 | Placeholder jadi satu-satunya sumber label saat form terisi |
| **Tipografi** | Tidak ada skala: 10.5, 11.5, 12, 12.5, 13, 13.5, 14, 14.5, 15, 16, 17, 21, 22px (6 angka pecahan). H1 22px vs body 15px (rasio 1.47). Judul card 15px = ukuran body | Hierarki lemah; replikasi antar layar akan mudah pecah |
| **Spacing** | Skala 4/8/12/16/24/32 ada, tapi `padding: 7px`, `6px 10px`, `9px 11px`, `11px 14px`, `12px 16px`, `22px`, `26px`, `30px` dipakai ad-hoc | Ritme visual tidak seragam antar komponen |
| **Radius** | 6/10/16 dipakai, tapi badge 999px, pill 999px, num 50% |cukup konsisten (kebab-case classes) |
| **Elevation** | 3 tingkat konsisten | Bagus |
| **Token bocor** | Hard-coded hex: `#fbdada`, `#a6b0c0`, `#f3c9c9`, `#f0dcb5`, `#c6dff2`, `#bfe6d6`, `#c7d0de`, `#6b8bff`, `rgba(...)` ×5 | Mode gelap & token semantik mustahil diturunkan |
| **Style inline** | `SettingsView.vue:88,105`, `AktivitasView.vue:162` | Melewati design system |
| **Duplikasi struktur** | Tabel & kartu mobile render ulang baris yang sama | Dua tempat harus selalu sinkron |
| **Dekorasi generik** | Numbered section 1–4, uppercase small labels, toast side-tab border, gradient logo, card border-left | Kombinasi ciri UI generik (terdeteksi di detector dan pada review manual) |

---

## 6. Lima Prioritas Perbaikan

### Prioritas 1 — P1 · Kontras & keterbacaan status (A11y)
- **Lokasi:** `style.css:14-31` (token), dipakai di `.muted`, `.field__hint`, `.data-table th`, `.recruiter-card__meta-label`, `.tone-*`.
- **Dampak:** elemen sekunder dan badge status — bagian yang paling sering dipindai petugas — berada di bawah ambang WCAG AA (4.5:1). Di layar HP outdoor, praktis tak terbaca.
- **Standar:** WCAG 2.1 SC 1.4.3 (Text Contrast), SC 1.4.11 (Non-text Contrast).
- **Perbaikan:** gelapkan `--muted` ke `~#5c6879` (≥4.6:1 di putih), dan setiap status punya dua token: `--status-*-fg` (teks) + `--status-*-bg` (latar). Badge & alert memakai pasangan yang sama agar konsisten.

### Prioritas 2 — P1 · Dialog konfirmasi, fokus, & umpan balik error belum aksesibel
- **Lokasi:** `components/ConfirmDialog.vue:7-24`, `views/FormView.vue` (error spans), `views/LoginView.vue` (error global), `style.css` (tidak ada `:focus-visible`).
- **Dampak:** pengguna keyboard dapat mengetik di belakang dialog; pengguna screen reader tidak diberi tahu dialog muncul maupun pesan error; fokus tidak pernah diarahkan ke field pertama yang gagal. Tidak ada satu pun aturan `:focus-visible` di stylesheet → indikator fokus hanya default browser yang hilang saat klik mouse, dan kontras fokus tidak dijamin.
- **Standar:** WCAG 2.1 SC 2.1.1 Keyboard, SC 2.4.7 Focus Visible, SC 2.4.3 Focus Order, SC 3.3.1 Error Identification, SC 4.1.3 Status Messages.
- **Perbaikan:** focus trap + fokus awal ke tombol konfirmasi + Esc untuk batal + `aria-describedby` ke pesan; tambahkan `aria-invalid="true"` + `aria-describedby` per field error; fokuskan field pertama yang gagal saat submit; `role="alert"` pada blok error; satu aturan fokus global (`:focus-visible { outline: 2px solid var(--color-focus); outline-offset: 2px }`).

### Prioritas 3 — P1 · Target sentuh & ergonomis input HP
- **Lokasi:** `.btn--sm` (tinggi ±27,6px), `.app-nav__link` (±30,8px), `.radio-option` (±32,8px), `.toast__close` (±17px), `.badge` (±20px).
- **Dampak:** pengguna utama (petugas lapangan) memakai HP; target di bawah 24px melanggar SC 2.5.8, jauh di bawah pedoman 44px, menyebabkan salah sentuh pada "Ubah/Hapus".
- **Standar:** WCAG 2.2 SC 2.5.8 Target Size (Minimum), 2.5.5 (AAA 44px).
- **Perbaikan:** tinggi minimum 44px untuk kontrol sentuh di breakpoint mobile (bisa tetap 32–36px di desktop), `min-width` 44px untuk tombol tutup toast, area sentuh radio/checkbox diperluas lewat pseudo-element, dan navigasi mobile diubah ke pola bottom-nav atau menu, bukan scroll horizontal tersembunyi.

### Prioritas 4 — P2 · Sistem tipografi & token yang belum tertutup
- **Lokasi:** `style.css` (seluruh skala font), `views/*.vue` (inline style).
- **Dampak:** 6 ukuran font pecahan tanpa peran semantik → hierarki lemah (h1 hanya 1.47× body, judul card = ukuran body) dan setiap layar baru berisiko berbeda sendiri. Font "Inter" hanya dideklarasikan di `--font` (`style.css:42`) tetapi tidak pernah dimuat (tanpa `@font-face`/Google Fonts) → semua perangkat jatuh ke system font, konsisten di satu perangkat, tidak di perangkat lain. Hex bocor + `style="..."` inline membuat dark mode & tema baru mustahil.
- **Perbaikan:** muat Inter secara eksplisit (`font-display: swap`) atau pilih satu system font stack; tetapkan type scale berbasis peran (`--text-xs/sm/body/lg/xl/2xl` + `--leading-*`), pindahkan semua hex ke token semantik (fg/bg/border per status & surface), hapus style inline, dan tulis `DESIGN.md` sebagai sumber kebenaran visual + daftarkan komponen (Button/Field/Card/Badge/Table/EmptyState/Alert/Toast/Dialog/Pagination).

### Prioritas 5 — P2 · Karakter produk & kejelasan alur
- **Lokasi:** `App.vue:14-22` (nav), `views/ListView.vue:166-390`, `views/SettingsView.vue:123-133`.
- **Dampak:** UI bisa dipakai aplikasi apa pun — belum ada pembeda "Rekrutmen Tangsel" (identitas logo/instansi, meta wilayah, hierarki borough/kecamatan, tone formalism). Ditambah empty state menyatukan dua kondisi, tabel 10 kolom tanpa bulk/ekspor/sort, dan halaman Pengaturan yang menampilkan key teknis mentah.
- **Perbaikan:** (a) brand mark + nama instansi yang konsisten, (b) pisahkan state "filter kosong" dari "belum ada data" dengan aksi "Reset filter" atau "Tambah data pertama", (c) rapikan tabel menjadi 6–7 kolom inti + panel detail/side sheet, (d) sediakan ekspor (CSV/Sheets) dan bulk action, (e) tampilkan informasi provider dalam bentuk terstruktur (nama, status koneksi, terakhir sinkron) — bukan dump JSON.

---

## 7. Rekomendasi Design System

### 7.1 Prinsip
Pertahankan identitas visual yang ada (biru primer, card putih di atas bg kebiruan, tipografi Inter) — fondasinya sudah layak. Yang diperbaiki adalah **perannya**: token semantik, skala tipografi, aksesibilitas, dan satu motif visual yang benar-benar milik domain (mis. peta wilayah/borough atau kartu data bergaya wilayah, bukan dekorasi generik).

### 7.2 Lapisan Token
```
Layer 1 — Primitif (nilai mentah, tanpa arti):  --blue-600, --gray-100, --red-500 ...
Layer 2 — Semantik (arti, gantiable):         --color-bg, --color-surface,
                                                --color-text / -soft / -muted,
                                                --color-border, --color-focus,
                                                --status-{success|warning|danger|info}-{fg,bg,border}
Layer 3 — Komponen:                            --btn-height-md/sm/touch, --field-height,
                                                --table-row-height, --card-radius, --card-padding
Layer 4 — Layout:                              --space-1..8 (skala 4px), --container-max,
                                                --z-{base,sticky,modal,toast}
```
Aturan: tidak ada hex di luar Layer 1; tidak ada angka px hard-coded di dalam komponen.

### 7.3 Type Scale (peran, bukan ukuran)
| Peran | Ukuran / line-height | Pemakaian |
|-------|----------------------|-----------|
| Display | 28 / 1.2 | Judul layar utama (jarang) |
| Title | 22 / 1.3 | `page-head__title` |
| Heading | 18 / 1.35 | `card__title`, judul seksi |
| Subheading | 16 / 1.4 | Subjudul, label grup |
| Body | 15 / 1.55 | Teks utama, isi input |
| Small | 13 / 1.5 | Label field, badge, tombol |
| Micro | 11.5 / 1.4 uppercase | Header tabel, meta label |

Rasio antar tingkat ≥1.2; judul layar minimal 1.6× body (perbaikan dari 1.47× saat ini).

### 7.4 Komponen Inti
- **Button**: `variant` (primary/secondary/ghost/danger), `size` (sm/md/touch), `icon`, `loading`, `block`. Tinggi: 36px desktop / 44px touch. Fokus: ring 2px `var(--color-focus)` + offset 2px pada semua elemen interaktif.
- **Field**: label + hint + error + counter, satu token tinggi (40px / 48px touch), `aria-describedby` otomatis menautkan hint & error.
- **RadioGroup**: pill untuk ≤3 opsi (UPL), card-select untuk opsi berdeskripsi (provider), bukan `<select>` panjang.
- **StatusBadge**: pasangan fg/bg semantik; ikon + teks (bukan warna saja); selalu menyimpan makna saat di-cetak/reader.
- **DataTable**: sticky header, zebra halus, kolom sortable, bulk checkbox, empty & skeleton state, tampilan kartu pada mobile yang **menggunakan data & aksi yang sama** (satu sumber render).
- **EmptyState**: dua varian — `kosong` (belum ada data) dengan CTA utama, dan `tersaring` (filter tidak cocok) dengan CTA "Reset filter".
- **Dialog**: focus trap, Esc, `aria-describedby`, fokus kembali ke pemicu saat ditutup.
- **Motion**: token durasi 120/180/240ms + easing; blokir via `prefers-reduced-motion`.

### 7.5 Dokumentasi & Proses
- `DESIGN.md` — token, skala, daftar komponen, aturan pakai (dipakai skill `impeccable document`).
- `PRODUCT.md` — pengguna, tugas utama, batasan (dipakai skill `impeccable init`).
- Aturan PR: tidak boleh menambah hex baru di luar Layer 1; komponen baru harus memakai token & gaya fokus standar.

### 7.6 Peta Jalan (Tahap 2+)
1. Init context (`PRODUCT.md`) & dokumentasi (`DESIGN.md`) — tanpa mengubah kode.
2. **Tiga arah visual** (gate Huashu-Design: 3 konsep × 2 kandidat, pilih satu) untuk karakter produk; tidak masuk implementasi sebelum arah disetujui.
3. Implementasi P1 (kontras, a11y dialog/form, target sentuh) — mandiri dari arah visual baru.
4. Implementasi P2 (type scale, token cleanup, alur & status layar).
5. Uji ulang `/impeccable audit` + `/impeccable critique` untuk mengukur perbaikan skor.

---

## 8. Catatan Deteksi Otomatis

| Temuan | File | Severity | Kategori |
|--------|------|----------|----------|
| Side-tab accent border (`border-left: 4px solid` pada toast) | `frontend/src/style.css:940` | warning | slop |

Temuan manual yang ditemukan dan tidak dilaporkan detector: kontras teks/badge di bawah WCAG AA, dialog tanpa focus trap, `aria-invalid`/`aria-describedby` yang hilang, `focus-visible` yang tidak ada, target sentuh < 24px, `prefers-reduced-motion` yang tidak ada, dan kebocoran hex/token.

### 8.1 Hasil cross-check review independen
Review kedua (design director terisolasi, tanpa melihat hasil detector) menguatkan temuan yang sama: kontras, target sentuh, dialog tanpa focus trap, hierarki tipografi, duplikasi CTA "Tambah", token bocor, dan `style` inline. Tambahan yang hanya muncul dari review kedua:

- **P1** — Tidak ada satupun aturan `:focus-visible` di seluruh stylesheet → indikator fokus tidak konsisten antar browser dan hilang pada interaksi mouse.
- **P2** — Font "Inter" dideklarasikan tapi tidak pernah dimuat → perangkat jatuh ke system font.
- **P2** — `ConfirmDialog` hanya bisa ditutup lewat klik backdrop atau tombol; tidak ada handler Esc.
- **P2** — Tombol "Batal" pada form (`FormView.vue:186-188`) membuang seluruh isian tanpa konfirmasi.
- **P3** — Nav "Tambah" tidak pernah aktif pada rute `/ubah/:id`, karena rute itu tidak ada di daftar nav.
- **P3** — Tombol "Hapus" pada tabel hanya menampilkan "..." saat memproses — tidak ada indikator memuat yang konsisten dengan pola spinner di aplikasi.
- **P3** — Formulir registrasi menampilkan 8 field dalam satu kartu 420px; pertanyaan/jawaban keamanan tidak dikelompokkan berpasangan.

### 8.2 Temuan domain (di luar lapisan visual, tapi berdampak ke UX)
- **Seed wilayah tidak lengkap:** `backend/internal/seed/wilayah.json` memuat **7 kecamatan / 54 kelurahan**. Kota Tangerang Selatan memiliki 13 kecamatan / 104 kelurahan (data ini perlu konfirmasi ke sumber resmi). Akibatnya dropdown Kecamatan/Kelurahan di form dan filter **tidak bisa dipilih petugas di wilayah yang tidak tercantum** — masalah data, bukan desain, tetapi tetap berisiko besar untuk adopsi. (Perlu verifikasi ulang.)
- **Data demo ikut tampil saat memakai store lokal:** `data/rekrutmen.json` berisi 20 record contoh (nama "Fulan", "Fulanah", "Agen Asuransi", dst.). Record contoh harus diganti dengan data riil sebelum rilis, atau diberi penanda jelas "Data contoh".

---

## 9. Catatan tentang integrasi skill

- **Impeccable** `context` tidak menemukan `PRODUCT.md`/`DESIGN.md`/surface brief, dan hook detector belum terpasang di `.claude/settings.local.json`.
- Perintah `impeccable audit-storage` **tidak tersedia** pada versi binary yang terpasang, jadi laporan tidak disimpan ke `.impeccable/`; laporan penuh tersimpan di `UI_UX_AUDIT.md` (root repo).
- **Huashu-Design**: gate "tiga arah" belum dijalankan karena tahap ini murni audit (tidak ada output visual baru). Gate tersebut wajib dijalankan sebelum Tahap 2 (perubahan visual) — lihat Peta Jalan 7.6.

---

## 9. Pertanyaan untuk Stakeholder

1. Prioritas mana yang lebih berat: kontras & aksesibilitas (P1 a11y), atau ergonomis & kejelasan alur (P1 sentuh + P2 alur)?
2. Apakah identitas visual "Rekrutmen Tangsel" perlu kuat (logo instansi, warna borough/motif wilayah), atau cukup utilitas internal yang netral?
3. Cakupan tahap 2: semua 5 prioritas, atau hanya P1 dulu sebelum ada perubahan visual besar?
