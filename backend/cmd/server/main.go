// Command server runs the Rekrutmen Tangsel backend.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"rekrutmen-tangsel/backend/internal/api"
	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/settings"
	"rekrutmen-tangsel/backend/internal/store/sheets"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	root := env("APP_ROOT", ".")
	dataDir := env("DATA_DIR", filepath.Join(root, "data"))
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
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
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("server berjalan di http://localhost%s (provider: %s)", addr, srv.Store().Name())
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server berhenti: %v", err)
	}
}
