package flags

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

func TestAddConfigPath(t *testing.T) {
	// Test with absolute path
	absPath := "/absolute/path"
	AddConfigPath(absPath)
	if !slices.Contains(configPaths, filepath.Clean(absPath)) {
		t.Fatalf("Expected %s to be in configPaths", absPath)
	}

	// Test with relative path
	relPath := "relative/path"
	expectedPath, _ := filepath.Abs(relPath)
	AddConfigPath(relPath)
	if !slices.Contains(configPaths, filepath.Clean(expectedPath)) {
		t.Errorf("Expected %s to be in configPaths", expectedPath)
	}

	// Test with $HOME
	homePath := "$HOME/config"
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}
	expectedHomePath := filepath.Join(home, "config")
	AddConfigPath(homePath)
	if !slices.Contains(configPaths, filepath.Clean(expectedHomePath)) {
		t.Errorf("Expected %s to be in configPaths", expectedHomePath)
	}
}
