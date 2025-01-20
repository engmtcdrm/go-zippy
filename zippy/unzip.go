package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// unzipFile extracts a single file from a zip archive.
// The file is extracted to the specified path.
//
// file is the file to extract.
//
// filePath is the path to extract the file to.
func unzipFile(file *zip.File, filePath string) error {
	zippedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zippedFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zipped file: %w", closeErr)
		}
	}()

	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	written, err := io.Copy(outFile, zippedFile)
	if err != nil {
		return err
	}

	return validateCopy(filePath, written, int64(file.UncompressedSize64))
}

// Unzip extracts the contents of a zip archive to a destination directory.
// The destination directory will be created if it does not exist.
// The file modification times will be preserved.
//
// zipFilePath is the path to the zip archive.
//
// dest is the destination directory.
func Unzip(zipFilePath string, dest string) error {
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zipReader.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	for _, file := range zipReader.File {
		filePath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return err
			}
		} else {
			// Ensure the directory for the inflated file exists
			fileDir := filepath.Dir(filePath)
			if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
				return err
			}

			if err := unzipFile(file, filePath); err != nil {
				return err
			}
		}

		// Preserve the file modification date
		if err := os.Chtimes(filePath, file.Modified, file.Modified); err != nil {
			return err
		}
	}

	return nil
}
