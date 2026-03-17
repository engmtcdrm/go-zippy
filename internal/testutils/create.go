package testutils

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// CreateTempFile creates a temporary file using [os.CreateTemp] in the given
// directory with the specified name pattern if it contains a wildcard "*". If
// the name pattern does not contain a wildcard, a temporary file is created
// using [os.OpenFile] with the exact name. The function writes the base name of
// the temporary file to its contents and returns the created temporary file.
func CreateTempFile(dir, name string) (*os.File, error) {
	var tempFile *os.File
	var err error

	if strings.Contains(name, "*") {
		tempFile, err = os.CreateTemp(dir, name)
	} else {
		tempFilePath := filepath.Join(dir, name)
		tempFile, err = os.OpenFile(tempFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	}

	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	_, err = tempFile.Write([]byte(filepath.Base(tempFile.Name())))
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("Created file %s", tempFile.Name()))

	return tempFile, err
}

// CreateTempFilename creates a temporary filename based on the given pattern.
// The last wildcard "*" in the pattern will be replaced with random digits.
func CreateTempFilename(pattern string) string {
	prefix, suffix := prefixAndSuffix(pattern)
	return prefix + genRandomDigits(10) + suffix
}

// CreateTempFiles creates the specified number of temporary files in the given
// directory with non-randomized names (e.g., test0.txt, test1.txt, etc.).
func CreateTempFiles(dir string, files int) ([]*os.File, error) {
	return createTempFiles(dir, files, "test%d.txt")
}

// CreateRandomTempFiles creates the specified number of temporary files in the given
// directory with randomized names (e.g., test0-1234567890.txt,
// test1-1829304829.txt, etc.).
func CreateRandomTempFiles(dir string, files int) ([]*os.File, error) {
	return createTempFiles(dir, files, "test%d-*.txt")
}

// CreateRandomTempFilesInSubdirs creates the specified number of temporary
// files in the given directory and in the specified number of subdirectories.
// File names are randomized (e.g., test0-1234567890.txt, test1-1829304829.txt,
// etc.).
func CreateRandomTempFilesInSubdirs(dir string, files int, subdirs int) ([]*os.File, error) {
	totalFiles := files + (files * subdirs)
	tempFiles := make([]*os.File, 0, totalFiles)

	for i := range subdirs {
		subdirPath, err := os.MkdirTemp(dir, fmt.Sprintf("subdir%d-*", i))
		if err != nil {
			return nil, err
		}

		slog.Debug(fmt.Sprintf("Created subdirectory %s", subdirPath))

		subTempFiles, err := CreateRandomTempFiles(subdirPath, files)
		if err != nil {
			return nil, err
		}

		tempFiles = append(tempFiles, subTempFiles...)
	}

	return tempFiles, nil
}

// CreateZipFileWithRandomFiles creates a zip file with the specified number of files and
// subdirectories.
func CreateZipFileWithRandomFiles(zipFilePath string, files int, subdirs int) (int, error) {
	// Step 1: Create a temporary directory to hold the files
	tempDir, err := os.MkdirTemp("", "zip-temp-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	expectedFiles := files + (files * subdirs)

	slog.Debug(fmt.Sprintf("Creating %d files/subdirectories in %s\n", expectedFiles+subdirs, tempDir))

	// Step 2: Create files and subdirectories in the temporary directory
	if _, err := CreateRandomTempFiles(tempDir, files); err != nil {
		return 0, err
	}

	if _, err := CreateRandomTempFilesInSubdirs(tempDir, files, subdirs); err != nil {
		return 0, err
	}

	// Step 3: Create the zip file
	zFile, err := os.Create(zipFilePath)
	if err != nil {
		return 0, err
	}
	defer zFile.Close()

	zWrite := zip.NewWriter(zFile)
	defer zWrite.Close()

	// Step 4: Walk through the temporary directory and add files to the zip archive
	err = filepath.Walk(tempDir, addFilesToZip(tempDir, zWrite))

	slog.Debug("")
	slog.Debug(fmt.Sprintf("Created zip file %s", zipFilePath))

	return expectedFiles, err
}

// addFilesToZip returns a filepath.WalkFunc that adds files and directories
// from the specified tempDir to the provided zip.Writer.
func addFilesToZip(tempDir string, zWrite *zip.Writer) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == tempDir {
			return nil
		}

		// Get the relative path to maintain the directory structure in the zip archive
		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Add a directory entry to the zip archive
			_, err := zWrite.Create(relPath + "/")
			return err
		}

		// Add a file entry to the zip archive
		zipFile, err := zWrite.Create(relPath)
		if err != nil {
			return err
		}

		// Open the file and copy its contents to the zip archive
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(zipFile, srcFile)
		return err
	}
}

// createTempFiles is a helper function that creates the specified number of
// temporary files in the given directory with names based on the provided
// format string. If the format string contains a wildcard "*", it will be
// replaced with random digits.
func createTempFiles(dir string, files int, format string) ([]*os.File, error) {
	if files < 0 {
		return nil, fmt.Errorf("invalid number of files: %d", files)
	}

	tempFiles := make([]*os.File, 0, files)

	for i := range files {
		file, err := CreateTempFile(dir, fmt.Sprintf(format, i))
		if err != nil {
			return nil, err
		}

		tempFiles = append(tempFiles, file)
	}

	return tempFiles, nil
}
