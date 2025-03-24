package testutils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
	// Step 1: Create a temporary directory to hold the files
	tempDir, err := os.MkdirTemp("", "zip-temp-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory

	// Step 2: Create files and subdirectories in the temporary directory
	for i := 0; i < files; i++ {
		filePath := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		if err := os.WriteFile(filePath, []byte(fmt.Sprintf("Test File %d", i)), 0644); err != nil {
			return err
		}
	}

	for i := 0; i < subdirs; i++ {
		subfolderPath := filepath.Join(tempDir, fmt.Sprintf("subfolder%d", i))
		if err := os.MkdirAll(subfolderPath, os.ModePerm); err != nil {
			return err
		}

		for j := 0; j < files; j++ {
			filePath := filepath.Join(subfolderPath, fmt.Sprintf("test%d.txt", j))
			if err := os.WriteFile(filePath, []byte(fmt.Sprintf("Test File %d", j)), 0644); err != nil {
				return err
			}
		}
	}

	// Step 3: Create the zip file
	zFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zFile.Close()

	zWrite := zip.NewWriter(zFile)
	defer zWrite.Close()

	// Step 4: Walk through the temporary directory and add files to the zip archive
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == tempDir {
			return nil
		}

		// Get the relative path to maintain the directory structure in the zip archive
		relPath, err := filepath.Rel(tempDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Add a directory entry to the zip archive
			_, err := zWrite.Create(relPath + "/")
			return err
		}

		// Add a file entry to the zip archive
		zipFile, err := zWrite.Create(relPath)
		if err != nil {
			return err
		}

		// Open the file and copy its contents to the zip archive
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(zipFile, srcFile)
		return err
	})

	return err
}

// PermissionTest is a helper function to wrap another function that requires a file to have specific permissions.
//
// filePermPath is the path to the file to change permissions on.
//
// funcToRun is the function to run that requires the file to have specific permissions.
//
// args are the arguments to pass to the function.
func PermissionTest(filePermPath string, funcToRun interface{}, args ...interface{}) error {
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
			if restoreErr := os.Chmod(filePermPath, 0755); restoreErr != nil {
				err = restoreErr
			}
		}()
	}

	// Use reflection to call the target function with the provided arguments
	funcValue := reflect.ValueOf(funcToRun)
	funcArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		funcArgs[i] = reflect.ValueOf(arg)
	}

	results := funcValue.Call(funcArgs)

	// Check if the function returned an error
	if len(results) > 0 && !results[len(results)-1].IsNil() {
		err = results[len(results)-1].Interface().(error)
	}

	return err
}
