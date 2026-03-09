package zippy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [Contents] function.
func TestContents(t *testing.T) {
	t.Run("Zip file exists", func(t *testing.T) {
		tempDir := t.TempDir()
		testZipFile := filepath.Join(tempDir, "test.zip")

		// Create a test zip file
		zipFile, err := os.Create(testZipFile)
		if err != nil {
			t.Fatalf("Failed to create test zip file: %v", err)
		}

		zipWriter := zip.NewWriter(zipFile)
		_, err = zipWriter.Create("testfile.txt")
		if err != nil {
			t.Fatalf("Failed to add file to test zip: %v", err)
		}
		zipWriter.Close()
		zipFile.Close()

		zipFiles, err := Contents(testZipFile)
		assert.Nil(t, err)
		assert.NotNil(t, zipFiles)
		assert.Len(t, zipFiles, 1)
	})

	t.Run("Zip file does not exist", func(t *testing.T) {
		zipFiles, err := Contents("nonexistent.zip")
		assert.NotNil(t, err)
		assert.Nil(t, zipFiles)
	})
}
