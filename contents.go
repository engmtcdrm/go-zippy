package zippy

import (
	"archive/zip"
)

// Contents returns a list of files in the zip archive.
//
// zipFile is the path to the zip archive.
func Contents(zipFile string) ([]*zip.File, error) {
	zipRead, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer zipRead.Close()

	return zipRead.File, err
}
