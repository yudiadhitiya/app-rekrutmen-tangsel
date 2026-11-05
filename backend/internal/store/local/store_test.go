package local

import (
	"path/filepath"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
	"rekrutmen-tangsel/backend/internal/store"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(filepath.Join(t.TempDir(), "rekrutmen.json"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func rec(noWA string) *model.Rekrutmen {
	return &model.Rekrutmen{
		ID:                noWA + "-id",
		Nama:              "Ahmad " + noWA,
		NoWA:              noWA,
		Kecamatan:         "Serpong",
		Kelurahan:         "Buaran",
		Aktivitas:         "Wiraswasta",
		TingkatPenerimaan: "Terbuka",
		Metode:            "Halaqoh",
		Status:            "Berhasil",
	}
}

// PRD: "Nomor Whatsapp bersifat unik untuk setiap orang".
func TestCreateRejectsDuplicateWhatsApp(t *testing.T) {
	s := newStore(t)
	if err := s.CreateRekrutmen(rec("628111111111")); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if err := s.CreateRekrutmen(rec("628111111111")); err != store.ErrDuplicateWA {
		t.Fatalf("expected ErrDuplicateWA, got %v", err)
	}
}

func TestUpdateRejectsDuplicateWhatsApp(t *testing.T) {
	s := newStore(t)
	if err := s.CreateRekrutmen(rec("628111111111")); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateRekrutmen(rec("628222222222")); err != nil {
		t.Fatal(err)
	}
	dup := rec("628222222222")
	dup.ID = "628111111111-id"
	if err := s.UpdateRekrutmen(dup); err != store.ErrDuplicateWA {
		t.Fatalf("expected ErrDuplicateWA, got %v", err)
	}
}

// Updating a record without changing its number must not trip the uniqueness check.
func TestUpdateWithSameWhatsAppSucceeds(t *testing.T) {
	s := newStore(t)
	r := rec("628111111111")
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatal(err)
	}
	r.Nama = "Ahmad Baru"
	if err := s.UpdateRekrutmen(r); err != nil {
		t.Fatalf("update with same wa failed: %v", err)
	}
	got, err := s.GetRekrutmen(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Nama != "Ahmad Baru" {
		t.Errorf("update not persisted, got %q", got.Nama)
	}
}

func TestListFilters(t *testing.T) {
	s := newStore(t)
	a := rec("628111111111")
	a.Nama = "Budi Santoso"
	a.Kecamatan = "Ciputat"
	a.Kelurahan = "Cipayung"
	a.Metode = "Program"
	a.Status = "Pending"
	if err := s.CreateRekrutmen(a); err != nil {
		t.Fatal(err)
	}
	b := rec("628222222222")
	b.Nama = "Siti Aminah"
	if err := s.CreateRekrutmen(b); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		f    model.Filter
		want int
	}{
		{"all", model.Filter{}, 2},
		{"by name", model.Filter{Q: "siti"}, 1},
		{"by wa", model.Filter{Q: "628111"}, 1},
		{"by kecamatan", model.Filter{Kecamatan: "Ciputat"}, 1},
		{"by kelurahan", model.Filter{Kelurahan: "Cipayung"}, 1},
		{"by metode", model.Filter{Metode: "Program"}, 1},
		{"by status", model.Filter{Status: "Pending"}, 1},
		{"no match", model.Filter{Q: "zzz"}, 0},
	}
	for _, c := range cases {
		data, total, err := s.ListRekrutmen(c.f)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if total != c.want {
			t.Errorf("%s: expected total %d, got %d", c.name, c.want, total)
		}
		if len(data) != c.want {
			t.Errorf("%s: expected %d rows, got %d", c.name, c.want, len(data))
		}
	}
}

func TestListPagination(t *testing.T) {
	s := newStore(t)
	for i := 0; i < 7; i++ {
		if err := s.CreateRekrutmen(rec("62811111111" + string(rune('0'+i)))); err != nil {
			t.Fatal(err)
		}
	}
	page1, total, err := s.ListRekrutmen(model.Filter{Page: 1, Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if total != 7 {
		t.Errorf("expected total 7, got %d", total)
	}
	if len(page1) != 3 {
		t.Errorf("expected 3 rows on page 1, got %d", len(page1))
	}
	page3, _, err := s.ListRekrutmen(model.Filter{Page: 3, Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(page3) != 1 {
		t.Errorf("expected 1 row on page 3, got %d", len(page3))
	}
	page9, _, err := s.ListRekrutmen(model.Filter{Page: 9, Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(page9) != 0 {
		t.Errorf("expected 0 rows past the end, got %d", len(page9))
	}
}

func TestFindByWhatsApp(t *testing.T) {
	s := newStore(t)
	r := rec("628111111111")
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatal(err)
	}
	got, err := s.FindRekrutmenByWA("628111111111")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != r.ID {
		t.Errorf("expected id %q, got %q", r.ID, got.ID)
	}
	if _, err := s.FindRekrutmenByWA("628999999999"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteRekrutmen(t *testing.T) {
	s := newStore(t)
	r := rec("628111111111")
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteRekrutmen(r.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetRekrutmen(r.ID); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

// PRD: "Terdapat master data Daftar Aktifitas/Pekerjaan yang dapat dikelola Admin".
func TestAktivitasMasterCRUD(t *testing.T) {
	s := newStore(t)
	a := &model.Aktivitas{ID: "a1", Nama: "Wiraswasta", Aktif: true}
	if err := s.CreateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListAktivitas()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Nama != "Wiraswasta" {
		t.Fatalf("unexpected list: %+v", list)
	}

	a.Nama = "Wiraswasta Tunggal"
	if err := s.UpdateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListAktivitas()
	if list[0].Nama != "Wiraswasta Tunggal" {
		t.Errorf("rename not applied: %q", list[0].Nama)
	}
}

// Renaming a master activity must rewrite the records that reference it.
func TestAktivitasRenameCascadesToRekrutmen(t *testing.T) {
	s := newStore(t)
	a := &model.Aktivitas{ID: "a1", Nama: "Wiraswasta", Aktif: true}
	if err := s.CreateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	r := rec("628111111111")
	r.Aktivitas = "Wiraswasta"
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatal(err)
	}

	a.Nama = "Wiraswasta Tunggal"
	if err := s.UpdateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetRekrutmen(r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Aktivitas != "Wiraswasta Tunggal" {
		t.Errorf("cascade rename failed, aktivitas=%q", got.Aktivitas)
	}
}

// An activity still in use must not be deleted.
func TestDeleteAktivitasInUseFails(t *testing.T) {
	s := newStore(t)
	a := &model.Aktivitas{ID: "a1", Nama: "Wiraswasta", Aktif: true}
	if err := s.CreateAktivitas(a); err != nil {
		t.Fatal(err)
	}
	r := rec("628111111111")
	r.Aktivitas = "Wiraswasta"
	if err := s.CreateRekrutmen(r); err != nil {
		t.Fatal(err)
	}
	if err := s.DeleteAktivitas("a1"); err != store.ErrAktivitasUsed {
		t.Fatalf("expected ErrAktivitasUsed, got %v", err)
	}
}

func TestCreateAktivitasRejectsDuplicateName(t *testing.T) {
	s := newStore(t)
	if err := s.CreateAktivitas(&model.Aktivitas{ID: "a1", Nama: "Pedagang"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateAktivitas(&model.Aktivitas{ID: "a2", Nama: "Pedagang"}); err != store.ErrDuplicateName {
		t.Fatalf("expected ErrDuplicateName, got %v", err)
	}
}

// Records must survive a restart.
func TestDataPersistsAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rekrutmen.json")
	s, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateRekrutmen(rec("628111111111")); err != nil {
		t.Fatal(err)
	}
	s2, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s2.GetRekrutmen("628111111111-id"); err != nil {
		t.Fatalf("record lost after reopen: %v", err)
	}
}
