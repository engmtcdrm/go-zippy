package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyZipFile(w *zip.Writer, f *zip.File) error {
	fWriter, err := w.CreateHeader(&f.FileHeader)
	if err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	written, err := io.Copy(fWriter, rc)
	if err != nil {
		return err
	}

	info := f.FileInfo()

	err = validateCopy(f.Name, written, info.Size())

	return err
}

// zipFile adds a file or directory to a zip archive.
//
// path is the file or directory to add.
//
// zipWriter is the zip.Writer to write to.
func zipFile(zipWriter *zip.Writer, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = convertToZipPath(path)

	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	if info.IsDir() {
		return nil
	}

	written, err := io.Copy(writer, file)
	if err != nil {
		return err
	}

	err = validateCopy(path, written, info.Size())

	return err
}

// zipFiles adds files or directories to a zip archive.
//
// zipWriter is the zip.Writer to write to.
//
// file is the file or directory to add.
func zipFiles(zipWriter *zip.Writer, files ...string) error {
	var err error

	for _, file := range files {
		fInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		// [ ] TODO: Add support for glob patterns
		if fInfo.IsDir() {
			err = filepath.WalkDir(file, func(path string, entry os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				return zipFile(zipWriter, path)
			})
		} else {
			err = zipFile(zipWriter, file)
		}

		if err != nil {
			return err
		}
	}

	return err
}

// Zip adds files or directories to a zip archive.
// The destination file will be created if it does not exist.
//
// dest is the destination zip archive path.
//
// files are the files or directories to compress.
func Zip(dest string, files ...string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	destZipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destZipFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(destZipFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	if err := zipFiles(zipWriter, files...); err != nil {
		return err
	}

	return nil
}

// Add adds files or directories to an existing zip archive.
// The destination file will be created if it does not exist.
//
// dest is the destination zip archive path.
//
// files are the files or directories to compress.
func Add(dest string, files ...string) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	destZipFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destZipFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	zipReader, err := zip.OpenReader(dest)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer func() {
		if closeErr := zipReader.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip reader: %w", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(destZipFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	// Copy existing files to the new zip archive if zip file exists
	if err == nil {
		for _, f := range zipReader.File {
			copyZipFile(zipWriter, f)
		}
	}

	if err := zipFiles(zipWriter, files...); err != nil {
		return err
	}

	return nil
}

// Delete deletes files or directories from an existing zip archive.
//
// dest is the destination zip archive path.
//
// files are the files or directories to delete. Glob patterns are supported.
func Delete(dest string, files ...string) error {
	_, err := os.Stat(dest)
	if err != nil {
		return err
	}

	zipReader, err := zip.OpenReader(dest)
	if err != nil && !os.IsNotExist(err) {

	}
	defer func() {
		if closeErr := zipReader.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip reader: %w", closeErr)
		}
	}()

	tempDir, err := os.MkdirTemp("", "zippy-")
	if err != nil {
		return err
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			err = fmt.Errorf("failed to remove temp dir: %w", removeErr)
		}
	}()

	newDestZipFile, err := os.Create(filepath.Join(tempDir, filepath.Base(dest)))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := newDestZipFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	zipWriter := zip.NewWriter(newDestZipFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	for i, file := range files {
		files[i] = convertToZipPath(file)
	}

	if err == nil {
		for _, f := range zipReader.File {
			shouldRemove := false
			for _, fileToRemove := range files {
				match, err := filepath.Match(fileToRemove, f.Name)
				if err != nil {
					return err
				}

				if match {
					shouldRemove = true
					break
				}
			}

			if shouldRemove {
				continue
			}

			copyZipFile(zipWriter, f)
		}
	}

	if err := zipReader.Close(); err != nil {
		return err
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	if err := newDestZipFile.Close(); err != nil {
		return err
	}

	// [ ] TODO: Handle cross-device move
	if err := os.Rename(newDestZipFile.Name(), dest); err != nil {
		return err
	}

	return err
}
