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
	copy_path := filepath.Join(cwd, fmt.Sprintf("copy_of_%s.zip", filepath.Base(cwd)))

	z := zippy.NewZippy(path)

	if err := z.Copy(copy_path); err != nil {
		panic(err)
	}
	fmt.Println("Copy of zip file created successfully.")
	fmt.Println()

	promptCleanupFiles(cwd, copy_path)
}
