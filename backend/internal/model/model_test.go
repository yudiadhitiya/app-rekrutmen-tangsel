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

func TestScopeFor(t *testing.T) {
	cases := []struct {
		name string
		user User
		want Scope
	}{
		{"admin penuh", User{Username: "admin", Role: RoleAdmin, Level: LevelAdmin}, Scope{}},
		{"admin kecamatan", User{Username: "kec", Role: RoleAdmin, Level: LevelAdminKecamatan, Kecamatan: "Serpong"}, Scope{Kecamatan: "Serpong"}},
		{"admin kecamatan tanpa district", User{Username: "kec", Role: RoleAdmin, Level: LevelAdminKecamatan}, Scope{CreatedBy: "kec"}},
		{"anggota", User{Username: "petugas", Role: RolePetugas, Level: LevelAnggota}, Scope{CreatedBy: "petugas"}},
	}
	for _, c := range cases {
		if got := c.user.ScopeFor(); got != c.want {
			t.Errorf("%s: harus %+v, dapat %+v", c.name, c.want, got)
		}
	}
}

func TestMatchScope(t *testing.T) {
	rec := Rekrutmen{CreatedBy: "petugas", Kecamatan: "Serpong"}
	cases := []struct {
		name  string
		scope Scope
		want  bool
	}{
		{"tanpa scope", Scope{}, true},
		{"pemilik sama", Scope{CreatedBy: "petugas"}, true},
		{"pemilik lain", Scope{CreatedBy: "lain"}, false},
		{"district sama", Scope{Kecamatan: "Serpong"}, true},
		{"district lain", Scope{Kecamatan: "Ciputat"}, false},
		{"district sama beda huruf", Scope{Kecamatan: "serpong"}, true},
	}
	for _, c := range cases {
		if got := MatchScope(c.scope, rec); got != c.want {
			t.Errorf("%s: harus %v, dapat %v", c.name, c.want, got)
		}
	}
}
