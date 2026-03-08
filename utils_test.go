package zippy

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [removeDriveLetter] function.
func Test_removeDriveLetter(t *testing.T) {
	t.Run("valid Windows drive letter removal", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows-specific test on non-Windows OS")
		}

		validPath := "C:\\windows"
		expectedPath := "windows"

		processedPath := removeDriveLetter(validPath)
		assert.Equal(t, expectedPath, processedPath)
	})

	t.Run("valid non-Windows path", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping non-Windows-specific test on Windows OS")
		}

		validPath := "/usr/local"
		expectedPath := "/usr/local"

		processedPath := removeDriveLetter(validPath)
		assert.Equal(t, expectedPath, processedPath)
	})

	t.Run("empty path", func(t *testing.T) {
		processedPath := removeDriveLetter("")
		assert.Equal(t, "", processedPath)
	})
}

// Tests for [validateCopy] function.
func Test_validateCopy(t *testing.T) {
	makeValidFile := func(t *testing.T) string {
		validPath := filepath.Join(t.TempDir(), "valid")
		err := os.WriteFile(validPath, []byte("Test File"), os.ModePerm)
		assert.NoError(t, err)

		return validPath
	}

	t.Run("valid copy", func(t *testing.T) {
		validPath := makeValidFile(t)

		err := validateCopy(validPath, 100, 100)
		assert.NoError(t, err)
	})

	t.Run("invalid path", func(t *testing.T) {
		invalidPath := "/invalid/path\0001"
		err := validateCopy(invalidPath, 0, 0)
		assert.Error(t, err)
	})

	t.Run("mismatched bytes", func(t *testing.T) {
		validPath := makeValidFile(t)

		err := validateCopy(validPath, 100, 200)
		assert.Error(t, err)
	})
}
