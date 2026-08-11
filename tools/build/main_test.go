package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageMetadataRejectsNonNumericVersion(t *testing.T) {
	directory := t.TempDir()
	infoPath := filepath.Join(directory, "INFO")

	if err := os.WriteFile(infoPath, []byte("package=\"WOLManager\"\nversion=\"1.0.2-beta\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, _, err := packageMetadata(infoPath)
	if err == nil || !strings.Contains(err.Error(), "numeric components") {
		t.Fatalf("packageMetadata() error = %v, want numeric version format error", err)
	}
}

func TestPackageMetadataAcceptsNumericVersion(t *testing.T) {
	directory := t.TempDir()
	infoPath := filepath.Join(directory, "INFO")

	if err := os.WriteFile(infoPath, []byte("package=\"WOLManager\"\nversion=\"1.0.2\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	packageName, version, err := packageMetadata(infoPath)
	if err != nil {
		t.Fatal(err)
	}
	if packageName != "WOLManager" || version != "1.0.2" {
		t.Fatalf("packageMetadata() = (%q, %q), want (%q, %q)", packageName, version, "WOLManager", "1.0.2")
	}
}
