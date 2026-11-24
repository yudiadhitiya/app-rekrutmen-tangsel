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
	"path/filepath"
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

	// UploadMaxAttempts is the number of batch uploads one account may
	// start in ten minutes. Zero means the default of five.
	UploadMaxAttempts int

	// AllowedOrigins are the origins that may call the API from a browser.
	// An empty list means only the same origin may do so.
	AllowedOrigins []string

	// LoginMaxAttempts is the number of failed sign ins one account or one
	// address may make before the sign in is refused. Zero means the default.
	LoginMaxAttempts int

	// TrustProxy says that a reverse proxy of this application writes the
	// real address of the caller in the X-Forwarded-For header. The header is
	// only read when this is true, because a caller can set it by itself.
	//
	// Without it the address of the connection is used. Behind a reverse
	// proxy that means every request shares the address of the proxy.
	TrustProxy bool

	// AllowSelfRegistration says that strangers may create their own account.
	//
	// A new account writes into the shared spreadsheet, so an open endpoint
	// lets anybody on the internet put rows into the production data. It is
	// therefore off unless the operator turns it on deliberately.
	AllowSelfRegistration bool
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

	// loginLimiter counts the failed sign ins per account and per address,
	// so the login cannot be used to guess a password.
	loginLimiter *loginLimiter

	// uploadLimiter counts the batch upload attempts per account.
	uploadLimiter *uploadLimiter

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
		registerLimiter: newRegisterLimiter(newLimiterStore(limiterStorePath(cfg.DataDir, "limiter_register.json"))),
		loginLimiter:    newLoginLimiter(cfg.LoginMaxAttempts, newLimiterStore(limiterStorePath(cfg.DataDir, "limiter_login.json"))),
		uploadLimiter:   newUploadLimiter(cfg.UploadMaxAttempts, newLimiterStore(limiterStorePath(cfg.DataDir, "limiter_upload.json"))),
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

// writeStoreError writes the answer for a failure of the server itself.
//
// The details of the failure go to the log of the server, not to the answer.
// A detail may hold a path of the disk, a query of the provider or a name of
// a file, and the browser has no need for it.
func writeStoreError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%s: %v", msg, err)
	writeError(w, http.StatusInternalServerError, msg)
}

// MaxJSONBody is the largest request body that is read as JSON.
//
// The limit keeps a single request from taking the memory of the server. The
// batch upload has its own, larger limit.
const MaxJSONBody = 1 << 20

// decodeJSON reads the body of a request as JSON.
//
// The body is limited, so a request that is too large is refused instead of
// being read into the memory of the server.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBody)
	return json.NewDecoder(r.Body).Decode(dst)
}

// bodyError writes the answer for a body that could not be read.
//
// A body that is too large gets its own status, so the client knows it must
// send less data.
func bodyError(w http.ResponseWriter, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeError(w, http.StatusRequestEntityTooLarge, "body terlalu besar")
		return
	}
	writeError(w, http.StatusBadRequest, "body tidak valid")
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

// requireAdmin wraps a handler so it only runs for an administrator.
//
// The access level decides, not the plain role. Only a full administrator
// (level 3) passes, which keeps the actions of the previous application that
// were reserved for its top administrators.
func (s *Server) requireAdmin(next func(http.ResponseWriter, *http.Request, model.User)) http.HandlerFunc {
	return s.requireAuth(func(w http.ResponseWriter, r *http.Request, user model.User) {
		if !user.IsFullAdmin() {
			writeError(w, http.StatusForbidden, "hanya admin yang dapat melakukan aksi ini")
			return
		}
		next(w, r, user)
	})
}

