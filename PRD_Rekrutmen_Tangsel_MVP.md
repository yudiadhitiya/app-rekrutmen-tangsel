# PRD --- Aplikasi Rekrutmen Tangsel

**Versi:** 0.1 (MVP)\
**Status:** Draft awal

## 1. Tujuan

Membangun aplikasi web untuk mencatat dan mengelola data rekrutmen di
wilayah Kota Tangerang Selatan. Fokus MVP adalah input, validasi,
penyimpanan, dan pengubahan data rekrutmen dengan UI modern,
profesional, dan responsif.

## 2. Pengguna

-   **Admin:** mengelola seluruh data rekrutmen dan konfigurasi sumber
    database.
-   **Petugas/Kader:** menambahkan dan mengubah data rekrutmen sesuai
    hak akses.

## 3. Fitur MVP

1.  Form tambah data rekrutmen.
2.  Daftar data rekrutmen dengan pencarian dan filter dasar.
3.  Ubah dan hapus data (hapus menggunakan konfirmasi).
4.  Validasi input wajib dan format data.
5.  Pilihan sumber penyimpanan: Supabase atau Google Sheets.
6.  Notifikasi berhasil/gagal saat menyimpan data.

## 4. Struktur Data dan Aturan Input

  -----------------------------------------------------------------------
  Field                   Tipe / Komponen         Ketentuan
  ----------------------- ----------------------- -----------------------
  Nama orang              Text                    Disarankan wajib;
                                                  validasi nama tidak
                                                  boleh kosong.

  No. WhatsApp            Text/angka              **Wajib**, hanya digit
                                                  0--9; simpan sebagai
                                                  string agar angka 0 di
                                                  awal tidak hilang.

  Alamat                  Text                    Opsional.

  RT                      Text/angka              Hanya digit 0--9.

  RW                      Text/angka              Hanya digit 0--9.

  Kecamatan               Picklist                **Wajib**; pilihan dari
                                                  7 kecamatan Kota
                                                  Tangerang Selatan:
                                                  Serpong, Serpong Utara,
                                                  Pondok Aren, Ciputat,
                                                  Ciputat Timur,
                                                  Pamulang, dan Setu.

  Kelurahan               Picklist bertingkat     Pilihan mengikuti
                                                  kecamatan yang dipilih;
                                                  reset pilihan jika
                                                  kecamatan berubah.

  Aktivitas               Picklist                **Wajib**; menggunakan
                                                  master kategori
                                                  pekerjaan/aktivitas
                                                  yang umum di Indonesia.
                                                  Master dapat dikelola
                                                  admin.

  Tingkat penerimaan      Picklist                **Wajib**; Sangat
                                                  Terbuka, Terbuka,
                                                  Netral, Kurang Terbuka,
                                                  Tertutup.

  Metode                  Picklist                **Wajib**; Halaqoh,
                                                  Program, Antara, Dakwah
                                                  Fardiyah.

  Status                  Picklist                **Wajib**; Berhasil,
                                                  Proses, Pending, Gagal.

  Keterangan              Text area               Opsional.

  UPLS                    Radio button            Opsional; hanya satu
                                                  pilihan. Opsi perlu
                                                  ditetapkan sebelum
                                                  implementasi.

  UPLM                    Radio button            Opsional; hanya satu
                                                  pilihan. Opsi perlu
                                                  ditetapkan sebelum
                                                  implementasi.
  -----------------------------------------------------------------------

**Catatan validasi:** Nama orang belum ditetapkan sebagai mandatory
dalam kebutuhan awal; pada MVP disarankan wajib. Nomor WhatsApp
divalidasi sebagai digit saja dan disimpan dalam format string. Validasi
panjang/format nomor dapat ditambahkan setelah aturan operasional
disepakati.

## 5. Alur Utama

1.  Pengguna membuka menu **Tambah Rekrutmen**.
2.  Pengguna mengisi data dan memilih kecamatan.
3.  Sistem memuat pilihan kelurahan sesuai kecamatan.
4.  Sistem memvalidasi field wajib dan format nomor.
5.  Pengguna menekan **Simpan**.
6.  Sistem menyimpan ke database aktif dan menampilkan notifikasi hasil.
7.  Data baru muncul pada daftar rekrutmen.

## 6. UI/UX

