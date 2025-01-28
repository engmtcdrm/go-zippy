package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/m/zippy"
	"example.com/m/zippy/testutils"
)

func main() {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	baseDir := ".test"
	tgtDir := filepath.Join(baseDir, "unzip")
	zipFile := filepath.Join(baseDir, "zip", "test.zip")

	os.RemoveAll(".test")

	if _, err := testutils.CreateTestFiles(tempDir, 3, 2); err != nil {
		panic(err)
	}

	if err := zippy.Zip(zipFile, tempDir); err != nil {
		panic(err)
	}

	tempFile, err := testutils.CreateTempFile(tempDir, "bubba-*.txt")
	if err != nil {
		panic(err)
	}

	if err := zippy.AddToZip(zipFile, tempFile.Name()); err != nil {
		panic(err)
	}

	if err := zippy.RemoveFromZip(zipFile, filepath.Base(tempFile.Name()), "subfolder0/"); err != nil {
		panic(err)
	}

	zFiles, err := zippy.Contents(zipFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Contents of zip file (%s):\n\n", zipFile)

	for _, zFile := range zFiles {
		fmt.Printf("    %s\n", zFile.Name)
	}

	if err := zippy.Unzip(zipFile, tgtDir); err != nil {
		panic(err)
	}
}
