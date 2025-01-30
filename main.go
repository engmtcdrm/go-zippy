package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"example.com/m/zippy"
	"example.com/m/zippy/testutils"
	pp "github.com/engmtcdrm/go-prettyprint"
)

func testAbsZip(zipFile string) error {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		return err
	}
	defer func() {
		if removeError := os.RemoveAll(tempDir); removeError != nil {
			err = fmt.Errorf("failed to remove temp dir: %w", removeError)
		}
	}()

	if _, err := testutils.CreateTestFiles(tempDir, 3, 2); err != nil {
		return err
	}

	if err := zippy.Zip(zipFile, tempDir); err != nil {
		return err
	}

	tempFile, err := testutils.CreateTempFile(tempDir, "bubba-*.txt")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(tempDir, "bubba"), os.ModePerm); err != nil {
		return err
	}

	_, err = testutils.CreateTempFile(filepath.Join(tempDir, "bubba"), "bubba2-*.txt")
	if err != nil {
		return err
	}

	if err := zippy.Add(zipFile, tempFile.Name(), filepath.Join(tempDir, "bubba")); err != nil {
		return err
	}

	if err := zippy.Delete(zipFile, filepath.Join(tempDir, "subfolder0"), filepath.Join(tempDir, "bubba")); err != nil {
		return err
	}

	return err
}

func testContents(zipFile string) error {
	zFiles, err := zippy.Contents(zipFile)
	if err != nil {
		return err
	}

	fmt.Printf("Contents of zip file (%s):\n\n", pp.Green(zipFile))

	var fileCnt = 0

	for _, zFile := range zFiles {
		fmt.Printf("    %s\n", pp.Green(zFile.Name))
		fileCnt++
	}

	fmt.Println("    -------")
	fmt.Printf("    Total files: %s\n", pp.Green(fileCnt))
	fmt.Println()

	return nil
}

func testRelPath(zipFile string) error {
	var err error
	zipFileDir := filepath.Dir(zipFile)
	relDir := filepath.Join(zipFileDir, "rel")

	os.Chdir(zipFileDir)

	if err := os.MkdirAll(relDir, os.ModePerm); err != nil {
		return err
	}

	if _, err := testutils.CreateTestFiles(relDir, 3, 2); err != nil {
		return err
	}

	if err := zippy.Zip(zipFile, relDir); err != nil {
		return err
	}

	return err
}

func main() {
	baseDir := flag.String("d", ".test", "Base directory for zip operations")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	absZip := filepath.Join(*baseDir, "abs-zip")
	relZip := filepath.Join(*baseDir, "rel-zip")
	absZipFile := filepath.Join(absZip, "test.zip")
	relZipFile := filepath.Join(relZip, "test.zip")

	fmt.Println("Current working directory:", pp.Green(cwd))
	fmt.Println("Base directory:", pp.Green(*baseDir))
	fmt.Println("Absolute path zip directory:", pp.Green(absZip))
	fmt.Println("Relative path zip directory:", pp.Green(relZip))
	fmt.Println("Absolute path zip file:", pp.Green(absZipFile))
	fmt.Println("Relative path zip file:", pp.Green(relZipFile))
	fmt.Println()

	if err := os.RemoveAll(*baseDir); err != nil {
		panic(err)
	}

	// Test absolute path
	if err := testAbsZip(absZipFile); err != nil {
		panic(err)
	}

	if err := testContents(absZipFile); err != nil {
		panic(err)
	}

	// Test unzipping absolute path zip file
	_, err = zippy.Unzip(absZipFile)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Unzipped files:\n\n")

	// for _, unzippedFile := range unzippedFiles {

	// 	fmt.Printf("    %s [%s]\n", pp.Green(unzippedFile.Name), filepath.Join(filepath.VolumeName(cwd), filepath.FromSlash(unzippedFile.Name)))
	// }

	_, err = zippy.UnzipTo(absZipFile, absZip+"2")
	if err != nil {
		panic(err)
	}

	// Test relative path
	// if err := testRelPath(relZipFile); err != nil {
	// 	panic(err)
	// }

	// if err := testContents(relZipFile); err != nil {
	// 	panic(err)
	// }
}
