package main

import (
	"os"
	"path/filepath"
	"testing"
)

// writeFile writes a file with a known content.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestBackupAndRestore checks that a backup of the data directory can be put
// back, so the data of the application survives a failure of the disk.
func TestBackupAndRestore(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	outDir := filepath.Join(dir, "backups")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dataDir, "users.json"), `{"admin":"hash"}`)
	writeFile(t, filepath.Join(dataDir, "settings.json"), `{"tema":"terang"}`)

	runBackup([]string{"-data", dataDir, "-out", outDir, "-keep", "3"})

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	var archive string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".gz" {
			archive = filepath.Join(outDir, e.Name())
		}
	}
	if archive == "" {
		t.Fatal("arsip tidak dibuat")
	}
	if _, err := os.Stat(archive + ".sha256"); err != nil {
		t.Fatalf("checksum tidak dibuat: %v", err)
	}

	// The live data changes and then the backup is put back.
	writeFile(t, filepath.Join(dataDir, "users.json"), "rusak")
	writeFile(t, filepath.Join(dataDir, "settings.json"), "rusak")

	runRestore([]string{"-data", dataDir, "-file", archive, "-force"})

	b, err := os.ReadFile(filepath.Join(dataDir, "users.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"admin":"hash"}` {
		t.Errorf("isi users.json setelah restore: %q", string(b))
	}
}

// TestChecksumDetectsTamperedArchive checks that a changed archive is
// refused, so a broken backup is noticed before it is restored.
func TestChecksumDetectsTamperedArchive(t *testing.T) {
	dir := t.TempDir()
	dataDir := filepath.Join(dir, "data")
	outDir := filepath.Join(dir, "backups")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dataDir, "users.json"), "awal")
	runBackup([]string{"-data", dataDir, "-out", outDir})

	entries, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	var archive string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".gz" {
			archive = filepath.Join(outDir, e.Name())
		}
	}

	// The archive is changed after the checksum was written.
	b, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	b[len(b)-1] ^= 0xff
	writeFile(t, archive, string(b))

	if err := verifyChecksum(archive); err == nil {
		t.Error("arsip yang diubah seharusnya ditolak")
	}
}

// TestPruneKeepsLimit checks that only the newest backups are kept.
func TestPruneKeepsLimit(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"rekrutmen-20240101-000000.tar.gz", "rekrutmen-20240102-000000.tar.gz", "rekrutmen-20240103-000000.tar.gz"} {
		writeFile(t, filepath.Join(dir, name), "x")
		writeFile(t, filepath.Join(dir, name+".sha256"), "x")
	}
	if err := prune(dir, 2); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var archives int
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".gz" {
			archives++
		}
	}
	if archives != 2 {
		t.Errorf("jumlah arsip setelah prune: %d, ingin 2", archives)
	}
}
