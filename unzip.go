package zippy

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

type UnzippyInterface interface {
	Extract() ([]*zip.File, error)
	ExtractTo(dest string) ([]*zip.File, error)
	ExtractFiles(files ...string) ([]*zip.File, error)
	ExtractFilesTo(dest string, files ...string) ([]*zip.File, error)
}

type UnzippyOptions struct {
	Junk      bool // Junk specifies whether to junk the path of files when extracting.
	Overwrite bool // Overwrite specifies whether to overwrite files when extracting.
}

type Unzippy struct {
	Path    string          // Path to the zip archive.
	Options *UnzippyOptions // Options to use when extracting files.
}

// NewUnzippy creates a new Unzippy instance.
//
// path is the path to the zip archive.
func NewUnzippy(path string, options *UnzippyOptions) *Unzippy {
	if options == nil {
		options = &UnzippyOptions{}
	}
	return &Unzippy{
		Path:    path,
		Options: options,
	}
}

// unzipFile extracts a single file from a zip archive.
// The file is extracted to the specified path. The file
// is validated using the CRC32 checksum and the size
// of the extracted file is validated against the expected
// size.
//
// file is the file to extract.
//
// filePath is the path to extract the file to.
func (u *Unzippy) unzipFile(file *zip.File, filePath string) error {
	zippedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zippedFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zipped file: %w", closeErr)
		}
	}()

	// Ensure the directory for the inflated file exists
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return err
	}

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

	// Copy the zipped file to the output file and calculate the checksum
	// using a TeeReader to read from the zipped file and write to the hash
	// at the same time.
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

// Extract all files from zip archive to the same directory as the archive.
func (u *Unzippy) Extract() ([]*zip.File, error) {
	return u.ExtractFiles()
}

// Extracts all files from the zip archive to a destination directory.
// The destination directory will be created if it does not exist.
// The file modification times will be preserved.
//
// dest is the destination directory.
func (u *Unzippy) ExtractTo(dest string) ([]*zip.File, error) {
	return u.ExtractFilesTo(dest)
}

// Extracts the specified files from the zip archive.
//
// files to be extracted. If no files are specified, all files will be extracted.
// Glob patterns are supported.
func (u *Unzippy) ExtractFiles(files ...string) ([]*zip.File, error) {
	return u.ExtractFilesTo(filepath.Dir(u.Path), files...)
}

// Extracts the specified files from the zip archive to a destination directory.
// The destination directory will be created if it does not exist.
// The file modification times will be preserved.
//
// dest is the destination directory.
//
// files to be extracted. If no files are specified, all files will be extracted.
// Glob patterns are supported.
func (u *Unzippy) ExtractFilesTo(dest string, files ...string) ([]*zip.File, error) {
	var err error

	if err = os.MkdirAll(dest, os.ModePerm); err != nil {
		return nil, err
	}

	zipReader, err := zip.OpenReader(u.Path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := zipReader.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	extFiles := zipReader.File

	// If we have files to extract, filter the files to extract
	if files != nil {
		extFiles = []*zip.File{}
		for _, file := range zipReader.File {
			for _, f := range files {
				match, err := filepath.Match(f, file.Name)
				if err != nil {
					return nil, err
				}

				if match {
					extFiles = append(extFiles, file)
				}
			}
		}
	}

	for _, file := range extFiles {
		if u.Options.Junk {
			file.Name = filepath.Base(file.Name)
		}

		filePath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			if err := u.unzipFile(file, filePath); err != nil {
				return nil, err
			}
		}

		// Preserve the file modification date
		if err := os.Chtimes(filePath, file.Modified, file.Modified); err != nil {
			return nil, err
		}
	}

	return extFiles, err
}
