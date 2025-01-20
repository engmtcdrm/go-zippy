package zippy

import (
	"fmt"
	"path/filepath"
)

func validateCopy(filePath string, written int64, expected int64) error {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	if written != expected {
		return fmt.Errorf("failed to copy '%s': expected %d bytes, got %d bytes", absFilePath, expected, written)
	}

	return nil
}
