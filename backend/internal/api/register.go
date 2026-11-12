package api

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"rekrutmen-tangsel/backend/internal/auth"
	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/seed"
	"rekrutmen-tangsel/backend/internal/store/sheets"
)

// Limits of a self-registration.
const (
	// minUsernameLength is the shortest username that may be registered.
	minUsernameLength = 4
	// minPasswordLength is the shortest password that may be registered.
	minPasswordLength = 6
	// minAnswerLength is the shortest accepted security answer.
	minAnswerLength = 2

	// registerAttempts is the number of registrations one address may send
	// within registerWindow.
	registerAttempts = 10
	// registerWindow is the period the attempts are counted in.
	registerWindow = 10 * time.Minute
)

// registerLimiter counts the registration attempts per address, so the
// public endpoint cannot be used to fill the sheet with accounts.
type registerLimiter struct {
	mu      sync.Mutex
	entries map[string][]time.Time
}

func newRegisterLimiter() *registerLimiter {
	return &registerLimiter{entries: map[string][]time.Time{}}
}

// allow reports whether another attempt from the address is permitted.
func (l *registerLimiter) allow(address string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	kept := l.entries[address][:0]
	for _, t := range l.entries[address] {
		if now.Sub(t) < registerWindow {
			kept = append(kept, t)
		}
	}
	if len(kept) >= registerAttempts {
		l.entries[address] = kept
		return false
	}
	l.entries[address] = append(kept, now)
	return true
}

// handleRegisterOptions returns what the registration form needs: the
// security questions and the districts a member may choose.
//
// It is public, because the form is shown before signing in.
func (s *Server) handleRegisterOptions(w http.ResponseWriter, r *http.Request) {
	kecamatan := make([]string, 0, 8)
	if list, err := seed.All(); err == nil {
		for _, k := range list {
			kecamatan = append(kecamatan, k.Nama)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"questions": model.SecurityQuestions,
		"kecamatan": kecamatan,
		"available": s.auth.CanRegister(),
	})
}

// handleRegister creates a member account.
//
// The account is written to the account store of the previous application
// with the same layout and hashing, so the member can sign in with it in
// both applications.
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Nama      string `json:"nama"`
		Kecamatan string `json:"kecamatan"`
		Question1 string `json:"question1"`
		Answer1   string `json:"answer1"`
		Question2 string `json:"question2"`
		Answer2   string `json:"answer2"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}

	if !s.registerLimiter.allow(registerAddress(r)) {
		writeError(w, http.StatusTooManyRequests, "terlalu banyak percobaan pendaftaran. Coba lagi nanti.")
		return
	}

	reg := model.Registration{
		Username:  strings.TrimSpace(body.Username),
		Password:  body.Password,
		Nama:      strings.TrimSpace(body.Nama),
		Kecamatan: strings.TrimSpace(body.Kecamatan),
		Question1: strings.TrimSpace(body.Question1),
		Answer1:   strings.TrimSpace(body.Answer1),
		Question2: strings.TrimSpace(body.Question2),
		Answer2:   strings.TrimSpace(body.Answer2),
	}

	if err := validateRegistration(reg); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.auth.RegisterUser(reg); err != nil {
		if errors.Is(err, sheets.ErrUsernameTaken) {
			writeError(w, http.StatusConflict, "username sudah digunakan")
			return
		}
		if errors.Is(err, auth.ErrNoRegistrar) {
			writeError(w, http.StatusServiceUnavailable, "pendaftaran anggota belum tersedia")
			return
		}
		log.Printf("gagal mendaftarkan anggota %s: %v", reg.Username, err)
		writeError(w, http.StatusInternalServerError, "gagal menyimpan akun: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":  "ACTIVE",
		"message": "Registrasi berhasil. Akun langsung aktif dan dapat digunakan untuk masuk.",
	})
}

// validateRegistration checks the registration input.
//
// The rules match the previous application, so a member registers the same
// way in both.
func validateRegistration(reg model.Registration) error {
	if reg.Username == "" {
		return errors.New("username wajib diisi")
	}
	if len([]rune(reg.Username)) < minUsernameLength {
		return errors.New("username minimal 4 karakter")
	}
	if strings.ContainsAny(reg.Username, " \t") {
		return errors.New("username tidak boleh memuat spasi")
	}
	if reg.Password == "" {
		return errors.New("password wajib diisi")
	}
	if len([]rune(reg.Password)) < minPasswordLength {
		return errors.New("password minimal 6 karakter")
	}
	if reg.Nama == "" {
		return errors.New("nama wajib diisi")
	}
	if reg.Kecamatan == "" {
		return errors.New("kecamatan wajib diisi")
	}
	if !seed.IsKecamatan(reg.Kecamatan) {
		return errors.New("kecamatan tidak valid")
	}
	if !isSecurityQuestion(reg.Question1) {
		return errors.New("pertanyaan keamanan 1 tidak valid")
	}
	if !isSecurityQuestion(reg.Question2) {
		return errors.New("pertanyaan keamanan 2 tidak valid")
	}
	if reg.Question1 == reg.Question2 {
		return errors.New("pilih dua pertanyaan keamanan yang berbeda")
	}
	if len([]rune(reg.Answer1)) < minAnswerLength {
		return errors.New("jawaban keamanan 1 minimal 2 karakter")
	}
	if len([]rune(reg.Answer2)) < minAnswerLength {
		return errors.New("jawaban keamanan 2 minimal 2 karakter")
	}
	return nil
}

// isSecurityQuestion reports whether the question is one of the offered ones.
func isSecurityQuestion(q string) bool {
	for _, known := range model.SecurityQuestions {
		if known == q {
			return true
		}
	}
	return false
}

// registerAddress returns the address the attempt comes from.
func registerAddress(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
