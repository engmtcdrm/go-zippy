package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/engmtcdrm/go-zippy"
)

func promptCleanupFiles(cwd string, copyPath string) {
	choice := "Y"

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to cleanup the copied file? (Y/n): ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	if text != "" {
		choice = text
	}

	fmt.Println()

	if choice == "y" || choice == "Y" {
		if err := os.RemoveAll(copyPath); err != nil {
			panic(err)
		}
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(cwd, fmt.Sprintf("%s.zip", filepath.Base(cwd)))
	copyPath := filepath.Join(cwd, fmt.Sprintf("copy_of_%s.zip", filepath.Base(cwd)))

	z := zippy.NewZippy(path)

	if err := z.Copy(copyPath, "f1", "d2/", "d2/d1/d2/f1"); err != nil {
		panic(err)
	}

	fmt.Printf("Copied files found in '%s':\n", copyPath)
	fmt.Println()

	zFiles, err := zippy.Contents(copyPath)
	if err != nil {
		panic(err)
	}

	pad := 0

	for _, zFile := range zFiles {
		if len(zFile.Name) > pad {
			pad = len(zFile.Name)
		}
	}

	for _, file := range zFiles {
		fmt.Printf(
			"Archive Path:\t%s\n",
			file.Name,
		)
	}

	fmt.Println()

	promptCleanupFiles(cwd, copyPath)
}
