package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// zipFile adds a file or directory to a zip archive.
//
// path is the file or directory to add.
//
// zipWriter is the zip.Writer to write to.
//
// baseDir is the base directory for relative paths.
func zipFile(path string, zipWriter *zip.Writer, baseDir string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	fmt.Println("baseDir:", baseDir)

	var relPath string

	if filepath.IsAbs(path) {
		relPath = path
	} else {
		// Ensure both paths are absolute or relative
		if !filepath.IsAbs(path) {
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
		}
		if !filepath.IsAbs(baseDir) {
			baseDir, err = filepath.Abs(baseDir)
			if err != nil {
				return err
			}
		}

		relPath, err = filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		fmt.Println("Path:", path)
		fmt.Println("RelPath:", relPath)
	}

	// Zip format always uses forward slash
	header.Name = normalizePath(relPath)
	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	fmt.Println("Header.Name:", header.Name)
	fmt.Println()

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
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
func zipFiles(zipWriter *zip.Writer, file string) error {
	fInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println("CWD:", baseDir)

	if fInfo.IsDir() {
		err = filepath.WalkDir(file, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			return zipFile(path, zipWriter, baseDir)
		})
	} else {
		err = zipFile(file, zipWriter, baseDir)
	}

	return err
}

// Zip adds multiple files or directories to a zip archive.
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

	for _, file := range files {
		if err := zipFiles(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

// AddToZip adds multiple files or directories to an existing zip archive.
// The destination file will be created if it does not exist.
//
// dest is the destination zip archive path.
//
// files are the files or directories to compress.
func AddToZip(dest string, files ...string) error {
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
			// TODO: Put into own function for defer reasons
			fWriter, err := zipWriter.CreateHeader(&f.FileHeader)
			if err != nil {
				return err
			}
			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(fWriter, rc)
			rc.Close()
			if err != nil {
				return err
			}
		}
	}

	for _, file := range files {
		if err := zipFiles(zipWriter, file); err != nil {
			return err
		}
	}

	return nil
}

// DeleteFromZip deletes multiple files or directories from an existing zip archive.
//
// dest is the destination zip archive path.
//
// files are the files or directories to delete.
func DeleteFromZip(dest string, files ...string) error {
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

	filesToRemove := make(map[string]string)
	for _, file := range files {
		filesToRemove[file] = normalizePath(file)
	}

	if err == nil {
		for _, f := range zipReader.File {
			shouldRemove := false
			for fileToRemove := range filesToRemove {
				if strings.HasPrefix(f.Name, filesToRemove[fileToRemove]+"/") || f.Name == filesToRemove[fileToRemove] {
					shouldRemove = true
					break
				}
			}
			if shouldRemove {
				continue
			}

			// TODO: Put into own function for defer reasons
			fWriter, err := zipWriter.CreateHeader(&f.FileHeader)
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(fWriter, rc)
			rc.Close()
			if err != nil {
				return err
			}
		}
	}

	if err := zipReader.Close(); err != nil {
		return err
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	zipOpen, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zipOpen.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	_, err = io.Copy(newDestZipFile, zipOpen)
	if err != nil {
		return err
	}

	// if err := os.Rename(newDestZipFile.Name(), dest); err != nil {
	// 	return err
	// }

	return err
}