// requireDelete wraps a handler so it only runs for an account that may
// delete a record.
//
// It follows the previous application: a member deletes the records it
// created itself, and a full administrator deletes any record. A district
// administrator may not delete, just as in the previous application.
//
// The record itself is checked by the handler, which hides a record outside
// the scope of the account.
func (s *Server) requireDelete(next func(http.ResponseWriter, *http.Request, model.User)) http.HandlerFunc {
	return s.requireAuth(func(w http.ResponseWriter, r *http.Request, user model.User) {
		switch user.LevelOrDefault() {
		case model.LevelAnggota, model.LevelAdmin:
			next(w, r, user)
		default:
			writeError(w, http.StatusForbidden, "hanya anggota penginput atau admin penuh yang dapat menghapus data")
		}
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
	mux.HandleFunc("PUT /api/me/password", s.requireAuth(s.handleChangePassword))

	// Self-registration of members.
	mux.HandleFunc("GET /api/register/options", s.handleRegisterOptions)
	mux.HandleFunc("POST /api/register", s.handleRegister)

	// Master wilayah.
	mux.HandleFunc("GET /api/wilayah", s.requireAuth(s.handleWilayah))

	// Dashboard rekapitulasi.
	mux.HandleFunc("GET /api/dashboard", s.requireAuth(s.handleDashboard))

	// Master aktivitas.
	mux.HandleFunc("GET /api/aktivitas", s.requireAuth(s.handleListAktivitas))
	mux.HandleFunc("POST /api/aktivitas", s.requireAdmin(s.handleCreateAktivitas))
	mux.HandleFunc("PUT /api/aktivitas/{id}", s.requireAdmin(s.handleUpdateAktivitas))
	mux.HandleFunc("DELETE /api/aktivitas/{id}", s.requireAdmin(s.handleDeleteAktivitas))

	// Rekrutmen.
	mux.HandleFunc("GET /api/rekrutmen", s.requireAuth(s.handleListRekrutmen))
	mux.HandleFunc("GET /api/rekrutmen/export", s.requireAuth(s.handleExportRekrutmen))
	mux.HandleFunc("GET /api/rekrutmen/check-wa", s.requireAuth(s.handleCheckWA))
	mux.HandleFunc("GET /api/rekrutmen/{id}", s.requireAuth(s.handleGetRekrutmen))
	mux.HandleFunc("POST /api/rekrutmen", s.requireAuth(s.handleCreateRekrutmen))
	mux.HandleFunc("POST /api/rekrutmen/batch", s.requireUpload(s.handleBatchCreateRekrutmen))
	mux.HandleFunc("PUT /api/rekrutmen/{id}", s.requireAuth(s.handleUpdateRekrutmen))
	mux.HandleFunc("DELETE /api/rekrutmen/{id}", s.requireDelete(s.handleDeleteRekrutmen))

	// Settings.
	mux.HandleFunc("GET /api/settings", s.requireAdmin(s.handleGetSettings))
	mux.HandleFunc("PUT /api/settings", s.requireAdmin(s.handlePutSettings))

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "provider": s.Store().Name()})
	})

	return withCORS(mux, s.cfg.AllowedOrigins)
}

