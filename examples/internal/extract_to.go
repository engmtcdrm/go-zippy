package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	fmt.Println()
	fmt.Println("Extracted files:")
	fmt.Println()

	pad := 0

	for _, zFile := range zFiles {
		if len(zFile.Name) > pad {
			pad = len(zFile.Name)
		}
	}

	for _, file := range zFiles {
		fmt.Printf(
			"Archive path:\t%s%s\tExtracted to:\t%s\n",
			file.Name,
			strings.Repeat(" ", pad-len(file.Name)),
			filepath.Join(outDir, file.Name),
		)
	}

	fmt.Println()
}
