# Deployment Plan — Rekrutmen Tangsel

**Status dokumen:** aktif, diperbarui per checkpoint
**Terakhir diperbarui:** 2026-11-06
**Commit baseline:** `3d7d05a`
**Rujukan:** [`PRD_Rekrutmen_Tangsel_MVP.md`](PRD_Rekrutmen_Tangsel_MVP.md) v0.1,
[`API_CONTRACT.md`](API_CONTRACT.md),
[`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md),
[`TESTING_PLAN.md`](TESTING_PLAN.md)

---

## 1. Cara memakai dokumen ini

Penanda status sama seperti dua dokumen lain:

| Penanda | Arti |
|---|---|
| ✅ | **Selesai** — sudah ada & terverifikasi |
| 🔧 | **Sedang dikerjakan** — sudah ada tapi belum diverifikasi/ter-commit |
| ⬜ | **Belum** — belum ada |
| ⏸️ | **Ditunda** — sengaja ditunda, dengan alasan |

Nomor item `P-xx` dipakai di dokumen ini. Item di dokumen lain tetap memakai
nomornya masing-masing (`D-xx`, `T-xx`).

---

## 2. Ringkasan

| Fase | Selesai | Belum | Total | Keterangan |
|---|---|---|---|---|
| P1 — Environ & konfigurasi | 6 | 2 | 8 | env vars & `.env.example` |
| P2 — Build & bundling | 4 | 3 | 7 | build Go + Vite |
| P3 — Server & reverse proxy | 2 | 6 | 8 | systemd + nginx |
| P4 — HTTPS & header keamanan | 0 | 5 | 5 | **blocker produksi** |
| P5 — Penyimpanan & data | 4 | 6 | 10 | backup/restore & migrasi |
| P6 — CI/CD | 0 | 4 | 4 | **belum ada sama sekali** |
| P7 — Monitoring & operasional | 2 | 7 | 9 | health check, log, backup otomatis |
| P8 — Kontainer & orkestrasi | 0 | 3 | 3 | opsional, non-blocking |
| **Total** | **18** | **36** | **54** | |

> P5 memuat 1 item berstatus 🔧 (P-33, migrasi skema — idempoten tapi belum
> berversi). Item itu dihitung sebagai "belum" pada kolom di atas.

### Status produksi

| Aspek | Status |
|---|---|
| Build & Artifact | ✅ bisa dibangun (Go + Vite) |
| Jalankan lokal | ✅ server boot, API & frontend terverifikasi |
| Jalankan di server | 🔧 bisa — butuh systemd/nginx |
| HTTPS | ⬜ belum ada |
| CI/CD | ⬜ belum ada |
| Backup otomatis | ⬜ belum ada |
| **Siap produksi?** | ❌ **belum** — P4 dan P6 adalah blocker |

---

## 3. Arsitektur target

```
                    Internet
                       │
                 ┌─────▼─────┐
                 │  nginx    │  TLS termination, HTTP/2,
                 │  :443     │  security headers, static files
                 └─────┬─────┘
                       │ proxy_pass → 127.0.0.1:8080
                 ┌─────▼──────────┐
                 │  rekrutmen     │  binary Go statis
                 │  systemd unit  │  127.0.0.1:8080 (tidak publik)
                 └─────┬──────────┘
                       │
        ┌──────────────┼───────────────┐
        │              │               │
   ┌────▼────┐   ┌─────▼─────┐   ┌────▼─────┐
   │  data/  │   │ Supabase  │   │  Google  │
   │ (JSON)  │   │PostgreSQL │   │  Sheets  │
   └─────────┘   └───────────┘   └──────────┘
```

Prinsip:
- Backend hanya mendengarkan di loopback; satu-satunya pintu masuk adalah nginx.
- Frontend disajikan sebagai file statis oleh nginx (bukan oleh Go), supaya Go
  fokus pada API. Backend tetap punya fallback penyaji statis untuk mode
  sederhana/dev.
- Rahasia hanya di environment systemd, tidak pernah di repo.

---

## 4. P1 — Environment & konfigurasi

| # | Item | Status | Detail |
|---|---|---|---|
| P-01 | Backend membaca konfigurasi dari environment | ✅ | `cmd/server/main.go:16-21` |
| P-02 | Variabel environment terdokumentasi | ✅ | README tabel env (9 variabel) |
| P-03 | Fallback aman bila `SESSION_SECRET` kosong | ✅ | generate + simpan `data/session.secret` (mode 0600) |
| P-04 | Data runtime terisolasi dari repo | ✅ | `.gitignore` blokir `data/`, `*.secret`, `*.env` |
| P-05 | Secret Google Sheets dari file di luar repo | ✅ | `GOOGLE_KEY_FILE` |
| P-06 | Kredensial Supabase dari environment | ✅ | `SUPABASE_URL` + `SUPABASE_API_KEY` (service role) |
| P-07 | **`.env.example`** | ⬜ | tidak ada; perlu sebagai template onboarding |
| P-08 | **Dokumentasi rotasi secret** | ⬜ | kapan & bagaimana mengganti `SESSION_SECRET` |

### Variabel environment

| Variabel | Default | Wajib di produksi | Keterangan |
|---|---|---|---|
| `ADDR` | `:8080` | ya | **set `127.0.0.1:8080`** di produksi |
| `DATA_DIR` | `./data` | ya | direktori data + backup |
| `FRONTEND_DIR` | `frontend/dist` | ya | path absolut ke hasil build |
| `SESSION_SECRET` | generate | **wajib** | hex panjang; sesi jadi invalid bila berubah |
| `SUPABASE_URL` | — | bila pakai Supabase | URL project |
| `SUPABASE_API_KEY` | — | bila pakai Supabase | service-role key (rahasia) |
| `GOOGLE_SHEET_ID` | — | bila pakai Sheets | ID spreadsheet |
| `GOOGLE_KEY_FILE` | — | bila pakai Sheets | path file JSON service account |
| `APP_ROOT` | `.` | ya | root aplikasi untuk path relatif |

---

## 5. P2 — Build & bundling

| # | Item | Status | Verifikasi |
|---|---|---|---|
| P-09 | Build backend → binary | ✅ | `go build -o rekrutmen ./cmd/server` |
| P-10 | Build frontend → bundle statis | ✅ | `npm run build` → `frontend/dist` |
| P-11 | Frontend served oleh backend (fallback) | ✅ | `main.go:76-80` |
| P-12 | Vite proxy ke backend saat dev | ✅ | `vite.config.js` |
| P-13 | **Build single-binary (embed frontend)** | ⬜ | opsional; embed `frontend/dist` ke dalam binary Go |
| P-14 | **Artifact ter-versioning (release + checksum)** | ⬜ | untuk rollback & audit |
| P-15 | **Build reproducible & vendor** | ⬜ | `go mod vendor`, lockfile Node |

### Build yang dipakai saat ini

```bash
# Backend — butuh Go 1.24
cd backend
export PATH=$PATH:/usr/local/go/bin
go build -trimpath -ldflags "-s -w" -o /opt/rekrutmen/bin/rekrutmen ./cmd/server

# Frontend — butuh Node 22
cd frontend
npm ci
npm run build          # → frontend/dist
```

> `-trimpath` dan `-ldflags "-s -w"` membuat binary lebih kecil dan path build
> tidak membocorkan struktur direktori server.

---

## 6. P3 — Server & reverse proxy

| # | Item | Status | Detail |
|---|---|---|---|
| P-16 | Server Start & boot | ✅ | `/api/health` → `{"ok":true}` |
| P-17 | Timeout header terlindungi | ✅ | `ReadHeaderTimeout: 10s` (`main.go:88`) |
| P-18 | **Unit systemd** | ⬜ | `rekrutmen-tangsel.service` |
| P-19 | **Konfigurasi nginx** | ⬜ | TLS + proxy + static |
| P-20 | **Graceful shutdown** | ⬜ | `main.go:91` — `ListenAndServe` → `Fatalf`, tanpa jeda |
| P-21 | **Watchdog / auto-restart** | ⬜ | `Restart=always` di systemd |
| P-22 | **Health check berkala** | ⬜ | polling `/api/health` |
| P-23 | **Batas ukuran request body** | ⬜ | `http.MaxBytesReader` — cegah payload besar membebani server |

### Contoh unit systemd

```ini
# /etc/systemd/system/rekrutmen-tangsel.service
[Unit]
Description=Aplikasi Rekrutmen Kota Tangerang Selatan
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=rekrutmen
Group=rekrutmen
WorkingDirectory=/opt/rekrutmen
EnvironmentFile=/etc/rekrutmen-tangsel.env
ExecStart=/opt/rekrutmen/bin/rekrutmen
Restart=always
RestartSec=5s

# Pengetatan
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/rekrutmen/data

[Install]
WantedBy=multi-user.target
```

### Contoh konfigurasi nginx

```nginx
server {
    listen 443 ssl http2;
    server_name rekrutmen.example.go.id;

    ssl_certificate     /etc/letsencrypt/live/rekrutmen.example.go.id/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/rekrutmen.example.go.id/privkey.pem;

    # Header keamanan (P-26)
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options    "nosniff" always;
    add_header X-Frame-Options           "DENY" always;
    add_header Referrer-Policy           "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy   "default-src 'self'" always;

    # Artefak frontend
    root /opt/rekrutmen/frontend;
    index index.html;

    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # SPA fallback
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API — proxy ke Go di loopback
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host              $host;
        proxy_set_header X-Real-IP         $remote_addr;
        proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name rekrutmen.example.go.id;
    return 301 https://$host$request_uri;
}
```

---

## 7. P4 — HTTPS & keamanan ⛔

**Semua item di fase ini adalah blocker produksi.**

| # | Item | Status | Detail |
|---|---|---|---|
| P-24 | **Sertifikat TLS** | ⬜ | Let's Encrypt + renew otomatis |
| P-25 | **Redirect HTTP → HTTPS** | ⬜ | contoh di atas |
| P-26 | **Header HSTS / X-Frame-Options / CSP / nosniff** | ⬜ | contoh di atas |
| P-27 | **Flag `Secure` pada cookie sesi** | ⬜ | `auth.go:113`, `auth.go:126` — tanpa ini sesi bisa dicuri |
| P-28 | **Cookie `Secure` honoris proxy** | ⬜ | bila `X-Forwarded-Proto` https |

> **P-27_detail:** saat ini cookie sesi di `backend/internal/auth/auth.go` di-set
> dengan `HttpOnly` dan `SameSite=Lax` tetapi **tanpa `Secure`**. Di belakang
> TLS, `Secure` wajib agar cookie tidak terkirim lewat koneksi HTTP biasa.
> Perbaikannya cukup satu baris per pemanggilan `http.SetCookie`.

---

## 8. P5 — Penyimpanan & data

| # | Item | Status | Detail |
|---|---|---|---|
| P-29 | Provider `local` (JSON) | ✅ | `internal/store/local` |
| P-30 | Provider `supabase` | ✅ | `internal/store/supabase` |
| P-31 | Provider `google_sheets` | ✅ | `internal/store/sheets` |
| P-32 | Skema Supabase + RLS | ✅ | `backend/supabase/schema.sql` |
| P-33 | Migrasi skema | 🔧 | `schema.sql` idempoten (`create table if not exists`), tapi belum ada versioning |
| P-34 | **Backup `data/` berkala** | ⬜ | `data/users.json` + `data/settings.json` + `data/session.secret` |
| P-35 | **Restore dari backup** | ⬜ | perlu diuji (T-74) |
| P-36 | **Backup Supabase** | ⬜ | dump otomatis Supabase |
| P-37 | **Prosedur ekspor/import** | ⬜ | PRD §7 — sebelum ganti provider |
| P-38 | **Rotasi `SESSION_SECRET`** | ⬜ | semua sesi invalid |

### Yang perlu dibackup

| Lokasi | Isi | Sensitif |
|---|---|---|
| `data/users.json` | hash password + salt | **tinggi** — jangan bocorkan |
| `data/session.secret` | kunci penandatangan cookie | **tinggi** |
| `data/settings.json` | provider aktif | rendah |
| `data/*.json` (local store) | data rekrutmen | **tinggi** — data pribadi |
| Supabase | data rekrutmen | **tinggi** |
| Google Sheets | data rekrutmen | **tinggi** |

> Data rekrutmen memuat nomor WhatsApp dan alamat — termasuk **data pribadi**.
> Backup wajib dienkripsi dan tempatnya harus punya kontrol akses.

---

## 9. P6 — CI/CD ⛔

**Belum ada sama sekali.** `.github/` tidak ada di repo.

| # | Item | Status | Detail |
|---|---|---|---|
| P-39 | **Workflow CI: build + test** | ⬜ | tiap push & PR |
| P-40 | **Workflow: build artefak rilis** | ⬜ | saat tag |
| P-41 | **Deploy otomatis ke server** | ⬜ | rsync/ssh saat tag |
| P-42 | **Rilis & rollback terdokumentasi** | ⬜ | prosedur + tombol kembali |

### CI yang direncanakan

```yaml
# .github/workflows/ci.yml
name: CI
on:
  push: { branches: [main] }
  pull_request:
jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.24' }
      - run: go vet ./...
        working-directory: backend
      - run: go test -race -count=1 -cover ./...
        working-directory: backend
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '22' }
      - run: npm ci
        working-directory: frontend
      - run: npm run build
        working-directory: frontend
```

> `go test -race` disarankan — toko `store` dipakai konkuren dan belum pernah
> diuji dengan detektor race.

---

## 10. P7 — Monitoring & operasional

| # | Item | Status | Detail |
|---|---|---|---|
| P-43 | Endpoint health | ✅ | `GET /api/health` → `{"ok":true,"provider":"..."}` |
| P-44 | Log ke stdout/stderr | ✅ | `log.Printf` |
| P-45 | **Health check berkala** | ⬜ | systemd watchdog atau monitor eksternal |
| P-46 | **Rotasi log** | ⬜ | logrotate |
| P-47 | **Monitoring uptime** | ⬜ | Uptime Kuma / Prometheus |
| P-48 | **Alert error** | ⬜ | notifikasi ke admin |
| P-49 | **Backup otomatis terjadwal** | ⬜ | cron |
| P-50 | **DokumentasiProcedure insiden** | ⬜ | apa yang dilakukan saat server mati |
| P-51 | **Disaster recovery teruji** | ⬜ | restore drill |

---

## 11. P8 — Kontainer (opsional, non-blocking)

| # | Item | Status | Detail |
|---|---|---|---|
| P-52 | **Dockerfile multi-stage** | ⬜ | build Go + Vite, runtime minimal |
| P-53 | **docker-compose** | ⬜ | untuk development & demo |
| P-54 | **Volume untuk `data/`** | ⬜ | persistensi data lokal |

---

## 12. Prosedur deploy

### 12.1 Rilis pertama

```bash
# 1. Siapkan server
sudo useradd -r -s /usr/sbin/nologin rekrutmen
sudo mkdir -p /opt/rekrutmen/{bin,data,frontend} /etc/rekrutmen-tangsel
sudo chown -R rekrutmen:rekrutmen /opt/rekrutmen

# 2. Salin artefak
#    backend/rekrutmen  → /opt/rekrutmen/bin/rekrutmen
#    frontend/dist/*   → /opt/rekrutmen/frontend/

# 3. Environment (WAJIB: SESSION_SECRET)
sudo tee /etc/rekrutmen-tangsel.env > /dev/null <<'EOF'
ADDR=127.0.0.1:8080
DATA_DIR=/opt/rekrutmen/data
FRONTEND_DIR=/opt/rekrutmen/frontend
SESSION_SECRET=<hex-acak-64-karakter>
SUPABASE_URL=
SUPABASE_API_KEY=
EOF
sudo chmod 600 /etc/rekrutmen-tangsel.env
sudo chown rekrutmen:rekrutmen /etc/rekrutmen-tangsel.env

# 4. Layani aplikasi
sudo systemctl daemon-reload
sudo systemctl enable --now rekrutmen-tangsel
sudo systemctl status rekrutmen-tangsel

# 5. Health check
curl -fsS http://127.0.0.1:8080/api/health

# 6. TLS + nginx
sudo certbot --nginx -d rekrutmen.example.go.id
sudo nginx -t && sudo systemctl reload nginx

# 7. Ganti password default sebelum dipakai!
```

### 12.2 Rilis berikutnya

```bash
# 1. Build & uji
go test -count=1 ./... && npm run build

# 2. Backup data sebelum menimpa
sudo -u rekrutmen tar czf /backup/data-$(date +%F-%H%M).tar.gz -C /opt/rekrutmen data

# 3. Deploy
sudo systemctl stop rekrutmen-tangsel
sudo cp rekrutmen /opt/rekrutmen/bin/rekrutmen
sudo rsync -a --delete frontend/dist/ /opt/rekrutmen/frontend/
sudo systemctl start rekrutmen-tangsel

# 4. Verifikasi
curl -fsS http://127.0.0.1:8080/api/health
```

### 12.3 Rollback

```bash
sudo systemctl stop rekrutmen-tangsel
sudo cp /backup/rekrutmen-<versi-lama> /opt/rekrutmen/bin/rekrutmen
sudo systemctl start rekrutmen-tangsel
curl -fsS http://127.0.0.1:8080/api/health
```

> Cadangkan **binary versi lama** sebelum menimpa. Data tidak berubah oleh
> release binary (kecuali migrasi skema — dan saat ini belum ada migrasi
> otomatis, `schema.sql` dijalankan manual).

---

## 13. Checklist pra-rilis

| # | Pemeriksaan | Rujukan |
|---|---|---|
| R-01 | `go build ./...` & `go test -count=1 ./...` hijau | TESTING_PLAN §3 |
| R-02 | `npm run build` sukses | P-10 |
| R-03 | `go vet ./...` bersih | — |
| R-04 | Tidak ada secret di commit (`git log -p \| grep`) | D-08 |
| R-05 | `SESSION_SECRET` di-set & berbeda dari default | P-03 |
| R-06 | Password default **sudah diganti** | D-85 |
| R-07 | `ADDR=127.0.0.1:8080` (tidak `0.0.0.0`) | P-01 |
| R-08 | HTTPS aktif, HTTP redirect ke HTTPS | P-25 |
| R-09 | Cookie sesi punya flag `Secure` | P-27 |
| R-10 | Header HSTS / CSP / nosniff aktif | P-26 |
| R-11 | `data/` tidak dapat dibaca user lain | P-34 |
| R-12 | Backup berjalan & **restore sudah diuji** | P-34, P-35 |
| R-13 | Health check terpasang | P-45 |
| R-14 | TLS certificate renewal otomatis | P-24 |
| R-15 | Smoke test manual §7 TESTING_PLAN lulus | TESTING_PLAN §7 |
| R-16 | `data/`, `*.env`, `*.secret` tidak ter-commit | P-04 |

---

## 14. Checklist pasca-rilis

| # | Pemeriksaan |
|---|---|
| R-17 | `/api/health` → `{"ok":true}` |
| R-18 | Login admin & petugas berhasil |
| R-19 | Tambah satu data uji, lalu hapus |
| R-20 | Data uji tidak tertinggal di produksi |
| R-21 | Log server bersih (tidak ada panic/500) |
| R-22 | Monitor uptime tidak mencatat insiden |
| R-23 | Backup berikutnya terjadwal & sukses |

---

## 15. Environment per stadium

| Aspek | Development | Staging | Produksi |
|---|---|---|---|
| `ADDR` | `:8080` | `127.0.0.1:8080` | `127.0.0.1:8080` |
| Provider | `local` | `local` | `supabase` atau `google_sheets` |
| HTTPS | tidak | ya | ya (wajib) |
| Data | dummy | dummy (salinan) | asli |
| Sumber frontend | `npm run dev` | build statis | build statis |
| User | `admin`/`admin123` | kredensial khusus | kredensial khusus |
| `SESSION_SECRET` | generate | tetap | tetap & rahasia |

---

## 16. Yang harus dikerjakan berikutnya

| Urutan | Item | Alasan |
|---|---|---|
| 1 | P-27 flag `Secure` + P-26 header keamanan | **blocker produksi** |
| 2 | P-24/P-25 TLS + redirect | **blocker produksi** |
| 3 | P-34/P-49 backup otomatis + P-35 restore | data pribadi |
| 4 | P-39 CI | mencegah regresi |
| 5 | P-18 systemd + P-19 nginx | operasional |
| 6 | P-20 graceful shutdown | deploy tanpa cutoff mendadak |
| 7 | P-07 `.env.example` + P-08 rotasi secret | onboarding |
| 8 | P-52 Docker (opsional) | portabilitas |

> P-01, P-02, P-03, P-09, P-10, P-11, P-12, P-16, P-17, P-29 s.d. P-32, P-43,
> P-44 sudah ✅ karena kode aplikasinya sudah mendukung — yang belum adalah
> *pengoperasian* di server, bukan kemampuan aplikasinya.

---

## 17. Checkpoint terakhir

| Item | Nilai |
|---|---|
| Commit terakhir | `3d7d05a` |
| Remote `main` | `3d7d05a` (selaras dengan lokal) |
| Working tree | bersih |
| `go build ./...` | ✅ |
| `go test -count=1 ./...` | ✅ 5 paket ok |
| `npm run build` | ✅ |
| `GET /api/health` | ✅ `{"ok":true,"provider":"local"}` |
| Artefak deploy | 2 (binary Go, `frontend/dist`) |
| Infrastruktur server | belum ada |
| CI/CD | belum ada |
| TLS | belum ada |

### Titik kembali kalau terputus

1. `git log --oneline` — commit terakhir yang ✅.
2. `git status` — pastikan tidak ada perubahan menggantung.
3. Lanjutkan dari item `⬜` pertama di bagian 16 (P-27).

---

**Dokumen terkait:** [`DEVELOPMENT_PLAN.md`](DEVELOPMENT_PLAN.md) ·
[`TESTING_PLAN.md`](TESTING_PLAN.md)
