package zippy

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

// unzipFile extracts a single file from a zip archive.
// The file is extracted to the specified path. The file
// is validated using the CRC32 checksum and the size
// of the extracted file is validated against the expected
// size.
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

	hash := crc32.NewIEEE()

	written, err := io.Copy(outFile, io.TeeReader(zippedFile, hash))
	if err != nil {
		return err
	}

	checksum := hash.Sum32()

	if checksum != file.CRC32 {
		return fmt.Errorf("failed to copy '%s': expected '%08x' checksum, got '%08x' checksum", filePath, file.CRC32, checksum)
	}

	err = validateCopy(filePath, written, int64(file.UncompressedSize64))

	return err
}

type UnzippyInterface interface {
	// Extract extracts the contents of a zip archive to the current working directory.
	Extract() error

	// ExtractTo extracts the contents of a zip archive to a destination directory.
	//
	// dest is the destination directory to extract to.
	ExtractTo(dest string) error

	// ExtractFile extracts a file from a zip archive.
	//
	// file are the files to extract.
	ExtractFiles(files ...string) error

	// ExtractFileTo extracts a file from a zip archive to a destination directory.
	//
	// file are the files to extract.
	// dest is the destination directory to extract to.
	ExtractFilesTo(dest string, files ...string) error

	// Freshen existing files, i.e., extract only those files that already exist on disk and that are newer than the disk copies.
	Freshen() error
}

type Unzippy struct {
	Path string // Path is the path to the zip archive.
	Junk bool   // Junk specifies whether to junk the path when extracting.
}

func NewUnzippy(path string) *Unzippy {
	return &Unzippy{
		Path: path,
		Junk: false,
	}
}

// Extract the contents of a zip archive to the same directory as the archive.
func (uz *Unzippy) Extract() ([]*zip.File, error) {
	return uz.ExtractTo(filepath.Dir(uz.Path))
}

// ExtractTo extracts the contents of a zip archive to a destination directory.
// The destination directory will be created if it does not exist.
// The file modification times will be preserved.
//
// path is the path to the zip archive.
//
// dest is the destination directory.
func (uz *Unzippy) ExtractTo(dest string) ([]*zip.File, error) {
	zipReader, err := zip.OpenReader(uz.Path)
	if err != nil {
		return nil, err
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
				return nil, err
			}
		} else {
			// Ensure the directory for the inflated file exists
			fileDir := filepath.Dir(filePath)
			if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
				return nil, err
			}

			if err := unzipFile(file, filePath); err != nil {
				return nil, err
			}
		}

		// Preserve the file modification date
		if err := os.Chtimes(filePath, file.Modified, file.Modified); err != nil {
			return nil, err
		}
	}

	return zipReader.File, err
}

// ExtractFiles extracts the specified files from a zip archive.
//
// files are the files to extract.
func (uz *Unzippy) ExtractFiles(files ...string) error {
	return uz.ExtractFilesTo(filepath.Dir(uz.Path), files...)
}

// ExtractFilesTo extracts the specified files from a zip archive to a destination directory.
// The destination directory will be created if it does not exist.
// The file modification times will be preserved.
//
// files are the files to extract.
//
// dest is the destination directory.
func (uz *Unzippy) ExtractFilesTo(dest string, files ...string) error {
	// [ ] TODO: Implement ExtractFilesTo
	return nil
}

// Freshen existing files, i.e., extract only those files that already exist on disk and that are newer than the disk copies.
func (uz *Unzippy) Freshen() error {
	// [ ] TODO: Implement Freshen
	return nil
}
