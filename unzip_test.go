package zippy

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/go-zippy/testutils"
	"github.com/stretchr/testify/assert"
)

const testZipFile = "test.zip"

// Tests [NewUnzippy] function.
func Test_NewUnzippy(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		u, err := NewUnzippy(testZipFile, nil)
		assert.NoError(t, err)
		assert.Equal(t, testZipFile, u.Path, fmt.Sprintf("Expected Path to be '%s'", testZipFile))
		assert.NotNil(t, u.Options)
		assert.False(t, u.Options.Junk)
		assert.False(t, u.Options.Overwrite)
	})

	t.Run("with options", func(t *testing.T) {
		options := &UnzippyOptions{Junk: true, Overwrite: true}
		u, err := NewUnzippy(testZipFile, options)
		assert.NoError(t, err)
		assert.Equal(t, testZipFile, u.Path)
		assert.Equal(t, options, u.Options)
	})

	t.Run("empty path", func(t *testing.T) {
		u, err := NewUnzippy("", nil)
		assert.Error(t, err)
		assert.Nil(t, u)
	})

}

// Tests for [Unzippy.unzipFile] function.
func Test_Unzippy_unzipFile(t *testing.T) {
	t.Run("valid unzip file", func(t *testing.T) {
		tempDir := t.TempDir()
		defer os.RemoveAll(tempDir)
		zipFilePath := filepath.Join(tempDir, testZipFile)

		testCreateValidZipFile(t, zipFilePath)

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

	t.Run("unzip file with bad permissions when writing", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFile)
		zipFileDirPath := filepath.Join(tempDir, "badperm")
		err := os.MkdirAll(zipFileDirPath, 0000)
		assert.NoError(t, err)
		defer os.Chmod(zipFileDirPath, os.ModePerm)
		defer os.RemoveAll(zipFileDirPath)

		testCreateValidZipFile(t, zipFilePath)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			err := u.unzipFile(file, filepath.Join(zipFileDirPath, "badsubperm", file.Name))
			assert.Error(t, err)
		}
	})

	t.Run("unzip file with bad permissions when writing", func(t *testing.T) {
		tempDir := t.TempDir()
		zipFilePath := filepath.Join(tempDir, testZipFile)
		zipFileDirPath := filepath.Join(tempDir, "badperm")
		err := os.MkdirAll(zipFileDirPath, 0000)
		assert.NoError(t, err)
		defer os.Chmod(zipFileDirPath, os.ModePerm)
		defer os.RemoveAll(zipFileDirPath)

		testCreateValidZipFile(t, zipFilePath)

		u, err := NewUnzippy(zipFilePath, nil)
		assert.NoError(t, err)
		assert.NotNil(t, u)

		zipReader, err := zip.OpenReader(zipFilePath)
		assert.NoError(t, err)
		assert.NotNil(t, zipReader)
		defer zipReader.Close()

		for _, file := range zipReader.File {
			err := u.unzipFile(file, filepath.Join(zipFileDirPath, file.Name))
			assert.Error(t, err)
		}
	})
}

func TestUnzip(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			t.Fatalf("Failed to remove temp dir: %v", removeErr)
		}
	}()

	tests := []struct {
		testName   string
		exists     bool
		isFile     bool
		filePath   string
		dest       string
		files      int
		subfolders int
		wantErr    bool
	}{
		{"Zip Exists", true, false, filepath.Join(tempDir, "test.zip"), filepath.Join(tempDir, "output"), 10, 0, false},
		{"Empty Zip Exists", true, false, filepath.Join(tempDir, "test2.zip"), filepath.Join(tempDir, "output"), 0, 0, false},
		{"Zip Exists w Subfolders", true, false, filepath.Join(tempDir, "test3.zip"), filepath.Join(tempDir, "output"), 10, 2, false},
		{"Zip Exists wo Files w Subfolders", true, false, filepath.Join(tempDir, "test4.zip"), filepath.Join(tempDir, "output"), 0, 2, false},
		{"Zip Does Not Exist", false, false, "nonexistent.zip", filepath.Join(tempDir, "output"), 1, 0, true},
		{"Not a Zip File", false, true, filepath.Join(tempDir, "not_a_zip.txt"), filepath.Join(tempDir, "output"), 1, 0, true},
		{"Bad Zip File Path, Invalid Character", false, false, "/invalid/path\0001", filepath.Join(tempDir, "output"), 1, 0, true},
		{"Bad Permissions Unzip Input Path", true, false, filepath.Join(tempDir, "bad-in-perm.zip"), filepath.Join(tempDir, "bad-in-perm"), 0, 0, true},
		{"Bad Permissions Unzip Output Path", true, false, filepath.Join(tempDir, "bad-out-perm.zip"), filepath.Join(tempDir, "bad-out-perm"), 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			defer func() {
				files, err := os.ReadDir(tempDir)
				if err != nil {
					t.Fatalf("Failed to read tempDir: %v", err)
				}
				for _, file := range files {
					err := os.RemoveAll(filepath.Join(tempDir, file.Name()))
					if err != nil {
						t.Fatalf("Failed to remove file or directory in tempDir: %v", err)
					}
				}
			}()

			if tt.exists {
				if err = testutils.CreateZipFile(tt.filePath, tt.files, tt.subfolders); err != nil {
					t.Fatalf("Failed to create test zip file: %v", err)
				}
			} else if tt.isFile {
				file, err := os.Create(tt.filePath)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer func() {
					if closeErr := file.Close(); closeErr != nil {
						t.Fatalf("Failed to close test file: %v", closeErr)
					}
				}()
			}

			uz, err := NewUnzippy(tt.filePath, nil)
			assert.NoError(t, err)

			if tt.testName == "Bad Permissions Unzip Input Path" {
				err = testutils.PermissionTest(tt.filePath, uz.Extract)
			} else if tt.testName == "Bad Permissions Unzip Output Path" {
				destDir := filepath.Dir(tt.dest)
				err = testutils.PermissionTest(destDir, uz.ExtractTo, tt.dest)
			} else {
				_, err = uz.ExtractTo(tt.dest)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// [ ] TODO: Add tests for UnzipTo

func testCreateValidZipFile(t *testing.T, zipPath string) {
	zipFile, err := os.Create(zipPath)
	assert.NoError(t, err)
	assert.NotNil(t, zipFile)

	zipWriter := zip.NewWriter(zipFile)
	fileWriter, err := zipWriter.Create("testfile.txt")
	assert.NoError(t, err)

	_, err = fileWriter.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	zipWriter.Close()
	zipFile.Close()
}
