//go:build linux || darwin
// +build linux darwin

package zippy

import (
	"os"

	"golang.org/x/sys/unix"
)

// isCrossDevice checks if two paths are on different devices on Unix-like systems.
func isCrossDevice(path1, path2 string) (bool, error) {
	info1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}

	info2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	stat1 := info1.Sys().(*unix.Stat_t)
	stat2 := info2.Sys().(*unix.Stat_t)
	return stat1.Dev != stat2.Dev, nil
}
