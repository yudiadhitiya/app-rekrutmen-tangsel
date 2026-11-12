// Package api exposes the REST endpoints of the application.
package api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/seed"
	"rekrutmen-tangsel/backend/internal/settings"
	"rekrutmen-tangsel/backend/internal/store"
	"rekrutmen-tangsel/backend/internal/store/local"
	"rekrutmen-tangsel/backend/internal/store/sheets"
	"rekrutmen-tangsel/backend/internal/store/supabase"
	"rekrutmen-tangsel/backend/internal/validate"
)

// Config holds everything the API needs to run.
type Config struct {
	DataDir string

	// Provider credentials, read from the environment.
	SupabaseURL    string
	SupabaseAPIKey string
	GoogleSheetID  string
	GoogleKeyFile  string

	// Google Sheets layout. Empty values fall back to the store defaults.
	GoogleSheetsMode         string
	GoogleSheetsTabRekrutmen string
	GoogleSheetsTabAktivitas string
	GoogleSheetsTabSumber    string
	GoogleSheetsTabPengguna  string

	// LegacyAuthMode says how accounts of the previous application are
	// handled:
	//
	//	"upgrade" - the first sign in with the old credentials stores the
	//	            account locally with a salted hash; later sign ins use it
	//	"sheet"   - every sign in is verified against the sheet
	//	"off"     - the old accounts are not used
	//
	// An empty value means "upgrade".
	LegacyAuthMode string
}

// Server holds the HTTP handlers and the active storage provider.
type Server struct {
	cfg      Config
	auth     *auth.Manager
	settings *settings.Store

	mu    sync.RWMutex
	store store.Store

	// registerLimiter counts the self-registration attempts per address.
	registerLimiter *registerLimiter

	// storeFactory builds a provider on demand, so an admin can switch
	// providers without restarting the process.
	storeFactory func(provider string) (store.Store, error)
}

// New builds a Server and selects the configured provider.
func New(cfg Config, authMgr *auth.Manager, setStore *settings.Store) (*Server, error) {
	s := &Server{
		cfg:             cfg,
		auth:            authMgr,
		settings:        setStore,
		registerLimiter: newRegisterLimiter(),
		storeFactory: func(provider string) (store.Store, error) {
			switch provider {
			case settings.ProviderLocal:
				return local.New(cfg.DataDir + "/rekrutmen.json")
			case settings.ProviderSupa:
				if cfg.SupabaseURL == "" || cfg.SupabaseAPIKey == "" {
					return nil, errors.New("supabase: SUPABASE_URL dan SUPABASE_API_KEY belum diatur")
				}
				return supabase.New(cfg.SupabaseURL, cfg.SupabaseAPIKey, ""), nil
			case settings.ProviderSheets:
				if cfg.GoogleSheetID == "" || cfg.GoogleKeyFile == "" {
					return nil, errors.New("google sheets: GOOGLE_SHEET_ID dan GOOGLE_KEY_FILE belum diatur")
				}
				return sheets.New(sheets.Config{
					SpreadsheetID: cfg.GoogleSheetID,
					KeyFile:       cfg.GoogleKeyFile,
					TabRekrutmen:  cfg.GoogleSheetsTabRekrutmen,
					TabAktivitas:  cfg.GoogleSheetsTabAktivitas,
					TabSumber:     cfg.GoogleSheetsTabSumber,
					Mode:          cfg.GoogleSheetsMode,
				})
			default:
				return nil, errors.New("provider tidak dikenal: " + provider)
			}
		},
	}

	st, err := s.storeFactory(setStore.Get().Provider)
	if err != nil {
		return nil, err
	}
	s.store = st
	// Ping the provider so a misconfigured provider is reported early. It is a
	// warning, not a fatal error, so the admin can still reach the settings
	// page and switch to a working provider.
	if err := st.Ping(); err != nil {
		log.Printf("provider %s belum siap: %v", st.Name(), err)
	}
	if err := seedAktivitasIfEmpty(st); err != nil {
		return nil, err
	}
	return s, nil
}

// seedAktivitasIfEmpty populates the master activity list on a fresh install.
func seedAktivitasIfEmpty(st store.Store) error {
	list, err := st.ListAktivitas()
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return nil
	}
	for _, name := range seed.Aktivitas {
		a := &model.Aktivitas{ID: newID(), Nama: name, Aktif: true}
		if err := st.CreateAktivitas(a); err != nil {
			return err
		}
	}
	return nil
}

// Store returns the active storage provider.
func (s *Server) Store() store.Store {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store
}

