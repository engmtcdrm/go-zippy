package testutils

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CreateTempFile creates a temporary file in the specified directory with the specified name.
// The file is created with a random name and the content of the file is the name of the file.
//
// dir is the directory where the file will be created.
//
// name is the name of the file.
func CreateTempFile(dir, name string) (*os.File, error) {
	tempFile, err := os.CreateTemp(dir, name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := tempFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close zipped file: %w", closeErr)
		}
	}()

	_, err = tempFile.Write([]byte(filepath.Base(tempFile.Name())))
	if err != nil {
		return nil, err
	}

	return tempFile, err
}

// CreateTestFiles creates a specified number of files and subdirectories in the specified directory.
//
// dir is the directory where the files and subdirectories will be created.
//
// files is the number of files to create.
//
// subdirs is the number of subdirectories to create. The subdirectories will contain the same number
// of files as the parent directory.
func CreateTestFiles(dir string, files int, subdirs int) ([]*os.File, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	var tempFiles []*os.File

	for i := 0; i < files; i++ {
		tempFile, err := CreateTempFile(dir, fmt.Sprintf("test%d-*.txt", i))
		if err != nil {
			return nil, err
		}

		tempFiles = append(tempFiles, tempFile)
	}

	for i := 0; i < subdirs; i++ {
		subfolder := filepath.Join(dir, fmt.Sprintf("subfolder%d", i))

		if err := os.MkdirAll(subfolder, os.ModePerm); err != nil {
			return nil, err
		}

		for j := 0; j < files; j++ {
			tempFile, err := CreateTempFile(subfolder, fmt.Sprintf("test%d-*.txt", j))
			if err != nil {
				return nil, err
			}

			tempFiles = append(tempFiles, tempFile)
		}
	}

	return tempFiles, nil
}

// CreateZipFile creates a zip file with the specified number of files and subdirectories.
//
// zipFilePath is the path to the zip file to create.
//
// files is the number of files to create.
//
// subdirs is the number of subdirectories to create. The subdirectories will contain the same number
// of files as the parent directory.
func CreateZipFile(zipFilePath string, files int, subdirs int) error {
	zFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := zFile.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	zWrite := zip.NewWriter(zFile)
	defer func() {
		if closeErr := zWrite.Close(); closeErr != nil {
			err = closeErr
		}
	}()

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
		subfolder := fmt.Sprintf("subfolder%d/", i)

		_, err := zWrite.Create(subfolder)
		if err != nil {
			return nil
		}

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

	return err
}

// PermissionTest1Arg is a helper function to wrap another function that requires a file to have specific permissions.
//
// filePermPath is the path to the file to change permissions on.
//
// funcToRun is the function to run that requires the file to have specific permissions.
//
// arg1 is the first argument to pass to the function.
func PermissionTest1Arg(filePermPath string, funcToRun func(string) ([]*zip.File, error), arg1 string) error {
	var err error

	if runtime.GOOS == "windows" {
		cmd := exec.Command("icacls", filePermPath, "/deny", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
		if err := cmd.Run(); err != nil {
			return err
		}
		defer func() {
			// Restore permissions after the test
			cmd := exec.Command("icacls", filePermPath, "/grant", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
			if runError := cmd.Run(); runError != nil {
				err = runError
			}
		}()
	} else {
		if err := os.Chmod(filePermPath, 0000); err != nil {
			return err
		}
		defer func() {
			// Restore permissions after the test
			if restoreErr := os.Chmod(filePermPath, 0755); err != nil {
				err = restoreErr
			}
		}()
	}

	_, err = funcToRun(arg1)

	return err
}

// PermissionTest2Args is a helper function to wrap another function that requires a file to have specific permissions.
//
// filePermPath is the path to the file to change permissions on.
//
// funcToRun is the function to run that requires the file to have specific permissions.
//
// arg1 is the first argument to pass to the function.
//
// arg2 is the second argument to pass to the function.
func PermissionTest2Args(filePermPath string, funcToRun func(string, string) error, arg1 string, arg2 string) error {
	var err error

	if runtime.GOOS == "windows" {
		cmd := exec.Command("icacls", filePermPath, "/deny", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
		if err := cmd.Run(); err != nil {
			return err
		}
		defer func() {
			// Restore permissions after the test
			cmd := exec.Command("icacls", filePermPath, "/grant", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
			if runError := cmd.Run(); runError != nil {
				err = runError
			}
		}()
	} else {
		if err := os.Chmod(filePermPath, 0000); err != nil {
			return err
		}
		defer func() {
			// Restore permissions after the test
			if restoreErr := os.Chmod(filePermPath, 0755); err != nil {
				err = restoreErr
			}
		}()
	}

	err = funcToRun(arg1, arg2)

	return err
}

// PermissionTestVariadic is a helper function to wrap another function that requires a file to have specific permissions.
// This version supports a variadic argument for funcToRun.
//
// filePermPath is the path to the file to change permissions on.
//
// funcToRun is the function to run that requires the file to have specific permissions.
//
// arg1 is the first argument to pass to the function.
//
// args are the variadic arguments to pass to the function.
func PermissionTestVariadic(filePermPath string, funcToRun func(string, ...string) error, arg1 string, args ...string) error {
	var err error

	if runtime.GOOS == "windows" {
		cmd := exec.Command("icacls", filePermPath, "/deny", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
		if err := cmd.Run(); err != nil {
			return err
		}
		defer func() {
			// Restore permissions after the test
			cmd := exec.Command("icacls", filePermPath, "/grant", fmt.Sprintf("%s:F", os.Getenv("USERNAME")))
			if runError := cmd.Run(); runError != nil {
				err = runError
			}
		}()
	} else {
		if err := os.Chmod(filePermPath, 0000); err != nil {
			return err
		}
		defer func() {
			// Restore permissions after the test
			if restoreErr := os.Chmod(filePermPath, 0755); err != nil {
				err = restoreErr
			}
		}()
	}

	err = funcToRun(arg1, args...)

	return err
}
