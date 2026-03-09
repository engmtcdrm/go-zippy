package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/go-zippy/testutils"
	"github.com/stretchr/testify/assert"
)

const testZipFileName = "test.zip"

// Tests [NewUnzippy] function.
func Test_NewUnzippy(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		u, err := NewUnzippy(testZipFileName, nil)
		assert.NoError(t, err)
		assert.Equal(t, testZipFileName, u.Path, fmt.Sprintf("Expected Path to be '%s'", testZipFileName))
		assert.NotNil(t, u.Options)
		assert.False(t, u.Options.Junk)
		assert.False(t, u.Options.Overwrite)
	})

	t.Run("with options", func(t *testing.T) {
		options := &UnzippyOptions{Junk: true, Overwrite: true}
		u, err := NewUnzippy(testZipFileName, options)
		assert.NoError(t, err)
		assert.Equal(t, testZipFileName, u.Path)
		assert.Equal(t, options, u.Options)
	})

	t.Run("empty path", func(t *testing.T) {
		u, err := NewUnzippy("", nil)
		assert.Error(t, err)
		assert.Nil(t, u)
	})

}

// Tests for [Unzippy.Extract] function.
func Test_Unzippy_Extract(t *testing.T) {
	t.Run("zip exists", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		err := testutils.CreateZipFile(zipFilePath, 10, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		_, err = u.Extract()
		assert.NoError(t, err)
	})
}

// Tests for [Unzippy.ExtractTo] function.
func Test_Unzippy_ExtractTo(t *testing.T) {
	t.Run("zip exists", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		err := testutils.CreateZipFile(zipFilePath, 10, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.NoError(t, err)
		assert.Len(t, files, 10)
	})

	t.Run("empty zip exists", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		err := testutils.CreateZipFile(zipFilePath, 0, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.NoError(t, err)
		assert.Len(t, files, 0)
	})

	t.Run("zip exists with subfolders", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		err := testutils.CreateZipFile(zipFilePath, 10, 2)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.NoError(t, err)
		assert.Len(t, files, 32)
	})

	t.Run("zip exists without files in subfolders", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		err := testutils.CreateZipFile(zipFilePath, 0, 2)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.NoError(t, err)
		assert.Len(t, files, 2)
	})

	t.Run("zip does not exist", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		dest := filepath.Join(tempDir, "output")

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.Error(t, err)
		assert.Nil(t, files)
	})

	t.Run("not a zip file", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, "not-a-zip.txt")
		dest := filepath.Join(tempDir, "output")

		err := os.WriteFile(zipFilePath, []byte("this is not a zip file"), 0644)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.Error(t, err)
		assert.Nil(t, files)
	})

	t.Run("bad zip file path", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := "/invalid/path\0001"
		dest := filepath.Join(tempDir, "output")

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		files, err := u.ExtractTo(dest)
		assert.Error(t, err)
		assert.Nil(t, files)
	})

	t.Run("bad permissions unzip input path", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, "bad-in-perm.zip")
		dest := filepath.Join(tempDir, "bad-in-perm")

		err := testutils.CreateZipFile(zipFilePath, 0, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		err = testutils.PermissionTest(zipFilePath, u.ExtractTo, dest)
		assert.Error(t, err)
	})

	t.Run("bad permissions unzip output path", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, "bad-out-perm.zip")
		dest := filepath.Join(tempDir, "bad-out-perm")

		err := testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		destDir := filepath.Dir(dest)
		err = testutils.PermissionTest(destDir, u.ExtractTo, dest)
		assert.Error(t, err)
	})
}

// [ ] TODO: Add tests for UnzipTo

