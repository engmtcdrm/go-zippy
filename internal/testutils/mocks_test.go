package testutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for [NewMockReader] function.
func Test_NewMockReader(t *testing.T) {
	t.Run("create valid mock reader", func(t *testing.T) {
		reader := strings.NewReader("your string here")
		mockReader := NewMockReader(reader)
		require.NotNil(t, mockReader)
	})

	t.Run("create nil mock reader", func(t *testing.T) {
		mockReader := NewMockReader(nil)
		require.NotNil(t, mockReader)
	})

}
