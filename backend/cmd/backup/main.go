// Command backup makes and restores a backup of the data directory of the
// application.
//
// The data directory holds the accounts, the settings and the session
// secret, so a backup of it is the only way to bring the application back
// after a failure of the disk.
//
// Usage:
//
//	backup backup  -data DATA_DIR -out BACKUP_DIR [-keep N]
//	backup restore -data DATA_DIR -file ARCHIVE [-force]
//	backup list    -out BACKUP_DIR
package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// defaultKeep is the number of backups that are kept when nothing else is
// asked for.
const defaultKeep = 14

func main() {
	log.SetFlags(0)
	log.SetPrefix("backup: ")

	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "backup":
		runBackup(os.Args[2:])
	case "restore":
		runRestore(os.Args[2:])
	case "list":
		runList(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		log.Fatalf("perintah tidak dikenal: %s", os.Args[1])
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `pemakaian:
  backup backup  -data DATA_DIR -out BACKUP_DIR [-keep N]
  backup restore -data DATA_DIR -file ARCHIVE [-force]
  backup list    -out BACKUP_DIR`)
	os.Exit(2)
}

// runBackup writes an archive of the data directory.
func runBackup(args []string) {
	fs := flag.NewFlagSet("backup", flag.ExitOnError)
	dataDir := fs.String("data", "data", "direktori data yang dicadangkan")
	outDir := fs.String("out", "backups", "direktori tempat arsip disimpan")
	keep := fs.Int("keep", defaultKeep, "jumlah arsip yang disimpan")
	_ = fs.Parse(args)

	files, err := listFiles(*dataDir)
	if err != nil {
		log.Fatalf("gagal membaca direktori data: %v", err)
	}
	if len(files) == 0 {
		log.Fatalf("direktori data %q kosong, tidak ada yang dicadangkan", *dataDir)
	}

	if err := os.MkdirAll(*outDir, 0o700); err != nil {
		log.Fatalf("gagal membuat direktori cadangan: %v", err)
	}

	name := fmt.Sprintf("rekrutmen-%s.tar.gz", time.Now().Format("20060102-150405"))
	path := filepath.Join(*outDir, name)
	if err := writeArchive(path, *dataDir, files); err != nil {
		log.Fatalf("gagal menulis arsip: %v", err)
	}
	sum, err := fileSHA256(path)
	if err != nil {
		log.Fatalf("gagal menghitung checksum: %v", err)
	}
	if err := os.WriteFile(path+".sha256", []byte(sum+"  "+name+"\n"), 0o600); err != nil {
		log.Fatalf("gagal menulis checksum: %v", err)
	}
	log.Printf("cadangan ditulis: %s (%d berkas, sha256 %s)", path, len(files), sum[:12])

	if *keep > 0 {
		if err := prune(*outDir, *keep); err != nil {
			log.Printf("gagal membersihkan cadangan lama: %v", err)
		}
	}
}

// runRestore puts the files of an archive back into the data directory.
func runRestore(args []string) {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	dataDir := fs.String("data", "data", "direktori data tujuan")
	file := fs.String("file", "", "arsip yang dipulihkan")
	force := fs.Bool("force", false, "timpa berkas yang sudah ada")
	_ = fs.Parse(args)

	if *file == "" {
		log.Fatalf("arsip belum diisi, pakai -file")
	}
	if _, err := os.Stat(*file); err != nil {
		log.Fatalf("arsip tidak ditemukan: %v", err)
	}
	if err := verifyChecksum(*file); err != nil {
		log.Fatalf("cadangan tidak dapat dipercaya: %v", err)
	}
	if err := os.MkdirAll(*dataDir, 0o700); err != nil {
		log.Fatalf("gagal membuat direktori data: %v", err)
	}

	// A restore overwrites the live data, so it refuses to do that without
	// being told to.
	if !*force {
		existing, err := listFiles(*dataDir)
		if err == nil && len(existing) > 0 {
			log.Fatalf("direktori data sudah berisi %d berkas; pakai -force untuk menimpanya", len(existing))
		}
	}

	written, err := readArchive(*file, *dataDir)
	if err != nil {
		log.Fatalf("gagal memulihkan arsip: %v", err)
	}
	log.Printf("dipulihkan %d berkas ke %s", written, *dataDir)
}

// runList shows the backups that exist.
func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	outDir := fs.String("out", "backups", "direktori cadangan")
	_ = fs.Parse(args)

	entries, err := os.ReadDir(*outDir)
	if err != nil {
		log.Fatalf("gagal membaca direktori cadangan: %v", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() > entries[j].Name() })
	found := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tar.gz") {
			continue
		}
		found = true
		info, err := e.Info()
		if err != nil {
			continue
		}
		fmt.Printf("%s  %s  %d byte\n", info.ModTime().Format("2006-01-02 15:04:05"), e.Name(), info.Size())
	}
	if !found {
		fmt.Println("belum ada cadangan")
	}
}

// listFiles returns the names of the regular files in a directory, without
// walking into subdirectories.
func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}
		out = append(out, e.Name())
	}
	return out, nil
}

// writeArchive writes the given files of a directory into a gzipped tar
// archive.
func writeArchive(path, dir string, names []string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(gz)
	for _, name := range names {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = name
		header.Uid, header.Gid = 0, 0
		header.Uname, header.Gname = "", ""
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		src, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, src); err != nil {
			src.Close()
			return err
		}
		src.Close()
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gz.Close()
}

// readArchive writes the files of an archive into a directory.
func readArchive(path, dir string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return 0, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	count := 0
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return count, err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		// The name of a file in an archive must not escape the target
		// directory.
		name := filepath.Clean(filepath.Base(header.Name))
		if name == "." || name == string(filepath.Separator) {
			continue
		}
		target := filepath.Join(dir, name)
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			return count, err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return count, err
		}
		out.Close()
		count++
	}
	return count, nil
}

// verifyChecksum compares the archive with the checksum file beside it.
func verifyChecksum(path string) error {
	sumFile := path + ".sha256"
	b, err := os.ReadFile(sumFile)
	if err != nil {
		if os.IsNotExist(err) {
			// A backup without a checksum file is old, so it is accepted
			// with a warning.
			log.Printf("peringatan: berkas checksum tidak ada, arsip tidak diverifikasi")
			return nil
		}
		return err
	}
	want := strings.Fields(string(b))
	if len(want) == 0 {
		return errors.New("berkas checksum kosong")
	}
	got, err := fileSHA256(path)
	if err != nil {
		return err
	}
	if !strings.EqualFold(want[0], got) {
		return fmt.Errorf("checksum tidak cocok: %s != %s", want[0], got)
	}
	return nil
}

// fileSHA256 returns the checksum of a file.
func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// prune removes the oldest archives, so the number of kept backups stays
// within the limit.
func prune(dir string, keep int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var archives []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tar.gz") {
			continue
		}
		archives = append(archives, e.Name())
	}
	if len(archives) <= keep {
		return nil
	}
	sort.Strings(archives)
	for _, name := range archives[:len(archives)-keep] {
		for _, suffix := range []string{"", ".sha256"} {
			if err := os.Remove(filepath.Join(dir, name+suffix)); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		log.Printf("arsip lama dihapus: %s", name)
	}
	return nil
}