-   Tampilan modern, profesional, bersih, dan responsif untuk desktop
    maupun tablet/mobile.
-   Form dikelompokkan: Identitas, Wilayah, Informasi Rekrutmen, dan
    Keterangan.
-   Tampilkan penanda field wajib, pesan validasi di dekat field, serta
    loading state saat menyimpan.
-   Daftar data menggunakan tabel pada layar lebar dan tampilan kartu
    pada layar kecil.
-   Sediakan pencarian berdasarkan nama/nomor WhatsApp serta filter
    kecamatan, kelurahan, metode, dan status.

## 7. Arsitektur dan Tech Stack

-   **Frontend:** Vue 3.
-   **Backend/API:** Golang (REST API).
-   **Database utama:** Supabase (PostgreSQL).
-   **Alternatif penyimpanan:** Google Sheets melalui Google Sheets API.
-   **Pola integrasi:** gunakan interface/repository penyimpanan di
    backend agar aplikasi dapat memilih provider `supabase` atau
    `google_sheets` melalui konfigurasi.

Pemilihan provider dilakukan oleh Admin melalui pengaturan aplikasi.
Kredensial Google disimpan aman di backend/environment secret, bukan di
frontend. Perubahan provider tidak otomatis memindahkan data; sediakan
prosedur ekspor/impor dan verifikasi sebelum perpindahan digunakan.

## 8. Data Master Wilayah dan Aktivitas

-   Master kecamatan berisi 7 kecamatan Kota Tangerang Selatan.
-   Master kelurahan berelasi ke satu kecamatan dan menjadi sumber
    picklist dinamis.
-   Master aktivitas/pekerjaan dikelola terpusat agar pilihan konsisten
    dan dapat diperbarui.
-   Hindari mengetik bebas untuk field yang telah ditetapkan sebagai
    picklist.

## 9. Keamanan dan Kualitas Minimum

-   Autentikasi pengguna dan pembatasan akses berdasarkan peran.
-   Validasi dilakukan di frontend dan backend.
-   Gunakan HTTPS; rahasia/kunci API hanya berada di backend.
-   Terapkan Row Level Security pada Supabase sesuai kebutuhan akses.
-   Catat waktu pembuatan dan perubahan data (`created_at`,
    `updated_at`) serta pengguna yang melakukan perubahan.
-   Cegah duplikasi nomor WhatsApp bila aturan bisnis mengharuskannya;
    default MVP: tampilkan peringatan sebelum menyimpan nomor yang sudah
    ada.

## 10. Kriteria Penerimaan MVP

-   Pengguna dapat menambah, melihat, mengubah, dan menghapus data.
-   Nomor WhatsApp wajib dan menolak karakter selain digit.
-   RT/RW menolak karakter selain digit.
-   Kecamatan wajib dan hanya menampilkan 7 kecamatan Tangsel.
-   Kelurahan berubah mengikuti kecamatan yang dipilih.
-   Aktivitas, tingkat penerimaan, metode, dan status wajib dipilih dari
    opsi yang ditentukan.
-   UPLS dan UPLM bersifat opsional dan masing-masing hanya menerima
    satu pilihan.
-   Data tersimpan pada provider yang sedang aktif (Supabase atau Google
    Sheets).
-   Form dan daftar data dapat digunakan pada desktop dan perangkat
    mobile.
-   Nama Orang wajib diisi.
-   Radio button UPLS dan UPLM maksudnya user secara opsional bisa memilih salah satu dari UPLS atau UPLM.
-   Nomor Whatsapp bersifat unik untuk setiap orang.
-   Terkait Google Sheet, Admin bisa mengganti provider dari menu pengaturan
-   Terdapat master data Daftar Aktifitas/Pekerjaan yang dapat dikelola Admin,

## 11. Di Luar Scope MVP

-   Dashboard analitik dan laporan lanjutan.
-   Impor massal dari Excel/CSV.
-   Integrasi WhatsApp untuk mengirim pesan.
-   Workflow persetujuan berjenjang.
-   Sinkronisasi dua arah otomatis antara Supabase dan Google Sheets.

------------------------------------------------------------------------

**Prinsip implementasi:** bangun MVP terlebih dahulu: form input +
validasi + daftar data + edit, dengan Supabase sebagai provider awal.
Tambahkan Google Sheets sebagai provider alternatif setelah alur CRUD
stabil.
