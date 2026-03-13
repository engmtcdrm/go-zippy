package testutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [PermissionTest] function.
func TestPermissionTest(t *testing.T) {
	t.Run("0-arg function", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile, err := CreateTempFile(tempDir, "test-*.txt")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)

		err = PermissionTest(tempFile.Name(), func() error { return nil })
		assert.NoError(t, err)
	})

	t.Run("1-arg function", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile, err := CreateTempFile(tempDir, "test-*.txt")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)

		err = PermissionTest(tempFile.Name(), func(arg1 string) error { return nil }, tempFile.Name())
		assert.NoError(t, err)
	})

	t.Run("2-arg function", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile, err := CreateTempFile(tempDir, "test-*.txt")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)

		err = PermissionTest(tempFile.Name(), func(arg1, arg2 string) error { return nil }, tempFile.Name(), "arg2")
		assert.NoError(t, err)
	})

	t.Run("variadic function", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile, err := CreateTempFile(tempDir, "test-*.txt")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)

		err = PermissionTest(tempFile.Name(), func(args ...string) error {
			if len(args) != 2 {
				return errors.New("expected 2 arguments")
			}
			return nil
		}, tempFile.Name(), "arg2")
		assert.NoError(t, err)
	})

	t.Run("function with error", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile, err := CreateTempFile(tempDir, "test-*.txt")
		assert.NoError(t, err)
		assert.NotNil(t, tempFile)

		err = PermissionTest(tempFile.Name(), func(args ...string) error {
			return errors.New("test error")
		}, tempFile.Name(), "arg2")
		assert.Error(t, err)
	})
}
