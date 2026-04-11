package testutils

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for [CreateTempFile] function.
func Test_CreateTempFile(t *testing.T) {
	t.Run("create a 1 temp file", func(t *testing.T) {
		tempDir := t.TempDir()

		tempFile, err := CreateTempFile(tempDir, "testfile-*.txt")
		require.NoError(t, err)
		require.NotNil(t, tempFile)
	})

	t.Run("create empty name", func(t *testing.T) {
		tempDir := t.TempDir()

		tempFile, err := CreateTempFile(tempDir, "")
		require.NoError(t, err)
		require.NotNil(t, tempFile)
	})

	t.Run("error from os.CreateTemp", func(t *testing.T) {
		tempFile, err := CreateTempFile(os.DevNull, "testfile-*.txt")
		require.Error(t, err)
		require.Nil(t, tempFile)
	})
}

// Tests for [CreateTempFiles] function.
func Test_CreateTempFiles(t *testing.T) {
	t.Run("create 3 temp files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles, err := CreateTempFiles(tempDir, 3)
		require.NoError(t, err)
		require.NotNil(t, testFiles)
		require.Len(t, testFiles, 3)

		for i, file := range testFiles {
			expectedName := fmt.Sprintf("test%d.txt", i)
			require.Equal(t, expectedName, filepath.Base(file.Name()))
		}
	})

	t.Run("create -1 temp files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles, err := CreateTempFiles(tempDir, -1)
		require.Error(t, err)
		require.Nil(t, testFiles)
	})

	t.Run("error from CreateTempFile", func(t *testing.T) {
		testFiles, err := CreateTempFiles(os.DevNull, 1)
		require.Error(t, err)
		require.Nil(t, testFiles)
	})
}

// Tests for [CreateRandomTempFiles] function.
func Test_CreateRandomTempFiles(t *testing.T) {
	t.Run("create 3 temp files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles, err := CreateRandomTempFiles(tempDir, 3)
		require.NoError(t, err)
		require.NotNil(t, testFiles)
		require.Len(t, testFiles, 3)
	})

	t.Run("create -1 temp files", func(t *testing.T) {
		tempDir := t.TempDir()

		testFiles, err := CreateRandomTempFiles(tempDir, -1)
		require.Error(t, err)
		require.Nil(t, testFiles)
	})

	t.Run("error from CreateTempFile", func(t *testing.T) {
		testFiles, err := CreateRandomTempFiles(os.DevNull, 1)
		require.Error(t, err)
		require.Nil(t, testFiles)
	})
}

// Tests for [CreateZipFileWithRandomFiles] function.
func Test_CreateZipFile(t *testing.T) {
	tempDir := t.TempDir()
	zipFilePath := filepath.Join(tempDir, "test.zip")

	expectedCount, err := CreateZipFileWithRandomFiles(zipFilePath, 3, 2)
	require.NoError(t, err)

	// Check if the zip file was created
	_, err = os.Stat(zipFilePath)
	require.NoError(t, err)

	// Check the contents of the zip file
	zFile, err := zip.OpenReader(zipFilePath)
	require.NoError(t, err)
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

	require.Equal(t, expectedCount, actualFiles)
	require.Equal(t, expectedDirs, actualDirs)
}