// newID returns a random RFC 4122 UUID string. The dashed form matters: the
// Supabase schema uses uuid columns, which reject a bare hex string.
func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Falling back to a timestamp-derived id is still unique enough here.
		return fmt.Sprintf("%08x-%04x-4%03x-8%03x-%012x",
			time.Now().UnixNano(), time.Now().Unix()&0xffff,
			time.Now().Nanosecond()&0xfff, time.Now().Unix()&0x3fff,
			time.Now().UnixNano()&0xffffffffffff)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}

func writeFieldError(w http.ResponseWriter, msg string, fields validate.Fields) {
	writeJSON(w, http.StatusBadRequest, map[string]any{
		"error":  msg,
		"fields": fields,
	})
}

// currentUser returns the authenticated user for a request.
func (s *Server) currentUser(r *http.Request) (model.User, bool) {
	return s.auth.UserFromRequest(r)
}

// requireAuth wraps a handler so it only runs for authenticated users.
func (s *Server) requireAuth(next func(http.ResponseWriter, *http.Request, model.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := s.currentUser(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "sesi tidak valid, silakan login kembali")
			return
		}
		next(w, r, user)
	}
}

// requireAdmin wraps a handler so it only runs for admins.
func (s *Server) requireAdmin(next func(http.ResponseWriter, *http.Request, model.User)) http.HandlerFunc {
	return s.requireAuth(func(w http.ResponseWriter, r *http.Request, user model.User) {
		if user.Role != model.RoleAdmin {
			writeError(w, http.StatusForbidden, "hanya admin yang dapat melakukan aksi ini")
			return
		}
		next(w, r, user)
	})
}

// ---------------------------------------------------------------------------
// Routing
// ---------------------------------------------------------------------------

// Handler returns the HTTP handler for the whole API.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	// Auth.
	mux.HandleFunc("POST /api/login", s.handleLogin)
	mux.HandleFunc("POST /api/logout", s.handleLogout)
	mux.HandleFunc("GET /api/me", s.requireAuth(s.handleMe))

	// Self-registration of members.
	mux.HandleFunc("GET /api/register/options", s.handleRegisterOptions)
	mux.HandleFunc("POST /api/register", s.handleRegister)

	// Master wilayah.
	mux.HandleFunc("GET /api/wilayah", s.requireAuth(s.handleWilayah))

	// Master aktivitas.
	mux.HandleFunc("GET /api/aktivitas", s.requireAuth(s.handleListAktivitas))
	mux.HandleFunc("POST /api/aktivitas", s.requireAdmin(s.handleCreateAktivitas))
	mux.HandleFunc("PUT /api/aktivitas/{id}", s.requireAdmin(s.handleUpdateAktivitas))
	mux.HandleFunc("DELETE /api/aktivitas/{id}", s.requireAdmin(s.handleDeleteAktivitas))

	// Rekrutmen.
	mux.HandleFunc("GET /api/rekrutmen", s.requireAuth(s.handleListRekrutmen))
	mux.HandleFunc("GET /api/rekrutmen/check-wa", s.requireAuth(s.handleCheckWA))
	mux.HandleFunc("GET /api/rekrutmen/{id}", s.requireAuth(s.handleGetRekrutmen))
	mux.HandleFunc("POST /api/rekrutmen", s.requireAuth(s.handleCreateRekrutmen))
	mux.HandleFunc("PUT /api/rekrutmen/{id}", s.requireAuth(s.handleUpdateRekrutmen))
	mux.HandleFunc("DELETE /api/rekrutmen/{id}", s.requireAuth(s.handleDeleteRekrutmen))

	// Settings.
	mux.HandleFunc("GET /api/settings", s.requireAdmin(s.handleGetSettings))
	mux.HandleFunc("PUT /api/settings", s.requireAdmin(s.handlePutSettings))

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "provider": s.Store().Name()})
	})

	return withCORS(mux)
}

// withCORS allows the Vite dev server to call the API.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------
// Auth handlers
// ---------------------------------------------------------------------------

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	user, fromLegacy, err := s.auth.VerifyLogin(body.Username, body.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "username atau password salah")
			return
		}
		// The account store could not be reached; that is not the user's
		// fault, so it is not reported as a wrong password.
		writeError(w, http.StatusInternalServerError, "gagal memverifikasi akun: "+err.Error())
		return
	}

	// The account came from the previous application. Store it locally, so
	// the next sign in uses the salted hash of this application.
	migrated := false
	if fromLegacy && s.cfg.LegacyAuthMode != "sheet" {
		if err := s.auth.MigrateLegacyUser(user, body.Password); err != nil {
			log.Printf("gagal memindahkan akun %s: %v", user.Username, err)
		} else {
			migrated = true
		}
	}

	if err := s.auth.SetSession(w, user); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal membuat sesi")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user, "migrated": migrated})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.auth.ClearSession(w)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request, user model.User) {
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

// ---------------------------------------------------------------------------
// Master wilayah
// ---------------------------------------------------------------------------

