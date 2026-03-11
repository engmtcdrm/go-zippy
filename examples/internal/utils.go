package internal

import (
	"os"
	"path/filepath"

	"github.com/engmtcdrm/go-zippy/testutils"
)

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