// ---------------------------------------------------------------------------
// Auth handlers
// ---------------------------------------------------------------------------

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		bodyError(w, err)
		return
	}
	username := strings.TrimSpace(body.Username)
	address := s.clientAddress(r)

	// The limit of the address refuses before the password is checked, so a
	// flood from one address costs the server nothing.
	if s.loginAddressLocked(w, r, username) {
		return
	}

	// An account that is being guessed at may only try its password every
	// few seconds, so a flood against one account cannot use up the server.
	if !s.loginLimiter.mayTry(username, address) {
		w.Header().Set("Retry-After", "10")
		writeError(w, http.StatusTooManyRequests,
			"terlalu banyak percobaan login untuk akun ini. Coba lagi sebentar.")
		return
	}

	user, fromLegacy, err := s.auth.VerifyLogin(body.Username, body.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			// The failed attempt is counted, so the password cannot be
			// guessed by trying many times.
			locked := s.loginLimiter.locked(username)
			s.loginLimiter.fail(username, address)
			if locked {
				// The account is being guessed at, so the wrong password is
				// refused. A password that is right always gets in, so an
				// attacker cannot lock the account out of the application.
				if wait := s.loginLimiter.lockedFor(username); wait > 0 {
					w.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds())+1))
				}
				writeError(w, http.StatusTooManyRequests,
					"terlalu banyak percobaan login untuk akun ini. Coba lagi nanti.")
				return
			}
			writeError(w, http.StatusUnauthorized, "username atau password salah")
			return
		}
		// The account store could not be reached; that is not the user's
		// fault, so it is not reported as a wrong password.
		writeStoreError(w, "gagal memverifikasi akun", err)
		return
	}
	s.loginLimiter.reset(username, address)

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

	if err := s.auth.SetSession(w, r, user); err != nil {
		writeError(w, http.StatusInternalServerError, "gagal membuat sesi")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": user, "migrated": migrated})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.auth.ClearSession(w, r)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request, user model.User) {
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

// handleChangePassword replaces the password of the signed in account.
//
// Every account may change its own password, so the access level is not
// checked. The current password must be given, so a borrowed session cannot
// take the account over.
func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request, user model.User) {
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		bodyError(w, err)
		return
	}
	if body.CurrentPassword == "" {
		writeError(w, http.StatusBadRequest, "password lama wajib diisi")
		return
	}
	if fields := validate.Password(body.NewPassword); len(fields) > 0 {
		writeFieldError(w, "validasi gagal", fields)
		return
	}
	if body.CurrentPassword == body.NewPassword {
		writeError(w, http.StatusBadRequest, "password baru harus berbeda dari password lama")
		return
	}

	switch err := s.auth.ChangePassword(user.Username, body.CurrentPassword, body.NewPassword); {
	case errors.Is(err, auth.ErrInvalidCredentials):
		writeError(w, http.StatusBadRequest, "password lama salah")
	case errors.Is(err, auth.ErrAccountNotFound):
		writeError(w, http.StatusNotFound, "akun tidak ditemukan")
	case err != nil:
		writeStoreError(w, "gagal mengubah password", err)
	default:
		// The password change invalidated the sessions that were made
		// before it, including this one. A fresh session is written, so the
		// browser that changed the password stays signed in.
		if err := s.auth.SetSession(w, r, user); err != nil {
			writeError(w, http.StatusInternalServerError, "gagal membuat sesi")
			return
		}
		// The password of the first administrator was written to a file when
		// it was made randomly. That file is useless once the password has
		// been replaced, so it is removed.
		s.removeInitialPasswordFile()
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	}
}

