package sheets

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// usersGrid builds a USERS tab with one existing account.
func usersGrid() [][]string {
	return [][]string{
		{"USERNAME", "PASSWORD_HASH", "ROLE", "NAMA ANGGOTA", "NAMA KECAMATAN", "AKTIF",
			"SECURITY_Q1", "SECURITY_A1_HASH", "SECURITY_Q2", "SECURITY_A2_HASH"},
		{"anggota_ciputat", "hash_lama", "ANGGOTA", "Anggota Ciputat", "Ciputat", "Y",
			"q1", "a1", "q2", "a2"},
	}
}

func registration() model.Registration {
	return model.Registration{
		Username:  "anggota_baru",
		Password:  "rahasia123",
		Nama:      "Anggota Baru",
		Kecamatan: "Ciputat",
		Question1: model.SecurityQuestions[0],
		Answer1:   "Boneka",
		Question2: model.SecurityQuestions[1],
		Answer2:   "Bandung",
	}
}

func TestRegisterUserAppendsRowWithLegacyLayout(t *testing.T) {
	st, fake := newTestStore(t)
	fake.grid["USERS"] = usersGrid()

	if err := st.RegisterUser(registration()); err != nil {
		t.Fatalf("register: %v", err)
	}

	grid := fake.grid["USERS"]
	if len(grid) != 3 {
		t.Fatalf("baris: %d, ingin 3", len(grid))
	}
	row := grid[2]
	if len(row) != len(grid[0]) {
		t.Fatalf("lebar baris: %d, ingin %d", len(row), len(grid[0]))
	}
	if row[0] != "anggota_baru" {
		t.Errorf("username: %q", row[0])
	}
	if row[2] != "ANGGOTA" {
		t.Errorf("role: %q, ingin ANGGOTA", row[2])
	}
	if row[3] != "Anggota Baru" {
		t.Errorf("nama: %q", row[3])
	}
	if row[4] != "Ciputat" {
		t.Errorf("kecamatan: %q", row[4])
	}
	if row[5] != "Y" {
		t.Errorf("aktif: %q, ingin Y", row[5])
	}
}

// The password is hashed the same way the previous application hashed it:
// a plain SHA-256 digest, without a salt.
func TestRegisterUserHashesPasswordWithSHA256(t *testing.T) {
	st, fake := newTestStore(t)
	fake.grid["USERS"] = usersGrid()

	if err := st.RegisterUser(registration()); err != nil {
		t.Fatalf("register: %v", err)
	}

	sum := sha256.Sum256([]byte("rahasia123"))
	want := hex.EncodeToString(sum[:])
	if got := fake.grid["USERS"][2][1]; got != want {
		t.Errorf("hash password: %q, ingin %q", got, want)
	}
}

// The security answers are lowercased and squeezed before they are hashed,
// exactly as the previous application did.
func TestRegisterUserNormalizesSecurityAnswers(t *testing.T) {
	st, fake := newTestStore(t)
	fake.grid["USERS"] = usersGrid()

	reg := registration()
	reg.Answer1 = "  Boneka   Muda "
	reg.Answer2 = "BANDUNG"
	if err := st.RegisterUser(reg); err != nil {
		t.Fatalf("register: %v", err)
	}

	row := fake.grid["USERS"][2]
	sum := sha256.Sum256([]byte("boneka muda"))
	if got, want := row[7], hex.EncodeToString(sum[:]); got != want {
		t.Errorf("hash jawaban 1: %q, ingin %q", got, want)
	}
	sum = sha256.Sum256([]byte("bandung"))
	if got, want := row[9], hex.EncodeToString(sum[:]); got != want {
		t.Errorf("hash jawaban 2: %q, ingin %q", got, want)
	}
	if row[6] != reg.Question1 || row[8] != reg.Question2 {
		t.Errorf("pertanyaan keamanan: %q / %q", row[6], row[8])
	}
}

func TestRegisterUserRejectsTakenUsername(t *testing.T) {
	st, fake := newTestStore(t)
	fake.grid["USERS"] = usersGrid()

	reg := registration()
	reg.Username = "Anggota_Ciputat"
	err := st.RegisterUser(reg)
	if err == nil {
		t.Fatal("harus gagal: username sudah dipakai")
	}
	if !strings.Contains(err.Error(), "sudah digunakan") {
		t.Errorf("pesan: %v", err)
	}
	if len(fake.grid["USERS"]) != 2 {
		t.Errorf("baris: %d, ingin 2 (tidak ada baris baru)", len(fake.grid["USERS"]))
	}
}

// An account registered here can sign in through the legacy check, because
// the hash format is the same.
func TestRegisteredUserCanSignIn(t *testing.T) {
	st, fake := newTestStore(t)
	fake.grid["USERS"] = usersGrid()

	if err := st.RegisterUser(registration()); err != nil {
		t.Fatalf("register: %v", err)
	}

	user, ok, err := st.VerifyLegacy("anggota_baru", "rahasia123")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !ok {
		t.Fatal("akun baru harus bisa masuk")
	}
	if user.Username != "anggota_baru" {
		t.Errorf("username: %q", user.Username)
	}
	if user.Role != model.RolePetugas {
		t.Errorf("role: %q, ingin %q", user.Role, model.RolePetugas)
	}
}
