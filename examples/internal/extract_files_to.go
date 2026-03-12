package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/engmtcdrm/go-zippy"
)

func ExtractFilesTo() {
	tempDir, zipPath, err := createZipFile(10, 1)
	if err != nil {
		fmt.Printf("failed to create zip file: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	contents, err := zippy.Contents(zipPath)
	if err != nil {
		fmt.Printf("failed to get zip contents: %v\n", err)
		return
	}

	uz, err := zippy.NewUnzippy(zipPath, nil)
	if err != nil {
		fmt.Printf("failed to create unzippy: %v\n", err)
		return
	}

	outDir := filepath.Join(tempDir, "extracted")

	zFiles, err := uz.ExtractFilesTo(outDir, contents[1].Name, contents[12].Name, contents[15].Name)
	if err != nil {
		fmt.Printf("failed to extract files: %v\n", err)
		return
	}

	displayExtractedFiles(outDir, zFiles)
}
