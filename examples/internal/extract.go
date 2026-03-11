package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			filepath.Join(tempDir, file.Name),
		)
	}

	fmt.Println()
}
