package testutils

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// CreateTempFile creates a temporary file using [os.CreateTemp] in the given
// directory with the specified name pattern and writes the base name of the
// temporary file to its contents.
func CreateTempFile(dir, name string) (*os.File, error) {
	tempFile, err := os.CreateTemp(dir, name)
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
// directory.
func CreateTempFiles(dir string, files int) ([]*os.File, error) {
	if files < 0 {
		return nil, fmt.Errorf("invalid number of files: %d", files)
	}

	tempFiles := make([]*os.File, 0, files)

	for i := range files {
		file, err := CreateTempFile(dir, fmt.Sprintf("test%d-*.txt", i))
		if err != nil {
			return nil, err
		}

		tempFiles = append(tempFiles, file)
	}

	return tempFiles, nil
}

// CreateTempFilesInSubdirs creates the specified number of temporary files in
// the given directory and in the specified number of subdirectories.
func CreateTempFilesInSubdirs(dir string, files int, subdirs int) ([]*os.File, error) {
	totalFiles := files + (files * subdirs)
	tempFiles := make([]*os.File, 0, totalFiles)

	for i := range subdirs {
		subdirPath, err := os.MkdirTemp(dir, fmt.Sprintf("subfolder%d-*", i))
		if err != nil {
			return nil, err
		}

		slog.Debug(fmt.Sprintf("Created subdirectory %s", subdirPath))

		subTempFiles, err := CreateTempFiles(subdirPath, files)
		if err != nil {
			return nil, err
		}

		tempFiles = append(tempFiles, subTempFiles...)
	}

	return tempFiles, nil
}

// CreateZipFile creates a zip file with the specified number of files and
// subdirectories.
func CreateZipFile(zipFilePath string, files int, subdirs int) (int, error) {
	// Step 1: Create a temporary directory to hold the files
	tempDir, err := os.MkdirTemp("", "zip-temp-")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	expectedFiles := files + (files * subdirs)

	slog.Debug(fmt.Sprintf("Creating %d files/subdirectories in %s\n", expectedFiles+subdirs, tempDir))

	// Step 2: Create files and subdirectories in the temporary directory
	if _, err := CreateTempFiles(tempDir, files); err != nil {
		return 0, err
	}

	if _, err := CreateTempFilesInSubdirs(tempDir, files, subdirs); err != nil {
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
