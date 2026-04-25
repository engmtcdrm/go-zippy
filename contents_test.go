package zippy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for [Contents] function.
func Test_Contents(t *testing.T) {
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
		require.Nil(t, err)
		require.NotNil(t, zipFiles)
		require.Len(t, zipFiles, 1)
	})

	t.Run("Zip file does not exist", func(t *testing.T) {
		zipFiles, err := Contents("nonexistent.zip")
		require.NotNil(t, err)
		require.Nil(t, zipFiles)
	})
}
