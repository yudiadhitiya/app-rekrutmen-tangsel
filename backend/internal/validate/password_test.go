package validate

import "testing"

// A password that protects the data of other people has to be long enough.
func TestPasswordNeedsTwelveCharacters(t *testing.T) {
	if f := Password("pendek1234"); len(f) == 0 {
		t.Error("password pendek harus ditolak")
	}
	if f := Password("panjang-sekarang-2024"); len(f) != 0 {
		t.Errorf("password panjang harus diterima, dapat %v", f)
	}
}

// A password that everybody knows is no password, whatever its length.
func TestPasswordRefusesCommonOnes(t *testing.T) {
	for _, weak := range []string{
		"admin12345678",
		"petugas123456",
		"password1234",
		"rahasia123456",
		"qwerty1234567",
		"123456789012",
		"aaaaaaaaaaaa",
	} {
		if f := Password(weak); len(f) == 0 {
			t.Errorf("password umum %q harus ditolak", weak)
		}
	}
}

// A password beyond the limit brings nothing but a slow hash.
func TestPasswordRefusesVeryLongOnes(t *testing.T) {
	long := ""
	for i := 0; i < 250; i++ {
		long += "x"
	}
	if f := Password(long); len(f) == 0 {
		t.Error("password kepanjangan harus ditolak")
	}
}

// The name of the account is known to everybody who knows the account, so it
// cannot be the password.
func TestPasswordDiffersFromUsername(t *testing.T) {
	if !PasswordDiffersFromUsername("anggota_baru", "anggota_baru") {
		t.Error("password yang sama dengan username harus ditolak")
	}
	if !PasswordDiffersFromUsername("Anggota_Baru", "anggota_baru") {
		t.Error("perbandingan harus mengabaikan huruf besar-kecil")
	}
	if PasswordDiffersFromUsername("rahasia-anggota-2024", "anggota_baru") {
		t.Error("password yang berbeda dari username harus diterima")
	}
}
