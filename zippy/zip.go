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
	// Add adds files or directories to a zip archive.
	//
	// files are the files or directories to archive. Glob patterns are supported.
	Add(files ...string) error

	// Delete deletes files or directories from an existing zip archive.
	//
	// files are the files or directories to delete. Glob patterns are supported.
	Delete(files ...string) error

	// Update updates files in a zip archive.
	//
	// files are the files or directories to update.
	Update(files ...string) error

	// Freshen updates files in a zip archive if the source file is newer.
	Freshen() error

	// Copy copies files from a zip archive to a destination directory.
	Copy(dest string, files ...string) error
}

type Zippy struct {
	Path          string // Path is the path to the zip archive.
	Junk          bool   // Junk specifies whether to junk the path when archiving.
	existingFiles map[string]*zip.File
	w             *zip.Writer
	rc            *zip.ReadCloser
}

func NewZippy(path string) *Zippy {
	return &Zippy{
		Path:          path,
		Junk:          false,
		existingFiles: make(map[string]*zip.File),
	}
}

// Add adds files or directories to a zip archive.
//
// files are the files or directories to archive. Glob patterns are supported.
func (z *Zippy) Add(files ...string) error {
	if err := os.MkdirAll(filepath.Dir(z.Path), os.ModePerm); err != nil {
		return err
	}

	var zipFile *os.File

	_, err := os.Stat(z.Path)
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

		z.rc, err = zip.OpenReader(z.Path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		defer func() {
			if closeErr := z.rc.Close(); closeErr != nil {
				err = fmt.Errorf("failed to close zip reader: %w", closeErr)
			}
		}()
	}

	z.w = zip.NewWriter(zipFile)
	defer func() {
		if closeErr := z.w.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	// Copy existing files to the new zip archive if zip file exists
	if err == nil && z.rc != nil {
		z.existingFiles = make(map[string]*zip.File)

		for _, f := range z.rc.File {
			z.existingFiles[f.Name] = f
		}

		for _, f := range z.rc.File {
			z.copyZipFile(f)
		}
	}

	if err := z.zipFiles(files...); err != nil {
		return err
	}

	return nil
}

// Delete deletes files or directories from an existing zip archive.
//
// files are the files or directories to delete. Glob patterns are supported.
func (z *Zippy) Delete(files ...string) error {
	_, err := os.Stat(z.Path)
	if err != nil {
		return err
	}

	z.rc, err = zip.OpenReader(z.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer func() {
		if closeErr := z.rc.Close(); closeErr != nil {
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

	newZipFile, err := os.Create(filepath.Join(tempDir, filepath.Base(z.Path)))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := newZipFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip file: %w", closeErr)
		}
	}()

	z.w = zip.NewWriter(newZipFile)
	defer func() {
		if closeErr := z.w.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zip writer: %w", closeErr)
		}
	}()

	for i, file := range files {
		files[i] = convertToZipPath(file)
	}

	// Copy existing files to the new zip archive if zip file exists
	if err == nil {
		if err := z.copyZipFilesRemove(z.rc.File, files); err != nil {
			return err
		}
	}

	if err := z.rc.Close(); err != nil {
		return err
	}

	if err := z.w.Close(); err != nil {
		return err
	}

	if err := newZipFile.Close(); err != nil {
		return err
	}

	// [ ] TODO: Handle cross-device move
	if err := os.Rename(newZipFile.Name(), z.Path); err != nil {
		return err
	}

	return err
}

// Update updates files in a zip archive.
func (z *Zippy) Update(files ...string) error {
	// Implementation of Update method
	return nil
}

// Freshen updates files in a zip archive if the source file is newer.
func (z *Zippy) Freshen() error {
	// Implementation of Freshen method
	return nil
}

// Copy copies files from existing zip archive to a new zip archive.
//
// dest is the new zip archive path.
// files are the files to copy. If no files are provided, all files will be copied.
func (z *Zippy) Copy(dest string, files ...string) error {
	// Implementation of Copy function
	return nil
}

// copyZipFile copies a file from a zip archive to another zip archive.
//
// zipWriter is the zip.Writer to write to.
//
// file is the file to copy.
func (z *Zippy) copyZipFile(file *zip.File) error {
	fWriter, err := z.w.CreateHeader(&file.FileHeader)
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

	err = validateCopy(file.Name, written, int64(file.UncompressedSize64))

	return err
}

// copyZipFilesRemove copies files from a zip archive to another zip archive, removing files that match the given patterns.
//
// zipWriter is the zip.Writer to write to.
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

// zipFile adds a file or directory to a zip archive.
//
// path is the file or directory to add.
func (z *Zippy) zipFile(path string) error {
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

	writer, err := z.w.CreateHeader(header)
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

// zipFiles adds files or directories to a zip archive.
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
