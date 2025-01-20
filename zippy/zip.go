package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func zipFile(path string, zipWriter *zip.Writer, baseDir string) error {
	relPath, err := filepath.Rel(baseDir, path)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = relPath
	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

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

	return validateCopy(path, written, info.Size())
}

// Zip compresses a file or directory to a zip archive.
// The destination file will be created if it does not exist.
//
// file is the file or directory to compress.
//
// dest is the destination zip archive path.
func Zip(file string, dest string) error {
	// Make sure the destination directory exists
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

	fInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	if fInfo.IsDir() {
		err = filepath.WalkDir(file, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			return zipFile(path, zipWriter, file)
		})
	} else {
		err = zipFile(file, zipWriter, filepath.Dir(file))
	}

	return err
}
