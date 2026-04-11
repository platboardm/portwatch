package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestVersionFlag builds the binary and checks --version output.
func TestVersionFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build test in short mode")
	}

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "portwatch")

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	var buf bytes.Buffer
	cmd := exec.Command(binPath, "--version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("--version exited with error: %v", err)
	}

	got := buf.String()
	if got == "" {
		t.Error("expected non-empty version output")
	}
}

// TestMissingConfig verifies a non-zero exit when config file is absent.
func TestMissingConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build test in short mode")
	}

	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "portwatch")

	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(binPath, "--config", filepath.Join(tmpDir, "nonexistent.yaml"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config, got nil")
	}
}
