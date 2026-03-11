package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
)

// PermissionTest is a helper function to wrap another function that requires a
// file to have specific permissions.
func PermissionTest(filePermPath string, fn interface{}, args ...interface{}) error {
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
		defer os.Chmod(filePermPath, 0755)
	}

	// Use reflection to call the target function with the provided arguments
	funcValue := reflect.ValueOf(fn)
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
