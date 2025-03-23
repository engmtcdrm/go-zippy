package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pp "github.com/engmtcdrm/go-prettyprint"
	"github.com/engmtcdrm/zippy-tmp"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, "contents.zip")

	zFiles, err := zippy.Contents(path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Contents of zip file (%s):\n\n", pp.Green(path))

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
			pp.Green(zFile.Name),
			strings.Repeat(" ", pad-len(zFile.Name)),
			zFile.Modified.Format("2006-01-02 15:04:05"),
		)
		fileCnt++
	}

	fmt.Println("    -------")
	fmt.Printf("    Total files: %s\n", pp.Green(fileCnt))
	fmt.Println()
}