// removeInitialPasswordFile removes the file that holds the password of the
// first administrator.
//
// The failure is not reported: the file is a convenience, and the password
// has already been replaced.
func (s *Server) removeInitialPasswordFile() {
	if s.cfg.DataDir == "" {
		return
	}
	path := filepath.Join(s.cfg.DataDir, "initial_admin_password.txt")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		log.Printf("gagal menghapus berkas password awal: %v", err)
	}
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
		writeStoreError(w, "gagal memuat daftar aktivitas", err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleCreateAktivitas(w http.ResponseWriter, r *http.Request, user model.User) {
	var body struct {
		Nama string `json:"nama"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		bodyError(w, err)
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
		writeStoreError(w, "gagal menyimpan aktivitas", err)
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
	if err := decodeJSON(w, r, &body); err != nil {
		bodyError(w, err)
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
			writeStoreError(w, "gagal memperbarui aktivitas", err)
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
			writeStoreError(w, "gagal menghapus aktivitas", err)
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

// ---------------------------------------------------------------------------
// Rekrutmen
// ---------------------------------------------------------------------------

func (s *Server) handleListRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	f := rekrutmenFilterFromQuery(r, user)
	// The scope of the account decides which records are reachable. It is
	// part of the filter, so it is applied before the rows are paginated and
	// the reported total stays correct.
	f.Scope = user.ScopeFor()
	data, total, err := s.Store().ListRekrutmen(f)
	if err != nil {
		writeStoreError(w, "gagal memuat data", err)
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
		writeStoreError(w, "gagal memuat data", err)
		return
	}
	if !canAccess(user, *rec) {
		writeError(w, http.StatusNotFound, "data tidak ditemukan")
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

// canAccess reports whether a user may read or change a record.
//
// The scope of the account decides, not the plain role: a full administrator
// reaches every record, a district administrator only the records of its own
// district, and a member only the records it created itself. A record without
// an owner stays with the administrators.
func canAccess(user model.User, rec model.Rekrutmen) bool {
	return model.MatchScope(user.ScopeFor(), rec)
}

// checkDistrictScope refuses a write that would put a record into the scope
// of another district.
//
// The scope of an account decides which records it may reach, and the scope
// of a district administrator is the district of the record. Without this
// check the district of the request body decides that, so an administrator
// of one district could move a record into the queue of another district, or
// hide one from their own.
//
// A full administrator reaches every district, so it may write any of them.
//
// A district administrator only reaches the records of its own district, so it
// may only write in that district.
//
// A member only reaches the records it created itself. It may name any
// district when it creates a record, but it may not move a record that
// already exists, because that would put the record into the queue of
// another district.
func checkDistrictScope(user model.User, rec *model.Rekrutmen, stored *model.Rekrutmen) error {
	switch user.LevelOrDefault() {
	case model.LevelAdmin:
		return nil
	case model.LevelAdminKecamatan:
		own := strings.TrimSpace(user.Kecamatan)
		if own == "" {
			// The account reaches only the records it created itself, so it
			// has no district of its own to write in.
			return checkMemberDistrict(rec, stored)
		}
		if !strings.EqualFold(strings.TrimSpace(rec.Kecamatan), own) {
			return errors.New("kecamatan data harus sama dengan kecamatan Anda")
		}
		return nil
	default:
		return checkMemberDistrict(rec, stored)
	}
}

// checkMemberDistrict keeps the district of a stored record where it is.
//
// The district of a record decides the queue of the district administrator
// that sees it, so moving it hands the record to another district.
func checkMemberDistrict(rec *model.Rekrutmen, stored *model.Rekrutmen) error {
	if stored == nil {
		return nil
	}
	if !strings.EqualFold(strings.TrimSpace(rec.Kecamatan), strings.TrimSpace(stored.Kecamatan)) {
		return errors.New("kecamatan data tidak boleh diubah")
	}
	return nil
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
		writeStoreError(w, "gagal memeriksa nomor", err)
		return
	}
	// A record outside the scope of the account is reported as existing
	// without its details, so a duplicate number is still refused while the
	// other district stays hidden.
	if !canAccess(user, *rec) {
		writeJSON(w, http.StatusOK, map[string]any{"exists": true})
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

func (s *Server) decodeRekrutmen(w http.ResponseWriter, r *http.Request) (*model.Rekrutmen, error) {
	var rec model.Rekrutmen
	if err := decodeJSON(w, r, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (s *Server) handleCreateRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	rec, err := s.decodeRekrutmen(w, r)
	if err != nil {
		bodyError(w, err)
		return
	}
	if err := checkDistrictScope(user, rec, nil); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
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
		writeStoreError(w, "gagal menyimpan data", err)
		return
	}
	writeJSON(w, http.StatusCreated, rec)
}

func (s *Server) handleUpdateRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	rec, err := s.decodeRekrutmen(w, r)
	if err != nil {
		bodyError(w, err)
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
		writeStoreError(w, "gagal memuat data", err)
		return
	}
	if !canAccess(user, *prev) {
		writeError(w, http.StatusNotFound, "data tidak ditemukan")
		return
	}
	// The district of the record decides whose scope it belongs to, so it may
	// not be moved into the scope of another district.
	if err := checkDistrictScope(user, rec, prev); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
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
			writeStoreError(w, "gagal memperbarui data", err)
		}
		return
	}
	writeJSON(w, http.StatusOK, rec)
}

func (s *Server) handleDeleteRekrutmen(w http.ResponseWriter, r *http.Request, user model.User) {
	rec, err := s.Store().GetRekrutmen(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeStoreError(w, "gagal memuat data", err)
		return
	}
	// A record outside the scope of the account does not exist for it.
	if !canAccess(user, *rec) {
		writeError(w, http.StatusNotFound, "data tidak ditemukan")
		return
	}
	if err := s.Store().DeleteRekrutmen(rec.ID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "data tidak ditemukan")
			return
		}
		writeStoreError(w, "gagal menghapus data", err)
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
	if err := decodeJSON(w, r, &body); err != nil {
		bodyError(w, err)
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
		log.Printf("provider %q tidak bisa disiapkan: %v", body.Provider, err)
		writeError(w, http.StatusBadRequest, "provider itu belum dikonfigurasi di server")
		return
	}
	if err := target.Ping(); err != nil {
		log.Printf("provider %q tidak dapat diakses: %v", body.Provider, err)
		writeError(w, http.StatusBadGateway, "provider tidak dapat diakses")
		return
	}
	// Seed the master activity list on the new provider, otherwise the picklist
	// would be empty and no record could be created after the switch.
	if err := seedAktivitasIfEmpty(target); err != nil {
		writeStoreError(w, "gagal menyiapkan master aktivitas", err)
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
