//go:build windows
// +build windows

package zippy

import (
	"os"
	"path/filepath"
)

// getVolumeName retrieves the volume name for the given path.
func getVolumeName(path string) (string, error) {
	volumeName := filepath.VolumeName(path)
	if volumeName == "" {
		return "", os.ErrInvalid
	}
	return volumeName, nil
}

// isCrossDevice checks if two paths are on different devices on Windows.
func isCrossDevice(path1, path2 string) (bool, error) {
	volume1, err := getVolumeName(path1)
	if err != nil {
		return false, err
	}

	volume2, err := getVolumeName(path2)
	if err != nil {
		return false, err
	}

	return volume1 != volume2, nil
}
