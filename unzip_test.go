package zippy

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/go-zippy/internal/testutils"
	"github.com/stretchr/testify/require"
)

const testZipFileName = "test.zip"

// Tests [NewUnzippy] function.
func Test_NewUnzippy(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		u, err := NewUnzippy(testZipFileName, nil)
		require.NoError(t, err)
		require.Equal(t, testZipFileName, u.Path, fmt.Sprintf("Expected Path to be '%s'", testZipFileName))
		require.NotNil(t, u.Options)
		require.False(t, u.Options.Junk)
		require.False(t, u.Options.Overwrite)
	})

	t.Run("with options", func(t *testing.T) {
		options := &UnzippyOptions{Junk: true, Overwrite: true}
		u, err := NewUnzippy(testZipFileName, options)
		require.NoError(t, err)
		require.Equal(t, testZipFileName, u.Path)
		require.Equal(t, options, u.Options)
	})

	t.Run("empty path", func(t *testing.T) {
		u, err := NewUnzippy("", nil)
		require.Error(t, err)
		require.Nil(t, u)
	})

}

// Tests for [Unzippy.Extract] and [Unzippy.ExtractFiles] function.
func Test_Unzippy_Extract_ExtractFiles(t *testing.T) {
	t.Run("zip exists", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 10, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		_, err = u.Extract()
		require.NoError(t, err)
	})
}

// Tests for [Unzip.ExtractFilesTo] function.
func Test_Unzippy_ExtractFilesTo(t *testing.T) {
	t.Run("valid zip", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 10, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		_, err = u.ExtractFilesTo(tempDir)
		require.NoError(t, err)
	})

	t.Run("invalid glob pattern", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 10, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		_, err = u.ExtractFilesTo(tempDir, "[")
		require.Error(t, err)
	})
}

// Tests for [Unzippy.ExtractTo] function.
func Test_Unzippy_ExtractTo(t *testing.T) {
	t.Run("zip exists", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 10, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.NoError(t, err)
		require.Len(t, files, 10)
	})

	t.Run("empty zip exists", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 0, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.NoError(t, err)
		require.Len(t, files, 0)
	})

	t.Run("zip exists with subfolders", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 10, 2)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.NoError(t, err)
		require.Len(t, files, 32)
	})

	t.Run("zip exists without files in subfolders", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 0, 2)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.NoError(t, err)
		require.Len(t, files, 2)
	})

	t.Run("zip does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.Error(t, err)
		require.Nil(t, files)
	})

	t.Run("not a zip file", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, "not-a-zip.txt")
		dest := filepath.Join(tempDir, "output")

		err := os.WriteFile(zipFilePath, []byte("this is not a zip file"), 0644)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.Error(t, err)
		require.Nil(t, files)
	})

	t.Run("bad zip file path", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := "/invalid/path\0001"
		dest := filepath.Join(tempDir, "output")

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		require.Error(t, err)
		require.Nil(t, files)
	})

	t.Run("bad permissions unzip input path", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, "bad-in-perm.zip")
		dest := filepath.Join(tempDir, "bad-in-perm")

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 0, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		err = testutils.PermissionTest(zipFilePath, u.ExtractTo, dest)
		require.Error(t, err)
	})

	t.Run("bad permissions unzip output path", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, "bad-out-perm.zip")
		dest := filepath.Join(tempDir, "bad-out-perm")

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		destDir := filepath.Dir(dest)
		err = testutils.PermissionTest(destDir, u.ExtractTo, dest)
		require.Error(t, err)
	})
}

// [ ] TODO: Add tests for UnzipTo

