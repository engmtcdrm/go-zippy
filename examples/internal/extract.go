package internal

import (
	"fmt"
	"os"

	"github.com/engmtcdrm/go-zippy"
)

func ExtractExample() {
	tempDir, zipPath, err := createZipFile(10, 0)
	if err != nil {
		fmt.Printf("failed to create zip file: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	u, err := zippy.NewUnzippy(zipPath, nil)
	if err != nil {
		fmt.Printf("failed to create unzippy instance: %v\n", err)
		return
	}

	zFiles, err := u.Extract()
	if err != nil {
		fmt.Printf("failed to extract zip file: %v\n", err)
		return
	}

	displayExtractedFiles(tempDir, zFiles)
}
