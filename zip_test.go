package zippy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/engmtcdrm/go-zippy/internal/testutils"
	"github.com/stretchr/testify/assert"
)

// Tests for [Zippy.Add] function.
func TestZippyAdd(t *testing.T) {
	// originalDir, err := os.Getwd()
	// if err != nil {
	// 	t.Fatalf("Failed to get current working directory: %v", err)
	// }

	// tempDir, err := os.MkdirTemp("", "test-")
	// if err != nil {
	// 	t.Fatalf("Failed to create temp dir: %v", err)
	// }
	// defer func() {
	// 	// Change back to the original working directory
	// 	os.Chdir(originalDir)
	// 	os.RemoveAll(tempDir)
	// }()

	// // Set the working directory to the temporary directory
	// if err := os.Chdir(tempDir); err != nil {
	// 	t.Fatalf("Failed to change working directory: %v", err)
	// }

	createFilesToZip := func(t *testing.T, files, subdirs int) (string, string, string) {
		tempDir := t.TempDir()
		zipFileName := testutils.CreateTempFilename("test-*.zip")
		zipFilePath := filepath.Join(tempDir, zipFileName)

		toCompressDir, err := os.MkdirTemp(tempDir, "to-compress-*")
		assert.NoError(t, err)

		baseFiles, err := testutils.CreateTempFiles(toCompressDir, files)
		assert.NoError(t, err)
		assert.Len(t, baseFiles, files)

		if subdirs > 0 {
			subdirFiles, err := testutils.CreateTempFilesInSubdirs(toCompressDir, files, subdirs)
			assert.NoError(t, err)
			_ = subdirFiles
		}

		return tempDir, zipFilePath, toCompressDir
	}

	calcExpectedCount := func(files, subdirs int) int {
		return files + subdirs + (files * subdirs) + 1 // + 1 for base directory
	}

	t.Run("zip 0 files", func(t *testing.T) {
		files := 0
		subdirs := 0
		expectedFilesCount := calcExpectedCount(files, subdirs)
		_, zipFilePath, toCompressDir := createFilesToZip(t, files, subdirs)

		z := NewZippy(zipFilePath)

		err := z.Add(toCompressDir)
		assert.NoError(t, err)

		zippedFiles, err := Contents(zipFilePath)
		assert.NoError(t, err)
		assert.Len(t, zippedFiles, expectedFilesCount)
	})

	t.Run("zip 1 file", func(t *testing.T) {
		files := 1
		subdirs := 0
		expectedFilesCount := calcExpectedCount(files, subdirs)
		_, zipFilePath, toCompressDir := createFilesToZip(t, files, subdirs)

		z := NewZippy(zipFilePath)

		err := z.Add(toCompressDir)
		assert.NoError(t, err)

		zippedFiles, err := Contents(zipFilePath)
		assert.NoError(t, err)
		assert.Len(t, zippedFiles, expectedFilesCount)
	})

	t.Run("zip 10 files", func(t *testing.T) {
		files := 10
		subdirs := 0
		expectedFilesCount := calcExpectedCount(files, subdirs)
		_, zipFilePath, toCompressDir := createFilesToZip(t, files, subdirs)

		z := NewZippy(zipFilePath)

		err := z.Add(toCompressDir)
		assert.NoError(t, err)

		zippedFiles, err := Contents(zipFilePath)
		assert.NoError(t, err)
		assert.Len(t, zippedFiles, expectedFilesCount)
	})

	t.Run("zip 0 files and 1 subdirectory", func(t *testing.T) {
		files := 0
		subdirs := 1
		expectedFilesCount := calcExpectedCount(files, subdirs)
		_, zipFilePath, toCompressDir := createFilesToZip(t, files, subdirs)

		z := NewZippy(zipFilePath)

		err := z.Add(toCompressDir)
		assert.NoError(t, err)

		zippedFiles, err := Contents(zipFilePath)
		assert.NoError(t, err)
		assert.Len(t, zippedFiles, expectedFilesCount)
	})

	t.Run("zip 10 files and 2 subdirectories", func(t *testing.T) {
		files := 10
		subdirs := 2
		expectedFilesCount := calcExpectedCount(files, subdirs)
		_, zipFilePath, toCompressDir := createFilesToZip(t, files, subdirs)

		z := NewZippy(zipFilePath)

		err := z.Add(toCompressDir)
		assert.NoError(t, err)

		zippedFiles, err := Contents(zipFilePath)
		assert.NoError(t, err)
		assert.Len(t, zippedFiles, expectedFilesCount)
	})

	// tests := []struct {
	// 	testName   string
	// 	exists     bool
	// 	filePath   string
	// 	dest       string
	// 	files      int
	// 	subfolders int
	// 	wantErr    bool
	// }{
	// 	// {"Zip 0 Files", true, filepath.Join(tempDir, "test0"), filepath.Join(tempDir, "test0.zip"), 0, 0, false},
	// 	// {"Zip 1 File", true, filepath.Join(tempDir, "test1"), filepath.Join(tempDir, "test1.zip"), 1, 0, false},
	// 	// {"Zip 10 Files", true, filepath.Join(tempDir, "test10"), filepath.Join(tempDir, "test10.zip"), 10, 0, false},
	// 	// {"Zip 0 Files and 1 Subdirectory", true, filepath.Join(tempDir, "test0sub"), filepath.Join(tempDir, "test0sub.zip"), 0, 1, false},
	// 	// {"Zip 1 File and 1 Subdirectory", true, filepath.Join(tempDir, "test1sub"), filepath.Join(tempDir, "test1sub.zip"), 1, 1, false},
	// 	// {"Zip 10 Files and 2 Subdirectories", true, filepath.Join(tempDir, "test10sub2"), filepath.Join(tempDir, "test10sub2.zip"), 10, 2, false},
	// 	{"Nonexistent Zip Input Path", false, "nonexistent", filepath.Join(tempDir, "nonexistent.zip"), 0, 0, true},
	// 	{"Invalid Zip Input Path", false, "/invalid/path\0001", filepath.Join(tempDir, "invalid.zip"), 0, 0, true},
	// 	{"Bad Permissions Zip Input Path", true, filepath.Join(tempDir, "bad-in-perm"), filepath.Join(tempDir, "bad-in-perm.zip"), 0, 0, true},
	// 	{"Bad Permissions Zip Output Path", true, filepath.Join(tempDir, "bad-out-perm"), filepath.Join(tempDir, "bad-out-perm", "test.zip"), 1, 0, true},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.testName, func(t *testing.T) {
	// 		if tt.exists {
	// 			if _, err = testutils.CreateTempFilesInSubdirs(tt.filePath, tt.files, tt.subfolders); err != nil {
	// 				t.Fatalf("Failed to create test files: %v", err)
	// 			}
	// 		}

	// 		var err error

	// 		z := NewZippy(tt.dest)

	// 		if tt.testName == "Bad Permissions Zip Input Path" {
	// 			// err = testutils.PermissionTestVariadic(tt.filePath, z.Zip, tt.dest, tt.filePath)
	// 			err = testutils.PermissionTest(tt.filePath, z.Add, tt.filePath)
	// 		} else if tt.testName == "Bad Permissions Zip Output Path" {
	// 			destDir := filepath.Dir(tt.dest)

	// 			// err = testutils.PermissionTestVariadic(destDir, z.Zip, tt.dest, tt.filePath)
	// 			err = testutils.PermissionTest(destDir, z.Add, tt.filePath)
	// 		} else {
	// 			err = z.Add(tt.filePath)
	// 		}

	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("Zippy.Add() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }
}

// func TestZippyDelete(t *testing.T) {
// 	originalDir, err := os.Getwd()
// 	if err != nil {
// 		t.Fatalf("Failed to get current working directory: %v", err)
// 	}

// 	tempDir, err := os.MkdirTemp("", "test-")
// 	if err != nil {
// 		t.Fatalf("Failed to create temp dir: %v", err)
// 	}
// 	defer func() {
// 		// Change back to the original working directory
// 		os.Chdir(originalDir)
// 		os.RemoveAll(tempDir)
// 	}()

// 	// Set the working directory to the temporary directory
// 	if err := os.Chdir(tempDir); err != nil {
// 		t.Fatalf("Failed to change working directory: %v", err)
// 	}

// 	tests := []struct {
// 		testName   string
// 		exists     bool
// 		filePath   string
// 		dest       string
// 		files      int
// 		subfolders int
// 		deleteGlob string
// 		wantErr    bool
// 	}{
// 		{"Delete 1 File", true, filepath.Join(tempDir, "test1"), filepath.Join(tempDir, "test1.zip"), 1, 0, "file0.txt", false},
// 		{"Delete 10 Files", true, filepath.Join(tempDir, "test10"), filepath.Join(tempDir, "test10.zip"), 10, 0, "file*.txt", false},
// 		{"Delete 1 Subdirectory", true, filepath.Join(tempDir, "test1sub"), filepath.Join(tempDir, "test1sub.zip"), 1, 1, "subdir0/", false},
// 		{"Delete All Files and Subdirectories", true, filepath.Join(tempDir, "test10sub2"), filepath.Join(tempDir, "test10sub2.zip"), 10, 2, "*", false},
// 		{"Delete from Nonexistent Zip", false, "nonexistent", filepath.Join(tempDir, "nonexistent.zip"), 0, 0, "*", true},
// 		{"Delete from Empty Zip", true, filepath.Join(tempDir, "empty"), filepath.Join(tempDir, "empty.zip"), 0, 0, "*", false},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.testName, func(t *testing.T) {
// 			if tt.exists {
// 				if _, err = testutils.CreateTempFilesInSubdirs(tt.filePath, tt.files, tt.subfolders); err != nil {
// 					t.Fatalf("Failed to create test files: %v", err)
// 				}

// 				z := NewZippy(tt.dest)
// 				if err := z.Add(tt.filePath); err != nil {
// 					t.Fatalf("Failed to add files to zip: %v", err)
// 				}
// 			}

// 			z := NewZippy(tt.dest)
// 			err := z.Delete(tt.deleteGlob)

// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Zippy.Delete() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			// Verify the zip file contents if no error is expected
// 			if !tt.wantErr {
// 				z.zReadCloser, err = zip.OpenReader(tt.dest)
// 				if err != nil {
// 					t.Fatalf("Failed to open zip file: %v", err)
// 				}
// 				defer z.zReadCloser.Close()

// 				for _, f := range z.zReadCloser.File {
// 					match, _ := filepath.Match(tt.deleteGlob, f.Name)
// 					if match {
// 						t.Errorf("File %s was not deleted as expected", f.Name)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// TODO: Add Tests for Zippy.Update

// TODO: Add Tests for Zippy.Copy