// Tests for [Unzippy.copyAndValidate] function.
func Test_Unzippy_copyAndValidate(t *testing.T) {
	initUnzippy := func(t *testing.T, zipFilePath string) (*Unzippy, *zip.ReadCloser) {
		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)

		return u, zipReader
	}

	t.Run("valid copy and validate", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		u, zipReader := initUnzippy(t, zipFilePath)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			destFilePath := filepath.Join(tempDir, "output", file.Name)
			err := os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm)
			require.NoError(t, err)

			destFile, err := os.Create(destFilePath)
			require.NoError(t, err)
			defer destFile.Close()

			zippedFileReader, err := file.Open()
			require.NoError(t, err)
			defer zippedFileReader.Close()

			err = u.copyAndValidate(zippedFileReader, file, filepath.Dir(destFilePath), destFile)
			require.NoError(t, err)
		}
	})

	t.Run("error from io.Copy", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		u, zipReader := initUnzippy(t, zipFilePath)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			destFilePath := filepath.Join(tempDir, "output", file.Name)
			err := os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm)
			require.NoError(t, err)

			destFile, err := os.Create(destFilePath)
			require.NoError(t, err)
			defer destFile.Close()

			zippedFileReader, err := file.Open()
			require.NoError(t, err)
			defer zippedFileReader.Close()

			// Invalidate the zip file
			file.UncompressedSize64 = uint64(12345)

			err = u.copyAndValidate(zippedFileReader, file, filepath.Dir(destFilePath), destFile)
			require.Error(t, err)
		}
	})

	t.Run("error from checksum mismatch", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		u, zipReader := initUnzippy(t, zipFilePath)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			destFilePath := filepath.Join(tempDir, "output", file.Name)
			err := os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm)
			require.NoError(t, err)

			destFile, err := os.Create(destFilePath)
			require.NoError(t, err)
			defer destFile.Close()

			zippedFileReader, err := file.Open()
			require.NoError(t, err)
			defer zippedFileReader.Close()

			// Wrap the zippedFileReader with a custom reader that modifies the
			// data. This will allow io.Copy to be successful, but checksum
			// after to fail.
			corruptedReader := testutils.NewMockReader(zippedFileReader)

			err = u.copyAndValidate(corruptedReader, file, filepath.Dir(destFilePath), destFile)
			require.Error(t, err)
			require.Contains(t, err.Error(), "checksum") // Ensure the error is due to checksum mismatch
		}
	})
}

// Tests for [Unzippy.unzipFile] function.
func Test_Unzippy_unzipFile(t *testing.T) {
	t.Run("valid unzip file", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			err := u.unzipFile(file, filepath.Join(tempDir, "test_output", file.Name))
			require.NoError(t, err)
		}
	})

	t.Run("error from zip.File.Open due to bad Method", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			// Set an invalid compression method to trigger an error
			file.Method = 54321
			err := u.unzipFile(file, filepath.Join(tempDir, "test_output", file.Name))
			require.Error(t, err)
		}
	})

	t.Run("error from os.MkdirAll due to bad permissions", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		zipFileDirPath := filepath.Join(tempDir, "badperm")

		err := os.MkdirAll(zipFileDirPath, os.ModePerm)
		require.NoError(t, err)

		_, err = testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 1)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			fileDest := filepath.Join(zipFileDirPath, "badsubperm", file.Name)
			err = testutils.PermissionTest(zipFileDirPath, u.unzipFile, file, fileDest)
			require.Error(t, err)
		}
	})

	t.Run("error from os.OpenFile due to bad permissions", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		zipFileDirPath := filepath.Join(tempDir, "badperm")

		err := os.MkdirAll(zipFileDirPath, os.ModePerm)
		require.NoError(t, err)

		_, err = testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			fileDest := filepath.Join(zipFileDirPath, file.Name)
			err = testutils.PermissionTest(zipFileDirPath, u.unzipFile, file, fileDest)
			require.Error(t, err)
		}
	})
}

// TODO: Tests for [Unzippy.unzipFiles] function.
func Test_Unzippy_unzipFiles(t *testing.T) {
	t.Run("valid unzip files", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		require.NoError(t, err)
		require.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)
		defer zipReader.Close()

		err = u.unzipFiles(tempDir, zipReader.File...)
		require.NoError(t, err)
	})

	t.Run("valid unzip files junked", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		_, err := testutils.CreateZipFileWithRandomFiles(zipFilePath, 1, 0)
		require.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, &UnzippyOptions{Junk: true})
		require.NoError(t, err)
		require.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		require.NoError(t, err)
		require.NotNil(t, zipReader)
		defer zipReader.Close()

		err = u.unzipFiles(tempDir, zipReader.File...)
		require.NoError(t, err)
	})
}
