package testutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for [NewMockReader] function.
func TestNewMockReader(t *testing.T) {
	t.Run("create valid mock reader", func(t *testing.T) {
		reader := strings.NewReader("your string here")
		mockReader := NewMockReader(reader)
		assert.NotNil(t, mockReader)
	})

	t.Run("create nil mock reader", func(t *testing.T) {
		mockReader := NewMockReader(nil)
		assert.NotNil(t, mockReader)
	})

}
