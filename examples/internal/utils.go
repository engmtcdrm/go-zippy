package internal

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/go-zippy/testutils"
)

// createZipFile creates a temporary zip file with the specified number of files
// and subdirectories.
func createZipFile(files, subdirs int) (string, string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), tempDirName)
	if err != nil {
		return "", "", err
	}

	zipPath := filepath.Join(tempDir, zipFileName)

	_, err = testutils.CreateZipFile(zipPath, files, subdirs)
	if err != nil {
		return "", "", err
	}

	return tempDir, zipPath, nil
}

// displayExtractedFiles prints the list of extracted files along with their
// destination paths.
func displayExtractedFiles(dest string, files []*zip.File) {
	fmt.Println()
	fmt.Println("Extracted files:")
	fmt.Println()

	pad := 0

	for _, file := range files {
		if len(file.Name) > pad {
			pad = len(file.Name)
		}
	}

	for _, file := range files {
		fmt.Printf(
			"Archive path:\t%s%s\tExtracted to:\t%s\n",
			file.Name,
			strings.Repeat(" ", pad-len(file.Name)),
			filepath.Join(dest, file.Name),
		)
	}

	fmt.Println()
}
