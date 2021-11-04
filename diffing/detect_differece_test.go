package diffing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectDifferences(t *testing.T) {

	t.Run("all nil", func(t *testing.T) {
		result := DetectDifferences(nil, nil)
		assert.Equal(t, empty(), result.Removed)
		assert.Equal(t, empty(), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})

	t.Run("before nil", func(t *testing.T) {
		result := DetectDifferences(nil, slice("a", "b"))
		assert.Equal(t, empty(), result.Removed)
		assert.Equal(t, slice("a", "b"), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})

	t.Run("after nil", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b"), nil)
		assert.Equal(t, slice("a", "b"), result.Removed)
		assert.Equal(t, empty(), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})

	t.Run("all empty", func(t *testing.T) {
		result := DetectDifferences(empty(), empty())
		assert.Equal(t, empty(), result.Removed)
		assert.Equal(t, empty(), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})

	t.Run("all same", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b"), slice("a", "b"))
		assert.Equal(t, empty(), result.Removed)
		assert.Equal(t, empty(), result.Added)
		assert.Equal(t, slice("a", "b"), result.Unchanged)
	})

	t.Run("one removed", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b", "c"), slice("a", "b"))
		assert.Equal(t, slice("c"), result.Removed)
		assert.Equal(t, empty(), result.Added)
		assert.Equal(t, slice("a", "b"), result.Unchanged)
	})

	t.Run("all removed", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b"), empty())
		assert.Equal(t, slice("a", "b"), result.Removed)
		assert.Equal(t, empty(), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})

	t.Run("one added", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b"), slice("a", "b", "c"))
		assert.Equal(t, empty(), result.Removed)
		assert.Equal(t, slice("c"), result.Added)
		assert.Equal(t, slice("a", "b"), result.Unchanged)
	})

	t.Run("all added", func(t *testing.T) {
		result := DetectDifferences(empty(), slice("a", "b"))
		assert.Equal(t, empty(), result.Removed)
		assert.Equal(t, slice("a", "b"), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})

	t.Run("one added, one removed", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b"), slice("b", "c"))
		assert.Equal(t, slice("a"), result.Removed)
		assert.Equal(t, slice("c"), result.Added)
		assert.Equal(t, slice("b"), result.Unchanged)
	})

	t.Run("all added, all removed", func(t *testing.T) {
		result := DetectDifferences(slice("a", "b"), slice("c", "d"))
		assert.Equal(t, slice("a", "b"), result.Removed)
		assert.Equal(t, slice("c", "d"), result.Added)
		assert.Equal(t, empty(), result.Unchanged)
	})
}

func slice(varargs ...string) []interface{} {
	result := []interface{}{}
	for _, arg := range varargs {
		result = append(result, arg)
	}
	return result
}

func empty() []interface{} {
	return []interface{}{}
}
