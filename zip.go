package zippy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	Path          string // The path, including the file name, to the zip archive.
	Junk          bool   // Specifies whether to junk the path when archiving.
	tempFile      string // Temp file name when working with zip archives.
	existingFiles map[string]*zip.File
	zWriter       *zip.Writer
	zReadCloser   *zip.ReadCloser
}

func NewZippy(path string) *Zippy {
	path = filepath.Clean(path)

	return &Zippy{
		Path:          path,
		Junk:          false,
		tempFile:      "zippy-*",
		existingFiles: make(map[string]*zip.File),
	}
}

func (z *Zippy) createTempZipWithFiles(dest string, files ...string) (tempZipPath string, err error) {
	z.zReadCloser, err = zip.OpenReader(z.Path)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := z.zReadCloser.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close reader: %v", err, closeErr)
		}
	}()

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return "", err
	}

	// Create a temporary zip file in the same directory as Zippy.Path
	tempZipFile, err := os.CreateTemp(filepath.Dir(dest), z.tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary zip file: %w", err)
	}
	defer func() {
		closeErr := tempZipFile.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close temporary zip file: %v", err, closeErr)
		}
	}()

	// Copy entire zip file if no files are provided to copy
	if files == nil {
		return z.copyEntireZip(tempZipFile.Name())
	}

	z.zWriter = zip.NewWriter(tempZipFile)
	defer func() {
		closeErr := z.zWriter.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close zip writer: %v", err, closeErr)
		}
	}()

	files = toZipPaths(files...)

	// Copy existing files to the new zip archive, excluding the ones to delete
	if err := z.copyZipFilesKeep(z.zReadCloser.File, files); err != nil {
		return "", err
	}

	return tempZipFile.Name(), err
}

// Copy files from current zip to a temporary zip file removing any provided files listed
//
// files are the files or directories to delete. Glob patterns are supported.
//
// returns the path to the temporary zip file as well as any errors
func (z *Zippy) createTempZipWithoutFiles(files ...string) (tempZipPath string, err error) {
	z.zReadCloser, err = zip.OpenReader(z.Path)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := z.zReadCloser.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close zip reader: %v", err, closeErr)
		}
	}()

	// Create a temporary zip file in the same directory as Zippy.Path
	tempZipFile, err := os.CreateTemp(filepath.Dir(z.Path), z.tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary zip file: %w", err)
	}
	defer func() {
		closeErr := tempZipFile.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close temporary zip file: %v", err, closeErr)
		}
	}()

	z.zWriter = zip.NewWriter(tempZipFile)
	defer func() {
		closeErr := z.zWriter.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close zip writer: %v", err, closeErr)
		}
	}()

	files = toZipPaths(files...)

	// Copy existing files to the new zip archive, excluding the ones to delete
	if err := z.copyZipFilesRemove(z.zReadCloser.File, files); err != nil {
		return "", err
	}

	return tempZipFile.Name(), err
}

// copyEntireZip creates a copy of the entire zip file
//
// tempZipPath is the path to the temporary zip file to create
//
// returns the path to the temporary zip file and any errors
func (z *Zippy) copyEntireZip(tempZipPath string) (string, error) {
	var err error

	// Close any existing readers
	if z.zReadCloser != nil {
		closeErr := z.zReadCloser.Close()
		if closeErr != nil {
			return "", fmt.Errorf("failed to close zip reader: %w", closeErr)
		}
		z.zReadCloser = nil
	}

	fReader, err := os.Open(z.Path)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := fReader.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close file reader: %v", err, closeErr)
		}
	}()

	fWriter, err := os.Create(tempZipPath)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := fWriter.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close file writer: %v", err, closeErr)
		}
	}()

	_, err = io.Copy(fWriter, fReader)
	if err != nil {
		return "", err
	}

	return tempZipPath, err
}

// Copies files from a zip archive to another zip archive, removing files that match the given patterns.
//
// files are the files to copy.
//
// patterns are the patterns to match files to remove.
func (z *Zippy) copyZipFilesRemove(files []*zip.File, patterns []string) error {
	// First identify which files should be removed
	filesToRemove := make(map[string]bool)
	for _, pattern := range patterns {
		for _, file := range files {
			match, err := filepath.Match(pattern, file.Name)
			if err != nil {
				return err
			}

			if match {
				filesToRemove[file.Name] = true
			}
		}
	}

	// Create a map to track directories and their contents
	dirContents := make(map[string][]string)

	// Group files by their parent directories
	for _, file := range files {
		if !filesToRemove[file.Name] {
			// Use strings.Split/Join to get parent directory with forward slashes
			parts := strings.Split(file.Name, "/")
			if len(parts) > 1 {
				dir := strings.Join(parts[:len(parts)-1], "/")
				dirContents[dir] = append(dirContents[dir], file.Name)
			}
		}
	}

	// Identify empty directories that should be removed
	emptyDirs := make(map[string]bool)

	// Mark directories as empty if all their contents are marked for removal
	// This is a recursive process that works bottom-up
	changed := true
	for changed {
		changed = false
		for dir, contents := range dirContents {
			// Skip if already marked as empty
			if emptyDirs[dir] {
				continue
			}

			isEmpty := true
			for _, content := range contents {
				// If content is not marked for removal and is not an empty dir, dir is not empty
				if !filesToRemove[content] && !emptyDirs[content] {
					isEmpty = false
					break
				}
			}

			if isEmpty {
				emptyDirs[dir] = true
				changed = true

				// Mark this directory for removal in its parent directory's context
				// Use strings.Split/Join to get parent directory with forward slashes
				parts := strings.Split(dir, "/")
				if len(parts) > 1 {
					parentDir := strings.Join(parts[:len(parts)-1], "/")
					// Add this dir to its parent's contents list if not already there
					found := false
					for _, content := range dirContents[parentDir] {
						if content == dir || content == dir+"/" {
							found = true
							break
						}
					}

					if !found {
						dirContents[parentDir] = append(dirContents[parentDir], dir)
					}
				}
			}
		}
	}

	// Copy files that should not be removed
	for _, file := range files {
		// Skip if file is marked for removal
		if filesToRemove[file.Name] {
			continue
		}

		// Skip if file is a directory that is marked as empty
		if file.FileInfo().IsDir() {
			dirName := strings.TrimSuffix(file.Name, "/")
			if emptyDirs[dirName] {
				continue
			}
		}

		if err := z.zWriter.Copy(file); err != nil {
			return err
		}
	}

	return nil
}

