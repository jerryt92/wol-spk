package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type stringList []string

func (values *stringList) String() string { return strings.Join(*values, ",") }

func (values *stringList) Set(value string) error {
	*values = append(*values, filepath.ToSlash(value))
	return nil
}

func main() {
	var source, output string
	var gzipOutput bool
	var executables stringList
	var normalizeLF stringList

	flag.StringVar(&source, "source", "", "directory to archive")
	flag.StringVar(&output, "output", "", "output archive path")
	flag.BoolVar(&gzipOutput, "gzip", false, "gzip-compress the tar archive")
	flag.Var(&executables, "executable", "relative path that must have mode 0755 (repeatable)")
	flag.Var(&normalizeLF, "normalize-lf", "relative text path whose line endings are normalized to LF (repeatable)")
	flag.Parse()

	if source == "" || output == "" {
		fmt.Fprintln(os.Stderr, "usage: tarpack -source DIR -output FILE [-gzip] [-executable PATH]")
		os.Exit(2)
	}

	executableSet := make(map[string]bool, len(executables))
	for _, path := range executables {
		executableSet[path] = true
	}
	normalizeLFSet := make(map[string]bool, len(normalizeLF))
	for _, path := range normalizeLF {
		normalizeLFSet[path] = true
	}

	entries, err := archiveEntries(source)
	if err != nil {
		fatal(err)
	}

	file, err := os.Create(output)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	var writer io.Writer = file
	var gzipWriter *gzip.Writer
	if gzipOutput {
		gzipWriter = gzip.NewWriter(file)
		writer = gzipWriter
	}
	tarWriter := tar.NewWriter(writer)

	for _, entry := range entries {
		if err := writeEntry(tarWriter, source, entry, executableSet[entry], normalizeLFSet[entry]); err != nil {
			fatal(err)
		}
	}
	if err := tarWriter.Close(); err != nil {
		fatal(err)
	}
	if gzipWriter != nil {
		if err := gzipWriter.Close(); err != nil {
			fatal(err)
		}
	}
}

func archiveEntries(source string) ([]string, error) {
	var entries []string
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == source {
			return nil
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		entries = append(entries, filepath.ToSlash(relative))
		return nil
	})
	sort.Slice(entries, func(i, j int) bool {
		left, right := archiveEntryPriority(entries[i]), archiveEntryPriority(entries[j])
		if left != right {
			return left < right
		}
		return entries[i] < entries[j]
	})
	return entries, err
}

// DSM's package parser reads the outer SPK metadata before inspecting the
// remaining payload. Keep its required files at the beginning of the archive
// instead of relying on filesystem traversal or lexical ordering.
func archiveEntryPriority(entry string) int {
	switch {
	case entry == "INFO":
		return 0
	case entry == "package.tgz":
		return 1
	case entry == "scripts" || strings.HasPrefix(entry, "scripts/"):
		return 2
	case entry == "conf" || strings.HasPrefix(entry, "conf/"):
		return 3
	case entry == "PACKAGE_ICON.PNG" || entry == "PACKAGE_ICON_256.PNG":
		return 4
	default:
		return 5
	}
}

func writeEntry(writer *tar.Writer, source, relative string, executable, normalizeLF bool) error {
	path := filepath.Join(source, filepath.FromSlash(relative))
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	var normalizedContents []byte
	if normalizeLF && info.Mode().IsRegular() {
		normalizedContents, err = os.ReadFile(path)
		if err != nil {
			return err
		}
		normalizedContents = bytes.ReplaceAll(bytes.ReplaceAll(normalizedContents, []byte("\r\n"), []byte("\n")), []byte("\r"), []byte("\n"))
		header.Size = int64(len(normalizedContents))
	}
	header.Name = relative
	header.Format = tar.FormatUSTAR
	header.Uid = 0
	header.Gid = 0
	header.Uname = ""
	header.Gname = ""
	header.ModTime = time.Unix(0, 0).UTC()
	header.AccessTime = time.Time{}
	header.ChangeTime = time.Time{}
	if info.IsDir() {
		header.Name += "/"
		header.Mode = 0755
	} else if executable {
		header.Mode = 0755
	} else {
		header.Mode = 0644
	}
	if err := writer.WriteHeader(header); err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}
	if normalizedContents != nil {
		_, err := writer.Write(normalizedContents)
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(writer, file)
	return err
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
