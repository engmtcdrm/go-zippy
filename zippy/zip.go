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

func createHeader(path string, junkPath bool) (*zip.FileHeader, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return nil, err
	}

	header.Name = convertToZipPath(path)

	if junkPath {
		header.Name = filepath.Base(header.Name)
	}

	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	return header, nil
}

// zipFile adds a file or directory to a zip archive.
//
// zipWriter is the zip.Writer to write to.
//
// path is the file or directory to add.
func zipFile(zipWriter *zip.Writer, path string) error {
	header, err := createHeader(path, false)
	if err != nil {
		return err
	}

	writer, err := zipWriter.CreateHeader(header)
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
		return err
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
}

type Zippy struct {
	Path          string // Path is the path to the zip archive.
	Junk          bool   // Junk specifies whether to junk the path when archiving or extracting.
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

// Extract extracts the contents of a zip archive to the current working directory.
func (z *Zippy) Extract() error {
	// Implementation of Extract method
	return nil
}

// ExtractTo extracts the contents of a zip archive to a destination directory.
func (z *Zippy) ExtractTo(dest string) error {
	// Implementation of ExtractTo method
	return nil
}

// ExtractFile extracts a file from a zip archive.
func (z *Zippy) ExtractFiles(files ...string) error {
	// Implementation of ExtractFile method
	return nil
}

// ExtractFileTo extracts a file from a zip archive to a destination directory.
func (z *Zippy) ExtractFilesTo(dest string, files ...string) error {
	// Implementation of ExtractFileTo method
	return nil
}

func (z *Zippy) Contents() ([]*zip.File, error) {
	var err error

	if z.rc == nil {
		z.rc, err = zip.OpenReader(z.Path)
		if err != nil {
			return nil, err
		}
		defer func() {
			if closeErr := z.rc.Close(); closeErr != nil {
				err = closeErr
			}
		}()
	}

	return z.rc.File, err
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
