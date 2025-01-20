package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Zip compresses a file or directory to a zip archive.
// The destination file will be created if it does not exist.
//
// file is the file or directory to compress.
//
// dest is the destination zip archive path.
func Zip(file string, dest string) error {
	fInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	if fInfo.IsDir() {
		return zipDir(file, dest)
	}

	return zipFile(file, dest)
}

// TODO: Refactor the actual read/copy process to its own function

func zipFile(file string, dest string) error {
	// Make sure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zipFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(file)

	zipFileWriter, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := fileReader.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file reader: %w", closeErr)
		}
	}()

	written, err := io.Copy(zipFileWriter, fileReader)
	if err != nil {
		return err
	}

	return validateCopy(file, written, fileInfo.Size())
}

func zipDir(dir string, dest string) error {
	// Make sure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zipFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	err = filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		if entry.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath

		zipFile, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		written, err := io.Copy(zipFile, file)
		if err != nil {
			return err
		}

		return validateCopy(path, written, info.Size())
	})

	return err
}
