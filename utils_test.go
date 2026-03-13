package zippy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [fileFound] function.
func Test_fileFound(t *testing.T) {
	zipFile := &zip.File{}
	zipFile.Name = "test.txt"

	t.Run("matching file found", func(t *testing.T) {
		match, err := fileFound(zipFile, "test.txt")
		assert.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("no matching file found", func(t *testing.T) {
		match, err := fileFound(zipFile, "other.txt")
		assert.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("bad glob pattern", func(t *testing.T) {
		_, err := fileFound(zipFile, "[")
		assert.Error(t, err)
	})
}

// Tests for [filterFiles] function.
func Test_filterFiles(t *testing.T) {
	zipFile1 := &zip.File{}
	zipFile1.Name = "test1.txt"
	zipFile2 := &zip.File{}
	zipFile2.Name = "test2.txt"
	zipFiles := []*zip.File{zipFile1, zipFile2}

	t.Run("filter with matching files", func(t *testing.T) {
		filteredFiles, err := filterFiles(zipFiles, "test1.txt")
		assert.NoError(t, err)
		assert.Len(t, filteredFiles, 1)
		assert.Equal(t, "test1.txt", filteredFiles[0].Name)
	})

	t.Run("zipFiles is nil", func(t *testing.T) {
		filteredFiles, err := filterFiles(nil, "test1.txt")
		assert.NoError(t, err)
		assert.Nil(t, filteredFiles)
	})

	t.Run("bad glob pattern", func(t *testing.T) {
		_, err := filterFiles(zipFiles, "[")
		assert.Error(t, err)
	})
}

// Tests for [removeDriveLetter] function.
func Test_removeDriveLetter(t *testing.T) {
	t.Run("valid Windows drive letter removal", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows-specific test on non-Windows OS")
		}

		validPath := "C:\\windows"
		expectedPath := "\\windows"

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

	t.Run("error from filepath.Abs", func(t *testing.T) {
		tempDir := t.TempDir()
		testTmpDir, err := os.MkdirTemp(tempDir, "test")
		assert.NoError(t, err)

		// Save current directory to restore later
		origDir, err := os.Getwd()
		assert.NoError(t, err)
		defer func() {
			_ = os.Chdir(origDir)
		}()

		// Change to the temp directory
		err = os.Chdir(testTmpDir)
		assert.NoError(t, err)

		// Remove the directory we're currently in
		err = os.Remove(testTmpDir)
		// Non-Windows systems should be able to remove the current directory
		// without error
		if runtime.GOOS != "windows" {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}

		err = validateCopy("relative/path", 100, 100)
		assert.Error(t, err)
	})
}
