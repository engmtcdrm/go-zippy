package testutils

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

func CreateTempFile(dir, name string) error {
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

func CreateTestFiles(dir string, files int, subdirs int) error {
	for i := 0; i < files; i++ {
		if err := CreateTempFile(dir, fmt.Sprintf("test%d-*.txt", i)); err != nil {
			return err
		}
	}

	for i := 0; i < subdirs; i++ {
		subfolder := filepath.Join(dir, fmt.Sprintf("subfolder%d", i))

		if err := os.Mkdir(subfolder, os.ModePerm); err != nil {
			return err
		}

		for j := 0; j < files; j++ {
			if err := CreateTempFile(subfolder, fmt.Sprintf("test%d-*.txt", j)); err != nil {
				return err
			}
		}
	}

	return nil
}

func CreateZipFile(zipFilePath string, files int, subdirs int) error {
	zFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zFile.Close()

	zWrite := zip.NewWriter(zFile)
	defer zWrite.Close()

	for i := 0; i < files; i++ {
		fWrite, err := zWrite.Create(fmt.Sprintf("test%d.txt", i))
		if err != nil {
			return err
		}

		_, err = fWrite.Write([]byte(fmt.Sprintf("Test File %d", i)))
		if err != nil {
			return err
		}
	}

	for i := 0; i < subdirs; i++ {
		subfolder := fmt.Sprintf("subfolder%d", i)

		for j := 0; j < files; j++ {
			fileWriter, err := zWrite.Create(filepath.Join(subfolder, fmt.Sprintf("test%d.txt", j)))
			if err != nil {
				return err
			}

			_, err = fileWriter.Write([]byte(fmt.Sprintf("Test File %d", j)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
