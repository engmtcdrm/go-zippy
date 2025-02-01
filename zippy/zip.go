package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyZipFile copies a file from a zip archive to another zip archive.
//
// zipWriter is the zip.Writer to write to.
//
// file is the file to copy.
func copyZipFile(zipWriter *zip.Writer, file *zip.File) error {
	fWriter, err := zipWriter.CreateHeader(&file.FileHeader)
	if err != nil {
		return err
	}

	rc, err := file.Open()
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

	info := file.FileInfo()

	err = validateCopy(file.Name, written, info.Size())

	return err
}

// copyZipFiles copies files from a zip archive to another zip archive.
//
// zipWriter is the zip.Writer to write to.
//
// files are the files to copy.
func copyZipFiles(zipWriter *zip.Writer, files []*zip.File) error {
	for _, file := range files {
		if err := copyZipFile(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

// copyZipFilesRemove copies files from a zip archive to another zip archive, removing files that match the given patterns.
//
// zipWriter is the zip.Writer to write to.
//
// files are the files to copy.
//
// patterns are the patterns to match files to remove.
func copyZipFilesRemove(zipWriter *zip.Writer, files []*zip.File, patterns []string) error {
	for _, file := range files {
		shouldRemove := false
		for _, fileToRemove := range patterns {
			match, err := filepath.Match(fileToRemove, file.Name)
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

		if err := copyZipFile(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

// zipFile adds a file or directory to a zip archive.
//
// zipWriter is the zip.Writer to write to.
//
// path is the file or directory to add.
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
// files are the files or directories to add. Glob patterns are supported.
func zipFiles(zipWriter *zip.Writer, files ...string) error {
	for _, file := range files {
		fileMatches, err := filepath.Glob(file)
		if err != nil {
			return fmt.Errorf("failed to glob pattern '%s': %v", file, err)
		}

		// If no matches found, treat file as a literal path
		if len(fileMatches) == 0 {
			fileMatches = append(fileMatches, file)
		}

		for _, fileMatch := range fileMatches {
			fInfo, err := os.Stat(fileMatch)
			if err != nil {
				return err
			}

			if fInfo.IsDir() {
				err = filepath.WalkDir(fileMatch, func(path string, entry os.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}

					return zipFile(zipWriter, path)
				})
			} else {
				err = zipFile(zipWriter, fileMatch)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Zip adds files or directories to a zip archive.
// The destination file will be created if it does not exist.
//
// dest is the destination zip archive path.
//
// files are the files or directories to archived. Glob patterns are supported.
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
// files are the files or directories to archive. Glob patterns are supported.
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

	// Copy existing files to the new zip archive if zip file exists
	if err == nil {
		if err := copyZipFilesRemove(zipWriter, zipReader.File, files); err != nil {
			return err
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
