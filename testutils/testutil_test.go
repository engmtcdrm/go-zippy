package testutils

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [CreateTempFile] function.
func TestCreateTempFile(t *testing.T) {
	t.Run("create a 1 temp file", func(t *testing.T) {
		tempDir := t.TempDir()

		tempFile, err := CreateTempFile(tempDir, "testfile-*.txt")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)
	})

	t.Run("create empty name", func(t *testing.T) {
		tempDir := t.TempDir()

		tempFile, err := CreateTempFile(tempDir, "")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)
	})

	t.Run("error from os.CreateTemp", func(t *testing.T) {
		tempFile, err := CreateTempFile(os.DevNull, "testfile-*.txt")
		assert.Error(t, err)
		assert.Nil(t, tempFile)
	})
}

// Tests for [CreateTempFiles] function.
func TestCreateTestFiles(t *testing.T) {
	t.Run("create 3 temp files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles, err := CreateTempFiles(tempDir, 3)
		assert.NoError(t, err)
		assert.NotNil(t, testFiles)
		assert.Len(t, testFiles, 3)
	})

	t.Run("create -1 temp files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles, err := CreateTempFiles(tempDir, -1)
		assert.Error(t, err)
		assert.Nil(t, testFiles)
	})

	t.Run("error from CreateTempFile", func(t *testing.T) {
		testFiles, err := CreateTempFiles(os.DevNull, 1)
		assert.Error(t, err)
		assert.Nil(t, testFiles)
	})
}

// Tests for [CreateZipFile] function.
func TestCreateZipFile(t *testing.T) {
	tempDir := t.TempDir()
	zipFilePath := filepath.Join(tempDir, "test.zip")

	expectedCount, err := CreateZipFile(zipFilePath, 3, 2)
	assert.NoError(t, err)

	// Check if the zip file was created
	_, err = os.Stat(zipFilePath)
	assert.NoError(t, err)

	// Check the contents of the zip file
	zFile, err := zip.OpenReader(zipFilePath)
	assert.NoError(t, err)
	defer zFile.Close()

	// expectedFiles := 3
	expectedDirs := 2
	actualFiles := 0
	actualDirs := 0

	for _, file := range zFile.File {
		if file.FileInfo().IsDir() {
			actualDirs++
		} else {
			actualFiles++
		}
	}

	assert.Equal(t, expectedCount, actualFiles)
	assert.Equal(t, expectedDirs, actualDirs)
}
