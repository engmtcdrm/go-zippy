package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/engmtcdrm/go-zippy"
)

func ExtractTo() {
	tempDir, zipPath, err := createZipFile(10, 1)
	if err != nil {
		fmt.Printf("failed to create zip file: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	uz, err := zippy.NewUnzippy(zipPath, nil)
	if err != nil {
		fmt.Printf("failed to create unzippy: %v\n", err)
		return
	}

	outDir := filepath.Join(tempDir, "extracted")

	zFiles, err := uz.ExtractTo(outDir)
	if err != nil {
		fmt.Printf("failed to extract files: %v\n", err)
		return
	}

	displayExtractedFiles(outDir, zFiles)
}
