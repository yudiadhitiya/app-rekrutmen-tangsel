package model

import "testing"

func TestLevelOrDefault(t *testing.T) {
	cases := []struct {
		name string
		user User
		want int
	}{
		{"level tersimpan dipakai", User{Role: RoleAdmin, Level: LevelAdminKecamatan}, LevelAdminKecamatan},
		{"admin tanpa level", User{Role: RoleAdmin}, LevelAdmin},
		{"petugas tanpa level", User{Role: RolePetugas}, LevelAnggota},
		{"tanpa role dan level", User{}, LevelAnggota},
	}
	for _, c := range cases {
		if got := c.user.LevelOrDefault(); got != c.want {
			t.Errorf("%s: harus %d, dapat %d", c.name, c.want, got)
		}
	}
}
