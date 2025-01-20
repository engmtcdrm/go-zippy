package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/m/zippy"
)

func createTempFile(dir, name string) error {
	tempFile, err := os.CreateTemp(dir, name)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = tempFile.Write([]byte(filepath.Base(tempFile.Name())))
	if err != nil {
		panic(err)
	}

	return nil
}

func createTestFiles(dir string, files int) error {
	for i := 0; i < files; i++ {
		if err := createTempFile(dir, fmt.Sprintf("test%d-*.txt", i)); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	tempDir, err := os.MkdirTemp("", "test-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	if err := createTestFiles(tempDir, 10); err != nil {
		panic(err)
	}

	if err := zippy.Zip(tempDir, "test/test.zip"); err != nil {
		panic(err)
	}

	if err := zippy.Unzip("test/test.zip", "test2"); err != nil {
		panic(err)
	}

	// err = zippy.Zip("test2", "test3/test2-bak.zip")
	// if err != nil {
	// 	panic(err)
	// }

	// err = zippy.Zip("test/test1", "test3/testfile1.zip")
	// if err != nil {
	// 	panic(err)
	// }

	// err = zippy.Zip("test/test2", "test3/testfile2.zip")
	// if err != nil {
	// 	panic(err)
	// }

	// err = zippy.Zip("test", "test3/test-dir.zip")
	// if err != nil {
	// 	panic(err)
	// }
}
