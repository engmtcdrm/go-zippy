package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/go-zippy"
)

func promptCleanupFiles(cwd string, zFiles []*zip.File) {
	choice := "Y"

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to cleanup the extracted file? (Y/n): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	if text != "" {
		choice = text
	}

	fmt.Println()

	if choice == "y" || choice == "Y" {
		for _, file := range zFiles {
			if err := os.RemoveAll(filepath.Join(cwd, file.Name)); err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, "extract.zip")

	uz := zippy.NewUnzippy(path, nil)

	zFiles, err := uz.Extract()
	if err != nil {
		panic(err)
	}

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
			filepath.Join(cwd, file.Name),
		)
	}

	fmt.Println()

	promptCleanupFiles(cwd, zFiles)
}
