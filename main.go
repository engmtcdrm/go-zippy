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

var (
	absZip *zippy.Zippy
	relZip *zippy.Zippy
)

func testAbsZip() error {
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

	if err := absZip.Add(tempDir); err != nil {
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

	if err := absZip.Add(tempFile.Name(), filepath.Join(tempDir, "bubba", "*")); err != nil {
		return err
	}

	if err := absZip.Delete(filepath.Join(tempDir, "subfolder0", "*"), filepath.Join(tempDir, "bubba/bubba*.txt")); err != nil {
		return err
	}

	// if err := zippy.Zip(zipFile, tempDir); err != nil {
	// 	return err
	// }

	// if err := zippy.Add(zipFile, tempFile.Name(), filepath.Join(tempDir, "bubba", "*")); err != nil {
	// 	return err
	// }

	// if err := zippy.Delete(zipFile, filepath.Join(tempDir, "subfolder0", "*"), "bubba/bubba*.txt"); err != nil {
	// 	return err
	// }

	return err
}

func testContents(z *zippy.Zippy) error {
	zFiles, err := zippy.Contents(z.Path)
	if err != nil {
		return err
	}

	fmt.Printf("Contents of zip file (%s):\n\n", pp.Green(z.Path))

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

	z := zippy.NewZippy(zipFile)

	if err := z.Add(relDir); err != nil {
		return err
	}

	// if err := zippy.Zip(zipFile, relDir); err != nil {
	// 	return err
	// }

	return err
}

func main() {
	baseDir := flag.String("d", ".test", "Base directory for zip operations")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	absBasePath := filepath.Join(*baseDir, "abs-zip")
	relBasePath := filepath.Join(*baseDir, "rel-zip")
	absZip = zippy.NewZippy(filepath.Join(absBasePath, "test.zip"))
	relZip = zippy.NewZippy(filepath.Join(relBasePath, "test.zip"))

	fmt.Println("Current working directory:", pp.Green(cwd))
	fmt.Println("Base directory:", pp.Green(*baseDir))
	fmt.Println("Absolute path zip directory:", pp.Green(absBasePath))
	fmt.Println("Relative path zip directory:", pp.Green(relBasePath))
	fmt.Println("Absolute path zip file:", pp.Green(absZip.Path))
	fmt.Println("Relative path zip file:", pp.Green(relZip.Path))
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
	if err := testAbsZip(); err != nil {
		panic(err)
	}

	if err := testContents(absZip); err != nil {
		panic(err)
	}

	uzippy := zippy.NewUnzippy(absZip.Path)

	// Test unzipping absolute path zip file
	_, err = uzippy.Extract()
	if err != nil {
		panic(err)
	}

	_, err = uzippy.ExtractTo(absBasePath + "2")
	if err != nil {
		panic(err)
	}

	uzippy.Junk = true

	_, err = uzippy.ExtractFilesTo(absBasePath+"3", "*test0*.txt")
	if err != nil {
		panic(err)
	}

	// Test relative path
	if err := testRelPath(relZip.Path); err != nil {
		panic(err)
	}

	if err := testContents(relZip); err != nil {
		panic(err)
	}
}
