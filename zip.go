package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZippyInterface defines the methods for working with zip archives.
type ZippyInterface interface {
	// Adds files or directories to a zip archive.
	//
	// files are the files or directories to archive. Glob patterns are supported.
	Add(files ...string) (err error)

	// Deletes files or directories from an existing zip archive.
	//
	// files are the files or directories to delete. Glob patterns are supported.
	Delete(files ...string) (err error)

	// Updates files in a zip archive.
	//
	// files are the files or directories to update. Glob patterns are supported.
	Update(files ...string) (err error)

	// Copies files from existing zip archive to a new zip archive.
	//
	// dest is the new zip archive path.
	//
	// files are the files to copy. Glob patterns are supported. If no files are provided, all files will be copied.
	Copy(dest string, files ...string) (err error)
}

type Zippy struct {
	Path          string // Path is the path to the zip archive.
	Junk          bool   // Junk specifies whether to junk the path when archiving.
	existingFiles map[string]*zip.File
	zWriter       *zip.Writer
	zReadCloser   *zip.ReadCloser
}

func NewZippy(path string) *Zippy {
	path = filepath.Clean(path)

	return &Zippy{
		Path:          path,
		Junk:          false,
		existingFiles: make(map[string]*zip.File),
	}
}

// Adds files or directories to a zip archive.
//
// files are the files or directories to archive. Glob patterns are supported.
func (z *Zippy) Add(files ...string) (err error) {
	if err := os.MkdirAll(filepath.Dir(z.Path), os.ModePerm); err != nil {
		return err
	}

	var zipFile *os.File

	_, err = os.Stat(z.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if os.IsNotExist(err) {
		zipFile, err = os.Create(z.Path)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := zipFile.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close zip file: %w", closeErr)
			}
		}()
	} else {
		zipFile, err = os.OpenFile(z.Path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := zipFile.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close zip file: %w", closeErr)
			}
		}()

		z.zReadCloser, err = zip.OpenReader(z.Path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		defer func() {
			if closeErr := z.zReadCloser.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close zip reader: %w", closeErr)
			}
		}()
	}

	z.zWriter = zip.NewWriter(zipFile)
	defer func() {
		if closeErr := z.zWriter.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	// Copy existing files to the new zip archive if zip file exists
	if err == nil && z.zReadCloser != nil {
		z.existingFiles = make(map[string]*zip.File)

		for _, f := range z.zReadCloser.File {
			z.existingFiles[f.Name] = f
		}

		for _, f := range z.zReadCloser.File {
			z.copyZipFile(f)
		}
	}

	if err := z.zipFiles(files...); err != nil {
		return err
	}

	return err
}

// Deletes files or directories from an existing zip archive.
//
// files are the files or directories to delete. Glob patterns are supported.
func (z *Zippy) Delete(files ...string) (err error) {
	_, err = os.Stat(z.Path)
	if err != nil {
		return err
	}

	z.zReadCloser, err = zip.OpenReader(z.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer func() {
		if z.zReadCloser != nil {
			if closeErr := z.zReadCloser.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close zip reader: %w", closeErr)
			}
		}
	}()

	// Create a temporary zip file in the same directory as Zippy.Path
	tempDir := filepath.Dir(z.Path)
	tempZipFile, err := os.CreateTemp(tempDir, ".zippy-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary zip file: %w", err)
	}
	defer func() {
		if tempZipFile != nil {
			if closeErr := tempZipFile.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close temporary zip file: %w", closeErr)
			}
		}
	}()

	tempZipPath := tempZipFile.Name()

	z.zWriter = zip.NewWriter(tempZipFile)
	defer func() {
		if z.zWriter != nil {
			if closeErr := z.zWriter.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close zip writer: %w", closeErr)
			}
		}
	}()

	// Convert file paths to zip-compatible paths
	for i, file := range files {
		files[i] = convertToZipPath(file)
	}

	// Copy existing files to the new zip archive, excluding the ones to delete
	if err := z.copyZipFilesRemove(z.zReadCloser.File, files); err != nil {
		return err
	}

	// Explicitly close resources before renaming the file
	if z.zWriter != nil {
		if closeErr := z.zWriter.Close(); closeErr != nil {
			return fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
		z.zWriter = nil
	}

	if tempZipFile != nil {
		if closeErr := tempZipFile.Close(); closeErr != nil {
			return fmt.Errorf("failed to close temporary zip file: %w", closeErr)
		}
		tempZipFile = nil
	}

	if z.zReadCloser != nil {
		if closeErr := z.zReadCloser.Close(); closeErr != nil {
			return fmt.Errorf("failed to close zip reader: %w", closeErr)
		}
		z.zReadCloser = nil
	}

	// Rename the temporary zip file to the original path
	if err := os.Rename(tempZipPath, z.Path); err != nil {
		return fmt.Errorf("failed to rename temporary zip file: %w", err)
	}

	return err
}

// Updates files in a zip archive.
//
// files are the files or directories to update.  Glob patterns are supported.
func (z *Zippy) Update(files ...string) (err error) {
	// TODO: Implementation of Update method
	return err
}

// Copies files from existing zip archive to a new zip archive.
//
// dest is the new zip archive path.
// files are the files to copy.  Glob patterns are supported. If no files are provided, all files will be copied.
func (z *Zippy) Copy(dest string, files ...string) (err error) {
	// TODO: Implementation of Copy function
	return err
}

// Copies a file from a zip archive to another zip archive.
//
// file is the file to copy.
func (z *Zippy) copyZipFile(file *zip.File) (err error) {
	fWriter, err := z.zWriter.CreateHeader(&file.FileHeader)
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

	_, err = io.Copy(fWriter, rc)
	if err != nil {
		return err
	}

	// err = validateCopy(file.Name, written, int64(file.UncompressedSize64))

	return err
}

// Copies files from a zip archive to another zip archive, removing files that match the given patterns.
//
// files are the files to copy.
//
// patterns are the patterns to match files to remove.
func (z *Zippy) copyZipFilesRemove(files []*zip.File, patterns []string) error {
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

		if err := z.copyZipFile(file); err != nil {
			return err
		}
	}

	return nil
}

// Adds a file or directory to a zip archive.
//
// path is the file or directory to add.
func (z *Zippy) zipFile(path string) error {
	path = filepath.Clean(path)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = convertToZipPath(path)

	if z.Junk {
		header.Name = filepath.Base(header.Name)
	}

	if header.FileInfo().IsDir() { // was info.IsDir()
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	// search z.existingFiles for matching header.Name
	// if found, skip
	if _, ok := z.existingFiles[header.Name]; ok {
		return nil
	}

	writer, err := z.zWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Open is done before checking if file is a directory to check permissions on the file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	if header.FileInfo().IsDir() {
		return nil
	}

	written, err := io.Copy(writer, file)
	if err != nil {
		return err
	}

	err = validateCopy(path, written, int64(header.UncompressedSize64))

	return err
}

// Adds files or directories to a zip archive.
//
// files are the files or directories to add. Glob patterns are supported.
func (z *Zippy) zipFiles(files ...string) error {
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
			fileMatch = filepath.Clean(fileMatch)

			fInfo, err := os.Stat(fileMatch)
			if err != nil {
				return err
			}

			if fInfo.IsDir() {
				err = filepath.WalkDir(fileMatch, func(path string, entry os.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}

					return z.zipFile(path)
				})
			} else {
				err = z.zipFile(fileMatch)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}