func (s *Server) handleWilayah(w http.ResponseWriter, r *http.Request, user model.User) {
	list, err := seed.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat master wilayah")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"kecamatan": list})
}

// ---------------------------------------------------------------------------
// Master aktivitas
// ---------------------------------------------------------------------------

func (s *Server) handleListAktivitas(w http.ResponseWriter, r *http.Request, user model.User) {
	list, err := s.Store().ListAktivitas()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat daftar aktivitas: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleCreateAktivitas(w http.ResponseWriter, r *http.Request, user model.User) {
	var body struct {
		Nama string `json:"nama"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if fields := validate.Aktivitas(body.Nama); len(fields) > 0 {
		writeFieldError(w, "validasi gagal", fields)
		return
	}
	a := &model.Aktivitas{ID: newID(), Nama: strings.TrimSpace(body.Nama), Aktif: true}
	if err := s.Store().CreateAktivitas(a); err != nil {
		if errors.Is(err, store.ErrDuplicateName) {
			writeError(w, http.StatusConflict, "nama aktivitas sudah digunakan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal menyimpan aktivitas: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) handleUpdateAktivitas(w http.ResponseWriter, r *http.Request, user model.User) {
	id := r.PathValue("id")
	var body struct {
		Nama  string `json:"nama"`
		Aktif *bool  `json:"aktif"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if fields := validate.Aktivitas(body.Nama); len(fields) > 0 {
		writeFieldError(w, "validasi gagal", fields)
		return
	}
	// Preserve the existing active flag when the client omits it.
	current, err := s.Store().ListAktivitas()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat aktivitas")
		return
	}
	aktif := true
	found := false
	for _, a := range current {
		if a.ID == id {
			aktif = a.Aktif
			found = true
			break
		}
	}
	if !found {
		writeError(w, http.StatusNotFound, "aktivitas tidak ditemukan")
		return
	}
	if body.Aktif != nil {
		aktif = *body.Aktif
	}
	a := &model.Aktivitas{ID: id, Nama: strings.TrimSpace(body.Nama), Aktif: aktif}
	if err := s.Store().UpdateAktivitas(a); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateName):
			writeError(w, http.StatusConflict, "nama aktivitas sudah digunakan")
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "aktivitas tidak ditemukan")
		default:
			writeError(w, http.StatusInternalServerError, "gagal memperbarui aktivitas: "+err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleDeleteAktivitas(w http.ResponseWriter, r *http.Request, user model.User) {
	id := r.PathValue("id")
	if err := s.Store().DeleteAktivitas(id); err != nil {
		switch {
		case errors.Is(err, store.ErrAktivitasUsed):
			writeError(w, http.StatusConflict, "aktivitas masih dipakai data rekrutmen")
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "aktivitas tidak ditemukan")
		default:
			writeError(w, http.StatusInternalServerError, "gagal menghapus aktivitas: "+err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// ---------------------------------------------------------------------------
// Rekrutmen
// ---------------------------------------------------------------------------

func (s *Server) handleListRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	f := model.Filter{
		Q:         q.Get("q"),
		Kecamatan: q.Get("kecamatan"),
		Kelurahan: q.Get("kelurahan"),
		Metode:    model.NormalizeMetode(q.Get("metode")),
		Status:    q.Get("status"),
		Page:      page,
		Limit:     limit,
	}
	data, total, err := s.Store().ListRekrutmen(f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal memuat data: "+err.Error())
		return
	}
	if data == nil {
		data = []model.Rekrutmen{}
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 25
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"total": total,
		"page":  f.Page,
		"limit": f.Limit,
		"data":  data,
	})
}

func (s *Server) handleGetRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	rec, err := s.Store().GetRekrutmen(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memuat data: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

// handleCheckWA reports whether a WhatsApp number is already registered, so the
// UI can warn before saving.
func (s *Server) handleCheckWA(w http.ResponseWriter, r *http.Request, user model.User) {
	no := strings.TrimSpace(r.URL.Query().Get("no"))
	if no == "" {
		writeJSON(w, http.StatusOK, map[string]any{"exists": false})
		return
	}
	rec, err := s.Store().FindRekrutmenByWA(no)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeJSON(w, http.StatusOK, map[string]any{"exists": false})
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memeriksa nomor: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"exists": true, "id": rec.ID, "nama": rec.Nama})
}

// aktivitasNames returns the master activity names allowed by the provider.
func (s *Server) aktivitasNames() []string {
	list, err := s.Store().ListAktivitas()
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, a := range list {
		out = append(out, a.Nama)
	}
	return out
}

func (s *Server) decodeRekrutmen(r *http.Request) (*model.Rekrutmen, error) {
	var rec model.Rekrutmen
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (s *Server) handleCreateRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	rec, err := s.decodeRekrutmen(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if fields := validate.Rekrutmen(rec, s.aktivitasNames()); len(fields) > 0 {
		writeFieldError(w, "validasi gagal", fields)
		return
	}
	now := time.Now().UTC()
	rec.ID = newID()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	rec.CreatedBy = user.Username
	rec.UpdatedBy = user.Username

	if err := s.Store().CreateRekrutmen(rec); err != nil {
		if errors.Is(err, store.ErrDuplicateWA) {
			writeError(w, http.StatusConflict, "nomor whatsapp sudah terdaftar")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal menyimpan data: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rec)
}

func (s *Server) handleUpdateRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	rec, err := s.decodeRekrutmen(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if fields := validate.Rekrutmen(rec, s.aktivitasNames()); len(fields) > 0 {
		writeFieldError(w, "validasi gagal", fields)
		return
	}
	id := r.PathValue("id")
	rec.ID = id

	// The audit fields are owned by the server, not the client: the request
	// body carries only the editable fields, so preserve the stored values
	// instead of overwriting them with zero values.
	prev, err := s.Store().GetRekrutmen(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal memuat data: "+err.Error())
		return
	}
	rec.CreatedAt = prev.CreatedAt
	rec.CreatedBy = prev.CreatedBy
	rec.UpdatedAt = time.Now().UTC()
	rec.UpdatedBy = user.Username

	if err := s.Store().UpdateRekrutmen(rec); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateWA):
			writeError(w, http.StatusConflict, "nomor whatsapp sudah terdaftar")
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
		default:
			writeError(w, http.StatusInternalServerError, "gagal memperbarui data: "+err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (s *Server) handleDeleteRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	if err := s.Store().DeleteRekrutmen(r.PathValue("id")); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal menghapus data: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// ---------------------------------------------------------------------------
// Settings
// ---------------------------------------------------------------------------

// handleGetSettings reports the runtime configuration to an admin.
func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request, user model.User) {
	cur := s.settings.Get()
	writeJSON(w, http.StatusOK, map[string]any{
		"provider":         cur.Provider,
		"providers":        settings.Providers,
		"active_provider":  s.Store().Name(),
		"supabase_ready":   s.cfg.SupabaseURL != "" && s.cfg.SupabaseAPIKey != "",
		"google_sheets_ok": s.cfg.GoogleSheetID != "" && s.cfg.GoogleKeyFile != "",
		"google_sheets": map[string]any{
			"mode":      sheetsMode(s.cfg.GoogleSheetsMode),
			"rekrutmen": sheetsTab(s.cfg.GoogleSheetsTabRekrutmen, sheets.DefaultTabRekrutmen),
			"aktivitas": sheetsTab(s.cfg.GoogleSheetsTabAktivitas, sheets.DefaultTabAktivitas),
			"sumber":    s.cfg.GoogleSheetsTabSumber,
		},
	})
}

// sheetsMode returns the configured Sheets mode, defaulting to app.
func sheetsMode(mode string) string {
	if mode == "" {
		return sheets.ModeApp
	}
	return mode
}

// sheetsTab returns the configured tab name, defaulting to the given value.
func sheetsTab(name, fallback string) string {
	if name == "" {
		return fallback
	}
	return name
}

func (s *Server) handlePutSettings(w http.ResponseWriter, r *http.Request, user model.User) {
	var body struct {
		Provider string `json:"provider"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if !settings.IsValidProvider(body.Provider) {
		writeError(w, http.StatusBadRequest, "provider tidak dikenal")
		return
	}
	// Build the target provider first so an unconfigured provider fails here
	// rather than leaving the app without storage.
	target, err := s.storeFactory(body.Provider)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := target.Ping(); err != nil {
		writeError(w, http.StatusBadGateway, "provider tidak dapat diakses: "+err.Error())
		return
	}
	// Seed the master activity list on the new provider, otherwise the picklist
	// would be empty and no record could be created after the switch.
	if err := seedAktivitasIfEmpty(target); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal menyiapkan master aktivitas: "+err.Error())
		return
	}
	if err := s.settings.Set(settings.Settings{Provider: body.Provider}); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal menyimpan pengaturan")
		return
	}
	s.mu.Lock()
	s.store = target
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{
		"provider": body.Provider,
		"note":     "provider diganti; data lama tidak dipindahkan otomatis",
	})
}

// ---------------------------------------------------------------------------
// Static frontend (production build)
// ---------------------------------------------------------------------------

// ServeFrontend serves the built frontend from disk when the directory exists.
func ServeFrontend(w http.ResponseWriter, r *http.Request, dir string) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	full := dir + "/" + path
	if info, err := os.Stat(full); err == nil && !info.IsDir() {
		http.ServeFile(w, r, full)
		return
	}
	// SPA fallback: unknown paths render the app shell.
	http.ServeFile(w, r, dir+"/index.html")
}
