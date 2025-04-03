package zippy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Validates the number of bytes written during a copy operation with
// expected number of bytes.
//
// path is the path of the file that was copied.
//
// written is the number of bytes written.
//
// expected is the expected number of bytes.
func validateCopy(path string, written int64, expected int64) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	if written != expected {
		return fmt.Errorf("failed to copy '%s': expected %d bytes, got %d bytes", absPath, expected, written)
	}

	return nil
}

// Removes the drive letter and colon from a Windows path.
func removeDriveLetter(path string) string {
	return strings.TrimPrefix(path, filepath.VolumeName(path))
}

// Converts a path to a zip-compatible path.
func toZipPath(path string) string {
	zipPath := removeDriveLetter(strings.ReplaceAll(path, string(filepath.Separator), "/"))
	return strings.TrimPrefix(zipPath, "/")
}

// Converts paths to zip-compatible paths
func toZipPaths(paths ...string) []string {
	for i, file := range paths {
		paths[i] = toZipPath(file)
	}

	return paths
}
