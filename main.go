package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	if err := zippy.Add(zipFile, tempFile.Name(), filepath.Join(tempDir, "bubba", "*")); err != nil {
		return err
	}

	if err := zippy.Delete(zipFile, filepath.Join(tempDir, "subfolder0", "*"), "bubba/bubba*.txt"); err != nil {
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

	testFiles := []string{
		"t1",
		"t2",
		"testdir/",
		"testdir/tt2",
		"testdir/tt1",
		"t1.txt",
		"t2.txt",
		"testdir/t2.txt",
		"testdir/t1.txt",
		"bubba/bubba/",
		"bubba-123456.txt",
		"bubba/",
		"bubba/bubba/bubba-19023850.txt",
		"bubba/bubba-1204023.txt",
		"testdir2/",
		"testdir2/bubba/",
		"testdir2/bubba/bubba-190238591835.txt",
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter some text: ")
	matchString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	matchString = strings.ReplaceAll(matchString, "\r\n", "")

	for _, testFile := range testFiles {
		match, err := filepath.Match(matchString, testFile)
		if err != nil {
			panic(err)
		}

		if match {
			fmt.Printf("Matched '%s' with '%s'\n", pp.Green(testFile), pp.Green(match))
		} else {
			fmt.Printf("Did not match '%s' with '%s'\n", pp.Red(testFile), pp.Red(match))
		}
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

	_, err = zippy.UnzipTo(absZipFile, absZip+"2")
	if err != nil {
		panic(err)
	}

	// Test relative path
	if err := testRelPath(relZipFile); err != nil {
		panic(err)
	}

	if err := testContents(relZipFile); err != nil {
		panic(err)
	}
}
