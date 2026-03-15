package zippy

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// copyZipFiles copies all files from the reader zip to the writer zip and
// returns a map of the existing files.
func copyZipFiles(reader *zip.ReadCloser, writer *zip.Writer) (map[string]*zip.File, error) {
	existingFiles := make(map[string]*zip.File, len(reader.File))
	for _, f := range reader.File {
		existingFiles[f.Name] = f

		if err := writer.Copy(f); err != nil {
			return nil, err
		}
	}

	return existingFiles, nil
}

// fileFound checks if a zip file matches any of the provided glob patterns.
func fileFound(zipFile *zip.File, files ...string) (bool, error) {
	for _, f := range files {
		match, err := filepath.Match(f, zipFile.Name)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}

// filterFiles filters the zip files based on the provided glob patterns. If no
// patterns are provided, all files are returned.
func filterFiles(zipFiles []*zip.File, files ...string) ([]*zip.File, error) {
	if zipFiles == nil {
		return nil, nil
	}

	// If we have files to extract, filter the files to extract
	if files == nil {
		return zipFiles, nil
	}

	extFiles := []*zip.File{}
	for _, file := range zipFiles {
		match, err := fileFound(file, files...)
		if err != nil {
			return nil, err
		}

		if match {
			extFiles = append(extFiles, file)
			continue
		}
	}

	return extFiles, nil
}

// openZipFile opens a zip file for reading and writing. If the file does not
// exist, it is created. If the file exists, it is opened for reading and
// and writing.
func openZipFile(path string) (*os.File, *zip.ReadCloser, error) {
	var zipFile *os.File
	var zipReadCloser *zip.ReadCloser

	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, err
	} else if os.IsNotExist(err) {
		zipFile, err = os.Create(path)
		if err != nil {
			return nil, nil, err
		}
	} else {
		zipFile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, nil, err
		}

		zipReadCloser, err = zip.OpenReader(path)
		if err != nil && !os.IsNotExist(err) {
			return nil, nil, err
		}
	}

	return zipFile, zipReadCloser, nil
}

// Removes the drive letter and colon from a Windows path.
func removeDriveLetter(path string) string {
	return strings.TrimPrefix(path, filepath.VolumeName(path))
}

// Converts a path to a zip-compatible path.
func toZipPath(path string) string {
	zipPath := removeDriveLetter(strings.ReplaceAll(path, string(filepath.Separator), "/"))
	return strings.TrimPrefix(zipPath, "/")
}

// Converts paths to zip-compatible paths
func toZipPaths(paths ...string) []string {
	for i, file := range paths {
		paths[i] = toZipPath(file)
	}

	return paths
}

// Validates the number of bytes written during a copy operation with
// expected number of bytes.
func validateCopy(path string, written int64, expected int64) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	if written != expected {
		return fmt.Errorf("failed to copy '%s': expected %d bytes, got %d bytes", absPath, expected, written)
	}

	return nil
}
