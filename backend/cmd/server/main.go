// Command server runs the Rekrutmen Tangsel backend.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"rekrutmen-tangsel/backend/internal/api"
	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/settings"
	"rekrutmen-tangsel/backend/internal/store/sheets"
)

// shutdownTimeout is the time the server gets to finish the requests that
// are running, after a stop signal.
const shutdownTimeout = 20 * time.Second

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// envInt reads a number from the environment, falling back when the value is
// missing or not a number.
func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("%s bukan angka, memakai %d", key, fallback)
		return fallback
	}
	return n
}

// splitList reads a comma separated list from the environment.
func splitList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// envBool reads a yes or no value from the environment.
func envBool(key string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return fallback
}

func main() {
	root := env("APP_ROOT", ".")
	dataDir := env("DATA_DIR", filepath.Join(root, "data"))
	// The directory holds the accounts and the session secret, so only the
	// user of the service may read it.
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		log.Fatalf("gagal membuat direktori data: %v", err)
	}

	// Session signing secret: from the environment, otherwise generated and
	// persisted so logins survive a restart during development.
	secret, err := auth.RandomSecret()
	if err != nil {
		log.Fatalf("gagal membuat secret: %v", err)
	}
	if env := os.Getenv("SESSION_SECRET"); env != "" {
		secret = []byte(env)
	} else if b, readErr := os.ReadFile(filepath.Join(dataDir, "session.secret")); readErr == nil {
		secret = b
	} else {
		if err := os.WriteFile(filepath.Join(dataDir, "session.secret"), secret, 0o600); err != nil {
			log.Fatalf("gagal menyimpan secret: %v", err)
		}
	}

	users, err := auth.ParseUsers(filepath.Join(dataDir, "users.json"))
	if err != nil {
		log.Fatalf("gagal memuat pengguna: %v", err)
	}
	authMgr := auth.NewManager(users, secret)
	authMgr.SetUsersPath(filepath.Join(dataDir, "users.json"))
	// The tokens of the sessions that a password change invalidated. They
	// are kept on disk, so the effect of a password change survives a
	// restart.
	authMgr.SetSessionVersionsPath(filepath.Join(dataDir, "session_versions.json"))
	// The session cookie carries the Secure flag whenever the request is
	// secure. COOKIE_SECURE=true forces it, which is the right setting for a
	// server that always sits behind TLS.
	if os.Getenv("COOKIE_SECURE") == "true" {
		authMgr.SetSecureCookies(true)
	}

	// Accounts of the previous application live in the USERS sheet of the
	// spreadsheet. When it is configured, a user can sign in with the
	// credentials they already had.
	legacyMode := env("LEGACY_AUTH_MODE", "upgrade")
	if legacyMode == "off" {
		log.Printf("login warisan dimatikan (LEGACY_AUTH_MODE=off)")
	} else if os.Getenv("GOOGLE_SHEET_ID") != "" && os.Getenv("GOOGLE_KEY_FILE") != "" {
		legacyStore, lerr := sheets.New(sheets.Config{
			SpreadsheetID: os.Getenv("GOOGLE_SHEET_ID"),
			KeyFile:       os.Getenv("GOOGLE_KEY_FILE"),
			TabPengguna:   os.Getenv("GOOGLE_SHEETS_TAB_PENGGUNA"),
		})
		if lerr != nil {
			log.Printf("login warisan tidak aktif: %v", lerr)
		} else {
			authMgr.SetLegacyAuthenticator(legacyStore)
			// New member accounts are written to the same sheet, so a member
			// can sign in with them in both applications.
			authMgr.SetRegistrar(legacyStore)
			if legacyMode == "sheet" {
				// Every sign in is verified against the sheet, so a local
				// copy of an account never decides a sign in.
				authMgr.SetSheetOnly(true)
				log.Printf("PERINGATAN: mode sheet memverifikasi password dengan hash SHA-256 tanpa salt milik app lama. Gunakan mode upgrade agar password diperiksa dengan hash ber-salt.")
			}
			tabPengguna := os.Getenv("GOOGLE_SHEETS_TAB_PENGGUNA")
			if tabPengguna == "" {
				tabPengguna = sheets.DefaultTabPengguna
			}
			log.Printf("login warisan aktif dari sheet %q (mode: %s)", tabPengguna, legacyMode)
			log.Printf("pendaftaran anggota menulis ke sheet %q", tabPengguna)
		}
	}

	setStore, err := settings.NewStore(filepath.Join(dataDir, "settings.json"))
	if err != nil {
		log.Fatalf("gagal memuat pengaturan: %v", err)
	}

	cfg := api.Config{
		DataDir:        dataDir,
		SupabaseURL:    os.Getenv("SUPABASE_URL"),
		SupabaseAPIKey: os.Getenv("SUPABASE_API_KEY"),
		GoogleSheetID:  os.Getenv("GOOGLE_SHEET_ID"),
		GoogleKeyFile:  os.Getenv("GOOGLE_KEY_FILE"),

		GoogleSheetsMode:         os.Getenv("GOOGLE_SHEETS_MODE"),
		GoogleSheetsTabRekrutmen: os.Getenv("GOOGLE_SHEETS_TAB_REKRUTMEN"),
		GoogleSheetsTabAktivitas: os.Getenv("GOOGLE_SHEETS_TAB_AKTIVITAS"),
		GoogleSheetsTabSumber:    os.Getenv("GOOGLE_SHEETS_TAB_SUMBER"),
		GoogleSheetsTabPengguna:  os.Getenv("GOOGLE_SHEETS_TAB_PENGGUNA"),
		LegacyAuthMode:           legacyMode,
		UploadMaxAttempts:        envInt("UPLOAD_MAX_ATTEMPTS", api.DefaultUploadMaxAttempts),
		LoginMaxAttempts:         envInt("LOGIN_MAX_ATTEMPTS", api.DefaultLoginMaxAttempts),
		AllowedOrigins:           splitList(os.Getenv("ALLOWED_ORIGINS")),
		TrustProxy:               envBool("TRUST_PROXY", false),
	}

	// The header X-Forwarded-For is only read when the reverse proxy of this
	// application is trusted. Without that, every caller behind the proxy
	// shares one address, so the counters of the addresses become the
	// counters of the whole application.
	if !cfg.TrustProxy {
		log.Printf("TRUST_PROXY tidak aktif: header X-Forwarded-For diabaikan dan batas percobaan memakai alamat koneksi")
	} else {
		log.Printf("TRUST_PROXY aktif: alamat pengirim diambil dari X-Forwarded-For")
	}

	srv, err := api.New(cfg, authMgr, setStore)
	if err != nil {
		log.Fatalf("gagal menyiapkan API: %v", err)
	}

	addr := env("ADDR", ":8080")
	mux := http.NewServeMux()
	mux.Handle("/api/", srv.Handler())

	// Serve the built frontend when it is present.
	frontendDir := env("FRONTEND_DIR", filepath.Join(root, "frontend", "dist"))
	if _, err := os.Stat(filepath.Join(frontendDir, "index.html")); err == nil {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			api.ServeFrontend(w, r, frontendDir)
		})
		log.Printf("frontend dilayani dari %s", frontendDir)
	} else {
		log.Printf("frontend belum dibangun; hanya API yang dilayani")
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           api.SecurityHeaders(mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// The server is stopped with a signal, for example when systemd restarts
	// the service. The requests that are running are finished first, so a
	// deployment does not cut them off.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Printf("sinyal diterima, menutup server")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("penutupan server tidak tuntas: %v", err)
		}
	}()

	log.Printf("server berjalan di http://localhost%s (provider: %s)", addr, srv.Store().Name())
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server berhenti: %v", err)
	}
	log.Printf("server ditutup dengan rapi")
}
