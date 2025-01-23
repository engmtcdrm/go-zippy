package zippy

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"example.com/m/zippy/testutils"
)

func TestZipFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			t.Fatalf("Failed to remove temp dir: %v", removeErr)
		}
	}()

	// Create a test file to zip
	testFilePath := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFilePath, []byte("Hello, World!"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a zip file
	zipFilePath := filepath.Join(tempDir, "test.zip")
	zFile, err := os.Create(zipFilePath)
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}
	defer func() {
		if closeErr := zFile.Close(); closeErr != nil {
			t.Fatalf("Failed to close file %v", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(zFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			t.Fatalf("Failed to close zip writer: %v", closeErr)
		}
	}()

	// Test zipFile function
	err = zipFile(testFilePath, zipWriter, tempDir)
	if err != nil {
		t.Errorf("zipFile() error = %v", err)
	}
}

func TestZip(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			t.Fatalf("Failed to remove temp dir: %v", removeErr)
		}
	}()

	// Test zip of 0 files
	// Test zip of 1 file
	// Test zip of 10 files
	// Test zip of 0 files and 1 subdirectory
	// Test zip of 1 file and 1 subdirectory
	// Test zip of 10 files and 2 subdirectories
	// Test nonexistent zip input path
	// Test invalid (bad char \000) zip input path
	// Test bad permissions zip input path
	// Test bad permissions zip output path

	tests := []struct {
		testName   string
		exists     bool
		filePath   string
		dest       string
		files      int
		subfolders int
		wantErr    bool
	}{
		{"Zip 0 Files", true, filepath.Join(tempDir, "test0"), filepath.Join(tempDir, "test0.zip"), 0, 0, false},
		{"Zip 1 File", true, filepath.Join(tempDir, "test1"), filepath.Join(tempDir, "test1.zip"), 1, 0, false},
		{"Zip 10 Files", true, filepath.Join(tempDir, "test10"), filepath.Join(tempDir, "test10.zip"), 10, 0, false},
		{"Zip 0 Files and 1 Subdirectory", true, filepath.Join(tempDir, "test0sub"), filepath.Join(tempDir, "test0sub.zip"), 0, 1, false},
		{"Zip 1 File and 1 Subdirectory", true, filepath.Join(tempDir, "test1sub"), filepath.Join(tempDir, "test1sub.zip"), 1, 1, false},
		{"Zip 10 Files and 2 Subdirectories", true, filepath.Join(tempDir, "test10sub2"), filepath.Join(tempDir, "test10sub2.zip"), 10, 2, false},
		{"Nonexistent Zip Input Path", false, "nonexistent", filepath.Join(tempDir, "nonexistent.zip"), 0, 0, true},
		{"Invalid Zip Input Path", false, "/invalid/path\0001", filepath.Join(tempDir, "invalid.zip"), 0, 0, true},
		{"Bad Permissions Zip Input Path", true, filepath.Join(tempDir, "bad-in-perm"), filepath.Join(tempDir, "bad-in-perm.zip"), 0, 0, true},
		{"Bad Permissions Zip Output Path", true, filepath.Join(tempDir, "bad-out-perm"), filepath.Join(tempDir, "bad-out-perm", "test.zip"), 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.exists {
				if err = testutils.CreateTestFiles(tt.filePath, tt.files, tt.subfolders); err != nil {
					t.Fatalf("Failed to create test files: %v", err)
				}
			}

			if runtime.GOOS == "windows" {
				if tt.testName == "Bad Permissions Zip Input Path" {
					cmd := exec.Command("icacls", tt.filePath, "/deny", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
					if err := cmd.Run(); err != nil {
						t.Fatalf("Failed to set permissions: %v", err)
					}
					defer func() {
						// Restore permissions after the test
						cmd := exec.Command("icacls", tt.filePath, "/grant", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
						cmd.Run()
					}()
				}

				if tt.testName == "Bad Permissions Zip Output Path" {
					destDir := filepath.Dir(tt.dest)

					cmd := exec.Command("icacls", destDir, "/deny", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
					if err := cmd.Run(); err != nil {
						t.Fatalf("Failed to set permissions: %v", err)
					}
					defer func() {
						// Restore permissions after the test
						cmd := exec.Command("icacls", destDir, "/grant", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
						cmd.Run()
					}()
				}
			} else {
				if tt.testName == "Bad Permissions Zip Input Path" {
					if err := os.Chmod(tt.filePath, 0000); err != nil {
						t.Fatalf("Failed to set permissions: %v", err)
					}
					defer func() {
						// Restore permissions after the test
						if err := os.Chmod(tt.filePath, 0755); err != nil {
							t.Fatalf("Failed to restore permissions: %v", err)
						}
					}()
				}

				if tt.testName == "Bad Permissions Zip Output Path" {
					destDir := filepath.Dir(tt.dest)

					if err := os.Mkdir(destDir, 0755); err != nil {
						t.Fatalf("Failed to create test directory: %v", err)
					}

					if err := os.Chmod(destDir, 0000); err != nil {
						t.Fatalf("Failed to set permissions: %v", err)
					}
					defer func() {
						// Restore permissions after the test
						if err := os.Chmod(destDir, 0755); err != nil {
							t.Fatalf("Failed to restore permissions: %v", err)
						}
					}()
				}
			}

			err := Zip(tt.filePath, tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Zip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
