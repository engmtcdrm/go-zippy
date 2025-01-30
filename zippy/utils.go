package zippy

import (
	"fmt"
	"path/filepath"
	"strings"
)

// validateCopy validates the number of bytes written during a copy operation with
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

	if written != expected {
		return fmt.Errorf("failed to copy '%s': expected %d bytes, got %d bytes", absPath, expected, written)
	}

	return nil
}

// removeDriveLetter removes the drive letter and colon from a Windows path.
func removeDriveLetter(path string) string {
	return strings.TrimPrefix(strings.TrimPrefix(path, filepath.VolumeName(path)), "/")
}

// convertToZipPath removes the drive letter and colon from a Windows path and replaces
// backslashes with forward slashes.
func convertToZipPath(path string) string {
	return removeDriveLetter(strings.ReplaceAll(path, string(filepath.Separator), "/"))
}
