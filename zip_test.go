package zippy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/go-zippy/testutils"
)

// func TestZipFile(t *testing.T) {
// 	tempDir, err := os.MkdirTemp("", "test-")
// 	if err != nil {
// 		t.Fatalf("Failed to create temp dir: %v", err)
// 	}
// 	defer func() {
// 		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
// 			t.Fatalf("Failed to remove temp dir: %v", removeErr)
// 		}
// 	}()

// 	// Create a test file to zip
// 	testFilePath := filepath.Join(tempDir, "testfile.txt")
// 	err = os.WriteFile(testFilePath, []byte("Hello, World!"), 0644)
// 	if err != nil {
// 		t.Fatalf("Failed to create test file: %v", err)
// 	}

// 	// Create a zip file
// 	zipFilePath := filepath.Join(tempDir, "test.zip")
// 	zFile, err := os.Create(zipFilePath)
// 	if err != nil {
// 		t.Fatalf("Failed to create zip file: %v", err)
// 	}
// 	defer func() {
// 		if closeErr := zFile.Close(); closeErr != nil {
// 			t.Fatalf("Failed to close file %v", closeErr)
// 		}
// 	}()

// 	zipWriter := zip.NewWriter(zFile)
// 	defer func() {
// 		if closeErr := zipWriter.Close(); closeErr != nil {
// 			t.Fatalf("Failed to close zip writer: %v", closeErr)
// 		}
// 	}()

// 	// Test zipFile function
// 	err = zipFile(zipWriter, testFilePath)
// 	if err != nil {
// 		t.Errorf("zipFile() error = %v", err)
// 	}
// }

func TestZippyAdd(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		// Change back to the original working directory
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("Failed to change back to original working directory: %v", err)
		}

		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			t.Fatalf("Failed to remove temp dir: %v", removeErr)
		}
	}()

	// Set the working directory to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

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
				if _, err = testutils.CreateTestFiles(tt.filePath, tt.files, tt.subfolders); err != nil {
					t.Fatalf("Failed to create test files: %v", err)
				}
			}

			var err error

			z := NewZippy(tt.dest)

			if tt.testName == "Bad Permissions Zip Input Path" {
				// err = testutils.PermissionTestVariadic(tt.filePath, z.Zip, tt.dest, tt.filePath)
				err = testutils.PermissionTest(tt.filePath, z.Add, tt.filePath)
			} else if tt.testName == "Bad Permissions Zip Output Path" {
				destDir := filepath.Dir(tt.dest)

				// err = testutils.PermissionTestVariadic(destDir, z.Zip, tt.dest, tt.filePath)
				err = testutils.PermissionTest(destDir, z.Add, tt.filePath)
			} else {
				err = z.Add(tt.filePath)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Zippy.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