// Tests for [Unzippy.copyAndValidate] function.
func Test_Unzippy_copyAndValidate(t *testing.T) {
	initUnzippy := func(t *testing.T, zipFilePath string) (*Unzippy, *zip.ReadCloser) {
		err := testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)

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
			assert.NoError(t, err)

			destFile, err := os.Create(destFilePath)
			assert.NoError(t, err)
			defer destFile.Close()

			zippedFileReader, err := file.Open()
			assert.NoError(t, err)
			defer zippedFileReader.Close()

			err = u.copyAndValidate(zippedFileReader, file, filepath.Dir(destFilePath), destFile)
			assert.NoError(t, err)
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
			assert.NoError(t, err)

			destFile, err := os.Create(destFilePath)
			assert.NoError(t, err)
			defer destFile.Close()

			zippedFileReader, err := file.Open()
			assert.NoError(t, err)
			defer zippedFileReader.Close()

			// Invalidate the zip file
			file.UncompressedSize64 = uint64(12345)

			err = u.copyAndValidate(zippedFileReader, file, filepath.Dir(destFilePath), destFile)
			assert.Error(t, err)
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
			assert.NoError(t, err)

			destFile, err := os.Create(destFilePath)
			assert.NoError(t, err)
			defer destFile.Close()

			zippedFileReader, err := file.Open()
			assert.NoError(t, err)
			defer zippedFileReader.Close()

			// Wrap the zippedFileReader with a custom reader that modifies the
			// data. This will allow io.Copy to be successful, but checksum
			// after to fail.
			corruptedReader := testutils.NewMockReader(zippedFileReader)

			err = u.copyAndValidate(corruptedReader, file, filepath.Dir(destFilePath), destFile)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "checksum") // Ensure the error is due to checksum mismatch
		}
	})
}

// Tests for [Unzippy.unzipFile] function.
func Test_Unzippy_unzipFile(t *testing.T) {
	t.Run("valid unzip file", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		err := testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			err := u.unzipFile(file, filepath.Join(tempDir, "test_output", file.Name))
			assert.NoError(t, err)
		}
	})

	t.Run("error from zip.File.Open due to bad Method", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		err := testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			// Set an invalid compression method to trigger an error
			file.Method = 54321
			err := u.unzipFile(file, filepath.Join(tempDir, "test_output", file.Name))
			assert.Error(t, err)
		}
	})

	t.Run("error from os.MkdirAll due to bad permissions", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		zipFileDirPath := filepath.Join(tempDir, "badperm")

		err := os.MkdirAll(zipFileDirPath, os.ModePerm)
		assert.NoError(t, err)

		err = testutils.CreateZipFile(zipFilePath, 1, 1)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			fileDest := filepath.Join(zipFileDirPath, "badsubperm", file.Name)
			err = testutils.PermissionTest(zipFileDirPath, u.unzipFile, file, fileDest)
			assert.Error(t, err)
		}
	})

	t.Run("error from os.OpenFile due to bad permissions", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)
		zipFileDirPath := filepath.Join(tempDir, "badperm")

		err := os.MkdirAll(zipFileDirPath, os.ModePerm)
		assert.NoError(t, err)

		err = testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			fileDest := filepath.Join(zipFileDirPath, file.Name)
			err = testutils.PermissionTest(zipFileDirPath, u.unzipFile, file, fileDest)
			assert.Error(t, err)
		}
	})
}

// mockZipFile is a custom struct that embeds zip.File and overrides the Open method.
type mockZipFile struct {
	*zip.File
}

// Open overrides the Open method to return a mock error.
func (m *mockZipFile) Open() (io.ReadCloser, error) {
	return nil, fmt.Errorf("mock open error")
}

// TODO: Tests for [Unzippy.unzipFiles] function.
func Test_Unzippy_unzipFiles(t *testing.T) {
	t.Run("valid unzip files", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		err := testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		err = u.unzipFiles(tempDir, zipReader.File...)
		assert.NoError(t, err)
	})

	t.Run("valid unzip files junked", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFileName)

		err := testutils.CreateZipFile(zipFilePath, 1, 0)
		assert.NoError(t, err)

		u, err := NewUnzippy(zipFilePath, &UnzippyOptions{Junk: true})
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		err = u.unzipFiles(tempDir, zipReader.File...)
		assert.NoError(t, err)
	})
}
