package zippy

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/zippy-tmp/zippy/testutils"
)

func TestUnzipFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test zip file
	zipFile, err := os.Create(filepath.Join(tempDir, "test.zip"))
	if err != nil {
		t.Fatalf("Failed to create test zip file: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)
	fileWriter, err := zipWriter.Create("testfile.txt")
	if err != nil {
		t.Fatalf("Failed to add file to test zip: %v", err)
	}
	_, err = fileWriter.Write([]byte("Hello, World!"))
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}
	zipWriter.Close()
	zipFile.Close()

	// Test unzipFile function
	zipReader, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		t.Fatalf("Failed to open test zip file: %v", err)
	}
	defer zipReader.Close()

	u := NewUnzippy(zipFile.Name())

	for _, file := range zipReader.File {
		err := u.unzipFile(file, filepath.Join(tempDir, "test_output", file.Name))
		if err != nil {
			t.Errorf("unzipFile() error = %v", err)
		}
	}
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

			uz := NewUnzippy(tt.filePath)

			var err error

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
