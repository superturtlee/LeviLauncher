package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func LauncherDir() string {
	exe, err := os.Executable()
	if err != nil {
		cwd, _ := os.Getwd()
		return cwd
	}
	return filepath.Dir(exe)
}

// BaseRoot returns the base directory for storing launcher data.
// In CLI mode, uses a local ./data directory relative to the executable location.
func BaseRoot() string {
	// Use ./data directory relative to the launcher location
	root := filepath.Join(LauncherDir(), "data")
	_ = os.MkdirAll(root, 0o755)
	return root
}

func GetInstallerDir() (string, error) {
	base := BaseRoot()
	dir := filepath.Join(base, "installers")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
			return "", mkErr
		}
	}
	return dir, nil
}

func GetVersionsDir() (string, error) {
	base := BaseRoot()
	dir := filepath.Join(base, "versions")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
			return "", mkErr
		}
	}
	return dir, nil
}

func CanWriteDir(p string) bool {
	v := strings.TrimSpace(p)
	if v == "" {
		return false
	}
	if !filepath.IsAbs(v) {
		return false
	}
	if err := os.MkdirAll(v, 0o755); err != nil {
		return false
	}
	tf := filepath.Join(v, ".ll_write_test.tmp")
	f, err := os.Create(tf)
	if err != nil {
		return false
	}
	_, werr := f.Write([]byte("ok"))
	cerr := f.Close()
	_ = os.Remove(tf)
	return werr == nil && cerr == nil
}
