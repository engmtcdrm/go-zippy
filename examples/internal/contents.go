package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/engmtcdrm/go-zippy"
)

func ExampleContents() {
	tempDir, zipPath, err := createZipFile(10, 1)
	if err != nil {
		fmt.Printf("failed to create zip file: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	zFiles, err := zippy.Contents(zipPath)
	if err != nil {
		fmt.Printf("failed to get contents of zip file: %v\n", err)
		return
	}

	fmt.Printf("\nContents of zip file (%s):\n\n", zipPath)

	var fileCnt = 0

	pad := 0

	for _, zFile := range zFiles {
		if len(zFile.Name) > pad {
			pad = len(zFile.Name)
		}
	}

	for _, zFile := range zFiles {
		// Use the pad variable to dynamically set the width for the filename
		fmt.Printf("    %s%s\t%s\n",
			zFile.Name,
			strings.Repeat(" ", pad-len(zFile.Name)),
			zFile.Modified.Format("2006-01-02 15:04:05"),
		)
		fileCnt++
	}

	fmt.Println("    -------")
	fmt.Printf("    Total files: %d\n", fileCnt)
	fmt.Println()
}