// Keeps only the files that match the given patterns and copies them to another zip archive.
//
// files are the files to copy.
//
// patterns are the patterns to match files to keep.
func (z *Zippy) copyZipFilesKeep(files []*zip.File, patterns []string) error {
	// Map to track directories that need to be included
	dirsToInclude := make(map[string]bool)
	filesToCopy := make(map[string]*zip.File)

	// First pass: identify files to keep and their parent directories
	for _, file := range files {
		shouldKeep := false
		for _, fileToKeep := range patterns {
			match, err := filepath.Match(fileToKeep, file.Name)
			if err != nil {
				return err
			}

			if match {
				shouldKeep = true
				break
			}
		}

		if shouldKeep {
			filesToCopy[file.Name] = file

			// Mark all parent directories for inclusion using ZIP paths
			// Use strings.Split/Join instead of filepath.Dir to maintain forward slashes
			parts := strings.Split(file.Name, "/")
			for i := len(parts) - 1; i > 0; i-- {
				dir := strings.Join(parts[:i], "/")
				if dir != "" {
					dirsToInclude[dir] = true
				}
			}
		}
	}

	// Second pass: copy all directories first
	for _, file := range files {
		if file.FileInfo().IsDir() {
			// If this is a directory that needs to be included
			dirName := strings.TrimSuffix(file.Name, "/")
			if dirsToInclude[dirName] {
				if err := z.zWriter.Copy(file); err != nil {
					return err
				}
			}
		}
	}

	// Third pass: copy all files
	for _, file := range filesToCopy {
		if !file.FileInfo().IsDir() {
			if err := z.zWriter.Copy(file); err != nil {
				return err
			}
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

	header.Name = toZipPath(path)

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
		closeErr := file.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close file: %v", err, closeErr)
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
			closeErr := zipFile.Close()
			if err == nil {
				err = closeErr
			} else if closeErr != nil {
				err = fmt.Errorf("primary error: %w, additionally failed to close zip file: %v", err, closeErr)
			}
		}()
	} else {
		zipFile, err = os.OpenFile(z.Path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer func() {
			closeErr := zipFile.Close()
			if err == nil {
				err = closeErr
			} else if closeErr != nil {
				err = fmt.Errorf("primary error: %w, additionally failed to close zip file: %v", err, closeErr)
			}
		}()

		z.zReadCloser, err = zip.OpenReader(z.Path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		defer func() {
			if z.zReadCloser != nil {
				closeErr := z.zReadCloser.Close()
				if err == nil {
					err = closeErr
				} else if closeErr != nil {
					err = fmt.Errorf("primary error: %w, additionally failed to close zip reader: %v", err, closeErr)
				}
			}
		}()
	}

	z.zWriter = zip.NewWriter(zipFile)
	defer func() {
		closeErr := z.zWriter.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("primary error: %w, additionally failed to close zip writer: %v", err, closeErr)
		}
	}()

	// Copy existing files to the new zip archive if zip file exists
	if err == nil && z.zReadCloser != nil {
		z.existingFiles = make(map[string]*zip.File)

		for _, f := range z.zReadCloser.File {
			z.existingFiles[f.Name] = f
		}

		for _, f := range z.zReadCloser.File {
			if err := z.zWriter.Copy(f); err != nil {
				return err
			}
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
	tempZipPath, err := z.createTempZipWithoutFiles(files...)
	if err != nil {
		// If the temp zip file was made, but we had an error happen after the fact
		// lets clean it up if it exists
		_, err = os.Stat(tempZipPath)
		if err == nil {
			if err := os.Remove(tempZipPath); err != nil {
				return err
			}
		}

		return err
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
//
// files are the files to copy.  Glob patterns are supported. If no files are provided, all files will be copied.
func (z *Zippy) Copy(dest string, files ...string) (err error) {
	tempZipPath, err := z.createTempZipWithFiles(dest, files...)
	if err != nil {
		return err
	}

	// Rename the temporary zip file to the destination path
	if err := os.Rename(tempZipPath, dest); err != nil {
		return fmt.Errorf("failed to rename temporary zip file: %w", err)
	}

	return err
}
